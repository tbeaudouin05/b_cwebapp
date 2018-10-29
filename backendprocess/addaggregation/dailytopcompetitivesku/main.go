package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/nosql/mongobulk"
	"github.com/thomas-bamilo/sql/connectdb"
)

func main() {

	start := time.Now()
	log.Println(`Start: ` + start.Format(`1 January 2006, 15:04:05`))

	// Connection URL
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
		"sku_name":               17,
		"sku_link":               17,
		"bi_category_one_name":   17,
		"bi_category_two_name":   17,
		"bi_category_three_name": 17,
		"avg_price":              17,
		"avg_special_price":      17,
		`fk_dgk_catalog_config`:  17}).All(&bmlAllResult)

	for _, bmlResult := range bmlAllResult {
		arrayOfFKDgkCatalogConfig = append(arrayOfFKDgkCatalogConfig, bmlResult.FKBmlCatalogConfig)
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

	categoryOneStatTable, categoryTwoStatTable, categoryThreeStatTable, allStat := GetStatTable(dbSqlite, bmlAllResult, dgkAllResult)

	var wg sync.WaitGroup
	wg.Add(1)
	go upsertStat(mongoSession, categoryOneStatTable, categoryTwoStatTable, categoryThreeStatTable, allStat, start, &wg)
	wg.Wait()

}

