package dbr

import (
	"context"
	"database/sql"

	"github.com/gocraft/dbr"
)

type insertStmt struct {
	*dbr.InsertStmt
}

type InsertStmt interface {
	Columns(column ...string) InsertStmt
	ExecContext(ctx context.Context) (sql.Result, error)
	Record(structValue interface{}) InsertStmt
}

func (stmt *insertStmt) Columns(column ...string) InsertStmt {
	return stmt.Columns(column...)
}

func (stmt *insertStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return stmt.ExecContext(ctx)
}

func (stmt *insertStmt) Record(structValue interface{}) InsertStmt {
	return stmt.Record(structValue)
}
