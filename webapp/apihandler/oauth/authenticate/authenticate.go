package authenticate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"

	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/credential"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var state string
var user useraccess.User

func AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	conf := credential.InitCred()

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	handleErr(c, err)

	client := conf.Client(oauth2.NoContext, tok)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	handleErr(c, err)
	defer userInfo.Body.Close()

	userInfoData, err := ioutil.ReadAll(userInfo.Body)
	handleErr(c, err)
	//log.Println("Email body: ", string(userInfoData))

	err = json.Unmarshal(userInfoData, &user)
	handleErr(c, err)
	c.Status(http.StatusOK)

	session.Set("userEmail", user.Email)
	err = session.Save()
	handleErr(c, err)

	http.Redirect(c.Writer, c.Request, `/`, http.StatusSeeOther)
}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
