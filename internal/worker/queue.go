//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package worker provides worker processes
package worker

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Queue represents the interface of queue.
type Queue interface {
	Start(ctx context.Context) (<-chan error, error)
	Push(ctx context.Context, job JobFunc) error
	Pop(ctx context.Context) (JobFunc, error)
	Len() uint64
}

type queue struct {
	buffer  int
	eg      errgroup.Group
	qcdur   time.Duration // queue check duration
	inCh    chan JobFunc
	outCh   chan JobFunc
	qLen    atomic.Value
	running atomic.Value
}

// NewQueue returns Queue if no error is occurred.
func NewQueue(opts ...QueueOption) (Queue, error) {
	q := new(queue)
	for _, opt := range append(defaultQueueOpts, opts...) {
		if err := opt(q); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	q.qLen.Store(uint64(0))
	q.running.Store(false)

	q.inCh = make(chan JobFunc, q.buffer)
	q.outCh = make(chan JobFunc, 1)

	return q, nil
}

// Start starts execute queueing if queue is not runnnig.
// If queue is already reunning, it returns error.
// It returns the error channel that the queueing job return.
func (q *queue) Start(ctx context.Context) (<-chan error, error) {
	if q.isRunning() {
		return nil, errors.ErrQueueIsAlreadyRunning
	}

	ech := make(chan error, 1)
	q.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		defer close(q.outCh)
		defer close(q.inCh)
		defer q.running.Store(false)
		s := make([]JobFunc, 0, q.buffer)
		tick := time.NewTicker(q.qcdur)
		defer tick.Stop()
		for {
			if len(s) > 0 && len(q.outCh) == 0 && cap(q.inCh) != len(q.inCh) {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case q.outCh <- s[0]:
					s = s[1:]
				case <-tick.C:
				}
			} else {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case in := <-q.inCh:
					s = append(s, in)
				case <-tick.C:
				}
			}
			q.qLen.Store(uint64(len(s)))
		}
	}))

	q.running.Store(true)

	return ech, nil
}

// isRunning returns true when queue is already running or false when queue is not running.
func (q *queue) isRunning() bool {
	return q.running.Load().(bool)
}

// Push sends JobFunc to channel, which will be used for stock JobFunc, when JobFunc is not nil and queue is running.
// If JobFunc is nil or queue is not running, Push returns error.
func (q *queue) Push(ctx context.Context, job JobFunc) error {
	if job == nil {
		return errors.ErrJobFuncIsNil
	}

	if !q.isRunning() {
		return errors.ErrQueueIsNotRunning
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.inCh <- job:
		return nil
	}
}

// Pop returns (JobFunc, nil) if the channnel, which will be used for queuing job, contains JobFunc.
// It returns (nil ,error) if it failed to pop from the job queue.
func (q *queue) Pop(ctx context.Context) (JobFunc, error) {
	tryCnt := int(q.Len()) + 1 // include the first try

	for i := 0; i < tryCnt; i++ {
		if !q.isRunning() {
			return nil, errors.ErrQueueIsNotRunning
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case job := <-q.outCh:
			if job != nil {
				return job, nil
			}
		}
	}
	return nil, errors.ErrJobFuncNotFound
}

// Len returns the length of queue.
func (q *queue) Len() uint64 {
	return q.qLen.Load().(uint64)
}
