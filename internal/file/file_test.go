//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package file provides file I/O functionality
package file

import (
	"os"
	"syscall"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestOpen(t *testing.T) {
	type args struct {
		path string
		flg  int
		perm os.FileMode
	}
	type want struct {
		want *os.File
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *os.File, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}

	defaultCheckFunc := func(w want, got *os.File, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if w.want == nil {
			if got != nil {
				return errors.New("got is not nil")
			}
		} else {
			if got, want := got.Name(), w.want.Name(); got != want {
				return errors.Errorf("got name = %s, want: %s")
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "returns *os.File when path is `test/data`",
			args: args{
				path: "test/data",
				flg:  os.O_CREATE,
				perm: os.ModePerm,
			},
			checkFunc: func(_ want, got *os.File, gotErr error) error {
				file, err := os.OpenFile("test/data", os.O_CREATE, os.ModePerm)
				if err != nil {
					return err
				}
				return defaultCheckFunc(want{
					want: file,
					err:  nil,
				}, got, gotErr)
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				if err := os.RemoveAll("test"); err != nil {
					t.Fatal(err)
				}
			},
		},

		{
			name: "returns *os.File when path is `test/test/data`",
			args: args{
				path: "test/test/data",
				flg:  os.O_CREATE,
				perm: os.ModePerm,
			},
			checkFunc: func(_ want, got *os.File, gotErr error) error {
				file, err := os.OpenFile("test/test/data", os.O_CREATE, os.ModePerm)
				if err != nil {
					return err
				}
				return defaultCheckFunc(want{
					want: file,
					err:  nil,
				}, got, gotErr)
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				if err := os.RemoveAll("test"); err != nil {
					t.Fatal(err)
				}
			},
		},

		{
			name: "returns *os.File when path is `file.go`",
			args: args{
				path: "file.go",
				flg:  os.O_RDONLY,
				perm: os.ModePerm,
			},
			want: want{
				want: func() *os.File {
					f, _ := os.OpenFile("file.go", os.O_RDONLY, os.ModePerm)
					return f
				}(),
				err: nil,
			},
		},

		{
			name: "returns (nil, error) when path is empty",
			args: args{
				flg:  os.O_CREATE,
				perm: os.ModeDir,
			},
			want: want{
				want: nil,
				err:  errors.ErrPathNotSpecified,
			},
		},

		{
			name: "returns (nil, error) when file does not exists and flag is not CREATE or APPEND",
			args: args{
				path: "dummy",
				flg:  os.O_RDONLY,
				perm: os.ModePerm,
			},
			want: want{
				want: nil,
				err: &os.PathError{
					Op:   "open",
					Path: "dummy",
					Err: func() error {
						_, err := syscall.Open("dummy", syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
						return err
					}(),
				},
			},
		},

		{
			name: "returns (nil, error) when the folder does not exists and flag is not CREATE or APPEND",
			args: args{
				path: "dummy/dummy",
				flg:  os.O_RDONLY,
				perm: os.ModePerm,
			},
			want: want{
				want: nil,
				err: &os.PathError{
					Op:   "open",
					Path: "dummy/dummy",
					Err: func() error {
						_, err := syscall.Open("dummy/", syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
						return err
					}(),
				},
			},
			afterFunc: func(t *testing.T, _ args) {
				t.Helper()
				if err := os.RemoveAll("dummy"); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := Open(test.args.path, test.args.flg, test.args.perm)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
