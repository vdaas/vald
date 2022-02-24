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
package assets

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/assets/x1b"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_loadLargeData(t *testing.T) {
	t.Parallel()
	type args struct {
		trainFileName       string
		queryFileName       string
		groundTruthFileName string
		distanceFileName    string
		name                string
		distanceType        string
		objectType          string
	}
	type want struct {
		want func() (Dataset, error)
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() (Dataset, error)) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() (Dataset, error)) error {
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
		           trainFileName: "",
		           queryFileName: "",
		           groundTruthFileName: "",
		           distanceFileName: "",
		           name: "",
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
		           args: args {
		           trainFileName: "",
		           queryFileName: "",
		           groundTruthFileName: "",
		           distanceFileName: "",
		           name: "",
		           distanceType: "",
		           objectType: "",
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

			got := loadLargeData(
				test.args.trainFileName,
				test.args.queryFileName,
				test.args.groundTruthFileName,
				test.args.distanceFileName,
				test.args.name,
				test.args.distanceType,
				test.args.objectType,
			)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_Train(t *testing.T) {
	t.Parallel()
	type args struct {
		i int
	}
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
		           i: 0,
		       },
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got, err := d.Train(test.args.i)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_TrainSize(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.TrainSize()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_Query(t *testing.T) {
	t.Parallel()
	type args struct {
		i int
	}
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
		           i: 0,
		       },
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got, err := d.Query(test.args.i)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_QuerySize(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.QuerySize()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_Distance(t *testing.T) {
	t.Parallel()
	type args struct {
		i int
	}
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
		           i: 0,
		       },
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got, err := d.Distance(test.args.i)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_DistanceSize(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.DistanceSize()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_Neighbor(t *testing.T) {
	t.Parallel()
	type args struct {
		i int
	}
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
	}
	type want struct {
		want []int
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []int, err error) error {
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
		           i: 0,
		       },
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got, err := d.Neighbor(test.args.i)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_NeighborSize(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.NeighborSize()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_Dimension(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.Dimension()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_DistanceType(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.DistanceType()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_ObjectType(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.ObjectType()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_largeDataset_Name(t *testing.T) {
	t.Parallel()
	type fields struct {
		dataset     *dataset
		train       x1b.BillionScaleVectors
		query       x1b.BillionScaleVectors
		groundTruth [][]int
		distances   x1b.FloatVectors
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
		           dataset: dataset{},
		           train: nil,
		           query: nil,
		           groundTruth: nil,
		           distances: nil,
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
			d := &largeDataset{
				dataset:     test.fields.dataset,
				train:       test.fields.train,
				query:       test.fields.query,
				groundTruth: test.fields.groundTruth,
				distances:   test.fields.distances,
			}

			got := d.Name()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
