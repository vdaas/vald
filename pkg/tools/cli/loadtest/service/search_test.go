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

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
	"go.uber.org/goleak"
)

func Test_searchRequestProvider(t *testing.T) {
	type args struct {
		dataset assets.Dataset
	}
	type want struct {
		want  func() interface{}
		want1 int
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() interface{}, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() interface{}, got1 int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
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

			got, got1, err := searchRequestProvider(test.args.dataset)
			if err := test.checkFunc(test.want, got, got1, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_newSearch(t *testing.T) {
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
		want loadFunc
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, loadFunc, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got loadFunc, err error) error {
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

			got, err := l.newSearch()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_loader_newStreamSearch(t *testing.T) {
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
		want loadFunc
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, loadFunc, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got loadFunc, err error) error {
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

			got, err := l.newStreamSearch()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
