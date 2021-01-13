//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	client "github.com/vdaas/vald/internal/client/v1/client/compressor"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/worker"
)

type Registerer interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) (<-chan error, error)
	PostStop(ctx context.Context) error
	Register(ctx context.Context, vec *payload.Backup_Vector) error
	RegisterMulti(ctx context.Context, vecs *payload.Backup_Vectors) error
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
	client     client.Client
	vecs       map[string]*payload.Backup_Vector
	vecsMu     sync.Mutex
}

func NewRegisterer(opts ...RegistererOption) (Registerer, error) {
	r := new(registerer)
	for _, opt := range append(defaultRegistererOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	r.vecs = make(map[string]*payload.Backup_Vector, 100)

	return r, nil
}

func (r *registerer) PreStart(ctx context.Context) (err error) {
	r.worker, err = worker.New(append(r.workerOpts, worker.WithErrGroup(r.eg))...)
	return err
}

func (r *registerer) Start(ctx context.Context) (<-chan error, error) {
	return r.worker.Start(ctx)
}

func (r *registerer) PostStop(ctx context.Context) (err error) {
	log.Info("compressor registerer service poststop processing...")

	r.worker.Pause()

	ech, err := r.client.Start(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(r.client.Stop(ctx), err.Error())
			return
		}
		err = r.client.Stop(ctx)
	}()

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
		log.Errorf("compressor registerer service poststop failed: %v", err)
		return err
	}

	log.Info("compressor registerer service poststop completed")

	return nil
}

func (r *registerer) Register(ctx context.Context, vec *payload.Backup_Vector) error {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := r.dispatch(ctx, vec)
	if err != nil && span != nil {
		span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
	}

	return err
}

func (r *registerer) RegisterMulti(ctx context.Context, vecs *payload.Backup_Vectors) error {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.RegisterMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var err, errs error
	for _, vec := range vecs.GetVectors() {
		err = r.Register(ctx, vec)
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

func (r *registerer) dispatch(ctx context.Context, vec *payload.Backup_Vector) error {
	r.vecsMu.Lock()
	r.vecs[vec.GetUuid()] = vec
	r.vecsMu.Unlock()

	return r.worker.Dispatch(ctx, r.registerProcessFunc(vec))
}

func (r *registerer) registerProcessFunc(vec *payload.Backup_Vector) worker.JobFunc {
	return func(ctx context.Context) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/manager-compressor/service/Registerer.Register.DispatchedJob")
		defer func() {
			r.vecsMu.Lock()
			delete(r.vecs, vec.GetUuid())
			r.vecsMu.Unlock()

			if span != nil {
				span.End()
			}
		}()

		var vector []byte

		vector, err = r.compressor.Compress(ctx, vec.GetVector())
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}

			return err
		}

		err = r.backup.Register(
			ctx,
			&payload.Backup_Compressed_Vector{
				Uuid:   vec.GetUuid(),
				Vector: vector,
				Ips:    vec.GetIps(),
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

	r.vecsMu.Lock()
	defer r.vecsMu.Unlock()

	log.Debugf("compressor registerer queued vec-vector count: %d", len(r.vecs))

	for uuid, vec := range r.vecs {
		log.Debugf("forwarding uuid %s", uuid)

		err = r.client.Register(ctx, vec)
		if err != nil {
			log.Errorf("compressor registerer failed to backup uuid %s: %v", uuid, err)
			errs = errors.Wrap(errs, err.Error())
		}
	}

	return errs
}
