package backendprocess

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	bmlconfigbackendprocess "github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig/backendprocess"
	mongobulk "github.com/thomas-bamilo/nosql/mongobulk"
)

func SetMatchedToFalse(mongoSession *mgo.Session) {

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config")

	c.UpdateAll(bson.M{}, bson.M{"$set": bson.M{"matched": false}})

}

func GetBmlCatalogConfigTableToMatch(mongoSession *mgo.Session, BmlCatalogConfigToMatch *[]bmlcatalogconfig.BmlCatalogConfig) {

	bmlCatalogConfigC := mongoSession.DB("competition_analysis").C("bml_catalog_config")

	bmlCatalogConfigC.Find(bson.M{"$and": []bson.M{
		bson.M{"good_match": bson.M{"$ne": true}},
		bson.M{"matched": bson.M{"$ne": true}},
	}}).All(BmlCatalogConfigToMatch)

}

func UpsertConfigMatch(mongoSession *mgo.Session, bamiloCatalogConfigMatchTable []bmlconfigbackendprocess.Match, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config")

	config := mongobulk.Config{OpsPerBatch: 950}

	mongoBulk := mongobulk.New(c, config)

	for _, bamiloCatalogConfigMatch := range bamiloCatalogConfigMatchTable {

		upsertConfigMatch(&bamiloCatalogConfigMatch, mongoBulk)

	}

	err := mongoBulk.Finish()
	checkError(err)

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)

}

func upsertConfigMatch(bamiloCatalogConfigMatch *bmlconfigbackendprocess.Match, mongoBulk *mongobulk.Bulk) {

	bamiloCatalogConfigMatchByte, err := bson.Marshal(bamiloCatalogConfigMatch)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())
	}
	bamiloCatalogConfigMatchBson := make(map[string]interface{})
	err = bson.Unmarshal(bamiloCatalogConfigMatchByte, bamiloCatalogConfigMatchBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())
	}
	colQuerier := bson.M{"_id": bamiloCatalogConfigMatch.IDBmlCatalogConfig}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": bamiloCatalogConfigMatchBson})
}

// sales ------------------------------------------------------------------------------

func UpsertConfigSales(mongoSession *mgo.Session, bamiloCatalogConfigSalesTable []bmlcatalogconfig.BmlCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	now := time.Now()

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config")

	config := mongobulk.Config{OpsPerBatch: 950}

	mongoBulk := mongobulk.New(c, config)

	for _, bamiloCatalogConfigSales := range bamiloCatalogConfigSalesTable {

		// only keep appropriate information for bmlcatalogconfig
		bamiloCatalogConfigSales := bmlconfigbackendprocess.MongoSales{
			IDBmlCatalogConfig: bamiloCatalogConfigSales.IDBmlCatalogConfig,
			ConfigSnapshotAt:   now,

			CountOfSoi:            bamiloCatalogConfigSales.CountOfSoi,
			SumOfUnitPrice:        bamiloCatalogConfigSales.SumOfUnitPrice,
			SumOfPaidPrice:        bamiloCatalogConfigSales.SumOfPaidPrice,
			SumOfCouponMoneyValue: bamiloCatalogConfigSales.SumOfCouponMoneyValue,
			SumOfCartRuleDiscount: bamiloCatalogConfigSales.SumOfCartRuleDiscount,
		}

		upsertConfigSales(&bamiloCatalogConfigSales, mongoBulk)

	}

	err := mongoBulk.Finish()
	checkError(err)

	end := time.Now()
	log.Println(`End time config sales Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config sales Mongo: `, duration.Minutes(), ` minutes`)

}

func upsertConfigSales(bamiloCatalogConfigSales *bmlconfigbackendprocess.MongoSales, mongoBulk *mongobulk.Bulk) {

	bamiloCatalogConfigSalesByte, err := bson.Marshal(bamiloCatalogConfigSales)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())

	}
	bamiloCatalogConfigSalesBson := make(map[string]interface{})
	err = bson.Unmarshal(bamiloCatalogConfigSalesByte, bamiloCatalogConfigSalesBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())

	}
	colQuerier := bson.M{"_id": bamiloCatalogConfigSales.IDBmlCatalogConfig}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": bamiloCatalogConfigSalesBson})

}

func UpsertConfigSalesHist(mongoSession *mgo.Session, bamiloCatalogConfigSalesHistTable []bmlcatalogconfig.BmlCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config_hist")

	config := mongobulk.Config{OpsPerBatch: 950}

	mongoBulk := mongobulk.New(c, config)

	for _, bamiloCatalogConfigSalesHist := range bamiloCatalogConfigSalesHistTable {

		iDBmlCatalogConfigHist, err := strconv.Atoi(strconv.Itoa(bamiloCatalogConfigSalesHist.IDBmlCatalogConfig) + bamiloCatalogConfigSalesHist.ConfigSnapshotAt.Format(`01022006`))
		checkError(err)
		// only keep appropriate information for bmlcatalogconfig
		bamiloCatalogConfigSalesHist := bmlconfigbackendprocess.MongoSalesHist{
			IDBmlCatalogConfigHist: iDBmlCatalogConfigHist,
			FKBmlCatalogConfig:     bamiloCatalogConfigSalesHist.IDBmlCatalogConfig,
			ConfigSnapshotAt:       bamiloCatalogConfigSalesHist.ConfigSnapshotAt,

			CountOfSoi:            bamiloCatalogConfigSalesHist.CountOfSoi,
			SumOfUnitPrice:        bamiloCatalogConfigSalesHist.SumOfUnitPrice,
			SumOfPaidPrice:        bamiloCatalogConfigSalesHist.SumOfPaidPrice,
			SumOfCouponMoneyValue: bamiloCatalogConfigSalesHist.SumOfCouponMoneyValue,
			SumOfCartRuleDiscount: bamiloCatalogConfigSalesHist.SumOfCartRuleDiscount,
		}

		upsertConfigSalesHist(&bamiloCatalogConfigSalesHist, mongoBulk)

	}

	err := mongoBulk.Finish()
	checkError(err)

	end := time.Now()
	log.Println(`End time config sales Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config sales Mongo: `, duration.Minutes(), ` minutes`)

}

func upsertConfigSalesHist(bamiloCatalogConfigSalesHist *bmlconfigbackendprocess.MongoSalesHist, mongoBulk *mongobulk.Bulk) {

	bamiloCatalogConfigSalesHistByte, err := bson.Marshal(bamiloCatalogConfigSalesHist)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())

	}
	bamiloCatalogConfigSalesHistBson := make(map[string]interface{})
	err = bson.Unmarshal(bamiloCatalogConfigSalesHistByte, bamiloCatalogConfigSalesHistBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())

	}
	colQuerier := bson.M{"_id": bamiloCatalogConfigSalesHist.IDBmlCatalogConfigHist}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": bamiloCatalogConfigSalesHistBson})

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
