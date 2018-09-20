package main

import (
	"context"
	"flag"
	"log"

	"github.com/iwanbk/gosqlbencher/plan"
)

var (
	planFile string
)

func main() {
	flag.StringVar(&planFile, "plan", "plan.yaml", "gosqlbencher plan file")
	flag.Parse()

	// read config
	pl, err := plan.Read(planFile)
	if err != nil {
		log.Fatalf("failed to read plan: %v", err)
	}

	// print some nice message
	log.Printf("Benchmarking\ndsn: %v\nNumWorker:%v\n",
		pl.DataSourceName, pl.NumWorker)

	// open DB
	db, err := openDB(pl.DriverName, pl.DataSourceName)
	if err != nil {
		log.Fatalf("failed to open database:%v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run the benchmarks
	for i, query := range pl.Queries {
		err = benchmarQuery(ctx, db, pl.NumWorker, query)
		if err != nil {
			log.Fatalf("benchmarck query #%v failed: %v", i, err)
		}
	}
}
