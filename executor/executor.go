package executor

import (
	"context"
	"database/sql"

	"github.com/iwanbk/gosqlbencher/query"
)

// Executor is interface that should be implemented by
// the query runner
type Executor interface {
	Execute(ctx context.Context, args ...interface{}) error
	Close() error
}

func New(db *sql.DB, q query.Query) (Executor, error) {
	switch {
	case q.Type == "exec" && q.Prepare:
		return newPreparer(db, q)
	default:
		return newExecer(db, q)
	}
}
