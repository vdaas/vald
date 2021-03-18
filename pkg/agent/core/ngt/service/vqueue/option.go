//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"github.com/vdaas/vald/internal/errgroup"
)

type Option func(n *vqueue) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithDeleteBufferPoolSize(1000),
	WithInsertBufferPoolSize(1000),
	WithDeleteBufferSize(100),
	WithInsertBufferSize(100),
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(v *vqueue) error {
		if eg != nil {
			v.eg = eg
		}

		return nil
	}
}

func WithInsertBufferSize(size int) Option {
	return func(v *vqueue) error {
		if size > 0 {
			v.ichSize = size
		}

		return nil
	}
}

func WithDeleteBufferSize(size int) Option {
	return func(v *vqueue) error {
		if size > 0 {
			v.dchSize = size
		}

		return nil
	}
}

func WithInsertBufferPoolSize(size int) Option {
	return func(v *vqueue) error {
		if size > 0 {
			v.iBufSize = size
		}

		return nil
	}
}

func WithDeleteBufferPoolSize(size int) Option {
	return func(v *vqueue) error {
		if size > 0 {
			v.dBufSize = size
		}

		return nil
	}
}
