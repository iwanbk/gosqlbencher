package query

// Query represents a database query
type Query struct {
	NumQuery        int     `yaml:"num_query"`
	Type            string  `yaml:"type"`
	QueryStr        string  `yaml:"query_str"`
	Params          []Param `yaml:"params"`
	Prepare         bool    `yaml:"prepare"`
	PrepareOnInit   bool    `yaml:"prepare_on_init"`
	WithPlaceholder bool    `yaml:"with_placeholder"`
}

// Param represents a database query param
type Param struct {
	DataType string `yaml:"data_type"`
	// random/sequential, currently unused, always sequential
	GenType string `yaml:"gen_type"`
	Prefix  string `yaml:"prefix"`
}
