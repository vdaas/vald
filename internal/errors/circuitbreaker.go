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
	var errstr string
	if e.err != nil {
		errstr = ": " + e.err.Error()
	}
	return "circuit breaker ignorable error" + errstr
}

func (e *ErrCircuitBreakerIgnorable) Unwrap() error {
	return e.err
}

type ErrCircuitBreakerMarkWithSuccess struct {
	err error
}

func NewErrCircuitBreakerMarkWithSuccess(err error) error {
	if err == nil {
		return nil
	}
	return &ErrCircuitBreakerMarkWithSuccess{
		err: err,
	}
}

func (e *ErrCircuitBreakerMarkWithSuccess) Error() string {
	var errstr string
	if e.err != nil {
		errstr = ": " + e.err.Error()
	}
	return "circuit breaker mark with success" + errstr
}

func (e *ErrCircuitBreakerMarkWithSuccess) Unwrap() error {
	return e.err
}
