package configcountcategorypage

import(
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo"
	"sort"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendBmlCategory1FilterOptionList(client *websocketserver.Client,data interface{}){

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
	var bmlCategory []string
	db.C("bml_catalog_config").Find(bson.M{}).Distinct("bi_category_one_name", &bmlCategory)
	sort.Strings(bmlCategory)
	if(bmlCategory[0]==""){
		bmlCategory = bmlCategory[1:len(bmlCategory)]
	}
	bmlCategory = append(bmlCategory,"All")
	sort.Strings(bmlCategory)
	var cat bmldgktable.Category 
	for _,bmlCategoryRng := range bmlCategory{
		tmpCat := bmldgktable.OptionList{
			OptionValue:bmlCategoryRng,
			OptionText:bmlCategoryRng,
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
	message.Name = "bmlFilterCategory1 get"
	message.Data = cat
	client.Send <- message
	}