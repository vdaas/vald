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

// Package ioutil provides utility function for I/O
package ioutil

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestReadFile(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		want []byte
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			const content = "test case 1"
			tempfile, err := ioutil.TempFile("", "")
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if _, err = tempfile.WriteString(content); err != nil {
				t.Errorf("error = %v", err)
			}
			path := tempfile.Name()
			return test{
				name: "success when output string to file",
				args: args{
					path: path,
				},
				want: want{
					want: []byte(content),
					err:  nil,
				},
				afterFunc: func(a args) {
					if err := os.Remove(a.path); err != nil {
						t.Errorf("error = %v", err)
					}
				},
			}
		}(),
		func() test {
			return test{
				name: "fail with empty path",
				args: args{
					path: "",
				},
				want: want{
					want: nil,
					err:  errors.New("the path is not specified"),
				},
			}
		}(),
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

			got, err := ReadFile(test.args.path)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
