package bmldgkmanualmatching

import (
	"fmt"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgkmanualmatching"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func ApplyManualMatching(client *websocketserver.Client, data interface{}) {

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

	iDBmlCatalogConfig, _ := strconv.Atoi(eventStr["IDBmlCatalogConfig"])
	fKDgkCatalogConfig := eventStr["FKDgkCatalogConfig"]
	dgkImgLink := eventStr["DgkImgLink"]
	email := eventStr["Email"]
	name := eventStr["Name"]

	bmldgkmanualmatching.ApplyManualMatching(iDBmlCatalogConfig, fKDgkCatalogConfig, email, name, mongoSession)
	var message websocketserver.Message
	message.Name = "matching status"
	message.Data = dgkImgLink

	client.Send <- message

}
