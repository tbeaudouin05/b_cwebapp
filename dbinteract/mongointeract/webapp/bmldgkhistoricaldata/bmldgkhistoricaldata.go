package bmldgkhistoricaldata

import (
	"fmt"
"sync"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgkskuhistory"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)

func FetchBmlDgkHistoricalData(bmlID int,
	BmlSumOfStockQuantity int,BmlSKUName string,BmlImgLink string,BmlSupplierName string,
	BmlBrand string,BmlConfigSnapshot string ,BmlSKULink string,DgkIDCatalogConfig string,
	DgkSKUName string,DgkImgLink string,DgkSupplierName string,DgkBrand string,DgkConfigSnapshot string,
	DgkStock string,DgkSKULink string, mongoSession *mgo.Session) bmldgktable.HistoricalChartData {

	var historicalChartData bmldgktable.HistoricalChartData




	historicalChartData.SKUName.BmlSKUName = BmlSKUName
	historicalChartData.SKUName.DgkSKUName = DgkSKUName

	var wg sync.WaitGroup
	wg.Add(2)
	var bamiloCatalogConfigHistStr []bmldgkskuhistory.BmlDgkConfigHist
	go findBmlCatalog(bmlID ,  mongoSession , &bamiloCatalogConfigHistStr,&wg)

	var dgkCatalogConfigHistStr []dgkcatalogconfig.DgkCatalogConfigHist
	go findDgkCatalog(DgkIDCatalogConfig,  mongoSession ,  &dgkCatalogConfigHistStr ,&wg)
	wg.Wait()

	initialPrice := -3
	
	for i, tmpBmlCatalogConfigHistStr := range bamiloCatalogConfigHistStr {
		historicalChartData.Value.DgkPrice = append(historicalChartData.Value.DgkPrice, initialPrice)
		tmpPrice := tmpBmlCatalogConfigHistStr.AvgSpecialPrice
		if (tmpPrice == 0){
			tmpPrice = tmpBmlCatalogConfigHistStr.AvgPrice
		}
		historicalChartData.Value.BmlPrice = append(historicalChartData.Value.BmlPrice, tmpPrice/10000)
		historicalChartData.Value.BmlSales = append(historicalChartData.Value.BmlSales, tmpBmlCatalogConfigHistStr.SumOfPaidPrice/10000)
		// instead of if for each loop, initialize and end the array, then apply other rules inside the array
		if i == 0 || i == len(bamiloCatalogConfigHistStr)-1 {
			historicalChartData.Label = append(historicalChartData.Label, tmpBmlCatalogConfigHistStr.ConfigSnapshotAt.Format("01/02 2006"))
		} else {
			historicalChartData.Label = append(historicalChartData.Label, "")
		}
	}
	if len(dgkCatalogConfigHistStr) > 0 {
		cnt := 0
		for i, tmpBmlCatalogConfigHistStr := range bamiloCatalogConfigHistStr {
			if cnt < len(dgkCatalogConfigHistStr) {
				if dgkCatalogConfigHistStr[cnt].ConfigSnapshotAt.Format("01/02 2006") == tmpBmlCatalogConfigHistStr.ConfigSnapshotAt.Format("01/02 2006") {	
					tmpPrice := dgkCatalogConfigHistStr[cnt].AvgSpecialPrice
					if (tmpPrice == 0){
						tmpPrice = dgkCatalogConfigHistStr[cnt].AvgPrice
					}
					historicalChartData.Value.DgkPrice[i] = tmpPrice/10000
					cnt++
				}
			} else {
				break
			}

		}
	}
	n := len(historicalChartData.Value.DgkPrice) - 1
	for i := 0; i <= n; i++ {
		if historicalChartData.Value.DgkPrice[i] == initialPrice && i > 0 {
			historicalChartData.Value.DgkPrice[i] = historicalChartData.Value.DgkPrice[i-1]
		}
	}
	for i := n; i >= 0; i-- {
		if historicalChartData.Value.DgkPrice[i] == initialPrice && i < n {
			historicalChartData.Value.DgkPrice[i] = historicalChartData.Value.DgkPrice[i+1]
		}
	}

	return historicalChartData
}


func findDgkCatalog( DgkIDCatalogConfig string,  mongoSession *mgo.Session ,dgkCatalogConfigHistStr *[]dgkcatalogconfig.DgkCatalogConfigHist, wg *sync.WaitGroup){
	defer wg.Done()
	stageSort := bson.M{"$sort": bson.M{"-config_snapshot_at": -1}}

	db := mongoSession.DB("competition_analysis")
	dgkCatalogConfigHist := db.C("dgk_catalog_config_hist")

	findDgkCatalogConfigHist := bson.M{"$match": bson.M{"fk_dgk_catalog_config": DgkIDCatalogConfig}}
	pipe := dgkCatalogConfigHist.Pipe([]bson.M{findDgkCatalogConfigHist, stageSort})
	err := pipe.All(dgkCatalogConfigHistStr)
	fmt.Println(dgkCatalogConfigHistStr)
	if err != nil {
		fmt.Println("err", err)
	}
}
func findBmlCatalog( bmlID int,  mongoSession *mgo.Session,bamiloCatalogConfigHistStr *[]bmldgkskuhistory.BmlDgkConfigHist , wg *sync.WaitGroup){
	defer wg.Done()

	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfigHist := db.C("bml_catalog_config_hist")

	stageSort := bson.M{"$sort": bson.M{"-config_snapshot_at": -1}}

	findBmlCatalogConfigHist := bson.M{"$match": bson.M{"fk_bml_catalog_config": bmlID}}
	pipe := bmlCatalogConfigHist.Pipe([]bson.M{findBmlCatalogConfigHist, stageSort})
	err := pipe.All(bamiloCatalogConfigHistStr)
	if err != nil {
		fmt.Println("err", err)
	}
}