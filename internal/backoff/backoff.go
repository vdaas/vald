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

package backoff

import (
	"context"
	"maps"
	"math"
	stdrand "math/rand/v2"
	"strconv"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// NOTE: This variable is for observability package.
//
//	This will be fixed when refactoring the observability package.
var (
	mu      sync.RWMutex
	metrics map[string]int64 = make(map[string]int64)
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
}

// Backoff represents an interface to handle backoff operation.
type Backoff interface {
	Do(context.Context, func(ctx context.Context) (any, bool, error)) (any, error)
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

	return b
}

func resolveBackoffName(ctx context.Context) string {
	name := FromBackoffName(ctx)
	if name == "" || name == "Value" || strings.TrimSpace(name) == "" {
		st := info.Get().StackTrace
		for _, f := range st {
			if f.FuncName != "" && !strings.Contains(f.FuncName, "/internal/backoff.") {
				return f.FuncName
			}
		}
		if len(st) > 0 && st[len(st)-1].FuncName != "" {
			return st[len(st)-1].FuncName
		}
		return "unknown backoff"
	}
	return name
}

func backoffContextErr(dctx context.Context, name string, err error) error {
	switch dctx.Err() {
	case context.DeadlineExceeded:
		log.Debugf("[backoff]\tfor: "+name+",\tDeadline Exceeded\terror: %v", err.Error())
		return errors.ErrBackoffTimeout(err)
	case context.Canceled:
		log.Debugf("[backoff]\tfor: "+name+",\tCanceled\terror: %v", err.Error())
		return err
	default:
		return errors.Join(dctx.Err(), err)
	}
}

func (b *backoff) recordRetryFailure(ctx context.Context, name string, err error) {
	if b.errLog && err != nil && err.Error() != "" {
		if IsQuietErrorLog(ctx) {
			log.Debugd("[backoff]\tfor: "+name+",\terror: "+err.Error(), info.Get())
		} else {
			log.Errord("[backoff]\tfor: "+name+",\terror: "+err.Error(), info.Get())
		}
	}
	// e.g. name = vald.v1.Exists/ip ...etc
	if name != "" {
		mu.Lock()
		metrics[name] += 1
		mu.Unlock()
	}
}

func (b *backoff) nextDuration(dur float64) (nextDur, nextJdur float64) {
	if dur >= b.durationLimit {
		return b.maxDuration, b.maxDuration
	}
	nextDur = dur * b.backoffFactor
	return nextDur, b.addJitter(nextDur)
}

func (b *backoff) retry(
	ctx context.Context, firstErr error, f func(ctx context.Context) (retryable bool, err error),
) (cleared bool, err error) {
	if firstErr == nil || f == nil || b.maxRetryCount <= 0 {
		return false, firstErr
	}
	ctx, running := isRunning(ctx)
	if running {
		return false, firstErr
	}

	sctx, span := trace.StartSpan(ctx, traceTag)
	defer trace.End(span)
	b.wg.Add(1)
	defer b.wg.Done()
	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	dur := b.initialDuration
	jdur := b.jittedInitialDuration
	name := resolveBackoffName(ctx)
	err = firstErr

	dctx, cancel := context.WithDeadline(sctx, time.Now().Add(b.backoffTimeLimit))
	defer cancel()
	for cnt := 0; cnt < b.maxRetryCount; cnt++ {
		select {
		case <-dctx.Done():
			return true, backoffContextErr(dctx, name, err)
		default:
			var ret bool
			if trace.Enabled() {
				ret, err = func() (retryable bool, err error) {
					ssctx, span := trace.StartSpan(dctx, traceTag+"/"+strconv.Itoa(cnt+1))
					defer trace.End(span)
					return f(ssctx)
				}()
			} else {
				ret, err = f(dctx)
			}
			if !ret {
				return false, err
			}
			if err == nil {
				return false, nil
			}
			if errors.Is(err, context.Canceled) ||
				errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, errors.ErrBackoffTimeout(nil)) {
				return true, err
			}
			b.recordRetryFailure(ctx, name, err)

			// No attempt is left (continue ends the loop), so waiting here would only delay (or, if the deadline hits, mask) the final result.
			if cnt == b.maxRetryCount-1 {
				continue
			}
			timer.Reset(time.Duration(jdur))
			select {
			case <-dctx.Done():
				return true, backoffContextErr(dctx, name, err)
			case <-timer.C:
				dur, jdur = b.nextDuration(dur)
			}
		}
	}
	select {
	case <-dctx.Done():
		return true, backoffContextErr(dctx, name, err)
	default:
		return false, err
	}
}

// RetryVoid executes the retry loop after attempt 0 already ran inline and
// returned (retryable=true, firstErr!=nil).
func (b *backoff) RetryVoid(
	ctx context.Context, firstErr error, f func(ctx context.Context) (retryable bool, err error),
) error {
	_, err := b.retry(ctx, firstErr, f)
	return err
}

// Do tries to backoff using the input function and returns the response and error.
func (b *backoff) Do(
	ctx context.Context, f func(ctx context.Context) (val any, retryable bool, err error),
) (res any, err error) {
	if f == nil {
		return res, err
	}
	var ret bool
	res, ret, err = f(ctx)
	if err == nil || !ret {
		return res, err
	}
	return b.doRetry(ctx, res, err, f)
}

//go:noinline

func (b *backoff) doRetry(
	ctx context.Context,
	res any,
	firstErr error,
	f func(ctx context.Context) (val any, retryable bool, err error),
) (any, error) {
	cleared, rerr := b.retry(ctx, firstErr, func(ictx context.Context) (bool, error) {
		var r bool
		var e error
		res, r, e = f(ictx)
		return r, e
	})
	if cleared {
		return nil, rerr
	}
	return res, rerr
}

func (b *backoff) addJitter(dur float64) float64 {
	hd := math.Min(dur/10, b.jitterLimit)
	if hd < 1 {
		return dur
	}
	return dur + stdrand.Float64()*hd*2 - hd
}

// Close wait for the backoff process to complete.
func (b *backoff) Close() {
	b.wg.Wait()
}

func Metrics(context.Context) map[string]int64 {
	mu.RLock()
	defer mu.RUnlock()

	if len(metrics) == 0 {
		return nil
	}

	return maps.Clone(metrics)
}
