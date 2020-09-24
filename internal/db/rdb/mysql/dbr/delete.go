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

	dbr "github.com/gocraft/dbr/v2"
)

type DeleteStmt interface {
	ExecContext(ctx context.Context) (sql.Result, error)
	Where(query interface{}, value ...interface{}) DeleteStmt
}

type deleteStmt struct {
	*dbr.DeleteStmt
}

func (stmt *deleteStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return stmt.DeleteStmt.ExecContext(ctx)
}

func (stmt *deleteStmt) Where(query interface{}, value ...interface{}) DeleteStmt {
	stmt.DeleteStmt = stmt.DeleteStmt.Where(query, value...)
	return stmt
}
