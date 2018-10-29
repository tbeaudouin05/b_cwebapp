package backendprocess

import (
	"database/sql"
	"log"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
)

func AddBestPriceTag(db *sql.DB) (bmlCatalogConfigTable []bmlcatalogconfig.BmlCatalogConfig) {

	/*// qualitative information
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
	)*/

	// create bmlAllResult
	query := `CREATE TABLE bamiloCatalogConfigTable (
		IDBmlCatalogConfig TEXT
		,SKU     TEXT    
		,SKUName TEXT
		,ImgLink   TEXT
		,SKULink   TEXT
		,Description TEXT
		,ShortDescription  TEXT
		,PackageContent      TEXT
		,ProductWarranty  TEXT

		,BiCategoryOneName  TEXT
		,BiCategoryTwoName  TEXT
		,BiCategoryThreeName  TEXT

		,BrandName  TEXT
		,BrandNameEn  TEXT

		,SupplierName  TEXT
		,SupplierNameEn  TEXT

		,VisibleInShop TEXT
		,AvgPrice INTEGER
		,AvgSpecialPrice INTEGER
		,SumOfStockQuantity INTEGER
		,MinOfStockQuantity INTEGER
		)`
	queryP, err := db.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO bamiloCatalogConfigTable (
		IDBmlCatalogConfig 
		,SKU    
		,SKUName     
		,ImgLink   
		,SKULink   
		,Description 
		,ShortDescription  
		,PackageContent      
		,ProductWarranty  

		,BiCategoryOneName  
		,BiCategoryTwoName  
		,BiCategoryThreeName  

		,BrandName  
		,BrandNameEn  

		,SupplierName  
		,SupplierNameEn  

		,VisibleInShop 
		,AvgPrice 
		,AvgSpecialPrice 
		,SumOfStockQuantity 
		,MinOfStockQuantity 
		) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	queryP, err = db.Prepare(query)
	checkError(err)
	for i := 0; i < len(bmlCatalogConfigTable); i++ {

		bamiloCatalogConfig := bmlCatalogConfigTable[i]

		queryP.Exec(
			// qualitative information
			bamiloCatalogConfig.IDBmlCatalogConfig,
			bamiloCatalogConfig.SKU,
			bamiloCatalogConfig.SKUName,
			bamiloCatalogConfig.ImgLink,
			bamiloCatalogConfig.SKULink,
			bamiloCatalogConfig.Description,
			bamiloCatalogConfig.ShortDescription,
			bamiloCatalogConfig.PackageContent,
			bamiloCatalogConfig.ProductWarranty,
			// category
			bamiloCatalogConfig.BiCategoryOneName,
			bamiloCatalogConfig.BiCategoryTwoName,
			bamiloCatalogConfig.BiCategoryThreeName,
			// brand
			bamiloCatalogConfig.BrandName,
			bamiloCatalogConfig.BrandNameEn,
			// supplier
			bamiloCatalogConfig.SupplierName,
			bamiloCatalogConfig.SupplierNameEn,
			// historical data
			bamiloCatalogConfig.VisibleInShop,
			bamiloCatalogConfig.AvgPrice,
			bamiloCatalogConfig.AvgSpecialPrice,
			bamiloCatalogConfig.SumOfStockQuantity,
			bamiloCatalogConfig.MinOfStockQuantity,
		)
		time.Sleep(1 * time.Millisecond)

	}

	// Create view with one price
	query = `
		CREATE VIEW bamiloCatalogConfigTableOnePrice AS
		SELECT 
		IDBmlCatalogConfig 
		,SKU    
		,SKUName     
		,ImgLink   
		,SKULink   
		,Description 
		,ShortDescription  
		,PackageContent      
		,ProductWarranty  

		,BiCategoryOneName  
		,BiCategoryTwoName  
		,BiCategoryThreeName  

		,BrandName  
		,BrandNameEn  

		,SupplierName  
		,SupplierNameEn  

		,VisibleInShop 
		,AvgPrice 
		,AvgSpecialPrice 
		,SumOfStockQuantity 
		,MinOfStockQuantity
		CASE 
			WHEN  bar.avg_special_price <= 0 THEN bar.avg_price
			ELSE bar.avg_special_price END 'bml_price',
		CASE 
			WHEN  dar.avg_special_price <= 0 THEN dar.avg_price
			ELSE dar.avg_special_price END 'dgk_price'
		FROM bamiloCatalogConfigTable  bcct)
		`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	/* have two tables: one with fk_marketplace_product_group, one without it. If fk_marketplace_product_group exists then take the best price among this group
	if it does not exist then just consider the product as being the best price
	*/

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
