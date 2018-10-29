package manualmatching

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/bmldgkmanualmatchingrequest"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

// Start loads the first web page of the application - GET request
func Start(c *gin.Context) {

	authorize.Authorize(c, &user)

	session := sessions.Default(c)

	bmlDgkManualMatchingRequestI := session.Get("bmlDgkManualMatchingRequest")
	bmlDgkManualMatchingRequest, ok := bmlDgkManualMatchingRequestI.(bmldgkmanualmatchingrequest.BmlDgkManualMatchingRequest)
	if !ok {
		err := errors.New("Could not get bmlDgkManualMatchingRequest from session")
		log.Println(`Object received: `, bmlDgkManualMatchingRequestI)
		handleErr(c, err)
	}

	bmlDgkManualMatchingRequest.Email = user.Email
	bmlDgkManualMatchingRequest.Name = user.Name

	bmlDgkManualMatchingRequest.Render(c, `frontend/bmldgkskumanualmatch/build/index.html`)

}

// GoToBmlDgkManualMatching set the user session variable IDBmlDgkSkuHistoryRequest to BmlIDCatalogConfig - POST request
func GoToBmlDgkManualMatching(c *gin.Context) {

	authorize.Authorize(c, &user)

	session := sessions.Default(c)

	r := c.Request

	bmlDgkManualMatchingRequest := &bmldgkmanualmatchingrequest.BmlDgkManualMatchingRequest{
		BmlIDCatalogConfig: r.FormValue("BmlIDCatalogConfig"),
		DgkImgLink:         r.FormValue("DgkImgLink"),
	}

	session.Set("bmlDgkManualMatchingRequest", bmlDgkManualMatchingRequest)
	err := session.Save()
	handleErr(c, err)

	log.Println(bmlDgkManualMatchingRequest.BmlIDCatalogConfig)

	http.Redirect(c.Writer, r, `/bmldgkskumanualmatching`, http.StatusSeeOther)

}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
