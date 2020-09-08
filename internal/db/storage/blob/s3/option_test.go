//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package s3

import (
	stderrs "errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}
)

func TestWithErrGroup(t *testing.T) {
	type T = client
	type args struct {
		eg errgroup.Group
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			eg := errgroup.Get()
			return test{
				name: "set success when eg is not nil",
				args: args{
					eg: eg,
				},
				want: want{
					obj: &T{
						eg: eg,
					},
				},
			}
		}(),

		func() test {
			return test{
				name: "set nothing when eg is nil",
				args: args{
					eg: nil,
				},
				want: want{
					obj: &T{
						eg: nil,
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithErrGroup(test.args.eg)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSession(t *testing.T) {
	type T = client
	type args struct {
		sess *session.Session
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			sess := new(session.Session)
			return test{
				name: "set success when sess is not nil",
				args: args{
					sess: sess,
				},
				want: want{
					obj: &T{
						session: sess,
					},
				},
			}
		}(),

		func() test {
			return test{
				name: "set nothing when sess is not nil",
				args: args{
					sess: nil,
				},
				want: want{
					obj: &T{
						session: nil,
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithSession(test.args.sess)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBucket(t *testing.T) {
	type T = client
	type args struct {
		bucket string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when bucket is `bucket`",
			args: args{
				bucket: "bucket",
			},
			want: want{
				obj: &T{
					bucket: "bucket",
				},
			},
		},

		{
			name: "set success when bucket is empty",
			args: args{
				bucket: "",
			},
			want: want{
				obj: new(T),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithBucket(test.args.bucket)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxPartSize(t *testing.T) {
	type T = client
	type args struct {
		size string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set nothing when size is `100G`",
			args: args{
				size: "100G",
			},
			want: want{
				obj: &T{
					maxPartSize: 107374182400,
				},
			},
		},

		{
			name: "set nothing when size is `1M`",
			args: args{
				size: "1M",
			},
			want: want{
				obj: &T{
					maxPartSize: 0,
				},
			},
		},

		{
			name: "set default and returns error when size given with invalid string",
			args: args{
				size: "a",
			},
			want: want{
				obj: &T{
					maxPartSize: 0,
				},
				err: func() (err error) {
					err = stderrs.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")
					err = errors.Wrap(err, errors.ErrParseUnitFailed("a").Error())
					return
				}(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithMaxPartSize(test.args.size)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxChunkSize(t *testing.T) {
	type T = client
	type args struct {
		size string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when size is `100G`",
			args: args{
				size: "100G",
			},
			want: want{
				obj: &T{
					maxChunkSize: 107374182400,
				},
				err: nil,
			},
		},

		{
			name: "set nothing when size is `1M`",
			args: args{
				size: "1M",
			},
			want: want{
				obj: &T{
					maxChunkSize: 0,
				},
				err: nil,
			},
		},

		{
			name: "returns error when size is `a`",
			args: args{
				size: "a",
			},
			want: want{
				obj: &T{
					maxChunkSize: 0,
				},
				err: func() (err error) {
					err = stderrs.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")
					err = errors.Wrap(err, errors.ErrParseUnitFailed("a").Error())
					return
				}(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithMaxChunkSize(test.args.size)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReaderBackoff(t *testing.T) {
	type T = client
	type args struct {
		enabled bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when enabled is `true`",
			args: args{
				enabled: true,
			},
			want: want{
				obj: &T{
					readerBackoffEnabled: true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithReaderBackoff(test.args.enabled)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReaderBackoffOpts(t *testing.T) {
	type T = client
	type args struct {
		opts []backoff.Option
	}
	type fields struct {
		readerBackoffOpts []backoff.Option
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []cmp.Option{
			cmp.AllowUnexported(*obj),
			cmp.AllowUnexported(*w.obj),
			cmp.Comparer(func(want, got []backoff.Option) bool {
				return len(got) == len(want)
			}),
			cmp.Comparer(func(want, got backoff.Option) bool {
				return reflect.ValueOf(got).Pointer() == reflect.ValueOf(want).Pointer()
			}),
		}
		if diff := cmp.Diff(w.obj, obj, opts...); diff != "" {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}

		return nil
	}

	tests := []test{
		{
			name: "set success when opts is not nil and c.readerBackoffOpts is not nil",
			args: args{
				opts: []backoff.Option{
					backoff.WithRetryCount(1),
				},
			},
			fields: fields{
				readerBackoffOpts: []backoff.Option{
					backoff.WithRetryCount(1),
				},
			},
			want: want{
				obj: &T{
					readerBackoffOpts: []backoff.Option{
						backoff.WithRetryCount(1),
						backoff.WithRetryCount(1),
					},
				},
				err: nil,
			},
		},

		{
			name: "set success when opts is not nil and r.readerBackoffOpts is nil",
			args: args{
				opts: []backoff.Option{
					backoff.WithRetryCount(1),
				},
			},
			want: want{
				obj: &T{
					readerBackoffOpts: []backoff.Option{
						backoff.WithRetryCount(1),
					},
				},
				err: nil,
			},
		},

		{
			name: "set nothing when opts is nil",
			args: args{
				opts: nil,
			},
			want: want{
				obj: new(T),
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithReaderBackoffOpts(test.args.opts...)
			obj := &T{
				readerBackoffOpts: test.fields.readerBackoffOpts,
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxChunkSize(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		size string
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
	           return errors.Errorf("got error = %v, want %v", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
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
		           size: "",
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
		           size: "",
		           },
		           want: want {
		               obj: new(T),
		           },
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := WithMaxChunkSize(test.args.size)
			   obj := new(T)
			   if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := WithMaxChunkSize(test.args.size)
			   obj := new(T)
			   got(obj)
			   if err := test.checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithReaderBackoff(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		enabled bool
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
	           return errors.Errorf("got error = %v, want %v", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
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
		           enabled: false,
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
		           enabled: false,
		           },
		           want: want {
		               obj: new(T),
		           },
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := WithReaderBackoff(test.args.enabled)
			   obj := new(T)
			   if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := WithReaderBackoff(test.args.enabled)
			   obj := new(T)
			   got(obj)
			   if err := test.checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithReaderBackoffOpts(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		opts []backoff.Option
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
	           return errors.Errorf("got error = %v, want %v", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := WithReaderBackoffOpts(test.args.opts...)
			   obj := new(T)
			   if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := WithReaderBackoffOpts(test.args.opts...)
			   obj := new(T)
			   got(obj)
			   if err := test.checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}
