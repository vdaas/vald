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

// Package gache provides implementation of cache using gache
package gache

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/errors"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantC *cache
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *cache) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotC *cache) error {
		if !reflect.DeepEqual(gotC, w.wantC) {
			return errors.Errorf("got = %v, want %v", gotC, w.wantC)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
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
		           opts: nil,
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

			gotC := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotC); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_cache_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		gache          gache.Gache
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
		           },
		           fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
			c := &cache{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}

			c.Start(test.args.ctx)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cache_Get(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		gache          gache.Gache
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
		want  interface{}
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           key: "",
		       },
		       fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
		           key: "",
		           },
		           fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
			c := &cache{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}

			got, got1 := c.Get(test.args.key)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_cache_Set(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	type fields struct {
		gache          gache.Gache
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           key: "",
		           val: nil,
		       },
		       fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
		           key: "",
		           val: nil,
		           },
		           fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
			c := &cache{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}

			c.Set(test.args.key, test.args.val)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cache_Delete(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		gache          gache.Gache
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           key: "",
		       },
		       fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
		           key: "",
		           },
		           fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
			c := &cache{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}

			c.Delete(test.args.key)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cache_GetAndDelete(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		gache          gache.Gache
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
		want  interface{}
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           key: "",
		       },
		       fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
		           key: "",
		           },
		           fields: fields {
		           gache: nil,
		           expireDur: nil,
		           expireCheckDur: nil,
		           expiredHook: nil,
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
			c := &cache{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}

			got, got1 := c.GetAndDelete(test.args.key)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
