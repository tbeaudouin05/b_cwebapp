package bmldgktable

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pmylund/sortutil"

	"github.com/SepidehKHH/formatting"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)

func SelectCompetitionAnalysisTable(numOfRow int, pgNumber int, searchedBy string, biCat string, cat1 string, cat2 string, cat3 string, sortedBy string, mongoSession *mgo.Session) (bmldgktable.Table, bool) {
	var selectedResult []bmlcatalogconfig.BmlCatalogConfig

	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")
	dgkCatalogConfig := db.C("dgk_catalog_config")

	// draft for filtering //
	//stageFilterForDgkScore := bson.M{"$match": bson.M{"dgk_score":bson.M{"$gte":0}}}
	//check which queries are required
	wrongSortFlag := false
	if sortedBy == "SkuRank" {
		sortedBy = "-count_of_soi"
		wrongSortFlag = true
	}

	biCategoryCheck := bson.M{"bi_category": biCat}
	category1Check := bson.M{"bi_category_one_name": cat1}
	category2Check := bson.M{"bi_category_two_name": cat2}
	category3Check := bson.M{"bi_category_three_name": cat3}

	if searchedBy == "" {

		if cat3 != "" {
			log.Println(`sortedBy`, sortedBy)
			bmlCatalogConfig.Find(category3Check).Sort(sortedBy).Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else if cat2 != "" {
			bmlCatalogConfig.Find(category2Check).Sort(sortedBy).Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else if cat1 != "" {
			bmlCatalogConfig.Find(category1Check).Sort(sortedBy).Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else if biCat != "" {
			bmlCatalogConfig.Find(biCategoryCheck).Sort(sortedBy).Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else {
			bmlCatalogConfig.Find(nil).Sort(sortedBy).Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		}
	} else {
		if cat3 != "" {
			bmlCatalogConfig.Find(bson.M{"$and": []bson.M{category3Check},
				"$text": bson.M{"$search": searchedBy},
			}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else if cat2 != "" {
			bmlCatalogConfig.Find(bson.M{"$and": []bson.M{category2Check},
				"$text": bson.M{"$search": searchedBy},
			}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else if cat1 != "" {
			bmlCatalogConfig.Find(bson.M{"$and": []bson.M{category1Check},
				"$text": bson.M{"$search": searchedBy},
			}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else if biCat != "" {
			bmlCatalogConfig.Find(bson.M{"$and": []bson.M{biCategoryCheck},
				"$text": bson.M{"$search": searchedBy},
			}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		} else {
			bmlCatalogConfig.Find(bson.M{
				"$text": bson.M{"$search": searchedBy},
			}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(numOfRow*pgNumber - numOfRow).Limit(numOfRow).All(&selectedResult)
		}

	}

	var checkIfTableIsNull bool
	if selectedResult == nil {
		checkIfTableIsNull = true
	}

	log.Println(`selectedResult`, selectedResult)

	var matchedDgkCatalogConfig []dgkcatalogconfig.DgkCatalogConfig
	for _, bmlCatalogConfigStruct := range selectedResult {

		var tmpMatchedDgkCatalogConfig dgkcatalogconfig.DgkCatalogConfig
		var choosedFK string
		if len(bmlCatalogConfigStruct.ManualFKBmlCatalogConfig) > 0 {
			choosedFK = bmlCatalogConfigStruct.ManualFKBmlCatalogConfig
		} else {
			choosedFK = bmlCatalogConfigStruct.FKBmlCatalogConfig
		}
		findDgkCatalogConfig := bson.M{"$match": bson.M{"_id": choosedFK}}
		pipe := dgkCatalogConfig.Pipe([]bson.M{findDgkCatalogConfig})
		err := pipe.One(&tmpMatchedDgkCatalogConfig)
		checkError(err)

		matchedDgkCatalogConfig = append(matchedDgkCatalogConfig, tmpMatchedDgkCatalogConfig)

	}

	selectedResultInBmlDgkStruct := convertCatalogConfigMatchToBmlDgkStruct(selectedResult, matchedDgkCatalogConfig)
	//fmt.Println(selectedResultInBmlDgkStruct)
	if wrongSortFlag {
		sortutil.AscByField(selectedResultInBmlDgkStruct, "SkuRank")
		for _, rng1 := range selectedResultInBmlDgkStruct {
			fmt.Println(rng1.SkuRank)
		}
	}
	//fmt.Println(selectedResultInBmlDgkStruct)

	selectedResultInTable := convertBmlDgkStructToTable(selectedResultInBmlDgkStruct)
	if checkIfTableIsNull == true {
		return selectedResultInTable, true
	} else {
		return selectedResultInTable, false
	}
}

func convertCatalogConfigMatchToBmlDgkStruct(inputArr []bmlcatalogconfig.BmlCatalogConfig, matchedDgkCatalogConfig []dgkcatalogconfig.DgkCatalogConfig) []bmldgktable.BmlDgkConfigTable {
	var outPut []bmldgktable.BmlDgkConfigTable
	for i, inputRow := range inputArr {
		var tempNewStruct bmldgktable.BmlDgkConfigTable
		//fmt.Println(inputRow.ConfigSnapshotAt) // change time to Tehran time zone
		inputRow.ConfigSnapshotAt = inputRow.ConfigSnapshotAt.Add(time.Hour*4 + time.Minute*30 + time.Second*0)
		//fmt.Println(inputRow.ConfigSnapshotAt)
		//fmt.Println("--------")
		secTempNewStruct := bmldgktable.BmlDgkConfigTable{
			SKU:                     inputRow.SKU,
			BamiloImgLink:           inputRow.ImgLink,
			BamiloSKULink:           inputRow.SKULink,
			ProductWarranty:         inputRow.ProductWarranty,
			BrandName:               inputRow.BrandName,
			BrandNameEn:             inputRow.BrandNameEn,
			SupplierName:            inputRow.SupplierName,
			SupplierNameEn:          inputRow.SupplierNameEn,
			VisibleInShop:           inputRow.VisibleInShop,
			BamiloAvgPrice:          strconv.Itoa(inputRow.AvgPrice),
			BamiloAvgSpecialPrice:   strconv.Itoa(inputRow.AvgSpecialPrice),
			SumOfStockQuantity:      strconv.Itoa(inputRow.SumOfStockQuantity),
			MinOfStockQuantity:      strconv.Itoa(inputRow.MinOfStockQuantity),
			BamiloCategoryOneName:   inputRow.BiCategoryOneName,
			BamiloCategoryTwoName:   inputRow.BiCategoryTwoName,
			BamiloCategoryThreeName: inputRow.BiCategoryThreeName,
			BamiloSKUName:           inputRow.SKUName,
			SumOfPaidPrice:          strconv.Itoa(inputRow.SumOfPaidPrice),
			IDBmlCatalogConfig:      inputRow.IDBmlCatalogConfig,
			DgkScore:                inputRow.DgkScore,
			BmlConfigSnapshotAt:     inputRow.ConfigSnapshotAt.Format("01/02 2006 at 3:04pm"),
			GoodMatch:               inputRow.GoodMatch,
		}

		if secTempNewStruct.BamiloAvgPrice == secTempNewStruct.BamiloAvgSpecialPrice {
			secTempNewStruct.BamiloAvgPrice = ""
		}

		tmpSpecialPrice := secTempNewStruct.BamiloAvgSpecialPrice
		if tmpSpecialPrice == "" || tmpSpecialPrice == "0" {
			secTempNewStruct.BamiloPrice = secTempNewStruct.BamiloAvgPrice
		} else {
			secTempNewStruct.BamiloPrice = tmpSpecialPrice
		}
		if len(matchedDgkCatalogConfig) > 0 {
			secTempNewStruct.DgkAvgPrice = strconv.Itoa(matchedDgkCatalogConfig[i].AvgPrice)
			secTempNewStruct.DgkAvgSpecialPrice = strconv.Itoa(matchedDgkCatalogConfig[i].AvgSpecialPrice)
			//fmt.Println("DgkAvgPrice",secTempNewStruct.DgkAvgPrice)
			//fmt.Println("DgkAvgSpecialPrice",secTempNewStruct.DgkAvgSpecialPrice)
			if secTempNewStruct.DgkAvgPrice == secTempNewStruct.DgkAvgSpecialPrice {
				secTempNewStruct.DgkAvgPrice = ""
			}
			secTempNewStruct.IDDgkCatalogConfig = matchedDgkCatalogConfig[i].IDDgkCatalogConfig
			secTempNewStruct.DgkCategoryOneName = matchedDgkCatalogConfig[i].DgkCategoryOneName
			secTempNewStruct.DgkCategoryTwoName = matchedDgkCatalogConfig[i].DgkCategoryTwoName
			secTempNewStruct.IsOutOfStock = strconv.FormatBool(matchedDgkCatalogConfig[i].IsOutOfStock)
			secTempNewStruct.Rating = strconv.Itoa(matchedDgkCatalogConfig[i].Rating)
			secTempNewStruct.DgkSKUName = matchedDgkCatalogConfig[i].SkuName
			secTempNewStruct.SkuRank = matchedDgkCatalogConfig[i].SkuRank
			secTempNewStruct.DgkSKULink = matchedDgkCatalogConfig[i].SkuLink
			secTempNewStruct.DgkImgLink = matchedDgkCatalogConfig[i].ImgLink
			secTempNewStruct.DgkConfigSnapshotAt = matchedDgkCatalogConfig[i].ConfigSnapshotAt.Format("01/02 2006 at 3:04pm")
		} else {
			secTempNewStruct.DgkAvgPrice = "NoMatch"
			secTempNewStruct.DgkAvgSpecialPrice = "NoMatch"
			secTempNewStruct.IDDgkCatalogConfig = "NoMatch"
			secTempNewStruct.DgkCategoryOneName = "NoMatch"
			secTempNewStruct.DgkCategoryTwoName = "NoMatch"
			secTempNewStruct.IsOutOfStock = "NoMatch"
			secTempNewStruct.Rating = "NoMatch"
			secTempNewStruct.DgkSKUName = "NoMatch"
			secTempNewStruct.SkuRank = 0
			secTempNewStruct.DgkSKULink = "NoMatch"
			secTempNewStruct.DgkSKULink = "NoMatch"
			secTempNewStruct.DgkImgLink = "NoMatch"
			secTempNewStruct.DgkConfigSnapshotAt = "NoMatch"
		}
		tmpSpecialPrice = secTempNewStruct.DgkAvgSpecialPrice
		if tmpSpecialPrice == "" || tmpSpecialPrice == "0" {
			secTempNewStruct.DgkPrice = secTempNewStruct.DgkAvgPrice
		} else {
			secTempNewStruct.DgkPrice = tmpSpecialPrice
		}
		tempNewStruct = secTempNewStruct
		outPut = append(outPut, tempNewStruct)
	}
	return outPut
}

func convertBmlDgkStructToTable(competitionAnalysisArr []bmldgktable.BmlDgkConfigTable) bmldgktable.Table {
	var arrOfRow []bmldgktable.Row
	n := len(competitionAnalysisArr)
	for i := 1; i < n+1; i++ {
		var row bmldgktable.Row
		row.RowKey = i

		row.RowValue.BmlIDCatalogConfig = strconv.Itoa(competitionAnalysisArr[i-1].IDBmlCatalogConfig)

		row.RowValue.BmlSKUName = competitionAnalysisArr[i-1].BamiloSKUName

		row.RowValue.BmlImgLink = competitionAnalysisArr[i-1].BamiloImgLink

		row.RowValue.BmlSKULink = competitionAnalysisArr[i-1].BamiloSKULink

		row.RowValue.BmlSKUPrice = formatting.ChangeNumberFormat(competitionAnalysisArr[i-1].BamiloPrice)

		row.RowValue.BmlAvgPrice = formatting.ChangeNumberFormat(competitionAnalysisArr[i-1].BamiloAvgPrice)

		row.RowValue.BmlAvgSpecialPrice = formatting.ChangeNumberFormat(competitionAnalysisArr[i-1].BamiloAvgSpecialPrice)

		row.RowValue.DgkIDCatalogConfig = competitionAnalysisArr[i-1].IDDgkCatalogConfig

		row.RowValue.DgkScore = int(competitionAnalysisArr[i-1].DgkScore)

		row.RowValue.DgkSKUName = competitionAnalysisArr[i-1].DgkSKUName

		row.RowValue.DgkImgLink = competitionAnalysisArr[i-1].DgkImgLink

		row.RowValue.DgkSKULink = competitionAnalysisArr[i-1].DgkSKULink

		row.RowValue.GoodMatch = competitionAnalysisArr[i-1].GoodMatch

		if competitionAnalysisArr[i-1].DgkPrice != "NoMatch" {
			competitionAnalysisArr[i-1].DgkPrice = formatting.ChangeNumberFormat(competitionAnalysisArr[i-1].DgkPrice)
		}
		if competitionAnalysisArr[i-1].DgkAvgSpecialPrice != "NoMatch" {
			competitionAnalysisArr[i-1].DgkAvgSpecialPrice = formatting.ChangeNumberFormat(competitionAnalysisArr[i-1].DgkAvgSpecialPrice)
		}
		competitionAnalysisArr[i-1].DgkAvgPrice = formatting.ChangeNumberFormat(competitionAnalysisArr[i-1].DgkAvgPrice)

		row.RowValue.DgkSKUPrice = competitionAnalysisArr[i-1].DgkPrice

		row.RowValue.DgkAvgPrice = competitionAnalysisArr[i-1].DgkAvgPrice

		row.RowValue.DgkAvgSpecialPrice = competitionAnalysisArr[i-1].DgkAvgSpecialPrice

		row.RowValue.BmlConfigSnapshotAt = competitionAnalysisArr[i-1].BmlConfigSnapshotAt
		row.RowValue.DgkConfigSnapshotAt = competitionAnalysisArr[i-1].DgkConfigSnapshotAt
		row.RowValue.BmlSupplierName = competitionAnalysisArr[i-1].SupplierName
		row.RowValue.BmlBrand = competitionAnalysisArr[i-1].BrandName
		row.RowValue.BmlMinOfStockQuantity = competitionAnalysisArr[i-1].MinOfStockQuantity
		row.RowValue.BmlSumOfStockQuantity = competitionAnalysisArr[i-1].SumOfStockQuantity
		if competitionAnalysisArr[i-1].IsOutOfStock == "true" {
			row.RowValue.DgkStock = "out-of-stock"
		} else if competitionAnalysisArr[i-1].IsOutOfStock == "false" {
			row.RowValue.DgkStock = "in-stock"
		}

		arrOfRow = append(arrOfRow, row)
	}
	var competitionAnalysisTable bmldgktable.Table
	competitionAnalysisTable.Table = arrOfRow
	return competitionAnalysisTable
}

func SetNullTabel() bmldgktable.Table {
	var arrOfRow []bmldgktable.Row
	n := 1 // one empty row
	for i := 1; i < n+1; i++ {
		var row bmldgktable.Row
		row.RowKey = i

		row.RowValue.BmlIDCatalogConfig = ""

		row.RowValue.BmlSKUName = ""

		row.RowValue.BmlImgLink = ""

		row.RowValue.BmlSKULink = ""

		row.RowValue.BmlSKUPrice = ""

		row.RowValue.DgkSKUName = ""

		row.RowValue.DgkImgLink = ""

		row.RowValue.DgkSKULink = ""
		row.RowValue.DgkSKUPrice = ""
		row.RowValue.BmlConfigSnapshotAt = ""
		row.RowValue.DgkConfigSnapshotAt = ""

		arrOfRow = append(arrOfRow, row)
	}
	var competitionAnalysisTable bmldgktable.Table
	competitionAnalysisTable.Table = arrOfRow
	return competitionAnalysisTable
}

func ShowPageNumber(pgNumber int, searchedBy string, biCat string, cat1 string, cat2 string, cat3 string, sortedBy string, mongoSession *mgo.Session) []int {
	//build array of pages to show
	var showPageArr []int
	//db := mongoSession.DB("competition_analysis")
	//bmlCatalogConfig := db.C("bml_catalog_config")

	if pgNumber > 1 {
		showPageArr = append(showPageArr, pgNumber-1)
	}
	showPageArr = append(showPageArr, pgNumber)
	for i := 1; i < 2; i++ {
		_, tableISNull := SelectCompetitionAnalysisTable(10, pgNumber+i, searchedBy, biCat, cat1, cat2, cat3, sortedBy, mongoSession)
		if !tableISNull {
			showPageArr = append(showPageArr, pgNumber+i)
			//count(bmlCatalogConfig.Find([]bson.M{stageFind,stageSort,stageSkip, stageLimit,stageMerge))

		} else {
			break
		}
	}

	return showPageArr
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
