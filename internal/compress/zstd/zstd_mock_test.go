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
package zstd

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestMockEncoder_Write(t *testing.T) {
	t.Parallel()
	type args struct {
		p []byte
	}
	type fields struct {
		WriteFunc    func(p []byte) (n int, err error)
		CloseFunc    func() error
		ReadFromFunc func(r io.Reader) (n int64, err error)
	}
	type want struct {
		wantN int
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           p: nil,
		       },
		       fields: fields {
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ReadFromFunc: nil,
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
		           p: nil,
		           },
		           fields: fields {
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ReadFromFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockEncoder{
				WriteFunc:    test.fields.WriteFunc,
				CloseFunc:    test.fields.CloseFunc,
				ReadFromFunc: test.fields.ReadFromFunc,
			}

			gotN, err := m.Write(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockEncoder_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		WriteFunc    func(p []byte) (n int, err error)
		CloseFunc    func() error
		ReadFromFunc func(r io.Reader) (n int64, err error)
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
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
		       fields: fields {
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ReadFromFunc: nil,
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
		           fields: fields {
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ReadFromFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockEncoder{
				WriteFunc:    test.fields.WriteFunc,
				CloseFunc:    test.fields.CloseFunc,
				ReadFromFunc: test.fields.ReadFromFunc,
			}

			err := m.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockEncoder_ReadFrom(t *testing.T) {
	t.Parallel()
	type args struct {
		r io.Reader
	}
	type fields struct {
		WriteFunc    func(p []byte) (n int, err error)
		CloseFunc    func() error
		ReadFromFunc func(r io.Reader) (n int64, err error)
	}
	type want struct {
		wantN int64
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int64, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           r: nil,
		       },
		       fields: fields {
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ReadFromFunc: nil,
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
		           r: nil,
		           },
		           fields: fields {
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ReadFromFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockEncoder{
				WriteFunc:    test.fields.WriteFunc,
				CloseFunc:    test.fields.CloseFunc,
				ReadFromFunc: test.fields.ReadFromFunc,
			}

			gotN, err := m.ReadFrom(test.args.r)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockDecoder_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		CloseFunc   func()
		ReadFunc    func(p []byte) (int, error)
		WriteToFunc func(w io.Writer) (int64, error)
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           CloseFunc: nil,
		           ReadFunc: nil,
		           WriteToFunc: nil,
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
		           fields: fields {
		           CloseFunc: nil,
		           ReadFunc: nil,
		           WriteToFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockDecoder{
				CloseFunc:   test.fields.CloseFunc,
				ReadFunc:    test.fields.ReadFunc,
				WriteToFunc: test.fields.WriteToFunc,
			}

			m.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockDecoder_Read(t *testing.T) {
	t.Parallel()
	type args struct {
		p []byte
	}
	type fields struct {
		CloseFunc   func()
		ReadFunc    func(p []byte) (int, error)
		WriteToFunc func(w io.Writer) (int64, error)
	}
	type want struct {
		want int
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got int, err error) error {
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
		           p: nil,
		       },
		       fields: fields {
		           CloseFunc: nil,
		           ReadFunc: nil,
		           WriteToFunc: nil,
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
		           p: nil,
		           },
		           fields: fields {
		           CloseFunc: nil,
		           ReadFunc: nil,
		           WriteToFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockDecoder{
				CloseFunc:   test.fields.CloseFunc,
				ReadFunc:    test.fields.ReadFunc,
				WriteToFunc: test.fields.WriteToFunc,
			}

			got, err := m.Read(test.args.p)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockDecoder_WriteTo(t *testing.T) {
	t.Parallel()
	type fields struct {
		CloseFunc   func()
		ReadFunc    func(p []byte) (int, error)
		WriteToFunc func(w io.Writer) (int64, error)
	}
	type want struct {
		want  int64
		wantW string
		err   error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int64, string, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int64, gotW string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(gotW, w.wantW) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotW, w.wantW)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           CloseFunc: nil,
		           ReadFunc: nil,
		           WriteToFunc: nil,
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
		           fields: fields {
		           CloseFunc: nil,
		           ReadFunc: nil,
		           WriteToFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockDecoder{
				CloseFunc:   test.fields.CloseFunc,
				ReadFunc:    test.fields.ReadFunc,
				WriteToFunc: test.fields.WriteToFunc,
			}
			w := &bytes.Buffer{}

			got, err := m.WriteTo(w)
			if err := test.checkFunc(test.want, got, w.String(), err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockZstd_NewWriter(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []zstd.EOption
	}
	type fields struct {
		NewWriterFunc func(w io.Writer, opts ...zstd.EOption) (Encoder, error)
		NewReaderFunc func(r io.Reader, opts ...zstd.DOption) (Decoder, error)
	}
	type want struct {
		want  Encoder
		wantW string
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Encoder, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Encoder, gotW string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(gotW, w.wantW) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotW, w.wantW)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
		       },
		       fields: fields {
		           NewWriterFunc: nil,
		           NewReaderFunc: nil,
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
		           opts: nil,
		           },
		           fields: fields {
		           NewWriterFunc: nil,
		           NewReaderFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockZstd{
				NewWriterFunc: test.fields.NewWriterFunc,
				NewReaderFunc: test.fields.NewReaderFunc,
			}
			w := &bytes.Buffer{}

			got, err := m.NewWriter(w, test.args.opts...)
			if err := test.checkFunc(test.want, got, w.String(), err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockZstd_NewReader(t *testing.T) {
	t.Parallel()
	type args struct {
		r    io.Reader
		opts []zstd.DOption
	}
	type fields struct {
		NewWriterFunc func(w io.Writer, opts ...zstd.EOption) (Encoder, error)
		NewReaderFunc func(r io.Reader, opts ...zstd.DOption) (Decoder, error)
	}
	type want struct {
		want Decoder
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Decoder, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Decoder, err error) error {
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
		           r: nil,
		           opts: nil,
		       },
		       fields: fields {
		           NewWriterFunc: nil,
		           NewReaderFunc: nil,
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
		           r: nil,
		           opts: nil,
		           },
		           fields: fields {
		           NewWriterFunc: nil,
		           NewReaderFunc: nil,
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
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &MockZstd{
				NewWriterFunc: test.fields.NewWriterFunc,
				NewReaderFunc: test.fields.NewReaderFunc,
			}

			got, err := m.NewReader(test.args.r, test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
