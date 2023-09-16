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

// Package service
package service

import (
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(i *index) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithIndexingConcurrency(1),
	WithIndexingDuration("1m"),
	WithIndexingDurationLimit("30m"),
	WithSaveIndexDurationLimit("3h"),
	WithMinUncommitted(100),
	WithCreationPoolSize(10000),
}

func WithIndexingConcurrency(c int) Option {
	return func(idx *index) error {
		if c != 0 {
			idx.createIndexConcurrency = c
		}
		return nil
	}
}

func WithSaveConcurrency(c int) Option {
	return func(idx *index) error {
		if c != 0 {
			idx.saveIndexConcurrency = c
		}
		return nil
	}
}

func WithIndexingDuration(dur string) Option {
	return func(idx *index) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		idx.indexDuration = d
		return nil
	}
}

func WithIndexingDurationLimit(dur string) Option {
	return func(idx *index) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		idx.indexDurationLimit = d
		return nil
	}
}

func WithSaveIndexDurationLimit(dur string) Option {
	return func(idx *index) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		idx.saveIndexDurationLimit = d
		return nil
	}
}

func WithMinUncommitted(n uint32) Option {
	return func(idx *index) error {
		if n > 0 {
			idx.minUncommitted = n
		}
		return nil
	}
}

func WithCreationPoolSize(size uint32) Option {
	return func(idx *index) error {
		if size > 0 {
			idx.creationPoolSize = size
		}
		return nil
	}
}

func WithDiscoverer(c discoverer.Client) Option {
	return func(idx *index) error {
		if c != nil {
			idx.client = c
		}
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(idx *index) error {
		if eg != nil {
			idx.eg = eg
		}
		return nil
	}
}
