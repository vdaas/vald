package dbr

import (
	dbr "github.com/gocraft/dbr/v2"
)

type tx struct {
	*dbr.Tx
}

type Tx = interface {
	Commit() error
	Rollback() error
	RollbackUnlessCommitted()
	InsertBySql(query string, value interface{}) InsertStmt
	InsertInto(table string) InsertStmt
	Select(column ...string) SelectStmt
	DeleteFrom(table string) DeleteStmt
}

func (t *tx) Commit() error {
	return t.Commit()
}

func (t *tx) Rollback() error {
	return t.Rollback()
}

func (t *tx) RollbackUnlessCommitted() {
	t.RollbackUnlessCommitted()
}

func (t *tx) InsertBySql(query string, value interface{}) InsertStmt {
	return t.InsertBySql(query, value)
}

func (t *tx) InsertInto(table string) InsertStmt {
	return t.InsertInto(table)
}

func (t *tx) Select(column ...string) SelectStmt {
	return t.Select(column...)
}

func (t *tx) DeleteFrom(table string) DeleteStmt {
	return t.DeleteFrom(table)
}
