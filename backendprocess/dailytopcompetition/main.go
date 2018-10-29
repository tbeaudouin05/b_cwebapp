package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gocarina/gocsv"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/export"
	"github.com/thomas-bamilo/email/goemail"
)

func main() {

	// Connection URL
	//jsekjndc
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)
	db := mongoSession.DB("competition_analysis")

	// Find distincted categories in each database
	var dgkCategory []string
	var bmlCategory []string
	db.C("dgk_catalog_config").Find(bson.M{}).Distinct("dgk_category_two_name", &dgkCategory)
	db.C("bml_catalog_config").Find(bson.M{}).Distinct("bi_category_two_name", &bmlCategory)

	// Select top 500 Digikala products in each category based on "sku_rank" and store it in "dgkAllResult"
	var dgkAllResult []dgkcatalogconfig.DgkCatalogConfig
	for _, cat := range dgkCategory {
		var dgkResult []dgkcatalogconfig.DgkCatalogConfig
		db.C("dgk_catalog_config").Find(bson.M{"dgk_category_two_name": cat}).Select(bson.M{
			"avg_price":             12,
			"avg_special_price":     12,
			"dgk_category_one_name": 12,
			"dgk_category_two_name": 12,
			"img_link":              12,
			"is_out_of_stock":       12,
			"rating":                12,
			"sku_name":              12,
			"sku_rank":              12,
			"sku_link":              12}).Sort("sku_rank").Limit(300).All(&dgkResult)

		dgkAllResult = append(dgkAllResult, dgkResult...)
	}

	// Select top 500 Bamilo products in each category based on "sum_of_paid_price" and store it in "bmlAllResult"
	var bmlAllResult []bmlcatalogconfig.BmlCatalogConfig
	for _, cat := range bmlCategory {
		var bmlResult []bmlcatalogconfig.BmlCatalogConfig

		db.C("bml_catalog_config").Find(bson.M{"bi_category_two_name": cat}).Select(bson.M{
			"sku":                    17,
			"sku_name":               17,
			"img_link":               17,
			"product_warranty":       17,
			"bi_category_one_name":   17,
			"bi_category_two_name":   17,
			"bi_category_three_name": 17,
			"brand_name":             17,
			"brand_name_en":          17,
			"supplier_name":          17,
			"supplier_name_en":       17,
			"visible_in_shop":        17,
			"avg_price":              17,
			"avg_special_price":      17,
			"sum_of_paid_price":      17,
			"sum_of_stock_quantity":  17,
			"min_of_stock_quantity":  17}).Sort("-sum_of_paid_price").Limit(300).All(&bmlResult)

		bmlAllResult = append(bmlAllResult, bmlResult...)
	}

	// change the format of the selected data and append together
	var top []export.DailyTopCompetition
	for _, result := range dgkAllResult {
		top = append(top, dgkToGeneral(result))
	}
	for _, result := range bmlAllResult {
		top = append(top, bmlToGeneral(result))
	}

	//output
	file, err := os.OpenFile("DailyTopCompetition.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	gocsv.MarshalFile(&top, file)
	time.Sleep(1 * time.Millisecond)
	goemail.GoEmail()

}

func dgkToGeneral(dgk dgkcatalogconfig.DgkCatalogConfig) export.DailyTopCompetition {
	mrg := export.DailyTopCompetition{
		SKU:                "",
		SKUName:            dgk.SkuName,
		ImgLink:            dgk.ImgLink,
		ProductWarranty:    "",
		CategoryOneName:    dgk.DgkCategoryOneName,
		CategoryTwoName:    dgk.DgkCategoryTwoName,
		CategoryThreeName:  "",
		BrandName:          "",
		BrandNameEn:        "",
		SupplierName:       "",
		SupplierNameEn:     "",
		VisibleInShop:      "",
		AvgPrice:           strconv.Itoa(dgk.AvgPrice),
		AvgSpecialPrice:    strconv.Itoa(dgk.AvgSpecialPrice),
		SumOfStockQuantity: "",
		MinOfStockQuantity: "",
		SkuLink:            dgk.SkuLink,
		SkuRank:            strconv.Itoa(dgk.SkuRank),
		Rating:             strconv.Itoa(dgk.Rating),
		IsOutOfStock:       strconv.FormatBool(dgk.IsOutOfStock),
		CompetitorName:     "Digikala",
	}

	return mrg
}
func bmlToGeneral(bml bmlcatalogconfig.BmlCatalogConfig) export.DailyTopCompetition {
	mrg := export.DailyTopCompetition{
		SKU:                bml.SKU,
		SKUName:            bml.SKUName,
		ImgLink:            bml.ImgLink,
		ProductWarranty:    "",
		CategoryOneName:    bml.BiCategoryOneName,
		CategoryTwoName:    bml.BiCategoryTwoName,
		CategoryThreeName:  bml.BiCategoryThreeName,
		BrandName:          bml.BrandName,
		BrandNameEn:        bml.BrandName,
		SupplierName:       bml.SupplierName,
		SupplierNameEn:     bml.SupplierNameEn,
		VisibleInShop:      strconv.FormatBool(bml.VisibleInShopBool),
		AvgPrice:           strconv.Itoa(bml.AvgPrice),
		AvgSpecialPrice:    strconv.Itoa(bml.AvgSpecialPrice),
		SumOfStockQuantity: strconv.Itoa(bml.SumOfStockQuantity),
		MinOfStockQuantity: strconv.Itoa(bml.MinOfStockQuantity),
		SumOfPaidPrice:     strconv.Itoa(bml.SumOfPaidPrice),
		SkuLink:            "",
		SkuRank:            "",
		Rating:             "",
		IsOutOfStock:       "",
		CompetitorName:     "Bamilo",
	}
	return mrg
}
