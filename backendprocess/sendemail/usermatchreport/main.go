package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/email/emailconf"
	sendemail "github.com/thomas-bamilo/email/sendemail"
)

//ready for release!
func main() {
	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	checkError(err)

	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	now := time.Now()
	date := now.Format(`Jan 2`)

	var emailconf emailconf.EmailConf

	emailconf.ReadYamlEmailConf()

	var userMatchDataTable []userMatchData
	getUserMatchData(mongoSession, &userMatchDataTable)

	log.Println(len(userMatchDataTable))

	emailTitle := emailconf.EmailTitle + ` ` + date
	emailBody := designTemplate(userMatchDataTable, date)

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
}

func getUserMatchData(mongoSession *mgo.Session, userMatchDataTable *[]userMatchData) {

	db := mongoSession.DB("competition_analysis")
	c := db.C("user")

	getAllDoc := bson.M{"$match": bson.M{}}
	sortByThisWeekSkuCount := bson.M{"$sort": bson.M{"total_matched_sku": -1}}

	pipe := c.Pipe([]bson.M{getAllDoc, sortByThisWeekSkuCount})
	err := pipe.All(userMatchDataTable)
	checkError(err)

}

type userMatchData struct {
	Email                    string `bson:"email"`
	Last7DayMatchedSku       int    `bson:"last_7_day_matched_sku"`
	Last14DayMatchedSku      int    `bson:"last_14_day_matched_sku"`
	TotalMatchedSku          int    `bson:"total_matched_sku"`
	Last7DayMatchedSupplier  int    `bson:"last_7_day_matched_supplier"`
	Last14DayMatchedSupplier int    `bson:"last_14_day_matched_supplier"`
	TotalMatchedSupplier     int    `bson:"total_matched_supplier"`
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func designTemplate(userMatchDataTable []userMatchData, date string) (template string) {

	userMatchDataTableF := userMatchDataTable[:0]

	for _, x := range userMatchDataTable {
		if !stringInSlice(x.Email, []string{
			`putik.dhiraramanti@bamilo.com`,
			`mohammadhadi.mojahedi@bamilo.com`,
			`aram.jafari@bamilo.com`,
			`sundeep.sahni@bamilo.com`,
			`cleland.robertson@bamilo.com`,
			`thomas.beaudouin@bamilo.com`,
			`sepideh.haghverdi@bamilo.com`,
		}) {
			userMatchDataTableF = append(userMatchDataTableF, x)
		}
	}

	template = `<p>Please find below the number of SKUs and suppliers matched by users for the past 7 and 14 days as of ` + date + `</p>

	<table style="border: 1px solid black; border-collapse: collapse;">
		<tr>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Email</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Last 7 days #SKU</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Last 14 days #SKU</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Total #SKU</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Last 7 days #Supplier</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Last 14 days #Supplier</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">Total #Supplier</th>
		</tr> `

	for _, userMatchData := range userMatchDataTableF {

		template += `
			<tr> 
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + userMatchData.Email + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(userMatchData.Last7DayMatchedSku) + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(userMatchData.Last14DayMatchedSku) + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(userMatchData.TotalMatchedSku) + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(userMatchData.Last7DayMatchedSupplier) + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(userMatchData.Last14DayMatchedSupplier) + `</th>
			<th style="text-align: center; padding: 5px; border: 1px solid black; border-collapse: collapse;">` + strconv.Itoa(userMatchData.TotalMatchedSupplier) + `</th>
			</tr> `
	}

	template += `
		</table>`

	return template
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
