package bmldgkmanualmatching

import (
	"fmt"
	"strconv"
	"github.com/SepidehKHH/formatting"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgkmanualmatching"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)


func RunManualMatchingPage(client *websocketserver.Client,data interface{}){
	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)
	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")

	// read received values
	var eventStr map[string]string
	var eventNum map[string]int
	mapstructure.Decode(data, &eventStr)
	mapstructure.Decode(data, &eventNum)
	fmt.Printf("%#v\n", eventStr)
	fmt.Printf("%#v\n", eventNum)
	bmlIDCatalogConfig,_ := strconv.Atoi(eventStr["bmlID"])
	pgNumber := eventNum["pageNumber"]
	cat1 :=  eventStr["category1"]
	cat2 :=  eventStr["category2"]
	searchedBy := eventStr["skuName"]
	minPrice,err1 := strconv.Atoi(eventStr["minPrice"])
	if (err1 !=nil){
		fmt.Println(err1)
		minPrice = -5
	}
	maxPrice,err2 := strconv.Atoi(eventStr["maxPrice"])
	if (err2 !=nil){
		fmt.Println(err1)
		minPrice = 9999999999
	}
	


	var bamiloCatalogConfigStr bmlcatalogconfig.BmlCatalogConfig
	findBamiloCatalogConfig := bson.M{"$match": bson.M{"_id": bmlIDCatalogConfig}}
	pipe := bmlCatalogConfig.Pipe([]bson.M{findBamiloCatalogConfig})
	err = pipe.One(&bamiloCatalogConfigStr)
	if (err !=nil){
		fmt.Println(err)
	}

	// buil table and set values to message and send
	var message websocketserver.Message
	
	var dgkTableInterface bmldgktable.DgkTableInterface

	dgkTable,isTableNull := bmldgkmanualmatching.RunManualMatchingPage(bmlIDCatalogConfig, pgNumber, searchedBy, cat1, cat2, minPrice, maxPrice, mongoSession)
	fmt.Println("is it null?",isTableNull)
	if(!isTableNull){
		dgkTableInterface.Table = dgkTable.Table
		dgkTableInterface.PageNumber = bmldgkmanualmatching.ShowPageNumber(bmlIDCatalogConfig, pgNumber, searchedBy, cat1, cat2, minPrice, maxPrice, mongoSession)
	}else{
		dgkTableInterface.Table = bmldgkmanualmatching.SetNullTabel()
		dgkTableInterface.PageNumber = append(dgkTableInterface.PageNumber,1)
	}
	dgkTableInterface.BmlSkuName = bamiloCatalogConfigStr.SKUName
	tempPrice := bamiloCatalogConfigStr.AvgSpecialPrice
	if (tempPrice==0){
		tempPrice = bamiloCatalogConfigStr.AvgPrice
	}
	dgkTableInterface.BmlPrice = formatting.ChangeNumberFormat(strconv.Itoa(tempPrice))
	dgkTableInterface.BmlImgLink = bamiloCatalogConfigStr.ImgLink
	dgkTableInterface.BmlBrand = bamiloCatalogConfigStr.BrandName
	dgkTableInterface.BmlSKULink = bamiloCatalogConfigStr.SKULink
	

	message.Name = "manualMatchingTablePage change"

	message.Data = dgkTableInterface

	client.Send <- message
}


