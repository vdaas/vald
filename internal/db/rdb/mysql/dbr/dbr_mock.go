//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package dbr

import (
	"context"
	"database/sql"
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
	BeginFunc       func() (Tx, error)
	CloseFunc       func() error
	PingContextFunc func(ctx context.Context) error
}

func (s *MockSession) Select(column ...string) SelectStmt {
	return s.SelectFunc(column...)
}

func (s *MockSession) Begin() (Tx, error) {
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

type MockInsert struct {
	ColumnsFunc     func(column ...string) InsertStmt
	ExecContextFunc func(ctx context.Context) (sql.Result, error)
	RecordFunc      func(structValue interface{}) InsertStmt
}

func (s *MockInsert) Columns(column ...string) InsertStmt {
	return s.ColumnsFunc(column...)
}

func (s *MockInsert) ExecContext(ctx context.Context) (sql.Result, error) {
	return s.ExecContextFunc(ctx)
}

func (s *MockInsert) Record(structValue interface{}) InsertStmt {
	return s.RecordFunc(structValue)
}

type MockDelete struct {
	ExecContextFunc func(ctx context.Context) (sql.Result, error)
	WhereFunc       func(query interface{}, value ...interface{}) DeleteStmt
}

func (s *MockDelete) ExecContext(ctx context.Context) (sql.Result, error) {
	return s.ExecContextFunc(ctx)
}

func (s *MockDelete) Where(query interface{}, value ...interface{}) DeleteStmt {
	return s.WhereFunc(query, value...)
}
