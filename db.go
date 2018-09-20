package main

import (
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

func openDB(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
