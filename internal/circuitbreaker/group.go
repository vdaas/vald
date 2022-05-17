package circuitbreaker

import (
	"context"
	"sync"
)

type CircuitBreaker interface {
	Do(ctx context.Context, key string, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error)
}

// breakerGroup (or breakerPannel) is a group of breakers by name.
type breakerGroup struct {
	m sync.Map // key: string, value: *breaker
}

func NewCircuitBreaker(opts ...Option) (CircuitBreaker, error) {
	return nil, nil
}

func (bg *breakerGroup) Do(ctx context.Context, key string, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	// 1. Get breaker from bg.m with given key. if breaker dose not exist, create new braker.
	// 2. Call braker.do function.
	return nil, nil
}
