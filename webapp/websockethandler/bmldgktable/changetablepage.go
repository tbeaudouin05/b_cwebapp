package bmldgktable

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgktable"
	bmldgktablestruct "github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func ChangeTablePage(client *websocketserver.Client, data interface{}) {

	// read received values
	var eventStr map[string]string
	var eventNum map[string]int
	mapstructure.Decode(data, &eventStr)
	mapstructure.Decode(data, &eventNum)

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	// building table
	var table bmldgktablestruct.Table
	pgNumber := eventNum["pageNumber"]
	biCat := eventStr["biCategory"]
	cat1 := eventStr["category1"]
	cat2 := eventStr["category2"]
	cat3 := eventStr["category3"]
	searchedBy := eventStr["skuName"]
	sortedBy := eventStr["sorting"]
	log.Println(`sortedBy12: `, sortedBy)
	if sortedBy == "" {
		sortedBy = "CountOfSoi"
	}
	log.Println(`searchedBy: `, searchedBy)
	log.Println(`FILTERS: `, pgNumber, cat1, cat2, cat3, searchedBy, sortedBy)
	var tableIsNull bool
	table, tableIsNull = bmldgktable.SelectCompetitionAnalysisTable(10, pgNumber, searchedBy, biCat, cat1, cat2, cat3, sortedBy, mongoSession)
	fmt.Println(tableIsNull)

	//check if table is null
	if tableIsNull {
		table = bmldgktable.SetNullTabel()
	}
	var tableInterface bmldgktablestruct.TableInterface
	tableInterface.Table = table.Table
	showPageArr := bmldgktable.ShowPageNumber(pgNumber, searchedBy, biCat, cat1, cat2, cat3, sortedBy, mongoSession)
	tableInterface.PageNumber = showPageArr

	// set values to message and send
	var message websocketserver.Message
	message.Name = "tablePage change"
	message.Data = tableInterface
	client.Send <- message
}

func checkNull(data interface{}) bool {

	if data == nil {
		return true
	}
	// add other stuff later

	return false
}
