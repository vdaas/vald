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
package grpc

import (
	"context"
	"reflect"
	"testing"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_LinearSearch(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
