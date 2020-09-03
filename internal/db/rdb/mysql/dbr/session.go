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

type Session interface {
	Select(column ...string) SelectStmt
	Begin() (*tx, error)
	Close() error
	PingContext(ctx context.Context) error
}

type session struct {
	*dbr.Session
}

func NewSession(conn *Connection, event EventReceiver) Session {
	return &session{
		conn.NewSession(event),
	}
}

func (sess *session) Select(column ...string) SelectStmt {
	return &selectStmt{
		sess.Session.Select(column...),
	}
}

func (sess *session) Begin() (*tx, error) {
	t, err := sess.Session.Begin()
	return &tx{
		t,
	}, err
}

func (sess *session) Close() error {
	return sess.Session.Close()
}

func (sess *session) PingContext(ctx context.Context) error {
	return sess.Session.PingContext(ctx)
}
