package biinteract

import (
	"database/sql"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig/backendprocess"
	"github.com/thomas-bamilo/nosql/mongobulk"
)

func GetBiData(dbBi *sql.DB, mongoSession *mgo.Session) {

	// store sellerCenterQuery in a string
	stmt, err := dbBi.Prepare(
		`
	SELECT 

	cat.id_catalog_config
	,COALESCE(cat.category_bi,'')
	,COALESCE(cat.category_1_en,'')
	,COALESCE(cat.category_2_en,'')
	,COALESCE(cat.category_3_en,'')
	,COALESCE(cat.category_4_en,'')
	,COALESCE(cat.category_5_en,'')
	,COALESCE(cat.category_6_en,'')
	,COALESCE(cat.department,'')
	,COALESCE(sup.key_account_manager,'')
	FROM StagingDB_Replica.Gathering.tblConfigCatalog cat
	LEFT JOIN StagingDB_Replica.Gathering.tblSupplierCatalog sup ON cat.fk_supplier = sup.id_supplier 	`)

	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	bamiloCatalogConfig := backendprocess.BiRow{}

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config")
	config := mongobulk.Config{OpsPerBatch: 950}
	mongoBulk := mongobulk.New(c, config)

	for rows.Next() {

		err := rows.Scan(

			&bamiloCatalogConfig.CatId,
			&bamiloCatalogConfig.CatBi,
			&bamiloCatalogConfig.CatEn1,
			&bamiloCatalogConfig.CatEn2,
			&bamiloCatalogConfig.CatEn3,
			&bamiloCatalogConfig.CatEn4,
			&bamiloCatalogConfig.CatEn5,
			&bamiloCatalogConfig.CatEn6,
			&bamiloCatalogConfig.Department,
			&bamiloCatalogConfig.KeyAccountManager,
		)
		checkError(err)

		updateConfigMongo(&bamiloCatalogConfig, mongoBulk)
	}
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
func updateConfigMongo(bamiloCatalogConfigMongo *backendprocess.BiRow, mongoBulk *mongobulk.Bulk) {

	bamiloCatalogConfigMongoByte, err := bson.Marshal(bamiloCatalogConfigMongo)
	if err != nil {
		log.Println("Error marshaling: ", err.Error())

	}
	bamiloCatalogConfigMongoBson := make(map[string]interface{})
	err = bson.Unmarshal(bamiloCatalogConfigMongoByte, bamiloCatalogConfigMongoBson)
	if err != nil {
		log.Println("Error unmarshaling: ", err.Error())

	}
	colQuerier := bson.M{"_id": bamiloCatalogConfigMongo.CatId}

	mongoBulk.Update(colQuerier, bson.M{"$set": bamiloCatalogConfigMongoBson})
}
