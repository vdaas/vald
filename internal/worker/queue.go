//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package worker provides worker processes
package worker

import (
	"context"
	"reflect"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
)

type Queue interface {
	Start(ctx context.Context) (<-chan error, error)
	Push(ctx context.Context, job JobFunc) error
	Pop(ctx context.Context) (JobFunc, error)
	Len() uint64
}

type queue struct {
	buffer  int
	eg      errgroup.Group
	inCh    chan JobFunc
	outCh   chan JobFunc
	qLen    atomic.Value
	running atomic.Value
}

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
	q.outCh = make(chan JobFunc)

	return q, nil
}

func (q *queue) Start(ctx context.Context) (<-chan error, error) {
	if q.isRunning() {
		return nil, errors.ErrQueueIsAlreadyRunning()
	}

	ech := make(chan error, 1)

	s := make([]JobFunc, 0, q.buffer)

	inFn := func(j JobFunc) {
		s = append(s, j)
		q.qLen.Store(uint64(len(s)))
	}
	outFn := func() {
		s = s[1:]
		q.qLen.Store(uint64(len(s)))
	}

	q.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(q.outCh)
		defer close(q.inCh)
		defer q.running.Store(false)

		for {
			if len(s) > 0 {
				j := s[0]
				select {
				case <-ctx.Done():
					return ctx.Err()
				case q.outCh <- j:
					outFn()
				case j := <-q.inCh:
					inFn(j)
				}
			} else {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case j := <-q.inCh:
					inFn(j)
				}
			}
		}
	}))

	q.running.Store(true)

	return ech, nil
}

func (q *queue) isRunning() bool {
	return q.running.Load().(bool)
}

func (q *queue) Push(ctx context.Context, job JobFunc) error {
	if job == nil {
		return errors.ErrJobFuncIsNil()
	}

	if !q.isRunning() {
		return errors.ErrQueueIsNotRunning()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.inCh <- job:
		return nil
	}
}

func (q *queue) Pop(ctx context.Context) (JobFunc, error) {
	if !q.isRunning() {
		return nil, errors.ErrQueueIsNotRunning()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case job := <-q.outCh:
		if job == nil {
			return nil, errors.ErrJobFuncIsNil()
		}

		return job, nil
	}
}

func (q *queue) Len() uint64 {
	return q.qLen.Load().(uint64)
}
