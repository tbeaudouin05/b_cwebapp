package configcountsupplierpage

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/configcountsupplierpage"
)

func ApplySellerUrl(client *websocketserver.Client, data interface{}) {
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
	mapstructure.Decode(data, &eventStr)
	fmt.Printf("%#v\n", eventStr)
	// extract message

	SellerUrl := eventStr["SellerUrl"] 
	
	var message websocketserver.Message
	message.Name = "sellerUrlValidation status"

		if len(SellerUrl)>33  && SellerUrl[:32] == "https://www.digikala.com/seller/" && SellerUrl[len(SellerUrl)-1:]=="/" {
			configcountsupplierpage.ApplySellerURL(SellerUrl, mongoSession)
			message.Data = "The URL successfully added."
		}else {
			message.Data = "Faild: Wrong URL"
		}
	
	client.Send <- message
	



}