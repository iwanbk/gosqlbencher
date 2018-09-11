package executor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iwanbk/gosqlbencher/query"
)

type preparer struct {
	stmt      *sql.Stmt
	queryType string
	stmtStr   string
	db        *sql.DB
}

func newPreparer(db *sql.DB, query query.Query) (*preparer, error) {
	p := preparer{
		stmtStr:   query.QueryStr,
		queryType: query.Type,
		db:        db,
	}
	if !query.PrepareOnInit {
		return &p, nil
	}

	if err := p.initStmt(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *preparer) initStmt() error {
	stmt, err := p.db.Prepare(p.stmtStr)
	if err != nil {
		return err
	}

	p.stmt = stmt
	return nil
}

// Execute implements Executor.Execute
func (p *preparer) Execute(ctx context.Context, args ...interface{}) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	if p.stmt == nil {
		stmt, err = p.db.Prepare(p.stmtStr)
		if err != nil {
			return err
		}
		defer stmt.Close()
	} else {
		stmt = p.stmt
	}

	switch p.queryType {
	case query.TypeExec:
		_, err = stmt.Exec(args...)
	case query.TypeExecContext:
		_, err = stmt.ExecContext(ctx, args...)
	case query.TypeQuery, query.TypeQueryContext:
		var rows *sql.Rows
		if p.queryType == query.TypeQuery {
			rows, err = stmt.Query(args...)
		} else {
			rows, err = stmt.QueryContext(ctx, args...)
		}

		if err != nil {
			return err
		}
		for rows.Next() {
		}
		err = rows.Close()
	default:
		err = fmt.Errorf("preparer: unsupported query type: %v", p.queryType)
	}

	return err
}

// Close implements Executor.Close
func (p *preparer) Close() error {
	if p.stmt == nil {
		return nil
	}
	return p.stmt.Close()
}
