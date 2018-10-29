package main 

import(
	"time"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	//"github.com/thomas-bamilo/commercial/competitionanalysis/backendprocess/sendemail/usermatchreport"
)

func main(){

// Connection URL
var url = `mongodb://localhost:27017/competition_analysis`
mongoSession, err := mgo.Dial(url)
checkErr(err)

defer mongoSession.Close()
mongoSession.SetMode(mgo.Monotonic, true)
db := mongoSession.DB("competition_analysis")
bmlCatalogConfig := db.C("bml_catalog_config")
bmlAggStatisticHist := db.C("bml_agg_statistic_hist")
bmlDgkAggStatisticHist := db.C("bml_dgk_agg_statistic_hist")

userCollection := db.C("user")

now := time.Now()
lastWeek := now.AddDate(0, 0, -7)
lastTwoWeek := lastWeek.AddDate(0, 0, -7)

var userEmail []string
userCollection.Find(nil).Distinct("email", &userEmail)
for _,rngUserEmail := range userEmail {
	cnt1, err1 := bmlCatalogConfig.Find(bson.M{"matched_by_email":rngUserEmail}).Count()
	cnt2, err2 := bmlAggStatisticHist.Find(bson.M{"matched_by_email":rngUserEmail}).Count()
	checkErr(err1)
	checkErr(err2)
	_,err :=userCollection.Upsert(bson.M{"email": rngUserEmail},bson.M{"$set":bson.M{"total_matched_sku":cnt1 ,"total_matched_supplier":cnt2 }})
	checkErr(err)
}

for _,rngUserEmail := range userEmail {
	
	cnt1, err1 := bmlDgkAggStatisticHist.Find(bson.M{"matched_by_email":rngUserEmail}).Count()
	checkErr(err1)
	_,err :=userCollection.Upsert(bson.M{"email": rngUserEmail},bson.M{"$set":bson.M{"total_not_competitive":cnt1 }})
	checkErr(err)

}



var userEmailMatchedSkuThisWeek  []string
bmlCatalogConfig.Find(bson.M{ "good_match_at": bson.M{"$gte": lastWeek} }).Distinct("matched_by_email", &userEmailMatchedSkuThisWeek)
for _,rngUserEmail := range userEmailMatchedSkuThisWeek {
	cnt, err1 := bmlCatalogConfig.Find(bson.M{"matched_by_email":rngUserEmail,"good_match_at": bson.M{"$gte": lastWeek}}).Count()
	checkErr(err1)
	err2:=userCollection.Update(bson.M{"email": rngUserEmail},bson.M{"$set":bson.M{"last_7_day_matched_sku":cnt }})
	checkErr(err2)
} 

var userEmailMatchedSkuLastWeek [] string
bmlCatalogConfig.Find(bson.M{ "good_match_at": bson.M{"$gte": lastTwoWeek, "$lt":lastWeek }}).Distinct("matched_by_email", &userEmailMatchedSkuLastWeek)
	for _,rngUserEmail := range userEmailMatchedSkuLastWeek {
	cnt, err1 := bmlCatalogConfig.Find(bson.M{"matched_by_email":rngUserEmail,"good_match_at": bson.M{"$gte": lastTwoWeek, "$lt":lastWeek}}).Count()
	checkErr(err1)
	err2:=userCollection.Update(bson.M{"email": rngUserEmail},bson.M{"$set":bson.M{ "last_14_day_matched_sku":cnt }})
	checkErr(err2)
	}

	var userEmailMatchedSupplierThisWeek [] string
	bmlAggStatisticHist.Find(bson.M{ "good_match_at": bson.M{"$gte": lastWeek}}).Distinct("matched_by_email", &userEmailMatchedSupplierThisWeek)
	for _,rngUserEmail := range userEmailMatchedSupplierThisWeek {
		cnt, err1 := bmlAggStatisticHist.Find(bson.M{"matched_by_email":rngUserEmail,"good_match_at": bson.M{"$gte": lastWeek}}).Count()
		checkErr(err1)
		err2:=userCollection.Update(bson.M{"email": rngUserEmail},bson.M{"$set":bson.M{ "last_7_day_matched_supplier":cnt }})
		checkErr(err2)
	} 
	
	var userEmailMatchedSupplierLastWeek [] string
	bmlAggStatisticHist.Find(bson.M{ "good_match_at": bson.M{"$gte": lastTwoWeek, "$lt":lastWeek }}).Distinct("matched_by_email", &userEmailMatchedSupplierLastWeek)
	for _,rngUserEmail := range userEmailMatchedSupplierLastWeek {
		cnt, err1 := bmlAggStatisticHist.Find(bson.M{"matched_by_email":rngUserEmail,"good_match_at": bson.M{"$gte": lastTwoWeek, "$lt":lastWeek}}).Count()
		checkErr(err1)
		err2:=userCollection.Update(bson.M{"email": rngUserEmail},bson.M{"$set":bson.M{"last_14_day_matched_supplier":cnt }})
		checkErr(err2)
	}

	//usermatchreport.SendUserMatchedReport(mongoSession)
}


	

func checkErr(err error){
	if err != nil {
		fmt.Println(err)
	}
}
