package circuitbreaker

import (
	"context"
	"sync"
)

// CircuitBreaker is a state machine to prevent doing processes that are likely to fail.
type CircuitBreaker interface {
	Do(ctx context.Context, key string, fn func(ctx context.Context) (interface{}, error)) (val interface{}, err error)
}

type breakerManager struct {
	m    sync.Map // breaker group. key: string, value: *breaker.
	opts []BreakerOption
}

// NewCircuitBreaker returns CircuitBreaker object if no error occurs.
func NewCircuitBreaker(opts ...Option) (CircuitBreaker, error) {
	bm := &breakerManager{}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(bm); err != nil {
			return nil, err
		}
	}
	return bm, nil
}

// Do invokes the breaker matching the given key.
func (bm *breakerManager) Do(ctx context.Context, key string, fn func(ctx context.Context) (interface{}, error)) (val interface{}, err error) {
	// Pre-loading to prevent a lot of object generation.
	val, ok := bm.m.Load(key)
	if ok {
		return val.(*breaker).do(ctx, fn)
	}

	b, err := newBreaker(bm.opts...)
	if err != nil {
		return nil, err
	}
	val, _ = bm.m.LoadOrStore(key, b)
	return val.(*breaker).do(ctx, fn)
}
