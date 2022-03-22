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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"testing"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		want Server
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Server, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Server, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_newLocations(t *testing.T) {
	t.Parallel()
	type args struct {
		uuids []string
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantLocs *payload.Object_Locations
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations) error {
		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuids: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           uuids: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotLocs := s.newLocations(test.args.uuids...)
			if err := checkFunc(test.want, gotLocs); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_newLocation(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		want *payload.Object_Location
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *payload.Object_Location) error {
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
		           uuid: "",
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           uuid: "",
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			got := s.newLocation(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Exists(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		uid *payload.Object_ID
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_ID
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_ID, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_ID, err error) error {
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
		           ctx: nil,
		           uid: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           uid: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.Exists(test.args.ctx, test.args.uid)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Search(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Search_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.Search(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SearchByID(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Search_IDRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_toSearchResponse(t *testing.T) {
	t.Parallel()
	type args struct {
		dists []model.Distance
		err   error
	}
	type want struct {
		wantRes *payload.Search_Response
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		           dists: nil,
		           err: nil,
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
		           dists: nil,
		           err: nil,
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

			gotRes, err := toSearchResponse(test.args.dists, test.args.err)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamSearch(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Search_StreamSearchServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamSearch(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamSearchByID(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Search_StreamSearchByIDServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamSearchByID(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiSearch(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiSearch(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiSearchByID(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiIDRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Insert(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *payload.Insert_Request
	}
	type fields struct {
		name string
		ip   string
		// ngt               service.NGT
		// eg                errgroup.Group
		streamConcurrency int

		svcCfg  *config.NGT
		svcOpts []service.Option
	}
	type want struct {
		wantRes *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*server)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err.Error(), w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}

	name := "vald-agent-ngt-1"         // agent name
	id := "uuid1"                      // insert request id
	ip := net.LoadLocalIP()            // agent ip address
	intVec := []float32{1, 2, 3}       // int vector of the insert request
	f32Vec := []float32{1.5, 2.3, 3.6} // float32 vector of the insert request
	intVecDim := 3                     // int vector dimension
	f32VecDim := 3                     // float32 vector dimension

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
		- Equivalence Class Testing
			- uint8, float32
				- case 1.1: Insert vector success (vector type is uint8)
				- case 1.2: Insert vector success (vector type is float32)
				- case 2.1: Insert vector with different dimension (vector type is uint8)
				- case 2.2: Insert vector with different dimension (vector type is float32)
				- case 3.1: Insert gaussian distributed vector success (vector type is uint8)
				- case 3.2: Insert gaussian distributed vector success (vector type is float32)
				- case 4.1: Insert uniform distributed vector success (vector type is uint8)
				- case 4.2: Insert uniform distributed vector success (vector type is float32)

		- Boundary Value Testing
			- uint8, float32
				- case 1.1: Insert vector with 0 value success (vector type is uint8)
				- case 1.1: Insert vector with 0 value success (vector type is float32)
				- case 2.1: Insert vector with min value success (vector type is uint8)
				- case 2.2: Insert vector with min value success (vector type is float32)
				- case 3.1: Insert vector with max value success (vector type is uint8)
				- case 3.2: Insert vector with max value success (vector type is float32)
				- case 4.1: Insert with empty UUID fail (vector type is uint8)
				- case 4.2: Insert with empty UUID fail (vector type is float32)

			- float32
				- case 5: Insert vector with NaN value fail (vector type is float32)

			- case 6: Insert nil insert request fail
				* IncompatibleDimensionSize error will be returned.
			- case 7: Insert nil vector fail
				* IncompatibleDimensionSize error will be returned.
			- case 8: Insert empty insert vector fail
				* IncompatibleDimensionSize error will be returned.

		- Decision Table Testing
			- duplicated ID, duplicated vector, duplicated ID & vector
				- case 1.1: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID)
					* AlreadyExists error will be returned.
				- case 1.2: Insert duplicated request success when SkipStrictExistCheck is false (duplicated vector)
				- case 1.3: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID & vector)
				- case 2.1: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID)
					* SkipStrictExistCheck flag is not used in agent handler, so the result is same as case 1.
				- case 2.2: Insert duplicated request success when SkipStrictExistCheck is true (duplicated vector)
				- case 2.3: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID & vector)
	*/
	tests := []test{
		// Equivalence Class Testing
		func() test {
			return test{
				name: "Equivalence Class Testing case 1.1: Insert vector success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: intVec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "Equivalence Class Testing case 1.2: Insert vector success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: f32Vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			vec := []float32{1, 2, 3, 4, 5, 6, 7}
			dim := 3
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vec,
				},
			}

			return test{
				name: "Equivalence Class Testing case 2.1: Insert vector with different dimension (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    dim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(vec), 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			vec := []float32{1.5, 2.3, 3.6, 4.5, 6.6, 7.7}
			dim := 3
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: vec,
				},
			}

			return test{
				name: "Equivalence Class Testing case 2.2: Insert vector with different dimension (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    dim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(len(vec), 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   id,
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			generator, _ := vector.Uint8VectorGenerator(vector.Gaussian)
			ivec := generator(1, 128)[0]
			vec := make([]float32, len(ivec))
			for i := range ivec {
				vec = append(vec, float32(i))
			}

			return test{
				name: "Equivalence Class Testing case 3.1: Insert gaussian distributed vector success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			generator, _ := vector.Float32VectorGenerator(vector.Gaussian)
			vec := generator(1, 128)[0]

			return test{
				name: "Equivalence Class Testing case 3.2: Insert gaussian distributed vector success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			generator, _ := vector.Uint8VectorGenerator(vector.Uniform)
			ivec := generator(1, 128)[0]
			vec := make([]float32, len(ivec))
			for i := range ivec {
				vec = append(vec, float32(i))
			}

			return test{
				name: "Equivalence Class Testing case 4.1: Insert uniform distributed vector success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			generator, _ := vector.Float32VectorGenerator(vector.Uniform)
			vec := generator(1, 128)[0]

			return test{
				name: "Equivalence Class Testing case 4.2: Insert uniform distributed vector success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),

		// Boundary Value Testing
		func() test {
			vec := []float32{0, 0, 0}

			return test{
				name: "Boundary Value Testing case 1.1: Insert vector with 0 value success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			vec := []float32{0, 0, 0}

			return test{
				name: "Boundary Value Testing case 1.2: Insert vector with 0 value success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			vec := []float32{math.MinInt, math.MinInt, math.MinInt}

			return test{
				name: "Boundary Value Testing case 2.1: Insert vector with min value success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			vec := []float32{math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32}

			return test{
				name: "Boundary Value Testing case 2.2: Insert vector with min value success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			vec := []float32{math.MaxInt, math.MaxInt, math.MaxInt}

			return test{
				name: "Boundary Value Testing case 3.1: Insert vector with max value success (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			vec := []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}

			return test{
				name: "Boundary Value Testing case 3.2: Insert vector with max value success (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			id := ""
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
			}

			return test{
				name: "Boundary Value Testing case 4.1: Insert with empty UUID fail (vector type is uint8)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrUUIDNotFound(0)
						err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", id), err,
							&errdetails.RequestInfo{
								RequestId:   req.GetVector().GetId(),
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "uuid",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			id := ""
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: f32Vec,
				},
			}

			return test{
				name: "Boundary Value Testing case 4.2: Insert with empty UUID fail (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrUUIDNotFound(0)
						err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", id), err,
							&errdetails.RequestInfo{
								RequestId:   req.GetVector().GetId(),
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "uuid",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			id := ""
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: f32Vec,
				},
			}

			return test{
				name: "Boundary Value Testing case 4.2: Insert with empty UUID fail (vector type is float32)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    f32VecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					err: func() error {
						err := errors.ErrUUIDNotFound(0)
						err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", id), err,
							&errdetails.RequestInfo{
								RequestId:   req.GetVector().GetId(),
								ServingData: errdetails.Serialize(req),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "uuid",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
						return err
					}(),
				},
			}
		}(),
		func() test {
			nan := float32(math.NaN())
			vec := []float32{nan, nan, nan}

			return test{
				name: "Boundary Value Testing case 5: Insert vector with NaN value fail (vector type is float32)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: vec,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    len(vec),
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			dim := 3

			return test{
				name: "Boundary Value Testing case 6: Insert nil insert request fail",
				args: args{
					ctx: ctx,
					req: nil,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    dim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					// IncompatibleDimensionSize error will be returned
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   "",
								ServingData: errdetails.Serialize(nil),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			id := "1"
			dim := 3

			return test{
				name: "Boundary Value Testing case 7: Insert nil vector fail",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: nil,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    dim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					// IncompatibleDimensionSize error will be returned
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   "1",
								ServingData: errdetails.Serialize(nil),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),
		func() test {
			id := "1"
			dim := 3

			return test{
				name: "Boundary Value Testing case 8: Insert empty insert vector fail",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: []float32{},
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    dim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Float.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				want: want{
					// IncompatibleDimensionSize error will be returned
					err: func() error {
						err := errors.ErrIncompatibleDimensionSize(0, 3)
						return status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
							err,
							&errdetails.RequestInfo{
								RequestId:   "1",
								ServingData: errdetails.Serialize(nil),
							},
							&errdetails.BadRequest{
								FieldViolations: []*errdetails.BadRequestFieldViolation{
									{
										Field:       "vector dimension size",
										Description: err.Error(),
									},
								},
							},
							&errdetails.ResourceInfo{
								ResourceType: ngtResourceType + "/ngt.Insert",
								ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
							})
					}(),
				},
			}
		}(),

		// Decision Table Testing
		func() test {
			vec2 := []float32{3, 2, 1} // used in beforeFunc
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: false,
				},
			}

			return test{
				name: "Decision Table Testing case 1.1: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, vec2)
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			id2 := "uuid2"             // use in beforeFunc
			vec2 := []float32{3, 2, 1} // use in beforeFunc

			return test{
				name: "Decision Table Testing case 1.2: Insert duplicated request success when SkipStrictExistCheck is false (duplicated vector)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: intVec,
						},
						Config: &payload.Insert_Config{
							SkipStrictExistCheck: false,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id2, vec2)
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: false,
				},
			}

			return test{
				name: "Decision Table Testing case 1.3: Insert duplicated request fail when SkipStrictExistCheck is false (duplicated ID & vector)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, intVec)
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			vec2 := []float32{3, 2, 1} // use in beforeFunc
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}

			return test{
				name: "Decision Table Testing case 2.1: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, vec2)
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
		func() test {
			id2 := "uuid2"

			return test{
				name: "Decision Table Testing case 2.2: Insert duplicated request success when SkipStrictExistCheck is true (duplicated vector)",
				args: args{
					ctx: ctx,
					req: &payload.Insert_Request{
						Vector: &payload.Object_Vector{
							Id:     id,
							Vector: intVec,
						},
						Config: &payload.Insert_Config{
							SkipStrictExistCheck: true,
						},
					},
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id2, intVec)
				},
				want: want{
					wantRes: &payload.Object_Location{
						Name: name,
						Uuid: id,
						Ips:  []string{ip},
					},
				},
			}
		}(),
		func() test {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     id,
					Vector: intVec,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
				},
			}

			return test{
				name: "Decision Table Testing case 2.3: Insert duplicated request fail when SkipStrictExistCheck is true (duplicated ID & vector)",
				args: args{
					ctx: ctx,
					req: req,
				},
				fields: fields{
					name: name,
					ip:   ip,
					svcCfg: &config.NGT{
						Dimension:    intVecDim,
						DistanceType: ngt.Angle.String(),
						ObjectType:   ngt.Uint8.String(),
						KVSDB:        &config.KVSDB{},
						VQueue:       &config.VQueue{},
					},
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
				},
				beforeFunc: func(s *server) {
					s.ngt.Insert(id, intVec)
				},
				want: want{
					err: status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", id), errors.ErrUUIDAlreadyExists(id),
						&errdetails.RequestInfo{
							RequestId:   req.GetVector().GetId(),
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: ngtResourceType + "/ngt.Insert",
							ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, name, ip),
						}),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts, service.WithErrGroup(eg))...)
			if err != nil {
				t.Error(err)
			}

			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               ngt,
				eg:                eg,
				streamConcurrency: test.fields.streamConcurrency,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(s)
			}

			gotRes, err := s.Insert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Insert_StreamInsertServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamInsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Insert_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Locations
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiInsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Update(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Update_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.Update(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Update_StreamUpdateServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamUpdate(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Update_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Locations
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Upsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Upsert_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpsert(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Upsert_StreamUpsertServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamUpsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Upsert_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Locations
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_Remove(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Remove_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.Remove(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamRemove(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Remove_StreamRemoveServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamRemove(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiRemove(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Remove_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Locations
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiRemove(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_GetObject(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  *payload.Object_VectorRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Vector
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Vector, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Vector, err error) error {
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
		           ctx: nil,
		           id: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           id: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.GetObject(test.args.ctx, test.args.id)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamGetObject(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Object_StreamGetObjectServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamGetObject(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_CreateIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		c   *payload.Control_CreateIndexRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
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
		           ctx: nil,
		           c: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           c: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.CreateIndex(test.args.ctx, test.args.c)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		in1 *payload.Empty
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
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
		           ctx: nil,
		           in1: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           in1: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.SaveIndex(test.args.ctx, test.args.in1)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_CreateAndSaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		c   *payload.Control_CreateIndexRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
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
		           ctx: nil,
		           c: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           c: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.CreateAndSaveIndex(test.args.ctx, test.args.c)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_IndexInfo(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		in1 *payload.Empty
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Info_Index_Count
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Index_Count, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Info_Index_Count, err error) error {
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
		           ctx: nil,
		           in1: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           in1: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.IndexInfo(test.args.ctx, test.args.in1)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_LinearSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_Request
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_LinearSearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_IDRequest
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamLinearSearch(t *testing.T) {
	type args struct {
		stream vald.Search_StreamLinearSearchServer
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			err := s.StreamLinearSearch(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamLinearSearchByID(t *testing.T) {
	type args struct {
		stream vald.Search_StreamLinearSearchByIDServer
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
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
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			err := s.StreamLinearSearchByID(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiLinearSearch(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiRequest
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
	}
	type want struct {
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.reqs)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiLinearSearchByID(t *testing.T) {
	type args struct {
		ctx  context.Context
		reqs *payload.Search_MultiIDRequest
	}
	type fields struct {
		name                     string
		ip                       string
		ngt                      service.NGT
		eg                       errgroup.Group
		streamConcurrency        int
		UnimplementedAgentServer agent.UnimplementedAgentServer
		UnimplementedValdServer  vald.UnimplementedValdServer
	}
	type want struct {
		wantRes *payload.Search_Responses
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Search_Responses, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
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
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           UnimplementedAgentServer: nil,
		           UnimplementedValdServer: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &server{
				name:                     test.fields.name,
				ip:                       test.fields.ip,
				ngt:                      test.fields.ngt,
				eg:                       test.fields.eg,
				streamConcurrency:        test.fields.streamConcurrency,
				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
			}

			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.reqs)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
