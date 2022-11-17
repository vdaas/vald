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
package errors

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestErrWatchDirNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrWatchDirNotFound error",
			want: want{
				want: New("fs watcher watch dir not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := ErrWatchDirNotFound
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFileAlreadyExists(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrFileAlreadyExists error with the file path is 'metadata.json'",
			args: args{
				path: "metadata.json",
			},
			want: want{
				want: New("file already exists: metadata.json"),
			},
		},
		{
			name: "return ErrFileAlreadyExists error with the file path is empty",
			args: args{
				path: "",
			},
			want: want{
				want: New("file already exists: "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrFileAlreadyExists(test.args.path)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrPathNotSpecified(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrPathNotSpecified error",
			want: want{
				want: New("the path is not specified"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := ErrPathNotSpecified
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrPathNotAllowed(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ErrPathNotAllowed error with the file path is '../x.json'",
			args: args{
				path: "../x.json",
			},
			want: want{
				want: New("the specified file path is not allowed: ../x.json"),
			},
		},
		{
			name: "return ErrPathNotAllowed error with the file path is empty",
			args: args{
				path: "",
			},
			want: want{
				want: New("the specified file path is not allowed: "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrPathNotAllowed(test.args.path)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_fitos(t *testing.T) {
	type args struct {
		path string
		fi   os.FileInfo
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           path: "",
		           fi: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           path: "",
		           fi: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := fitos(test.args.path, test.args.fi)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
