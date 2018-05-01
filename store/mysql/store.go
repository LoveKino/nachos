package mysql

import (
	"database/sql"
	. "github.com/LoveKino/nachos/store"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

/**
 * query sql store
 */
func query(db *sql.DB, sqlText string) QueryStoreType {
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		log.Fatal("Fail to prepare sql " + sqlText)
		panic(err)
	}

	return func(params ...interface{}) (Records, error) {
		rows, rerr := stmt.Query(params...)
		if rerr != nil {
			return nil, nil
		}
		return rowsToStringMap(rows)
	}
}

func exec(db *sql.DB, sqlText string) ExecStoreType {
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		log.Fatal("Fail to prepare sql: " + sqlText)
		panic(err)
	}

	return func(params ...interface{}) (interface{}, error) {
		return stmt.Exec(params...)
	}
}

func GetStoreConstructor(dsn string) (StoreConstructor, error) {
	var storeConstructor StoreConstructor

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return storeConstructor, nil
	}

	getQueryStore := func(sqlText string) QueryStoreType {
		return query(db, sqlText)
	}

	getExecStore := func(sqlText string) ExecStoreType {
		return exec(db, sqlText)
	}

	storeConstructor = StoreConstructor{
		GetQueryStore: getQueryStore,
		GetExecStore:  getExecStore,
	}

	return storeConstructor, nil
}

/**
 * convert queried rows to a string map
 */
func rowsToStringMap(rows *sql.Rows) ([]interface{}, error) {
	var result []interface{}
	cols, _ := rows.Columns()

	defer rows.Close()
	for rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			m[colName] = columns[i]
		}

		result = append(result, m)
	}

	return result, nil
}
