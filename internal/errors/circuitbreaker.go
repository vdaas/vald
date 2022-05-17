package errors

import "errors"

var (
	// ErrCircuitBreakerTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests.
	ErrCircuitBreakerTooManyRequests = errors.New("too many requests")
	// ErrCircuitBreakerOpenState is returned when the CB state is open.
	ErrCircuitBreakerOpenState = errors.New("circuit breaker is open")
)
