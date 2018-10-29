package bmldgkskuhistory

import "time"

// BmlDgkConfigHist is the data sent to frontend to generate the chart of bmldgkskuhistory
type BmlDgkConfigHist struct {
	FKBmlCatalogConfig int       `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	SumOfPaidPrice     int       `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	AvgPrice           int       `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice           int       `json:"avg_special_price" bson:"avg_special_price"`
	ConfigSnapshotAt   time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"`
}
