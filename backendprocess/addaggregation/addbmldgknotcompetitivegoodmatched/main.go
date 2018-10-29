package main

import (
	
	"time"
	"log"
	"database/sql"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/sql/connectdb"
) 

func main(){
	start := time.Now()
	log.Println(`Start: ` + start.Format(`1 January 2006, 15:04:05`))
	
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	db := mongoSession.DB("competition_analysis")
	
	var bmlAllResult []bmlcatalogconfig.BmlCatalogConfig
	var arrayOfFKDgkCatalogConfig []string

	db.C("bml_catalog_config").Find(bson.M{"good_match": true}).Select(bson.M{
		"id_bml_catalog_config":  17,
		"sku_name":               17,
		"sku_link":               17,
		"sku":17,
		"supplier_name" : 17,
		"bi_category_one_name":   17,
		"bi_category_two_name":   17,
		"bi_category_three_name": 17,
		"avg_price":              17,
		"avg_special_price":      17,
		`manual_fk_dgk_catalog_config`:  17}).All(&bmlAllResult)

	for _, bmlResult := range bmlAllResult {
		arrayOfFKDgkCatalogConfig = append(arrayOfFKDgkCatalogConfig, bmlResult.ManualFKBmlCatalogConfig)
	}

	var dgkAllResult []dgkcatalogconfig.DgkCatalogConfig
	db.C("dgk_catalog_config").Find(bson.M{"_id": bson.M{"$in": arrayOfFKDgkCatalogConfig}}).Select(bson.M{
		"id_dgk_catalog_config": 10,
		"avg_price":             10,
		"avg_special_price":     10,
		"dgk_category_one_name": 10,
		"dgk_category_two_name": 10,
		"rating":                10,
		"sku_name":              10,
		"sku_link":              10}).All(&dgkAllResult)

	dbSqlite := connectdb.ConnectToSQLite()
	defer dbSqlite.Close()

	db.C("bml_dgk_agg_statistic_hist").RemoveAll(bson.M{"type":"NotCompetitiveGoodMatched"})

	GetNotCompetitiveTable(mongoSession, dbSqlite, bmlAllResult, dgkAllResult)

	
	
}

