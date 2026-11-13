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

package observer

import (
	"testing"
	"time"
)

// time.NewTicker panics for a non-positive interval, and both values feed a ticker,
// so non-positive input must leave the previously configured (default) value untouched.
func TestObserverDurationOptions_NonPositive(t *testing.T) {
	t.Parallel()
	const initial = 7 * time.Second
	tests := []struct {
		name  string
		opt   func(string) Option
		field func(*observer) *time.Duration
		dur   string
		want  time.Duration
	}{
		{
			name:  "backup duration zero is ignored",
			opt:   WithBackupDuration,
			field: func(o *observer) *time.Duration { return &o.checkDuration },
			dur:   "0s",
			want:  initial,
		},
		{
			name:  "backup duration negative is ignored",
			opt:   WithBackupDuration,
			field: func(o *observer) *time.Duration { return &o.checkDuration },
			dur:   "-1s",
			want:  initial,
		},
		{
			name:  "backup duration positive is applied",
			opt:   WithBackupDuration,
			field: func(o *observer) *time.Duration { return &o.checkDuration },
			dur:   "2s",
			want:  2 * time.Second,
		},
		{
			name:  "post stop timeout zero is ignored",
			opt:   WithPostStopTimeout,
			field: func(o *observer) *time.Duration { return &o.postStopTimeout },
			dur:   "0s",
			want:  initial,
		},
		{
			name:  "post stop timeout negative is ignored",
			opt:   WithPostStopTimeout,
			field: func(o *observer) *time.Duration { return &o.postStopTimeout },
			dur:   "-1s",
			want:  initial,
		},
		{
			name:  "post stop timeout positive is applied",
			opt:   WithPostStopTimeout,
			field: func(o *observer) *time.Duration { return &o.postStopTimeout },
			dur:   "2s",
			want:  2 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := new(observer)
			*tt.field(o) = initial
			if err := tt.opt(tt.dur)(o); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := *tt.field(o); got != tt.want {
				t.Errorf("duration = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestWithHooks_accumulates pins that repeated WithHooks calls append instead
// of replacing the hooks registered by an earlier call.
func TestWithHooks_accumulates(t *testing.T) {
	t.Parallel()

	h1, h2, h3 := new(recordHook), new(recordHook), new(recordHook)
	tests := []struct {
		name  string
		calls [][]Hook
		want  []Hook
	}{
		{name: "single call", calls: [][]Hook{{h1}}, want: []Hook{h1}},
		{name: "second call keeps the first hooks", calls: [][]Hook{{h1}, {h2, h3}}, want: []Hook{h1, h2, h3}},
		{name: "nil call is a no-op", calls: [][]Hook{{h1}, nil}, want: []Hook{h1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := new(observer)
			for _, hooks := range tt.calls {
				if err := WithHooks(hooks...)(o); err != nil {
					t.Fatalf("WithHooks returned error: %v", err)
				}
			}
			if len(o.hooks) != len(tt.want) {
				t.Fatalf("hooks = %v, want %v", o.hooks, tt.want)
			}
			for i := range tt.want {
				if o.hooks[i] != tt.want[i] {
					t.Errorf("hooks[%d] = %p, want %p", i, o.hooks[i], tt.want[i])
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestWithBackupDuration(t *testing.T) {
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
// 			   got := WithBackupDuration(test.args.dur)
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
// 			   got := WithBackupDuration(test.args.dur)
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
// func TestWithPostStopTimeout(t *testing.T) {
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
// 			   got := WithPostStopTimeout(test.args.dur)
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
// 			   got := WithPostStopTimeout(test.args.dur)
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
// func TestWithWatch(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		enabled bool
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
// 		           enabled:false,
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
// 		           enabled:false,
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
// 			   got := WithWatch(test.args.enabled)
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
// 			   got := WithWatch(test.args.enabled)
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
// func TestWithTicker(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		enabled bool
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
// 		           enabled:false,
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
// 		           enabled:false,
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
// 			   got := WithTicker(test.args.enabled)
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
// 			   got := WithTicker(test.args.enabled)
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
// func TestWithDir(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		dir string
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
// 		           dir:"",
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
// 		           dir:"",
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
// 			   got := WithDir(test.args.dir)
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
// 			   got := WithDir(test.args.dir)
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
// func TestWithBlobStorage(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		storage storage.Storage
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
// 		           storage:nil,
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
// 		           storage:nil,
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
// 			   got := WithBlobStorage(test.args.storage)
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
// 			   got := WithBlobStorage(test.args.storage)
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
// func TestWithHooks(t *testing.T) {
// 	// Change interface type to the type of object you are testing
// 	type T = any
// 	type args struct {
// 		hooks []Hook
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
// 		           hooks:nil,
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
// 		           hooks:nil,
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
// 			   got := WithHooks(test.args.hooks...)
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
// 			   got := WithHooks(test.args.hooks...)
// 			   obj := new(T)
// 			   got(obj)
// 			   if err := checkFunc(test.want, obj); err != nil {
// 			       tt.Errorf("error = %v", err)
// 			   }
// 			*/
// 		})
// 	}
// }
