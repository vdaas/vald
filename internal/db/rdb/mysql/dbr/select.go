package dbr

import (
	"context"

	dbr "github.com/gocraft/dbr/v2"
)

type selectStmt struct {
	*dbr.SelectStmt
}

type SelectStmt interface {
	From(table interface{}) SelectStmt
	Where(query interface{}, value ...interface{}) SelectStmt
	Limit(n uint64) SelectStmt
	LoadContext(ctx context.Context, value interface{}) (int, error)
}

func (stmt *selectStmt) From(table interface{}) SelectStmt {
	return stmt.From(table)
}

func (stmt *selectStmt) Where(query interface{}, value ...interface{}) SelectStmt {
	return stmt.Where(query, value)
}

func (stmt *selectStmt) Limit(n uint64) SelectStmt {
	return stmt.Limit(n)
}

func (stmt *selectStmt) LoadContext(ctx context.Context, value interface{}) (int, error) {
	return stmt.LoadContext(ctx, value)
}
