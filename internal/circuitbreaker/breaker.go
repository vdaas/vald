package circuitbreaker

import (
	"context"
	"sync/atomic"
)

type counts struct {
	consecutiveSuccesses uint32
	consecutiveFailures  uint32
}

func (c *counts) onSuccess() {
	// 1. Increment consecutiveSuccesses by using atomic.
	// 2. Clear consecutiveFailures by using atomic.
}
func (c *counts) onFailure() {
	// 1. Increment consecutiveFailures by using atomic.
	// 2. Clear consecutiveFailures by using atomic.
}
func (c *counts) clear() {
	// 1. Clear consecutiveSuccesses and consecutiveFailures by using atomic
}

type breaker struct {
	st      atomic.Value // type: state
	count   atomic.Value // type: *count
	tripped int32
}

func (b *breaker) do(ctx context.Context, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	// 1. Check current state. If current state is open, returns an error.

	// 2. Execute the given function variable fn.

	// 3. Get an error of above function calling.

	// 4. Call success or failuer function based on above error.
	//    - if error is nil or error is ErrCircuitBreakerIgnorable, calls success function.
	//    - if error is not nil, call fail function.
	return nil, nil
}

func (b *breaker) ready() (ok bool) {
	// 1. Ready will return true if the circuit breaker is ready(closed or half-open) to call the function.
	return true
}

func (b *breaker) currentState() (st state) {
	// 1. Returns current state
	if atomic.LoadInt32(&b.tripped) == 1 {

		return stateHalfOpen // or return stateOpen
	}
	return stateClosed
}

func (b *breaker) success() {
	// 1. Increment the consecutiveSuccesses and clear consecutiveFailures. (Call count.onSuccess function)
	// 2. Call onSuccess function (Notify success event to current state) of st object.
}

func (b *breaker) fail() {
	// 1. Increment the consecutiveFailures and clear consecutiveSuccesses. (Call count.onFailuer function)
	// 2. Call onFailuer function (Notify failuer event to current state) of st object.
}

func (b *breaker) setState(st stater) {
	// 1. Change current state with given state.
}
