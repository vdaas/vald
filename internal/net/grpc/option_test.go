//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package grpc

// NOT IMPLEMENTED BELOW
//
// func TestWithAddrs(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		addrs []string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           addrs:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           addrs:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithAddrs(test.args.addrs...)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithAddrs(test.args.addrs...)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithHealthCheckDuration(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithHealthCheckDuration(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithHealthCheckDuration(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithConnectionPoolRebalanceDuration(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithConnectionPoolRebalanceDuration(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithConnectionPoolRebalanceDuration(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithResolveDNS(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		flg bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           flg:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           flg:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithResolveDNS(test.args.flg)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithResolveDNS(test.args.flg)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithEnableConnectionPoolRebalance(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		flg bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           flg:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           flg:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithEnableConnectionPoolRebalance(test.args.flg)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithEnableConnectionPoolRebalance(test.args.flg)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithConnectionPoolSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithConnectionPoolSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithConnectionPoolSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithBackoffMaxDelay(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithBackoffMaxDelay(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithBackoffMaxDelay(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithBackoffBaseDelay(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithBackoffBaseDelay(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithBackoffBaseDelay(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithBackoffMultiplier(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		m float64
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           m:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           m:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithBackoffMultiplier(test.args.m)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithBackoffMultiplier(test.args.m)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithBackoffJitter(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		j float64
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           j:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           j:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithBackoffJitter(test.args.j)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithBackoffJitter(test.args.j)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMinConnectTimeout(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMinConnectTimeout(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMinConnectTimeout(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithErrGroup(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		eg errgroup.Group
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           eg:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           eg:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithErrGroup(test.args.eg)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithErrGroup(test.args.eg)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithBackoff(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		bo backoff.Backoff
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           bo:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           bo:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithBackoff(test.args.bo)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithBackoff(test.args.bo)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithCircuitBreaker(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		cb circuitbreaker.CircuitBreaker
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cb:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           cb:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithCircuitBreaker(test.args.cb)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithCircuitBreaker(test.args.cb)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithCallOptions(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		opts []grpc.CallOption
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithCallOptions(test.args.opts...)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithCallOptions(test.args.opts...)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithCallContentSubtype(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		contentSubtype string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           contentSubtype:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           contentSubtype:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithCallContentSubtype(test.args.contentSubtype)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithCallContentSubtype(test.args.contentSubtype)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMaxRecvMsgSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMaxRecvMsgSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMaxRecvMsgSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMaxSendMsgSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMaxSendMsgSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMaxSendMsgSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMaxRetryRPCBufferSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMaxRetryRPCBufferSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMaxRetryRPCBufferSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithWaitForReady(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		flg bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           flg:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           flg:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithWaitForReady(test.args.flg)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithWaitForReady(test.args.flg)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithDialOptions(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		opts []grpc.DialOption
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithDialOptions(test.args.opts...)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithDialOptions(test.args.opts...)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithWriteBufferSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithWriteBufferSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithWriteBufferSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithReadBufferSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithReadBufferSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithReadBufferSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithInitialWindowSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int32
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithInitialWindowSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithInitialWindowSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithInitialConnectionWindowSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int32
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithInitialConnectionWindowSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithInitialConnectionWindowSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMaxMsgSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMaxMsgSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMaxMsgSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithInsecure(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		flg bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           flg:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           flg:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithInsecure(test.args.flg)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithInsecure(test.args.flg)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithKeepaliveParams(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		t                   string
// 		to                  string
// 		permitWithoutStream bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           t:"",
// 		           to:"",
// 		           permitWithoutStream:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           t:"",
// 		           to:"",
// 		           permitWithoutStream:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithKeepaliveParams(test.args.t, test.args.to, test.args.permitWithoutStream)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithKeepaliveParams(test.args.t, test.args.to, test.args.permitWithoutStream)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithDialer(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		network string
// 		der     net.Dialer
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           network:"",
// 		           der:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           network:"",
// 		           der:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithDialer(test.args.network, test.args.der)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithDialer(test.args.network, test.args.der)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithTLSConfig(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		cfg *tls.Config
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           cfg:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithTLSConfig(test.args.cfg)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithTLSConfig(test.args.cfg)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithAuthority(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		a string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           a:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           a:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithAuthority(test.args.a)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithAuthority(test.args.a)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithDisableRetry(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		disable bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           disable:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           disable:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithDisableRetry(test.args.disable)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithDisableRetry(test.args.disable)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithIdleTimeout(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithIdleTimeout(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithIdleTimeout(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMaxCallAttempts(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		n int
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           n:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           n:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMaxCallAttempts(test.args.n)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMaxCallAttempts(test.args.n)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithMaxHeaderListSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size uint32
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           size:0,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           size:0,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithMaxHeaderListSize(test.args.size)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithMaxHeaderListSize(test.args.size)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithSharedWriteBuffer(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		enable bool
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           enable:false,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           enable:false,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithSharedWriteBuffer(test.args.enable)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithSharedWriteBuffer(test.args.enable)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithUserAgent(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		ua string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ua:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           ua:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithUserAgent(test.args.ua)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithUserAgent(test.args.ua)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithClientInterceptors(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		names []string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           names:nil,
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           names:nil,
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithClientInterceptors(test.args.names...)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithClientInterceptors(test.args.names...)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
//
// func TestWithOldConnCloseDelay(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dur string
// 	}
// 	type want struct {
// 		obj *T
// 		// Uncomment this line if the option returns an error, otherwise delete it
// 		// err error
// 	}
// 	type test struct {
// 		name string
// 		args args
// 		want want
// 		// Use the first line if the option returns an error. otherwise use the second line
// 		// checkFunc  func(want, *T, error) error
// 		// checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
//
// 	// Uncomment this block if the option returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T, err error) error {
// 	       if !errors.Is(err, w.err) {
// 	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 	       }
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	// Uncomment this block if the option do not returns an error, otherwise delete it
// 	/*
// 	   defaultCheckFunc := func(w want, obj *T) error {
// 	       if !reflect.DeepEqual(obj, w.obj) {
// 	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
// 	       }
// 	       return nil
// 	   }
// 	*/
//
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           dur:"",
// 		       },
// 		       want: want {
// 		           obj: new(T),
// 		       },
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
// 		           dur:"",
// 		           },
// 		           want: want {
// 		               obj: new(T),
// 		           },
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
//
// 			// Uncomment this block if the option returns an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
//
// 			   got := WithOldConnCloseDelay(test.args.dur)
// 			   obj := new(T)
// 			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
//
// 			// Uncomment this block if the option do not return an error, otherwise delete it
// 			/*
// 			   checkFunc := test.checkFunc
// 			   if test.checkFunc == nil {
// 			       checkFunc = defaultCheckFunc
// 			   }
// 			   got := WithOldConnCloseDelay(test.args.dur)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
