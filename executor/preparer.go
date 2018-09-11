package executor

import (
	"context"
	"database/sql"

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

	_, err = stmt.ExecContext(ctx, args...)
	return err
}

func (p *preparer) Close() error {
	if p.stmt == nil {
		return nil
	}
	return p.stmt.Close()
}