func GetStatTable(db *sql.DB, bmlAllResult []bmlcatalogconfig.BmlCatalogConfig, dgkAllResult []dgkcatalogconfig.DgkCatalogConfig) (categoryOneStatTable []CategoryOneStat, categoryTwoStatTable []CategoryTwoStat, categoryThreeStatTable []CategoryThreeStat, allStat AllStat) {

	// create bmlAllResult
	query := `CREATE TABLE bmlAllResult (
		sku_name TEXT,
		sku_link     TEXT,        
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
		sku_name,
		sku_link,        
		bi_category_one_name,
		bi_category_two_name   ,
		bi_category_three_name ,
		avg_price              ,
		avg_special_price      ,
		fk_dgk_catalog_config  
		) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	queryP, err = db.Prepare(query)
	checkError(err)
	for i := 0; i < len(bmlAllResult); i++ {

		queryP.Exec(
			bmlAllResult[i].SKUName,
			bmlAllResult[i].SKULink,
			bmlAllResult[i].BiCategoryOneName,
			bmlAllResult[i].BiCategoryTwoName,
			bmlAllResult[i].BiCategoryThreeName,
			bmlAllResult[i].AvgPrice,
			bmlAllResult[i].AvgSpecialPrice,
			bmlAllResult[i].FKBmlCatalogConfig,
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
	bar.sku_name 'bml_sku_name',
	dar.sku_name 'dgk_sku_name',
	bar.sku_link 'bml_sku_link',
	dar.sku_link 'dgk_sku_link',
	bar.bi_category_one_name,
	bar.bi_category_two_name,
	bar.bi_category_three_name,
	CASE 
		WHEN  bar.avg_special_price <= 0 THEN bar.avg_price
		ELSE bar.avg_special_price END 'bml_price',
	CASE 
		WHEN  dar.avg_special_price <= 0 THEN dar.avg_price
		ELSE dar.avg_special_price END 'dgk_price'
	FROM bmlAllResult  bar
	JOIN dgkAllResult dar
	USING(fk_dgk_catalog_config)
	`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// categoryOneStat
	query = `
	CREATE VIEW categoryOneStat AS
	SELECT 
	ar.bi_category_one_name,
	SUM(CASE
		WHEN ar.bml_price < ar.dgk_price THEN 1
		ELSE 0 END) * 100 / COUNT(*) 'competitiveness',
	COUNT(*) good_match_count,
	'CategoryOneCompetitiveAnalysis' 'type'
	FROM allResult  ar
	GROUP BY ar.bi_category_one_name
	`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// categoryTwoStat
	query = `
		CREATE VIEW categoryTwoStat AS
		SELECT 
		ar.bi_category_one_name,
		ar.bi_category_two_name,
		SUM(CASE
			WHEN ar.bml_price < ar.dgk_price THEN 1
			ELSE 0 END) * 100 / COUNT(*) 'competitiveness',
		COUNT(*) good_match_count,
		'CategoryTwoCompetitiveAnalysis' 'type'
		FROM allResult  ar
		GROUP BY ar.bi_category_two_name
		`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// categoryThreeStat
	query = `
		CREATE VIEW categoryThreeStat AS
		SELECT 
		ar.bi_category_one_name,
		ar.bi_category_two_name,
		ar.bi_category_three_name,
		SUM(CASE
			WHEN ar.bml_price < ar.dgk_price THEN 1
			ELSE 0 END) * 100 / COUNT(*) 'competitiveness',
		COUNT(*) good_match_count,
		'CategoryThreeCompetitiveAnalysis' 'type'
		FROM allResult  ar
		GROUP BY ar.bi_category_three_name
		`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// allStat
	query = `
		CREATE VIEW allStat AS
		SELECT 
		SUM(CASE
			WHEN ar.bml_price < ar.dgk_price THEN 1
			ELSE 0 END) * 100 / COUNT(*) 'competitiveness',
		COUNT(*) good_match_count,
		'AllCompetitiveAnalysis' 'type'
		FROM allResult  ar
		`
	queryP, err = db.Prepare(query)
	checkError(err)
	queryP.Exec()

	// return categoryOneStatTable
	query = `
	SELECT
	cs.bi_category_one_name,
	cs.competitiveness,
	cs.good_match_count,
	cs.type
	FROM categoryOneStat cs`

	rows, err := db.Query(query)
	checkError(err)

	var CategoryOneStat1 CategoryOneStat
	for rows.Next() {
		err := rows.Scan(&CategoryOneStat1.BiCategoryOneName, &CategoryOneStat1.Competitiveness, &CategoryOneStat1.GoodMatchCount, &CategoryOneStat1.Type)
		checkError(err)
		categoryOneStatTable = append(categoryOneStatTable, CategoryOneStat1)
	}

	// return categoryTwoStatTable
	query = `
		SELECT
		cs.bi_category_one_name,
		cs.bi_category_two_name,
		cs.competitiveness,
		cs.good_match_count,
		cs.type
		FROM categoryTwoStat cs`

	rows, err = db.Query(query)
	checkError(err)

	var CategoryTwoStat1 CategoryTwoStat
	for rows.Next() {
		err := rows.Scan(&CategoryTwoStat1.BiCategoryOneName, &CategoryTwoStat1.BiCategoryTwoName, &CategoryTwoStat1.Competitiveness, &CategoryTwoStat1.GoodMatchCount, &CategoryTwoStat1.Type)
		checkError(err)
		categoryTwoStatTable = append(categoryTwoStatTable, CategoryTwoStat1)
	}

	// return categoryThreeStatTable
	query = `
		SELECT
		cs.bi_category_one_name,
		cs.bi_category_two_name,
		cs.bi_category_Three_name,
		cs.competitiveness,
		cs.good_match_count,
		cs.type
		FROM categoryThreeStat cs`

	rows, err = db.Query(query)
	checkError(err)

	var CategoryThreeStat1 CategoryThreeStat
	for rows.Next() {
		err := rows.Scan(&CategoryThreeStat1.BiCategoryOneName, &CategoryThreeStat1.BiCategoryTwoName, &CategoryThreeStat1.BiCategoryThreeName, &CategoryThreeStat1.Competitiveness, &CategoryThreeStat1.GoodMatchCount, &CategoryThreeStat1.Type)
		checkError(err)
		categoryThreeStatTable = append(categoryThreeStatTable, CategoryThreeStat1)
	}

	// return categoryThreeStatTable
	query = `
		SELECT
		cs.competitiveness,
		cs.good_match_count,
		cs.type
		FROM allStat cs`

	row := db.QueryRow(query)
	checkError(err)

	err = row.Scan(&allStat.Competitiveness, &allStat.GoodMatchCount, &allStat.Type)
	checkError(err)

	return categoryOneStatTable, categoryTwoStatTable, categoryThreeStatTable, allStat
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type CategoryOneStat struct {
	BiCategoryOneName string `bson:"bi_category_one_name"`
	Competitiveness   int    `bson:"competitiveness"`
	GoodMatchCount    int    `bson:"good_match_count"`
	Type              string `bson:"type"`
}

type CategoryTwoStat struct {
	BiCategoryOneName string `bson:"bi_category_one_name"`
	BiCategoryTwoName string `bson:"bi_category_two_name"`
	Competitiveness   int    `bson:"competitiveness"`
	GoodMatchCount    int    `bson:"good_match_count"`
	Type              string `bson:"type"`
}

type CategoryThreeStat struct {
	BiCategoryOneName   string `bson:"bi_category_one_name"`
	BiCategoryTwoName   string `bson:"bi_category_two_name"`
	BiCategoryThreeName string `bson:"bi_category_three_name"`
	Competitiveness     int    `bson:"competitiveness"`
	GoodMatchCount      int    `bson:"good_match_count"`
	Type                string `bson:"type"`
}

type AllStat struct {
	Competitiveness int    `bson:"competitiveness"`
	GoodMatchCount  int    `bson:"good_match_count"`
	Type            string `bson:"type"`
}

func upsertStat(mongoSession *mgo.Session, categoryOneStatTable []CategoryOneStat, categoryTwoStatTable []CategoryTwoStat, categoryThreeStatTable []CategoryThreeStat, allStat AllStat, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	c := mongoSession.DB("competition_analysis").C("bml_dgk_agg_statistic_hist")

	config := mongobulk.Config{OpsPerBatch: 950}

	mongoBulk := mongobulk.New(c, config)

	upsertAllStat(&allStat, mongoBulk)

	for _, categoryOneStat := range categoryOneStatTable {
		upsertCategoryOneStat(&categoryOneStat, mongoBulk)
	}
	for _, categoryTwoStat := range categoryTwoStatTable {
		upsertCategoryTwoStat(&categoryTwoStat, mongoBulk)
	}
	for _, categoryThreeStat := range categoryThreeStatTable {
		upsertCategoryThreeStat(&categoryThreeStat, mongoBulk)
	}

	err := mongoBulk.Finish()
	checkError(err)

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)

}

