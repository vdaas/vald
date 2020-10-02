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
	"context"

	dbr "github.com/gocraft/dbr/v2"
)

type selectStmt struct {
	*dbr.SelectStmt
}

type SelectStmt interface {
	// From(table interface{}) SelectStmt
	// Where(query interface{}, value ...interface{}) SelectStmt
	// Limit(n uint64) SelectStmt
	LoadContext(ctx context.Context, value interface{}) (int, error)
}

func (stmt *selectStmt) From(table interface{}) SelectStmt {
	return stmt.SelectStmt.From(table)
}

func (stmt *selectStmt) Where(query interface{}, value ...interface{}) SelectStmt {
	return stmt.SelectStmt.Where(query, value...)
}

func (stmt *selectStmt) Limit(n uint64) SelectStmt {
	return stmt.SelectStmt.Limit(n)
}

func (stmt *selectStmt) LoadContext(ctx context.Context, value interface{}) (int, error) {
	return stmt.SelectStmt.LoadContext(ctx, value)
}
