package bmldgkhistory

import (
	"fmt"
	"log"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgkhistoricaldata"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
	
)

func SendBmlDgkHistoricalData(client *websocketserver.Client, data interface{}) {
	log.Println(`data`, data)
	// read received values
	var event map[string]string
	mapstructure.Decode(data, &event)
	fmt.Printf("%#v\n", event)

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	// extract message
	bmlIDCatalogConfig, _ := strconv.Atoi(event["BmlIDCatalogConfig"])
	BmlSumOfStockQuantity, _ := strconv.Atoi(event["BmlSumOfStockQuantity"])
	BmlSKUName := event["BmlSKUName"]
	BmlImgLink := event["BmlImgLink"]
	BmlSupplierName := event["BmlSupplierName"]
	BmlBrand := event["BmlBrand"]
	BmlConfigSnapshot := event["BmlConfigSnapshot"]
	BmlSKULink := event["BmlSKULink"]
	DgkIDCatalogConfig := event["DgkIDCatalogConfig"]
	DgkSKUName := event["DgkSKUName"]
	DgkImgLink := event["DgkImgLink"]
	DgkSupplierName := event["DgkSupplierName"]
	DgkBrand := event["DgkBrand"]
	DgkConfigSnapshot := event["DgkConfigSnapshot"]
	DgkStock := event["DgkStock"]
	DgkSKULink := event["DgkSKULink"]
	

	historicalChartData := bmldgkhistoricaldata.FetchBmlDgkHistoricalData(bmlIDCatalogConfig,
		BmlSumOfStockQuantity,BmlSKUName ,BmlImgLink ,BmlSupplierName ,
		BmlBrand ,BmlConfigSnapshot ,BmlSKULink,DgkIDCatalogConfig,
		DgkSKUName,DgkImgLink ,DgkSupplierName,DgkBrand,DgkConfigSnapshot,DgkStock ,DgkSKULink,
		 mongoSession)

	// set values to message and send
	var message websocketserver.Message
	message.Name = "BmlDgkSKUHistoricalData send"
	message.Data = historicalChartData
	client.Send <- message
}
