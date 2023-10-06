//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want vald.ServerWithFilter
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, vald.ServerWithFilter) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got vald.ServerWithFilter) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
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
// 		           opts:nil,
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
//
// 			got := New(test.args.opts...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_SearchObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_ObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.SearchObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiSearchObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.MultiSearchObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamSearchObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamSearchObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamSearchObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_LinearSearchObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_ObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		want *payload.Search_Response
// 		err  error
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
// 	defaultCheckFunc := func(w want, got *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			got, err := s.LinearSearchObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearchObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.MultiLinearSearchObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearchObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamSearchObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamLinearSearchObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_InsertObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_ObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			got, err := s.InsertObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamInsertObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamInsertObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamInsertObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiInsertObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiInsertObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_UpdateObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_ObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			got, err := s.UpdateObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamUpdateObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamUpdateObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamUpdateObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiUpdateObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiUpdateObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_UpsertObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Upsert_ObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		want *payload.Object_Location
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			got, err := s.UpsertObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamUpsertObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Filter_StreamUpsertObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamUpsertObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiUpsertObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Upsert_MultiObjectRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiUpsertObject(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_Exists(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		meta *payload.Object_ID
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		want *payload.Object_ID
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_ID, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
// 		           meta:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           meta:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			got, err := s.Exists(test.args.ctx, test.args.meta)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_Search(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.Search(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.SearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchByIDServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.MultiSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
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
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_LinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_LinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamLinearSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchByIDServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamLinearSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_Insert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLoc, err := s.Insert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamInsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Insert_StreamInsertServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamInsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiInsert(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiInsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_Update(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLoc, err := s.Update(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamUpdate(t *testing.T) {
// 	type args struct {
// 		stream vald.Update_StreamUpdateServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamUpdate(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiUpdate(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_Upsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Upsert_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamUpsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Upsert_StreamUpsertServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamUpsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiUpsert(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Upsert_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_Remove(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_Request
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLoc *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLoc, err := s.Remove(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLoc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamRemove(t *testing.T) {
// 	type args struct {
// 		stream vald.Remove_StreamRemoveServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamRemove(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_MultiRemove(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Remove_MultiRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotLocs, err := s.MultiRemove(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_GetObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_VectorRequest
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
// 	}
// 	type want struct {
// 		wantVec *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			gotVec, err := s.GetObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_StreamGetObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Object_StreamGetObjectServer
// 	}
// 	type fields struct {
// 		eg                                errgroup.Group
// 		defaultVectorizer                 string
// 		defaultFilters                    []string
// 		name                              string
// 		ip                                string
// 		ingress                           ingress.Client
// 		egress                            egress.Client
// 		gateway                           client.Client
// 		copts                             []grpc.CallOption
// 		streamConcurrency                 int
// 		Vectorizer                        string
// 		DistanceFilters                   []*config.DistanceFilterConfig
// 		ObjectFilters                     []*config.ObjectFilterConfig
// 		SearchFilters                     []string
// 		InsertFilters                     []string
// 		UpdateFilters                     []string
// 		UpsertFilters                     []string
// 		UnimplementedValdServerWithFilter vald.UnimplementedValdServerWithFilter
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 		           eg:nil,
// 		           defaultVectorizer:"",
// 		           defaultFilters:nil,
// 		           name:"",
// 		           ip:"",
// 		           ingress:nil,
// 		           egress:nil,
// 		           gateway:nil,
// 		           copts:nil,
// 		           streamConcurrency:0,
// 		           Vectorizer:"",
// 		           DistanceFilters:nil,
// 		           ObjectFilters:nil,
// 		           SearchFilters:nil,
// 		           InsertFilters:nil,
// 		           UpdateFilters:nil,
// 		           UpsertFilters:nil,
// 		           UnimplementedValdServerWithFilter:nil,
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
// 				eg:                                test.fields.eg,
// 				defaultVectorizer:                 test.fields.defaultVectorizer,
// 				defaultFilters:                    test.fields.defaultFilters,
// 				name:                              test.fields.name,
// 				ip:                                test.fields.ip,
// 				ingress:                           test.fields.ingress,
// 				egress:                            test.fields.egress,
// 				gateway:                           test.fields.gateway,
// 				copts:                             test.fields.copts,
// 				streamConcurrency:                 test.fields.streamConcurrency,
// 				Vectorizer:                        test.fields.Vectorizer,
// 				DistanceFilters:                   test.fields.DistanceFilters,
// 				ObjectFilters:                     test.fields.ObjectFilters,
// 				SearchFilters:                     test.fields.SearchFilters,
// 				InsertFilters:                     test.fields.InsertFilters,
// 				UpdateFilters:                     test.fields.UpdateFilters,
// 				UpsertFilters:                     test.fields.UpsertFilters,
// 				UnimplementedValdServerWithFilter: test.fields.UnimplementedValdServerWithFilter,
// 			}
//
// 			err := s.StreamGetObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
