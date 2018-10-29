package user

type UserInfo struct {
	
	Email  string `json:"email" bson:"email"`
	Name   string `json:"name" bson:"name"`
	TotalSKU  int `json:"totalSKU" bson:"total_matched_sku"`
	TotalSupplier   int `json:"totalSupplier" bson:"total_matched_supplier"`
	ThisWeekSKU  int `json:"thisWeekSKU" bson:"last_7_day_matched_sku"`
	LastWeekSKU  int `json:"lastWeekSKU" bson:"last_14_day_matched_sku"`
	TotalNotCompetitive int `json:"TotalNotCompetitive" bson:"total_not_competitive"`
	Message string `json:"Message"`
}
