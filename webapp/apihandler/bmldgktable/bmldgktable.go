package bmldgktable

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

// Start loads the first web page of the application - GET request
func Start(c *gin.Context) {

	authorize.Authorize(c, &user)

	user.Render(c, `frontend/bmldgktable/build/index.html`)

}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
