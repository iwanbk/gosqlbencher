package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	numRows = int(5e4)
	cfg     = config{
		NumWorker:      30,
		DataSourceName: "postgres://127.0.0.1:5432/example?sslmode=disable",
	}
)

func main() {

	db := initDB(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := query{
		Type:     "exec",
		QueryStr: "insert into pgbench_accounts (aid,bid, abalance, filler)values($1, $2, $3,$4)",
		Params: []queryParam{
			queryParam{
				DataType: integerDataType,
			},
			queryParam{
				DataType: integerDataType,
			},
			queryParam{
				DataType: integerDataType,
			},
			queryParam{
				DataType: stringDataType,
				Prefix:   "name_",
			},
		},
	}

	wp := &workProducer{}
	argsCh := wp.run(ctx, 100, query.Params)

	// insert to table
	func() {
		log.Println("Insert data")
		start := time.Now()

		q := "insert into pgbench_accounts (aid,bid, abalance, filler)values($1, $2, $3,$4)"
		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			log.Fatalf("prepare failed: %v", err)
		}

		for args := range argsCh {
			_, err = stmt.ExecContext(ctx, args...)
			if err != nil {
				log.Fatalf("error insert: %v", err)
			}
		}
		log.Println("Insert data - finished")
		log.Printf("TPS = %v", float64(numRows)/time.Since(start).Seconds())
	}()

	// select with query

	// delete
}

func initDB(cfg config) *sql.DB {
	log.Println("Open DB")
	db, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	log.Println("Ping DB")
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	return db
}
