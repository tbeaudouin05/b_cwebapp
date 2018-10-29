package bmldgkskuhistoryrequest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type BmlDgkSkuHistoryRequest struct {
	BmlIDCatalogConfig    string
	BmlSKUName            string
	BmlImgLink            string
	BmlSupplierName       string
	BmlBrand              string
	BmlConfigSnapshot     string
	BmlSumOfStockQuantity string
	BmlSKULink            string

	DgkIDCatalogConfig string
	DgkSKUName         string
	DgkImgLink         string
	DgkSupplierName    string
	DgkBrand           string
	DgkConfigSnapshot  string
	DgkStock           string
	DgkSKULink         string

	Success string
	Error   string
}

// Validate validates the data of the purchase request sent by the user
func (bmlDgkSkuHistoryRequest *BmlDgkSkuHistoryRequest) Validate() bool {

	bmlDgkSkuHistoryRequest.Error = ""

	// define validation of each field of the purchase request
	err := validation.ValidateStruct(bmlDgkSkuHistoryRequest,
		validation.Field(&bmlDgkSkuHistoryRequest.BmlIDCatalogConfig, validation.Required),
		validation.Field(&bmlDgkSkuHistoryRequest.DgkIDCatalogConfig, validation.Required),
	)

	// add potential error text to bmlDgkSkuHistoryRequest.Error
	if err != nil {
		bmlDgkSkuHistoryRequest.Error = err.Error()
	}

	// return true if no error, false otherwise
	return bmlDgkSkuHistoryRequest.Error == ""
}

// Render the web page itself given the html template and the bmlDgkSkuHistoryRequest
func (bmlDgkSkuHistoryRequest *BmlDgkSkuHistoryRequest) Render(c *gin.Context, htmlTemplate string) {
	// fetch the htmlTemplate
	tmpl, err := template.ParseFiles(htmlTemplate)
	handleErr(c, err)
	// render the htmlTemplate given the bmlDgkSkuHistoryRequest
	err = tmpl.Execute(c.Writer, map[string]interface{}{
		`BmlIDCatalogConfig`:    bmlDgkSkuHistoryRequest.BmlIDCatalogConfig,
		`BmlSKUName`:            bmlDgkSkuHistoryRequest.BmlSKUName,
		`BmlImgLink`:            bmlDgkSkuHistoryRequest.BmlImgLink,
		`BmlSupplierName`:       bmlDgkSkuHistoryRequest.BmlSupplierName,
		`BmlBrand`:              bmlDgkSkuHistoryRequest.BmlBrand,
		`BmlConfigSnapshot`:     bmlDgkSkuHistoryRequest.BmlConfigSnapshot,
		`BmlSumOfStockQuantity`: bmlDgkSkuHistoryRequest.BmlSumOfStockQuantity,
		`BmlSKULink`:            bmlDgkSkuHistoryRequest.BmlSKULink,

		`DgkIDCatalogConfig`: bmlDgkSkuHistoryRequest.DgkIDCatalogConfig,
		`DgkSKUName`:         bmlDgkSkuHistoryRequest.DgkSKUName,
		`DgkImgLink`:         bmlDgkSkuHistoryRequest.DgkImgLink,
		`DgkSupplierName`:    bmlDgkSkuHistoryRequest.DgkSupplierName,
		`DgkBrand`:           bmlDgkSkuHistoryRequest.DgkBrand,
		`DgkConfigSnapshot`:  bmlDgkSkuHistoryRequest.DgkConfigSnapshot,
		`DgkStock`:           bmlDgkSkuHistoryRequest.DgkStock,
		`DgkSKULink`:         bmlDgkSkuHistoryRequest.DgkSKULink,

		`Success`: bmlDgkSkuHistoryRequest.Success,
		`Error`:   bmlDgkSkuHistoryRequest.Error,
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
