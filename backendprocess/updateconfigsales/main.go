package main

import (
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	bobinteract "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract/backendprocess"
	mongointeract "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/backendprocess"
	"github.com/thomas-bamilo/sql/connectdb"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()
	bamiloCatalogConfigSalesTable := bobinteract.GetBmlCatalogConfigSalesTable(dbBob)
	bamiloCatalogConfigSalesHistTable := bobinteract.GetBmlCatalogConfigSalesHistTable(dbBob)

	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	checkError(err)
	defer mongoSession.Close()

	// Optional. Switch the session to a monotonic behavior.
	mongoSession.SetMode(mgo.Monotonic, true)

	var wg sync.WaitGroup
	wg.Add(2)

	go mongointeract.UpsertConfigSales(mongoSession, bamiloCatalogConfigSalesTable, start, &wg)
	go mongointeract.UpsertConfigSalesHist(mongoSession, bamiloCatalogConfigSalesHistTable, start, &wg)

	wg.Wait()

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
