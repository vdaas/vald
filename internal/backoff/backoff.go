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

// Package backoff provides backoff function controller
package backoff

import (
	"context"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/ctxkey"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/rand"
)

type backoff struct {
	wg                    sync.WaitGroup
	backoffFactor         float64
	initialDuration       float64
	jittedInitialDuration float64
	jitterLimit           float64
	durationLimit         float64
	maxDuration           float64
	maxRetryCount         int
	backoffTimeLimit      time.Duration
	errLog                bool
	metricsEnabled        bool

	mu      sync.RWMutex
	metrics map[string]int64
}

// Backoff represents an interface to handle backoff operation.
type Backoff interface {
	Do(context.Context, func(ctx context.Context) (interface{}, bool, error)) (interface{}, error)
	Metrics(ctx context.Context) map[string]int64
	Close()
}

const traceTag = "vald/internal/backoff/Backoff.Do/retry"

// New creates the new backoff with option.
func New(opts ...Option) Backoff {
	b := new(backoff)
	for _, opt := range append(defaultOptions, opts...) {
		opt(b)
	}
	if b.backoffFactor < 1 {
		b.backoffFactor = 1.1
	}
	b.durationLimit = b.maxDuration / b.backoffFactor
	b.jittedInitialDuration = b.addJitter(b.initialDuration)
	b.metrics = make(map[string]int64)

	return b
}

// Do tries to backoff using the input function and returns the response and error.
func (b *backoff) Do(ctx context.Context, f func(ctx context.Context) (val interface{}, retryable bool, err error)) (res interface{}, err error) {
	res, ret, err := f(ctx)
	if err == nil || !ret {
		return
	}

	sctx, span := trace.StartSpan(ctx, traceTag)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	b.wg.Add(1)
	defer b.wg.Done()
	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	dur := b.initialDuration
	jdur := b.jittedInitialDuration

	dctx, cancel := context.WithDeadline(sctx, time.Now().Add(b.backoffTimeLimit))
	defer cancel()
	for cnt := 0; cnt < b.maxRetryCount; cnt++ {
		select {
		case <-dctx.Done():
			switch dctx.Err() {
			case context.DeadlineExceeded:
				return nil, errors.ErrBackoffTimeout(err)
			case context.Canceled:
				return nil, err
			default:
				return nil, errors.Wrap(err, dctx.Err().Error())
			}
		default:
			res, ret, err = func() (val interface{}, retryable bool, err error) {
				ssctx, span := trace.StartSpan(dctx, traceTag+"/"+strconv.Itoa(cnt+1))
				defer func() {
					if span != nil {
						span.End()
					}
				}()
				return f(ssctx)
			}()

			// e.g. name = v1.vald.Exists/10.0.0.0 ...etc
			if name := ctxkey.FromBackoffName(ctx); len(name) != 0 && b.metricsEnabled {
				b.mu.Lock()
				if v, ok := b.metrics[name]; !ok || v >= math.MaxInt64 {
					b.metrics[name] = 0
				} else {
					b.metrics[name] += 1
				}
				b.mu.Unlock()
			}

			if !ret {
				return res, err
			}
			if err == nil {
				return res, nil
			}
			if b.errLog {
				log.Error(err)
			}
			timer.Reset(time.Duration(jdur))
			select {
			case <-dctx.Done():
				switch dctx.Err() {
				case context.DeadlineExceeded:
					return nil, errors.ErrBackoffTimeout(err)
				case context.Canceled:
					return nil, err
				default:
					return nil, errors.Wrap(dctx.Err(), err.Error())
				}
			case <-timer.C:
				if dur >= b.durationLimit {
					dur = b.maxDuration
					jdur = b.maxDuration
				} else {
					dur *= b.backoffFactor
					jdur = b.addJitter(dur)
				}
			}
		}
	}
	return res, err
}

func (b *backoff) addJitter(dur float64) float64 {
	hd := math.Min(dur/10, b.jitterLimit)
	return dur + float64(rand.LimitedUint32(uint64(hd))) - hd
}

func (b *backoff) Metrics(_ context.Context) map[string]int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.metrics) == 0 {
		return nil
	}

	m := make(map[string]int64, len(b.metrics))
	for name, cnt := range b.metrics {
		m[name] = cnt
	}
	return m
}

// Close wait for the backoff process to complete.
func (b *backoff) Close() {
	b.wg.Wait()
}
