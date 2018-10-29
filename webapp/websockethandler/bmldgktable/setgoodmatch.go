package bmldgktable

import (
	"fmt"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SetGoodMatch(client *websocketserver.Client, data interface{}) {

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
	var eventbool map[string]bool
	var eventStr map[string]string
	mapstructure.Decode(data, &eventbool)
	mapstructure.Decode(data, &eventStr)

	fmt.Printf("%#v\n", eventbool)
	fmt.Printf("%#v\n", eventStr)
	// extract message

	iDBmlCatalogConfig, _ := strconv.Atoi(eventStr["ID"])
	isGoodMatched := eventbool["IsGoodMatched"]
	email := eventStr["Email"]
	name := eventStr["Name"]

	bmldgktable.ApplyGoodMatch(iDBmlCatalogConfig, isGoodMatched, email, name, mongoSession)

}
