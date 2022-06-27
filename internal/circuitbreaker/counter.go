package circuitbreaker

import "sync/atomic"

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
