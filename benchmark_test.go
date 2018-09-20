package main

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"

	"github.com/iwanbk/gosqlbencher/query"
)

func TestBenchmarkQuery(t *testing.T) {
	// setup db
	db, err := sql.Open("ramsql", "TestBenchmarkQuery")
	require.NoError(t, err)
	defer db.Close()

	// setup table
	createTableStr := "CREATE TABLE address (id bigserial PRIMARY KEY, street TEXT, number INT)"
	_, err = db.Exec(createTableStr)
	require.NoError(t, err)

	// setup benchmark test
	const (
		prefix = "jalan_"
	)
	ctx := context.Background()
	que := query.Query{
		Name:            "test benchmark",
		NumQuery:        100,
		Type:            query.TypeExecContext,
		QueryStr:        `insert into address(street, number) values($1, $2)`,
		WithPlaceholder: true,
		Args: []query.Arg{
			{
				DataType: query.DataTypeString,
				GenType:  query.GenTypeSequential,
				Prefix:   prefix,
			},
			{
				DataType: query.DataTypeInteger,
				GenType:  query.GenTypeSequential,
			},
		},
	}

	// execute benchmark
	err = benchmarQuery(ctx, db, 10, que)
	require.NoError(t, err)

	// check the db, make sure the benchmark did things it supposed to do
	addressMap := make(map[int]string)
	rows, err := db.Query("select id, street from address")
	require.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var (
			id     int
			street string
		)
		err = rows.Scan(&id, &street)
		require.NoError(t, err)
		addressMap[id] = street
	}
	require.NoError(t, rows.Close())

	// check we have correct number of address
	require.Equal(t, que.NumQuery, len(addressMap))
	for i := 1; i <= que.NumQuery; i++ {
		street, ok := addressMap[i]
		require.True(t, ok)
		require.Contains(t, street, prefix)
	}
}
