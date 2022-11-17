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
package compress

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMockCompressor_CompressVector(t *testing.T) {
	type args struct {
		vector []float32
	}
	type fields struct {
		CompressVectorFunc   func(vector []float32) (bytes []byte, err error)
		DecompressVectorFunc func(bytes []byte) (vector []float32, err error)
		ReaderFunc           func(src io.ReadCloser) (io.ReadCloser, error)
		WriterFunc           func(dst io.WriteCloser) (io.WriteCloser, error)
	}
	type want struct {
		wantBytes []byte
		err       error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotBytes []byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotBytes, w.wantBytes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotBytes, w.wantBytes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vector: nil,
		       },
		       fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
		           vector: nil,
		           },
		           fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
			m := &MockCompressor{
				CompressVectorFunc:   test.fields.CompressVectorFunc,
				DecompressVectorFunc: test.fields.DecompressVectorFunc,
				ReaderFunc:           test.fields.ReaderFunc,
				WriterFunc:           test.fields.WriterFunc,
			}

			gotBytes, err := m.CompressVector(test.args.vector)
			if err := checkFunc(test.want, gotBytes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockCompressor_DecompressVector(t *testing.T) {
	type args struct {
		bytes []byte
	}
	type fields struct {
		CompressVectorFunc   func(vector []float32) (bytes []byte, err error)
		DecompressVectorFunc func(bytes []byte) (vector []float32, err error)
		ReaderFunc           func(src io.ReadCloser) (io.ReadCloser, error)
		WriterFunc           func(dst io.WriteCloser) (io.WriteCloser, error)
	}
	type want struct {
		wantVector []float32
		err        error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotVector []float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVector, w.wantVector) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVector, w.wantVector)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           bytes: nil,
		       },
		       fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
		           bytes: nil,
		           },
		           fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
			m := &MockCompressor{
				CompressVectorFunc:   test.fields.CompressVectorFunc,
				DecompressVectorFunc: test.fields.DecompressVectorFunc,
				ReaderFunc:           test.fields.ReaderFunc,
				WriterFunc:           test.fields.WriterFunc,
			}

			gotVector, err := m.DecompressVector(test.args.bytes)
			if err := checkFunc(test.want, gotVector, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockCompressor_Reader(t *testing.T) {
	type args struct {
		src io.ReadCloser
	}
	type fields struct {
		CompressVectorFunc   func(vector []float32) (bytes []byte, err error)
		DecompressVectorFunc func(bytes []byte) (vector []float32, err error)
		ReaderFunc           func(src io.ReadCloser) (io.ReadCloser, error)
		WriterFunc           func(dst io.WriteCloser) (io.WriteCloser, error)
	}
	type want struct {
		want io.ReadCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.ReadCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.ReadCloser, err error) error {
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
		           src: nil,
		       },
		       fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
		           src: nil,
		           },
		           fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
			m := &MockCompressor{
				CompressVectorFunc:   test.fields.CompressVectorFunc,
				DecompressVectorFunc: test.fields.DecompressVectorFunc,
				ReaderFunc:           test.fields.ReaderFunc,
				WriterFunc:           test.fields.WriterFunc,
			}

			got, err := m.Reader(test.args.src)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockCompressor_Writer(t *testing.T) {
	type args struct {
		dst io.WriteCloser
	}
	type fields struct {
		CompressVectorFunc   func(vector []float32) (bytes []byte, err error)
		DecompressVectorFunc func(bytes []byte) (vector []float32, err error)
		ReaderFunc           func(src io.ReadCloser) (io.ReadCloser, error)
		WriterFunc           func(dst io.WriteCloser) (io.WriteCloser, error)
	}
	type want struct {
		want io.WriteCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.WriteCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.WriteCloser, err error) error {
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
		           dst: nil,
		       },
		       fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
		           dst: nil,
		           },
		           fields: fields {
		           CompressVectorFunc: nil,
		           DecompressVectorFunc: nil,
		           ReaderFunc: nil,
		           WriterFunc: nil,
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
			m := &MockCompressor{
				CompressVectorFunc:   test.fields.CompressVectorFunc,
				DecompressVectorFunc: test.fields.DecompressVectorFunc,
				ReaderFunc:           test.fields.ReaderFunc,
				WriterFunc:           test.fields.WriterFunc,
			}

			got, err := m.Writer(test.args.dst)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockReadCloser_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		ReadFunc  func(p []byte) (n int, err error)
		CloseFunc func() error
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
		           ReadFunc: nil,
		           CloseFunc: nil,
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
			m := &MockReadCloser{
				ReadFunc:  test.fields.ReadFunc,
				CloseFunc: test.fields.CloseFunc,
			}

			got, err := m.Read(test.args.p)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockReadCloser_Close(t *testing.T) {
	type fields struct {
		ReadFunc  func(p []byte) (n int, err error)
		CloseFunc func() error
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			m := &MockReadCloser{
				ReadFunc:  test.fields.ReadFunc,
				CloseFunc: test.fields.CloseFunc,
			}

			err := m.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockWriteCloser_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		WriteFunc func(p []byte) (n int, err error)
		CloseFunc func() error
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
		           WriteFunc: nil,
		           CloseFunc: nil,
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
			m := &MockWriteCloser{
				WriteFunc: test.fields.WriteFunc,
				CloseFunc: test.fields.CloseFunc,
			}

			got, err := m.Write(test.args.p)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestMockWriteCloser_Close(t *testing.T) {
	type fields struct {
		WriteFunc func(p []byte) (n int, err error)
		CloseFunc func() error
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			m := &MockWriteCloser{
				WriteFunc: test.fields.WriteFunc,
				CloseFunc: test.fields.CloseFunc,
			}

			err := m.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
