package bmlcatalogconfig

import (
	"log"
	"math/rand"
	"time"
)

type IdBmlCatalogConfig struct {
	// qualitative information
	IDBmlCatalogConfig int `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
}

// BmlCatalogConfig is the most generic struct representing one BmlCatalogConfig. It is used in dailytopcompetition, elasticinteract, mongointeract & bobinteract backendprocess and  mongointeract > weboscketserver
type BmlCatalogConfig struct {
	// qualitative information
	IDBmlCatalogConfig        int       `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	FKBmlCatalogConfig        string    `json:"fk_dgk_catalog_config" bson:"fk_dgk_catalog_config"`
	FKCatalogConfigGroup      string    `json:"fk_catalog_config_group" bson:"fk_catalog_config_group"`
	FKMarketplaceProductGroup string    `json:"fk_marketplace_product_group" bson:"fk_marketplace_product_group"`
	ManualFKBmlCatalogConfig  string    `json:"manual_fk_dgk_catalog_config" bson:"manual_fk_dgk_catalog_config"`
	IDBmlCatalogConfigHist    int       `json:"id_bml_catalog_config_hist" bson:"id_bml_catalog_config_hist"`
	ConfigSnapshotAt          time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	SKU                       string    `json:"sku" bson:"sku"`
	SKUName                   string    `json:"sku_name" bson:"sku_name"`
	ImgLink                   string    `json:"img_link" bson:"img_link"`
	SKULink                   string    `json:"sku_link" bson:"sku_link"`
	Description               string    `json:"description" bson:"description"`
	ShortDescription          string    `json:"short_description" bson:"short_description"`
	PackageContent            string    `json:"package_content" bson:"package_content"`
	ProductWarranty           string    `json:"product_warranty" bson:"product_warranty"`
	// category
	BiCategory          string `json:"bi_category" bson:"bi_category"`
	BiCategoryOneName   string `json:"bi_category_one_name" bson:"bi_category_one_name"`
	BiCategoryTwoName   string `json:"bi_category_two_name" bson:"bi_category_two_name"`
	BiCategoryThreeName string `json:"bi_category_three_name" bson:"bi_category_three_name"`
	BiCategoryFourName  string `json:"bi_category_four_name" bson:"bi_category_four_name"`
	BiCategoryFiveName  string `json:"bi_category_five_name" bson:"bi_category_five_name"`
	BiCategorySixName   string `json:"bi_category_six_name" bson:"bi_category_six_name"`
	// brand
	BrandName   string `json:"brand_name" bson:"brand_name"`
	BrandNameEn string `json:"brand_name_en" bson:"brand_name_en"`
	// department
	Department string `json:"department" bson:"department"`
	// department
	KeyAccountManager string `json:"key_account_manager" bson:"key_account_manager"`
	// supplier
	SupplierName   string `json:"supplier_name" bson:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en" bson:"supplier_name_en"`
	// historical visibility
	VisibleInShop     string
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           int `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    int `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity int `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity int `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
	// sales
	CountOfSoi            int                 `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        int                 `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        int                 `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue int                 `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount int                 `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
	Match                 []map[string]string `bson:"match" json:"match"`
	DgkScore              float64             `json:"dgk_score" bson:"dgk_score"`
	StoredDgkScore        float64             `json:"stored_dgk_score" bson:"stored_dgk_score"`
	GoodMatch             bool                `json:"good_match" bson:"good_match"`
	//user
	MatchedByEmail		  string `json:"MatchedByEmail" bson:"matched_by_email"`
}

// DivideBmlCatalogConfigTableInRandomChunk divides an array of BmlCatalogConfig into several random arrays of BmlCatalogConfig
func DivideBmlCatalogConfigTableInRandomChunk(bamiloCatalogConfigTable []BmlCatalogConfig) (arrayOfBmlCatalogConfigTableRandomChunk [][]BmlCatalogConfig) {

	bamiloCatalogConfigTableRand := make([]BmlCatalogConfig, len(bamiloCatalogConfigTable))
	permutation := rand.Perm(len(bamiloCatalogConfigTable))
	for originalOrder, permutatedOrder := range permutation {
		bamiloCatalogConfigTableRand[permutatedOrder] = bamiloCatalogConfigTable[originalOrder]
	}

	chunkSize := 950

	for i := 0; i < len(bamiloCatalogConfigTableRand); i += chunkSize {
		end := i + chunkSize

		if end > len(bamiloCatalogConfigTableRand) {
			end = len(bamiloCatalogConfigTableRand)
		}

		arrayOfBmlCatalogConfigTableRandomChunk = append(arrayOfBmlCatalogConfigTableRandomChunk, bamiloCatalogConfigTableRand[i:end])
	}

	return arrayOfBmlCatalogConfigTableRandomChunk
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// elastic -----------------------------------------------------------------------------------------------------------------------------------------------

/*type BmlCatalogConfigElastic struct {
	// qualitative information
	IDBmlCatalogConfig int    `json:"id_bml_catalog_config"`
	SKUName            string `json:"sku_name"`
	SKUNameEn          string `json:"sku_name_en"`
	ImgLink            string `json:"img_link"`
	Description        string `json:"description"`
	ShortDescription   string `json:"short_description"`
	PackageContent     string `json:"package_content"`
	ProductWarranty    string `json:"product_warranty"`
	// category
	BiCategoryOneName   string `json:"bi_category_one_name"`
	BiCategoryTwoName   string `json:"bi_category_two_name"`
	BiCategoryThreeName string `json:"bi_category_three_name"`
	// brand
	BrandName   string `json:"brand_name"`
	BrandNameEn string `json:"brand_name_en"`
	// supplier
	SupplierName   string `json:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en"`
	// historical price and quantity
	AvgPrice        int `json:"avg_price"`
	AvgSpecialPrice int `json:"avg_special_price"`
}*/
