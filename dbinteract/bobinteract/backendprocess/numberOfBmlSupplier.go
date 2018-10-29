package backendprocess

import (
	"database/sql"

)

func NumberOfBmlSupplier(dbBob *sql.DB) int {

	stmt, err := dbBob.Prepare(`
	-- it is fine if several SKUs have the same name, keep structure similar to BOB - plus it will be much clearer for business users
	SELECT
	  count(DISTINCT s.name) count_of_supplier
	  FROM supplier s
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()
	var count int

	for rows.Next() {
		err := rows.Scan(
			&count,
		)
		checkError(err)

	}

	return count

}
