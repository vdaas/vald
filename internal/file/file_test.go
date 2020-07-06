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

// Package file provides file I/O functionality
package file

import (
	"os"
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
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *os.File) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}

	defaultCheckFunc := func(w want, got *os.File) error {
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
				perm: os.ModeDir,
			},
			checkFunc: func(_ want, got *os.File) error {
				file, err := os.OpenFile("test/data", os.O_CREATE, os.ModeDir)
				if err != nil {
					return err
				}
				return defaultCheckFunc(want{
					want: file,
				}, got)
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
			},
		},

		{
			name: "returns nil when path is empty",
			args: args{
				flg:  os.O_CREATE,
				perm: os.ModeDir,
			},
			want: want{
				want: nil,
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

			got := Open(test.args.path, test.args.flg, test.args.perm)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
