package executor

import (
	"context"
	"database/sql"
	"testing"

	"github.com/iwanbk/gosqlbencher/query"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

func TestPreparer(t *testing.T) {
	db, err := sql.Open("ramsql", "TestPrepare")
	require.NoError(t, err)
	defer db.Close()

	createTableStr := "CREATE TABLE address (id INT, street TEXT, number INT)"
	_, err = db.Exec(createTableStr)
	require.NoError(t, err)

	_, err = db.Exec(`insert into address( id, street, number)values(1,'jalan', 1)`)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		queryType     query.TypeType
		queryStr      string
		args          []interface{}
		prepareOnInit bool
	}{
		{
			name:      "query",
			queryType: query.TypeQuery,
			queryStr:  "select * from address where id=$1",
			args:      []interface{}{1},
		},
		{
			name:      "query context",
			queryType: query.TypeQueryContext,
			queryStr:  "select * from address where id=$1",
			args:      []interface{}{1},
		},
		{
			name:      "exec",
			queryType: query.TypeExec,
			queryStr:  `insert into address(id,street, number) values($1, $2, $3)`,
			args:      []interface{}{2, "jalan2", 2},
		},
		{
			name:          "exec context",
			queryType:     query.TypeExecContext,
			queryStr:      `insert into address(id,street, number) values($1, $2, $3)`,
			args:          []interface{}{3, "jalan3", 3},
			prepareOnInit: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			que := query.Query{
				Type:          tc.queryType,
				QueryStr:      tc.queryStr,
				PrepareOnInit: tc.prepareOnInit,
			}

			p, err := newPreparer(db, que)
			require.NoError(t, err)
			defer p.Close()

			err = p.Execute(context.Background(), tc.args...)
			require.NoError(t, err)
		})
	}

}
