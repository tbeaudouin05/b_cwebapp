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

func SendBiCategoryFilterOptionList(client *websocketserver.Client, data interface{}) {

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
	var bmlBiCategory []string

	db.C("bml_catalog_config").Find(bson.M{}).Distinct("bi_category", &bmlBiCategory)

	sort.Strings(bmlBiCategory)

	var cat bmldgktable.Category
	cat.OptionList = append(cat.OptionList,
		bmldgktable.OptionList{
			OptionValue: "",
			OptionText:  "All",
		})
	for i := 0; i < len(bmlBiCategory); i++ {
		if bmlBiCategory[i] != "" {
			cat.OptionList = append(cat.OptionList,
				bmldgktable.OptionList{
					OptionValue: bmlBiCategory[i],
					OptionText:  bmlBiCategory[i],
				})
		}

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
	message.Name = "filterBiCategory get"
	message.Data = cat
	client.Send <- message
}
