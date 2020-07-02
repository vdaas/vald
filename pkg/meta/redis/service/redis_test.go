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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/tcp"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.Redis
	}
	type want struct {
		want Redis
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Redis, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Redis, err error) error {
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
		           cfg: nil,
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
		           cfg: nil,
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

			got, err := New(test.args.cfg)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Disconnect(t *testing.T) {
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
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
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			err := c.Disconnect()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Connect(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			err := c.Connect(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Get(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           key: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.Get(test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_GetMultiple(t *testing.T) {
	type args struct {
		keys []string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		wantVals []string
		err      error
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
	defaultCheckFunc := func(w want, gotVals []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVals, w.wantVals) {
			return errors.Errorf("got = %v, want %v", gotVals, w.wantVals)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           keys: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           keys: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			gotVals, err := c.GetMultiple(test.args.keys...)
			if err := test.checkFunc(test.want, gotVals, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_GetInverse(t *testing.T) {
	type args struct {
		val string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           val: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           val: "",
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.GetInverse(test.args.val)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_GetInverseMultiple(t *testing.T) {
	type args struct {
		vals []string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           vals: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           vals: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.GetInverseMultiple(test.args.vals...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_appendPrefix(t *testing.T) {
	type args struct {
		prefix string
		key    string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           prefix: "",
		           key: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           prefix: "",
		           key: "",
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got := c.appendPrefix(test.args.prefix, test.args.key)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_get(t *testing.T) {
	type args struct {
		prefix string
		key    string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		wantVal string
		err     error
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
	defaultCheckFunc := func(w want, gotVal string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVal, w.wantVal) {
			return errors.Errorf("got = %v, want %v", gotVal, w.wantVal)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           prefix: "",
		           key: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           prefix: "",
		           key: "",
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			gotVal, err := c.get(test.args.prefix, test.args.key)
			if err := test.checkFunc(test.want, gotVal, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_getMulti(t *testing.T) {
	type args struct {
		prefix string
		keys   []string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		wantVals []string
		err      error
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
	defaultCheckFunc := func(w want, gotVals []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVals, w.wantVals) {
			return errors.Errorf("got = %v, want %v", gotVals, w.wantVals)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           prefix: "",
		           keys: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           prefix: "",
		           keys: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			gotVals, err := c.getMulti(test.args.prefix, test.args.keys...)
			if err := test.checkFunc(test.want, gotVals, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Set(t *testing.T) {
	type args struct {
		key string
		val string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           key: "",
		           val: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           val: "",
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			err := c.Set(test.args.key, test.args.val)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_SetMultiple(t *testing.T) {
	type args struct {
		kvs map[string]string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           kvs: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           kvs: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			err := c.SetMultiple(test.args.kvs)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Delete(t *testing.T) {
	type args struct {
		key string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           key: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.Delete(test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_DeleteMultiple(t *testing.T) {
	type args struct {
		keys []string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           keys: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           keys: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.DeleteMultiple(test.args.keys...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_DeleteInverse(t *testing.T) {
	type args struct {
		val string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           val: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           val: "",
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.DeleteInverse(test.args.val)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_DeleteInverseMultiple(t *testing.T) {
	type args struct {
		vals []string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
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
		           vals: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           vals: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			got, err := c.DeleteInverseMultiple(test.args.vals...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_delete(t *testing.T) {
	type args struct {
		pfx    string
		pfxInv string
		key    string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		wantVal string
		err     error
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
	defaultCheckFunc := func(w want, gotVal string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVal, w.wantVal) {
			return errors.Errorf("got = %v, want %v", gotVal, w.wantVal)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           pfx: "",
		           pfxInv: "",
		           key: "",
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           pfx: "",
		           pfxInv: "",
		           key: "",
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			gotVal, err := c.delete(test.args.pfx, test.args.pfxInv, test.args.key)
			if err := test.checkFunc(test.want, gotVal, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_deleteMulti(t *testing.T) {
	type args struct {
		pfx    string
		pfxInv string
		keys   []string
	}
	type fields struct {
		db              redis.Redis
		opts            []redis.Option
		topts           []tcp.DialerOption
		kvPrefix        string
		vkPrefix        string
		prefixDelimiter string
	}
	type want struct {
		wantVals []string
		err      error
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
	defaultCheckFunc := func(w want, gotVals []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVals, w.wantVals) {
			return errors.Errorf("got = %v, want %v", gotVals, w.wantVals)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           pfx: "",
		           pfxInv: "",
		           keys: nil,
		       },
		       fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
		           pfx: "",
		           pfxInv: "",
		           keys: nil,
		           },
		           fields: fields {
		           db: nil,
		           opts: nil,
		           topts: nil,
		           kvPrefix: "",
		           vkPrefix: "",
		           prefixDelimiter: "",
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
			c := &client{
				db:              test.fields.db,
				opts:            test.fields.opts,
				topts:           test.fields.topts,
				kvPrefix:        test.fields.kvPrefix,
				vkPrefix:        test.fields.vkPrefix,
				prefixDelimiter: test.fields.prefixDelimiter,
			}

			gotVals, err := c.deleteMulti(test.args.pfx, test.args.pfxInv, test.args.keys...)
			if err := test.checkFunc(test.want, gotVals, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
