package bmldgkmanualmatching

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
)

func ApplyManualMatching(iDBmlCatalogConfig int, fKDgkCatalogConfig string, email string, name string, mongoSession *mgo.Session) {

	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")
	//queries

	var tmpBmlCatalogConfig []bmlcatalogconfig.BmlCatalogConfig
	bmlCatalogConfig.Find(bson.M{"_id": iDBmlCatalogConfig}).All(&tmpBmlCatalogConfig)

	err := bmlCatalogConfig.UpdateId(iDBmlCatalogConfig, bson.M{"$set": bson.M{"manual_fk_dgk_catalog_config": fKDgkCatalogConfig, "matched_by_email": email, "matched_by_name": name, "good_match_at": time.Now()}})
	CheckErr(err)
	if tmpBmlCatalogConfig != nil {
		// check to prevent bugs if they match manually and click good match at the same time
		if tmpBmlCatalogConfig[0].DgkScore != 2000000000 {
			err = bmlCatalogConfig.UpdateId(iDBmlCatalogConfig, bson.M{"$set": bson.M{"stored_dgk_score": tmpBmlCatalogConfig[0].DgkScore}})
			fmt.Println(err)
		}
	} else {
		fmt.Println("No doc found--manual match")
	}
	CheckErr(err)
	err = bmlCatalogConfig.UpdateId(iDBmlCatalogConfig, bson.M{"$set": bson.M{"dgk_score": 2000000000}})
	CheckErr(err)
	err = bmlCatalogConfig.UpdateId(iDBmlCatalogConfig, bson.M{"$set": bson.M{"good_match": true}})
	CheckErr(err)
}

func UpdateFrequency(IDBmlCatalogConfig int, checkCycle int, mongoSession *mgo.Session) bool {
	fmt.Println(checkCycle, IDBmlCatalogConfig)
	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")

	var tmpBmlCatalogConfig []bmlcatalogconfig.BmlCatalogConfig
	bmlCatalogConfig.Find(bson.M{"good_match": bson.M{"$exists": true}, "_id": IDBmlCatalogConfig}).All(&tmpBmlCatalogConfig)
	if tmpBmlCatalogConfig != nil && tmpBmlCatalogConfig[0].GoodMatch == true {
		err := bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$set": bson.M{"update_frequency": checkCycle}})
		CheckErr(err)
		return true
	} else {
		return false
	}

}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
