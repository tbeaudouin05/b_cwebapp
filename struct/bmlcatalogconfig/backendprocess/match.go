package backendprocess

// config info -------------------------------------------------------------

// Match is the data uploaded to MongoDB to update Bamilo documents with fk_dgk_catalog_config and dgk_score which represent the best Digikala match for each BamiloCatalogConfig
type Match struct {
	IDBmlCatalogConfig int     `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	IDDgkCatalogConfig string  `json:"id_dgk_catalog_config" bson:"fk_dgk_catalog_config"`
	DgkScore           float64 `json:"dgk_score" bson:"dgk_score"`
	Matched            bool    `json:"matched" bson:"matched"`
}
