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
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"
)

func Test_insertRequestProvider(t *testing.T) {
	t.Parallel()
	type args struct {
		dataset   assets.Dataset
		batchSize int
	}
	type want struct {
		wantF    func() interface{}
		wantSize int
		err      error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() interface{}, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotF func() interface{}, gotSize int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotF, w.wantF) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotF, w.wantF)
		}
		if !reflect.DeepEqual(gotSize, w.wantSize) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSize, w.wantSize)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           dataset: nil,
		           batchSize: 0,
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
		           dataset: nil,
		           batchSize: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			gotF, gotSize, err := insertRequestProvider(test.args.dataset, test.args.batchSize)
			if err := test.checkFunc(test.want, gotF, gotSize, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_objectVectorProvider(t *testing.T) {
	t.Parallel()
	type args struct {
		dataset assets.Dataset
	}
	type want struct {
		want  func() interface{}
		want1 int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() interface{}, int) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() interface{}, got1 int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           dataset: nil,
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
		           dataset: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got, got1 := objectVectorProvider(test.args.dataset)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_objectVectorsProvider(t *testing.T) {
	t.Parallel()
	type args struct {
		dataset assets.Dataset
		n       int
	}
	type want struct {
		want  func() interface{}
		want1 int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() interface{}, int) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() interface{}, got1 int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           dataset: nil,
		           n: 0,
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
		           dataset: nil,
		           n: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got, got1 := objectVectorsProvider(test.args.dataset, test.args.n)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_agent(t *testing.T) {
	t.Parallel()
	type args struct {
		conn *grpc.ClientConn
	}
	type want struct {
		want inserter
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, inserter) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got inserter) error {
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
		           conn: nil,
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
		           conn: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := agent(test.args.conn)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_gateway(t *testing.T) {
	t.Parallel()
	type args struct {
		conn *grpc.ClientConn
	}
	type want struct {
		want inserter
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, inserter) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got inserter) error {
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
		           conn: nil,
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
		           conn: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := gateway(test.args.conn)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_insert(t *testing.T) {
	t.Parallel()
	type args struct {
		c func(*grpc.ClientConn) inserter
	}
	type want struct {
		want loadFunc
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, loadFunc) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got loadFunc) error {
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
		           c: nil,
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
		           c: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := insert(test.args.c)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bulkInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		c func(*grpc.ClientConn) inserter
	}
	type want struct {
		want loadFunc
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, loadFunc) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got loadFunc) error {
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
		           c: nil,
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
		           c: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := bulkInsert(test.args.c)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_newInsert(t *testing.T) {
	t.Parallel()
	type fields struct {
		eg               errgroup.Group
		client           grpc.Client
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
		wantF loadFunc
		err   error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, loadFunc, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotF loadFunc, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotF, w.wantF) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotF, w.wantF)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
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
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
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

			gotF, err := l.newInsert()
			if err := test.checkFunc(test.want, gotF, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_newStreamInsert(t *testing.T) {
	t.Parallel()
	type fields struct {
		eg               errgroup.Group
		client           grpc.Client
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
		wantF loadFunc
		err   error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, loadFunc, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotF loadFunc, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotF, w.wantF) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotF, w.wantF)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
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
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
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

			gotF, err := l.newStreamInsert()
			if err := test.checkFunc(test.want, gotF, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