func upsertAllStat(allStat *AllStat, mongoBulk *mongobulk.Bulk) {

	allStatByte, err := bson.Marshal(allStat)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())
	}
	allStatBson := make(map[string]interface{})
	err = bson.Unmarshal(allStatByte, allStatBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())
	}
	colQuerier := bson.M{"type": allStat.Type}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": allStatBson})
}

func upsertCategoryOneStat(categoryOneStatTable *CategoryOneStat, mongoBulk *mongobulk.Bulk) {

	categoryOneStatTableByte, err := bson.Marshal(categoryOneStatTable)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())
	}
	categoryOneStatTableBson := make(map[string]interface{})
	err = bson.Unmarshal(categoryOneStatTableByte, categoryOneStatTableBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())
	}
	colQuerier := bson.M{"bi_category_one_name": categoryOneStatTable.BiCategoryOneName}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": categoryOneStatTableBson})
}

func upsertCategoryTwoStat(categoryTwoStatTable *CategoryTwoStat, mongoBulk *mongobulk.Bulk) {

	categoryTwoStatTableByte, err := bson.Marshal(categoryTwoStatTable)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())
	}
	categoryTwoStatTableBson := make(map[string]interface{})
	err = bson.Unmarshal(categoryTwoStatTableByte, categoryTwoStatTableBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())
	}
	colQuerier := bson.M{"bi_category_Two_name": categoryTwoStatTable.BiCategoryTwoName}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": categoryTwoStatTableBson})
}

func upsertCategoryThreeStat(categoryThreeStatTable *CategoryThreeStat, mongoBulk *mongobulk.Bulk) {

	categoryThreeStatTableByte, err := bson.Marshal(categoryThreeStatTable)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())
	}
	categoryThreeStatTableBson := make(map[string]interface{})
	err = bson.Unmarshal(categoryThreeStatTableByte, categoryThreeStatTableBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())
	}
	colQuerier := bson.M{"bi_category_Three_name": categoryThreeStatTable.BiCategoryThreeName}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": categoryThreeStatTableBson})
}
