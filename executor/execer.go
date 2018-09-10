package executor

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/iwanbk/gosqlbencher/query"
)

var (
	execerSupportedType = map[string]struct{}{
		"exec": struct{}{},
	}
)

type execer struct {
	queryStr  string
	queryType string
	db        *sql.DB
	execute   func(ctx context.Context, args ...interface{}) error
}

func newExecer(db *sql.DB, query query.Query) (*execer, error) {
	qt := strings.ToLower(query.Type)

	if _, ok := execerSupportedType[qt]; !ok {
		return nil, fmt.Errorf("execer: unsupported query type: %s", query.Type)
	}

	ex := &execer{
		db:        db,
		queryStr:  query.QueryStr,
		queryType: qt,
	}
	if query.WithPlaceholder {
		ex.execute = ex.executeWithPlaceholder
	} else {
		ex.execute = ex.executeNoPlaceholder
	}
	return ex, nil
}

func (ex *execer) Execute(ctx context.Context, args ...interface{}) error {
	return ex.execute(ctx, args...)
}

func (ex *execer) executeNoPlaceholder(ctx context.Context, args ...interface{}) error {
	q := fmt.Sprintf(ex.queryStr, args...)
	return ex.executeGeneric(ctx, q)
}

func (ex *execer) executeWithPlaceholder(ctx context.Context, args ...interface{}) error {
	return ex.executeGeneric(ctx, ex.queryStr, args...)
}

func (ex *execer) executeGeneric(ctx context.Context, query string, args ...interface{}) error {
	var err error
	switch ex.queryType {
	case "exec":
		_, err = ex.db.ExecContext(ctx, query, args...)
	default:
		err = fmt.Errorf("execer: unsupported query type: %v", ex.queryType)
	}
	return err
}
