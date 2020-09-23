//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package mock

import (
	"context"
	"database/sql"
	"time"

	"github.com/vdaas/vald/internal/db/rdb/mysql/dbr"
)

type MockDBR struct {
	OpenFunc func(driver, dsn string, log dbr.EventReceiver) (dbr.Connection, error)
	EqFunc   func(col string, val interface{}) dbr.Builder
}

func (d *MockDBR) Open(driver, dsn string, log dbr.EventReceiver) (dbr.Connection, error) {
	return d.OpenFunc(driver, dsn, log)
}

func (d *MockDBR) Eq(col string, val interface{}) dbr.Builder {
	return d.EqFunc(col, val)
}

type MockSession struct {
	SelectFunc      func(column ...string) dbr.SelectStmt
	BeginFunc       func() (dbr.Tx, error)
	CloseFunc       func() error
	PingContextFunc func(ctx context.Context) error
}

func (s *MockSession) Select(column ...string) dbr.SelectStmt {
	return s.SelectFunc(column...)
}

func (s *MockSession) Begin() (dbr.Tx, error) {
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
	InsertBySqlFunc             func(query string, value ...interface{}) dbr.InsertStmt
	InsertIntoFunc              func(table string) dbr.InsertStmt
	SelectFunc                  func(column ...string) dbr.SelectStmt
	DeleteFromFunc              func(table string) dbr.DeleteStmt
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

func (t *MockTx) InsertBySql(query string, value ...interface{}) dbr.InsertStmt {
	return t.InsertBySqlFunc(query, value...)
}

func (t *MockTx) InsertInto(table string) dbr.InsertStmt {
	return t.InsertIntoFunc(table)
}

func (t *MockTx) Select(column ...string) dbr.SelectStmt {
	return t.SelectFunc(column...)
}

func (t *MockTx) DeleteFrom(table string) dbr.DeleteStmt {
	return t.DeleteFromFunc(table)
}

type MockConn struct {
	NewSessionFunc         func(event dbr.EventReceiver) dbr.Session
	SetConnMaxLifetimeFunc func(d time.Duration)
	SetMaxIdleConnsFunc    func(n int)
	SetMaxOpenConnsFunc    func(n int)
}

func (c *MockConn) NewSession(event dbr.EventReceiver) dbr.Session {
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
	FromFunc        func(table interface{}) dbr.SelectStmt
	WhereFunc       func(query interface{}, value ...interface{}) dbr.SelectStmt
	LimitFunc       func(n uint64) dbr.SelectStmt
	LoadContextFunc func(ctx context.Context, value interface{}) (int, error)
}

func (s *MockSelect) From(table interface{}) dbr.SelectStmt {
	return s.FromFunc(table)
}

func (s *MockSelect) Where(query interface{}, value ...interface{}) dbr.SelectStmt {
	return s.WhereFunc(query, value...)
}

func (s *MockSelect) Limit(n uint64) dbr.SelectStmt {
	return s.LimitFunc(n)
}

func (s *MockSelect) LoadContext(ctx context.Context, value interface{}) (int, error) {
	return s.LoadContextFunc(ctx, value)
}

type MockInsert struct {
	ColumnsFunc     func(column ...string) dbr.InsertStmt
	ExecContextFunc func(ctx context.Context) (sql.Result, error)
	RecordFunc      func(structValue interface{}) dbr.InsertStmt
}

func (s *MockInsert) Colums(column ...string) dbr.InsertStmt {
	return s.ColumnsFunc(column...)
}

func (s *MockInsert) ExecContext(ctx context.Context) (sql.Result, error) {
	return s.ExecContextFunc(ctx)
}

func (s *MockInsert) Record(structValue interface{}) dbr.InsertStmt {
	return s.RecordFunc(structValue)
}
