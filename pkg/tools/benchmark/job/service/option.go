//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/test/data/hdf5"
)

type Option func(j *job) error

var defaultOpts = []Option{
	WithDimension(748),
	// TODO: set default config for client
}

func WithDimension(dim int) Option {
	return func(j *job) error {
		if dim > 0 {
			j.dimension = dim
		}
		return nil
	}
}

func WithInsertConfig(c *v1.InsertConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.insertConfig = c
		}
		return nil
	}
}

func WithUpdateConfig(c *v1.UpdateConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.updateConfig = c
		}
		return nil
	}
}

func WithUpsertConfig(c *v1.UpsertConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.upsertConfig = c
		}
		return nil
	}
}

func WithSearchConfig(c *v1.SearchConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.searchConfig = c
		}
		return nil
	}
}

func WithRemoveConfig(c *v1.RemoveConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.removeConfig = c
		}
		return nil
	}
}

func WithValdClient(c vald.Client) Option {
	return func(j *job) error {
		if c == nil {
			return errors.NewErrInvalidOption("client", c)
		}
		j.client = c
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(j *job) error {
		if eg == nil {
			return errors.NewErrInvalidOption("client", eg)
		}
		j.eg = eg
		return nil
	}
}

func WithHdf5(d hdf5.Data) Option {
	return func(j *job) error {
		if d == nil {
			return errors.NewErrInvalidOption("hdf5", d)
		}
		j.hdf5 = d
		return nil
	}
}

func WithDataset(d *v1.BenchmarkDataset) Option {
	return func(j *job) error {
		if d == nil {
			return errors.NewErrInvalidOption("dataset", d)
		}
		j.dataset = d
		return nil
	}
}

func WithJobTypeByString(t string) Option {
	var jt jobType
	switch t {
	case "search":
		jt = SEARCH
	}
	return WithJobType(jt)
}

func WithJobType(jt jobType) Option {
	return func(j *job) error {
		switch jt {
		case SEARCH:
			j.jobType = jt
		default:
			return errors.NewErrInvalidOption("jobType", jt)
		}
		return nil
	}
}

func WithJobFunc(jf func(context.Context, chan error) error) Option {
	return func(j *job) error {
		if jf == nil {
			return errors.NewErrInvalidOption("jobFunc", jf)
		}
		j.jobFunc = jf
		return nil
	}
}
