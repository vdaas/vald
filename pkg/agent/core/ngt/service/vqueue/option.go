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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for vqueue.
type Option func(n *vqueue) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithDeleteBufferPoolSize(1000),
	WithInsertBufferPoolSize(1000),
}

// WithErrGroup returns the option to set the errgroup.
func WithErrGroup(eg errgroup.Group) Option {
	return func(v *vqueue) error {
		if eg == nil {
			return errors.NewErrInvalidOption("errgroup", eg)
		}
		v.eg = eg

		return nil
	}
}

// WithInsertBufferPoolSize returns the option to set the pool size of the insert buffer.
func WithInsertBufferPoolSize(size int) Option {
	return func(v *vqueue) error {
		if size <= 0 {
			return errors.NewErrInvalidOption("insertBufferPoolSize", size)
		}
		v.iBufSize = size

		return nil
	}
}

// WithDeleteBufferPoolSize returns the option to set the pool size of the delete buffer.
func WithDeleteBufferPoolSize(size int) Option {
	return func(v *vqueue) error {
		if size <= 0 {
			return errors.NewErrInvalidOption("deleteBufferPoolSize", size)
		}
		v.dBufSize = size

		return nil
	}
}
