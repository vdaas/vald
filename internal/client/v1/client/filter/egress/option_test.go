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

package egress

import (
	"testing"
	"time"
)

func TestWithIdleTimeout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "unset keeps the default", dur: "", want: defaultIdleTimeout},
		{name: "valid duration is used", dur: "5s", want: 5 * time.Second},
		{name: "invalid duration keeps the default", dur: "invalid", want: defaultIdleTimeout},
		{name: "non-positive duration keeps the default", dur: "-1s", want: defaultIdleTimeout},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			c, err := New(WithAddrs("configured:8081"), WithClient(newConnTracker().mock()), WithIdleTimeout(tc.dur))
			if err != nil {
				tt.Fatal(err)
			}
			cl, ok := c.(*client)
			if !ok {
				tt.Fatalf("New returned %T, want *client", c)
			}
			if cl.idle != tc.want {
				tt.Errorf("idle = %v, want %v", cl.idle, tc.want)
			}
		})
	}
}

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
// func TestWithClient(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		cl grpc.Client
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
// 		           cl:nil,
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
// 		           cl:nil,
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
// 			   got := WithClient(test.args.cl)
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
// 			   got := WithClient(test.args.cl)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
