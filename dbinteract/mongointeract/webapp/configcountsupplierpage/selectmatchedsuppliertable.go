package configcountsupplierpage

import (

	//"github.com/pmylund/sortutil"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/supplier"
)

///
func SelectMatchedSupplierTable(BmlSearchedBy string, DgkSearchedBy string, mongoSession *mgo.Session) (selectedResult []supplier.BmlDgkSupplier) {

	db := mongoSession.DB("competition_analysis")
	dgkAggStatisticHist := db.C("bml_agg_statistic_hist")

	typeName := "BmlSupplierConfigCountTopPageLast"
	// draft for filtering //
	//stageFilterForDgkScore := bson.M{"$match": bson.M{"dgk_score":bson.M{"$gte":0}}}

	//check which queries are required
	if BmlSearchedBy == "" && DgkSearchedBy == "" {
		dgkAggStatisticHist.Find(bson.M{"type": bson.M{"$eq": typeName}, "good_match": true}).Sort("supplier_name").All(&selectedResult)

	} else if BmlSearchedBy != "" && DgkSearchedBy == "" {
		dgkAggStatisticHist.Find(bson.M{"type": typeName, "good_match": true,
			"$text": bson.M{"$search": BmlSearchedBy},
		}).Select(bson.M{"score": bson.M{"$meta": "textScore"}, "good_match": true}).Sort("$textScore:score").Sort("supplier_name").All(&selectedResult)
	} else if BmlSearchedBy == "" && DgkSearchedBy != "" {
		dgkAggStatisticHist.Find(bson.M{"type": typeName,
			"$text": bson.M{"$search": DgkSearchedBy},
		}).Select(bson.M{"score": bson.M{"$meta": "textScore"}, "good_match": true}).Sort("$textScore:score").Sort("supplier_name").All(&selectedResult)
	}

	//sortutil.DescByField(selectedResult,"ConfigCountInt")

	if selectedResult == nil {
		var nullSupplier supplier.BmlDgkSupplier
		nullSupplier.BmlSupplierName = ""
		nullSupplier.DgkSupplierName = ""
		selectedResult = append(selectedResult, nullSupplier)
	}
	return selectedResult

}
