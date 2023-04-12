//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package rand

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func clearPool() {
	pool = sync.Pool{
		New: func() interface{} {
			return new(rand).init()
		},
	}
}

func TestUint32(t *testing.T) {
	type test struct {
		name       string
		beforeFunc func(*testing.T)
	}

	tests := []test{
		{
			name: "returns random number when pooled rand instance is nil",
			beforeFunc: func(t *testing.T) {
				t.Helper()
				clearPool()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}

			_ = Uint32()
			if atomic.LoadUint32(pool.Get().(*rand).x) == 0 {
				t.Error("r.x is 0")
			}
		})
	}
}

func TestLimitedUint32(t *testing.T) {
	type test struct {
		name       string
		beforeFunc func(*testing.T)
		max        uint64
	}

	tests := []test{
		{
			name: "returns random number less than max",
			beforeFunc: func(t *testing.T) {
				t.Helper()
				clearPool()
			},
			max: 100,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}

			got := LimitedUint32(test.max)
			if got > uint32(test.max) {
				t.Errorf("more than %v. got: %v", test.max, got)
			}
		})
	}
}

func Test_rand_Uint32(t *testing.T) {
	type test struct {
		name string
		x    *uint32
	}

	tests := []test{
		{
			name: "returns rand number when r.x is 0",
			x:    new(uint32),
		},

		func() test {
			x := uint32(100)
			return test{
				name: "returns rand number when r.x is not 0",
				x:    &x,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rand{
				x: tt.x,
			}

			_ = r.Uint32()
			if atomic.LoadUint32(r.x) == 0 {
				t.Error("r.x is 0")
			}
		})
	}
}

func Test_rand_init(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(*rand) error
	}

	tests := []test{
		{
			name: "returns rand instance",
			checkFunc: func(r *rand) error {
				if r == nil {
					return errors.New("rand is nil")
				}

				if r.x == nil {
					return errors.New("rand.x is nil")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(rand)
			if err := tt.checkFunc(r.init()); err != nil {
				t.Error(err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
