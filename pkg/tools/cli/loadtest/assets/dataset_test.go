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
package assets

import (
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func Test_identity(t *testing.T) {
	type args struct {
		dim int
	}
	type want struct {
		want func(tb testing.TB) Dataset
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func(tb testing.TB) Dataset) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func(tb testing.TB) Dataset) error {
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
		           dim: 0,
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
		           dim: 0,
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

			got := identity(test.args.dim)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_datasetDir(t *testing.T) {
	type args struct {
		tb testing.TB
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           tb: nil,
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
		           tb: nil,
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

			got := datasetDir(test.args.tb)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestData(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		want func(testing.TB) Dataset
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func(testing.TB) Dataset) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func(testing.TB) Dataset) error {
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
		           name: "",
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
		           name: "",
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

			got := Data(test.args.name)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_Train(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]float32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]float32) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.Train()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_TrainAsFloat64(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]float64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]float64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]float64) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.TrainAsFloat64()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_Query(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]float32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]float32) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.Query()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_QueryAsFloat64(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]float64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]float64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]float64) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.QueryAsFloat64()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_Distances(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]float32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]float32) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.Distances()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_DistancesAsFloat64(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]float64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]float64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]float64) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.DistancesAsFloat64()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_Neighbors(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want [][]int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, [][]int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got [][]int) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.Neighbors()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_IDs(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want []string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []string) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.IDs()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_Name(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.Name()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_Dimension(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.Dimension()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_DistanceType(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.DistanceType()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dataset_ObjectType(t *testing.T) {
	type fields struct {
		train              [][]float32
		trainAsFloat64     [][]float64
		trainOnce          sync.Once
		query              [][]float32
		queryAsFloat64     [][]float64
		queryOnce          sync.Once
		distances          [][]float32
		distancesAsFloat64 [][]float64
		distancesOnce      sync.Once
		neighbors          [][]int
		ids                []string
		name               string
		dimension          int
		distanceType       string
		objectType         string
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
		           train: nil,
		           trainAsFloat64: nil,
		           trainOnce: sync.Once{},
		           query: nil,
		           queryAsFloat64: nil,
		           queryOnce: sync.Once{},
		           distances: nil,
		           distancesAsFloat64: nil,
		           distancesOnce: sync.Once{},
		           neighbors: nil,
		           ids: nil,
		           name: "",
		           dimension: 0,
		           distanceType: "",
		           objectType: "",
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
			d := &dataset{
				train:              test.fields.train,
				trainAsFloat64:     test.fields.trainAsFloat64,
				trainOnce:          test.fields.trainOnce,
				query:              test.fields.query,
				queryAsFloat64:     test.fields.queryAsFloat64,
				queryOnce:          test.fields.queryOnce,
				distances:          test.fields.distances,
				distancesAsFloat64: test.fields.distancesAsFloat64,
				distancesOnce:      test.fields.distancesOnce,
				neighbors:          test.fields.neighbors,
				ids:                test.fields.ids,
				name:               test.fields.name,
				dimension:          test.fields.dimension,
				distanceType:       test.fields.distanceType,
				objectType:         test.fields.objectType,
			}

			got := d.ObjectType()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_float32To64(t *testing.T) {
	type args struct {
		x [][]float32
	}
	type want struct {
		wantY [][]float64
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotY [][]float64) error {
		if !reflect.DeepEqual(gotY, w.wantY) {
			return errors.Errorf("got = %v, want %v", gotY, w.wantY)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           x: nil,
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
		           x: nil,
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

			gotY := float32To64(test.args.x)
			if err := test.checkFunc(test.want, gotY); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
