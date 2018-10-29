package backendprocess

import (
	"database/sql"
	"log"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/bmlcatalogconfig"
)

// GetBmlCatalogConfigSalesTable retrieves the sum of sales data across a certain period for each SKU from BOB
func GetBmlCatalogConfigSalesTable(dbBob *sql.DB) []bmlcatalogconfig.BmlCatalogConfig {

	stmt, err := dbBob.Prepare(`
		SELECT

		cc.id_catalog_config

		,COUNT(DISTINCT soi.id_sales_order_item) count_of_soi
		,COALESCE(FLOOR(SUM(soi.paid_price))) sum_of_paid_price
		,COALESCE(FLOOR(SUM(soi.unit_price))) sum_of_unit_price
		,COALESCE(FLOOR(SUM(soi.coupon_money_value))) sum_of_coupon_money_value
		,COALESCE(FLOOR(SUM(soi.cart_rule_discount))) sum_of_cart_rule_discount
	
		FROM sales_order_item soi
		JOIN catalog_simple cs
		ON soi.sku = cs.sku
		JOIN catalog_config cc
		ON cs.fk_catalog_config = cc.id_catalog_config
	
		WHERE CAST(soi.created_at AS DATE) = CAST(NOW()-INTERVAL 30 DAY AS DATE)
	
		GROUP BY cc.id_catalog_config;
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	var bamiloCatalogConfigSalesTable []bmlcatalogconfig.BmlCatalogConfig
	bamiloCatalogConfigSales := bmlcatalogconfig.BmlCatalogConfig{}

	for rows.Next() {

		err := rows.Scan(
			&bamiloCatalogConfigSales.IDBmlCatalogConfig,

			&bamiloCatalogConfigSales.CountOfSoi,
			&bamiloCatalogConfigSales.SumOfPaidPrice,
			&bamiloCatalogConfigSales.SumOfUnitPrice,
			&bamiloCatalogConfigSales.SumOfCouponMoneyValue,
			&bamiloCatalogConfigSales.SumOfCartRuleDiscount,
		)
		checkError(err)

		bamiloCatalogConfigSalesTable = append(bamiloCatalogConfigSalesTable, bamiloCatalogConfigSales)

	}

	//log.Println(` Length of the table: ` + strconv.Itoa(len(bamiloCatalogConfigSalesTable)))

	return bamiloCatalogConfigSalesTable

}

// GetBmlCatalogConfigSalesHistTable retrieves the sum of sales for each day and each SKU from BOB
func GetBmlCatalogConfigSalesHistTable(dbBob *sql.DB) []bmlcatalogconfig.BmlCatalogConfig {

	stmt, err := dbBob.Prepare(`
		SELECT

		cc.id_catalog_config
		,DATE_FORMAT(soi.created_at, "%m/%d/%Y") config_snapshot_at

		,COUNT(DISTINCT soi.id_sales_order_item) count_of_soi
		,COALESCE(FLOOR(SUM(soi.paid_price))) sum_of_paid_price
		,COALESCE(FLOOR(SUM(soi.unit_price))) sum_of_unit_price
		,COALESCE(FLOOR(SUM(soi.coupon_money_value))) sum_of_coupon_money_value
		,COALESCE(FLOOR(SUM(soi.cart_rule_discount))) sum_of_cart_rule_discount
	
		FROM sales_order_item soi
		JOIN catalog_simple cs
		ON soi.sku = cs.sku
		JOIN catalog_config cc
		ON cs.fk_catalog_config = cc.id_catalog_config
	
		WHERE CAST(soi.created_at AS DATE) >= CAST(NOW()-INTERVAL 5 DAY AS DATE)
	
		GROUP BY cc.id_catalog_config, CAST(soi.created_at AS DATE);
	 `)

	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	var bamiloCatalogConfigSalesHistTable []bmlcatalogconfig.BmlCatalogConfig
	var configSnapshotAtStr string
	bamiloCatalogConfigSalesHist := bmlcatalogconfig.BmlCatalogConfig{}

	for rows.Next() {

		err := rows.Scan(
			&bamiloCatalogConfigSalesHist.IDBmlCatalogConfig,
			&configSnapshotAtStr,

			&bamiloCatalogConfigSalesHist.CountOfSoi,
			&bamiloCatalogConfigSalesHist.SumOfPaidPrice,
			&bamiloCatalogConfigSalesHist.SumOfUnitPrice,
			&bamiloCatalogConfigSalesHist.SumOfCouponMoneyValue,
			&bamiloCatalogConfigSalesHist.SumOfCartRuleDiscount,
		)
		checkError(err)

		bamiloCatalogConfigSalesHist.ConfigSnapshotAt, err = time.Parse(`01/02/2006`, configSnapshotAtStr)
		checkError(err)

		bamiloCatalogConfigSalesHistTable = append(bamiloCatalogConfigSalesHistTable, bamiloCatalogConfigSalesHist)

	}

	return bamiloCatalogConfigSalesHistTable

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
