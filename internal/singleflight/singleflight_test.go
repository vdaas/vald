//go:build !race

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

// Package singleflight represents zero time caching
package singleflight

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type want struct {
		want Group[any]
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, Group[any]) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Group[any]) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns Group implementation",
			want: want{
				want: &group[any]{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := New[any]()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_Do(t *testing.T) {
	type args[V any] struct {
		ctx context.Context
		key string
		fn  func() (V, error)
	}
	type want[V any] struct {
		wantV      V
		wantShared bool
		err        error
	}
	type test[V any] struct {
		name       string
		args       args[V]
		want       want[V]
		beforeFunc func(args[V])
		execFunc   func(*testing.T, args[V]) (V, bool, error)
		checkFunc  func(want[V], V, bool, error) error
		afterFunc  func(args[V])
	}
	tests := []test[string]{
		func() test[string] {
			// routine1
			key1 := "req_1"
			var cnt1 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn1 := func() (string, error) {
				atomic.AddUint32(&cnt1, 1)
				return "res_1", nil
			}

			// routine 2
			key2 := "req_2"
			var cnt2 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn2 := func() (string, error) {
				atomic.AddUint32(&cnt2, 1)
				return "res_2", nil
			}

			return test[string]{
				name: "returns (v, false, nil) when Do is called with another key",
				args: args[string]{
					key: key1,
					ctx: context.Background(),
					fn:  fn1,
				},
				want: want[string]{
					wantV:      "res_1",
					wantShared: false,
					err:        nil,
				},
				execFunc: func(t *testing.T, a args[string]) (got string, gotShared bool, err error) {
					t.Helper()
					g := New[string]()

					wg := new(sync.WaitGroup)
					wg.Add(1)
					go func() {
						got, gotShared, err = g.Do(a.ctx, a.key, a.fn)
						wg.Done()
					}()

					wg.Add(1)
					go func() {
						_, _, _ = g.Do(a.ctx, key2, fn2)
						wg.Done()
					}()

					wg.Wait()
					return got, gotShared, err
				},
				checkFunc: func(w want[string], gotV string, gotShared bool, err error) error {
					if got, want := int(atomic.LoadUint32(&cnt1)), 1; got != want {
						return errors.Errorf("cnt got = %d, want = %d", got, want)
					}
					if got, want := int(atomic.LoadUint32(&cnt2)), 1; got != want {
						return errors.Errorf("cnt got = %d, want = %d", got, want)
					}
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					if !reflect.DeepEqual(gotV, w.wantV) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
					}
					if !reflect.DeepEqual(gotShared, w.wantShared) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotShared, w.wantShared)
					}
					return nil
				},
			}
		}(),
		func() test[string] {
			// routine1
			var cnt1 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn1 := func() (string, error) {
				atomic.AddUint32(&cnt1, 1)
				time.Sleep(time.Millisecond * 500)
				return "res_1", nil
			}

			// routine 2
			var cnt2 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn2 := func() (string, error) {
				atomic.AddUint32(&cnt2, 1)
				return "res_2", nil
			}

			w := want[string]{
				wantV:      "res_1",
				wantShared: true,
				err:        nil,
			}

			checkFunc := func(w want[string], gotV string, gotShared bool, err error) error {
				c1 := int(atomic.LoadUint32(&cnt1))
				c2 := int(atomic.LoadUint32(&cnt2))
				// since there is a chance that the go routine 2 is executed before routine 1, we need to check if either one is executed
				if !((c1 == 1 && c2 == 0) || (c1 == 0 && c2 == 1)) {
					return errors.Errorf("cnt1 and cnt2 is executed, %d, %d", c1, c2)
				}
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if !reflect.DeepEqual(gotV, w.wantV) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
				}
				if !reflect.DeepEqual(gotShared, w.wantShared) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotShared, w.wantShared)
				}
				return nil
			}

			return test[string]{
				name: "returns (v, true, nil) when Do is called with the same key",
				args: args[string]{
					key: "req_1",
					ctx: context.Background(),
					fn:  fn1,
				},
				want: w,
				execFunc: func(t *testing.T, a args[string]) (string, bool, error) {
					t.Helper()

					g := New[string]()
					wg := new(sync.WaitGroup)
					var got, got1 string
					var gotShared, gotShared1 bool
					var err, err1 error

					wg.Add(1)
					go func() {
						got, gotShared, err = g.Do(a.ctx, a.key, fn1)
						wg.Done()
					}()

					// call with the same key but with another function
					wg.Add(1)
					time.Sleep(time.Millisecond * 100)
					go func() {
						got1, gotShared1, err1 = g.Do(a.ctx, a.key, fn2)
						wg.Done()
					}()

					wg.Wait()

					if err := checkFunc(w, got1, gotShared1, err1); err != nil {
						t.Fatal(err)
					}

					return got, gotShared, err
				},
				checkFunc: checkFunc,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}

			gotV, gotShared, err := test.execFunc(t, test.args)

			if err := test.checkFunc(test.want, gotV, gotShared, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
