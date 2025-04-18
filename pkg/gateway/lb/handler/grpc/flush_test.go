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
// func Test_server_Flush(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Flush_Request
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
// 		wantCnts *payload.Info_Index_Count
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Index_Count, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotCnts *payload.Info_Index_Count, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotCnts, w.wantCnts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotCnts, w.wantCnts)
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
// 			gotCnts, err := s.Flush(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotCnts, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
