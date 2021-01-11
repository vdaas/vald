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

// Package gongt provides gongt
package gongt

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/yahoojapan/gongt"

	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want algorithm.Bit64
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, algorithm.Bit64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got algorithm.Bit64, err error) error {
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
		           opts: nil,
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
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_Search(t *testing.T) {
	type args struct {
		vec     []float64
		size    int
		epsilon float32
		radius  float32
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
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
		           vec: nil,
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           vec: nil,
		           size: 0,
		           epsilon: 0,
		           radius: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			got, err := c.Search(test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_Insert(t *testing.T) {
	type args struct {
		vec []float64
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
	}
	type want struct {
		want uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got uint, err error) error {
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
		           vec: nil,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           vec: nil,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			got, err := c.Insert(test.args.vec)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_InsertCommit(t *testing.T) {
	type args struct {
		vec      []float64
		poolSize uint32
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
	}
	type want struct {
		want uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got uint, err error) error {
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
		           vec: nil,
		           poolSize: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           vec: nil,
		           poolSize: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			got, err := c.InsertCommit(test.args.vec, test.args.poolSize)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_BulkInsert(t *testing.T) {
	type args struct {
		vecs [][]float64
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
	}
	type want struct {
		want  []uint
		want1 []error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []uint, []error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []uint, got1 []error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vecs: nil,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           vecs: nil,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			got, got1 := c.BulkInsert(test.args.vecs)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_BulkInsertCommit(t *testing.T) {
	type args struct {
		vecs     [][]float64
		poolSize uint32
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
	}
	type want struct {
		want  []uint
		want1 []error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []uint, []error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []uint, got1 []error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vecs: nil,
		           poolSize: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           vecs: nil,
		           poolSize: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			got, got1 := c.BulkInsertCommit(test.args.vecs, test.args.poolSize)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_CreateAndSaveIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
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
		           poolSize: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           poolSize: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			err := c.CreateAndSaveIndex(test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_CreateIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
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
		           poolSize: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           poolSize: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			err := c.CreateIndex(test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_Remove(t *testing.T) {
	type args struct {
		id uint
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
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
		           id: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           id: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			err := c.Remove(test.args.id)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_BulkRemove(t *testing.T) {
	type args struct {
		ids []uint
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
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
		           ids: nil,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           ids: nil,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			err := c.BulkRemove(test.args.ids...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_GetVector(t *testing.T) {
	type args struct {
		id uint
	}
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
	}
	type want struct {
		want []float64
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float64, err error) error {
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
		           id: 0,
		       },
		       fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           id: 0,
		           },
		           fields: fields {
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			got, err := c.GetVector(test.args.id)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_core_Close(t *testing.T) {
	type fields struct {
		indexPath  string
		tmpdir     string
		objectType ObjectType
		dimension  int
		NGT        *gongt.NGT
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
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
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
		           indexPath: "",
		           tmpdir: "",
		           objectType: nil,
		           dimension: 0,
		           NGT: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &core{
				indexPath:  test.fields.indexPath,
				tmpdir:     test.fields.tmpdir,
				objectType: test.fields.objectType,
				dimension:  test.fields.dimension,
				NGT:        test.fields.NGT,
			}

			c.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_toUint(t *testing.T) {
	type args struct {
		in []int
	}
	type want struct {
		wantOut []uint
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []uint) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotOut []uint) error {
		if !reflect.DeepEqual(gotOut, w.wantOut) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOut, w.wantOut)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           in: nil,
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
		           in: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotOut := toUint(test.args.in)
			if err := test.checkFunc(test.want, gotOut); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
