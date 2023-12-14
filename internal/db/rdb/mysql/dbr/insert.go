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
	"database/sql"

	dbr "github.com/gocraft/dbr/v2"
)

// InsertStmt represents the interface to insert data to database.
type InsertStmt interface {
	Columns(column ...string) InsertStmt
	ExecContext(ctx context.Context) (sql.Result, error)
	Record(structValue interface{}) InsertStmt
}

type insertStmt struct {
	*dbr.InsertStmt
}

// Columns set colums to the insertStmt.
func (stmt *insertStmt) Columns(column ...string) InsertStmt {
	stmt.InsertStmt = stmt.InsertStmt.Columns(column...)
	return stmt
}

// ExecContext execure inserting to the database.
func (stmt *insertStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return stmt.InsertStmt.ExecContext(ctx)
}

// Record adds a tuple for columns from a struct.
func (stmt *insertStmt) Record(structValue interface{}) InsertStmt {
	stmt.InsertStmt = stmt.InsertStmt.Record(structValue)
	return stmt
}
