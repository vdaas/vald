package circuitbreaker

import (
	"context"
	"sync"
)

type CircuitBreaker interface {
	Do(ctx context.Context, key string, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error)
}

// breakerGroup is a group of breakers by name.
type breakerGroup struct {
	m sync.Map // key: string, value: *breaker
}

func NewCircuitBreaker(opts ...Option) (CircuitBreaker, error) {
	return nil, nil
}

func (bg *breakerGroup) Do(ctx context.Context, key string, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	// The breakerGroup manages many breakers.
	// So this function invokes the breaker matching the key.
	// The key is RPC name or endpoint etc.

	// 1. Get breaker from bg.m with given key. If breaker dose not exist, create new braker.
	// 2. Call braker.do function.

	// e.g. the following is an example flow code.
	/**
	    val, _ := bg.m.LoadOrStore(key, newBreaker())
	    return val.(*breaker).do(ctx, fn)
	**/
	return
}
