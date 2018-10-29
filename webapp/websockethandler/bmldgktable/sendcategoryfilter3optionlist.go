package bmldgktable

import (
	"fmt"
	"sort"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendCategoryFilter3OptionList(client *websocketserver.Client, data interface{}) {

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

	// category choices
	db := mongoSession.DB("competition_analysis")
	var bmlCategory3 []string

	db.C("bml_catalog_config").Find(bson.M{}).Distinct("bi_category_three_name", &bmlCategory3)

	sort.Strings(bmlCategory3)

	var cat bmldgktable.Category

	flag := false // to add All if "" doesnt exist in category3
	for i := 0; i < len(bmlCategory3); i++ {
		if bmlCategory3[i] == "" {
			cat.OptionList = append(cat.OptionList,
				bmldgktable.OptionList{
					OptionValue: "",
					OptionText:  "All",
				})
			flag = true
		} else {
			cat.OptionList = append(cat.OptionList,
				bmldgktable.OptionList{
					OptionValue: bmlCategory3[i],
					OptionText:  bmlCategory3[i],
				})
		}

	}
	if !flag {
		cat.OptionList = append(cat.OptionList,
			bmldgktable.OptionList{
				OptionValue: "",
				OptionText:  "All",
			})
	}

	// null checker
	if cat.OptionList == nil {
		cat.OptionList = append(cat.OptionList,
			bmldgktable.OptionList{
				OptionValue: "",
				OptionText:  "",
			})
	}

	// set values to message and send
	var message websocketserver.Message
	message.Name = "filterCategory3 get"
	message.Data = cat
	client.Send <- message
}
