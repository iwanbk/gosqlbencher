package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"github.com/iwanbk/gosqlbencher/executor"
	"github.com/iwanbk/gosqlbencher/query"
)

func main() {
	pl, err := readPlan("plan.yaml")
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

func benchmarQuery(parent context.Context, db *sql.DB, pl plan, query query.Query) error {
	var (
		wp         = &workProducer{}
		argsCh     = wp.run(parent, query.NumQuery, query.Params)
		group, ctx = errgroup.WithContext(parent)
		_, cancel  = context.WithCancel(ctx)
	)
	log.Printf("--------\n name: %v\n query string: %v\n num_query: %v\n "+
		"prepare: %v\n prepare_on_init: %v\n with_placeholder: %v\n",
		query.Name, query.QueryStr, query.NumQuery,
		query.Prepare, query.PrepareOnInit, query.WithPlaceholder)

	defer cancel()
	start := time.Now()

	for i := 0; i < pl.NumWorker; i++ {
		group.Go(func() error {
			runner, err := executor.New(db, query)
			if err != nil {
				return err
			}
			defer runner.Close()

			for {
				select {
				case <-ctx.Done():
					return nil
				case args, ok := <-argsCh:
					if !ok { // channel is closed, no more work
						return nil
					}

					err = runner.Execute(ctx, args...)
					if err != nil {
						return err
					}
				}
			}
		})
	}
	err := group.Wait()
	if err != nil {
		return err
	}

	log.Printf("TPS = %v", float64(query.NumQuery)/time.Since(start).Seconds())
	return nil
}

func initDB(pl plan) *sql.DB {
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
