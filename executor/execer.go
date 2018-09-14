package executor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iwanbk/gosqlbencher/query"
)

type execer struct {
	queryStr  string
	queryType query.TypeType
	db        *sql.DB
	execute   func(ctx context.Context, args ...interface{}) error
}

func newExecer(db *sql.DB, query query.Query) (*execer, error) {
	ex := &execer{
		db:        db,
		queryStr:  query.QueryStr,
		queryType: query.Type,
	}
	if query.WithPlaceholder {
		ex.execute = ex.executeWithPlaceholder
	} else {
		ex.execute = ex.executeNoPlaceholder
	}
	return ex, nil
}

// Execute implements Executor.Execute
func (ex *execer) Execute(ctx context.Context, args ...interface{}) error {
	return ex.execute(ctx, args...)
}

// Close implements Executor.Close
func (ex *execer) Close() error {
	return nil
}

// executeNoPlaceholder execute query with no placeholder in the query
func (ex *execer) executeNoPlaceholder(ctx context.Context, args ...interface{}) error {
	q := fmt.Sprintf(ex.queryStr, args...)
	return ex.executeGeneric(ctx, q)
}

// executeWithPlaceholder execute query with  placeholder in the query
func (ex *execer) executeWithPlaceholder(ctx context.Context, args ...interface{}) error {
	return ex.executeGeneric(ctx, ex.queryStr, args...)
}

func (ex *execer) executeGeneric(ctx context.Context, queryStr string, args ...interface{}) error {
	var err error
	switch ex.queryType {
	case query.TypeExec:
		_, err = ex.db.Exec(queryStr, args...)
	case query.TypeExecContext:
		_, err = ex.db.ExecContext(ctx, queryStr, args...)
	case query.TypeQuery, query.TypeQueryContext:
		var rows *sql.Rows
		if ex.queryType == query.TypeQuery {
			rows, err = ex.db.Query(queryStr, args...)
		} else {
			rows, err = ex.db.QueryContext(ctx, queryStr, args...)
		}
		if err != nil {
			return err
		}
		for rows.Next() {
		}
		err = rows.Close()
	default:
		err = fmt.Errorf("execer 1: unsupported query type: %v", ex.queryType)
	}
	return err
}
