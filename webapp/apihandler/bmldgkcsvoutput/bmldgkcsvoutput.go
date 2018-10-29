package bmldgkcsvoutput

import (
	"bytes"
	"encoding/csv"

	"fmt"
	"net/http"

	"errors"

	"log"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/useraccess"

	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/apiserver/bmldgktablecsvoutputrequest"
	bmldgktablestruct "github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/thomas-bamilo/commercial/competitionanalysis/webapp/apihandler/oauth/authorize"
)

var user useraccess.User

func GoDownloadCsvOutput(c *gin.Context) {

	authorize.Authorize(c, &user)

	session := sessions.Default(c)

	r := c.Request

	bmlDgkCsvOutputRequest := &bmldgktablecsvoutputrequest.BmlDgkCsvOutputRequest{
		NumberOfRow: r.FormValue("NumberOfRow"),
		Category1:   r.FormValue("Category1"),
		Category2:   r.FormValue("Category2"),
		Category3:   r.FormValue("Category3"),
		SearchedBy:  r.FormValue("SearchedBy"),
		PageNumber:  r.FormValue("PageNumber"),
		SortedBy:    r.FormValue("SortedBy"),
	}
	/////////////////////////////validate here

	fmt.Println("GoDownloadCsvOutput CHECK!")
	fmt.Println("numberOfRow11", bmlDgkCsvOutputRequest.NumberOfRow)
	fmt.Println("numberOfRowrr", r.FormValue("NumberOfRow"))
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	session.Set("bmlDgkCsvOutputRequest", bmlDgkCsvOutputRequest)
	err := session.Save()
	handleErr(c, err)

	http.Redirect(c.Writer, r, `/bmldgktablecsvoutput`, http.StatusSeeOther)

}

func DownloadCsvOutput(c *gin.Context) {

	authorize.Authorize(c, &user)

	fmt.Println(" DownloadCsvOutput check")

	session := sessions.Default(c)
	bmlDgkCsvOutputRequestI := session.Get("bmlDgkCsvOutputRequest")
	fmt.Println("start", bmlDgkCsvOutputRequestI)
	bmlDgkCsvOutputRequest, ok := bmlDgkCsvOutputRequestI.(bmldgktablecsvoutputrequest.BmlDgkCsvOutputRequest)
	if !ok {
		err := errors.New("Could not get bmlDgkCsvOutputRequest from session")
		log.Println(`Object received: `, bmlDgkCsvOutputRequestI)
		handleErr(c, err)
	}

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()
	mongoSession.SetMode(mgo.Monotonic, true)

	// building table
	var table bmldgktablestruct.Table
	var tableIsNull bool
	numberOfRow, err := strconv.Atoi(bmlDgkCsvOutputRequest.NumberOfRow)
	pageNumber, err1 := strconv.Atoi(bmlDgkCsvOutputRequest.PageNumber)
	checkErr(err)
	checkErr(err1)
	/*fmt.Println("numberOfRow", bmlDgkCsvOutputRequest.NumberOfRow, " pageNumber", bmlDgkCsvOutputRequest.PageNumber, " SearchedBy", bmlDgkCsvOutputRequest.SearchedBy, " Category",
		bmlDgkCsvOutputRequest.Category, " sortedBy", bmlDgkCsvOutputRequest.SortedBy)
	fmt.Println("numberOfRow err:", err)*/
	table, tableIsNull = bmldgktable.SelectCompetitionAnalysisTable(numberOfRow, pageNumber, bmlDgkCsvOutputRequest.SearchedBy, bmlDgkCsvOutputRequest.Category1, bmlDgkCsvOutputRequest.Category2, bmlDgkCsvOutputRequest.Category3, bmlDgkCsvOutputRequest.SortedBy, mongoSession)
	fmt.Println(tableIsNull)
	//fmt.Println(table)

	//check if table ise null
	if tableIsNull {
		table = bmldgktable.SetNullTabel()
	}
	csvOutput := table

	//get csv information
	//////write handler here

	// download csv
	var csvOutputStr [][]string

	csvOutputStr = append(csvOutputStr, []string{`BmlIDCatalogConfig`,
		`BmlSKUName`,
		`BmlImgLink`,
		`BmlSKULink`,
		`BmlSKUPrice`,
		`DgkIDCatalogConfig`,
		`DgkScore`,
		`DgkSKUName`,
		`DgkImgLink`,
		`DgkSKULink`,
		`DgkSKUPrice`,
		`BmlConfigSnapshotAt`,
		`DgkConfigSnapshotAt`,
		`BmlSupplierName`,
		`BmlBrand`,
		`BmlMinOfStockQuantity`,
		`BmlSumOfStockQuantity`,
		`DgkStock`,
	})

	for i := 0; i < len(csvOutput.Table); i++ {

		csvOutputStr = append(csvOutputStr, []string{
			csvOutput.Table[i].RowValue.BmlIDCatalogConfig,
			csvOutput.Table[i].RowValue.BmlSKUName,
			csvOutput.Table[i].RowValue.BmlImgLink,
			csvOutput.Table[i].RowValue.BmlSKULink,
			csvOutput.Table[i].RowValue.BmlSKUPrice,
			csvOutput.Table[i].RowValue.DgkIDCatalogConfig,
			strconv.Itoa(csvOutput.Table[i].RowValue.DgkScore),
			csvOutput.Table[i].RowValue.DgkSKUName,
			csvOutput.Table[i].RowValue.DgkImgLink,
			csvOutput.Table[i].RowValue.DgkSKULink,
			csvOutput.Table[i].RowValue.DgkSKUPrice,
			csvOutput.Table[i].RowValue.BmlConfigSnapshotAt,
			csvOutput.Table[i].RowValue.DgkConfigSnapshotAt,
			csvOutput.Table[i].RowValue.BmlSupplierName,
			csvOutput.Table[i].RowValue.BmlBrand,
			csvOutput.Table[i].RowValue.BmlMinOfStockQuantity,
			csvOutput.Table[i].RowValue.BmlSumOfStockQuantity,
			csvOutput.Table[i].RowValue.DgkStock,
		})
	}

	fmt.Println("table", table)
	b := &bytes.Buffer{}   // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.

	for i := 0; i < len(csvOutputStr); i++ {
		wr.Write(csvOutputStr[i]) // converts array of string to comma seperated values for 1 row.
	}
	wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=bmldgktable.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func handleErr(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(`Error: %v`, err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func checkErr(err error) {
	if err != nil {
		fmt.Println("error", err)
	}
}
