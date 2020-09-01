package dbr

import (
	dbr "github.com/gocraft/dbr/v2"
)

type (
	Builder       = dbr.Builder
	Connection    = dbr.Connection
	EventReceiver = dbr.EventReceiver
)

type db struct{}

type DBR interface {
	Open(driver, dsn string, log EventReceiver) (*Connection, error)
	Eq(col string, val interface{}) Builder
}

func New() DBR {
	return new(db)
}

func (*db) Open(driver string, dsn string, log EventReceiver) (*Connection, error) {
	return dbr.Open(driver, dsn, log)
}

func (*db) Eq(col string, val interface{}) Builder {
	return dbr.Eq(col, val)
}
