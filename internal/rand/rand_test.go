// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package rand

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/goleak"
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

func TestFloat32(t *testing.T) {
	type want struct {
		min float32
		max float32
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, float32) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got float32) error {
		if w.min > got || w.max < got {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w)
		}
		return nil
	}
	count := 1000
	tests := func() []test {
		tests := make([]test, count)
		for idx := range tests {
			tests[idx] = test{
				name: fmt.Sprint(idx),
				want: want{
					min: 0.0,
					max: 1.0,
				},
			}
		}
		return tests
	}()
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Float32()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
