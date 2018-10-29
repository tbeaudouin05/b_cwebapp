package configcountcategorypage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/SepidehKHH/formatting"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pmylund/sortutil"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
)

func SelectDgkCategoryAssortmentTable(cat string, mongoSession *mgo.Session) []bmldgktable.CategoryAssortmentTable {
	var selectedResult []bmldgktable.DgkCategorryAssortment

	fmt.Println("mongocheck-cat")

	db := mongoSession.DB("competition_analysis")
	dgkAggStatisticHist := db.C("dgk_agg_statistic_hist")

	typeName := "DgkCategoryTwoConfigCountTopPage"

	// draft for filtering //
	//stageFilterForDgkScore := bson.M{"$match": bson.M{"dgk_score":bson.M{"$gte":0}}}

	//check which queries are required
	now := time.Now()
	then := now.AddDate(0, 0, -15)
	if cat == "" || cat == "All" {
		dgkAggStatisticHist.Find(bson.M{"type": bson.M{"$eq": typeName}, "statistic_snapshot_at": bson.M{"$gte": then}}).Sort("dgk_category_two_name").All(&selectedResult)

	} else {
		dgkAggStatisticHist.Find(bson.M{"type": typeName,
			"dgk_category_one_name": cat, "statistic_snapshot_at": bson.M{"$gte": then},
		}).Sort("dgk_category_two_name").All(&selectedResult)
	}
	//fmt.Println("Categories Selected Result: ",selectedResult)
	var distinctSelectedResult []bmldgktable.DgkCategorryAssortment
	var cat2 string
	var statisticSnapshotAt time.Time
	if selectedResult != nil {
		statisticSnapshotAt = selectedResult[0].StatisticSnapshotAt
	}
	for _, rngSelectedResult := range selectedResult {
		if rngSelectedResult.CategoryTwoName != cat2 {
			distinctSelectedResult = append(distinctSelectedResult, rngSelectedResult)
			cat2 = rngSelectedResult.CategoryTwoName
		} else {
			lastIndex := len(distinctSelectedResult) - 1
			if distinctSelectedResult[lastIndex].StatisticSnapshotAt.Before(statisticSnapshotAt) {
				distinctSelectedResult[lastIndex].StatisticSnapshotAt = statisticSnapshotAt
			}
			statisticSnapshotAt = rngSelectedResult.StatisticSnapshotAt
		}
	}
	//fmt.Println("distinctSelectedResult: ",distinctSelectedResult)

	//var dgkCategoryAssortmentTable []bmldgktable.DgkCategorryAssortmentTable

	var categoryAssortmentTable []bmldgktable.CategoryAssortmentTable

	if distinctSelectedResult == nil {
		var tmpCategoryAssortmentTable bmldgktable.CategoryAssortmentTable
		tmpCategoryAssortmentTable.CategoryOneName = " "
		tmpCategoryAssortmentTable.RowKey = 0
		tmpCategoryAssortmentTable.ConfigCountCat1 = " "
		var tmpCatTwo bmldgktable.CategoryTwo
		tmpCatTwo.CategoryTwoName = " "
		tmpCatTwo.ConfigCountCat2 = " "
		tmpCatTwo.RowKey = 0
		tmpCategoryAssortmentTable.CategoryTwo = append(tmpCategoryAssortmentTable.CategoryTwo, tmpCatTwo)
		categoryAssortmentTable = append(categoryAssortmentTable, tmpCategoryAssortmentTable)
		return categoryAssortmentTable
	}

	///////////////////////sort by cat1 again :(((
	sortutil.AscByField(distinctSelectedResult, "CategoryOneName")
	var cat1 string // := selectedResult[0].CategoryOneName
	i := 0
	j := 0
	for i < len(distinctSelectedResult) {

		for cat1 != distinctSelectedResult[i].CategoryOneName {
			var tmpCategoryAssortmentTable bmldgktable.CategoryAssortmentTable
			tmpCategoryAssortmentTable.CategoryOneName = distinctSelectedResult[i].CategoryOneName
			tmpCategoryAssortmentTable.ConfigCountCat1Int = distinctSelectedResult[i].ConfigCount
			var tmpCatTwo bmldgktable.CategoryTwo
			tmpCatTwo.CategoryTwoName = distinctSelectedResult[i].CategoryTwoName
			tmpCatTwo.ConfigCountCat2 = formatting.ChangeNumberFormat(strconv.Itoa(distinctSelectedResult[i].ConfigCount))
			tmpCatTwo.RowKey = i
			tmpCategoryAssortmentTable.CategoryTwo = append(tmpCategoryAssortmentTable.CategoryTwo, tmpCatTwo)
			tmpCategoryAssortmentTable.RowKey = j
			categoryAssortmentTable = append(categoryAssortmentTable, tmpCategoryAssortmentTable)
			cat1 = distinctSelectedResult[i].CategoryOneName
			i++
			j++
			if i == len(distinctSelectedResult) {
				break
			}
		}

		if i == len(distinctSelectedResult) {
			break
		}
		for cat1 == distinctSelectedResult[i].CategoryOneName {
			categoryAssortmentTable[j-1].ConfigCountCat1Int = categoryAssortmentTable[j-1].ConfigCountCat1Int + distinctSelectedResult[i].ConfigCount
			var tmpCatTwo bmldgktable.CategoryTwo
			tmpCatTwo.CategoryTwoName = distinctSelectedResult[i].CategoryTwoName
			tmpCatTwo.ConfigCountCat2 = formatting.ChangeNumberFormat(strconv.Itoa(distinctSelectedResult[i].ConfigCount))
			tmpCatTwo.ConfigCountCat2Int = distinctSelectedResult[i].ConfigCount
			tmpCatTwo.RowKey = i
			categoryAssortmentTable[j-1].CategoryTwo = append(categoryAssortmentTable[j-1].CategoryTwo, tmpCatTwo)
			cat1 = distinctSelectedResult[i].CategoryOneName
			i++
			if i == len(distinctSelectedResult) {
				break
			}
		}
	}

	for cnt, rngCategoryAssortmentTable := range categoryAssortmentTable {
		categoryAssortmentTable[cnt].ConfigCountCat1 = formatting.ChangeNumberFormat(strconv.Itoa(categoryAssortmentTable[cnt].ConfigCountCat1Int))
		sortutil.DescByField(rngCategoryAssortmentTable.CategoryTwo, "ConfigCountCat2Int")
	}
	sortutil.DescByField(categoryAssortmentTable, "ConfigCountCat1Int")
	//fmt.Println(categoryAssortmentTable)
	return categoryAssortmentTable
}
