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
package x1b

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func Test_open(t *testing.T) {
	type args struct {
		fname       string
		elementSize int
	}
	type want struct {
		wantF *file
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *file, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotF *file, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotF, w.wantF) {
			return errors.Errorf("got = %v, want %v", gotF, w.wantF)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fname: "",
		           elementSize: 0,
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
		           fname: "",
		           elementSize: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			gotF, err := open(test.args.fname, test.args.elementSize)
			if err := test.checkFunc(test.want, gotF, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_file_Close(t *testing.T) {
	type fields struct {
		mem   []byte
		dim   int
		size  int
		block int
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
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
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			f := &file{
				mem:   test.fields.mem,
				dim:   test.fields.dim,
				size:  test.fields.size,
				block: test.fields.block,
			}

			err := f.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_file_load(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		mem   []byte
		dim   int
		size  int
		block int
	}
	type want struct {
		want []byte
		err  error
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
	defaultCheckFunc := func(w want, got []byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
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
		           i: 0,
		           },
		           fields: fields {
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			f := &file{
				mem:   test.fields.mem,
				dim:   test.fields.dim,
				size:  test.fields.size,
				block: test.fields.block,
			}

			got, err := f.load(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_file_Dimension(t *testing.T) {
	type fields struct {
		mem   []byte
		dim   int
		size  int
		block int
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
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
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			f := &file{
				mem:   test.fields.mem,
				dim:   test.fields.dim,
				size:  test.fields.size,
				block: test.fields.block,
			}

			got := f.Dimension()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_file_Size(t *testing.T) {
	type fields struct {
		mem   []byte
		dim   int
		size  int
		block int
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
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
		           mem: nil,
		           dim: 0,
		           size: 0,
		           block: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			f := &file{
				mem:   test.fields.mem,
				dim:   test.fields.dim,
				size:  test.fields.size,
				block: test.fields.block,
			}

			got := f.Size()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bvecs_LoadUint8(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		file *file
	}
	type want struct {
		want []uint8
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []uint8, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []uint8, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           file: file{},
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
		           i: 0,
		           },
		           fields: fields {
		           file: file{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			bv := &bvecs{
				file: test.fields.file,
			}

			got, err := bv.LoadUint8(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_bvecs_Load(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		file *file
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           file: file{},
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
		           i: 0,
		           },
		           fields: fields {
		           file: file{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			bv := &bvecs{
				file: test.fields.file,
			}

			got, err := bv.Load(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_fvecs_LoadFloat32(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		file *file
	}
	type want struct {
		want []float32
		err  error
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
	defaultCheckFunc := func(w want, got []float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           file: file{},
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
		           i: 0,
		           },
		           fields: fields {
		           file: file{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			fv := &fvecs{
				file: test.fields.file,
			}

			got, err := fv.LoadFloat32(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_fvecs_Load(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		file *file
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           file: file{},
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
		           i: 0,
		           },
		           fields: fields {
		           file: file{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			fv := &fvecs{
				file: test.fields.file,
			}

			got, err := fv.Load(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ivecs_LoadInt32(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		file *file
	}
	type want struct {
		want []int32
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []int32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []int32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           file: file{},
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
		           i: 0,
		           },
		           fields: fields {
		           file: file{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			iv := &ivecs{
				file: test.fields.file,
			}

			got, err := iv.LoadInt32(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ivecs_Load(t *testing.T) {
	type args struct {
		i int
	}
	type fields struct {
		file *file
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: 0,
		       },
		       fields: fields {
		           file: file{},
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
		           i: 0,
		           },
		           fields: fields {
		           file: file{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			iv := &ivecs{
				file: test.fields.file,
			}

			got, err := iv.Load(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestNewBVecs(t *testing.T) {
	type args struct {
		fname string
	}
	type want struct {
		want Uint8Vectors
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Uint8Vectors, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Uint8Vectors, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fname: "",
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
		           fname: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			got, err := NewUint8Vectors(test.args.fname)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestNewFVecs(t *testing.T) {
	type args struct {
		fname string
	}
	type want struct {
		want FloatVectors
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, FloatVectors, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got FloatVectors, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fname: "",
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
		           fname: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			got, err := NewFloatVectors(test.args.fname)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestNewIVecs(t *testing.T) {
	type args struct {
		fname string
	}
	type want struct {
		want Int32Vectors
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Int32Vectors, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Int32Vectors, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fname: "",
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
		           fname: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			got, err := NewInt32Vectors(test.args.fname)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestOpen(t *testing.T) {
	type args struct {
		fname string
	}
	type want struct {
		want BillionScaleVectors
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, BillionScaleVectors, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got BillionScaleVectors, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fname: "",
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
		           fname: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			got, err := Open(test.args.fname)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
