package bmldgktable

import (
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)

func SendSortingOptionList(client *websocketserver.Client, data interface{}) {

	// choices
	sortingList := []string{"Bamilo best sales", "Matching score"} // "Digikala best sales",
	sortingValueList := []string{"-count_of_soi", "-dgk_score"}    //"SkuRank"
	var list []bmldgktable.OptionList
	for i, sortingListRng := range sortingList {
		tmp := bmldgktable.OptionList{
			OptionText:  sortingListRng,
			OptionValue: sortingValueList[i],
		}
		list = append(list, tmp)
	}

	// null checker
	if list == nil {
		var tmp bmldgktable.OptionList
		tmp.OptionText = ""
		tmp.OptionValue = ""
		list = append(list, tmp)
	}

	// set values to message and send
	var message websocketserver.Message
	message.Name = "sortingFilter get"
	message.Data = list
	client.Send <- message
}
