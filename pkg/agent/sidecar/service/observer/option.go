//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package observer provides storage observer
package observer

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type Option func(o *observer) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithBackupDuration("10m"),
	WithPostStopTimeout("2m"),
	WithWatch(true),
	WithTicker(true),
}

func WithBackupDuration(dur string) Option {
	return func(o *observer) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return nil
		}
		o.checkDuration = d
		return nil
	}
}

func WithPostStopTimeout(dur string) Option {
	return func(o *observer) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return nil
		}
		o.postStopTimeout = d
		return nil
	}
}

func WithWatch(enabled bool) Option {
	return func(o *observer) error {
		o.watchEnabled = enabled

		return nil
	}
}

func WithTicker(enabled bool) Option {
	return func(o *observer) error {
		o.tickerEnabled = enabled

		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(o *observer) error {
		if eg != nil {
			o.eg = eg
		}
		return nil
	}
}

func WithDir(dir string) Option {
	return func(o *observer) error {
		if dir == "" {
			return nil
		}

		o.dir = dir
		o.metadataPath = file.Join(dir, metadata.AgentMetadataFileName)

		return nil
	}
}

func WithBlobStorage(storage storage.Storage) Option {
	return func(o *observer) error {
		if storage != nil {
			o.storage = storage
		}
		return nil
	}
}

func WithHooks(hooks ...Hook) Option {
	return func(o *observer) error {
		if hooks == nil {
			return nil
		}

		if o.hooks != nil {
			o.hooks = append(o.hooks, hooks...)
		}

		o.hooks = hooks

		return nil
	}
}
