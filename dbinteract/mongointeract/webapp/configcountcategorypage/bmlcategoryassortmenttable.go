package configcountcategorypage
import (
	"time"
	"strconv"

	"github.com/SepidehKHH/formatting"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"
	"github.com/pmylund/sortutil"
)

func SelectBmlCategoryAssortmentTable(cat string, mongoSession *mgo.Session) [] bmldgktable.CategoryAssortmentTable {
	var selectedResult []bmldgktable.BmlCategorryAssortment

	db := mongoSession.DB("competition_analysis")
	bmlAggStatisticHist := db.C("bml_agg_statistic_hist")

 typeName := "BmlCategoryConfigCountTopPage"

	// draft for filtering //
	//stageFilterForDgkScore := bson.M{"$match": bson.M{"dgk_score":bson.M{"$gte":0}}}

	//check which queries are required
	now := time.Now()
	then := now.AddDate(0,0,-20)
	if cat == "" || cat =="All" {
		bmlAggStatisticHist.Find(bson.M{"type": bson.M{"$eq": typeName}, "statistic_snapshot_at": bson.M{"$gte":then},}).Sort("bi_category_two_name").All(&selectedResult)
		
	} else {
		bmlAggStatisticHist.Find(bson.M{"type": typeName,
			"bi_category_one_name": cat, "statistic_snapshot_at": bson.M{"$gte":then},
		}).Sort("bi_category_two_name").All(&selectedResult)
	}
	//fmt.Println("Categories Selected Result: ",selectedResult)
	var distinctSelectedResult []bmldgktable.BmlCategorryAssortment
	var cat2 string 
	var statisticSnapshotAt time.Time
	if(selectedResult!=nil){
		statisticSnapshotAt = selectedResult[0].StatisticSnapshotAt
	}
	for _,rngSelectedResult := range selectedResult {
		if rngSelectedResult.CategoryTwoName != cat2 {
			distinctSelectedResult = append(distinctSelectedResult,rngSelectedResult)
			cat2 = rngSelectedResult.CategoryTwoName 
		}else{
			
			lastIndex := len(distinctSelectedResult)-1
			if(lastIndex>=0){
				if distinctSelectedResult[lastIndex].StatisticSnapshotAt.Before(statisticSnapshotAt) {
					distinctSelectedResult[lastIndex].StatisticSnapshotAt = statisticSnapshotAt
			}
			}
			
		statisticSnapshotAt = rngSelectedResult.StatisticSnapshotAt
	}
	}
//fmt.Println("distinctSelectedResult: ",distinctSelectedResult)


	var CategoryAssortmentTable [] bmldgktable.CategoryAssortmentTable

	if distinctSelectedResult == nil{
		var tmpBmlCategoryAssortmentTable bmldgktable.CategoryAssortmentTable
		tmpBmlCategoryAssortmentTable.CategoryOneName = " " 
		tmpBmlCategoryAssortmentTable.RowKey = 0
		tmpBmlCategoryAssortmentTable.ConfigCountCat1 =" "
		var tmpCatTwo bmldgktable.CategoryTwo
			tmpCatTwo.CategoryTwoName = " "
			tmpCatTwo.ConfigCountCat2 = " "
			tmpCatTwo.RowKey = 0
			tmpBmlCategoryAssortmentTable.CategoryTwo = append(tmpBmlCategoryAssortmentTable.CategoryTwo,tmpCatTwo)
		CategoryAssortmentTable = append(CategoryAssortmentTable, tmpBmlCategoryAssortmentTable)
		return CategoryAssortmentTable
	}

///////////////////////sort by cat1 again :(((
	sortutil.AscByField(distinctSelectedResult, "CategoryOneName")
	var cat1 string // := selectedResult[0].CategoryOneName
	i := 0
	j := 0
	for i<len(distinctSelectedResult) {

		for cat1!=distinctSelectedResult[i].CategoryOneName {
			var tmpBmlCategoryAssortmentTable bmldgktable.CategoryAssortmentTable
			tmpBmlCategoryAssortmentTable.CategoryOneName = distinctSelectedResult[i].CategoryOneName
			tmpBmlCategoryAssortmentTable.ConfigCountCat1Int = distinctSelectedResult[i].CategoryOneConfigCount
			var tmpCatTwo bmldgktable.CategoryTwo
			tmpCatTwo.CategoryTwoName = distinctSelectedResult[i].CategoryTwoName
			tmpCatTwo.ConfigCountCat2 = formatting.ChangeNumberFormat(strconv.Itoa(distinctSelectedResult[i].CategoryTwoConfigCount))
			tmpCatTwo.RowKey = i
			tmpBmlCategoryAssortmentTable.CategoryTwo = append(tmpBmlCategoryAssortmentTable.CategoryTwo,tmpCatTwo)
			tmpBmlCategoryAssortmentTable.RowKey = j
			CategoryAssortmentTable = append(CategoryAssortmentTable,tmpBmlCategoryAssortmentTable)
			cat1 = distinctSelectedResult[i].CategoryOneName
			i++
			j++
			if i==len(distinctSelectedResult){
				break
			}
		}

		if i==len(distinctSelectedResult){
			break
		}
		for cat1==distinctSelectedResult[i].CategoryOneName{
			CategoryAssortmentTable[j-1].ConfigCountCat1Int = CategoryAssortmentTable[j-1].ConfigCountCat1Int + distinctSelectedResult[i].CategoryTwoConfigCount
			var tmpCatTwo bmldgktable.CategoryTwo
			tmpCatTwo.CategoryTwoName = distinctSelectedResult[i].CategoryTwoName
			tmpCatTwo.ConfigCountCat2 = formatting.ChangeNumberFormat(strconv.Itoa(distinctSelectedResult[i].CategoryTwoConfigCount))
			tmpCatTwo.ConfigCountCat2Int = distinctSelectedResult[i].CategoryTwoConfigCount
			tmpCatTwo.RowKey = i
			CategoryAssortmentTable[j-1].CategoryTwo = append(CategoryAssortmentTable[j-1].CategoryTwo,tmpCatTwo)
			cat1 = distinctSelectedResult[i].CategoryOneName
			i++
			if i==len(distinctSelectedResult){
				break
			}
		}
	}
	for cnt,rngCategoryAssortmentTable := range CategoryAssortmentTable{
		CategoryAssortmentTable[cnt].ConfigCountCat1 = formatting.ChangeNumberFormat(strconv.Itoa(CategoryAssortmentTable[cnt].ConfigCountCat1Int))
		sortutil.DescByField(rngCategoryAssortmentTable.CategoryTwo, "ConfigCountCat2Int")
	}
	sortutil.DescByField(CategoryAssortmentTable, "ConfigCountCat1Int")
	
	return CategoryAssortmentTable
}