func GetNotCompetitiveTable(mongoSession *mgo.Session, db *sql.DB, bmlAllResult []bmlcatalogconfig.BmlCatalogConfig, dgkAllResult []dgkcatalogconfig.DgkCatalogConfig)  {
	
	// create bmlAllResult
	query := `CREATE TABLE bmlAllResult (
		id_bml_catalog_config INTEGER,
		sku_name TEXT,
		sku_link     TEXT,       
		sku 			TEXT,
		supplier_name 	TEXT, 
		bi_category_one_name   TEXT,
		bi_category_two_name   TEXT,
		bi_category_three_name TEXT,
		avg_price              INTEGER,
		avg_special_price      INTEGER,
		fk_dgk_catalog_config  TEXT
		)`
	queryP, err := db.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO bmlAllResult (
		id_bml_catalog_config,
		sku_name,
		sku_link,        
		sku 			,
		supplier_name 	, 
		bi_category_one_name,
		bi_category_two_name   ,
		bi_category_three_name ,
		avg_price              ,
		avg_special_price      ,
		fk_dgk_catalog_config  
		) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,? ,? )`
	queryP, err = db.Prepare(query)
	checkError(err)
	for i := 0; i < len(bmlAllResult); i++ {

		queryP.Exec(
			bmlAllResult[i].IDBmlCatalogConfig,
			bmlAllResult[i].SKUName,
			bmlAllResult[i].SKULink,
			bmlAllResult[i].SKU,
			bmlAllResult[i].SupplierName,
			bmlAllResult[i].BiCategoryOneName,
			bmlAllResult[i].BiCategoryTwoName,
			bmlAllResult[i].BiCategoryThreeName,
			bmlAllResult[i].AvgPrice,
			bmlAllResult[i].AvgSpecialPrice,
			bmlAllResult[i].ManualFKBmlCatalogConfig,
		)
		time.Sleep(1 * time.Millisecond)
	}

	// create dgkAllResult
	query = `CREATE TABLE dgkAllResult (
		fk_dgk_catalog_config TEXT,
		avg_price             INTEGER,
		avg_special_price     INTEGER,
		sku_name              TEXT,
		sku_link              TEXT
		)`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO dgkAllResult (
		fk_dgk_catalog_config ,
		avg_price             ,
		avg_special_price     ,
		sku_name              ,
		sku_link               
		) 
	VALUES (?, ?, ?, ?, ?)`
	queryP, err = db.Prepare(query)
	checkError(err)
	for i := 0; i < len(dgkAllResult); i++ {

		queryP.Exec(
			dgkAllResult[i].IDDgkCatalogConfig,
			dgkAllResult[i].AvgPrice,
			dgkAllResult[i].AvgSpecialPrice,
			dgkAllResult[i].SkuName,
			dgkAllResult[i].SkuLink,
		)

		time.Sleep(1 * time.Millisecond)

	}

	// join tables
	query = `
	CREATE VIEW allResult AS
	SELECT 
	bar.id_bml_catalog_config 'id',
	bar.sku_name 'bml_sku_name',
	dar.sku_name 'dgk_sku_name',
	bar.sku_link 'bml_sku_link',
	dar.sku_link 'dgk_sku_link',
	bar.sku 	  'bml_sku' ,
	bar.supplier_name 	'bml_supplier_name', 
	bar.bi_category_one_name,
	bar.bi_category_two_name,
	bar.bi_category_three_name,
	CASE 
		WHEN  bar.avg_special_price <= 0 THEN bar.avg_price
		ELSE bar.avg_special_price END 'bml_price',
	CASE 
		WHEN  dar.avg_special_price <= 0 THEN dar.avg_price
		ELSE dar.avg_special_price END 'dgk_price'
	FROM bmlAllResult bar
	JOIN dgkAllResult dar
	USING(fk_dgk_catalog_config)
	`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// categoryOneStat
	query = `
	CREATE VIEW notCompetitiveSKU AS
	SELECT 
	ar.id,
	ar.bml_sku_name,
	ar.dgk_sku_name,
	ar.bml_sku_link,
	ar.dgk_sku_link,
	ar.bml_sku ,
	ar.bml_supplier_name, 
	ar.bi_category_one_name,
	ar.bi_category_two_name,
	ar.bi_category_three_name,
	ar.bml_price,
	ar.dgk_price,
	'NotCompetitiveGoodMatched' 'type'

	FROM allResult ar
	WHERE ar.bml_price > ar.dgk_price

	`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()


	// return notCompetitiveSKU
	query = `
	SELECT
	cs.id,
	cs.bml_sku_name,
	cs.dgk_sku_name,
	cs.bml_sku_link,
	cs.dgk_sku_link,
	cs.bml_sku ,
	cs.bml_supplier_name, 
	cs.bi_category_one_name,
	cs.bi_category_two_name,
	cs.bi_category_three_name,
	cs.bml_price,
	cs.dgk_price,
	cs.type

	FROM notCompetitiveSKU cs`

	rows, err := db.Query(query)
	checkError(err)

	mongoDB := mongoSession.DB("competition_analysis")

	var notCompetitiveMatchedSKU NotCompetitiveMatchedSKU
	for rows.Next() {
		err := rows.Scan(&notCompetitiveMatchedSKU.ID ,&notCompetitiveMatchedSKU.BmlSkuName ,&notCompetitiveMatchedSKU.DgkSkuName ,&notCompetitiveMatchedSKU.BmlSkuLink ,&notCompetitiveMatchedSKU.DgkSkuLink,&notCompetitiveMatchedSKU.BmlSku ,&notCompetitiveMatchedSKU.BmlSupplierName ,&notCompetitiveMatchedSKU.BiCategoryOneName,&notCompetitiveMatchedSKU.BiCategoryTwoName,&notCompetitiveMatchedSKU.BiCategoryThreeName, &notCompetitiveMatchedSKU.BmlPrice, &notCompetitiveMatchedSKU.DgkPrice, &notCompetitiveMatchedSKU.Type)
		checkError(err)
		mongoDB.C("bml_dgk_agg_statistic_hist").UpsertId(notCompetitiveMatchedSKU.ID, notCompetitiveMatchedSKU)
	}

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type NotCompetitiveMatchedSKU struct {
	ID int `bson:"_id"`
	BmlSkuName string `bson:"bml_sku_name"`
	DgkSkuName string `bson:"dgk_sku_name"`
	BmlSkuLink string `bson:"bml_sku_link"`
	DgkSkuLink string `bson:"dgk_sku_link"`
	BmlSku string `bson:"bml_sku"`
	BmlSupplierName string`bson:"bml_supplier_name"`
	BiCategoryOneName string `bson:"bi_category_one_name"`
	BiCategoryTwoName string `bson:"bi_category_two_name"`
	BiCategoryThreeName string `bson:"bi_category_three_name"`
	BmlPrice int `bson:"bml_price"`
	DgkPrice int `bson:"dgk_price"`

	Type              string `bson:"type"`
}
