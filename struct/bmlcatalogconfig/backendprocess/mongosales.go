package backendprocess

import (
	"time"
)

// config sales -------------------------------------------------------------------------

// MongoSales represents the sales data uploaded to MongoDB. Sales cannot be uploaded at the same time as other data due to the structure of the source MySQL database
type MongoSales struct {
	// qualitative information
	IDBmlCatalogConfig int       `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	ConfigSnapshotAt   time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"`
	// sales
	CountOfSoi            int `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        int `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        int `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue int `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount int `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
}

// config sales hist -----------------------

// MongoSalesHist represents the historical sales data uploaded to MongoDB. Sales cannot be uploaded at the same time as other data due to the structure of the source MySQL database
type MongoSalesHist struct {
	// qualitative information
	IDBmlCatalogConfigHist int       `json:"id_bml_catalog_config_hist" bson:"id_bml_catalog_config_hist"`
	FKBmlCatalogConfig     int       `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	ConfigSnapshotAt       time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	// sales
	CountOfSoi            int `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        int `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        int `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue int `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount int `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
}
