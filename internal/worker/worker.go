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
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
)

type WorkerJobFunc func(context.Context) error

type Worker interface {
	Start(ctx context.Context) <-chan error
	PreStop(ctx context.Context) error
	Wait()
	Pause()
	Resume()
	IsRunning() bool
	Name() string
	Len() int
	Dispatch(ctx context.Context, f WorkerJobFunc) error
}

type worker struct {
	name       string
	limitation int
	buffer     int
	running    atomic.Value
	eg         errgroup.Group
	wg         *sync.WaitGroup
	jobCh      chan WorkerJobFunc
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

	w.wg = new(sync.WaitGroup)
	w.wg.Add(1)

	w.jobCh = make(chan WorkerJobFunc, w.buffer)

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
				w.wg.Add(1)
				eg.Go(safety.RecoverFunc(func() (err error) {
					err = job(ctx)
					if err != nil {
						log.Debug(err)
						w.wg.Done()
						runtime.Gosched()
						ech <- err
						return err
					}
					w.wg.Done()
					return nil
				}))
			}
		}
	}))

	return ech
}

func (w *worker) PreStop(ctx context.Context) error {
	w.Pause()
	w.wg.Done()

	return nil
}

func (w *worker) Wait() {
	w.wg.Wait()
}

func (w *worker) Pause() {
	w.running.Store(false)
}

func (w *worker) Resume() {
	w.running.Store(true)
}

func (w *worker) IsRunning() bool {
	return w.running.Load().(bool)
}

func (w *worker) Name() string {
	return w.name
}

func (w *worker) Len() int {
	return len(w.jobCh)
}

func (w *worker) Dispatch(ctx context.Context, f WorkerJobFunc) error {
	ctx, span := trace.StartSpan(ctx, "vald/internal/worker/Worker.Dispatch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if f != nil {
		select {
		case w.jobCh <- f:
		default:
			err := errors.ErrWorkerChannelIsFull(w.name)
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return err
		}
	}
	return nil
}
