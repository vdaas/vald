// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(_ *export) error

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
func WithGateway(gw vc.Client) Option {
	return func(e *export) error {
		if gw == nil {
			return errors.NewErrCriticalOption("gateway", gw)
		}
		e.gateway = gw
		return nil
	}
}
