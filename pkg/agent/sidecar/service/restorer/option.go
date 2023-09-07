//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package restorer provides restorer service
package restorer

import (
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type Option func(r *restorer) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithBackoff(false),
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(r *restorer) error {
		if eg != nil {
			r.eg = eg
		}
		return nil
	}
}

func WithDir(dir string) Option {
	return func(r *restorer) error {
		if dir == "" {
			return nil
		}

		r.dir = dir

		return nil
	}
}

func WithBlobStorage(storage storage.Storage) Option {
	return func(r *restorer) error {
		if storage != nil {
			r.storage = storage
		}
		return nil
	}
}

func WithBackoff(enabled bool) Option {
	return func(r *restorer) error {
		r.backoffEnabled = enabled
		return nil
	}
}

func WithBackoffOpts(opts ...backoff.Option) Option {
	return func(r *restorer) error {
		if r.backoffOpts == nil {
			r.backoffOpts = opts
		}

		r.backoffOpts = append(r.backoffOpts, opts...)

		return nil
	}
}
