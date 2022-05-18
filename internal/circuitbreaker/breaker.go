package circuitbreaker

import (
	"context"
	"sync/atomic"
	"time"
)

// counts holds the number of successes/failures.
type counts struct {
	consecutiveSuccesses uint32
	consecutiveFailures  uint32
}

func (c *counts) onSuccess() (n uint32) {
	// This function is called when user defined function succeds.

	// 1. Increment consecutiveSuccesses by using atomic.
	// 2. Clear consecutiveFailures by using atomic.

	// e.g. the following is an example flow.
	/**
		n = atomic.AddUint32(&c.consecutiveSuccesses, 1)
		atomic.StoreUint32(&c.consecutiveFailures, 0)
	**/
	return n
}
func (c *counts) onFailure() (n uint32) {
	// This function is called when user defined function fails.

	// 1. Increment consecutiveFailures by using atomic.
	// 2. Clear consecutiveFailures by using atomic.

	// e.g. the following is an example flow.
	/**
		n = atomic.AddUint32(&c.consecutiveFailures, 1)
		atomic.StoreUint32(&c.consecutiveSuccesses, 0)
	**/
	return n
}
func (c *counts) reset() {
	// 1. Clear consecutiveSuccesses and consecutiveFailures by using atomic

	// e.g. the following is an example flow.
	atomic.StoreUint32(&c.consecutiveFailures, 0)
	atomic.StoreUint32(&c.consecutiveSuccesses, 0)
}

type breaker struct {
	st      atomic.Value // type: state
	count   atomic.Value // type: *count
	tripped int32

	halfOpenTimeout time.Duration
	openTimeout     time.Duration
	expire          time.Time
}

func (b *breaker) do(ctx context.Context, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	// This function is a main flow of breaker.
	// When the current state is "Closed" or "Half-Open", this function executes the given user defined function variable.
	// If no error occurs, above function is treated as successful.

	// 1. Check current state. If current state is "Open", this function returns an error.

	// 2. Execute the given user defined function variable fn.

	// 3. Get an error of above function calling.

	// 4. Call success or fail function based on above error.
	//    - if error is nil or error is ErrCircuitBreakerIgnorable, calls success function.
	//    - if error is not nil, call fail function.

	// e.g. the following is an example flow.
	/**
		if !b.ready() { // when state is "Open"
		    return nil, errors.ErrCircuitBreakerOpenState
		}

		val, err = fn(ctx)
		if err != nil {
		    if !errors.Is(err, ErrCircuitBreakerIgnorable) {
		        b.fail() <------- mark with fail
		        return nil, err.Unwrap()
		    }
		}

		b.success() <-------- mark with success
		return val, nil
	**/
	return nil, nil
}

func (b *breaker) ready() (ok bool) {
	// This function determines breaker is ready or not to call user defined function.

	// 1. If the circuit breaker state is "Closed" or "Half-Open", this function returns true.

	// e.g. the following is an example flow.
	/**
	    st := b.currentState()
	    return st == stateClosed || st == stateHalfOpen
	**/
	return true
}

func (b *breaker) success() {
	// This function focus on the "Half-Open" state.
	// When current state is "Half-Open" state, this function records the number of successful attempts.
	// And the circuit breaker changes to the "Closed" state after a specified number of consecutive operation invocations have been successful.

	// 1. Increment the consecutiveSuccesses and clear consecutiveFailures by using count.onSuccess function
	// 2. Call b.unTrip function

	// e.g. the following is an example flow.
	/**
		switch b.currentState() {
	    case stateHalfOpen:
		    cnt := b.count.Load().(*counts).onSuccess()

		    if cnt >= b.halfOpenSuccessThreshold {
		        b.unTrip() // To "Closed"
		    }
		}
	**/
}

func (b *breaker) fail() {
	// This function focus on the "Closed" and "Half-Open" state.
	// When current state is "Closed" state, this function records the number of failuer attempts.
	// And the circuit breaker changes to the "Open" state if a specified number of consecutive operation invocations have been failed.
	// When current state is "Half-Open" state, the circuit breaker changes to the "Open" state.

	// 1. Increment the consecutiveFailures and clear consecutiveSuccesses by using count.onFailuer function
	// 2. Call b.trip function

	// e.g. the following is an example flow.
	/**
		switch b.currentState() {
		case stateClosed:
		    cnt := b.count.Load().(*counts).onFailure()

		    if cnt >= b.halfClosedFailureThreshold {
			    b.trip() // To "Open"
		    }
		case stateHalfOpen:
	        if b.st.CompareAndSwap(stateHalfOpen, stateOpen) {
		        b.reset(time.Now()) // Reset all counter
		    } else {
		        // The state has been changed to "Closed" during this function calling.
		    }
		}
	**/
}

func (b *breaker) currentState() (st state) {
	// This function returns current state.

	// 1. Returns current state

	// e.g. the following is an example flow.
	/**
		if atomic.LoadInt32(&b.tripped) == 1 {
			now := time.Now()

			switch b.st.Load().(state) {
			case stateOpen:
				if b.expire.Before(now) && b.unTrip() {
					b.reset(now)
					return stateClosed
				}
				return stateOpen
			case stateHalfOpen:
				return stateHalfOpen

			// case stateClosed:
				// When the unTrip method is called immediately after atomic.LoadInt32(&b.tripped)
			}
		}
		return stateClosed
	**/
	return
}

func (b *breaker) reset(now time.Time) {
	// This function resets all counter and expire time.
	/**
		b.count.(*count).reset()
	    b.expire.Store(now)
	**/
}

func (b *breaker) trip() (ok bool) {
	// This function change from untrip state to trip state.

	// 1. change tripped state by using cas operation
	// 2. Force change to "Open" state

	// e.g. the following is an example flow.
	/**
		if atomic.CompareAndSwapInt32(&b.tripped, 0, 1) {
			b.st.Store(stateOpen)
	        return true
		}
	    return false
	**/
	return
}

func (b *breaker) unTrip() (ok bool) {
	// This function change from trip state to untrip state.

	// 1. change tripped state by using cas operation
	// 2. Force change to "Close" state

	// e.g. the following is an example flow.
	/**
		if atomic.CompareAndSwapInt32(&b.tripped, 1, 0) {
			b.st.Store(stateClosed)
	        return true
		}
	    return false
	**/
	return
}
