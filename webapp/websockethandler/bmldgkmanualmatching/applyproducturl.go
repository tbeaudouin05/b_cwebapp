package bmldgkmanualmatching

import (
	"fmt"
	"strings"
	"github.com/globalsign/mgo"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgkmanualmatching"
)

func ApplyProductUrl(client *websocketserver.Client, data interface{}) {
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

	ProductUrl := eventStr["ProductUrl"] 
	BmlID ,_ := strconv.Atoi(eventStr["BmlID"])

	var message websocketserver.Message
	message.Name = "productUrlValidation status"


	if(len(ProductUrl)>0){
		if(ProductUrl[len(ProductUrl)-1:] == "?"){
			ProductUrl = ProductUrl[:len(ProductUrl)-1]
		}

	}
	fmt.Println(ProductUrl)
		if len(ProductUrl)>34  && ProductUrl[:33] == "https://www.digikala.com/product/" && strings.Contains(ProductUrl,"%") {
			resultVal,skuName := bmldgkmanualmatching.ApplyProductURL(ProductUrl, BmlID, mongoSession)
			if resultVal==0{
				message.Data = []string{"The URL successfully added."}
			}else if resultVal==1{
				message.Data = []string{"The Product already exists, the sku name is:" ,skuName}
				
			}else if resultVal==5{
				message.Data = []string{"There is an error, please try later."}
			}else{
				message.Data = []string{"The Product URL is added."}
			}
		}else {
			message.Data = []string{"Faild: Wrong URL"}
		}
		fmt.Println(message.Data )
	
	client.Send <- message
	



}