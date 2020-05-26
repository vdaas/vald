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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"reflect"
	"testing"
	"time"

	agent "github.com/vdaas/vald/apis/grpc/agent/core"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/gateway/meta/service"

	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want vald.ValdServer
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, vald.ValdServer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got vald.ValdServer) error {
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

			got := New(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_Exists(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *payload.Object_ID
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
		streamConcurrency int
	}
	type want struct {
		want *payload.Object_ID
		err  error
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
	defaultCheckFunc := func(w want, got *payload.Object_ID, err error) error {
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
		           ctx: nil,
		           meta: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           meta: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			got, err := s.Exists(test.args.ctx, test.args.meta)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_Search(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_Request
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.Search(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_SearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_IDRequest
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_search(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *payload.Search_Config
		f   func(ctx context.Context, ac agent.AgentClient, copts ...grpc.CallOption) (*payload.Search_Response, error)
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           cfg: nil,
		           f: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           cfg: nil,
		           f: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.search(test.args.ctx, test.args.cfg, test.args.f)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamSearch(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamSearchServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamSearch(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamSearchByID(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamSearchByIDServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamSearchByID(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Object_Vector
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
		streamConcurrency int
	}
	type want struct {
		wantCe *payload.Empty
		err    error
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
	defaultCheckFunc := func(w want, gotCe *payload.Empty, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotCe, w.wantCe) {
			return errors.Errorf("got = %v, want %v", gotCe, w.wantCe)
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
		           vec: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           vec: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotCe, err := s.Insert(test.args.ctx, test.args.vec)
			if err := test.checkFunc(test.want, gotCe, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamInsert(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamInsertServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamInsert(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_MultiInsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Object_Vectors
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           vecs: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           vecs: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiInsert(test.args.ctx, test.args.vecs)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Object_Vector
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           vec: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           vec: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.Update(test.args.ctx, test.args.vec)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamUpdate(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamUpdateServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamUpdate(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_MultiUpdate(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Object_Vectors
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           vecs: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           vecs: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiUpdate(test.args.ctx, test.args.vecs)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_Upsert(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Object_Vector
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
		streamConcurrency int
	}
	type want struct {
		want *payload.Empty
		err  error
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
	defaultCheckFunc := func(w want, got *payload.Empty, err error) error {
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
		           ctx: nil,
		           vec: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           vec: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			got, err := s.Upsert(test.args.ctx, test.args.vec)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamUpsert(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamUpsertServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamUpsert(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_MultiUpsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Object_Vectors
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
		streamConcurrency int
	}
	type want struct {
		want *payload.Empty
		err  error
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
	defaultCheckFunc := func(w want, got *payload.Empty, err error) error {
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
		           ctx: nil,
		           vecs: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           vecs: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			got, err := s.MultiUpsert(test.args.ctx, test.args.vecs)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		id  *payload.Object_ID
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
		streamConcurrency int
	}
	type want struct {
		want *payload.Empty
		err  error
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
	defaultCheckFunc := func(w want, got *payload.Empty, err error) error {
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
		           ctx: nil,
		           id: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			got, err := s.Remove(test.args.ctx, test.args.id)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamRemove(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamRemoveServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamRemove(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_MultiRemove(t *testing.T) {
	type args struct {
		ctx context.Context
		ids *payload.Object_IDs
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
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
		           ids: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           ids: nil,
		           },
		           fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiRemove(test.args.ctx, test.args.ids)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		id  *payload.Object_ID
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
		streamConcurrency int
	}
	type want struct {
		wantVec *payload.Backup_MetaVector
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Backup_MetaVector, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotVec *payload.Backup_MetaVector, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got = %v, want %v", gotVec, w.wantVec)
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotVec, err := s.GetObject(test.args.ctx, test.args.id)
			if err := test.checkFunc(test.want, gotVec, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_server_StreamGetObject(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamGetObjectServer
	}
	type fields struct {
		eg                errgroup.Group
		gateway           service.Gateway
		metadata          service.Meta
		backup            service.Backup
		timeout           time.Duration
		filter            service.Filter
		replica           int
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
		           stream: nil,
		       },
		       fields: fields {
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
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
		           eg: nil,
		           gateway: nil,
		           metadata: nil,
		           backup: nil,
		           timeout: nil,
		           filter: nil,
		           replica: 0,
		           streamConcurrency: 0,
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
			s := &server{
				eg:                test.fields.eg,
				gateway:           test.fields.gateway,
				metadata:          test.fields.metadata,
				backup:            test.fields.backup,
				timeout:           test.fields.timeout,
				filter:            test.fields.filter,
				replica:           test.fields.replica,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamGetObject(test.args.stream)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
