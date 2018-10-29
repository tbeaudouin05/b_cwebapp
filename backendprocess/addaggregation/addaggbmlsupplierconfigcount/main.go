package main

import (
	"log"
	"time"
	"strconv"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/sql/connectdb"
	supplierstr"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/supplier"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract/backendprocess"

)
func main(){
	// Mongo Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	db := mongoSession.DB("competition_analysis")
	bmlAggStatisticHist := db.C("bml_agg_statistic_hist")

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()
	
	bmlSupplierConfigCount := backendprocess.AddAggBmlSupplier(dbBob)
	for i,rngBmlSupplierConfigCount := range bmlSupplierConfigCount{
	
		var supplier supplierstr.SupplierHist
		now := time.Now()
		id,_ :=  strconv.Atoi( strconv.Itoa(rngBmlSupplierConfigCount.SupplierID) + now.Format("20060102") ) 
		supplier.ID =  id
		supplier.SupplierName = rngBmlSupplierConfigCount.SupplierName
		supplier.SupplierNameEn = rngBmlSupplierConfigCount.SupplierNameEn
		supplier.ConfigCount = rngBmlSupplierConfigCount.ConfigCount
		supplier.InStockConfigCount = rngBmlSupplierConfigCount.InStockConfigCount
		supplier.StatisticSnapshotAt = now
		supplier.TypeName = "BmlSupplierConfigCountTopPage"
		err := bmlAggStatisticHist.Insert(supplier)
	
		var checkSupplier [] supplierstr.Supplier
		bmlAggStatisticHist.FindId(rngBmlSupplierConfigCount.SupplierID).All(&checkSupplier)
		if checkSupplier != nil {
			err = bmlAggStatisticHist.Update(bson.M{"_id": rngBmlSupplierConfigCount.SupplierID},bson.M{"$set":bson.M{
				"supplier_name" : rngBmlSupplierConfigCount.SupplierName,
				"supplier_name_en": rngBmlSupplierConfigCount.SupplierNameEn,
				"config_count": rngBmlSupplierConfigCount.ConfigCount,
				"in_stock_config_count": rngBmlSupplierConfigCount.InStockConfigCount,
				"id_supplier":rngBmlSupplierConfigCount.SupplierID,
				"supplier_code":rngBmlSupplierConfigCount.SupplierCode,
				"rowkey": i+1,
				"type":"BmlSupplierConfigCountTopPageLast",
				"statistic_snapshot_at": now,
				 }})
				checkErr(err)
		}else{
			err = bmlAggStatisticHist.Insert(bson.M{"_id": rngBmlSupplierConfigCount.SupplierID,
				"supplier_name" : rngBmlSupplierConfigCount.SupplierName,
				"supplier_name_en": rngBmlSupplierConfigCount.SupplierNameEn,
				"config_count": rngBmlSupplierConfigCount.ConfigCount,
				"in_stock_config_count": rngBmlSupplierConfigCount.InStockConfigCount,
				"id_supplier":rngBmlSupplierConfigCount.SupplierID,
				"supplier_code":rngBmlSupplierConfigCount.SupplierCode,
				"rowkey": i+1,
				"type":"BmlSupplierConfigCountTopPageLast",
				"statistic_snapshot_at": now,
				 })
				checkErr(err)
		}
		
	}



	//save number of bml suppliers in db
	count := backendprocess.NumberOfBmlSupplier(dbBob)
	var supplierCount SupplierCount
	now := time.Now()
	id,_ := strconv.Atoi( "555"+now.Format("20060102"))
	supplierCount.ID = id
	supplierCount.Count = count
	supplierCount.TypeName = "BmlSupplierCount"
	supplierCount.StatisticSnapshotAt = now
	err = bmlAggStatisticHist.Insert(supplierCount)
	//checkErr(err)
	supplierCount.ID = 123456789
	supplierCount.TypeName = "BmlSupplierCountLast"
	_,err = bmlAggStatisticHist.UpsertId(supplierCount.ID ,supplierCount)


}
type SupplierCount struct {
	ID                  int       `bson:"_id"`
	Count               int       `bson:"supplier_count"`
	TypeName            string    `bson:"type"`
	StatisticSnapshotAt time.Time `bson:"statistic_snapshot_at"`
}

func checkErr(err error){
	if err != nil {
		log.Println("Error " ,err)
	}
}