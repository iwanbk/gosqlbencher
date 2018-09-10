package main

type query struct {
	Type     string
	QueryStr string
	Params   []queryParam
	Prepare  bool
}

type queryParam struct {
	DataType string
	GenType  string // random/sequential, currently unused, always sequential
	Prefix   string
}
