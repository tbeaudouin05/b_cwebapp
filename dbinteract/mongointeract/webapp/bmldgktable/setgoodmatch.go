package bmldgktable

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgkmanualmatching"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
)

func ApplyGoodMatch(IDBmlCatalogConfig int, goodMatched bool, email string, name string, mongoSession *mgo.Session) {

	fmt.Println("check interact")
	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")

	//queries
	var tmpBmlCatalogConfig []bmlcatalogconfig.BmlCatalogConfig

	bmlCatalogConfig.Find(bson.M{"_id": IDBmlCatalogConfig}).All(&tmpBmlCatalogConfig)

	if tmpBmlCatalogConfig != nil {
		err := bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$set": bson.M{"good_match": goodMatched, "matched_by_email": email, "matched_by_name": name, "good_match_at": time.Now()}})
		fmt.Println(err)
		if goodMatched {
			// check to prevent bugs if they match manually and click good match at the same time
			if tmpBmlCatalogConfig[0].DgkScore != 2000000000 {
				err = bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$set": bson.M{"stored_dgk_score": tmpBmlCatalogConfig[0].DgkScore}})
				fmt.Println(err)
			}
			err = bmlCatalogConfig.UpdateId(IDBmlCatalogConfig, bson.M{"$set": bson.M{"dgk_score": 2000000000}})
			fmt.Println(err)
		} else {
			bmldgkmanualmatching.ApplyUnmatch(IDBmlCatalogConfig, mongoSession)
			fmt.Println(err)
		}
	} else {
		fmt.Println("No doc found--goodmatch")
	}

}
