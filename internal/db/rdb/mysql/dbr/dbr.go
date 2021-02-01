//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// DBR repreesnts the interface to create connection to MySQL.
type DBR interface {
	Open(driver, dsn string, log EventReceiver) (Connection, error)
	Eq(col string, val interface{}) Builder
}

type (
	// Builder is a type alias of dbr.Builder.
	Builder = dbr.Builder

	// EventReceiver is a type alias of dbr.EventReceiver.
	EventReceiver = dbr.EventReceiver

	// TracingEventReceiver is a type alias of dbr.TracingEventReceiver.
	TracingEventReceiver = dbr.TracingEventReceiver

	// NullEventReceiver is a type alias of dbr.NullEventReceiver.
	NullEventReceiver = dbr.NullEventReceiver
)

type db struct{}

// New returns the new db struct.
func New() DBR {
	return new(db)
}

// Open returns the connection of db.
// When any error occures, it will return the error.
func (*db) Open(driver, dsn string, log EventReceiver) (Connection, error) {
	conn, err := dbr.Open(driver, dsn, log)
	return &connection{
		conn,
	}, err
}

// Eq returns the built SQL statement made from col name and the value.
func (*db) Eq(col string, val interface{}) Builder {
	return dbr.Eq(col, val)
}
