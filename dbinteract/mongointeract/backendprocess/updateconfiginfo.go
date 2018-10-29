package backendprocess

import (
	"log"

	"github.com/globalsign/mgo/bson"
	bmlconfigbackendprocess "github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig/backendprocess"
	mongobulk "github.com/thomas-bamilo/nosql/mongobulk"
)

func UpsertConfigMongo(bamiloCatalogConfigMongo *bmlconfigbackendprocess.Mongo, mongoBulk *mongobulk.Bulk) {

	bamiloCatalogConfigMongoByte, err := bson.Marshal(bamiloCatalogConfigMongo)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())

	}
	bamiloCatalogConfigMongoBson := make(map[string]interface{})
	err = bson.Unmarshal(bamiloCatalogConfigMongoByte, bamiloCatalogConfigMongoBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())

	}
	colQuerier := bson.M{"_id": bamiloCatalogConfigMongo.IDBmlCatalogConfig}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": bamiloCatalogConfigMongoBson})
}

func UpsertConfigMongoHist(bamiloCatalogConfigMongoHist *bmlconfigbackendprocess.MongoHist, mongoBulk *mongobulk.Bulk) {

	bamiloCatalogConfigMongoHistByte, err := bson.Marshal(bamiloCatalogConfigMongoHist)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())

	}
	bamiloCatalogConfigMongoHistBson := make(map[string]interface{})
	err = bson.Unmarshal(bamiloCatalogConfigMongoHistByte, bamiloCatalogConfigMongoHistBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())

	}
	colQuerier := bson.M{"_id": bamiloCatalogConfigMongoHist.IDBmlCatalogConfigHist}

	mongoBulk.Upsert(colQuerier, bson.M{"$set": bamiloCatalogConfigMongoHistBson})

}
