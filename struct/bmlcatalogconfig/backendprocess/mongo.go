package backendprocess

import (
	"time"
)

// config info -------------------------------------------------------------

// Mongo represents the data uploaded to MongoDB
type Mongo struct {
	// qualitative information
	IDBmlCatalogConfig int       `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	ConfigSnapshotAt   time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	SKU                string    `json:"sku" bson:"sku"`
	SKUName            string    `json:"sku_name" bson:"sku_name"`
	ImgLink            string    `json:"img_link" bson:"img_link"`
	SKULink            string    `json:"sku_link" bson:"sku_link"`
	Description        string    `json:"description" bson:"description"`
	ShortDescription   string    `json:"short_description" bson:"short_description"`
	PackageContent     string    `json:"package_content" bson:"package_content"`
	ProductWarranty    string    `json:"product_warranty" bson:"product_warranty"`
	// category
	BiCategoryOneName   string `json:"bi_category_one_name" bson:"bi_category_one_name"`
	BiCategoryTwoName   string `json:"bi_category_two_name" bson:"bi_category_two_name"`
	BiCategoryThreeName string `json:"bi_category_three_name" bson:"bi_category_three_name"`
	// brand
	BrandName   string `json:"brand_name" bson:"brand_name"`
	BrandNameEn string `json:"brand_name_en" bson:"brand_name_en"`
	// supplier
	SupplierName   string `json:"supplier_name" bson:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en" bson:"supplier_name_en"`
	// historical visibility
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           int `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    int `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity int `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity int `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
}

// SetVisibleInShopTrue sets VisibleInShop of Mongo to true. It is used to convert the integer of the generic BamiloCatalogConfig to MongoDB format of boolean.
func (bamiloCatalogConfigMongo *Mongo) SetVisibleInShopTrue() {

	bamiloCatalogConfigMongo.VisibleInShopBool = true

}

// config info hist ------------------------------------------

// MongoHist represents the historical data uploaded to MongoDB
type MongoHist struct {
	// qualitative information
	IDBmlCatalogConfigHist int       `json:"id_bml_catalog_config_hist" bson:"id_bml_catalog_config_hist"`
	FKBmlCatalogConfig     int       `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	ConfigSnapshotAt       time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at

	// historical visibility
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           int `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    int `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity int `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity int `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
}

// SetVisibleInShopTrue sets VisibleInShop of MongoHist to true. It is used to convert the integer of the generic BamiloCatalogConfig to MongoDB format of boolean
func (bamiloCatalogConfigMongoHist *MongoHist) SetVisibleInShopTrue() {

	bamiloCatalogConfigMongoHist.VisibleInShopBool = true

}
