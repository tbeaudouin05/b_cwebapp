package configcountsupplierpage

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/configcountsupplierpage"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendMatchedSupplierTable(client *websocketserver.Client, data interface{}) {

	// read received values
	var eventStr map[string]string
	mapstructure.Decode(data, &eventStr)
	fmt.Printf("%#v\n", eventStr)

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	// building table
	dgkSearchedBy := eventStr["DgkSearchedByMatch"] 
	bmlSearchedBy := eventStr["BmlSearchedByMatch"] 
	table := configcountsupplierpage.SelectMatchedSupplierTable(bmlSearchedBy,dgkSearchedBy, mongoSession)


	//fmt.Println("handler check",table)

	// set values to message and send
	var message websocketserver.Message
	message.Name = "matchedSellerList get"
	message.Data = table
	client.Send <- message
}

