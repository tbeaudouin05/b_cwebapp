package configcountsupplierpage

import (

	//"github.com/pmylund/sortutil"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/supplier"
)

///
func SelectDgkSupplierTable(searchedBy string, mongoSession *mgo.Session) (selectedResult []supplier.DgkSupplier) {

	//AddAggDgkSupplier()

	db := mongoSession.DB("competition_analysis")
	dgkAggStatisticHist := db.C("dgk_agg_statistic_hist")

	typeName := "DgkSupplierConfigCountTopPageLast"
	// draft for filtering //
	//stageFilterForDgkScore := bson.M{"$match": bson.M{"dgk_score":bson.M{"$gte":0}}}

	//check which queries are required
	if searchedBy == "" {
		dgkAggStatisticHist.Find(bson.M{"type": bson.M{"$eq": typeName}}).Sort("supplier_name").All(&selectedResult)

	} else {
		dgkAggStatisticHist.Find(bson.M{"type": typeName,
			"$text": bson.M{"$search": searchedBy},
		}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Sort("supplier_name").All(&selectedResult)
	}

	//sortutil.DescByField(selectedResult,"ConfigCountInt")

	if selectedResult == nil {
		var nullSupplier supplier.DgkSupplier
		nullSupplier.SupplierName = ""
		nullSupplier.RowKey = 1
		selectedResult = append(selectedResult, nullSupplier)
	}
	return selectedResult

}
