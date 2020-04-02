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

package service

import (
	"context"
	"reflect"
	"runtime"
	"sync/atomic"

	gcomp "github.com/vdaas/vald/apis/grpc/manager/compressor"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/worker"
)

type Registerer interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) (<-chan error, error)
	PreStop(ctx context.Context) error
	Register(ctx context.Context, meta *payload.Backup_MetaVector) error
	RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) error
	Len() int
	WorkerLen() int
}

type registerer struct {
	worker     worker.Worker
	workerOpts []worker.WorkerOption
	eg         errgroup.Group
	backup     Backup
	compressor Compressor
	addr       string
	client     grpc.Client
	ch         chan *payload.Backup_MetaVector
	buffer     int
	running    atomic.Value
	backoff    backoff.Backoff
}

func NewRegisterer(opts ...RegistererOption) (Registerer, error) {
	r := new(registerer)
	for _, opt := range append(defaultRegistererOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	r.running.Store(false)

	return r, nil
}

func (r *registerer) PreStart(ctx context.Context) error {
	r.ch = make(chan *payload.Backup_MetaVector, r.buffer)

	w, err := worker.NewWorker(append(r.workerOpts, worker.WithErrGroup(r.eg))...)
	if err != nil {
		return err
	}

	r.worker = w

	return nil
}

func (r *registerer) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)

	var wech, rech, cech <-chan error

	rech = r.startDispatcherLoop(ctx)

	wech = r.worker.Start(ctx)

	cech = r.startConnectionMonitor(ctx)

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-rech:
			case err = <-wech:
			case err = <-cech:
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

	r.running.Store(true)

	return ech, nil
}

func (r *registerer) PreStop(ctx context.Context) error {
	log.Info("compressor registerer service prestop processing...")

	r.running.Store(false)
	r.worker.Pause(ctx)

	err := r.forwardMetas(ctx)
	if err != nil {
		log.Errorf("compressor registerer service prestop failed: %v", err)
		return err
	}

	log.Info("compressor registerer service prestop completed")

	return nil
}

func (r *registerer) Register(ctx context.Context, meta *payload.Backup_MetaVector) error {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if !r.worker.IsRunning() {
		err := errors.ErrWorkerIsNotRunning(r.worker.Name())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
		}

		return err
	}

	err := r.sendToChWithRunningCheck(ctx, meta)
	if err != nil && span != nil {
		span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
	}

	return err
}

func (r *registerer) RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) error {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.RegisterMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var err, errs error
	for _, meta := range metas.GetVectors() {
		err = r.Register(ctx, meta)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}

	if errs != nil && span != nil {
		span.SetStatus(trace.StatusCodeUnavailable(errs.Error()))
	}

	return errs
}

func (r *registerer) Len() int {
	return len(r.ch)
}

func (r *registerer) WorkerLen() int {
	return r.worker.Len()
}

func (r *registerer) startDispatcherLoop(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(r.ch)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case meta := <-r.ch:
				err = r.worker.Dispatch(ctx, r.registerJob(meta))
				if err != nil {
					err = errors.Wrap(r.sendToCh(ctx, meta), err.Error())
					runtime.Gosched()
					ech <- err
				}
			}
		}
	}))

	return ech
}

func (r *registerer) startConnectionMonitor(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		ch, err := r.client.StartConnectionMonitor(ctx)
		if err != nil {
			ech <- err
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-ch:
				if err != nil {
					runtime.Gosched()
					ech <- err
				}
			}
		}
	}))

	return ech
}

func (r *registerer) sendToChWithRunningCheck(ctx context.Context, meta *payload.Backup_MetaVector) error {
	if !r.running.Load().(bool) {
		return errors.ErrCompressorRegistererIsNotRunning()
	}

	return r.sendToCh(ctx, meta)
}

func (r *registerer) sendToCh(ctx context.Context, meta *payload.Backup_MetaVector) error {
	f := func() error {
		select {
		case r.ch <- meta:
		default:
			return errors.ErrCompressorRegistererChannelIsFull()
		}
		return nil
	}

	if r.backoff != nil {
		_, err := r.backoff.Do(ctx, func() (interface{}, error) {
			return nil, f()
		})
		return err
	}

	return f()
}

func (r *registerer) registerJob(meta *payload.Backup_MetaVector) worker.WorkerJobFunc {
	return func(ctx context.Context) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.Register.DispatchedJob")
		defer func() {
			if err != nil {
				err = errors.Wrap(r.sendToCh(ctx, meta), err.Error())
			}
			if span != nil {
				span.End()
			}
		}()

		vector, err := r.compressor.Compress(ctx, meta.GetVector())
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return err
		}

		err = r.backup.Register(
			ctx,
			&payload.Backup_Compressed_MetaVector{
				Uuid:   meta.GetUuid(),
				Meta:   meta.GetMeta(),
				Vector: vector,
				Ips:    meta.GetIps(),
			},
		)
		if err != nil && span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}

		return err
	}
}

func (r *registerer) forwardMetas(ctx context.Context) (errs error) {
	var err error

	func() {
		for {
			select {
			case v, ok := <-r.ch:
				if !ok {
					return
				}
				log.Debugf("forwarding uuid %s", v.GetUuid())
				_, err = r.client.Do(
					ctx,
					r.addr,
					func(
						ctx context.Context,
						conn *grpc.ClientConn,
						copts ...grpc.CallOption,
					) (interface{}, error) {
						return gcomp.NewBackupClient(conn).Register(
							ctx,
							v,
							copts...,
						)
					},
				)
				if err != nil {
					log.Errorf(
						"compressor registerer failed to backup uuid %s: %v",
						v.GetUuid(),
						err,
					)
					errs = errors.Wrap(errs, err.Error())
				}
			default:
				return
			}
		}
	}()

	return errs
}
