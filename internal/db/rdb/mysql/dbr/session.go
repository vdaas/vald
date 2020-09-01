package dbr

import (
	dbr "github.com/gocraft/dbr/v2"
)

type session struct{
	*dbr.Session
}

type Session interface {
	Select(column ...string) SelectStmt
	Begin() (tx, error)
}

func NewSession(conn *dbr.Connection, event EventReceiver) Session {
	return &session{
		conn.NewSession(event),
	}
}

func (s *session) Select(column ...string) SelectStmt {
	return s.Select(column...)
}

func (s *session) Begin() (tx, error) {
	return s.Begin()
}

