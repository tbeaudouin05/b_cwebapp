package bmldgktable

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func GetGoodMatchCount(cat1 string, cat2 string, cat3 string, mongoSession *mgo.Session) int {

	db := mongoSession.DB("competition_analysis")
	bmlCatalogConfig := db.C("bml_catalog_config")
	var goodMatchCount goodMatchCount

	if cat1 == "" && cat2 == "" {
		goodMatchCountValue, err := bmlCatalogConfig.Find(bson.M{"$and": []bson.M{bson.M{"bi_category_three_name": cat3}, bson.M{"good_match": true}}}).Count()
		checkError(err)
		goodMatchCount.setGoodMatchCount(goodMatchCountValue)
	} else if cat1 == "" {
		goodMatchCountValue, err := bmlCatalogConfig.Find(bson.M{"$and": []bson.M{bson.M{"bi_category_two_name": cat2}, bson.M{"bi_category_three_name": cat3}, bson.M{"good_match": true}}}).Count()
		checkError(err)
		goodMatchCount.setGoodMatchCount(goodMatchCountValue)
	} else if cat2 == "" {
		goodMatchCountValue, err := bmlCatalogConfig.Find(bson.M{"$and": []bson.M{bson.M{"bi_category_one_name": cat1}, bson.M{"bi_category_three_name": cat3}, bson.M{"good_match": true}}}).Count()
		checkError(err)
		goodMatchCount.setGoodMatchCount(goodMatchCountValue)
	} else {
		goodMatchCountValue, err := bmlCatalogConfig.Find(bson.M{"$and": []bson.M{bson.M{"bi_category_one_name": cat1}, bson.M{"bi_category_two_name": cat2}, bson.M{"bi_category_three_name": cat3}, bson.M{"good_match": true}}}).Count()
		checkError(err)
		goodMatchCount.setGoodMatchCount(goodMatchCountValue)
	}

	return goodMatchCount.goodMatchCount

}

type goodMatchCount struct{ goodMatchCount int }

func (goodMatchCount *goodMatchCount) setGoodMatchCount(goodMatchCountValue int) {
	goodMatchCount.goodMatchCount = goodMatchCountValue
}
