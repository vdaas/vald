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
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
)

type JobFunc func(context.Context) error

type Job struct {
	Name       string
	Fn         JobFunc
	Data       interface{}
	Retry      bool
	RetryCount int
}

type Worker interface {
	Start(ctx context.Context) (<-chan error, error)
	Wait()
	Pause()
	Resume()
	IsRunning() bool
	Name() string
	Len() uint64
	Jobs() []*Job
	Dispatch(ctx context.Context, f *Job) error
}

type worker struct {
	name       string
	limitation int
	running    atomic.Value
	eg         errgroup.Group
	wg         *sync.WaitGroup
	inCh       chan *Job
	q          []*Job
	qLen       uint64
	jobCh      chan *Job
}

func New(opts ...WorkerOption) (Worker, error) {
	w := new(worker)
	for _, opt := range append(defaultWorkerOpts, opts...) {
		if err := opt(w); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	w.running.Store(false)

	return w, nil
}

func (w *worker) Start(ctx context.Context) (<-chan error, error) {
	w.inCh = make(chan *Job, w.limitation)
	w.jobCh = make(chan *Job)

	w.startQueueLoop(ctx)

	return w.startJobLoop(ctx)
}

func (w *worker) startQueueLoop(ctx context.Context) {
	w.running.Store(true)

	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(w.inCh)
		for {
			select {
			case <-ctx.Done():
				w.running.Store(false)
				return ctx.Err()
			case j := <-w.inCh:
				w.q = append(w.q, j)
				atomic.AddUint64(&w.qLen, 1)
			default:
			}

			if len(w.q) > 0 {
				j := w.q[0]
				select {
				case w.jobCh <- j:
					w.q = w.q[1:]
					atomic.AddUint64(&w.qLen, ^uint64(0))
				default:
				}
			}
		}
	}))
}

func (w *worker) startJobLoop(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 1)

	eg, ctx := errgroup.New(ctx)
	eg.Limitation(w.limitation)

	w.wg = new(sync.WaitGroup)
	w.wg.Add(1)

	retryFn := func(j *Job) error {
		j.RetryCount++

		if !w.IsRunning() {
			log.Debugf("worker %s: append job [%s] to queue slice", w.Name(), j.Name)

			w.q = append(w.q, j)

			return nil
		}

		log.Debugf("worker %s: retrying job [%s]", w.Name(), j.Name)

		return w.Dispatch(ctx, j)
	}

	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(w.jobCh)
		for {
			select {
			case <-ctx.Done():
				w.wg.Done()
				if err = ctx.Err(); err != nil {
					return errors.Wrap(eg.Wait(), err.Error())
				}
				return eg.Wait()
			case job := <-w.jobCh:
				w.wg.Add(1)
				eg.Go(safety.RecoverFunc(func() (err error) {
					defer w.wg.Done()

					err = job.Fn(ctx)
					if err != nil {
						log.Debug(err)

						if job.Retry {
							err = errors.Wrap(retryFn(job), err.Error())
						}
					}

					return err
				}))
			}
		}
	}))

	return ech, nil
}

func (w *worker) Wait() {
	log.Debugf("worker %s: waiting for rest jobs to be completed...", w.Name())

	w.wg.Wait()

	log.Debugf("worker %s: all of the queued worker jobs completed.", w.Name())
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
	return atomic.LoadUint64(&w.qLen)
}

func (w *worker) Jobs() []*Job {
	return w.q
}

func (w *worker) Dispatch(ctx context.Context, f *Job) error {
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

	if f != nil && f.Fn != nil {
		w.inCh <- f
	}

	return nil
}
