package supplierconfigcount

import (
	"github.com/gin-gonic/gin"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/supplierconfigcountrequest"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

func GoToSupplierConfigCountPage(c *gin.Context) {

	authorize.Authorize(c, &user)

	configCountRequest := supplierconfigcountrequest.ConfigCountRequest{
		Email: user.Email,
		Name:  user.Name,
	}

	configCountRequest.Render(c, `frontend/suppliermatching/build/index.html`)

}
