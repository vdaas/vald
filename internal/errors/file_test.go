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
package errors

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestErrWatchDirNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
		name       string
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
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
		name       string
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
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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

type fitosFakeFileInfo struct {
	name string
	size int64
	mode fs.FileMode
}

func (f fitosFakeFileInfo) Name() string       { return f.name }
func (f fitosFakeFileInfo) Size() int64        { return f.size }
func (f fitosFakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f fitosFakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fitosFakeFileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f fitosFakeFileInfo) Sys() any           { return nil }

func Test_fitos(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "missing")
	const format = "{name: %s, size: %d, mode: %s, mode_int: %d, is_dir: %v}"

	tests := []struct {
		fi   os.FileInfo
		name string
		path string
		want string
	}{
		{
			// a given fi must be used as-is: the (nonexistent) path is never stat'ed,
			// otherwise the stat failure would produce the "unknown file info" string.
			name: "given fi is used without stat",
			path: missing,
			fi:   fitosFakeFileInfo{name: "fake", size: 7, mode: 0o640},
			want: fmt.Sprintf(format, "fake", 7, fs.FileMode(0o640).String(), fs.FileMode(0o640), false),
		},
		{
			name: "given dir fi",
			path: missing,
			fi:   fitosFakeFileInfo{name: "d", size: 0, mode: fs.ModeDir | 0o755},
			want: fmt.Sprintf(format, "d", 0, (fs.ModeDir | 0o755).String(), fs.ModeDir|0o755, true),
		},
		{
			name: "nil fi and existing path is stat'ed",
			path: file,
			want: fmt.Sprintf(format, "f.txt", 3, fs.FileMode(0o600).String(), fs.FileMode(0o600), false),
		},
		{
			name: "nil fi and missing path",
			path: missing,
			want: fmt.Sprintf("unknown file info: for %s\t%v", missing, os.FileInfo(nil)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := fitos(tc.path, tc.fi); got != tc.want {
				t.Errorf("fitos() = %q, want %q", got, tc.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_fitos(t *testing.T) {
// 	type args struct {
// 		path string
// 		fi   os.FileInfo
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		           fi:nil,
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
// 		           path:"",
// 		           fi:nil,
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
//
// 			got := fitos(test.args.path, test.args.fi)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
