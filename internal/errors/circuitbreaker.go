package errors

import (
	"errors"
)

var (
	// ErrCircuitBreakerTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests.
	ErrCircuitBreakerTooManyRequests = errors.New("too many requests")
	// ErrCircuitBreakerOpenState is returned when the CB state is open.
	ErrCircuitBreakerOpenState = errors.New("circuit breaker is open")
)

type ErrCircuitBreakerIgnorable struct {
	err error
}

func NewErrCircuitBreakerIgnorable(err error) error {
	if err == nil {
		return nil
	}
	return &ErrCircuitBreakerIgnorable{
		err: err,
	}
}

func (e *ErrCircuitBreakerIgnorable) Error() string {
	return "dose not mark error: " + e.Error()
}

func (e *ErrCircuitBreakerIgnorable) Unwrap() error {
	return e.err
}
