package main

import (
	"os"

	"github.com/iwanbk/gosqlbencher/query"
	"gopkg.in/yaml.v2"
)

type config struct {
	DataSourceName string
	NumWorker      int
	NumQuery       int
}

type plan struct {
	DataSourceName string
	NumWorker      int
	Queries        []query.Query
}

func readPlan(filename string) (*plan, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var p plan

	err = yaml.NewDecoder(f).Decode(&p)
	return &p, err
}
