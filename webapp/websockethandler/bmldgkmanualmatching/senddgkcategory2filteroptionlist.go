package bmldgkmanualmatching

import(
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo"
	"sort"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendDgkCategory2FilterOptionList(client *websocketserver.Client,data interface{}){


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
	category1 :=  eventStr["category1"]



	// category choices
	db := mongoSession.DB("competition_analysis")
	var dgkCategory2 []string
	if category1!="All"{
		db.C("dgk_catalog_config").Find(bson.M{"dgk_category_one_name" : category1,}).Distinct("dgk_category_two_name", &dgkCategory2)
		sort.Strings(dgkCategory2)
	}else{
		db.C("dgk_catalog_config").Find(nil).Distinct("dgk_category_two_name", &dgkCategory2)
		sort.Strings(dgkCategory2)
	}
	
	if(dgkCategory2[0]==""){
		dgkCategory2 = dgkCategory2[1:len(dgkCategory2)]
	}
	dgkCategory2 = append(dgkCategory2,"All")
	sort.Strings(dgkCategory2)
	var cat bmldgktable.Category 
	for _,dgkCategory2Rng := range dgkCategory2{
		tmpCat := bmldgktable.OptionList{
			OptionValue:dgkCategory2Rng,
			OptionText:dgkCategory2Rng,
		}
		cat.OptionList = append(cat.OptionList,tmpCat)
	}

	// null checker	
	if cat.OptionList==nil{
		var tmpCat bmldgktable.OptionList
		tmpCat.OptionText= ""
		tmpCat.OptionValue= ""
		cat.OptionList = append(cat.OptionList,tmpCat)
	} 

	//fmt.Println(cat)
	// set values to message and send
	var message websocketserver.Message
	message.Name = "dgkFilterCategory2 get"
	message.Data = cat
	client.Send <- message
	}