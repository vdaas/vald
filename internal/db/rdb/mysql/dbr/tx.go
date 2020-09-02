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
package dbr

import (
	dbr "github.com/gocraft/dbr/v2"
)

type tx struct {
	*dbr.Tx
}

type Tx interface {
	Commit() error
	Rollback() error
	RollbackUnlessCommitted()
	InsertBySql(query string, value ...interface{}) InsertStmt
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

func (t *tx) InsertBySql(query string, value ...interface{}) InsertStmt {
	return t.InsertBySql(query, value...)
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
