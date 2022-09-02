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

// JobFunc represents the function of a job that works in the worker.
type JobFunc func(context.Context) error

// Worker represents the worker interface to execute jobs.
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
	qopts          []QueueOption
	requestedCount uint64
	completedCount uint64
}

// New initializes and return the worker, or return initialization error if occurred.
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
		append([]QueueOption{
			WithQueueBuffer(w.limitation),
			WithQueueErrGroup(w.eg),
		}, w.qopts...)...,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Start starts execute jobs in the worker queue.
// It returns the error channel that the job return, and the error if start failed.
func (w *worker) Start(ctx context.Context) (<-chan error, error) {
	if w.IsRunning() {
		return nil, errors.ErrWorkerIsAlreadyRunning(w.Name())
	}

	ech := make(chan error, 2)

	var wech, qech <-chan error
	var err error

	qech, err = w.queue.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	wech = w.startJobLoop(ctx)

	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		defer w.running.Store(false)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-qech:
			case err = <-wech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))

	w.running.Store(true)

	return ech, nil
}

func (w *worker) startJobLoop(ctx context.Context) <-chan error {
	ech := make(chan error, w.limitation)

	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		eg, ctx := errgroup.New(ctx)
		eg.Limitation(w.limitation)

		limitation := make(chan struct{}, w.limitation)
		defer close(limitation)

		for {
			select {
			case <-ctx.Done():
				if err = ctx.Err(); err != nil {
					return errors.Wrap(eg.Wait(), err.Error())
				}
				return eg.Wait()
			case limitation <- struct{}{}:
			}

			job, err := w.queue.Pop(ctx)
			if err != nil {
				ech <- err
				select {
				case <-limitation:
				case <-ctx.Done():
				}
				continue
			}

			if job != nil {
				eg.Go(safety.RecoverFunc(func() (err error) {
					defer atomic.AddUint64(&w.completedCount, 1)
					if err = job(ctx); err != nil {
						log.Debugf("an error occurred while executing a job: %s", err)
						ech <- err
					}
					select {
					case <-limitation:
					case <-ctx.Done():
					}
					return err
				}))
			} else {
				select {
				case <-limitation:
				case <-ctx.Done():
				}
			}
		}
	}))

	return ech
}

// Pause stops allowing new job to be dispatched to the worker.
func (w *worker) Pause() {
	w.running.Store(false)
}

// Resume resumes to allow new jobs to be dispatched to the worker.
func (w *worker) Resume() {
	w.running.Store(true)
}

// IsRunning returns if the worker is running or not.
func (w *worker) IsRunning() bool {
	return w.running.Load().(bool)
}

// Name returns the worker name.
func (w *worker) Name() string {
	return w.name
}

// Len returns the length of the worker queue.
func (w *worker) Len() uint64 {
	return w.queue.Len()
}

// TotalRequested returns the number of jobs that dispatched to the worker.
func (w *worker) TotalRequested() uint64 {
	return atomic.LoadUint64(&w.requestedCount)
}

// TotalCompleted returns the number of completed job.
func (w *worker) TotalCompleted() uint64 {
	return atomic.LoadUint64(&w.completedCount)
}

// Dispatch dispatches the job to the worker and waiting for the worker to process it.
// The job error is pushed to the error channel that Start() return.
// This function will return an error if the job cannot be dispatch to the worker queue, or the worker is not running.
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
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
