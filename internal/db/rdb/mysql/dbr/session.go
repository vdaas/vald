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
	"context"

	dbr "github.com/gocraft/dbr/v2"
)

// Session represents the interface to handle session.
type Session interface {
	Select(column ...string) SelectStmt
	Begin() (Tx, error)
	Close() error
	PingContext(ctx context.Context) error
}

type session struct {
	*dbr.Session
}

// NewSession creates the session with event and returns the Session interface.
func NewSession(conn Connection, event EventReceiver) Session {
	return conn.NewSession(event)
}

// SeleSelect creates and returns the SelectStmt.
func (sess *session) Select(column ...string) SelectStmt {
	return &selectStmt{
		sess.Session.Select(column...),
	}
}

// Begin creates the transaction using given session.
func (sess *session) Begin() (Tx, error) {
	t, err := sess.Session.Begin()
	return &tx{
		Tx: t,
	}, err
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server to finish.
// Close returns the errro if something goes worng during close.
func (sess *session) Close() error {
	return sess.Session.Close()
}

// PingContext verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (sess *session) PingContext(ctx context.Context) error {
	return sess.Session.PingContext(ctx)
}
