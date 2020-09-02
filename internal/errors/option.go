package errors

var (
	ErrCassandraInvalidPort = func(p int) error {
		return Errorf("invalid port number: %s", p)
	}
)
