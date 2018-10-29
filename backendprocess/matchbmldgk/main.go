package main

import (
	"context"
	"log"
	"math"
	"os"
	"sync"

	"time"

	"github.com/globalsign/mgo"
	elasticinteract "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/elasticinteract/backendprocess"
	mongointeract "github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract/backendprocess"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	bmlconfigbackendprocess "github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig/backendprocess"
	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	checkError(err)
	defer mongoSession.Close()

	var bamiloCatalogConfigTableToMatch []bmlcatalogconfig.BmlCatalogConfig
	mongointeract.GetBmlCatalogConfigTableToMatch(mongoSession, &bamiloCatalogConfigTableToMatch)
	log.Println(`bamiloCatalogConfigTableToMatch: `, len(bamiloCatalogConfigTableToMatch))

	if len(bamiloCatalogConfigTableToMatch) == 0 {
		mongointeract.SetMatchedToFalse(mongoSession)
	}

	arrayOfBmlCatalogConfigTableRandomChunk := bmlcatalogconfig.DivideBmlCatalogConfigTableInRandomChunk(bamiloCatalogConfigTableToMatch)

	log.Println(`arrayOfBmlCatalogConfigTableRandomChunk`, len(arrayOfBmlCatalogConfigTableRandomChunk))

	elasticClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	checkError(err)
	ctx := context.Background()

	var wg sync.WaitGroup
	var arrayOfBmlCatalogConfigMatchTableChunk [][]bmlconfigbackendprocess.Match
	for i := 0; i < len(arrayOfBmlCatalogConfigTableRandomChunk); i++ {
		arrayOfBmlCatalogConfigMatchTableChunk = append(arrayOfBmlCatalogConfigMatchTableChunk, []bmlconfigbackendprocess.Match{})
	}
	log.Println(`arrayOfBmlCatalogConfigMatchTableChunk`, len(arrayOfBmlCatalogConfigMatchTableChunk))

	limit := math.Max(float64(len(arrayOfBmlCatalogConfigTableRandomChunk)), 20)
	for i := 0; i < int(limit); i++ {
		arrayOfBmlCatalogConfigMatchTableChunk[i] = elasticinteract.MatchBmlDgkConfig(elasticClient, ctx, arrayOfBmlCatalogConfigTableRandomChunk[i])
		log.Println(`arrayOfBmlCatalogConfigMatchTableChunk[i]`, len(arrayOfBmlCatalogConfigMatchTableChunk[i]))
		wg.Add(1)
		go mongointeract.UpsertConfigMatch(mongoSession, arrayOfBmlCatalogConfigMatchTableChunk[i], start, &wg)
		wg.Wait()

	}

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
