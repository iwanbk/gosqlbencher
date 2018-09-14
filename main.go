package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"

	"github.com/iwanbk/gosqlbencher/plan"
)

var (
	planFile string
)

func main() {
	flag.StringVar(&planFile, "plan", "plan.yaml", "gosqlbencher plan file")
	flag.Parse()

	pl, err := plan.Read(planFile)
	if err != nil {
		log.Fatalf("failed to read plan: %v", err)
	}

	log.Printf("Benchmarking\ndsn: %v\nNumWorker:%v\n",
		pl.DataSourceName, pl.NumWorker)

	db := initDB(pl)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, query := range pl.Queries {
		err = benchmarQuery(ctx, db, pl, query)
		if err != nil {
			log.Fatalf("benchmarck query #%v failed: %v", i, err)
		}
	}
}

func initDB(pl plan.Plan) *sql.DB {
	log.Println("Open DB")
	db, err := sql.Open("postgres", pl.DataSourceName)
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
