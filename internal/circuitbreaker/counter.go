package circuitbreaker

import "sync/atomic"

type Counter interface {
	Successes() int64
	Fails() int64
}

type count struct {
	successes int64
	failures  int64
}

func (c *count) Successes() (n int64) {
	return atomic.LoadInt64(&c.successes)
}

func (c *count) Fails() (n int64) {
	return atomic.LoadInt64(&c.failures)
}

func (c *count) onSuccess() (n int64) {
	n = atomic.AddInt64(&c.successes, 1)
	return n
}

func (c *count) onFail() (n int64) {
	n = atomic.AddInt64(&c.failures, 1)
	return n
}

func (c *count) reset() {
	atomic.StoreInt64(&c.failures, 0)
	atomic.StoreInt64(&c.successes, 0)
}

var _ Counter = (*count)(nil)
