// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package grpc

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

// NOT IMPLEMENTED BELOW

func Test_server_aggregationSearch(t *testing.T) {
	type args struct {
		ctx  context.Context
		aggr Aggregator
		cfg  *payload.Search_Config
		f    func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)
	}
	type fields struct {
		eg                      errgroup.Group
		gateway                 service.Gateway
		timeout                 time.Duration
		replica                 int
		streamConcurrency       int
		multiConcurrency        int
		name                    string
		ip                      string
		UnimplementedValdServer vald.UnimplementedValdServer
	}
	type want struct {
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx:nil,
		           aggr:nil,
		           cfg:nil,
		           f:nil,
		       },
		       fields: fields {
		           eg:nil,
		           gateway:nil,
		           timeout:nil,
		           replica:0,
		           streamConcurrency:0,
		           multiConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServer:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           aggr:nil,
		           cfg:nil,
		           f:nil,
		           },
		           fields: fields {
		           eg:nil,
		           gateway:nil,
		           timeout:nil,
		           replica:0,
		           streamConcurrency:0,
		           multiConcurrency:0,
		           name:"",
		           ip:"",
		           UnimplementedValdServer:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			s := &server{
				eg:                      test.fields.eg,
				gateway:                 test.fields.gateway,
				timeout:                 test.fields.timeout,
				replica:                 test.fields.replica,
				streamConcurrency:       test.fields.streamConcurrency,
				multiConcurrency:        test.fields.multiConcurrency,
				name:                    test.fields.name,
				ip:                      test.fields.ip,
				UnimplementedValdServer: test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.aggregationSearch(test.args.ctx, test.args.aggr, test.args.cfg, test.args.f)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_newStd(t *testing.T) {
	type args struct {
		num     int
		replica int
	}
	type want struct {
		want Aggregator
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Aggregator) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Aggregator) error {
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
		           num:0,
		           replica:0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           num:0,
		           replica:0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := newStd(test.args.num, test.args.replica)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdStdAggr_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		num     int
		dch     chan DistPayload
		maxDist atomic.Value
		result  []*payload.Object_Distance
		cancel  context.CancelFunc
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           num:0,
		           dch:nil,
		           maxDist:nil,
		           result:nil,
		           cancel:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           num:0,
		           dch:nil,
		           maxDist:nil,
		           result:nil,
		           cancel:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdStdAggr{
				num:     test.fields.num,
				dch:     test.fields.dch,
				maxDist: test.fields.maxDist,
				result:  test.fields.result,
				cancel:  test.fields.cancel,
			}

			v.Start(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdStdAggr_Send(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *payload.Search_Response
	}
	type fields struct {
		num     int
		dch     chan DistPayload
		maxDist atomic.Value
		result  []*payload.Object_Distance
		cancel  context.CancelFunc
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           data:nil,
		       },
		       fields: fields {
		           num:0,
		           dch:nil,
		           maxDist:nil,
		           result:nil,
		           cancel:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           data:nil,
		           },
		           fields: fields {
		           num:0,
		           dch:nil,
		           maxDist:nil,
		           result:nil,
		           cancel:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdStdAggr{
				num:     test.fields.num,
				dch:     test.fields.dch,
				maxDist: test.fields.maxDist,
				result:  test.fields.result,
				cancel:  test.fields.cancel,
			}

			v.Send(test.args.ctx, test.args.data)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdStdAggr_Result(t *testing.T) {
	type fields struct {
		num     int
		dch     chan DistPayload
		maxDist atomic.Value
		result  []*payload.Object_Distance
		cancel  context.CancelFunc
	}
	type want struct {
		want *payload.Search_Response
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *payload.Search_Response) error {
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
		           num:0,
		           dch:nil,
		           maxDist:nil,
		           result:nil,
		           cancel:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           num:0,
		           dch:nil,
		           maxDist:nil,
		           result:nil,
		           cancel:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdStdAggr{
				num:     test.fields.num,
				dch:     test.fields.dch,
				maxDist: test.fields.maxDist,
				result:  test.fields.result,
				cancel:  test.fields.cancel,
			}

			got := v.Result()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_newPairingHeap(t *testing.T) {
	type args struct {
		num     int
		replica int
	}
	type want struct {
		want Aggregator
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Aggregator) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Aggregator) error {
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
		           num:0,
		           replica:0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           num:0,
		           replica:0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := newPairingHeap(test.args.num, test.args.replica)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdPairingHeapAggr_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		num    int
		ph     *PairingHeap
		result []*payload.Object_Distance
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           num:0,
		           ph:PairingHeap{},
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           num:0,
		           ph:PairingHeap{},
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdPairingHeapAggr{
				num:    test.fields.num,
				ph:     test.fields.ph,
				result: test.fields.result,
			}

			v.Start(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdPairingHeapAggr_Send(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *payload.Search_Response
	}
	type fields struct {
		num    int
		ph     *PairingHeap
		result []*payload.Object_Distance
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           data:nil,
		       },
		       fields: fields {
		           num:0,
		           ph:PairingHeap{},
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           data:nil,
		           },
		           fields: fields {
		           num:0,
		           ph:PairingHeap{},
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdPairingHeapAggr{
				num:    test.fields.num,
				ph:     test.fields.ph,
				result: test.fields.result,
			}

			v.Send(test.args.ctx, test.args.data)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdPairingHeapAggr_Result(t *testing.T) {
	type fields struct {
		num    int
		ph     *PairingHeap
		result []*payload.Object_Distance
	}
	type want struct {
		want *payload.Search_Response
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *payload.Search_Response) error {
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
		           num:0,
		           ph:PairingHeap{},
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           num:0,
		           ph:PairingHeap{},
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdPairingHeapAggr{
				num:    test.fields.num,
				ph:     test.fields.ph,
				result: test.fields.result,
			}

			got := v.Result()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_newSlice(t *testing.T) {
	type args struct {
		num     int
		replica int
	}
	type want struct {
		want Aggregator
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Aggregator) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Aggregator) error {
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
		           num:0,
		           replica:0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           num:0,
		           replica:0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := newSlice(test.args.num, test.args.replica)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdSliceAggr_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		num    int
		result []*DistPayload
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           num:0,
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           num:0,
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdSliceAggr{
				num:    test.fields.num,
				result: test.fields.result,
			}

			v.Start(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdSliceAggr_Send(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *payload.Search_Response
	}
	type fields struct {
		num    int
		result []*DistPayload
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           data:nil,
		       },
		       fields: fields {
		           num:0,
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           data:nil,
		           },
		           fields: fields {
		           num:0,
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdSliceAggr{
				num:    test.fields.num,
				result: test.fields.result,
			}

			v.Send(test.args.ctx, test.args.data)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdSliceAggr_Result(t *testing.T) {
	type fields struct {
		num    int
		result []*DistPayload
	}
	type want struct {
		wantRes *payload.Search_Response
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response) error {
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           num:0,
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           num:0,
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdSliceAggr{
				num:    test.fields.num,
				result: test.fields.result,
			}

			gotRes := v.Result()
			if err := checkFunc(test.want, gotRes); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_newPoolSlice(t *testing.T) {
	type args struct {
		num     int
		replica int
	}
	type want struct {
		want Aggregator
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Aggregator) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Aggregator) error {
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
		           num:0,
		           replica:0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           num:0,
		           replica:0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := newPoolSlice(test.args.num, test.args.replica)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdPoolSliceAggr_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		num    int
		result []*DistPayload
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		       },
		       fields: fields {
		           num:0,
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           num:0,
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdPoolSliceAggr{
				num:    test.fields.num,
				result: test.fields.result,
			}

			v.Start(test.args.ctx)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdPoolSliceAggr_Send(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *payload.Search_Response
	}
	type fields struct {
		num    int
		result []*DistPayload
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           ctx:nil,
		           data:nil,
		       },
		       fields: fields {
		           num:0,
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           data:nil,
		           },
		           fields: fields {
		           num:0,
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdPoolSliceAggr{
				num:    test.fields.num,
				result: test.fields.result,
			}

			v.Send(test.args.ctx, test.args.data)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_valdPoolSliceAggr_Result(t *testing.T) {
	type fields struct {
		num    int
		result []*DistPayload
	}
	type want struct {
		wantRes *payload.Search_Response
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Response) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response) error {
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           num:0,
		           result:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           num:0,
		           result:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			v := &valdPoolSliceAggr{
				num:    test.fields.num,
				result: test.fields.result,
			}

			gotRes := v.Result()
			if err := checkFunc(test.want, gotRes); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
