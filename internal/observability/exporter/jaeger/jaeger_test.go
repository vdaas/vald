//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantJ Jaeger
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Jaeger, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotJ Jaeger, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotJ, w.wantJ) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotJ, w.wantJ)
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

			gotJ, err := New(test.args.opts...)
			if err := checkFunc(test.want, gotJ, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_export_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		tp                  *trace.TracerProvider
		exp                 *jaeger.Exporter
		collectorEndpoint   string
		client              *http.Client
		collectorPassword   string
		collectorUserName   string
		agentHost           string
		agentPort           string
		agentReconnInterval time.Duration
		agentMaxPacketSize  int
		serviceName         string
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
		       },
		       fields: fields {
		           tp: nil,
		           exp: nil,
		           collectorEndpoint: "",
		           client: nil,
		           collectorPassword: "",
		           collectorUserName: "",
		           agentHost: "",
		           agentPort: "",
		           agentReconnInterval: nil,
		           agentMaxPacketSize: 0,
		           serviceName: "",
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
		           tp: nil,
		           exp: nil,
		           collectorEndpoint: "",
		           client: nil,
		           collectorPassword: "",
		           collectorUserName: "",
		           agentHost: "",
		           agentPort: "",
		           agentReconnInterval: nil,
		           agentMaxPacketSize: 0,
		           serviceName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			e := &export{
				tp:                  test.fields.tp,
				exp:                 test.fields.exp,
				collectorEndpoint:   test.fields.collectorEndpoint,
				client:              test.fields.client,
				collectorPassword:   test.fields.collectorPassword,
				collectorUserName:   test.fields.collectorUserName,
				agentHost:           test.fields.agentHost,
				agentPort:           test.fields.agentPort,
				agentReconnInterval: test.fields.agentReconnInterval,
				agentMaxPacketSize:  test.fields.agentMaxPacketSize,
				serviceName:         test.fields.serviceName,
			}

			err := e.Start(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_export_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		tp                  *trace.TracerProvider
		exp                 *jaeger.Exporter
		collectorEndpoint   string
		client              *http.Client
		collectorPassword   string
		collectorUserName   string
		agentHost           string
		agentPort           string
		agentReconnInterval time.Duration
		agentMaxPacketSize  int
		serviceName         string
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           tp: nil,
		           exp: nil,
		           collectorEndpoint: "",
		           client: nil,
		           collectorPassword: "",
		           collectorUserName: "",
		           agentHost: "",
		           agentPort: "",
		           agentReconnInterval: nil,
		           agentMaxPacketSize: 0,
		           serviceName: "",
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
		           tp: nil,
		           exp: nil,
		           collectorEndpoint: "",
		           client: nil,
		           collectorPassword: "",
		           collectorUserName: "",
		           agentHost: "",
		           agentPort: "",
		           agentReconnInterval: nil,
		           agentMaxPacketSize: 0,
		           serviceName: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			e := &export{
				tp:                  test.fields.tp,
				exp:                 test.fields.exp,
				collectorEndpoint:   test.fields.collectorEndpoint,
				client:              test.fields.client,
				collectorPassword:   test.fields.collectorPassword,
				collectorUserName:   test.fields.collectorUserName,
				agentHost:           test.fields.agentHost,
				agentPort:           test.fields.agentPort,
				agentReconnInterval: test.fields.agentReconnInterval,
				agentMaxPacketSize:  test.fields.agentMaxPacketSize,
				serviceName:         test.fields.serviceName,
			}

			e.Stop(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
