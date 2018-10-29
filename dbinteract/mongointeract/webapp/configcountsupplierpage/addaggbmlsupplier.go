package configcountsupplierpage

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/supplier"
)

func AddAggBmlSupplier() {

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

	var supplierName []string
	bmlCatalogConfig.Find(bson.M{}).Distinct("supplier_name", &supplierName)

	//fmt.Println("supplierName:  ",supplierName) /// works --- alot of suppliers :)

	for _, rngSupplierName := range supplierName {
		num, err := bmlCatalogConfig.Find(bson.M{"supplier_name": rngSupplierName}).Count()
		checkErr(err)
		var tmpSupplier supplier.Supplier
		//tmpSupplier.ID =  string(bson.NewObjectId())
		tmpSupplier.SupplierName = rngSupplierName
		tmpSupplier.ConfigCount = num
		tmpSupplier.StatisticSnapshotAt = time.Now()
		tmpSupplier.TypeName = "BmlSupplierConfigCountTopPage"
		err = bmlAggStatisticHist.Insert(tmpSupplier)
		checkErr(err)

	}

}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error ", err)
	}
}
