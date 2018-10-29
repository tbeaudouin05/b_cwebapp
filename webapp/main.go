package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"encoding/gob"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/websocketserver"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/adduser"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/websockethandler/bmldgkhistory"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/websockethandler/bmldgkmanualmatching"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/websockethandler/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/websockethandler/configcountcategorypage"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/websockethandler/configcountsupplierpage"
	homepagewebsockethandler "github.com/thomas-bamilo/commercial/competitionanalysis/webapp/websockethandler/homepage"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/bmldgkmanualmatchingrequest"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/bmldgkskuhistoryrequest"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/bmldgktablecsvoutputrequest"
	bmldgkskuhistfrontend "github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/bmldgkskuhist"
	bmldgktablefrontend "github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/categoryconfigcount"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/homepage"
	manualmatchingfrontend "github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/manualmatching"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authenticate"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/login"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/supplierconfigcount"
)

func main() {

	config := loadEnvConfig(`envConfig.json`)

	go startWebsocketServer(config)

	startFrontEndServer(config)

}

// -------------------------------------------------------------------------------------------------

func startWebsocketServer(config config) {
	routerWebsocket := websocketserver.NewRouter()

	routerWebsocket.Handle("tablePage request", bmldgktable.ChangeTablePage)
	routerWebsocket.Handle("goodMatchCount request", bmldgktable.SendGoodMatchCount)
	routerWebsocket.Handle("biCategoryFilterOptionList request", bmldgktable.SendBiCategoryFilterOptionList)
	routerWebsocket.Handle("categoryFilter1OptionList request", bmldgktable.SendCategoryFilter1OptionList)
	routerWebsocket.Handle("categoryFilter2OptionList request", bmldgktable.SendCategoryFilter2OptionList)
	routerWebsocket.Handle("categoryFilter3OptionList request", bmldgktable.SendCategoryFilter3OptionList)
	routerWebsocket.Handle("BmlDgkSKUHistoricalData request", bmldgkhistory.SendBmlDgkHistoricalData)
	routerWebsocket.Handle("ManualMatchingTablePage request", bmldgkmanualmatching.RunManualMatchingPage)
	routerWebsocket.Handle("dgkCategoryFilter1OptionList request", bmldgkmanualmatching.SendDgkCategory1FilterOptionList)
	routerWebsocket.Handle("ApplyManualMatching request", bmldgkmanualmatching.ApplyManualMatching)
	routerWebsocket.Handle("UpdateFrequency request", bmldgkmanualmatching.UpdateFrequency)
	routerWebsocket.Handle("FrequencyOptionList request", bmldgkmanualmatching.SendFrequencyOptionList)
	routerWebsocket.Handle("dgkCategoryFilter2OptionList request", bmldgkmanualmatching.SendDgkCategory2FilterOptionList)
	routerWebsocket.Handle("ApplyUnmatch request", bmldgkmanualmatching.ApplyUnmatch)
	routerWebsocket.Handle("ApplySellerUrl request", configcountsupplierpage.ApplySellerUrl)
	routerWebsocket.Handle("ApplyProductUrl request", bmldgkmanualmatching.ApplyProductUrl)
	routerWebsocket.Handle("DgkSellerList request", configcountsupplierpage.ChangeDgkSupplierTable)
	routerWebsocket.Handle("SetGoodMatch request", bmldgktable.SetGoodMatch)
	routerWebsocket.Handle("dgkCategoryAssortment request", configcountcategorypage.ChangeDgkCategoryAssortment)
	routerWebsocket.Handle("bmlCategoryAssortment request", configcountcategorypage.ChangeBmlCategoryAssortment)
	routerWebsocket.Handle("bmlCategoryFilter1OptionList request", configcountcategorypage.SendBmlCategory1FilterOptionList) ////change it
	routerWebsocket.Handle("BmlSellerList request", configcountsupplierpage.ChangeBmlSupplierTable)
	routerWebsocket.Handle("requestSortingFilterOptionList request", bmldgktable.SendSortingOptionList)
	routerWebsocket.Handle("MatchSeller request", configcountsupplierpage.MatchSupplier)
	routerWebsocket.Handle("MatchedSellerListEvent request", configcountsupplierpage.SendMatchedSupplierTable)
	routerWebsocket.Handle("PersonalizedInfo request", homepagewebsockethandler.SendPersonalizedInfo)

	http.Handle("/", routerWebsocket)
	http.ListenAndServe(config.WebsocketIP+`:`+config.WebsocketPort, nil)

}

// ---------------------------------------------------------------------------------------------------------

func startFrontEndServer(config config) {

	gob.Register(bmldgkmanualmatchingrequest.BmlDgkManualMatchingRequest{})
	gob.Register(bmldgkskuhistoryrequest.BmlDgkSkuHistoryRequest{})
	gob.Register(bmldgktablecsvoutputrequest.BmlDgkCsvOutputRequest{})

	routerAPI := gin.Default()

	// creating cookie store
	store := sessions.NewCookieStore([]byte(randToken(64)))
	store.Options(sessions.Options{
		Path:   `/`,
		MaxAge: 86400 * 7,
	})

	// using the cookie store:
	routerAPI.Use(sessions.Sessions(`goquestsession`, store))

	// Serve frontend static files
	routerAPI.Use(static.Serve("/frontend/", static.LocalFile("./frontend", true)))

	routerAPI.GET(`/`, homepage.Start)

	routerAPI.GET(`/login`, login.LoginHandler)
	routerAPI.GET(`/auth`, authenticate.AuthHandler)
	routerAPI.GET(`/unauthorized`, login.UnauthorizedHandler)
	routerAPI.GET(`/adduser`, adduser.Start)
	routerAPI.POST(`/adduser`, adduser.AnswerForm)
	routerAPI.GET(`/adduserconfirmation`, adduser.ConfirmForm)

	routerAPI.GET(`/bmldgkskuhist`, bmldgkskuhistfrontend.Start)
	routerAPI.GET(`/bmldgktable`, bmldgktablefrontend.Start)
	routerAPI.POST(`/seebmldgkskuhist`, bmldgkskuhistfrontend.SeeBmlDgkSkuHistory)
	routerAPI.GET(`/bmldgkskumanualmatching`, manualmatchingfrontend.Start)
	routerAPI.POST(`/bmldgkskumanualmatching`, manualmatchingfrontend.GoToBmlDgkManualMatching)
	routerAPI.GET(`/suppliermatching`, supplierconfigcount.GoToSupplierConfigCountPage)
	routerAPI.GET(`/categoryconfigcount`, categoryconfigcount.GoToCategoryConfigCountPage)

	//routerAPI.POST(`/bmldgktablecsvoutput`, bmldgkcsvoutputfrontend.GoDownloadCsvOutput)
	//routerAPI.GET(`/bmldgktablecsvoutput`, bmldgkcsvoutputfrontend.DownloadCsvOutput)

	routerAPI.Run(config.APIHandlerIP + `:` + config.APIHandlerPort)

}

// randToken returns a random token of i bytes
func randToken(i int) string {
	b := make([]byte, i)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

type config struct {
	APIHandlerIP   string `json:"APIHandlerIP"`
	WebsocketIP    string `json:"WebsocketIP"`
	APIHandlerPort string `json:"APIHandlerPort"`
	WebsocketPort  string `json:"WebsocketPort"`
}

func loadEnvConfig(file string) config {
	var config config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
