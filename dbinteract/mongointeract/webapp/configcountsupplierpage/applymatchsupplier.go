package configcountsupplierpage

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func ApplyMatchSupplier(iDBmlSupplier int, bmlMatchedSupplierName string, fKDgkSupplier string, dgkMatchedSupplierName string, email string, name string, mongoSession *mgo.Session) string {

	if bmlMatchedSupplierName == "" {
		return "Choose bamilo supplier"
	}
	if dgkMatchedSupplierName == "" {
		return "Choose digikala supplier"
	}

	db := mongoSession.DB("competition_analysis")
	bmlAggStatisticHist := db.C("bml_agg_statistic_hist")

	changeInfo, err1 := bmlAggStatisticHist.UpsertId(iDBmlSupplier, bson.M{"$set": bson.M{"fk_dgk_supplier": fKDgkSupplier,
		"dgk_matched_supplier_name": dgkMatchedSupplierName,
		"good_match":                true,
		"matched_by_email":          email,
		"matched_by_name":           name,
		"good_match_at":             time.Now(),
	}})
	CheckErr(err1)
	if err1 != nil {
		return "Error! try again later."
	}
	fmt.Println("changeInfo: ", changeInfo)
	return bmlMatchedSupplierName + " and " + dgkMatchedSupplierName + " are matched."

}
