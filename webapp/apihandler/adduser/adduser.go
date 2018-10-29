package adduser

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"

	mongouser "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

// Start loads the purchase request form web page - GET request
func Start(c *gin.Context) {

	authorize.Authorize(c, &user)

	// only pass form as addUserFormInput since we only want a blank form at start
	addUserFormInput := &useraccess.User{}

	// render the web page itself given the html frontend and the addUserFormInput
	addUserFormInput.Render(c, `frontend/adduser/build/adduser.html`)
}

// AnswerForm retrieves user inputs, validate them and upload them to database - POST request
func AnswerForm(c *gin.Context) {

	authorize.Authorize(c, &user)

	r := c.Request

	// pass all the form values input by the user
	// since we want to validate these values and upload them to database
	// in case validation fails, we also want to return these values to the form for good user experience
	addUserFormInput := &useraccess.User{
		Email: r.FormValue(`Email`),
		Name:  r.FormValue(`Name`),
	}

	// Validate validates the addUserFormInput form user inputs
	// if validation fails, reload the purchase request form page with the initial user inputs and error messages
	if addUserFormInput.Validate() == false {

		addUserFormInput.Render(c, `frontend/adduser/build/adduser.html`)
		return
	}

	mongoSession, err := mgo.Dial(`mongodb://localhost:27017/competition_analysis`)
	handleErr(c, err)
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	mongouser.AddUserLogin(addUserFormInput.Email, addUserFormInput.Name, mongoSession)

	// if everything goes well, redirect user to adduserconfirmation web page
	http.Redirect(c.Writer, r, `/adduserconfirmation`, http.StatusSeeOther)
}

// ConfirmForm loads the purchase request adduserconfirmation web page - GET request
func ConfirmForm(c *gin.Context) {

	addUserFormInput := &useraccess.User{}

	// render adduserconfirmation web page
	addUserFormInput.Render(c, `frontend/adduser/build/adduserconfirmation.html`)
}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
