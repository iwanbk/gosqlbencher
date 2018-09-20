package plan

import (
	"os"

	"github.com/iwanbk/gosqlbencher/query"
	"gopkg.in/yaml.v2"
)

// Plan define the gosqlbencher benchmark plan
type Plan struct {
	// DriverName is the name of the sql driver.
	// Supported drivers:
	// - postgres : https://github.com/lib/pq
	// - pgx : https://github.com/jackc/pgx
	DriverName string `yaml:"driver_name"`

	// DataSourceName is data source or connection string of the database
	// being tested
	DataSourceName string `yaml:"data_source_name"`

	// NumWorker is number of goroutines which executes the benchmarked queries.
	// It simulates number of concurrent queries in our system
	NumWorker int `yaml:"num_worker"`

	// Queries is sequence of queries we want to test
	Queries []query.Query `yaml:"queries"`
}

// Read reads and parse Plan from the given filename
func Read(filename string) (p Plan, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&p)
	return
}
