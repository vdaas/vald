// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// NOT IMPLEMENTED BELOW
//
// func Test_server_exists(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		uuid string
// 	}
// 	type fields struct {
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 	}
// 	type want struct {
// 		wantId *payload.Object_ID
// 		err    error
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
// 	defaultCheckFunc := func(w want, gotId *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotId, w.wantId) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotId, w.wantId)
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
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 			}
//
// 			gotId, err := s.exists(test.args.ctx, test.args.uuid)
// 			if err := checkFunc(test.want, gotId, err); err != nil {
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
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 	}
// 	type want struct {
// 		wantId *payload.Object_ID
// 		err    error
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
// 	defaultCheckFunc := func(w want, gotId *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotId, w.wantId) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotId, w.wantId)
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 			}
//
// 			gotId, err := s.Exists(test.args.ctx, test.args.meta)
// 			if err := checkFunc(test.want, gotId, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_getObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		uuid string
// 	}
// 	type fields struct {
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
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
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 			}
//
// 			gotVec, err := s.getObject(test.args.ctx, test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
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
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
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
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
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
//
// func Test_server_StreamListObject(t *testing.T) {
// 	type args struct {
// 		req    *payload.Object_List_Request
// 		stream vald.Object_StreamListObjectServer
// 	}
// 	type fields struct {
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
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
// 		           req:nil,
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           req:nil,
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 			}
//
// 			err := s.StreamListObject(test.args.req, test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_server_GetTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_TimestampRequest
// 	}
// 	type fields struct {
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 		name                    string
// 		ip                      string
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 	}
// 	type want struct {
// 		wantTs *payload.Object_Timestamp
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Timestamp, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotTs *payload.Object_Timestamp, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotTs, w.wantTs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTs, w.wantTs)
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 		           gateway:nil,
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           name:"",
// 		           ip:"",
// 		           UnimplementedValdServer:nil,
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
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 			}
//
// 			gotTs, err := s.GetTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotTs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
