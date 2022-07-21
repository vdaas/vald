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

// Package search manages the main logic of search job.
package search

import (
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/test/data/hdf5"
)

type Option func(s *searchJob) error

var defaultOpts = []Option{
	WithDimension(748),
	WithNum(10),
	WithMinNum(10),
	WithRadius(-1),
	WithEpsilon(0.1),
	WithTimeout("1s"),
}

func WithDimension(dim int) Option {
	return func(s *searchJob) error {
		s.dimension = dim
		return nil
	}
}

func WithNum(num uint32) Option {
	return func(s *searchJob) error {
		s.num = num
		return nil
	}
}

func WithMinNum(minNum uint32) Option {
	return func(s *searchJob) error {
		s.minNum = minNum
		return nil
	}
}

func WithRadius(radius float64) Option {
	return func(s *searchJob) error {
		s.radius = radius
		return nil
	}
}

func WithEpsilon(epsilon float64) Option {
	return func(s *searchJob) error {
		s.epsilon = epsilon
		return nil
	}
}

func WithTimeout(timeout string) Option {
	return func(s *searchJob) error {
		s.timeout = timeout
		return nil
	}
}

func WithValdClient(c vald.Client) Option {
	return func(s *searchJob) error {
		s.client = c
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(s *searchJob) error {
		s.eg = eg
		return nil
	}
}

func WithHdf5(d hdf5.Data) Option {
	return func(s *searchJob) error {
		s.hdf5 = d
		return nil
	}
}
