//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package service

import (
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(_ *export) error

var defaultOpts = []Option{
	WithStreamListConcurrency(200),
	WithKVSSyncInterval("5s"),
	WithKVSCompactionInterval("5s"),
	WithIndexPath("/var/export/index"),
	WithErrGroup(errgroup.Get()),
}

// WithStreamListConcurrency returns Option that sets streamListConcurrency.
func WithStreamListConcurrency(num int) Option {
	return func(e *export) error {
		if num <= 0 {
			return errors.NewErrInvalidOption("streamListConcurrency", num)
		}
		e.streamListConcurrency = num
		return nil
	}
}

// WithKVSSyncInterval returns Option that sets interval for background file sync.
func WithKVSSyncInterval(dur string) Option {
	return func(e *export) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		e.backgroundSyncInterval = d
		return nil
	}
}

// WithKVSCompactionInterval returns Option that sets interval for background file compaction.
func WithKVSCompactionInterval(dur string) Option {
	return func(e *export) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		e.backgroundCompactionInterval = d
		return nil
	}
}

// WithIndexPath returns Option that sets indexPath.
func WithIndexPath(path string) Option {
	return func(e *export) error {
		if path == "" {
			return errors.NewErrInvalidOption("indexPath", path)
		}
		e.indexPath = path
		return nil
	}
}

// WithGateway returns Option that sets gateway client.
func WithGateway(client vald.Client) Option {
	return func(e *export) error {
		if client == nil {
			return errors.NewErrCriticalOption("gateway", client)
		}
		e.gateway = client
		return nil
	}
}

// WithErrGroup returns Option that set errgroup.
func WithErrGroup(eg errgroup.Group) Option {
	return func(e *export) error {
		if eg != nil {
			e.eg = eg
		}
		return nil
	}
}
