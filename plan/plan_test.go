package plan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	pl, err := Read("../examples/insert.plan.postgre.yaml")
	require.NoError(t, err)

	require.NotEmpty(t, pl.Queries)
	require.NotEmpty(t, pl.DataSourceName)
	require.NotZero(t, pl.NumWorker)
}

func TestReadFailed(t *testing.T) {
	_, err := Read("plan.yaml")
	require.Error(t, err)
}
