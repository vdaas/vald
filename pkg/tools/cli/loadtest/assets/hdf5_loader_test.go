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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"gonum.org/v1/hdf5"
)

func Test_loadFloat32(t *testing.T) {
	t.Parallel()
	type args struct {
		dset    *hdf5.Dataset
		npoints int
		row     int
		dim     int
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
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
		           dset: nil,
		           npoints: 0,
		           row: 0,
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
		           dset: nil,
		           npoints: 0,
		           row: 0,
		           dim: 0,
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

			got, err := loadFloat32(test.args.dset, test.args.npoints, test.args.row, test.args.dim)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_loadInt(t *testing.T) {
	t.Parallel()
	type args struct {
		dset    *hdf5.Dataset
		npoints int
		row     int
		dim     int
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
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
		           dset: nil,
		           npoints: 0,
		           row: 0,
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
		           dset: nil,
		           npoints: 0,
		           row: 0,
		           dim: 0,
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

			got, err := loadInt(test.args.dset, test.args.npoints, test.args.row, test.args.dim)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_loadDataset(t *testing.T) {
	t.Parallel()
	type args struct {
		file *hdf5.File
		name string
		f    loaderFunc
	}
	type want struct {
		wantDim int
		wantVec interface{}
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, int, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotDim int, gotVec interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotDim, w.wantDim) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDim, w.wantDim)
		}
		if !reflect.DeepEqual(gotVec, w.wantVec) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           file: nil,
		           name: "",
		           f: nil,
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
		           file: nil,
		           name: "",
		           f: nil,
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

			gotDim, gotVec, err := loadDataset(test.args.file, test.args.name, test.args.f)
			if err := checkFunc(test.want, gotDim, gotVec, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		wantTrain     [][]float32
		wantTest      [][]float32
		wantDistances [][]float32
		wantNeighbors [][]int
		wantDim       int
		err           error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32, [][]float32, [][]float32, [][]int, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotTrain [][]float32, gotTest [][]float32, gotDistances [][]float32, gotNeighbors [][]int, gotDim int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotTrain, w.wantTrain) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTrain, w.wantTrain)
		}
		if !reflect.DeepEqual(gotTest, w.wantTest) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTest, w.wantTest)
		}
		if !reflect.DeepEqual(gotDistances, w.wantDistances) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDistances, w.wantDistances)
		}
		if !reflect.DeepEqual(gotNeighbors, w.wantNeighbors) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNeighbors, w.wantNeighbors)
		}
		if !reflect.DeepEqual(gotDim, w.wantDim) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDim, w.wantDim)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           path: "",
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
		           path: "",
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

			gotTrain, gotTest, gotDistances, gotNeighbors, gotDim, err := Load(test.args.path)
			if err := checkFunc(test.want, gotTrain, gotTest, gotDistances, gotNeighbors, gotDim, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
