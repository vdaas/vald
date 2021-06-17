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

// Package pool provides pool functionality for pooling buffer or others
package pool

import (
	"bytes"
	"sync/atomic"
)

// Option represents the functional option for pool.
type Option func(*pool)

var defaultOptions = []Option{
	WithSize(64),
	WithConstructor(func(size uint64) interface{} {
		return bytes.NewBuffer(make([]byte, 0, size))
	}),
}

func WithSize(size uint64) Option {
	return func(p *pool) {
		if size > 0 {
			atomic.StoreUint64(&p.size, size)
		}
	}
}

func WithLimit(limit uint64) Option {
	return func(p *pool) {
		atomic.StoreUint64(&p.limit, limit)
	}
}

func WithConstructor(f func(size uint64) interface{}) Option {
	return func(p *pool) {
		if f != nil {
			p.new = f
		}
	}
}
