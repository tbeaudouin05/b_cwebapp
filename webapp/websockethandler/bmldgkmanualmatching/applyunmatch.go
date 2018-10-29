package bmldgkmanualmatching

import (
	"fmt"
	"strconv"
	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgkmanualmatching"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func ApplyUnmatch(client *websocketserver.Client, data interface{}) {

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
	var eventInt map[string]int
	var eventStr map[string]string
	mapstructure.Decode(data, &eventInt)
	mapstructure.Decode(data, &eventStr)
	fmt.Printf("%#v\n", eventInt)
	fmt.Printf("%#v\n", eventStr)
	// extract message

	IDBmlCatalogConfig,_ := strconv.Atoi(eventStr["IDBmlCatalogConfig"])
	//DgkImgLink := eventStr["DgkImgLink"]
	DgkImgLink := bmldgkmanualmatching.ApplyUnmatch(IDBmlCatalogConfig, mongoSession)
	var message websocketserver.Message
	message.Name = "unmatch status"
	message.Data = DgkImgLink

	client.Send <- message

}