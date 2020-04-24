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
	"sync"
	"testing"

	"github.com/cockroachdb/errors"
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           size: 0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           size: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
		mu sync.RWMutex
		m  map[string]*call
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
		checkFunc  func(want, interface{}, bool, error) error
		beforeFunc func(args)
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

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
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
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &group{
				mu: test.fields.mu,
				m:  test.fields.m,
			}

			// TODO: refactor singleflight.Do
			gotV, err, gotShared := g.Do(test.args.ctx, test.args.key, test.args.fn)
			if err := test.checkFunc(test.want, gotV, gotShared, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
