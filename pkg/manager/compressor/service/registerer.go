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
	"sync"

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
	TotalRequested() uint64
	TotalCompleted() uint64
}

type registerer struct {
	worker     worker.Worker
	workerOpts []worker.WorkerOption
	eg         errgroup.Group
	backup     Backup
	compressor Compressor
	addr       string
	client     grpc.Client
	metas      map[string]*payload.Backup_MetaVector
	metasMux   sync.Mutex
}

func NewRegisterer(opts ...RegistererOption) (Registerer, error) {
	r := new(registerer)
	for _, opt := range append(defaultRegistererOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	r.metas = make(map[string]*payload.Backup_MetaVector, 0)

	return r, nil
}

func (r *registerer) PreStart(ctx context.Context) (err error) {
	r.worker, err = worker.New(append(r.workerOpts, worker.WithErrGroup(r.eg))...)
	return err
}

func (r *registerer) Start(ctx context.Context) (<-chan error, error) {
	return r.worker.Start(ctx)
}

func (r *registerer) PreStop(ctx context.Context) error {
	log.Info("compressor registerer service prestop processing...")

	r.worker.Pause()

	ech, err := r.client.StartConnectionMonitor(ctx)
	if err != nil {
		return err
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, cctx := errgroup.New(cctx)
	eg.Go(safety.RecoverFunc(func() error {
		for {
			select {
			case <-cctx.Done():
				return cctx.Err()
			case err := <-ech:
				if err != nil {
					log.Error(err)
				}
			}
		}
	}))

	err = r.forwardMetas(ctx)
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

	err := r.dispatch(ctx, meta)
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

func (r *registerer) TotalRequested() uint64 {
	return r.worker.TotalRequested()
}

func (r *registerer) TotalCompleted() uint64 {
	return r.worker.TotalCompleted()
}

func (r *registerer) dispatch(ctx context.Context, meta *payload.Backup_MetaVector) error {
	r.metasMux.Lock()
	r.metas[meta.GetUuid()] = meta
	r.metasMux.Unlock()

	return r.worker.Dispatch(ctx, r.registerProcessFunc(meta))
}

func (r *registerer) registerProcessFunc(meta *payload.Backup_MetaVector) worker.JobFunc {
	return func(ctx context.Context) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.Register.DispatchedJob")
		defer func() {
			r.metasMux.Lock()
			delete(r.metas, meta.GetUuid())
			r.metasMux.Unlock()

			if span != nil {
				span.End()
			}
		}()

		var vector []byte

		vector, err = r.compressor.Compress(ctx, meta.GetVector())
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

	r.metasMux.Lock()

	log.Debugf("compressor registerer queued meta-vector count: %d", len(r.metas))

	for uuid, meta := range r.metas {
		log.Debugf("forwarding uuid %s", uuid)

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
			log.Errorf("compressor registerer failed to backup uuid %s: %v", uuid, err)
			errs = errors.Wrap(errs, err.Error())
		}
	}

	r.metasMux.Unlock()

	return errs
}
