package query

// Query represents a database query
type Query struct {
	Type          string
	QueryStr      string
	Params        []Param
	Prepare       bool
	PrepareOnInit bool
}

// Param represents a database query param
type Param struct {
	DataType string
	GenType  string // random/sequential, currently unused, always sequential
	Prefix   string
}
