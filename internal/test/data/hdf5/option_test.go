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

// Package hdf5 is load hdf5 file
package hdf5

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithNameByString(t *testing.T) {
	// Change interface type to the type of object you are testing
	type args struct {
		n string
	}
	type want struct {
		obj data
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want want, got data) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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

	defaultCheckFunc := func(w want, obj data) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set original",
			args: args{
				n: "original",
			},
			want: want{
				obj: data{
					name: Original,
				},
			},
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
		{
			name: "set fashion-mnist",
			args: args{
				n: "fashion-mnist",
			},
			want: want{
				obj: data{
					name: FashionMNIST784Euclidean,
				},
			},
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := WithNameByString(test.args.n)
			obj := new(data)
			_ = got(obj)
			if err := checkFunc(test.want, *obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = any
	type args struct {
		dn DatasetName
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           dn:nil,
		       },
		       want: want {
		           obj: new(T),
		       },
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           dn:nil,
		           },
		           want: want {
		               obj: new(T),
		           },
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithName(test.args.dn)
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
			   got := WithName(test.args.dn)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithFilePath(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = any
	type args struct {
		f string
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
		           f:"",
		       },
		       want: want {
		           obj: new(T),
		       },
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           f:"",
		           },
		           want: want {
		               obj: new(T),
		           },
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   checkFunc := test.checkFunc
			   if test.checkFunc == nil {
			       checkFunc = defaultCheckFunc
			   }

			   got := WithFilePath(test.args.f)
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
			   got := WithFilePath(test.args.f)
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

// NOT IMPLEMENTED BELOW
