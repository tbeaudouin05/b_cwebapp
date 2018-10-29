package bmldgkmanualmatching

import (
	"strconv"
	"github.com/SepidehKHH/formatting"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/dgkcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/bmldgktable"

)

func RunManualMatchingPage(bmlID int,pgNumber int,searchedBy string,cat1 string,cat2 string, minPrice int, maxPrice int, mongoSession *mgo.Session) (bmldgktable.DgkTable,bool) {

	var dgkTable bmldgktable.DgkTable

	var dgkCatalogConfigStr [] dgkcatalogconfig.DgkCatalogConfig

	db := mongoSession.DB("competition_analysis")
	dgkCatalogConfig := db.C("dgk_catalog_config")
	
	 // we can set defult searchedBy = BmlSKU but it's gonna be slower //


	 all := "All"
	 if(searchedBy==""){
		 if cat1==all{ 
			if cat2==all{ //nothing
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice}}).Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}else {					  //cat2
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"dgk_category_two_name" : cat2,
				}).Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}
		 }else{
			if  cat2==all{ //cat1
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"dgk_category_one_name" : cat1,
				}).Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}else {					  //cat1 , cat2
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"dgk_category_one_name" : cat1,
				"dgk_category_two_name" : cat2,
				}).Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}
		 }
		
	 }else {
		 if cat1 ==all{
			if cat2==all{ //search
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"$text": bson.M{"$search": searchedBy},
				}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}else{				    //cat2 ,search
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"dgk_category_two_name" : cat2,
				"$text": bson.M{"$search": searchedBy},
					}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}
		 }else{
			if cat2==all{ //cat1 , search
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"dgk_category_one_name" : cat1,
				"$text": bson.M{"$search": searchedBy},
				}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}else{					//cat1 ,cat2 ,search
				dgkCatalogConfig.Find(bson.M{"avg_price":bson.M{"$gte":minPrice , "$lte":maxPrice},"dgk_category_one_name" : cat1,
				"dgk_category_two_name" : cat2,
				"$text": bson.M{"$search": searchedBy},
					}).Select(bson.M{"score": bson.M{"$meta": "textScore"}}).Sort("$textScore:score").Skip(10*pgNumber-10).Limit(10).All(&dgkCatalogConfigStr)
			}
		 }
	 }
	

	dgkTable = fillDgkTableByDgkCatalogConfigData(dgkCatalogConfigStr,minPrice ,maxPrice)
	if(dgkCatalogConfigStr==nil){
		return dgkTable,true
	}
	return dgkTable,false

}

func fillDgkTableByDgkCatalogConfigData(dgkCatalogConfigStr [] dgkcatalogconfig.DgkCatalogConfig, minPrice int ,maxPrice int) bmldgktable.DgkTable{
	var dgkTable bmldgktable.DgkTable
	for i,tmpDgkCatalogConfigStr := range dgkCatalogConfigStr{
		var tmpRowValue bmldgktable.DgkRow
		tmpRowValue.RowValue.DgkSKUName = tmpDgkCatalogConfigStr.SkuName
		tmpRowValue.RowValue.DgkImgLink = tmpDgkCatalogConfigStr.ImgLink
		tmpRowValue.RowValue.DgkSKULink = tmpDgkCatalogConfigStr.SkuLink
		tmpPrice := tmpDgkCatalogConfigStr.AvgSpecialPrice
		if(tmpPrice == 0 ){ 
			tmpPrice = tmpDgkCatalogConfigStr.AvgPrice
		}
		tmpRowValue.RowValue.DgkSKUPrice = formatting.ChangeNumberFormat(strconv.Itoa(tmpPrice))
		tmpRowValue.RowValue.DgkIDCatalogConfig = tmpDgkCatalogConfigStr.IDDgkCatalogConfig
		tmpRowValue.RowKey =i 
	//	if  minPrice <= tmpPrice && tmpPrice <= maxPrice{
			dgkTable.Table = append(dgkTable.Table,tmpRowValue)
	//	}
	}
	
	return dgkTable
}

func ShowPageNumber(bmlID int,pgNumber int, searchedBy string, cat1 string, cat2 string, minPrice int, maxPrice int, mongoSession *mgo.Session)[]int{
	//build array of pages to show 
	var showPageArr [] int
	//db := mongoSession.DB("competition_analysis")
	//bmlCatalogConfig := db.C("bml_catalog_config")
	
	if pgNumber>1{
		showPageArr = append(showPageArr,pgNumber-1)
	}
		showPageArr = append(showPageArr,pgNumber)
		for i:=1; i<2; i++{
			_,tableISNull := RunManualMatchingPage(bmlID, pgNumber+i, searchedBy, cat1, cat2, minPrice, maxPrice, mongoSession)
			if(!tableISNull){
				showPageArr = append(showPageArr,pgNumber+i)
				//count(bmlCatalogConfig.Find([]bson.M{stageFind,stageSort,stageSkip, stageLimit,stageMerge))
	
			}else{ break}
		}
	
		return showPageArr
	}
func SetNullTabel() ([]bmldgktable.DgkRow){
		var arrOfRow []bmldgktable.DgkRow
		n := 1 // one empty row 
		for i:=1; i<n+1; i++{
			var row bmldgktable.DgkRow
			row.RowKey = i
	
			row.RowValue.DgkSKUName = "null"
			row.RowValue.DgkImgLink = ""
			row.RowValue.DgkSKULink = ""
			row.RowValue.DgkSKUPrice = "null"
			
			arrOfRow = append(arrOfRow, row)
		}

		return arrOfRow
}