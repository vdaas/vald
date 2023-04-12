//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
//
//	This will be fixed when refactoring the observability package.
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

	var br *breaker
	// Pre-loading to prevent a lot of object generation.
	obj, ok := bm.m.Load(key)
	if !ok {
		br, err = newBreaker(key, bm.opts...)
		if err != nil {
			return nil, err
		}
		obj, _ = bm.m.LoadOrStore(key, br)
	}
	br, ok = obj.(*breaker)
	if !ok {
		br, err = newBreaker(key, bm.opts...)
		if err != nil {
			return nil, err
		}
		bm.m.Store(key, br)
	}
	val, st, err = br.do(ctx, fn)
	if err != nil {
		switch st {
		case StateClosed:
			err = errors.Wrapf(err, "circuitbreaker state is %s, this error is not caused by circuitbreaker", st.String())
		case StateOpen:
			if !errors.Is(err, errors.ErrCircuitBreakerOpenState) {
				err = errors.Wrap(err, errors.ErrCircuitBreakerOpenState.Error())
			}
		case StateHalfOpen:
			if !errors.Is(err, errors.ErrCircuitBreakerHalfOpenFlowLimitation) {
				err = errors.Wrap(err, errors.ErrCircuitBreakerHalfOpenFlowLimitation.Error())
			}
		}
		return val, err
	}
	return val, nil
}

func Metrics(context.Context) (ms map[string]map[State]int64) {
	mu.RLock()
	defer mu.RUnlock()

	if len(metrics) == 0 {
		return nil
	}
	ms = make(map[string]map[State]int64, len(metrics))
	for name, state := range metrics {
		sts := make(map[State]int64, len(state))
		for st, cnt := range state {
			sts[st] = cnt
		}
		ms[name] = sts
	}
	return ms
}
