package executor

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"

	"github.com/iwanbk/gosqlbencher/query"
)

func TestExecer(t *testing.T) {
	db, err := sql.Open("ramsql", "TestExecer")
	require.NoError(t, err)
	defer db.Close()

	createTableStr := "CREATE TABLE address (id BIGSERIAL PRIMARY KEY, street TEXT, number INT)"
	_, err = db.Exec(createTableStr)
	require.NoError(t, err)

	_, err = db.Exec(`insert into address( street, number)values('jalan', 1)`)
	require.NoError(t, err)

	testCases := []struct {
		name            string
		queryType       string
		queryStr        string
		withPlaceHolder bool
		args            []interface{}
	}{
		{
			name:      "select without placeholder",
			queryType: query.TypeQuery,
			queryStr:  "select * from address where id=%d",
			args:      []interface{}{1},
		},
		{
			name:            "select with placeholder",
			queryType:       query.TypeQueryContext,
			queryStr:        "select * from address where id=$1",
			withPlaceHolder: true,
			args:            []interface{}{1},
		},
		{
			name:      "exec",
			queryType: query.TypeExec,
			queryStr:  `insert into address(street, number) values(%s, %d)`,
			args:      []interface{}{"jalan2", 2},
		},
		{
			name:            "exec context",
			queryType:       query.TypeExecContext,
			queryStr:        `insert into address(street, number) values($1, $2)`,
			args:            []interface{}{"jalan3", 3},
			withPlaceHolder: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			que := query.Query{
				Type:            tc.queryType,
				QueryStr:        tc.queryStr,
				WithPlaceholder: tc.withPlaceHolder,
			}

			ex, err := newExecer(db, que)
			require.NoError(t, err)
			defer ex.Close()

			err = ex.Execute(context.Background(), tc.args...)
			require.NoError(t, err)
		})
	}

}
