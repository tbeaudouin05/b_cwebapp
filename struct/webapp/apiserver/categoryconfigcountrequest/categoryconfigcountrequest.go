package categoryconfigcountrequest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	//validation "github.com/go-ozzo/ozzo-validation"
)


type ConfigCountRequest struct {
	

	Success string
	Error   string
}

// Validate validates the data of the purchase request sent by the user
/*
func (configCountRequest *ConfigCountRequest) Validate() bool {

	configCountRequest.Error = ""

	// define validation of each field of the purchase request
	err := validation.ValidateStruct(configCountRequest,
		validation.Field(&configCountRequest.BmlIDCatalogConfig, validation.Required),
		validation.Field(&configCountRequest.DgkImgLink, validation.Required),
	)

	// add potential error text to configCountRequest.Error
	if err != nil {
		configCountRequest.Error = err.Error()
	}

	// return true if no error, false otherwise
	return configCountRequest.Error == ""
} */

// Render the web page itself given the html template and the configCountRequest
func (configCountRequest *ConfigCountRequest) Render(c *gin.Context, htmlTemplate string) {
	// fetch the htmlTemplate
	tmpl, err := template.ParseFiles(htmlTemplate)
	handleErr(c, err)
	// render the htmlTemplate given the configCountRequest
	err = tmpl.Execute(c.Writer, map[string]interface{}{
		
		`Success`: configCountRequest.Success,
		`Error`:   configCountRequest.Error,
	})
	handleErr(c, err)
}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
