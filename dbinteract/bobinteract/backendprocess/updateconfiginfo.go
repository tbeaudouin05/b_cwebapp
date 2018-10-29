package backendprocess

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	mongointeract "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/backendprocess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	bmlconfigbackendprocess "github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig/backendprocess"
	"github.com/thomas-bamilo/nosql/mongobulk"
)

func UpdateConfigInfoAndHist(dbBob *sql.DB, mongoSession *mgo.Session) {

	stmt, err := dbBob.Prepare(`
	-- it is fine if several SKUs have the same name, keep structure similar to BOB - plus it will be much clearer for business users
	SELECT
-- qualitative information
	  cc.id_catalog_config
	  ,cc.sku
	  ,COALESCE(cc.name,'') sku_name
	  ,COALESCE(CONCAT('https://media.bamilo.com/p/',cb.url_key,'-',RIGHT(UNIX_TIMESTAMP(cpi.updated_at),4),'-',REVERSE(cs.fk_catalog_config),'-1-product.jpg'),'') img_link
	  ,COALESCE(CONCAT('https://www.bamilo.com/',LOWER(cb.name),'-',REPLACE(LOWER(cc.name),' ','-'),'-',cc.id_catalog_config,'.html' ),'') sku_link
	  ,COALESCE(cc.description,'') description
	  ,COALESCE(cc.short_description,'') short_description
	  ,COALESCE(cc.package_content,'') package_content
	  ,COALESCE(cc.product_warranty,'') product_warranty
-- category
	  ,COALESCE(bi_one.name,'') bi_category_one_name
	  ,COALESCE(bi_two.name,'') bi_category_two_name
	  ,COALESCE(bi_three.name,'') bi_category_three_name
-- brand
	  ,COALESCE(cb.name,'') brand_name
	  ,COALESCE(cb.name_en,'') brand_name_en
-- supplier
	  ,COALESCE(s.name,'') supplier_name
	  ,COALESCE(s.name_en,'') supplier_name_en
-- historical data
  ,COALESCE(vccv.visible_in_shop,0)
	  ,COALESCE(FLOOR(AVG(cs.price)),0) avg_price
	  ,COALESCE(FLOOR(AVG(cs.special_price)),0) avg_special_price
	  ,COALESCE(FLOOR(SUM(cs2.quantity)),0) sum_of_stock_quantity
	  ,COALESCE(FLOOR(MIN(cs2.quantity)),0) min_of_stock_quantity

	  FROM catalog_simple cs
	  JOIN catalog_config cc
	  ON cs.fk_catalog_config = cc.id_catalog_config
	  JOIN catalog_brand cb
	  ON cc.fk_catalog_brand = cb.id_catalog_brand
	  JOIN catalog_product_image cpi
	  ON cpi.fk_catalog_config = cc.id_catalog_config
  JOIN bob_live_ir.view_catalog_config_visibility vccv
	  ON vccv.sku = cc.sku
  LEFT JOIN catalog_category_bi bi_one
	  ON cc.bi_category_one = bi_one.id_catalog_category_bi
	  LEFT JOIN catalog_category_bi bi_two
	  ON cc.bi_category_two = bi_two.id_catalog_category_bi
	  LEFT JOIN catalog_category_bi bi_three
	  ON cc.bi_category_three = bi_three.id_catalog_category_bi -- some products do not have bi_category_three, just to make sure, also left join other bi_category
	  LEFT JOIN catalog_source cs1
	  ON cs.id_catalog_simple = cs1.fk_catalog_simple
	  LEFT JOIN catalog_stock cs2
	  ON cs1.id_catalog_source = cs2.fk_catalog_source
  LEFT JOIN supplier s
	  ON cs1.fk_supplier = s.id_supplier
  LEFT JOIN (
  SELECT 
   soi.sku
   ,soi.created_at created_at -- not possible to refresh sales here, should be done with sales script
  FROM sales_order_item soi
  WHERE soi.created_at >= NOW()-INTERVAL 90 DAY -- 90
  GROUP BY soi.sku, soi.created_at
  ) sales_order_sku
  ON sales_order_sku.sku = cs.sku
 
  WHERE (sales_order_sku.created_at >= NOW()-INTERVAL 90 DAY -- 90
  OR cc.created_at >= NOW()-INTERVAL 15 DAY -- 15
  OR vccv.visible_in_shop = 1)

	  GROUP BY cc.id_catalog_config;
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	bamiloCatalogConfig := bmlcatalogconfig.BmlCatalogConfig{}

	cHist := mongoSession.DB("competition_analysis").C("bml_catalog_config_hist")
	c := mongoSession.DB("competition_analysis").C("bml_catalog_config")
	config := mongobulk.Config{OpsPerBatch: 950}
	mongoBulkHist := mongobulk.New(cHist, config)
	mongoBulk := mongobulk.New(c, config)

	for rows.Next() {

		err := rows.Scan(
			// qualitative information
			&bamiloCatalogConfig.IDBmlCatalogConfig,
			&bamiloCatalogConfig.SKU,
			&bamiloCatalogConfig.SKUName,
			&bamiloCatalogConfig.ImgLink,
			&bamiloCatalogConfig.SKULink,
			&bamiloCatalogConfig.Description,
			&bamiloCatalogConfig.ShortDescription,
			&bamiloCatalogConfig.PackageContent,
			&bamiloCatalogConfig.ProductWarranty,
			// category
			&bamiloCatalogConfig.BiCategoryOneName,
			&bamiloCatalogConfig.BiCategoryTwoName,
			&bamiloCatalogConfig.BiCategoryThreeName,
			// brand
			&bamiloCatalogConfig.BrandName,
			&bamiloCatalogConfig.BrandNameEn,
			// supplier
			&bamiloCatalogConfig.SupplierName,
			&bamiloCatalogConfig.SupplierNameEn,
			// historical data
			&bamiloCatalogConfig.VisibleInShop,
			&bamiloCatalogConfig.AvgPrice,
			&bamiloCatalogConfig.AvgSpecialPrice,
			&bamiloCatalogConfig.SumOfStockQuantity,
			&bamiloCatalogConfig.MinOfStockQuantity,
		)
		checkError(err)

		bamiloCatalogConfigHist := createbamiloCatalogConfigMongoHist(bamiloCatalogConfig)
		mongointeract.UpsertConfigMongoHist(&bamiloCatalogConfigHist, mongoBulkHist)

		bamiloCatalogConfig := createbamiloCatalogConfigMongo(bamiloCatalogConfig)
		mongointeract.UpsertConfigMongo(&bamiloCatalogConfig, mongoBulk)
	}

	err = mongoBulkHist.Finish()
	checkError(err)
	err = mongoBulk.Finish()
	checkError(err)

}

func createbamiloCatalogConfigMongo(bamiloCatalogConfig bmlcatalogconfig.BmlCatalogConfig) (bamiloCatalogConfigMongo bmlconfigbackendprocess.Mongo) {

	now := time.Now()

	// only keep appropriate information for bmlcatalogconfig
	bamiloCatalogConfigMongo = bmlconfigbackendprocess.Mongo{
		IDBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
		ConfigSnapshotAt:   now,
		SKU:                bamiloCatalogConfig.SKU,
		SKUName:            bamiloCatalogConfig.SKUName,
		ImgLink:            bamiloCatalogConfig.ImgLink,
		SKULink:            bamiloCatalogConfig.SKULink,
		Description:        bamiloCatalogConfig.Description,
		ShortDescription:   bamiloCatalogConfig.ShortDescription,
		PackageContent:     bamiloCatalogConfig.PackageContent,
		ProductWarranty:    bamiloCatalogConfig.ProductWarranty,

		BiCategoryOneName:   bamiloCatalogConfig.BiCategoryOneName,
		BiCategoryTwoName:   bamiloCatalogConfig.BiCategoryTwoName,
		BiCategoryThreeName: bamiloCatalogConfig.BiCategoryThreeName,

		BrandName:   bamiloCatalogConfig.BrandName,
		BrandNameEn: bamiloCatalogConfig.BrandNameEn,

		SupplierName:   bamiloCatalogConfig.SupplierName,
		SupplierNameEn: bamiloCatalogConfig.SupplierNameEn,

		VisibleInShopBool: bamiloCatalogConfig.VisibleInShopBool,

		AvgPrice:           bamiloCatalogConfig.AvgPrice,
		AvgSpecialPrice:    bamiloCatalogConfig.AvgSpecialPrice,
		SumOfStockQuantity: bamiloCatalogConfig.SumOfStockQuantity,
		MinOfStockQuantity: bamiloCatalogConfig.MinOfStockQuantity,
	}

	if bamiloCatalogConfig.VisibleInShop == `1` {
		bamiloCatalogConfigMongo.SetVisibleInShopTrue()
	}

	return bamiloCatalogConfigMongo
}

func createbamiloCatalogConfigMongoHist(bamiloCatalogConfig bmlcatalogconfig.BmlCatalogConfig) (bamiloCatalogConfigMongoHist bmlconfigbackendprocess.MongoHist) {

	now := time.Now()

	iDBmlCatalogConfigHist, err := strconv.Atoi(strconv.Itoa(bamiloCatalogConfig.IDBmlCatalogConfig) + now.Format(`01022006`))
	checkError(err)

	// only keep appropriate information for bmlcatalogconfig
	bamiloCatalogConfigMongoHist = bmlconfigbackendprocess.MongoHist{
		IDBmlCatalogConfigHist: iDBmlCatalogConfigHist,
		FKBmlCatalogConfig:     bamiloCatalogConfig.IDBmlCatalogConfig,
		ConfigSnapshotAt:       now,

		VisibleInShopBool:  bamiloCatalogConfig.VisibleInShopBool,
		AvgPrice:           bamiloCatalogConfig.AvgPrice,
		AvgSpecialPrice:    bamiloCatalogConfig.AvgSpecialPrice,
		SumOfStockQuantity: bamiloCatalogConfig.SumOfStockQuantity,
		MinOfStockQuantity: bamiloCatalogConfig.MinOfStockQuantity,
	}

	if bamiloCatalogConfig.VisibleInShop == `1` {
		bamiloCatalogConfigMongoHist.SetVisibleInShopTrue()
	}

	return bamiloCatalogConfigMongoHist
}
