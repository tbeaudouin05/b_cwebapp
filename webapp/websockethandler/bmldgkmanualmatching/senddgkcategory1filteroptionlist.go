package bmldgkmanualmatching

import(
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo"
	"sort"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendDgkCategory1FilterOptionList(client *websocketserver.Client,data interface{}){

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	
	// category choices
	db := mongoSession.DB("competition_analysis")
	var dgkCategory []string
	db.C("dgk_catalog_config").Find(bson.M{}).Distinct("dgk_category_one_name", &dgkCategory)
	sort.Strings(dgkCategory)
	if(dgkCategory[0]==""){
		dgkCategory = dgkCategory[1:len(dgkCategory)]
	}
	dgkCategory = append(dgkCategory,"All")
	sort.Strings(dgkCategory)
	var cat bmldgktable.Category 
	for _,dgkCategoryRng := range dgkCategory{
		tmpCat := bmldgktable.OptionList{
			OptionValue:dgkCategoryRng,
			OptionText:dgkCategoryRng,
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

	// set values to message and send
	var message websocketserver.Message
	message.Name = "dgkFilterCategory1 get"
	message.Data = cat
	client.Send <- message
	}