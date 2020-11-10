package cassandra

import "github.com/gocql/gocql"

type MockClusterConfig struct {
	CreateSessionFunc func() (*gocql.Session, error)
}

func (m *MockClusterConfig) CreateSession() (*gocql.Session, error) {
	return m.CreateSessionFunc()
}
