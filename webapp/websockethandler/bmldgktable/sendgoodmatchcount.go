package bmldgktable

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendGoodMatchCount(client *websocketserver.Client, data interface{}) {

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
	mapstructure.Decode(data, &eventStr)

	fmt.Printf("%#v\n", eventStr)
	// extract message

	goodMatchCount := bmldgktable.GetGoodMatchCount(eventStr["category1"], eventStr["category2"], eventStr["category3"], mongoSession)

	// set values to message and send
	var message websocketserver.Message
	message.Name = "goodMatchCount get"
	message.Data = goodMatchCount
	client.Send <- message

}
