// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package bbolt

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

type Option func(*bolt.Options) error

func WithTimeout(dur time.Duration) Option {
	return func(opts *bolt.Options) error {
		opts.Timeout = dur
		return nil
	}
}

func WithNoGrowSync(noGrowSync bool) Option {
	return func(opts *bolt.Options) error {
		opts.NoGrowSync = noGrowSync
		return nil
	}
}

func WithNoFreeListSync(noFreeListSync bool) Option {
	return func(opts *bolt.Options) error {
		opts.NoFreelistSync = noFreeListSync
		return nil
	}
}

func WithPreLoadFreelist(preloadFreelist bool) Option {
	return func(opts *bolt.Options) error {
		opts.NoFreelistSync = preloadFreelist
		return nil
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(opts *bolt.Options) error {
		opts.ReadOnly = readOnly
		return nil
	}
}

func WithMmapFlags(flags int) Option {
	return func(opts *bolt.Options) error {
		opts.MmapFlags = flags
		return nil
	}
}

func WithInitialMmapSize(initialMmapSize int) Option {
	return func(opts *bolt.Options) error {
		opts.InitialMmapSize = initialMmapSize
		return nil
	}
}

func WithPageSize(pageSize int) Option {
	return func(opts *bolt.Options) error {
		opts.PageSize = pageSize
		return nil
	}
}

func WithNoSync(noSync bool) Option {
	return func(opts *bolt.Options) error {
		opts.NoSync = noSync
		return nil
	}
}

func WithMlock(mlock bool) Option {
	return func(opts *bolt.Options) error {
		opts.Mlock = mlock
		return nil
	}
}
