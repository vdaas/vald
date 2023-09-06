//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	dbr "github.com/gocraft/dbr/v2"
)

// SelectStmt represents the interface to get data from database.
type SelectStmt interface {
	From(table interface{}) SelectStmt
	Where(query interface{}, value ...interface{}) SelectStmt
	Limit(n uint64) SelectStmt
	LoadContext(ctx context.Context, value interface{}) (int, error)
}

type selectStmt struct {
	*dbr.SelectStmt
}

// From specifies table to select from.
func (stmt *selectStmt) From(table interface{}) SelectStmt {
	stmt.SelectStmt = stmt.SelectStmt.From(table)
	return stmt
}

// Where adds a where condition.
func (stmt *selectStmt) Where(query interface{}, value ...interface{}) SelectStmt {
	stmt.SelectStmt = stmt.SelectStmt.Where(query, value...)
	return stmt
}

// Limit adds a limit condition.
func (stmt *selectStmt) Limit(n uint64) SelectStmt {
	stmt.SelectStmt = stmt.SelectStmt.Limit(n)
	return stmt
}

// LoadContext gets the result of select.
func (stmt *selectStmt) LoadContext(ctx context.Context, value interface{}) (int, error) {
	return stmt.SelectStmt.LoadContext(ctx, value)
}
