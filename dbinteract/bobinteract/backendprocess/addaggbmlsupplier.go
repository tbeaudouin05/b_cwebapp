package backendprocess

import (
	"database/sql"

	"github.com/thomas-bamilo/commercial/competitionanalysis/struct/webapp/supplier"
)

func AddAggBmlSupplier(dbBob *sql.DB) []supplier.Supplier {

	stmt, err := dbBob.Prepare(`
	-- it is fine if several SKUs have the same name, keep structure similar to BOB - plus it will be much clearer for business users
	SELECT
	  COALESCE(s.name,'') supplier_name
	  ,COALESCE(s.name_en,'') supplier_name_en
	  ,COALESCE(COUNT(cc.id_catalog_config),0) config_count
	  ,COALESCE(SUM(case when cs2.quantity > 0 then 1 else 0 end),0) in_stock_config_count
	  ,COALESCE(s.supplier_code,0) supplier_code 
	  ,COALESCE(s.id_supplier,0) id_supplier

	  FROM catalog_simple cs
	  LEFT JOIN catalog_source cs1
	  ON cs.id_catalog_simple = cs1.fk_catalog_simple
	  JOIN catalog_config cc
	  ON cs.fk_catalog_config = cc.id_catalog_config
	  LEFT JOIN catalog_stock cs2
	  ON cs1.id_catalog_source = cs2.fk_catalog_source
 	  LEFT JOIN supplier s
	  ON cs1.fk_supplier = s.id_supplier

	  GROUP BY s.name
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	var bamiloSupplierConfigCount []supplier.Supplier
	tmpBamiloSupplierConfigCount := supplier.Supplier{}

	for rows.Next() {

		err := rows.Scan(
			&tmpBamiloSupplierConfigCount.SupplierName,
			&tmpBamiloSupplierConfigCount.SupplierNameEn,
			&tmpBamiloSupplierConfigCount.ConfigCount,
			&tmpBamiloSupplierConfigCount.InStockConfigCount,
			&tmpBamiloSupplierConfigCount.SupplierCode,
			&tmpBamiloSupplierConfigCount.SupplierID,
		)
		checkError(err)

		bamiloSupplierConfigCount = append(bamiloSupplierConfigCount, tmpBamiloSupplierConfigCount)

	}

	return bamiloSupplierConfigCount

}
