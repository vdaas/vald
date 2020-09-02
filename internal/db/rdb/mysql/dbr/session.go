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

type session struct {
	*dbr.Session
}

type Session interface {
	Select(column ...string) SelectStmt
	Begin() (*dbr.Tx, error)
	Close()
	PingContext(ctx context.Context) error
}

func NewSession(conn *dbr.Connection, event EventReceiver) Session {
	return &session{
		conn.NewSession(event),
	}
}

func (s *session) Select(column ...string) SelectStmt {
	return s.Session.Select(column...)
}

func (s *session) Begin() (*dbr.Tx, error) {
	return s.Session.Begin()
}

func (s *session) Close() {
	s.Session.Close()
}

func (s *session) PingContext(ctx context.Context) error {
	return s.Session.PingContext(ctx)
}
