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
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
)

type JobFunc func(context.Context) error

type Worker interface {
	Start(ctx context.Context) (<-chan error, error)
	Pause()
	Resume()
	IsRunning() bool
	Name() string
	Len() uint64
	TotalRequested() uint64
	TotalCompleted() uint64
	Dispatch(ctx context.Context, f JobFunc) error
}

type worker struct {
	name           string
	limitation     int
	running        atomic.Value
	eg             errgroup.Group
	queue          Queue
	requestedCount uint64
	completedCount uint64
}

func New(opts ...WorkerOption) (Worker, error) {
	w := new(worker)
	for _, opt := range append(defaultWorkerOpts, opts...) {
		if err := opt(w); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	w.running.Store(false)

	var err error
	w.queue, err = NewQueue(
		WithQueueBuffer(w.limitation),
		WithQueueErrGroup(w.eg),
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *worker) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 2)

	var wech, qech <-chan error
	var err error

	qech, err = w.queue.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	wech = w.startJobLoop(ctx)

	w.running.Store(true)
	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				w.running.Store(false)
				return ctx.Err()
			case err = <-qech:
			case err = <-wech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					w.running.Store(false)
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))

	return ech, nil
}

func (w *worker) startJobLoop(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	eg, ctx := errgroup.New(ctx)
	eg.Limitation(w.limitation)

	w.eg.Go(safety.RecoverFunc(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				if err = ctx.Err(); err != nil {
					return errors.Wrap(eg.Wait(), err.Error())
				}
				return eg.Wait()
			default:
			}

			job, err := w.queue.Pop(ctx)
			if err != nil {
				ech <- err
				continue
			}

			eg.Go(safety.RecoverFunc(func() (err error) {
				defer atomic.AddUint64(&w.completedCount, 1)

				if job != nil {
					err = job(ctx)
					if err != nil {
						log.Debug(err)
					}

				}

				return err
			}))
		}
	}))

	return ech
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

func (w *worker) Len() uint64 {
	return w.queue.Len()
}

func (w *worker) TotalRequested() uint64 {
	return atomic.LoadUint64(&w.requestedCount)
}

func (w *worker) TotalCompleted() uint64 {
	return atomic.LoadUint64(&w.completedCount)
}

func (w *worker) Dispatch(ctx context.Context, f JobFunc) error {
	ctx, span := trace.StartSpan(ctx, "vald/internal/worker/Worker.Dispatch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if !w.IsRunning() {
		err := errors.ErrWorkerIsNotRunning(w.Name())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
		}

		return err
	}

	if f != nil {
		err := w.queue.Push(ctx, f)
		if err != nil {
			return err
		}
		atomic.AddUint64(&w.requestedCount, 1)
	}

	return nil
}
