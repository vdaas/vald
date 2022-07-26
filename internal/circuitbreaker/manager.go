//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package circuitbreaker

import (
	"context"
	"reflect"
	"sync"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

// NOTE: This variable is for observability package.
//       This will be fixed when refactoring the observability package.
var (
	mu      sync.RWMutex
	metrics = make(map[string]map[State]int64) // map[breaker_name]map[state]count
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
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return bm, nil
}

// Do invokes the breaker matching the given key.
func (bm *breakerManager) Do(ctx context.Context, key string, fn func(ctx context.Context) (interface{}, error)) (val interface{}, err error) {
	var st State
	defer func() {
		mu.Lock()
		if _, ok := metrics[key]; !ok {
			metrics[key] = map[State]int64{
				st: 1,
			}
		} else {
			metrics[key][st]++
		}
		mu.Unlock()
	}()

	// Pre-loading to prevent a lot of object generation.
	obj, ok := bm.m.Load(key)
	if ok {
		val, st, err = obj.(*breaker).do(ctx, fn)
		return val, err
	}

	b, err := newBreaker(key, bm.opts...)
	if err != nil {
		return nil, err
	}
	obj, _ = bm.m.LoadOrStore(key, b)
	val, st, err = obj.(*breaker).do(ctx, fn)
	return val, err
}

func Metrics(_ context.Context) map[string]map[State]int64 {
	mu.RLock()
	defer mu.RUnlock()

	if len(metrics) == 0 {
		return nil
	}
	m := make(map[string]map[State]int64, len(metrics))
	for name, sts := range metrics {
		m[name] = sts
	}
	return m
}
