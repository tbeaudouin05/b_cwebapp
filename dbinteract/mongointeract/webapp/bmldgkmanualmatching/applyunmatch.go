package bmldgkmanualmatching

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"

)

func ApplyUnmatch(IDBmlCatalogConfig int, mongoSession *mgo.Session) string {

	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")
	dgkCatalogConfig := db.C("dgk_catalog_config")

	var dgkImgLink string
	//queries
	var tmpBmlCatalogConfig [] bmlcatalogconfig.BmlCatalogConfig
	var tmpDgkCatalogConfig [] bmlcatalogconfig.BmlCatalogConfig

	bmlCatalogConfig.Find(bson.M{"_id":IDBmlCatalogConfig }).All(&tmpBmlCatalogConfig)
	
	err := bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$unset": bson.M{"manual_fk_dgk_catalog_config":"" }})
	CheckErr(err)
	err = bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$set": bson.M{"good_match": false}})
	CheckErr(err)
	
	if tmpBmlCatalogConfig != nil{
		err = bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$set": bson.M{"dgk_score": tmpBmlCatalogConfig[0].StoredDgkScore}})
		dgkCatalogConfig.Find(bson.M{"_id":tmpBmlCatalogConfig[0].FKBmlCatalogConfig }).All(&tmpDgkCatalogConfig)
			//bmlCatalogConfig.Find(bson.M{"_id":tmpBmlCatalogConfig[0].FKBmlCatalogConfig }).All(&tmpDgkCatalogConfig)

	}else {
		fmt.Println("No doc found--unmatch")
	}
	CheckErr(err)



	if tmpBmlCatalogConfig != nil{
		dgkImgLink = tmpDgkCatalogConfig[0].ImgLink
	}else {
		fmt.Println("No doc found--unmatch--dgk image")
	}

	return dgkImgLink
}

