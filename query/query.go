package query

const (
	// TypeExec is type of query which will executes db.Exec command
	// https://golang.org/pkg/database/sql/#DB.Exec
	TypeExec = "exec"

	// TypeExecContext is type of query which will executes db.ExecContext command
	// https://golang.org/pkg/database/sql/#DB.ExecContext
	TypeExecContext = "exec_context"

	// TypeQuery is type of query which will executes db.Query command
	// https://golang.org/pkg/database/sql/#DB.Query
	TypeQuery = "query"

	// TypeQueryContext is type of query which will executes db.QueryContext command
	// https://golang.org/pkg/database/sql/#DB.QueryContext
	TypeQueryContext = "query_context"
)

const (
	// DataTypeInteger is query argument with data type = integer
	DataTypeInteger = "integer"

	// DataTypeString is query argument with data type = string
	DataTypeString = "string"
)

const (
	// GenTypeSequential define argument generator in sequential mode
	GenTypeSequential = "sequential"

	// GenTypeRandom define argument generator in random mode
	GenTypeRandom = "random"
)

// Query represents a database query
type Query struct {
	// Name of the query, only used for informational purpose
	Name string `yaml:"name"`

	// Total number of queries which will be executed
	NumQuery int `yaml:"num_query"`

	// Query type
	Type string `yaml:"type"`

	// Query string
	QueryStr string `yaml:"query_str"`

	// Query arguments
	Args []Arg `yaml:"args"`

	// True if we want to make it a prepared statement
	Prepare bool `yaml:"prepare"`

	// True if we want to prepare it once during the initialization.
	// Otherwise we prepare it before executing the query
	PrepareOnInit bool `yaml:"prepare_on_init"`

	// True if the query contains placeholder
	WithPlaceholder bool `yaml:"with_placeholder"`
}

// Arg represents a database query arg
type Arg struct {
	// arg data type, currently supported type:
	//	- integer
	// 	- string
	DataType string `yaml:"data_type"`

	// arg generation mode:
	// sequential:
	//		integer: query number
	//		string : prefix + query number
	// random:
	//		integer: between specified range
	GenType string `yaml:"gen_type"`

	// Prefix used to generate string arg
	// the argument will be prefixed with the given prefix
	// appended with query number
	Prefix string `yaml:"prefix"`

	// max value of generated random value
	RandomRangeMax int `yaml:"random_range_hi"`

	// min value of generated random value
	RandomRangeMin int `yaml:"random_range_low"`
}
