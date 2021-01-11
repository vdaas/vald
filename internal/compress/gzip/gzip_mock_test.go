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
package gzip

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestMockReader_Read(t *testing.T) {
	t.Parallel()
	type args struct {
		p []byte
	}
	type fields struct {
		ReadFunc        func(p []byte) (n int, err error)
		CloseFunc       func() error
		ResetFunc       func(r io.Reader) error
		MultistreamFunc func(ok bool)
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
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
			m := &MockReader{
				ReadFunc:        test.fields.ReadFunc,
				CloseFunc:       test.fields.CloseFunc,
				ResetFunc:       test.fields.ResetFunc,
				MultistreamFunc: test.fields.MultistreamFunc,
			}

			gotN, err := m.Read(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockReader_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		ReadFunc        func(p []byte) (n int, err error)
		CloseFunc       func() error
		ResetFunc       func(r io.Reader) error
		MultistreamFunc func(ok bool)
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
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
			m := &MockReader{
				ReadFunc:        test.fields.ReadFunc,
				CloseFunc:       test.fields.CloseFunc,
				ResetFunc:       test.fields.ResetFunc,
				MultistreamFunc: test.fields.MultistreamFunc,
			}

			err := m.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockReader_Reset(t *testing.T) {
	t.Parallel()
	type args struct {
		r io.Reader
	}
	type fields struct {
		ReadFunc        func(p []byte) (n int, err error)
		CloseFunc       func() error
		ResetFunc       func(r io.Reader) error
		MultistreamFunc func(ok bool)
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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
		           r: nil,
		       },
		       fields: fields {
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
			m := &MockReader{
				ReadFunc:        test.fields.ReadFunc,
				CloseFunc:       test.fields.CloseFunc,
				ResetFunc:       test.fields.ResetFunc,
				MultistreamFunc: test.fields.MultistreamFunc,
			}

			err := m.Reset(test.args.r)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockReader_Multistream(t *testing.T) {
	t.Parallel()
	type args struct {
		ok bool
	}
	type fields struct {
		ReadFunc        func(p []byte) (n int, err error)
		CloseFunc       func() error
		ResetFunc       func(r io.Reader) error
		MultistreamFunc func(ok bool)
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ok: false,
		       },
		       fields: fields {
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
		           ok: false,
		           },
		           fields: fields {
		           ReadFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           MultistreamFunc: nil,
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
			m := &MockReader{
				ReadFunc:        test.fields.ReadFunc,
				CloseFunc:       test.fields.CloseFunc,
				ResetFunc:       test.fields.ResetFunc,
				MultistreamFunc: test.fields.MultistreamFunc,
			}

			m.Multistream(test.args.ok)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockWriter_Write(t *testing.T) {
	t.Parallel()
	type args struct {
		p []byte
	}
	type fields struct {
		WriteFunc func(p []byte) (n int, err error)
		CloseFunc func() error
		ResetFunc func(w io.Writer)
		FlushFunc func() error
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
			m := &MockWriter{
				WriteFunc: test.fields.WriteFunc,
				CloseFunc: test.fields.CloseFunc,
				ResetFunc: test.fields.ResetFunc,
				FlushFunc: test.fields.FlushFunc,
			}

			gotN, err := m.Write(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockWriter_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		WriteFunc func(p []byte) (n int, err error)
		CloseFunc func() error
		ResetFunc func(w io.Writer)
		FlushFunc func() error
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
			m := &MockWriter{
				WriteFunc: test.fields.WriteFunc,
				CloseFunc: test.fields.CloseFunc,
				ResetFunc: test.fields.ResetFunc,
				FlushFunc: test.fields.FlushFunc,
			}

			err := m.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockWriter_Reset(t *testing.T) {
	t.Parallel()
	type fields struct {
		WriteFunc func(p []byte) (n int, err error)
		CloseFunc func() error
		ResetFunc func(w io.Writer)
		FlushFunc func() error
	}
	type want struct {
		wantW string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotW string) error {
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
		           WriteFunc: nil,
		           CloseFunc: nil,
		           ResetFunc: nil,
		           FlushFunc: nil,
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
			m := &MockWriter{
				WriteFunc: test.fields.WriteFunc,
				CloseFunc: test.fields.CloseFunc,
				ResetFunc: test.fields.ResetFunc,
				FlushFunc: test.fields.FlushFunc,
			}
			w := &bytes.Buffer{}

			m.Reset(w)
			if err := test.checkFunc(test.want, w.String()); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockWriter_Flush(t *testing.T) {
	t.Parallel()
	type fields struct {
		WriteFunc func(p []byte) (n int, err error)
		CloseFunc func() error
		ResetFunc func(w io.Writer)
		FlushFunc func() error
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
		           ResetFunc: nil,
		           FlushFunc: nil,
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
			m := &MockWriter{
				WriteFunc: test.fields.WriteFunc,
				CloseFunc: test.fields.CloseFunc,
				ResetFunc: test.fields.ResetFunc,
				FlushFunc: test.fields.FlushFunc,
			}

			err := m.Flush()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockGzip_NewWriterLevel(t *testing.T) {
	t.Parallel()
	type args struct {
		level int
	}
	type fields struct {
		NewWriterLevelFunc func(w io.Writer, level int) (Writer, error)
		NewReaderFunc      func(r io.Reader) (Reader, error)
	}
	type want struct {
		want  Writer
		wantW string
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Writer, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Writer, gotW string, err error) error {
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
		           level: 0,
		       },
		       fields: fields {
		           NewWriterLevelFunc: nil,
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
		           level: 0,
		           },
		           fields: fields {
		           NewWriterLevelFunc: nil,
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
			m := &MockGzip{
				NewWriterLevelFunc: test.fields.NewWriterLevelFunc,
				NewReaderFunc:      test.fields.NewReaderFunc,
			}
			w := &bytes.Buffer{}

			got, err := m.NewWriterLevel(w, test.args.level)
			if err := test.checkFunc(test.want, got, w.String(), err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockGzip_NewReader(t *testing.T) {
	t.Parallel()
	type args struct {
		r io.Reader
	}
	type fields struct {
		NewWriterLevelFunc func(w io.Writer, level int) (Writer, error)
		NewReaderFunc      func(r io.Reader) (Reader, error)
	}
	type want struct {
		want Reader
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Reader, err error) error {
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
		       },
		       fields: fields {
		           NewWriterLevelFunc: nil,
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
		           },
		           fields: fields {
		           NewWriterLevelFunc: nil,
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
			m := &MockGzip{
				NewWriterLevelFunc: test.fields.NewWriterLevelFunc,
				NewReaderFunc:      test.fields.NewReaderFunc,
			}

			got, err := m.NewReader(test.args.r)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
