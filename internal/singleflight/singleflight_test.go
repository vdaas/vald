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
	"fmt"
	"reflect"
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
	type want struct {
		wantV      interface{}
		wantShared bool
		err        error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		cond       *sync.Cond
		mu         *sync.Mutex
		wg         *sync.WaitGroup
		checkFunc  func(want, interface{}, bool, error) error
		beforeFunc func(Group, args)
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           key: "",
		           fn: nil,
		       },
		       fields: fields {
		           mu: sync.RWMutex{},
		           m: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		func() test {
			mu := new(sync.Mutex)
			cond := sync.NewCond(mu)
			cnt := uint32(0)
			wg := new(sync.WaitGroup)
			return test{
				mu:   mu,
				cond: cond,
				wg:   wg,
				name: "returns (v, shared, nil) when Do is called with another key",
				args: args{
					ctx: context.Background(),
					key: "req_1",
					fn: func() (interface{}, error) {
						atomic.AddUint32(&cnt, 1)
						return "res_1", nil
					},
				},
				fields: fields{
					m: make(map[string]*call, 2),
				},
				want: want{
					wantV:      "res_1",
					wantShared: false,
					err:        nil,
				},
				beforeFunc: func(g Group, args args) {
					wg.Add(1)
					go func() {
						mu.Lock()
						defer mu.Unlock()
						cond.Wait()
						g.Do(context.Background(), "req_2", func() (interface{}, error) {
							defer wg.Done()
							atomic.AddUint32(&cnt, 1)
							return "res_2", nil
						})
					}()
				},
				checkFunc: func(want, interface{}, bool, error) error {
					if got, want := int(atomic.LoadUint32(&cnt)), 2; got != want {
						return errors.Errorf("cnt got = %d, want = %d", got, want)
					}
					return nil
				},
			}
		}(),
		func() test {
			mu := new(sync.Mutex)
			cond := sync.NewCond(mu)
			cnt := uint32(0)
			wg := new(sync.WaitGroup)
			return test{
				name: "returns (v, shared, nil) when Do is called with same key",
				mu:   mu,
				cond: cond,
				wg:   wg,
				args: args{
					ctx: context.Background(),
					key: "req_1",
					fn: func() (interface{}, error) {
						fmt.Println("args")
						atomic.AddUint32(&cnt, 1)
						return "res_1", nil
					},
				},
				fields: fields{
					m: make(map[string]*call, 2),
				},
				want: want{
					wantV:      "res_1",
					wantShared: true,
					err:        nil,
				},
				beforeFunc: func(g Group, args args) {
					wg.Add(1)
					ch := make(chan struct{})
					go func() {
						g.Do(context.Background(), "req_1", func() (interface{}, error) {
							ch <- struct{}{}
							fmt.Println("test")
							defer wg.Done()
							time.Sleep(time.Second * 10)
							return "res_1", nil
						})
					}()
					<- ch
					for i := 0; i < 10; i++ {
						wg.Add(1)
						go func(i int) {
							mu.Lock()
							defer mu.Unlock()
							cond.Wait()
							defer wg.Done()
							fmt.Println(i)
							g.Do(context.Background(), "req_1", func() (interface{}, error) {
								atomic.AddUint32(&cnt, 1)
								return "res_1", nil
							})
						}(i)
					}
				},
				checkFunc: func(want, interface{}, bool, error) error {
					if got, want := int(atomic.LoadUint32(&cnt)), 1; got != want {
						return errors.Errorf("cnt got = %d, want = %d", got, want)
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

			var gotV interface{}
			var gotShared bool
			var err error
			test.wg.Add(1)
			go func() {
				test.mu.Lock()
				defer test.mu.Unlock()
				test.cond.Wait()
				defer test.wg.Done()
				gotV, gotShared, err = g.Do(test.args.ctx, test.args.key, test.args.fn)
			}()
			time.Sleep(time.Second)
			test.cond.Broadcast()
			test.wg.Wait()
			if err := test.checkFunc(test.want, gotV, gotShared, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
