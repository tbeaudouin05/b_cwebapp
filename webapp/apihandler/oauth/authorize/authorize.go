package authorize

import (
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/mgo"

	mongouser "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context, user *useraccess.User) {

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")
	if userEmail == nil {
		http.Redirect(c.Writer, c.Request, `/login`, http.StatusSeeOther)
		return
	}

	user.Email = userEmail.(string)

	mongoSession, err := mgo.Dial(`mongodb://localhost:27017/competition_analysis`)
	handleErr(c, err)
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	mongouser.GetUserAccess(user, mongoSession)

	log.Println(`user.Access: ` + user.Access)
	if user.Access == `` {
		handleErr(c, err)
		http.Redirect(c.Writer, c.Request, `/unauthorized`, http.StatusSeeOther)
		return
	}

	log.Println(`here18`)

}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
