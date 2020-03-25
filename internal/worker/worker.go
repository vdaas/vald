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
	"runtime"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type Worker interface {
	Start(ctx context.Context) <-chan error
	IsRunning() bool
	Name() string
	Dispatch(ctx context.Context, f func() error) error
}

type worker struct {
	name       string
	limitation int
	buffer     int
	running    atomic.Value
	eg         errgroup.Group
	jobCh      chan func() error
}

func NewWorker(opts ...WorkerOption) (Worker, error) {
	w := new(worker)
	for _, opt := range append(defaultWorkerOpts, opts...) {
		if err := opt(w); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	w.running.Store(false)

	return w, nil
}

func (w *worker) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	eg, ctx := errgroup.New(ctx)
	eg.Limitation(w.limitation)

	w.jobCh = make(chan func() error, w.buffer)

	w.running.Store(true)
	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(w.jobCh)
		for {
			select {
			case <-ctx.Done():
				w.running.Store(false)
				if err = ctx.Err(); err != nil {
					return errors.Wrap(eg.Wait(), err.Error())
				}
				return eg.Wait()
			case job := <-w.jobCh:
				eg.Go(safety.RecoverFunc(func() (err error) {
					err = job()
					if err != nil {
						log.Debug(err)
						runtime.Gosched()
						ech <- err
					}
					return nil
				}))
			}
		}
	}))

	return ech
}

func (w *worker) IsRunning() bool {
	return w.running.Load().(bool)
}

func (w *worker) Name() string {
	return w.name
}

func (w *worker) Dispatch(ctx context.Context, f func() error) error {
	if f != nil {
		select {
		case w.jobCh <- f:
		default:
			return errors.ErrWorkerChannelIsFull(w.name)
		}
	}
	return nil
}
