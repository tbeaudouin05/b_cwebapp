package bmldgkmanualmatching

import (
	"time"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
"fmt"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"

)

func ApplyProductURL(DgkUrl string, BmlID int, mongoSession *mgo.Session) (int,string) {

	db := mongoSession.DB("competition_analysis")
	dgkUrlCol := db.C("dgk_url")

	//var dgkUrlStr string
	//queries
	var dgkUrlStr [] bmldgktable.DgkProductUrl


	dgkUrlCol.Find(bson.M{"url":DgkUrl }).All(&dgkUrlStr)
	if dgkUrlStr!=nil{
		return 2 ,""
	}else{
		var dgkSku [] dgkcatalogconfig.DgkCatalogConfig
		dgkCatalogConfig := db.C("dgk_catalog_config")
		dgkCatalogConfig.Find(bson.M{"sku_link":DgkUrl}).All(&dgkSku)
		fmt.Println("URL ", DgkUrl)
		if dgkSku==nil {
			var dgkUrlAdd bmldgktable.DgkProductUrl
			dgkUrlAdd.ID = time.Now().Format("20060102150405.9999990700")
			dgkUrlAdd.Url = DgkUrl
			dgkUrlAdd.BmlID = BmlID
			dgkUrlAdd.TypeName = "DgkProductUrl"
			err := dgkUrlCol.Insert(dgkUrlAdd)
			CheckErr(err)
			if err == nil{
				return 0 ,""
			}else {
				return 5 ,""
			}
		}else{
			return 1 , dgkSku[0].SkuName
		}
		
	}

}
