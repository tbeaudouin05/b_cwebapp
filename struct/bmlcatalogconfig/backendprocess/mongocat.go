package backendprocess

type BiRow struct {
	// BI Center
	CatId             int    `bson:"_id"`
	CatBi             string `bson:"bi_category"`
	CatEn1            string `bson:"bi_category_one_name"`
	CatEn2            string `bson:"bi_category_two_name"`
	CatEn3            string `bson:"bi_category_three_name"`
	CatEn4            string `bson:"bi_category_four_name"`
	CatEn5            string `bson:"bi_category_five_name"`
	CatEn6            string `bson:"bi_category_six_name"`
	Department        string `bson:"department"`
	KeyAccountManager string `bson:"key_account_manager"`
}
