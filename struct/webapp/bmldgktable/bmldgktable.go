package bmldgktable

import (
	"time"
	//"github.com/globalsign/mgo/bson"
)

// BmlDgkConfigTable is the data sent to frontend to generate the table with both Bamilo and Digikala data
type BmlDgkConfigTable struct {
	IDBmlCatalogConfig      int
	IDDgkCatalogConfig      string
	SKU                     string
	DgkSKUName              string
	BamiloSKUName           string
	DgkImgLink              string
	BamiloImgLink           string
	ProductWarranty         string
	DgkCategoryOneName      string
	DgkCategoryTwoName      string
	BamiloCategoryOneName   string
	BamiloCategoryTwoName   string
	BamiloCategoryThreeName string
	BrandName               string
	BrandNameEn             string
	SupplierName            string
	SupplierNameEn          string
	VisibleInShop           string
	BamiloAvgPrice          string
	BamiloAvgSpecialPrice   string
	BamiloPrice             string
	DgkAvgPrice             string
	DgkAvgSpecialPrice      string
	DgkPrice                string
	SumOfStockQuantity      string
	MinOfStockQuantity      string
	BamiloSKULink           string
	DgkSKULink              string
	SkuRank                 int
	Rating                  string
	IsOutOfStock            string
	SumOfPaidPrice          string
	DgkScore                float64
	BmlConfigSnapshotAt     string
	DgkConfigSnapshotAt     string
	GoodMatch               bool
	CountOfSoi              int
}

type OptionList struct {
	OptionValue string `json:"optionValue"`
	OptionText  string `json:"optionText"`
}
type FrequencyOptionList struct {
	OptionValue int    `json:"OptionValue"`
	OptionText  string `json:"OptionText"`
}

type Frequency struct {
	FrequencyOptionList []FrequencyOptionList `json:"FrequencyOptionList"`
}
type Category struct {
	OptionList []OptionList `json:"OptionList"`
}
type RowValue struct {
	BmlIDCatalogConfig    string `json:"BmlIDCatalogConfig"`
	BmlSKUName            string `json:"BmlSKUName"`
	BmlImgLink            string `json:"BmlImgLink"`
	BmlSKULink            string `json:"BmlSKULink"`
	BmlSKUPrice           string `json:"BmlSKUPrice"`
	BmlAvgPrice           string `json:"BmlAvgPrice"`
	BmlAvgSpecialPrice    string `json:"BmlAvgSpecialPrice"`
	DgkIDCatalogConfig    string `json:"DgkIDCatalogConfig"`
	DgkScore              int    `json:"DgkScore"`
	DgkSKUName            string `json:"DgkSKUName"`
	DgkImgLink            string `json:"DgkImgLink"`
	DgkSKULink            string `json:"DgkSKULink"`
	DgkSKUPrice           string `json:"DgkSKUPrice"`
	DgkAvgPrice           string `json:"DgkAvgPrice"`
	DgkAvgSpecialPrice    string `json:"DgkAvgSpecialPrice"`
	BmlConfigSnapshotAt   string `json:"BmlConfigSnapshot"`
	DgkConfigSnapshotAt   string `json:"DgkConfigSnapshot"`
	BmlSupplierName       string `json:"BmlSupplierName"`
	BmlBrand              string `json:"BmlBrand"`
	BmlMinOfStockQuantity string `json:"BmlMinOfStockQuantity"`
	BmlSumOfStockQuantity string `json:"BmlSumOfStockQuantity"`
	DgkStock              string `json:"DgkStock"`
	GoodMatch             bool   `json:"GoodMatch"`
}
type Row struct {
	RowValue RowValue `json:"RowValue"`
	RowKey   int      `json:"RowKey"`
}
type Table struct {
	Table []Row `json:""`
}
type TableInterface struct {
	Table      []Row `json:"Table"`
	PageNumber []int `json:"PageNumberList"`
}
type IDTable struct {
	ID []int `json:"_id" bson:"_id"`
}

type SKUName struct {
	BmlSKUName string `json:"BmlSKUName"`
	DgkSKUName string `json:"DgkSKUName"`
}
type Value struct {
	BmlPrice []int `json:"BmlPrice"`
	BmlSales []int `json:"BmlSales"`
	DgkPrice []int `json:"DgkPrice"`
}
type HistoricalChartData struct {
	Label   []string `json:"Label"`
	Value   Value    `json:"Value"`
	SKUName SKUName  `json:"SKUName"`
}
type DgkRowValue struct {
	//DgkScore *float64 `json:"dgkScore"`
	DgkIDCatalogConfig string `json:"DgkIDCatalogConfig"`
	DgkSKUName         string `json:"DgkSKUName"`
	DgkImgLink         string `json:"DgkImgLink"`
	DgkSKULink         string `json:"DgkSKULink"`
	DgkSKUPrice        string `json:"DgkSKUPrice"`
}
type DgkRow struct {
	RowValue DgkRowValue `json:"DgkRowValue"`
	RowKey   int         `json:"RowKey"`
}
type DgkTable struct {
	Table []DgkRow `json:""`
}
type DgkTableInterface struct {
	Table      []DgkRow `json:"Table"`
	PageNumber []int    `json:"PageNumberList"`
	BmlSkuName string   `json:"BmlSkuName"`
	BmlPrice   string   `json:"BmlPrice"`
	BmlImgLink string   `json:"BmlImgLink"`
	BmlBrand   string   `json:"BmlBrand"`
	BmlSKULink string   `json:"BmlSKULink"`
}

type DgkSellerUrl struct {
	Url      string `bson:"url"`
	Priority int    `bson:"priority"`
	ID       int    `bson:"_id"`
	TypeName string `bson:"type"`
}
type DgkProductUrl struct {
	Url      string `bson:"url"`
	ID       string    `bson:"_id"`
	BmlID    int    `bson:"bml_sku_id"`
	TypeName string `bson:"type"`
}

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
