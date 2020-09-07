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
package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
	"go.uber.org/goleak"
)

func TestNewLoader(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Loader
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Loader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Loader, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := NewLoader(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_Prepare(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	type fields struct {
		eg               errgroup.Group
		client           igrpc.Client
		addr             string
		concurrency      int
		batchSize        int
		dataset          string
		progressDuration time.Duration
		loaderFunc       loadFunc
		dataProvider     func() interface{}
		dataSize         int
		service          config.Service
		operation        config.Operation
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           in0: nil,
		       },
		       fields: fields {
		           eg: nil,
		           client: nil,
		           addr: "",
		           concurrency: 0,
		           batchSize: 0,
		           dataset: "",
		           progressDuration: nil,
		           loaderFunc: nil,
		           dataProvider: nil,
		           dataSize: 0,
		           service: nil,
		           operation: nil,
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
		           in0: nil,
		           },
		           fields: fields {
		           eg: nil,
		           client: nil,
		           addr: "",
		           concurrency: 0,
		           batchSize: 0,
		           dataset: "",
		           progressDuration: nil,
		           loaderFunc: nil,
		           dataProvider: nil,
		           dataSize: 0,
		           service: nil,
		           operation: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &loader{
				eg:               test.fields.eg,
				client:           test.fields.client,
				addr:             test.fields.addr,
				concurrency:      test.fields.concurrency,
				batchSize:        test.fields.batchSize,
				dataset:          test.fields.dataset,
				progressDuration: test.fields.progressDuration,
				loaderFunc:       test.fields.loaderFunc,
				dataProvider:     test.fields.dataProvider,
				dataSize:         test.fields.dataSize,
				service:          test.fields.service,
				operation:        test.fields.operation,
			}

			err := l.Prepare(test.args.in0)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_Do(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		eg               errgroup.Group
		client           igrpc.Client
		addr             string
		concurrency      int
		batchSize        int
		dataset          string
		progressDuration time.Duration
		loaderFunc       loadFunc
		dataProvider     func() interface{}
		dataSize         int
		service          config.Service
		operation        config.Operation
	}
	type want struct {
		want <-chan error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
		           eg: nil,
		           client: nil,
		           addr: "",
		           concurrency: 0,
		           batchSize: 0,
		           dataset: "",
		           progressDuration: nil,
		           loaderFunc: nil,
		           dataProvider: nil,
		           dataSize: 0,
		           service: nil,
		           operation: nil,
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
		           eg: nil,
		           client: nil,
		           addr: "",
		           concurrency: 0,
		           batchSize: 0,
		           dataset: "",
		           progressDuration: nil,
		           loaderFunc: nil,
		           dataProvider: nil,
		           dataSize: 0,
		           service: nil,
		           operation: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &loader{
				eg:               test.fields.eg,
				client:           test.fields.client,
				addr:             test.fields.addr,
				concurrency:      test.fields.concurrency,
				batchSize:        test.fields.batchSize,
				dataset:          test.fields.dataset,
				progressDuration: test.fields.progressDuration,
				loaderFunc:       test.fields.loaderFunc,
				dataProvider:     test.fields.dataProvider,
				dataSize:         test.fields.dataSize,
				service:          test.fields.service,
				operation:        test.fields.operation,
			}

			got := l.Do(test.args.ctx)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_do(t *testing.T) {
	type args struct {
		ctx    context.Context
		f      func(interface{}, error)
		notify func(context.Context, error)
	}
	type fields struct {
		eg               errgroup.Group
		client           igrpc.Client
		addr             string
		concurrency      int
		batchSize        int
		dataset          string
		progressDuration time.Duration
		loaderFunc       loadFunc
		dataProvider     func() interface{}
		dataSize         int
		service          config.Service
		operation        config.Operation
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
		           f: nil,
		           notify: nil,
		       },
		       fields: fields {
		           eg: nil,
		           client: nil,
		           addr: "",
		           concurrency: 0,
		           batchSize: 0,
		           dataset: "",
		           progressDuration: nil,
		           loaderFunc: nil,
		           dataProvider: nil,
		           dataSize: 0,
		           service: nil,
		           operation: nil,
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
		           f: nil,
		           notify: nil,
		           },
		           fields: fields {
		           eg: nil,
		           client: nil,
		           addr: "",
		           concurrency: 0,
		           batchSize: 0,
		           dataset: "",
		           progressDuration: nil,
		           loaderFunc: nil,
		           dataProvider: nil,
		           dataSize: 0,
		           service: nil,
		           operation: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &loader{
				eg:               test.fields.eg,
				client:           test.fields.client,
				addr:             test.fields.addr,
				concurrency:      test.fields.concurrency,
				batchSize:        test.fields.batchSize,
				dataset:          test.fields.dataset,
				progressDuration: test.fields.progressDuration,
				loaderFunc:       test.fields.loaderFunc,
				dataProvider:     test.fields.dataProvider,
				dataSize:         test.fields.dataSize,
				service:          test.fields.service,
				operation:        test.fields.operation,
			}

			err := l.do(test.args.ctx, test.args.f, test.args.notify)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
