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
	Len() uint64
}

type registerer struct {
	worker     worker.Worker
	workerOpts []worker.WorkerOption
	eg         errgroup.Group
	backup     Backup
	compressor Compressor
	addr       string
	client     grpc.Client
	running    atomic.Value
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

func (r *registerer) PreStart(ctx context.Context) (err error) {
	r.worker, err = worker.New(append(r.workerOpts, worker.WithErrGroup(r.eg))...)
	return err
}

func (r *registerer) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)

	var wech, cech <-chan error

	wech, err := r.worker.Start(ctx)
	if err != nil {
		return nil, err
	}

	cech = r.startConnectionMonitor(ctx)

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
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

	r.worker.Pause()
	r.worker.Wait()

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

	err := r.worker.Dispatch(ctx, r.registerProcessFunc(meta))
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

func (r *registerer) Len() uint64 {
	return r.worker.Len()
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

func (r *registerer) registerProcessFunc(meta *payload.Backup_MetaVector) *worker.Job {
	f := func(ctx context.Context) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.Register.DispatchedJob")
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		var vector []byte

		vector, err = r.compressor.Compress(ctx, meta.GetVector())
		if err != nil {
			log.Debugf("re-enqueueing uuid %s", meta.GetUuid())

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
		if err != nil {
			log.Debugf("re-enqueueing uuid %s", meta.GetUuid())

			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
		}

		return err
	}

	return &worker.Job{
		Fn:   f,
		Data: meta,
	}
}

func (r *registerer) forwardMetas(ctx context.Context) (errs error) {
	var err error

	log.Debugf("compressor registerer queued meta-vector count: %d", r.Len())

	// TODO read all metas from worker queue
	var metas []*payload.Backup_MetaVector

	for _, meta := range metas {
		log.Debugf("forwarding uuid %s", meta.GetUuid())

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
					meta,
					copts...,
				)
			},
		)
		if err != nil {
			log.Errorf(
				"compressor registerer failed to backup uuid %s: %v",
				meta.GetUuid(),
				err,
			)
			errs = errors.Wrap(errs, err.Error())
		}
	}

	return errs
}
