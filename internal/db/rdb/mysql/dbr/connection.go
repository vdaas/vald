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
	"time"

	dbr "github.com/gocraft/dbr/v2"
)

// Connection represents the interface to handle connection of database.
type Connection interface {
	NewSession(event EventReceiver) Session
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

type connection struct {
	*dbr.Connection
}

// NewSession instantiates a Session from Connection.
func (conn *connection) NewSession(event EventReceiver) Session {
	return &session{
		conn.Connection.NewSession(event),
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (conn *connection) SetConnMaxLifetime(d time.Duration) {
	conn.Connection.SetConnMaxLifetime(d)
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
func (conn *connection) SetMaxIdleConns(n int) {
	conn.Connection.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (conn *connection) SetMaxOpenConns(n int) {
	conn.Connection.SetMaxOpenConns(n)
}
