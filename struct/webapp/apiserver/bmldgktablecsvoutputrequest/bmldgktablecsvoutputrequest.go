package bmldgktablecsvoutputrequest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type BmlDgkCsvOutputRequest struct {
	NumberOfRow string
	Category1   string
	Category2   string
	Category3   string
	SearchedBy  string
	PageNumber  string
	SortedBy    string

	Success string
	Error   string
}

// Validate validates the data of the purchase request sent by the user
func (bmlDgkCsvOutputRequest *BmlDgkCsvOutputRequest) Validate() bool {

	bmlDgkCsvOutputRequest.Error = ""

	// define validation of each field of the purchase request
	err := validation.ValidateStruct(bmlDgkCsvOutputRequest,
		validation.Field(&bmlDgkCsvOutputRequest.NumberOfRow, validation.Required),
		validation.Field(&bmlDgkCsvOutputRequest.Category1, validation.Required),
		validation.Field(&bmlDgkCsvOutputRequest.SearchedBy, validation.Required),
	)

	// add potential error text to bmlDgkCsvOutputRequest.Error
	if err != nil {
		bmlDgkCsvOutputRequest.Error = err.Error()
	}

	// return true if no error, false otherwise
	return bmlDgkCsvOutputRequest.Error == ""
}

// Render the web page itself given the html template and the bmlDgkCsvOutputRequest
func (bmlDgkCsvOutputRequest *BmlDgkCsvOutputRequest) Render(c *gin.Context, htmlTemplate string) {
	// fetch the htmlTemplate
	tmpl, err := template.ParseFiles(htmlTemplate)
	handleErr(c, err)
	// render the htmlTemplate given the bmlDgkCsvOutputRequest
	err = tmpl.Execute(c.Writer, map[string]interface{}{
		`NumberOfRow`: bmlDgkCsvOutputRequest.NumberOfRow,
		`Category1`:   bmlDgkCsvOutputRequest.Category1,
		`Category2`:   bmlDgkCsvOutputRequest.Category2,
		`Category3`:   bmlDgkCsvOutputRequest.Category3,
		`SearchedBy`:  bmlDgkCsvOutputRequest.SearchedBy,
		`PageNumber`:  bmlDgkCsvOutputRequest.PageNumber,
		`SortedBy`:    bmlDgkCsvOutputRequest.SortedBy,

		`Success`: bmlDgkCsvOutputRequest.Success,
		`Error`:   bmlDgkCsvOutputRequest.Error,
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
