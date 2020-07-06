//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		size int
	}
	type want struct {
		want Group
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Group) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Group) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns Group when size is 0",
			want: want{
				want: &group{
					m: make(map[string]*call, 1),
				},
			},
		},
		{
			name: "returns Group when size is 1",
			args: args{
				size: 1,
			},
			want: want{
				want: &group{
					m: make(map[string]*call, 1),
				},
			},
		},
		{
			name: "returns Group when size is over than 1",
			args: args{
				size: 2,
			},
			want: want{
				want: &group{
					m: make(map[string]*call, 2),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := New(test.args.size)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_group_Do(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		fn  func() (interface{}, error)
	}
	type fields struct {
		m map[string]*call
	}
	type util struct {
		mu         *sync.Mutex
		wg         *sync.WaitGroup
		cond       *sync.Cond
		condWaitFn func()
	}
	type want struct {
		wantV      interface{}
		wantShared bool
		err        error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		util       util
		want       want
		checkFunc  func(want, interface{}, bool, error) error
		beforeFunc func(Group, args)
		execFunc   func(*testing.T, Group)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotV interface{}, gotShared bool, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotV, w.wantV) {
			return errors.Errorf("got = %v, want %v", gotV, w.wantV)
		}
		if !reflect.DeepEqual(gotShared, w.wantShared) {
			return errors.Errorf("got = %v, want %v", gotShared, w.wantShared)
		}
		return nil
	}
	tests := []test{
		func() test {
			var calledCnt uint32

			var (
				mu         = new(sync.Mutex)
				cond       = sync.NewCond(mu)
				wg         = new(sync.WaitGroup)
				condWaitFn = func() {
					mu.Lock()
					defer mu.Unlock()
					cond.Wait()
				}
			)

			return test{
				name: "returns (v, shared, nil) when Do is called with another key",
				fields: fields{
					m: make(map[string]*call),
				},
				util: util{
					mu:         mu,
					cond:       cond,
					wg:         wg,
					condWaitFn: condWaitFn,
				},
				args: args{
					key: "req_1",
					ctx: context.Background(),
					fn: func() (interface{}, error) {
						atomic.AddUint32(&calledCnt, 1)
						return "res_1", nil
					},
				},
				want: want{
					wantV:      "res_1",
					wantShared: false,
					err:        nil,
				},
				beforeFunc: func(g Group, args args) {
					gcnt := 10
					ch := make(chan struct{}, gcnt)

					for i := 0; i < gcnt; i++ {
						wg.Add(1)
						go func(i int) {
							ch <- struct{}{}
							defer wg.Done()
							condWaitFn()

							g.Do(context.Background(), strconv.Itoa(i), func() (interface{}, error) {
								time.Sleep(time.Nanosecond * 100)
								atomic.AddUint32(&calledCnt, 1)
								return "vdaas/vald", nil
							})
						}(i)
					}

					for i := 0; i < gcnt; i++ {
						<-ch
					}
					close(ch)
				},
				checkFunc: func(w want, gotV interface{}, gotShared bool, err error) error {
					if got, want := int(atomic.LoadUint32(&calledCnt)), 11; got != want {
						return errors.Errorf("calledCnt got = %d, want = %d", got, want)
					}

					if err := defaultCheckFunc(w, gotV, gotShared, err); err != nil {
						return err
					}
					return nil
				},
			}
		}(),

		func() test {
			var calledCnt uint32

			var (
				mu         = new(sync.Mutex)
				cond       = sync.NewCond(mu)
				wg         = new(sync.WaitGroup)
				condWaitFn = func() {
					mu.Lock()
					defer mu.Unlock()
					cond.Wait()
				}
			)

			return test{
				name: "returns (v, shared, nil) when Do is called with same key",
				args: args{
					key: "req_1",
					ctx: context.Background(),
					fn: func() (interface{}, error) {
						atomic.AddUint32(&calledCnt, 1)
						return "res_1", nil
					},
				},
				fields: fields{
					m: make(map[string]*call),
				},
				util: util{
					mu:         mu,
					cond:       cond,
					wg:         wg,
					condWaitFn: condWaitFn,
				},
				want: want{
					wantV:      "res_1",
					wantShared: true,
					err:        nil,
				},
				beforeFunc: func(g Group, args args) {
					wg.Add(1)
					go func() {
						defer wg.Done()
						g.Do(context.Background(), args.key, func() (interface{}, error) {
							time.Sleep(3 * time.Second)
							return args.fn()
						})
					}()

					gcnt := 10
					ch := make(chan struct{}, gcnt)

					for i := 0; i < gcnt; i++ {
						wg.Add(1)
						go func() {
							ch <- struct{}{}
							defer wg.Done()
							condWaitFn()

							g.Do(context.Background(), args.key, func() (interface{}, error) {
								atomic.AddUint32(&calledCnt, 1)
								return "vdaas/vald", nil
							})
						}()
					}

					for i := 0; i < gcnt; i++ {
						<-ch
					}
					close(ch)
				},
				checkFunc: func(w want, gotV interface{}, gotShared bool, err error) error {
					if got, want := int(atomic.LoadUint32(&calledCnt)), 1; got != want {
						return errors.Errorf("calledCnt got = %d, want = %d", got, want)
					}

					if err := defaultCheckFunc(w, gotV, gotShared, err); err != nil {
						return err
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			g := &group{
				m: test.fields.m,
			}

			if test.beforeFunc != nil {
				test.beforeFunc(g, test.args)
			}

			var (
				gotV      interface{}
				gotShared bool
				err       error
			)

			test.util.wg.Add(1)
			go func() {
				defer test.util.wg.Done()
				gotV, gotShared, err = g.Do(context.Background(), test.args.key, test.args.fn)
			}()

			test.util.cond.Broadcast()
			test.util.wg.Wait()

			if err := test.checkFunc(test.want, gotV, gotShared, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
