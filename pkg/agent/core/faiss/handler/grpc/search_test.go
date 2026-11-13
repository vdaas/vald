// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// Test_server_Search_NaNOrInf checks that the core rejection is mapped to InvalidArgument.
func Test_server_Search_NaNOrInf(t *testing.T) {
	t.Parallel()
	vec := []float32{0, 1}
	s := &server{faiss: &stubFaiss{err: errors.ErrNaNOrInfVector, dim: len(vec)}}
	_, err := s.Search(context.Background(), &payload.Search_Request{Vector: vec})
	wantInvalidArgument(t, err)
}

// Test_server_Search_rejectsNaNOrInfUpFront checks that a NaN / +-Inf query is rejected with
// InvalidArgument before it reaches the core, which would answer Aborted while indexing.
func Test_server_Search_rejectsNaNOrInfUpFront(t *testing.T) {
	t.Parallel()
	for name, b := range nanOrInfValues() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// the stub core would answer Aborted, as the real one does while creating indices.
			s := &server{faiss: &stubFaiss{err: errors.ErrCreateIndexingIsInProgress, dim: 2}}
			_, err := s.Search(context.Background(), &payload.Search_Request{Vector: []float32{0, b}})
			wantInvalidArgument(t, err)
		})
	}
}

// Test_server_Search_concurrentCoreError calls Search concurrently with a faiss core error.
// The errors.As target used to be a package-level variable, which is a data race under -race.
func Test_server_Search_concurrentCoreError(t *testing.T) {
	t.Parallel()

	const callers = 32
	tests := []struct {
		err     error
		name    string
		wantMsg string
	}{
		// the core returns FaissError by value (errors.NewFaissError), which the As target must match.
		{
			name:    "faiss core error takes the core error branch",
			err:     errors.NewFaissError("boom"),
			wantMsg: "Search API failed to process search request due to faiss core process returned error",
		},
		{
			name:    "other error takes the default branch",
			err:     errors.New("boom"),
			wantMsg: "Search API failed to process search request",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vec := []float32{0, 1}
			s := &server{faiss: &stubFaiss{err: test.err, dim: len(vec)}}
			req := &payload.Search_Request{Vector: vec}

			start := make(chan struct{})
			errs := make(chan error, callers)
			for range callers {
				go func() {
					<-start
					_, err := s.Search(context.Background(), req)
					errs <- err
				}()
			}
			close(start)

			timeout := time.After(10 * time.Second)
			for range callers {
				select {
				case err := <-errs:
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("got non-status error: %v", err)
						continue
					}
					if st.Code() != codes.Internal {
						t.Errorf("code = %v, want %v", st.Code(), codes.Internal)
					}
					if st.Message() != test.wantMsg {
						t.Errorf("message = %q, want %q", st.Message(), test.wantMsg)
					}
				case <-timeout:
					t.Fatal("timed out waiting for Search callers")
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_Search(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.Search(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_SearchByID(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchServer
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchByIDServer
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		faiss                    service.Faiss
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           faiss:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				faiss:                    test.fields.faiss,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
