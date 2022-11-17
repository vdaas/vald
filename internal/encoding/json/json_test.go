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
package json

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestEncode(t *testing.T) {
	type args struct {
		w    io.Writer
		data interface{}
	}

	type test struct {
		name      string
		args      args
		checkFunc func(err error) error
	}

	tests := []test{
		func() test {
			buf := new(bytes.Buffer)
			data := map[string]string{
				"name": "vald",
			}

			return test{
				name: "returns nil",
				args: args{
					w:    buf,
					data: data,
				},
				checkFunc: func(err error) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := buf.String(), "{\"name\":\"vald\"}\n"; got != want {
						return errors.Errorf("output data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		{
			name: "returns error when type is invalid",
			args: args{
				w:    new(bytes.Buffer),
				data: make(chan struct{}),
			},
			checkFunc: func(err error) error {
				if err == nil {
					return errors.New("err is nil")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Encode(tt.args.w, tt.args.data)
			if err := tt.checkFunc(err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		r    io.Reader
		data map[string]string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(err error, data map[string]string) error
	}

	tests := []test{
		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)

			return test{
				name: "returns nil",
				args: args{
					r:    buf,
					data: make(map[string]string, 1),
				},
				checkFunc: func(err error, data map[string]string) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return errors.Errorf("read data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`1`)

			wantData := make(map[string]string, 1)

			return test{
				name: "returns error when type is invalid",
				args: args{
					r:    buf,
					data: wantData,
				},
				checkFunc: func(err error, data map[string]string) error {
					if err == nil {
						return errors.New("err is nil")
					}

					if !reflect.DeepEqual(data, wantData) {
						return errors.Errorf("data not equals. want: %v, got: %v", wantData, data)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Decode(tt.args.r, &tt.args.data)
			if err := tt.checkFunc(err, tt.args.data); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	type args struct {
		data interface{}
		pref string
		ind  string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(data []byte, err error) error
	}

	tests := []test{
		{
			name: "returns data and nil",
			args: args{
				data: map[string]string{
					"name": "vald",
				},
				pref: "",
				ind:  "",
			},
			checkFunc: func(data []byte, err error) error {
				if err != nil {
					return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
				}

				if got, want := data, []byte("{\n\"name\": \"vald\"\n}"); !reflect.DeepEqual(got, want) {
					return errors.Errorf("data not equals. want: %v, got: %v", string(want), string(got))
				}

				return nil
			},
		},

		{
			name: "returns error when type is invalid",
			args: args{
				data: make(chan struct{}),
				pref: "",
				ind:  "",
			},
			checkFunc: func(data []byte, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}

				if len(data) != 0 {
					return errors.New("data is not empty")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalIndent(tt.args.data, tt.args.pref, tt.args.ind)
			if err := tt.checkFunc(data, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
		i    interface{}
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           data: nil,
		           i: nil,
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
		           data: nil,
		           i: nil,
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

			err := Unmarshal(test.args.data, test.args.i)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	type args struct {
		data interface{}
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           data: nil,
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
		           data: nil,
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

			got, err := Marshal(test.args.data)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
