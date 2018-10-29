package bmldgkmanualmatching

import(
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
)

func SendFrequencyOptionList(client *websocketserver.Client,data interface{}){


	frequencyListArr := []string {"monthly", "two times monthly","weekly","every 3 days","every 2 days",
	"daily","two times daily","three times daily"}
	frequencyIdArr := []int{1,2,4,10,15,30,60,90}

	var frequncyOptionList bmldgktable.Frequency
	for i,tmpFrequencyListArr := range frequencyListArr{
		var frequncyOption bmldgktable.FrequencyOptionList
		frequncyOption.OptionText = tmpFrequencyListArr
		frequncyOption.OptionValue = frequencyIdArr[i]
		frequncyOptionList.FrequencyOptionList = append(frequncyOptionList.FrequencyOptionList,frequncyOption)
	}
	


	// set values to message and send
	var message websocketserver.Message
	message.Name = "FrequencyOptionList get"
	message.Data = frequncyOptionList
	client.Send <- message
}