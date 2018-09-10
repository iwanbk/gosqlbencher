package main

import (
	"os"

	"github.com/iwanbk/gosqlbencher/query"
	"gopkg.in/yaml.v2"
)

type plan struct {
	DataSourceName string        `yaml:"data_source_name"`
	NumWorker      int           `yaml:"num_worker"`
	Queries        []query.Query `yaml:"queries"`
}

func readPlan(filename string) (p plan, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&p)
	return
}
