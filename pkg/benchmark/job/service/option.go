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
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/data/hdf5"
)

type Option func(j *job) error

var defaultOpts = []Option{
	WithDimension(748),
	WithNum(10),
	WithMinNum(10),
	WithRadius(-1),
	WithEpsilon(0.1),
	WithTimeout("1s"),
}

func WithDimension(dim int) Option {
	return func(j *job) error {
		if dim > 0 {
			j.dimension = dim
		}
		return nil
	}
}

func WithIter(iter int) Option {
	return func(j *job) error {
		if iter > 0 {
			j.iter = iter
		}
		return nil
	}
}

func WithNum(num uint32) Option {
	return func(j *job) error {
		if num > 0 {
			j.num = num
		}
		return nil
	}
}

func WithMinNum(minNum uint32) Option {
	return func(j *job) error {
		if minNum > 0 {
			j.minNum = minNum
		}
		return nil
	}
}

func WithRadius(radius float64) Option {
	return func(j *job) error {
		j.radius = radius
		return nil
	}
}

func WithEpsilon(epsilon float64) Option {
	return func(j *job) error {
		j.epsilon = epsilon
		return nil
	}
}

func WithTimeout(timeout string) Option {
	return func(j *job) error {
		dur, err := time.ParseDuration(timeout)
		if err != nil {
			return errors.NewErrInvalidOption("timeout", timeout, err)
		}
		j.timeout = dur
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
