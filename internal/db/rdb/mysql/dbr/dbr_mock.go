package dbr

import (
	"context"
	"time"
)

type MockDBR struct {
	OpenFunc func(driver, dsn string, log EventReceiver) (Connection, error)
	EqFunc   func(col string, val interface{}) Builder
}

func (d *MockDBR) Open(driver, dsn string, log EventReceiver) (Connection, error) {
	return d.OpenFunc(driver, dsn, log)
}

func (d *MockDBR) Eq(col string, val interface{}) Builder {
	return d.EqFunc(col, val)
}

type MockSession struct {
	SelectFunc      func(column ...string) SelectStmt
	BeginFunc       func() (*tx, error)
	CloseFunc       func() error
	PingContextFunc func(ctx context.Context) error
}

func (s *MockSession) Select(column ...string) SelectStmt {
	return s.SelectFunc(column...)
}

func (s *MockSession) Begin() (*tx, error) {
	return s.BeginFunc()
}

func (s *MockSession) Close() error {
	return s.CloseFunc()
}

func (s *MockSession) PingContext(ctx context.Context) error {
	return s.PingContextFunc(ctx)
}

type MockTx struct {
	CommitFunc                  func() error
	RollbackFunc                func() error
	RollbackUnlessCommittedFunc func()
	InsertBySqlFunc             func(query string, value ...interface{}) InsertStmt
	InsertIntoFunc              func(table string) InsertStmt
	SelectFunc                  func(column ...string) SelectStmt
	DeleteFromFunc              func(table string) DeleteStmt
}

func (t *MockTx) Commit() error {
	return t.CommitFunc()
}

func (t *MockTx) Rollback() error {
	return t.RollbackFunc()
}

func (t *MockTx) RollbackUnlessCommitted() {
	t.RollbackUnlessCommittedFunc()
}

func (t *MockTx) InsertBySql(query string, value ...interface{}) InsertStmt {
	return t.InsertBySqlFunc(query, value...)
}

func (t *MockTx) InsertInto(table string) InsertStmt {
	return t.InsertIntoFunc(table)
}

func (t *MockTx) Select(column ...string) SelectStmt {
	return t.SelectFunc(column...)
}

func (t *MockTx) DeleteFrom(table string) DeleteStmt {
	return t.DeleteFromFunc(table)
}

type MockConn struct {
	NewSessionFunc         func(event EventReceiver) Session
	SetConnMaxLifetimeFunc func(d time.Duration)
	SetMaxIdleConnsFunc    func(n int)
	SetMaxOpenConnsFunc    func(n int)
}

func (c *MockConn) NewSession(event EventReceiver) Session {
	return c.NewSessionFunc(event)
}

func (c *MockConn) SetConnMaxLifetime(d time.Duration) {
	c.SetConnMaxLifetimeFunc(d)
}

func (c *MockConn) SetMaxIdleConns(n int) {
	c.SetMaxIdleConnsFunc(n)
}

func (c *MockConn) SetMaxOpenConns(n int) {
	c.SetMaxOpenConnsFunc(n)
}

type MockSelect struct {
	FromFunc        func(table interface{}) SelectStmt
	WhereFunc       func(query interface{}, value ...interface{}) SelectStmt
	LimitFunc       func(n uint64) SelectStmt
	LoadContextFunc func(ctx context.Context, value interface{}) (int, error)
}

func (s *MockSelect) From(table interface{}) SelectStmt {
	return s.FromFunc(table)
}

func (s *MockSelect) Where(query interface{}, value ...interface{}) SelectStmt {
	return s.WhereFunc(query, value...)
}

func (s *MockSelect) Limit(n uint64) SelectStmt {
	return s.LimitFunc(n)
}

func (s *MockSelect) LoadContext(ctx context.Context, value interface{}) (int, error) {
	return s.LoadContextFunc(ctx, value)
}

