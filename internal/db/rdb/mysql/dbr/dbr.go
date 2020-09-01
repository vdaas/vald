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
