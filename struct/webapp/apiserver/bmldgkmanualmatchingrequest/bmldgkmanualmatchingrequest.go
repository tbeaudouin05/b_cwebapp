package bmldgkmanualmatchingrequest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type BmlDgkManualMatchingRequest struct {
	BmlIDCatalogConfig string
	Email              string
	Name               string
	DgkImgLink         string

	Success string
	Error   string
}

// Validate validates the data of the purchase request sent by the user
func (bmlDgkManualMatchingRequest *BmlDgkManualMatchingRequest) Validate() bool {

	bmlDgkManualMatchingRequest.Error = ""

	// define validation of each field of the purchase request
	err := validation.ValidateStruct(bmlDgkManualMatchingRequest,
		validation.Field(&bmlDgkManualMatchingRequest.BmlIDCatalogConfig, validation.Required),
		validation.Field(&bmlDgkManualMatchingRequest.DgkImgLink, validation.Required),
	)

	// add potential error text to bmlDgkManualMatchingRequest.Error
	if err != nil {
		bmlDgkManualMatchingRequest.Error = err.Error()
	}

	// return true if no error, false otherwise
	return bmlDgkManualMatchingRequest.Error == ""
}

// Render the web page itself given the html template and the bmlDgkManualMatchingRequest
func (bmlDgkManualMatchingRequest *BmlDgkManualMatchingRequest) Render(c *gin.Context, htmlTemplate string) {
	// fetch the htmlTemplate
	tmpl, err := template.ParseFiles(htmlTemplate)
	handleErr(c, err)
	// render the htmlTemplate given the bmlDgkManualMatchingRequest
	err = tmpl.Execute(c.Writer, map[string]interface{}{
		`BmlIDCatalogConfig`: bmlDgkManualMatchingRequest.BmlIDCatalogConfig,
		`Email`:              bmlDgkManualMatchingRequest.Email,
		`Name`:               bmlDgkManualMatchingRequest.Name,
		`DgkImgLink`:         bmlDgkManualMatchingRequest.DgkImgLink,

		`Success`: bmlDgkManualMatchingRequest.Success,
		`Error`:   bmlDgkManualMatchingRequest.Error,
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
