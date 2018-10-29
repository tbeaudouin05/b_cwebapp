package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/sql/connectdb"
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

	db := mongoSession.DB("competition_analysis")

	var bmlAllResult []bmlcatalogconfig.BmlCatalogConfig
	var arrayOfFKDgkCatalogConfig []string
	var arrayOfFKDgkCatalogManualConfig []string
	db.C("bml_catalog_config").Find(bson.M{"good_match": true}).Select(bson.M{
		"id_bml_catalog_config":        17,
		"sku_name":                     17,
		"sku_link":                     17,
		"sku":                          17,
		"supplier_name":                17,
		"brand_name":                   17,
		"department":                   17,
		"bi_category":                  17,
		"bi_category_one_name":         17,
		"bi_category_two_name":         17,
		"bi_category_three_name":       17,
		"bi_category_four_name":        17,
		"bi_category_five_name":        17,
		"bi_category_six_name":         17,
		"avg_price":                    17,
		"avg_special_price":            17,
		"key_account_manager":          17,
		"manual_fk_dgk_catalog_config": 17,
		"fk_dgk_catalog_config":        17,
		"matched_by_email":             17}).All(&bmlAllResult)

	for _, bmlResult := range bmlAllResult {
		arrayOfFKDgkCatalogManualConfig = append(arrayOfFKDgkCatalogManualConfig, bmlResult.ManualFKBmlCatalogConfig)
	}
	var r []string
	for _, str := range arrayOfFKDgkCatalogManualConfig {
		if str != "" {
			r = append(r, str)
		}
	}

	for _, bmlResult := range bmlAllResult {
		arrayOfFKDgkCatalogConfig = append(arrayOfFKDgkCatalogConfig, bmlResult.FKBmlCatalogConfig)
	}

	arrayOfIds := append(arrayOfFKDgkCatalogConfig, r...)

	var dgkAllResult []dgkcatalogconfig.DgkCatalogConfig
	db.C("dgk_catalog_config").Find(bson.M{"id_dgk_catalog_config": bson.M{"$in": arrayOfIds}}).Select(bson.M{
		"id_dgk_catalog_config": 10,
		"avg_price":             10,
		"avg_special_price":     10,
		"dgk_category_one_name": 10,
		"dgk_category_two_name": 10,
		"rating":                10,
		"sku_name":              10,
		"sku_link":              10}).All(&dgkAllResult)

	var arrayOfFKDgkSkuName []string

	for _, bmlCatResult := range dgkAllResult {
		arrayOfFKDgkSkuName = append(arrayOfFKDgkSkuName, bmlCatResult.SkuName)
	}

	var dgkCatLvlResult []dgkcatalogconfig.DgkCatalogConfig
	db.C("dgk_catalog_config").Find(bson.M{"sku_name": bson.M{"$in": arrayOfFKDgkSkuName}, "supplier_name": bson.M{"$exists": false}}).Select(bson.M{
		"id_dgk_catalog_config":         10,
		"avg_price":                     10,
		"avg_special_price":             10,
		"dgk_category_one_name":         10,
		"dgk_category_two_name":         10,
		"rating":                        10,
		"sku_name":                      10,
		"sku_link":                      10,
		"best_price":                    10,
		"best_price_update_snapshot_at": 10,
		"config_snapshot_at":            10,
		"supplier_name":                 10}).All(&dgkCatLvlResult)

	dbSqlite := connectdb.ConnectToSQLite()
	defer dbSqlite.Close()

	db.C("bml_dgk_agg_statistic_hist").RemoveAll(bson.M{"type": "NotCompetitiveGoodMatched"})
	CreateNotCompetitiveTable(mongoSession, dbSqlite, bmlAllResult, dgkAllResult, dgkCatLvlResult)

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)
}
func CreateNotCompetitiveTable(mongoSession *mgo.Session, db *sql.DB, bmlAllResult []bmlcatalogconfig.BmlCatalogConfig, dgkAllResult []dgkcatalogconfig.DgkCatalogConfig, dgkCatLvlResult []dgkcatalogconfig.DgkCatalogConfig) {

	// create bmlAllResult
	query := `CREATE TABLE bmlAllResult (
		id_bml_catalog_config INTEGER,
		sku_name                 TEXT,
		sku_link                 TEXT,       
		sku 		             TEXT,
		supplier_name 	         TEXT,
		brand_name               TEXT,
		department               TEXT,
		bi_category              TEXT,
		bi_category_one_name     TEXT,
		bi_category_two_name     TEXT,
		bi_category_three_name   TEXT,
		bi_category_four_name    TEXT,
		bi_category_five_name    TEXT,
		bi_category_six_name     TEXT,
		avg_price             INTEGER,
		avg_special_price     INTEGER,
		key_account_manager      TEXT,
		manual_fk_dgk_catalog_config    TEXT,
		fk_dgk_catalog_config    TEXT,
		matched_by_email 		 TEXT
		)`
	queryP, err := db.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO bmlAllResult (
		id_bml_catalog_config,
		sku_name,
		sku_link,        
		sku,
		supplier_name, 
		brand_name,
		department,
		bi_category,
		bi_category_one_name,
		bi_category_two_name,
		bi_category_three_name,
		bi_category_four_name,
		bi_category_five_name,
		bi_category_six_name,
		avg_price,
		avg_special_price,
		key_account_manager,
		manual_fk_dgk_catalog_config,
		fk_dgk_catalog_config,
		matched_by_email
		) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,? , ?, ?, ?, ?, ?, ?, ? , ?, ?,?)`
	queryP, err = db.Prepare(query)
	checkError(err)
	for i := 0; i < len(bmlAllResult); i++ {

		queryP.Exec(
			bmlAllResult[i].IDBmlCatalogConfig,
			bmlAllResult[i].SKUName,
			bmlAllResult[i].SKULink,
			bmlAllResult[i].SKU,
			bmlAllResult[i].SupplierName,
			bmlAllResult[i].BrandName,
			bmlAllResult[i].Department,
			bmlAllResult[i].BiCategory,
			bmlAllResult[i].BiCategoryOneName,
			bmlAllResult[i].BiCategoryTwoName,
			bmlAllResult[i].BiCategoryThreeName,
			bmlAllResult[i].BiCategoryFourName,
			bmlAllResult[i].BiCategoryFiveName,
			bmlAllResult[i].BiCategorySixName,
			bmlAllResult[i].AvgPrice,
			bmlAllResult[i].AvgSpecialPrice,
			bmlAllResult[i].KeyAccountManager,
			bmlAllResult[i].ManualFKBmlCatalogConfig,
			bmlAllResult[i].FKBmlCatalogConfig,
			bmlAllResult[i].MatchedByEmail,
		)
		time.Sleep(1 * time.Millisecond)
	}

	query = `
	CREATE VIEW bmlAllResultFk AS
	SELECT 
	ar.id_bml_catalog_config,
	ar.sku_name,
	ar.sku_link,
	ar.sku,
	ar.supplier_name, 
	ar.brand_name,
	ar.department,
	ar.bi_category,
	ar.bi_category_one_name,
	ar.bi_category_two_name,
	ar.bi_category_three_name,
	ar.bi_category_four_name,
	ar.bi_category_five_name,
	ar.bi_category_six_name,
	ar.avg_price,
	ar.avg_special_price,
	ar.key_account_manager,
	CASE WHEN ar.manual_fk_dgk_catalog_config='' THEN ar.fk_dgk_catalog_config ELSE ar.manual_fk_dgk_catalog_config END 'Fk',
	ar.matched_by_email

	FROM bmlAllResult ar
	

	`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// create dgkAllResult
	query = `CREATE TABLE dgkAllResult (
		id_dgk_catalog_config    TEXT,
		avg_price             INTEGER,
		avg_special_price     INTEGER,
		sku_name                 TEXT,
		sku_link                 TEXT
		)`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO dgkAllResult (
		id_dgk_catalog_config ,
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
	// create dgkCatLvlResult
	query = `CREATE TABLE dgkCatLvlResult (
		id_dgk_catalog_config         TEXT,
		sku_name                      TEXT,
		sku_link                      TEXT,       
		best_price_update_snapshot_at TEXT,
		avg_price                  INTEGER,
		best_price                 INTEGER,
		avg_special_price          INTEGER,
		color_of_best_price         String,
        supplier_of_best_price      String,
        warranty_of_best_price      String,
		config_snapshot_at            TEXT
		)`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO dgkCatLvlResult (
		id_dgk_catalog_config,
		sku_name,
		sku_link,        
		best_price_update_snapshot_at,
		avg_price,
		best_price,
		avg_special_price,
		color_of_best_price,
        supplier_of_best_price,
        warranty_of_best_price,
		config_snapshot_at      
		  
		) 
	VALUES (?, ?, ?, ?, ?, ?, ?,?,?,?,? )`
	queryP, err = db.Prepare(query)
	checkError(err)
	for i := 0; i < len(dgkCatLvlResult); i++ {

		queryP.Exec(
			dgkCatLvlResult[i].IDDgkCatalogConfig,
			dgkCatLvlResult[i].SkuName,
			dgkCatLvlResult[i].SkuLink,
			dgkCatLvlResult[i].BestPriceUpdateSnapshotAt,
			dgkCatLvlResult[i].AvgPrice,
			dgkCatLvlResult[i].BestPrice,
			dgkCatLvlResult[i].AvgSpecialPrice,
			dgkCatLvlResult[i].ColorOfBestPrice,
			dgkCatLvlResult[i].SupplierOfBestPrice,
			dgkCatLvlResult[i].WarrantyOfBestPrice,
			dgkCatLvlResult[i].ConfigSnapshotAt,
		)
		time.Sleep(1 * time.Millisecond)
	}

	// join tables
	query = `
	CREATE VIEW allResult AS
	SELECT 
	bar.id_bml_catalog_config                       'id'
	,COALESCE(bar.sku_name,'')                          'bml_sku_name'
	,COALESCE(dar.sku_name,'')                          'dgk_sku_name'
	,COALESCE(bar.sku_link,'')                          'bml_sku_link'
	,COALESCE(dar.sku_link,'')                          'dgk_sku_link'
	,COALESCE(bar.sku,'') 	                                 'bml_sku' 
	,COALESCE(bar.supplier_name,'')                'bml_supplier_name'
	,COALESCE(bar.brand_name,'')                          'brand_name'
	,COALESCE(bar.department,'')                          'department'
	,COALESCE(bar.bi_category,'')                        'bi_category'
	,COALESCE(bar.bi_category_one_name,'')      'bi_category_one_name'
	,COALESCE(bar.bi_category_two_name,'')      'bi_category_two_name'
	,COALESCE(bar.bi_category_three_name,'')  'bi_category_three_name'
	,COALESCE(bar.bi_category_four_name,'')    'bi_category_four_name'
	,COALESCE(bar.bi_category_five_name,'')    'bi_category_five_name'
	,COALESCE(bar.bi_category_six_name,'')      'bi_category_six_name'
	,COALESCE(bar.matched_by_email,'')              'matched_by_email'
	,CASE 
		WHEN  bar.key_account_manager='' 
		OR bar.key_account_manager='vendors' 
		OR bar.key_account_manager='vendors@bamilo.com' 
		OR bar.key_account_manager='vendor@bamilo.com' 
		OR bar.key_account_manager='Vendors@Bamilo.com' 
		OR bar.key_account_manager='Vendors@bamilo.com' 
		OR bar.key_account_manager='vendors@Bamilo.com' THEN 'VM'
		ELSE bar.key_account_manager END 'key_account_manager',
	COALESCE(car.color_of_best_price,'') 'color_of_best_price',
	COALESCE(car.supplier_of_best_price,'') 'supplier_of_best_price',
	COALESCE(car.warranty_of_best_price,'') 'warranty_of_best_price',
	CASE 
		WHEN  dar.avg_special_price <= 0 THEN COALESCE(dar.avg_price,0)
		ELSE COALESCE(dar.avg_special_price,0) END 'dgk_price',
	CASE 
		WHEN  bar.avg_special_price <= 0 THEN COALESCE(bar.avg_price,0)
		ELSE COALESCE(bar.avg_special_price,0) END 'bml_price',
	CASE 
		WHEN  COALESCE(car.best_price,0)<=0 OR  (julianday(car.best_price_update_snapshot_at)- julianday(car.config_snapshot_at)) >=3 THEN CASE WHEN COALESCE(dar.avg_special_price,0)<=0 THEN COALESCE(dar.avg_price,0) ELSE COALESCE(dar.avg_special_price,0) END
		ELSE COALESCE(car.best_price,0) END 'dgk_best_price'
		
	FROM bmlAllResultFk bar
	    LEFT JOIN dgkAllResult dar  ON bar.Fk = dar.id_dgk_catalog_config
		LEFT JOIN dgkCatLvlResult car ON bar.sku_name = car.sku_name
	
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
	ar.brand_name,
	ar.department,
	ar.bi_category,
	ar.bi_category_one_name,
	ar.bi_category_two_name,
	ar.bi_category_three_name,
	ar.bi_category_four_name,
	ar.bi_category_five_name,
	ar.bi_category_six_name,
	ar.key_account_manager,
	ar.color_of_best_price,
	ar.supplier_of_best_price,
	ar.warranty_of_best_price,
	ar.bml_price,
	ar.dgk_price,
	COALESCE(ar.dgk_best_price,ar.dgk_price) 'dgk_best_price',
	'NotCompetitiveGoodMatched' 'type',
	ar.matched_by_email

	FROM allResult ar
	WHERE ar.bml_price > COALESCE(ar.dgk_best_price,ar.dgk_price) AND COALESCE(ar.dgk_best_price,ar.dgk_price)>0

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
	cs.brand_name,
	cs.department,
	cs.bi_category, 
	cs.bi_category_one_name,
	cs.bi_category_two_name,
	cs.bi_category_three_name,
	cs.bi_category_four_name,
	cs.bi_category_five_name,
	cs.bi_category_six_name,
	cs.key_account_manager,
	cs.color_of_best_price,
	cs.supplier_of_best_price,
	cs.warranty_of_best_price,
	cs.bml_price,
	cs.dgk_price,
	cs.dgk_best_price,
	cs.type,
	cs.matched_by_email

	FROM notCompetitiveSKU cs`

	rows, err := db.Query(query)
	checkError(err)

	mongoDB := mongoSession.DB("competition_analysis")

	var notCompetitiveMatchedSKU NotCompetitiveMatchedSKU
	for rows.Next() {
		err := rows.Scan(&notCompetitiveMatchedSKU.ID,
			&notCompetitiveMatchedSKU.BmlSkuName,
			&notCompetitiveMatchedSKU.DgkSkuName,
			&notCompetitiveMatchedSKU.BmlSkuLink,
			&notCompetitiveMatchedSKU.DgkSkuLink,
			&notCompetitiveMatchedSKU.BmlSku,
			&notCompetitiveMatchedSKU.BmlSupplierName,
			&notCompetitiveMatchedSKU.BrandName,
			&notCompetitiveMatchedSKU.Department,
			&notCompetitiveMatchedSKU.BiCategory,
			&notCompetitiveMatchedSKU.BiCategoryOneName,
			&notCompetitiveMatchedSKU.BiCategoryTwoName,
			&notCompetitiveMatchedSKU.BiCategoryThreeName,
			&notCompetitiveMatchedSKU.BiCategoryFourName,
			&notCompetitiveMatchedSKU.BiCategoryFiveName,
			&notCompetitiveMatchedSKU.BiCategorySixName,
			&notCompetitiveMatchedSKU.KeyAccountManager,
			&notCompetitiveMatchedSKU.ColorOfBestPrice,
			&notCompetitiveMatchedSKU.SupplierOfBestPrice,
			&notCompetitiveMatchedSKU.WarrantyOfBestPrice,
			&notCompetitiveMatchedSKU.BmlPrice,
			&notCompetitiveMatchedSKU.DgkPrice,
			&notCompetitiveMatchedSKU.DgkBestPrice,
			&notCompetitiveMatchedSKU.Type,
			&notCompetitiveMatchedSKU.MatchedByEmail,
		)
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
	ID                  int    `bson:"_id"`
	BmlSkuName          string `bson:"bml_sku_name"`
	DgkSkuName          string `bson:"dgk_sku_name"`
	BmlSkuLink          string `bson:"bml_sku_link"`
	DgkSkuLink          string `bson:"dgk_sku_link"`
	BmlSku              string `bson:"bml_sku"`
	BmlSupplierName     string `bson:"bml_supplier_name"`
	BrandName           string `bson:"brand_name"`
	Department          string `bson:"department"`
	BiCategory          string `bson:"bi_category"`
	BiCategoryOneName   string `bson:"bi_category_one_name"`
	BiCategoryTwoName   string `bson:"bi_category_two_name"`
	BiCategoryThreeName string `bson:"bi_category_three_name"`
	BiCategoryFourName  string `bson:"bi_category_four_name"`
	BiCategoryFiveName  string `bson:"bi_category_five_name"`
	BiCategorySixName   string `bson:"bi_category_six_name"`
	BmlPrice            int    `bson:"bml_price"`
	DgkPrice            int    `bson:"dgk_price"`
	ColorOfBestPrice    string `bson:"color_of_best_price"`
	SupplierOfBestPrice string `bson:"supplier_of_best_price"`
	WarrantyOfBestPrice string `bson:"warranty_of_best_price"`
	DgkBestPrice        int    `bson:"dgk_best_price"`
	Type                string `bson:"type"`
	KeyAccountManager   string `bson:"key_account_manager"`
	MatchedByEmail      string `bson:"matched_by_email"`
}
