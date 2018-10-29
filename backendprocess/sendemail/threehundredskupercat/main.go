package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/email/emailconf"
	sendemail "github.com/thomas-bamilo/email/sendemail"
	"github.com/thomas-bamilo/sql/connectdb"
)

//ready for release!
func main() {
	start := time.Now()
	now := time.Now()
	date := now.Format(`Jan 2`)
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	db := mongoSession.DB("competition_analysis")

	var asdData []BmlAllResult

	db.C("top_300_bml_sku").Find(nil).Select(bson.M{
		"id_bml_catalog_config": 17,
		"bi_category":           17,
		"sku_name":              17,
		"count_of_soi":          17,
		"supplier_name":         17,
		"good_match":            17,
		"sku":                   17,
		"brand_name":            17}).All(&asdData)

	dbSqlite := connectdb.ConnectToSQLite()
	defer dbSqlite.Close()
	query := `CREATE TABLE asdTable (
		    bi_category              TEXT,
			good_match               TEXT,
			count                 INTEGER
			)`
	queryP, err := dbSqlite.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO asdTable (
		bi_category,
		good_match,	
		count  
			) 
		VALUES (?, ?,?)`
	queryP, err = dbSqlite.Prepare(query)
	checkError(err)
	for i := 0; i < len(asdData); i++ {

		queryP.Exec(
			asdData[i].BiCategory,
			asdData[i].GoodMatch,
			"",
		)
		time.Sleep(1 * time.Millisecond)
	}

	query = `SELECT bi_category, COUNT(*) 'count'
	FROM asdTable
	WHERE good_match = '1'
	GROUP BY  bi_category
	ORDER BY
    count DESC;
		`
	rows, err := dbSqlite.Query(query)
	checkError(err)

	var asdData1 BiCatMatch
	var asdDatatable []BiCatMatch

	for rows.Next() {
		err := rows.Scan(&asdData1.BiCategory,
			&asdData1.Count)
		checkError(err)

		asdDatatable = append(asdDatatable, asdData1)
	}

	var arrOfStr [][]string

	arrOfStr = append(arrOfStr, []string{
		`SKU`,
		`Bi Category`,
		`SKU Name`,
		`Count Of Soi`,
		`Supplier Name`,
		`Brand Name`,
		`Good Match`,
	})

	for i := 0; i < len(asdData); i++ {
		arrOfStr = append(arrOfStr, []string{
			string(asdData[i].Sku),
			string(asdData[i].BiCategory),
			string(asdData[i].SKUName),
			strconv.Itoa(asdData[i].CountOfSoi),
			string(asdData[i].SupplierName),
			string(asdData[i].BrandName),
			strconv.FormatBool(asdData[i].GoodMatch),
		})
	}

	file, err := os.Create("MatchedSkusPerCat.csv")
	checkError(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range arrOfStr {
		err := writer.Write(value)
		checkError(err)
	}
	var emailconf emailconf.EmailConf

	emailconf.ReadYamlEmailConf()

	emailTitle := emailconf.EmailTitle + ` ` + date

	emailBody := designTemplate(asdDatatable)

	sendemail.SendHTMLEmail(
		emailconf.EmailAttachPath,                    // attachPath
		strings.Split(emailconf.EmailRecipient, ","), // toEmail
		emailTitle, // title
		emailBody,  // message
		emailconf.SenderIdentity,
		emailconf.SenderName,
		emailconf.SenderEmail,
		emailconf.SenderPw,
		emailconf.SMTPAddr,
		emailconf.SMTPPort,
	)
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

func designTemplate(asdDatatable []BiCatMatch) (template string) {
	template = `<p>Please find below the number of matched SKUs amoung top 300 sales order Item during last 30 days </p>

	<table style="border: 1px solid black; border-collapse: collapse;">
		<tr>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Category</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;"># Matched SKUs Within Top 300</th>
		</tr> `

	for _, asdDatatable := range asdDatatable {

		template += `
			<tr> 
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + asdDatatable.BiCategory + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(asdDatatable.Count) + `</th>
			</tr> `
	}

	template += `
		</table>`

	return template
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

type BiCatMatch struct {
	BiCategory string `json:"bi_category"`
	Count      int    `json:"count"`
}
