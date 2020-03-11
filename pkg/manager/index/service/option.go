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

// Package service
package service

import (
	"time"

	"github.com/vdaas/vald/internal/client/discoverer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(i *index) error

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithIndexingConcurrency(1),
		WithIndexingDuration("1m"),
		WithIndexingDurationLimit("30m"),
		WithMinUncommitted(100),
	}
)

func WithIndexingConcurrency(c int) Option {
	return func(idx *index) error {
		if c != 0 {
			idx.concurrency = c
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
			d = time.Minute
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
			d = time.Minute * 30
		}
		idx.indexDurationLimit = d
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
