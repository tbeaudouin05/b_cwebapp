package backendprocess

import (
	"fmt"
	"time"


	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)
func AddAggBmlCategory()  {

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")
	bmlAggStatisticHist := db.C("bml_agg_statistic_hist")


	var categoryOneName [] string
	bmlCatalogConfig.Find(bson.M{}).Distinct("bi_category_one_name", &categoryOneName)
	fmt.Println(categoryOneName)
	
	//fmt.Println("supplierName:  ",supplierName) /// works --- alot of suppliers :)

	for _,rngCategoryOneName := range categoryOneName { 
		fmt.Println(":)")
		numCatOne,err := bmlCatalogConfig.Find(bson.M{"bi_category_one_name":rngCategoryOneName}).Count()
		checkErr(err)	
		var categoryTwoName [] string
		bmlCatalogConfig.Find(bson.M{"bi_category_one_name":rngCategoryOneName}).Distinct("bi_category_two_name", &categoryTwoName)
		for _,rngCategoryTwoName := range categoryTwoName{
			numCatTwo,err2 := bmlCatalogConfig.Find(bson.M{"bi_category_one_name":rngCategoryOneName,"bi_category_two_name":rngCategoryTwoName}).Count()
			checkErr(err2)	
			now := time.Now()
			var tmpCategoryAssortment bmldgktable.BmlCategorryAssortment
			tmpCategoryAssortment.ID = rngCategoryOneName+"/"+rngCategoryTwoName +now.Format("20060102")
			tmpCategoryAssortment.CategoryOneName = rngCategoryOneName
			tmpCategoryAssortment.CategoryTwoName = rngCategoryTwoName
			tmpCategoryAssortment.CategoryOneConfigCount = numCatOne
			tmpCategoryAssortment.CategoryTwoConfigCount = numCatTwo
			tmpCategoryAssortment.StatisticSnapshotAt = now
			tmpCategoryAssortment.TypeName = "BmlCategoryConfigCountTopPage"
		err = bmlAggStatisticHist.Insert(tmpCategoryAssortment)
		checkErr(err)
		tmpCategoryAssortment.TypeName = "BmlCategoryConfigCountTopPageLast"
		tmpCategoryAssortment.ID = rngCategoryOneName+"/"+rngCategoryTwoName 
		_,err = bmlAggStatisticHist.Upsert(bson.M{"id": tmpCategoryAssortment.ID} ,tmpCategoryAssortment)
		}
	}

}

func checkErr(err error){
	if err != nil {
		fmt.Println("Error " ,err)
	}
}