package supplier

import (
	"time"
	//"github.com/globalsign/mgo/bson"
)

type SupplierHist struct {
	SupplierName        string    `bson:"supplier_name" json:"SellerName"`
	SupplierNameEn      string    `bson:"supplier_name_en"`
	ID                  int       `bson:"_id"  json:"ID"  `
	ConfigCount         int       `bson:"config_count" json:"ConfigCount"`
	InStockConfigCount  int       `bson:"in_stock_config_count" json:"InStockConfigCount"`
	TypeName            string    `bson:"type"`
	StatisticSnapshotAt time.Time `bson:"statistic_snapshot_at"`
}
type Supplier struct {
	SupplierName       string `bson:"supplier_name" json:"SellerName"`
	SupplierNameEn     string `bson:"supplier_name_en"`
	ConfigCount        int    `bson:"config_count" json:"ConfigCount"`
	InStockConfigCount int    `bson:"in_stock_config_count" json:"InStockConfigCount"`
	TypeName           string `bson:"type"`
	//IDSring string `bson:"_id"  json:"IDString"  `
	ID                  int       `bson:"_id"  json:"ID"  `
	SupplierCode        int       `bson:"supplier_code"`
	SupplierID          int       `bson:"id_supplier"`
	RowKey              int       `bson:"rowkey" json:"RowKey"`
	StatisticSnapshotAt time.Time `bson:"statistic_snapshot_at"`
}

type BmlDgkSupplier struct {
	BmlSupplierName string `bson:"supplier_name" json:"BmlSellerName"`
	DgkSupplierName string `bson:"dgk_matched_supplier_name" json:"DgkSellerName"`
}
type DgkSupplier struct {
	SupplierName       string `bson:"supplier_name" json:"SellerName"`
	SupplierNameEn     string `bson:"supplier_name_en"`
	ConfigCount        int    `bson:"config_count" json:"ConfigCount"`
	InStockConfigCount int    `bson:"in_stock_config_count" json:"InStockConfigCount"`
	TypeName           string `bson:"type"`
	//IDSring string `bson:"_id"  json:"IDString"  `
	ID                  string    `bson:"_id"  json:"ID"  `
	SupplierCode        int       `bson:"supplier_code"`
	RowKey              int       `bson:"rowkey" json:"RowKey"`
	StatisticSnapshotAt time.Time `bson:"statistic_snapshot_at"`
}


/*type SupplierMatch struct {
	BmlSupplierID   int    `bson:"_id"`
	BmlSupplierName string `bson:"bml_supplier_name"`
	DgkSupplierName string `bson:"dgk_supplier_name"`
	DgkSupplierID   string `bson:"dgk_id"`
	Type            string `bson:"type"`
}*/

type DgkCategorryAssortment struct {
	CategoryOneName     string    `bson:"dgk_category_one_name"`
	CategoryTwoName     string    `bson:"dgk_category_two_name"`
	ConfigCount         int       `bson:"config_count"`
	TypeName            string    `bson:"type"`
	ID                  string    `bson:"_id"`
	StatisticSnapshotAt time.Time `bson:"statistic_snapshot_at"`
}

type DgkScrapedCategorryAssortment struct {
	CategoryOneName        string    `bson:"dgk_category_one_name"`
	CategoryTwoName        string    `bson:"dgk_category_two_name"`
	CategoryOneConfigCount int       `bson:"category_one_config_count"`
	CategoryTwoConfigCount int       `bson:"category_two_config_count"`
	TypeName               string    `bson:"type"`
	ID                     string    `bson:"id"`
	StatisticSnapshotAt    time.Time `bson:"statistic_snapshot_at"`
}

type BmlCategorryAssortment struct {
	CategoryOneName        string `bson:"bi_category_one_name"`
	CategoryTwoName        string `bson:"bi_category_two_name"`
	CategoryOneConfigCount int    `bson:"category_one_config_count"`
	CategoryTwoConfigCount int    `bson:"category_two_config_count"`
	TypeName               string `bson:"type"`
	//ID bson.ObjectId `bson:"_id"`
	ID                  string    `bson:"id"`
	StatisticSnapshotAt time.Time `bson:"statistic_snapshot_at"`
}

type CategoryTwo struct {
	CategoryTwoName    string `json:"CategoryTwoName"`
	ConfigCountCat2    string `json:"ConfigCountCat2"`
	ConfigCountCat2Int int    `json:"ConfigCountCat2Int"`

	RowKey int `json:"RowKey"`
}

type CategoryAssortmentTable struct {
	CategoryOneName     string        `json:"CategoryOneName"`
	ConfigCountCat1Int  int           `json:"ConfigCountCat1Int"`
	ConfigCountCat1     string        `json:"ConfigCountCat1"`
	CategoryTwo         []CategoryTwo `json:"CategoryTwo"`
	StatisticSnapshotAt time.Time     `bson:"statistic_snapshot_at"`

	RowKey int `json:"RowKey"`
}
