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

package params

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithConfigFilePathKeys(t *testing.T) {
	type T = parser
	type args struct {
		keys []string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when keys is `key1`",
			args: args{
				keys: []string{
					"key1",
				},
			},
			want: want{
				obj: &T{
					filePath: struct {
						keys        []string
						defaultPath string
						description string
					}{
						keys: []string{
							"key1",
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithConfigFilePathKeys(test.args.keys...)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithConfigFilePathDefault(t *testing.T) {
	type T = parser
	type args struct {
		path string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when path is `path`",
			args: args{
				path: "path",
			},
			want: want{
				obj: &T{
					filePath: struct {
						keys        []string
						defaultPath string
						description string
					}{
						defaultPath: "path",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithConfigFilePathDefault(test.args.path)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithConfigFileDescription(t *testing.T) {
	type T = parser
	type args struct {
		desc string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when desc is `desc`",
			args: args{
				desc: "desc",
			},
			want: want{
				obj: &T{
					filePath: struct {
						keys        []string
						defaultPath string
						description string
					}{
						description: "desc",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithConfigFileDescription(test.args.desc)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithVersionKeys(t *testing.T) {
	type T = parser
	type args struct {
		keys []string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when keys is `key1`",
			args: args{
				keys: []string{
					"key1",
				},
			},
			want: want{
				obj: &T{
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						keys: []string{
							"key1",
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithVersionKeys(test.args.keys...)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithVersionFlagDefault(t *testing.T) {
	type T = parser
	type args struct {
		flag bool
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when flag is true",
			args: args{
				flag: true,
			},
			want: want{
				obj: &T{
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						defaultFlag: true,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithVersionFlagDefault(test.args.flag)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithVersionDescription(t *testing.T) {
	type T = parser
	type args struct {
		desc string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when desc is true",
			args: args{
				desc: "desc",
			},
			want: want{
				obj: &T{
					version: struct {
						keys        []string
						defaultFlag bool
						description string
					}{
						description: "desc",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithVersionDescription(test.args.desc)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestWithName(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		name string
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
// 		           name:"",
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
// 		           name:"",
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
// 			   got := WithName(test.args.name)
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
// 			   got := WithName(test.args.name)
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
// func TestWithOverrideDefault(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		flag bool
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
// 		           flag:false,
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
// 		           flag:false,
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
// 			   got := WithOverrideDefault(test.args.flag)
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
// 			   got := WithOverrideDefault(test.args.flag)
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
// func TestWithArgumentFilters(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		filters []func(string) bool
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
// 		           filters:nil,
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
// 		           filters:nil,
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
// 			   got := WithArgumentFilters(test.args.filters...)
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
// 			   got := WithArgumentFilters(test.args.filters...)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
