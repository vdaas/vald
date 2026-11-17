// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package circuitbreaker

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// NOTE: This variable is for observability package.
//
//	This will be fixed when refactoring the observability package.
var (
	// metrics counts Do calls per breaker name and resulting state. Do runs on
	// every gRPC call that has a breaker, so the counters are lock-free per key
	// instead of sharing one package-wide mutex.
	metrics sync.Map[string, *stateCounts] //nolint:gochecknoglobals // process-wide counters shared by every breakerManager and read by the observability package; see the NOTE above
)

// CircuitBreaker is a state machine to prevent doing processes that are likely to fail.
type CircuitBreaker interface {
	Do(ctx context.Context, key string, fn func(ctx context.Context) (any, error)) (val any, err error)
}

// stateCounts holds one counter per State (StateUnknown..StateHalfOpen).
type stateCounts [StateHalfOpen + 1]atomic.Int64

type breakerManager struct {
	m    sync.Map[string, *breaker]
	opts []BreakerOption
}

// New returns CircuitBreaker object if no error occurs.
func New(opts ...Option) (CircuitBreaker, error) {
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
func (bm *breakerManager) Do(
	ctx context.Context, key string, fn func(ctx context.Context) (any, error),
) (val any, err error) {
	var (
		st State
		br *breaker
	)
	defer func() {
		// EvictAddr ran while this call was in flight: do not re-create the
		// counters of the evicted key.
		if br != nil {
			if cur, ok := bm.m.Load(key); !ok || cur != br {
				return
			}
		}
		c, ok := metrics.Load(key)
		if !ok {
			c, _ = metrics.LoadOrStore(key, new(stateCounts))
		}
		if int(st) < len(c) {
			c[st].Add(1)
		}
	}()

	// Pre-loading to prevent a lot of object generation.
	br, ok := bm.m.Load(key)
	if !ok {
		var nb *breaker
		nb, err = newBreaker(key, bm.opts...)
		if err != nil {
			return nil, err
		}
		br, _ = bm.m.LoadOrStore(key, nb)
	}

	val, st, err = br.do(ctx, fn)
	if err != nil {
		switch st {
		case StateClosed:
			err = errors.Wrapf(err, "circuitbreaker state is %s, this error is not caused by circuitbreaker", st.String())
		case StateOpen:
			if !errors.Is(err, errors.ErrCircuitBreakerOpenState) {
				err = errors.Join(err, errors.ErrCircuitBreakerOpenState)
			}
		case StateHalfOpen:
			if !errors.Is(err, errors.ErrCircuitBreakerHalfOpenFlowLimitation) {
				err = errors.Join(err, errors.ErrCircuitBreakerHalfOpenFlowLimitation)
			}
		}
		return val, err
	}
	return val, nil
}

// EvictAddr removes the breakers and metrics counters of every key that targets
// addr, so the state of removed addresses does not pile up. The gRPC client
// calls it when it disconnects addr. Its keys are "<addr>" or "<method>/<addr>"
// (grpc.WrapGRPCMethod). A later Do for an evicted key starts a new, closed
// breaker and new counters.
func (bm *breakerManager) EvictAddr(addr string) {
	if addr == "" {
		return
	}
	suffix := "/" + addr
	targets := func(key string) bool {
		return key == addr || strings.HasSuffix(key, suffix)
	}
	bm.m.Range(func(key string, _ *breaker) bool {
		if targets(key) {
			bm.m.Delete(key)
		}
		return true
	})
	metrics.Range(func(key string, _ *stateCounts) bool {
		if targets(key) {
			metrics.Delete(key)
		}
		return true
	})
}

func Metrics(context.Context) (ms map[string]map[State]int64) {
	metrics.Range(func(name string, c *stateCounts) bool {
		// Report only the states that occurred, like the previous map did.
		var sts map[State]int64
		for st := range c {
			if n := c[st].Load(); n > 0 {
				if sts == nil {
					sts = make(map[State]int64, len(c))
				}
				sts[State(st)] = n
			}
		}
		if sts != nil {
			if ms == nil {
				ms = make(map[string]map[State]int64)
			}
			ms[name] = sts
		}
		return true
	})
	return ms
}
