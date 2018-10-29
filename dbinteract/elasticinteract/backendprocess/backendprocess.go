package backendprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig/backendprocess"
	"gopkg.in/olivere/elastic.v5"
)

// match dgk ----------------------------------------------------

func MatchBmlDgkConfig(elasticClient *elastic.Client, ctx context.Context, bamiloCatalogConfigTable []bmlcatalogconfig.BmlCatalogConfig) (bamiloCatalogConfigMatchTable []backendprocess.Match) {

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {

		bamiloCatalogConfigMatch := matchBmlDgkConfig(&bamiloCatalogConfig, elasticClient, ctx)

		bamiloCatalogConfigMatchTable = append(bamiloCatalogConfigMatchTable, bamiloCatalogConfigMatch)

	}

	return bamiloCatalogConfigMatchTable
}

func matchBmlDgkConfig(bamiloCatalogConfig *bmlcatalogconfig.BmlCatalogConfig, elasticClient *elastic.Client, ctx context.Context) (bamiloCatalogConfigMatch backendprocess.Match) {

	defer elasticRecovery()

	multiQuerySkuName := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.SKUName,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2).Boost(30)

	multiQueryDescription := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.Description,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	multiQueryShortDescription := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.ShortDescription,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	multiQueryPackageContent := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.PackageContent,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	multiQueryProductWarranty := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.ProductWarranty,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	multiQueryBiCategoryOneName := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.BiCategoryOneName,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	multiQueryBiCategoryTwoName := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.BiCategoryTwoName,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	multiQueryBiCategoryThreeName := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.BiCategoryThreeName,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	// CHANGE THIS BRAND!
	multiQueryBrandName := elastic.NewMultiMatchQuery(
		bamiloCatalogConfig.BiCategoryThreeName,
		"sku_name", "sku_name.persian", "brand_name", "dgk_category_one_name", "dgk_category_two_name",
	).Fuzziness("AUTO").PrefixLength(2)

	percentAvgPrice := bamiloCatalogConfig.AvgPrice / 10
	gaussDecayFunction := elastic.NewGaussDecayFunction().
		Origin(bamiloCatalogConfig.AvgPrice).
		Scale(percentAvgPrice).
		FieldName("avg_price")
	functionScoreQueryAvgPrice := elastic.NewFunctionScoreQuery().
		AddScoreFunc(gaussDecayFunction).
		Boost(3).
		BoostMode("multiply")

	query := elastic.NewBoolQuery().Should(
		multiQuerySkuName,
		multiQueryDescription,
		multiQueryShortDescription,
		multiQueryPackageContent,
		multiQueryProductWarranty,
		multiQueryBiCategoryOneName,
		multiQueryBiCategoryTwoName,
		multiQueryBiCategoryThreeName,
		multiQueryBrandName,
		functionScoreQueryAvgPrice)

	searchResult, err := elasticClient.Search(`dgk_catalog_config`).
		Type(`dgk_catalog_config`).
		Query(query).Size(1).
		Do(ctx)
	printError(err)

	// attempt to unmarshall into string
	err = json.Unmarshal(*searchResult.Hits.Hits[0].Source, &bamiloCatalogConfigMatch)
	printError(err)
	// get the score of the result
	var dgkScoreNotPointer = *searchResult.Hits.Hits[0].Score
	bamiloCatalogConfigMatch.DgkScore = dgkScoreNotPointer
	// identify bamilo sku in bamiloCatalogConfigMatch
	bamiloCatalogConfigMatch.IDBmlCatalogConfig = bamiloCatalogConfig.IDBmlCatalogConfig

	bamiloCatalogConfigMatch.Matched = true

	return bamiloCatalogConfigMatch

}

func elasticRecovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered: ", r)
		debug.PrintStack()
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func printError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

/*func UpsertConfigInfo(elasticClient *elastic.Client, ctx context.Context, bamiloCatalogConfigTable []bmlcatalogconfig.BmlCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()
	// make it faster: erase all data then re-insert everything, do not provide custom ID, just add IDBmlCatalogConfig as a field
	// and / or use bulk insert (your nemesis)

	_, err := elasticClient.DeleteIndex(`bml_catalog_config`).Do(ctx)
	checkError(err)

	_, err = elasticClient.CreateIndex(`bml_catalog_config`).Do(ctx)
	checkError(err)

	// Setup a bulk processor
	bulkProcessor, err := elasticClient.BulkProcessor().
		Name("bulkProcessor").
		Workers(4).       // number of workers
		BulkActions(950). // commit if # requests >= 950
		Do(ctx)
	if err != nil {
		log.Println(err)
	}

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {
		// only keep appropriate information for elastic
		bamiloCatalogConfigElastic := bmlcatalogconfig.BmlCatalogConfigElastic{
			IDBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
			SKUName:            bamiloCatalogConfig.SKUName,
			Description:        bamiloCatalogConfig.Description,
			ShortDescription:   bamiloCatalogConfig.ShortDescription,
			PackageContent:     bamiloCatalogConfig.PackageContent,
			ProductWarranty:    bamiloCatalogConfig.ProductWarranty,

			BiCategoryOneName:   bamiloCatalogConfig.BiCategoryOneName,
			BiCategoryTwoName:   bamiloCatalogConfig.BiCategoryTwoName,
			BiCategoryThreeName: bamiloCatalogConfig.BiCategoryThreeName,

			BrandName:   bamiloCatalogConfig.BrandName,
			BrandNameEn: bamiloCatalogConfig.BrandNameEn,

			SupplierName:   bamiloCatalogConfig.SupplierName,
			SupplierNameEn: bamiloCatalogConfig.SupplierNameEn,

			AvgPrice:        bamiloCatalogConfig.AvgPrice,
			AvgSpecialPrice: bamiloCatalogConfig.AvgSpecialPrice,
		}

		bulkIndexRequest := elastic.NewBulkIndexRequest().
			Index("bml_catalog_config").
			Type("bml_catalog_config").
			Id(strconv.Itoa(bamiloCatalogConfigElastic.IDBmlCatalogConfig)).
			Doc(bamiloCatalogConfigElastic)

		bulkProcessor.Add(bulkIndexRequest)

	}

	end := time.Now()
	log.Println(`End time Elastic: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed Elastic: `, duration.Minutes(), ` minutes`)

	err = bulkProcessor.Close()
	checkError(err)

}*/
