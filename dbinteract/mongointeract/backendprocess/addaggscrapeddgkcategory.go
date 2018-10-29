package backendprocess

import (
	"fmt"
	"time"


	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)
func AddAggScrapedDgkCategory()  {

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	db := mongoSession.DB("competition_analysis")
	dgkCatalogConfig := db.C("dgk_catalog_config")
	dgkAggStatisticHist := db.C("dgk_agg_statistic_hist")


	var categoryOneName [] string
	dgkCatalogConfig.Find(bson.M{}).Distinct("dgk_category_one_name", &categoryOneName)
	fmt.Println(categoryOneName)
	
	//fmt.Println("supplierName:  ",supplierName) /// works --- alot of suppliers :)

	for _,rngCategoryOneName := range categoryOneName { 
		fmt.Println(":)")
		numCatOne,err := dgkCatalogConfig.Find(bson.M{"dgk_category_one_name":rngCategoryOneName}).Count()
		checkErr(err)	
		var categoryTwoName [] string
		dgkCatalogConfig.Find(bson.M{"dgk_category_one_name":rngCategoryOneName}).Distinct("dgk_category_two_name", &categoryTwoName)
		for _,rngCategoryTwoName := range categoryTwoName{
			numCatTwo,err2 := dgkCatalogConfig.Find(bson.M{"dgk_category_one_name":rngCategoryOneName,"dgk_category_two_name":rngCategoryTwoName}).Count()
			checkErr(err2)	
			var tmpCategoryAssortment bmldgktable.DgkScrapedCategorryAssortment
			tmpCategoryAssortment.ID = rngCategoryOneName+"/"+rngCategoryTwoName
			tmpCategoryAssortment.CategoryOneName = rngCategoryOneName
			tmpCategoryAssortment.CategoryTwoName = rngCategoryTwoName
			tmpCategoryAssortment.CategoryOneConfigCount = numCatOne
			tmpCategoryAssortment.CategoryTwoConfigCount = numCatTwo
			tmpCategoryAssortment.StatisticSnapshotAt = time.Now()
			tmpCategoryAssortment.TypeName = "DgkScrapedCategoryConfigCountTopPage"
		err = dgkAggStatisticHist.Insert(tmpCategoryAssortment)
		checkErr(err)
		}
	}

}
