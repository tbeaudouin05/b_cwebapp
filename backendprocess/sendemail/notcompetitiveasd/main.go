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

	var asdData []AsdNotCompetitiveMatchedSKU

	db.C("bml_dgk_agg_statistic_hist").Find(bson.M{"type": "NotCompetitiveGoodMatched"}).Select(bson.M{
		"bml_sku_name":           17,
		"dgk_sku_name":           17,
		"bml_sku_link":           17,
		"sku":                    17,
		"supplier_name":          17,
		"dgk_sku_link":           17,
		"bml_sku":                17,
		"bml_supplier_name":      17,
		"brand_name":             17,
		"department":             17,
		"bi_category":            17,
		"bml_price":              17,
		"dgk_price":              17,
		"color_of_best_price":    17,
		"supplier_of_best_price": 17,
		"warranty_of_best_price": 17,
		"dgk_best_price":         17,
		`key_account_manager`:    17}).All(&asdData)

	dbSqlite := connectdb.ConnectToSQLite()
	defer dbSqlite.Close()
	query := `CREATE TABLE asdTable (
			asd                      TEXT,
			count                 INTEGER
			)`
	queryP, err := dbSqlite.Prepare(query)
	checkError(err)
	queryP.Exec()

	query = `INSERT INTO asdTable (
		asd,	
		count  
			) 
		VALUES (?, ?)`
	queryP, err = dbSqlite.Prepare(query)
	checkError(err)
	for i := 0; i < len(asdData); i++ {

		queryP.Exec(
			asdData[i].KeyAccountManager,
			"",
		)
		time.Sleep(1 * time.Millisecond)
	}

	query = `SELECT asd, COUNT(*) 'count'
	FROM asdTable
	GROUP BY asd
	ORDER BY
    count DESC;
		`
	rows, err := dbSqlite.Query(query)
	checkError(err)

	var asdData1 AsdData
	var asdDatatable []AsdData

	for rows.Next() {
		err := rows.Scan(&asdData1.AsdName,
			&asdData1.Count)
		checkError(err)

		asdDatatable = append(asdDatatable, asdData1)
	}

	var arrOfStr [][]string

	arrOfStr = append(arrOfStr, []string{
		`Bamilo Sku Name`,
		`Digikala Sku Name`,
		`Bamilo Sku Link`,
		`Digikala Sku Link`,
		`Bamilo Sku`,
		`Bamilo Supplier Name`,
		`Brand Name`,
		`Department`,
		`Bi Category`,
		`Bamilo Price`,
		`Digikala Price`,
		`Color Of Best Price`,
		`Supplier Of Best Price`,
		`Warranty Of Best Price`,
		`Digikala Best Price`,
		`Key Account Manager`,
	})

	for i := 0; i < len(asdData); i++ {
		arrOfStr = append(arrOfStr, []string{
			string(asdData[i].BmlSkuName),
			string(asdData[i].DgkSkuName),
			string(asdData[i].BmlSkuLink),
			string(asdData[i].DgkSkuLink),
			string(asdData[i].BmlSku),
			string(asdData[i].BmlSupplierName),
			string(asdData[i].BrandName),
			string(asdData[i].Department),
			string(asdData[i].BiCategory),
			strconv.Itoa(asdData[i].BmlPrice),
			strconv.Itoa(asdData[i].DgkPrice),
			string(asdData[i].ColorOfBestPrice),
			string(asdData[i].SupplierOfBestPrice),
			string(asdData[i].WarrantyOfBestPrice),
			strconv.Itoa(asdData[i].DgkBestPrice),
			string(asdData[i].KeyAccountManager),
		})
	}

	file, err := os.Create("NonCompetitiveSkusPerASD.csv")
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

type AsdNotCompetitiveMatchedSKU struct {
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
	BmlPrice            int    `bson:"bml_price"`
	DgkPrice            int    `bson:"dgk_price"`
	ColorOfBestPrice    string `bson:"color_of_best_price"`
	SupplierOfBestPrice string `bson:"supplier_of_best_price"`
	WarrantyOfBestPrice string `bson:"warranty_of_best_price"`
	DgkBestPrice        int    `bson:"dgk_best_price"`
	KeyAccountManager   string `bson:"key_account_manager"`
}

type AsdData struct {
	AsdName string `json:"asd"`
	Count   int    `json:"count"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func designTemplate(asdDatatable []AsdData) (template string) {
	template = `<p>Please find below the number of nonCompetitive SKUs per ASD </p>

	<table style="border: 1px solid black; border-collapse: collapse;">
		<tr>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">ASD</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">number of nonCompetitive SKUs</th>
		</tr> `

	for _, asdDatatable := range asdDatatable {

		template += `
			<tr> 
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + asdDatatable.AsdName + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(asdDatatable.Count) + `</th>
			</tr> `
	}

	template += `
		</table>`

	return template
}
