package bmldgkskuhist

import (
	"errors"
	"fmt"
	//"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/bmldgkskuhistoryrequest"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

// Start loads the first web page of the application - GET request
func Start(c *gin.Context) {

	authorize.Authorize(c, &user)

	// change here!
	session := sessions.Default(c)
	bmlDgkSkuHistoryRequestI := session.Get("bmlDgkSkuHistoryRequest")
	bmlDgkSkuHistoryRequest, ok := bmlDgkSkuHistoryRequestI.(bmldgkskuhistoryrequest.BmlDgkSkuHistoryRequest)
	fmt.Println()
	fmt.Println(`Object received: `, bmlDgkSkuHistoryRequestI)
	if !ok {
		err := errors.New("Could not get bmlDgkSkuHistoryRequest from session")
		fmt.Println(`Object received: `, bmlDgkSkuHistoryRequestI)
		handleErr(c, err)
	}

	bmlDgkSkuHistoryRequest.Render(c, `frontend/bmldgkskuhist/build/index.html`)

}

// SeeBmlDgkSkuHistory set the user session variable IDBmlDgkSkuHistoryRequest to BmlIDCatalogConfig - POST request
func SeeBmlDgkSkuHistory(c *gin.Context) {

	authorize.Authorize(c, &user)

	session := sessions.Default(c)

	r := c.Request

	bmlDgkSkuHistoryRequest := &bmldgkskuhistoryrequest.BmlDgkSkuHistoryRequest{
		BmlIDCatalogConfig:    r.FormValue("BmlIDCatalogConfig"),
		BmlSKUName:            r.FormValue("BmlSKUName"),
		BmlImgLink:            r.FormValue("BmlImgLink"),
		BmlSupplierName:       r.FormValue("BmlSupplierName"),
		BmlBrand:              r.FormValue("BmlBrand"),
		BmlConfigSnapshot:     r.FormValue("BmlConfigSnapshot"),
		BmlSumOfStockQuantity: r.FormValue("BmlSumOfStockQuantity"),
		BmlSKULink:            r.FormValue("BmlSKULink"),

		DgkIDCatalogConfig: r.FormValue("DgkIDCatalogConfig"),
		DgkSKUName:         r.FormValue("DgkSKUName"),
		DgkImgLink:         r.FormValue("DgkImgLink"),
		DgkSupplierName:    r.FormValue("DgkSupplierName"),
		DgkBrand:           r.FormValue("DgkBrand"),
		DgkConfigSnapshot:  r.FormValue("DgkConfigSnapshot"),
		DgkStock:           r.FormValue("DgkStock"),
		DgkSKULink:         r.FormValue("DgkSKULink"),
	}

	session.Set("bmlDgkSkuHistoryRequest", bmlDgkSkuHistoryRequest)
	err := session.Save()
	handleErr(c, err)

	fmt.Println(bmlDgkSkuHistoryRequest.BmlIDCatalogConfig)

	http.Redirect(c.Writer, r, `/bmldgkskuhist`, http.StatusSeeOther)

}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
