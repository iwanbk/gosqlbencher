package executor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iwanbk/gosqlbencher/query"
)

// Executor is
type Executor interface {
	Execute(ctx context.Context, args ...interface{}) error
}

func New(db *sql.DB, q query.Query) (Executor, error) {
	switch {
	case q.Type == "exec" && q.Prepare:
		return newPreparer(db, q)
	default:
		return nil, fmt.Errorf("unsupported query: %v", q)
	}
}
