package dbr

import (
	"context"
	"database/sql"

	dbr "github.com/gocraft/dbr/v2"
)

type deleteStmt struct {
	*dbr.DeleteStmt
}

type DeleteStmt interface {
	ExecContext(ctx context.Context) (sql.Result, error)
	Where(query interface{}, value ...interface{}) deleteStmt
}

func (stmt *deleteStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return stmt.ExecContext(ctx)
}

func (stmt *deleteStmt) Where(query string, value ...interface{}) deleteStmt {
	return stmt.Where(query, value)
}
