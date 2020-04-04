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
	Start(ctx context.Context)
	InCh() chan<- JobFunc
	OutCh() <-chan JobFunc
	Len() uint64
}

type queue struct {
	buffer int
	eg     errgroup.Group
	inCh   chan JobFunc
	outCh  chan JobFunc
	qLen   atomic.Value
}

func NewQueue(opts ...QueueOption) (Queue, error) {
	q := new(queue)
	for _, opt := range append(defaultQueueOpts, opts...) {
		if err := opt(q); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	q.qLen.Store(uint64(0))

	return q, nil
}

func (q *queue) Start(ctx context.Context) {
	q.inCh = make(chan JobFunc, q.buffer)
	q.outCh = make(chan JobFunc)

	s := make([]JobFunc, 0)

	q.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(q.inCh)
		defer close(q.outCh)

		for {
			if len(s) > 0 {
				j := s[0]
				select {
				case <-ctx.Done():
					return ctx.Err()
				case q.outCh <- j:
					s = s[1:]
					q.qLen.Store(uint64(len(s)))
				case j := <-q.inCh:
					s = append(s, j)
					q.qLen.Store(uint64(len(s)))
				}
			} else {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case j := <-q.inCh:
					s = append(s, j)
					q.qLen.Store(uint64(len(s)))
				}
			}
		}
	}))
}

func (q *queue) InCh() chan<- JobFunc {
	return q.inCh
}

func (q *queue) OutCh() <-chan JobFunc {
	return q.outCh
}

func (q *queue) Len() uint64 {
	return q.qLen.Load().(uint64)
}
