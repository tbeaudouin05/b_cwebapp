package configcountsupplierpage

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"

)

func ApplySellerURL(DgkUrl string, mongoSession *mgo.Session) {

	db := mongoSession.DB("competition_analysis")
	dgkUrlCol := db.C("dgk_url")

	//var dgkUrlStr string
	//queries
	var dgkUrlStr [] bmldgktable.DgkSellerUrl


	dgkUrlCol.Find(bson.M{"url":DgkUrl }).All(&dgkUrlStr)
	if dgkUrlStr!=nil{
		err := dgkUrlCol.UpdateId(dgkUrlStr[0].ID, bson.M{"$set": bson.M{"priority":1 }})
		CheckErr(err)
	}else{
		var dgkUrlAdd bmldgktable.DgkSellerUrl
		dgkUrlAdd.Url = DgkUrl
		dgkUrlAdd.Priority = 1
		err := dgkUrlCol.Insert(dgkUrlAdd)
		CheckErr(err)
	}

}

func CheckErr(err error){
	if err!=nil{
		fmt.Println(err)
	}
}