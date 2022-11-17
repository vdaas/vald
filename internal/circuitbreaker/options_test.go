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
package circuitbreaker

import (
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithBreakerOpts(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		opts []BreakerOption
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
		       },
		       want: want {
		           obj: new(T),
		       },
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
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithBreakerOpts(test.args.opts...)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithBreakerOpts(test.args.opts...)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithClosedErrorRate(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		f float32
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           f: 0,
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           f: 0,
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithClosedErrorRate(test.args.f)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithClosedErrorRate(test.args.f)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithClosedErrorTripper(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		tp Tripper
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           tp: nil,
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           tp: nil,
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithClosedErrorTripper(test.args.tp)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithClosedErrorTripper(test.args.tp)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithHalfOpenErrorRate(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		f float32
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           f: 0,
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           f: 0,
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithHalfOpenErrorRate(test.args.f)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithHalfOpenErrorRate(test.args.f)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithHalfOpenTripper(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		tp Tripper
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           tp: nil,
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           tp: nil,
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithHalfOpenTripper(test.args.tp)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithHalfOpenTripper(test.args.tp)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithMinSamples(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		min int64
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           min: 0,
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           min: 0,
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithMinSamples(test.args.min)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithMinSamples(test.args.min)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithOpenTimeout(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		timeout string
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           timeout: "",
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           timeout: "",
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithOpenTimeout(test.args.timeout)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithOpenTimeout(test.args.timeout)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithClosedRefreshTimeout(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		timeout string
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           timeout: "",
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           timeout: "",
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithClosedRefreshTimeout(test.args.timeout)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithClosedRefreshTimeout(test.args.timeout)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithHalfOpenErrorTripper(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		tp Tripper
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           tp: nil,
		       },
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           tp: nil,
		           },
		           want: want {
		               obj: new(T),
		           },
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

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithHalfOpenErrorTripper(test.args.tp)
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }
			   got := WithHalfOpenErrorTripper(test.args.tp)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}
