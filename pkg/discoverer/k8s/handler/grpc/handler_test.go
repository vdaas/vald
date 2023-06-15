//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/singleflight"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

// NOT IMPLEMENTED BELOW

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantDs DiscovererServer
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, DiscovererServer, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotDs DiscovererServer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotDs, w.wantDs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDs, w.wantDs)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts:nil,
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
		           opts:nil,
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

			gotDs, err := New(test.args.opts...)
			if err := checkFunc(test.want, gotDs, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Start(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	type fields struct {
		dsc                           service.Discoverer
		pgroup                        singleflight.Group[*payload.Info_Pods]
		ngroup                        singleflight.Group[*payload.Info_Nodes]
		ip                            string
		name                          string
		UnimplementedDiscovererServer discoverer.UnimplementedDiscovererServer
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
		           in0:nil,
		       },
		       fields: fields {
		           dsc:nil,
		           pgroup:nil,
		           ngroup:nil,
		           ip:"",
		           name:"",
		           UnimplementedDiscovererServer:nil,
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
		           in0:nil,
		           },
		           fields: fields {
		           dsc:nil,
		           pgroup:nil,
		           ngroup:nil,
		           ip:"",
		           name:"",
		           UnimplementedDiscovererServer:nil,
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
				dsc:                           test.fields.dsc,
				pgroup:                        test.fields.pgroup,
				ngroup:                        test.fields.ngroup,
				ip:                            test.fields.ip,
				name:                          test.fields.name,
				UnimplementedDiscovererServer: test.fields.UnimplementedDiscovererServer,
			}

			s.Start(test.args.in0)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Pods(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Discoverer_Request
	}
	type fields struct {
		dsc                           service.Discoverer
		pgroup                        singleflight.Group[*payload.Info_Pods]
		ngroup                        singleflight.Group[*payload.Info_Nodes]
		ip                            string
		name                          string
		UnimplementedDiscovererServer discoverer.UnimplementedDiscovererServer
	}
	type want struct {
		want *payload.Info_Pods
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Pods, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *payload.Info_Pods, err error) error {
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
		           ctx:nil,
		           req:nil,
		       },
		       fields: fields {
		           dsc:nil,
		           pgroup:nil,
		           ngroup:nil,
		           ip:"",
		           name:"",
		           UnimplementedDiscovererServer:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           dsc:nil,
		           pgroup:nil,
		           ngroup:nil,
		           ip:"",
		           name:"",
		           UnimplementedDiscovererServer:nil,
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
				dsc:                           test.fields.dsc,
				pgroup:                        test.fields.pgroup,
				ngroup:                        test.fields.ngroup,
				ip:                            test.fields.ip,
				name:                          test.fields.name,
				UnimplementedDiscovererServer: test.fields.UnimplementedDiscovererServer,
			}

			got, err := s.Pods(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Nodes(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Discoverer_Request
	}
	type fields struct {
		dsc                           service.Discoverer
		pgroup                        singleflight.Group[*payload.Info_Pods]
		ngroup                        singleflight.Group[*payload.Info_Nodes]
		ip                            string
		name                          string
		UnimplementedDiscovererServer discoverer.UnimplementedDiscovererServer
	}
	type want struct {
		want *payload.Info_Nodes
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Nodes, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *payload.Info_Nodes, err error) error {
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
		           ctx:nil,
		           req:nil,
		       },
		       fields: fields {
		           dsc:nil,
		           pgroup:nil,
		           ngroup:nil,
		           ip:"",
		           name:"",
		           UnimplementedDiscovererServer:nil,
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
		           req:nil,
		           },
		           fields: fields {
		           dsc:nil,
		           pgroup:nil,
		           ngroup:nil,
		           ip:"",
		           name:"",
		           UnimplementedDiscovererServer:nil,
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
				dsc:                           test.fields.dsc,
				pgroup:                        test.fields.pgroup,
				ngroup:                        test.fields.ngroup,
				ip:                            test.fields.ip,
				name:                          test.fields.name,
				UnimplementedDiscovererServer: test.fields.UnimplementedDiscovererServer,
			}

			got, err := s.Nodes(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_singleflightKey(t *testing.T) {
	type args struct {
		pref string
		req  *payload.Discoverer_Request
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           pref:"",
		           req:nil,
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
		           pref:"",
		           req:nil,
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

			got := singleflightKey(test.args.pref, test.args.req)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
