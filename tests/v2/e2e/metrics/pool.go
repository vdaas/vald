//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package metrics provides a metrics collection and aggregation system for E2E tests.
package metrics

import "github.com/vdaas/vald/internal/sync"

type PoolProvider[V any] interface {
	Get() *V
	Put(*V)
}
type poolProvider[V any] struct {
	pool sync.Pool
}

func newPoolProvider[V any](gen func() *V) PoolProvider[V] {
	return &poolProvider[V]{
		pool: sync.Pool{
			New: func() any {
				return gen()
			},
		},
	}
}

func (p *poolProvider[V]) Get() (v *V) {
	pv := p.pool.Get()
	if pv == nil {
		return nil
	}
	var ok bool
	vp, ok := pv.(V)
	if !ok {
		return nil
	}
	return &vp
}

func (p *poolProvider[V]) Put(v *V) {
	p.pool.Put(v)
}

var (
	requestResultPool = newPoolProvider(func() *RequestResult {
		return new(RequestResult)
	})
)
