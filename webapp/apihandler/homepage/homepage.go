package homepage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

// Start loads the first web page of the application - GET request
func Start(c *gin.Context) {

	// authorize also populates user with info about the user
	authorize.Authorize(c, &user)

	log.Println(`here10`)

	//user.Success = `Welcome back ` + user.Name + `!`

	log.Println(`here11`)

	user.Render(c, `frontend/homepage/build/index.html`)

}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
