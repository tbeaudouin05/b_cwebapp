package configcountsupplierpage

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"

	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/configcountsupplierpage"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func MatchSupplier(client *websocketserver.Client, data interface{}) {
	fmt.Println("check handler")
	fmt.Println("check handler")
	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	// read received values
	var eventStr map[string]string
	var eventInt map[string]int
	mapstructure.Decode(data, &eventStr)
	fmt.Printf("%#v\n", eventStr)
	mapstructure.Decode(data, &eventInt)
	fmt.Printf("%#v\n", eventInt)
	// extract message

	selectedDgkSellerID := eventStr["SelectedDgkSellerID"]
	selectedDgkSellerName := eventStr["SelectedDgkSellerName"]
	selectedBmlSellerID := eventInt["SelectedBmlSellerID"]
	selectedBmlSellerName := eventStr["SelectedBmlSellerName"]
	email := eventStr["Email"]
	name := eventStr["Name"]

	var message websocketserver.Message
	message.Name = "MatchSupplier status"
	message.Data = configcountsupplierpage.ApplyMatchSupplier(selectedBmlSellerID, selectedBmlSellerName, selectedDgkSellerID, selectedDgkSellerName, email, name, mongoSession)

	client.Send <- message

}
