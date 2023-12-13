//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package gache provides implementation of cache using gache
package gache

import (
	"context"
	"reflect"
	"testing"
	"time"

	gache "github.com/kpango/gache/v2"
	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option[any]
	}
	type want struct {
		wantC cacher.Cache[any]
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, cacher.Cache[any]) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got cacher.Cache[any]) error {
		wc := reflect.ValueOf(w.wantC.(*cache[any]))
		gc := reflect.ValueOf(got.(*cache[any]))
		flag := false
		for i := 0; i < reflect.Indirect(gc).NumField(); i++ {
			flag = reflect.DeepEqual(reflect.Indirect(gc).Field(i), reflect.Indirect(wc).Field(i))
		}
		if flag {
			return errors.Errorf("got: \"%#v\",\n\t\t\twant: \"%#v\"", got, w.wantC)
		}
		return nil
	}
	tests := []test{
		func() test {
			c := new(cache[any])
			for _, opt := range defaultOptions[any]() {
				opt(c)
			}
			c.gache.SetDefaultExpire(c.expireDur)
			return test{
				name: "set success when opts is nil",
				want: want{
					wantC: c,
				},
			}
		}(),
		func() test {
			expiredHook := func(context.Context, string) {}
			c := new(cache[any])
			for _, opt := range append(defaultOptions[any](), WithExpiredHook[any](expiredHook)) {
				opt(c)
			}
			c.gache.SetDefaultExpire(c.expireDur)
			if c.expiredHook != nil {
				c.gache = c.gache.SetExpiredHook(c.expiredHook).EnableExpiredHook()
			}
			return test{
				name: "set success when opts is not nil",
				args: args{
					opts: []Option[any]{
						WithExpiredHook[any](expiredHook),
					},
				},
				want: want{
					wantC: c,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			gotC := New(test.args.opts...)
			if err := checkFunc(test.want, gotC); err != nil {
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
		gache          gache.Gache[any]
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct{}
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
		{
			name: "Call Start",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					return ctx
				}(),
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &cache[any]{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}
			c.Start(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
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
		gache          gache.Gache[any]
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
		beforeFunc func(*testing.T, args, *cache[any])
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		{
			name: "Call Get when gache is empty",
			args: args{
				key: "vdaas",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				want:  nil,
				want1: false,
			},
		},
		{
			name: "Call Get when gache is not empty",
			args: args{
				key: "vdaas",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				want:  "vald",
				want1: true,
			},
			beforeFunc: func(t *testing.T, args args, c *cache[any]) {
				t.Helper()
				c.Set(args.key, "vald")
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			c := &cache[any]{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, c)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got, got1 := c.Get(test.args.key)
			if err := checkFunc(test.want, got, got1); err != nil {
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
		gache          gache.Gache[any]
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
		key   string
		want  interface{}
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *cache[any]) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, c *cache[any]) error {
		got, got1 := c.Get(w.key)
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want = %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want = %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		{
			name: "Call Set",
			args: args{
				key: "vdaas",
				val: "vald",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				key:   "vdaas",
				want:  "vald",
				want1: true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &cache[any]{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}

			c.Set(test.args.key, test.args.val)
			if err := checkFunc(test.want, c); err != nil {
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
		gache          gache.Gache[any]
		expireDur      time.Duration
		expireCheckDur time.Duration
		expiredHook    func(context.Context, string)
	}
	type want struct {
		key   string
		want  interface{}
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *cache[any]) error
		beforeFunc func(*testing.T, args, *cache[any])
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, c *cache[any]) error {
		got, got1 := c.Get(w.key)
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want = %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want = %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		{
			name: "Call Delete when gache is empty",
			args: args{
				key: "vdaas",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				key:   "vdaas",
				want:  nil,
				want1: false,
			},
		},
		{
			name: "Call Delete when gache is not empty",
			args: args{
				key: "vdaas",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				key:   "vdaas",
				want:  nil,
				want1: false,
			},
			beforeFunc: func(t *testing.T, args args, c *cache[any]) {
				t.Helper()
				c.Set(args.key, "vald")
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			c := &cache[any]{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, c)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			c.Delete(test.args.key)
			if err := checkFunc(test.want, c); err != nil {
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
		gache          gache.Gache[any]
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
		beforeFunc func(*testing.T, args, *cache[any])
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		{
			name: "Call GetAndDelete when gache is empty",
			args: args{
				key: "vdaas",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				want:  nil,
				want1: false,
			},
		},
		{
			name: "Call GetAndDelete when gache is not empty",
			args: args{
				key: "vdaas",
			},
			fields: fields{
				gache:          gache.New[any](),
				expireDur:      1 * time.Second,
				expireCheckDur: 1 * time.Second,
				expiredHook:    nil,
			},
			want: want{
				want:  "vald",
				want1: true,
			},
			beforeFunc: func(t *testing.T, args args, c *cache[any]) {
				t.Helper()
				c.Set(args.key, "vald")
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			c := &cache[any]{
				gache:          test.fields.gache,
				expireDur:      test.fields.expireDur,
				expireCheckDur: test.fields.expireCheckDur,
				expiredHook:    test.fields.expiredHook,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, c)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, got1 := c.GetAndDelete(test.args.key)
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
