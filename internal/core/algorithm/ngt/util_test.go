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

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

import (
	"os"
	"reflect"
	"testing"

	"go.uber.org/goleak"

	"github.com/vdaas/vald/internal/errors"
)

func Test_fileExists(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	const (
		baseDir      = "./utiltest"
		testDirPath  = baseDir + "/index"
		testFilePath = baseDir + "/ngt-meta.kvsdb"
		testSym      = "sym"

		testFailsDirPath  = baseDir + "/fails-index"
		testFailsFilePath = baseDir + "/ngt-meta-fails.kvsdb"
	)

	defaultAfterFunc := func(t *testing.T, args args) {
		t.Helper()
		if err := os.RemoveAll(baseDir); err != nil {
			t.Error(err)
		}
		if err := os.RemoveAll(args.path); err != nil {
			t.Error(err)
		}
	}

	tests := []test{
		{
			name: "return true when the directory exists",
			args: args{
				path: testDirPath,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(args.path, 0o750); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the directory exists and the type is symbolic link",
			args: args{
				path: testSym,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(testDirPath, 0o750); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(testDirPath, testSym); err != nil {
					t.Error(err)
				}
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the file exists",
			args: args{
				path: testFilePath,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(baseDir, 0o750); err != nil {
					t.Fatal(err)
				}

				f, err := os.Create(args.path)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						t.Error(err)
					}
				}()
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the file exists and the type is symbolic link",
			args: args{
				path: testSym,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(baseDir, 0o750); err != nil {
					t.Fatal(err)
				}

				f, err := os.Create(testFilePath)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						t.Error(err)
					}
				}()

				if err := os.Symlink(testFilePath, testSym); err != nil {
					t.Error(err)
				}
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when the directory does not exist",
			args: args{
				path: testFailsDirPath,
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the file exists",
			args: args{
				path: testFailsFilePath,
			},
			want: want{
				want: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			defer test.afterFunc(tt, test.args)

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := fileExists(test.args.path)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
