package export

// DailyTopCompetition represent on row of the CSV export in the backend process dailytopcompetition
type DailyTopCompetition struct {
	CompetitorName     string
	SKU                string
	SKUName            string
	ImgLink            string
	ProductWarranty    string
	CategoryOneName    string
	CategoryTwoName    string
	CategoryThreeName  string
	BrandName          string
	BrandNameEn        string
	SupplierName       string
	SupplierNameEn     string
	VisibleInShop      string
	AvgPrice           string
	AvgSpecialPrice    string
	SumOfStockQuantity string
	MinOfStockQuantity string
	SkuLink            string
	SkuRank            string
	Rating             string
	IsOutOfStock       string
	SumOfPaidPrice     string
}
