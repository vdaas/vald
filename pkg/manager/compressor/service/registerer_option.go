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

// Package service
package service

import (
	client "github.com/vdaas/vald/internal/client/v1/client/compressor"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/worker"
)

type RegistererOption func(r *registerer) error

var defaultRegistererOpts = []RegistererOption{
	WithRegistererWorker(),
	WithRegistererErrGroup(errgroup.Get()),
}

func WithRegistererWorker(opts ...worker.WorkerOption) RegistererOption {
	return func(r *registerer) error {
		r.workerOpts = opts
		return nil
	}
}

func WithRegistererErrGroup(eg errgroup.Group) RegistererOption {
	return func(r *registerer) error {
		if eg != nil {
			r.eg = eg
		}
		return nil
	}
}

func WithRegistererBackup(b Backup) RegistererOption {
	return func(r *registerer) error {
		if b != nil {
			r.backup = b
		}
		return nil
	}
}

func WithRegistererCompressor(c Compressor) RegistererOption {
	return func(r *registerer) error {
		if c != nil {
			r.compressor = c
		}
		return nil
	}
}

func WithRegistererClient(c client.Client) RegistererOption {
	return func(r *registerer) error {
		if c != nil {
			r.client = c
		}
		return nil
	}
}
