//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package pool

// NOT IMPLEMENTED BELOW
//
// func TestWithAddr(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		addr string
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
// 		           addr:"",
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
// 		           addr:"",
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
// 			   got := WithAddr(test.args.addr)
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
// 			   got := WithAddr(test.args.addr)
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
// func TestWithHost(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		host string
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
// 		           host:"",
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
// 		           host:"",
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
// 			   got := WithHost(test.args.host)
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
// 			   got := WithHost(test.args.host)
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
// func TestWithPort(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		port int
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
// 		           port:0,
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
// 		           port:0,
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
// 			   got := WithPort(test.args.port)
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
// 			   got := WithPort(test.args.port)
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
// func TestWithStartPort(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		port int
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
// 		           port:0,
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
// 		           port:0,
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
// 			   got := WithStartPort(test.args.port)
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
// 			   got := WithStartPort(test.args.port)
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
// func TestWithEndPort(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		port int
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
// 		           port:0,
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
// 		           port:0,
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
// 			   got := WithEndPort(test.args.port)
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
// 			   got := WithEndPort(test.args.port)
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
// 			   got := WithResolveDNS(test.args.enable)
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
// 			   got := WithResolveDNS(test.args.enable)
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
// func TestWithSize(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		size uint64
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
// 			   got := WithSize(test.args.size)
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
// 			   got := WithSize(test.args.size)
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
// 		opts []DialOption
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
// func TestWithDialTimeout(t *testing.T) {
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
// 			   got := WithDialTimeout(test.args.dur)
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
// 			   got := WithDialTimeout(test.args.dur)
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
