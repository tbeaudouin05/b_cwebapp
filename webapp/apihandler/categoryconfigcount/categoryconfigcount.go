package categoryconfigcount

import (
	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/categoryconfigcountrequest"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

func GoToCategoryConfigCountPage(c *gin.Context) {

	authorize.Authorize(c, &user)

	var configCountRequest categoryconfigcountrequest.ConfigCountRequest
	configCountRequest.Render(c, `frontend/categorypage/build/index.html`)

}
