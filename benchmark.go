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

func benchmarQuery(parent context.Context, db *sql.DB, pl plan, query query.Query) error {
	var (
		ap         = newArgsProducer()
		argsCh     = ap.run(parent, query.NumQuery, query.Args)
		group, ctx = errgroup.WithContext(parent)
		_, cancel  = context.WithCancel(ctx)
	)

	printBenchmarkQueryHeader(query)

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

	timeNeeded := time.Since(start)
	log.Printf("Query time       : %v", timeNeeded.String())
	log.Printf("Query per second : %v", float64(query.NumQuery)/timeNeeded.Seconds())
	return nil
}

func printBenchmarkQueryHeader(query query.Query) {
	log.Printf("Benchmarking query\n"+
		"\tname            : %v\n"+
		"\tquery string    : %v\n"+
		"\tnum_query       : %v\n"+
		"\tprepare         : %v\n"+
		"\tprepare_on_init : %v\n"+
		"\twith_placeholder: %v\n",
		query.Name, query.QueryStr, query.NumQuery,
		query.Prepare, query.PrepareOnInit, query.WithPlaceholder)

}
