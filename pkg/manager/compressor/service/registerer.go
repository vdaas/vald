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

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/worker"
)

type Registerer interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	PreStop(ctx context.Context) error
	Register(ctx context.Context, meta *payload.Backup_MetaVector) error
	RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) error
}

type registerer struct {
	worker     worker.Worker
	workerOpts []worker.WorkerOption
	eg         errgroup.Group
	backup     Backup
	compressor Compressor
}

func NewRegisterer(opts ...RegistererOption) (Registerer, error) {
	r := new(registerer)
	for _, opt := range append(defaultRegistererOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return r, nil
}

func (r *registerer) PreStart(ctx context.Context) error {
	w, err := worker.NewWorker(append(r.workerOpts, worker.WithErrGroup(r.eg))...)
	if err != nil {
		return err
	}

	r.worker = w

	return nil
}

func (r *registerer) Start(ctx context.Context) <-chan error {
	return r.worker.Start(ctx)
}

func (r *registerer) PreStop(ctx context.Context) error {
	// TODO backup all index data here
	return nil
}

func (r *registerer) Register(ctx context.Context, meta *payload.Backup_MetaVector) error {
	if !r.worker.IsRunning() {
		return errors.ErrWorkerIsNotRunning(r.worker.Name())
	}

	return r.worker.Dispatch(ctx, func(ctx context.Context) error {
		vector, err := r.compressor.Compress(ctx, meta.GetVector())
		if err != nil {
			return err
		}

		return r.backup.Register(
			ctx,
			&payload.Backup_Compressed_MetaVector{
				Uuid:   meta.GetUuid(),
				Meta:   meta.GetMeta(),
				Vector: vector,
				Ips:    meta.GetIps(),
			},
		)
	})
}

func (r *registerer) RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) error {
	var err, errs error
	for _, meta := range metas.GetVectors() {
		err = r.Register(ctx, meta)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}
	return errs
}
