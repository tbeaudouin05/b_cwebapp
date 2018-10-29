package dgkcatalogconfig

import "time"

// DgkCatalogConfig is the most generic struct representing one DgkCatalogConfig. It is used in dailytopcompetition and mongointeract > websocketserver
type DgkCatalogConfig struct {
	IDDgkCatalogConfig        string    `csv:"id_dgk_catalog_config" bson:"id_dgk_catalog_config"` // removed for export
	Version                   int       `csv:"__v" bson:"__v"`
	AvgPrice                  int       `csv:"avg_price" bson:"avg_price"`
	AvgSpecialPrice           int       `csv:"avg_special_price" bson:"avg_special_price"`
	DgkCategoryOneName        string    `csv:"dgk_category_one_name" bson:"dgk_category_one_name"`
	DgkCategoryTwoName        string    `csv:"dgk_category_two_name" bson:"dgk_category_two_name"`
	ImgLink                   string    `csv:"img_link" bson:"img_link"`
	IsOutOfStock              bool      `csv:"is_out_of_stock" bson:"is_out_of_stock"`
	Rating                    int       `csv:"rating" bson:"rating"`
	SkuName                   string    `csv:"sku_name" bson:"sku_name"`
	SkuRank                   int       `csv:"sku_rank" bson:"sku_rank"`
	SkuLink                   string    `csv:"sku_link" bson:"sku_link"`
	ConfigSnapshotAt          time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"`
	BestPrice                 int       `csv:"best_price" bson:"best_price"`
	BestPriceUpdateSnapshotAt time.Time `csv:"best_price_update_snapshot_at" bson:"best_price_update_snapshot_at"`
	ColorOfBestPrice          string    `csv:"color_of_best_price" bson:"color_of_best_price"`
	SupplierOfBestPrice       string    `csv:"supplier_of_best_price" bson:"supplier_of_best_price"`
	WarrantyOfBestPrice       string    `csv:"warranty_of_best_price" bson:"warranty_of_best_price"`
}

type DgkCatalogConfigHist struct {
	AvgPrice         int       `csv:"avg_price" bson:"avg_price"`
	ConfigSnapshotAt time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"`
	AvgSpecialPrice  int       `csv:"avg_special_price" bson:"avg_special_price"`
}
