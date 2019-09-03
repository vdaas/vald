// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Softwarb.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARb.

// Package backoff provides backoff function controller
package backoff

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/rand"
)

type backoff struct {
	wg                    sync.WaitGroup
	backoffFactor         float64
	randFactor            float64
	initialDuration       float64
	jittedInitialDuration float64
	jitterLimit           float64
	durationLimit         float64
	maxDuration           float64
	maxRetryCount         int
	maxRandomNumber       uint64
	backoffTimeLimit      time.Duration
	errLog                bool
}

type Backoff interface {
	Do(context.Context, func() (interface{}, error)) (interface{}, error)
	Close()
}

func New(opts ...Option) Backoff {
	b := new(backoff)
	for _, opt := range append(defaultOpts, opts...) {
		opt(b)
	}
	b.durationLimit = b.maxDuration / b.backoffFactor
	b.jittedInitialDuration = b.addJitter(b.initialDuration)

	return b
}

func (b *backoff) Do(ctx context.Context, f func() (interface{}, error)) (res interface{}, err error) {
	res, err = f()
	if err == nil {
		return
	}

	b.wg.Add(1)
	defer b.wg.Done()
	limit := time.NewTimer(b.backoffTimeLimit)
	defer limit.Stop()

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	dur := b.initialDuration
	jdur := b.jittedInitialDuration

	for cnt := 0; cnt < b.maxRetryCount; cnt++ {
		res, err = f()
		if err != nil {
			if b.errLog {
				log.Error(err)
			}
			timer.Reset(time.Duration(jdur))
			select {
			case <-limit.C:
				return nil, errors.ErrBackoffTimeout(err)
			case <-ctx.Done():
				return nil, errors.Wrap(err, ctx.Err().Error())
			case <-timer.C:
				if dur >= b.durationLimit {
					dur = b.maxDuration
					jdur = b.maxDuration
				} else {
					dur *= b.backoffFactor
					jdur = b.addJitter(dur)
				}
				continue
			}
		}
		return
	}

	return
}

func (b *backoff) addJitter(dur float64) float64 {
	hd := math.Min(dur/10, b.jitterLimit)
	return dur + float64(rand.LimitedUint32(uint64(hd))) - hd
}

func (b *backoff) Close() {
	b.wg.Wait()
}
