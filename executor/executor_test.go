package executor

import (
	"database/sql"
	"testing"

	"github.com/iwanbk/gosqlbencher/query"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

func TestNewPreparer(t *testing.T) {
	db, err := sql.Open("ramsql", "TestPrepare")
	require.NoError(t, err)
	defer db.Close()

	// check it
	ex, err := New(db, query.Query{
		Type:    query.TypeQuery,
		Prepare: true,
	})
	require.NoError(t, err)
	require.IsType(t, &preparer{}, ex)
}

func TestNewExecer(t *testing.T) {
	db, err := sql.Open("ramsql", "TestPrepare")
	require.NoError(t, err)
	defer db.Close()

	// check it
	ex, err := New(db, query.Query{
		Type:    query.TypeQuery,
		Prepare: false,
	})
	require.NoError(t, err)
	require.IsType(t, &execer{}, ex)
}
