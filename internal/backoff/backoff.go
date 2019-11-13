//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"sync"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
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
