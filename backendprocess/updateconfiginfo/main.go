package main

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	bobinteract "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract/backendprocess"
	"github.com/thomas-bamilo/sql/connectdb"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	dbBob := connectdb.ConnectToBob()
	dbBob.SetMaxIdleConns(0)
	defer dbBob.Close()

	mongoSession, err := mgo.Dial(`mongodb://localhost:27017/competition_analysis`)
	checkError(err)
	defer mongoSession.Close()

	bobinteract.UpdateConfigInfoAndHist(dbBob, mongoSession)

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
