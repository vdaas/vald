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
	st    atomic.Value // type: *stater
	count atomic.Value // type: *count

	// or mu sync.Mutex
}

func (bg *breaker) do(ctx context.Context, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	// 1. Check current state. If current state is open, returns an error.

	// 2. Execute the given function variable fn.

	// 3. Get an error of above function calling.

	// 4. Call success or failuer function based on above error.
	//    - if error is nil or error is ErrCircuitBreakerIgnorable, calls success function.
	//    - if error is not nil, call failuer function.
	return nil, nil
}

func (bg *breaker) success() {
	// 1. Increment the consecutiveSuccesses and clear consecutiveFailures. (Call count.onSuccess function)
	// 2. Call onSuccess function (Notify success event to current state) of st object.
}

func (bg *breaker) failuer() {
	// 1. Increment the consecutiveFailures and clear consecutiveSuccesses. (Call count.onFailuer function)
	// 2. Call onFailuer function (Notify failuer event to current state) of st object.
}

func (bg *breaker) setState(st stater) {
	// 1. Change current state with given state.
}

func (bg *breaker) isReachedSuccessThreshold() (ok bool) {
	// 1. Check whether success thresholds have been reached
	return
}
