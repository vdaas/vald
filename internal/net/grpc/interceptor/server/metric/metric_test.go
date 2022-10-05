// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package metric

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMetricInterceptors(t *testing.T) {
	type want struct {
		want  grpc.UnaryServerInterceptor
		want1 grpc.StreamServerInterceptor
		err   error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got grpc.UnaryServerInterceptor, got1 grpc.StreamServerInterceptor, err error) error {
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
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, got1, err := MetricInterceptors()
			if err := checkFunc(test.want, got, got1, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_attributesFromError(t *testing.T) {
	type T = []attribute.KeyValue
	type args struct {
		method string
		err    error
	}
	type want struct {
		obj T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(w want, obj T) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}

	defaultCheckFunc := func(w want, obj T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "return []attribute.KeyValue when err is nil",
			args: args{
				method: "InsertRPC",
				err:    nil,
			},
			want: want{
				obj: []attribute.KeyValue{
					attribute.String(gRPCMethodKeyName, "InsertRPC"),
					attribute.String(gRPCStatus, codes.OK.String()),
				},
			},
		},
		{
			name: "return []attribute.KeyValue when err is not nil",
			args: args{
				method: "InsertRPC",
				err:    context.DeadlineExceeded,
			},
			want: want{
				obj: []attribute.KeyValue{
					attribute.String(gRPCMethodKeyName, "InsertRPC"),
					attribute.String(gRPCStatus, codes.DeadlineExceeded.String()),
				},
			},
		},
		{
			name: "return []attribute.KeyValue when err type is wrapped error",
			args: args{
				method: "InsertRPC",
				err:    status.WrapWithInvalidArgument("Insert API failed", errors.ErrIncompatibleDimensionSize(100, 940)),
			},
			want: want{
				obj: []attribute.KeyValue{
					attribute.String(gRPCMethodKeyName, "InsertRPC"),
					attribute.String(gRPCStatus, codes.InvalidArgument.String()),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := attributesFromError(test.args.method, test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
