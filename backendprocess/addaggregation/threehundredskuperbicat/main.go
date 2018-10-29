package main

import (
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func main() {
	start := time.Now()
	log.Println(`Start: ` + start.Format(`1 January 2006, 15:04:05`))
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	mongoSession.SetMode(mgo.Monotonic, true)

	mongoDB := mongoSession.DB("competition_analysis")

	db := mongoSession.DB("competition_analysis")

	db.C("top_300_bml_sku").RemoveAll(nil)
	var bmlAllResult []BmlAllResult

	var biCats []string
	db.C("bml_catalog_config").Find(nil).Distinct(`bi_category`, &biCats)
	fmt.Println(biCats)

	for j := 0; j < len(biCats); j++ {

		db.C("bml_catalog_config").Find(bson.M{"bi_category": biCats[j]}).Sort("-count_of_soi").Limit(300).Select(bson.M{
			"id_bml_catalog_config": 17,
			"bi_category":           17,
			"sku_name":              17,
			"count_of_soi":          17,
			"supplier_name":         17,
			"good_match":            17,
			"sku":                   17,
			"brand_name":            17}).All(&bmlAllResult)

		fmt.Println(j)

		for i := 0; i < len(bmlAllResult); i++ {

			mongoDB.C("top_300_bml_sku").UpsertId(bmlAllResult[i].IDBmlCatalogConfig, bmlAllResult[i])
		}

	}

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type BiCatString struct {
	BiCatString string `bson:"bi_category"`
}

type BmlAllResult struct {
	IDBmlCatalogConfig int    `bson:"id_bml_catalog_config"`
	BiCategory         string `bson:"bi_category"`
	SKUName            string `bson:"sku_name"`
	CountOfSoi         int    `bson:"count_of_soi"`
	SupplierName       string `bson:"supplier_name"`
	BrandName          string `bson:"brand_name"`
	GoodMatch          bool   `bson:"good_match"`
	Sku                string `bson:"sku"`
}
