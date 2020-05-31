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

// Package service provides meta service
package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"go.uber.org/goleak"
)

func TestNewMeta(t *testing.T) {
	type args struct {
		opts []MetaOption
	}
	type want struct {
		wantMi Meta
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Meta, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotMi Meta, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotMi, w.wantMi) {
			return errors.Errorf("got = %v, want %v", gotMi, w.wantMi)
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

			gotMi, err := NewMeta(test.args.opts...)
			if err := test.checkFunc(test.want, gotMi, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_Exists(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want bool
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		           meta: "",
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           meta: "",
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.Exists(test.args.ctx, test.args.meta)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_GetMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		wantV string
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotV string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotV, w.wantV) {
			return errors.Errorf("got = %v, want %v", gotV, w.wantV)
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
		           uuid: "",
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           uuid: "",
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			gotV, err := m.GetMeta(test.args.ctx, test.args.uuid)
			if err := test.checkFunc(test.want, gotV, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_GetMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want []string
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		           uuids: nil,
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           uuids: nil,
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.GetMetas(test.args.ctx, test.args.uuids...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_GetUUID(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		wantK string
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotK string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotK, w.wantK) {
			return errors.Errorf("got = %v, want %v", gotK, w.wantK)
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
		           meta: "",
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           meta: "",
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			gotK, err := m.GetUUID(test.args.ctx, test.args.meta)
			if err := test.checkFunc(test.want, gotK, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_GetUUIDs(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want []string
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		           metas: nil,
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           metas: nil,
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.GetUUIDs(test.args.ctx, test.args.metas...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_SetUUIDandMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
		meta string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           uuid: "",
		           meta: "",
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           uuid: "",
		           meta: "",
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			err := m.SetUUIDandMeta(test.args.ctx, test.args.uuid, test.args.meta)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_SetUUIDandMetas(t *testing.T) {
	type args struct {
		ctx context.Context
		kvs map[string]string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           kvs: nil,
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           kvs: nil,
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			err := m.SetUUIDandMetas(test.args.ctx, test.args.kvs)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_DeleteMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		wantV string
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotV string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotV, w.wantV) {
			return errors.Errorf("got = %v, want %v", gotV, w.wantV)
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
		           uuid: "",
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           uuid: "",
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			gotV, err := m.DeleteMeta(test.args.ctx, test.args.uuid)
			if err := test.checkFunc(test.want, gotV, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_DeleteMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want []string
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		           uuids: nil,
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           uuids: nil,
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.DeleteMetas(test.args.ctx, test.args.uuids...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_DeleteUUID(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want string
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		           meta: "",
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           meta: "",
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.DeleteUUID(test.args.ctx, test.args.meta)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_meta_DeleteUUIDs(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []string
	}
	type fields struct {
		addr                string
		client              grpc.Client
		cache               cache.Cache
		enableCache         bool
		expireCheckDuration string
		expireDuration      string
	}
	type want struct {
		want []string
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
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
		           ctx: nil,
		           metas: nil,
		       },
		       fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
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
		           metas: nil,
		           },
		           fields: fields {
		           addr: "",
		           client: nil,
		           cache: nil,
		           enableCache: false,
		           expireCheckDuration: "",
		           expireDuration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			m := &meta{
				addr:                test.fields.addr,
				client:              test.fields.client,
				cache:               test.fields.cache,
				enableCache:         test.fields.enableCache,
				expireCheckDuration: test.fields.expireCheckDuration,
				expireDuration:      test.fields.expireDuration,
			}

			got, err := m.DeleteUUIDs(test.args.ctx, test.args.metas...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
