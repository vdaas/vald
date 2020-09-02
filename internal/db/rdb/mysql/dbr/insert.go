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
	"database/sql"

	"github.com/gocraft/dbr"
)

type insertStmt struct {
	*dbr.InsertStmt
}

type InsertStmt interface {
	// Columns(column ...string) InsertStmt
	ExecContext(ctx context.Context) (sql.Result, error)
	// Record(structValue interface{}) InsertStmt
}

func (stmt *insertStmt) Columns(column ...string) InsertStmt {
	return stmt.InsertStmt.Columns(column...)
}

func (stmt *insertStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return stmt.ExecContext(ctx)
}

func (stmt *insertStmt) Record(structValue interface{}) InsertStmt {
	return stmt.InsertStmt.Record(structValue)
}
