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

// Package ioutil provides utility function for I/O
package ioutil

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/compress/gob"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.uber.org/goleak"
)

func genStr() []string {
	str := make([]string, 0)
	for i := 0; i < 100; i++ {
		str = append(str, "vdaas.vald.org\n")
	}
	return str
}

func genNonPermittedFile(path string) {
	fp, err := os.Create(path)
	if err != nil {
		log.Error(err)
	}
	err = fp.Chmod(0333)
	if err != nil {
		log.Error(err)
	}
	defer fp.Close()
}

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
			fName := "empty_ioutil_test.txt"
			return test{
				name: "return ([]byte, nil error) when path is exist file name and file is empty",
				args: args{
					path: fName,
				},
				want: want{
					want: []byte{},
				},
				beforeFunc: func(args) {
					fp, err := os.Create(fName)
					if err != nil {
						log.Error(err)
					}
					defer fp.Close()
				},
				afterFunc: func(args) {
					if err := os.Remove(fName); err != nil {
						log.Error(err)
					}
				},
			}
		}(),
		func() test {
			fName := "ioutil_test.txt"
			return test{
				name: "return ([]byte, nil error) when path is exist file name and file is not empty",
				args: args{
					path: fName,
				},
				want: want{
					want: func() []byte {
						strs := genStr()
						buf := &bytes.Buffer{}
						if err := gob.New().NewEncoder(buf).Encode(strs); err != nil {
							log.Error(err)
							return []byte{}
						}
						return buf.Bytes()
					}(),
				},
				beforeFunc: func(args) {
					fp, err := os.Create(fName)
					if err != nil {
						log.Error(err)
					}
					strs := genStr()
					buf := &bytes.Buffer{}
					if err = gob.New().NewEncoder(buf).Encode(strs); err != nil {
						log.Error(err)
					}
					if _, err = fp.Write(buf.Bytes()); err != nil {
						log.Error(err)
					}
					defer fp.Close()
				},
				afterFunc: func(args) {
					if err := os.Remove(fName); err != nil {
						log.Error(err)
					}
				},
			}
		}(),
		func() test {
			fName := "cannot_read_ioutil_test.txt"
			genNonPermittedFile(fName)
			return test{
				name: "return (nil, error) when path is exist file cannot be opend due to permission",
				args: args{
					path: fName,
				},
				want: want{
					err: func() error {
						_, err := os.OpenFile(fName, os.O_RDONLY, os.ModePerm)
						return err
					}(),
				},
				afterFunc: func(args) {
					if err := os.Remove(fName); err != nil {
						log.Error(err)
					}
				},
			}
		}(),
		func() test {
			return test{
				name: "return (nil, error) when path is empty",
				args: args{},
				want: want{
					err: func() error {
						_, err := os.OpenFile("", os.O_RDONLY, os.ModePerm)
						return err
					}(),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when path is not exist file name",
				args: args{
					path: "notexist.txt",
				},
				want: want{
					err: func() error {
						_, err := os.OpenFile("notexist.txt", os.O_RDONLY, os.ModePerm)
						return err
					}(),
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := ReadFile(test.args.path)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
