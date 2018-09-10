package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/iwanbk/gosqlbencher/executor"
	"github.com/iwanbk/gosqlbencher/query"
)

var (
	numRows = int(5e4)
	cfg     = config{
		NumWorker:      30,
		DataSourceName: "postgres://127.0.0.1:5432/example?sslmode=disable",
		NumQuery:       10000,
	}
)

func main() {

	db := initDB(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := query.Query{
		Type:     "exec",
		QueryStr: "insert into pgbench_accounts (aid,bid, abalance, filler)values($1, $2, $3,$4)",
		Params: []query.Param{
			query.Param{
				DataType: integerDataType,
			},
			query.Param{
				DataType: integerDataType,
			},
			query.Param{
				DataType: integerDataType,
			},
			query.Param{
				DataType: stringDataType,
				Prefix:   "name_",
			},
		},
		Prepare:       true,
		PrepareOnInit: true,
	}

	wp := &workProducer{}

	runner, err := executor.New(db, query)
	if err != nil {
		log.Fatalf("failed to create executor: %v", err)
	}

	argsCh := wp.run(ctx, cfg.NumQuery, query.Params)

	// insert to table
	func() {
		log.Println("Insert data")
		start := time.Now()

		for args := range argsCh {
			err = runner.Execute(ctx, args...)
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
