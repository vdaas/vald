// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"path/filepath"
	"testing"

	"gonum.org/v1/hdf5"
)

// writeHDF5 writes float32 datasets with the given shapes into a new hdf5 file.
func writeHDF5(t *testing.T, path string, shapes map[string][]uint) {
	t.Helper()
	f, err := hdf5.CreateFile(path, hdf5.F_ACC_TRUNC)
	if err != nil {
		t.Fatalf("failed to create hdf5 file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("failed to close hdf5 file: %v", err)
		}
	}()
	for name, shape := range shapes {
		sp, err := hdf5.CreateSimpleDataspace(shape, nil)
		if err != nil {
			t.Fatalf("failed to create dataspace for %s: %v", name, err)
		}
		d, err := f.CreateDataset(name, hdf5.T_NATIVE_FLOAT, sp)
		if err != nil {
			t.Fatalf("failed to create dataset %s: %v", name, err)
		}
		n := uint(1)
		for _, s := range shape {
			n *= s
		}
		data := make([]float32, n)
		for i := range data {
			data[i] = float32(i)
		}
		if err := d.Write(&data); err != nil {
			t.Fatalf("failed to write dataset %s: %v", name, err)
		}
		if err := d.Close(); err != nil {
			t.Fatalf("failed to close dataset %s: %v", name, err)
		}
		if err := sp.Close(); err != nil {
			t.Fatalf("failed to close dataspace for %s: %v", name, err)
		}
	}
}

func Test_load_datasetRank(t *testing.T) {
	// not parallel: the hdf5 C library is not guaranteed to be built thread-safe.
	// dataset names load reads from the hdf5 file.
	const (
		trainDataset = "train"
		testDataset  = "test"
	)
	tests := []struct {
		shapes    map[string][]uint
		name      string
		wantTrain [2]int
		wantTest  [2]int
		wantErr   bool
	}{
		{
			name: "returns error instead of panicking on a 1-D train dataset",
			shapes: map[string][]uint{
				trainDataset: {6},
				testDataset:  {2, 3},
			},
			wantErr: true,
		},
		{
			name: "returns error instead of panicking on a 3-D test dataset",
			shapes: map[string][]uint{
				trainDataset: {2, 3},
				testDataset:  {1, 2, 3},
			},
			wantErr: true,
		},
		{
			name: "loads 2-D datasets",
			shapes: map[string][]uint{
				trainDataset: {4, 3},
				testDataset:  {2, 3},
			},
			wantTrain: [2]int{4, 3},
			wantTest:  [2]int{2, 3},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			path := filepath.Join(tt.TempDir(), "data.hdf5")
			writeHDF5(tt, path, test.shapes)

			var (
				train, tst [][]float32
				err        error
			)
			func() {
				defer func() {
					if r := recover(); r != nil {
						tt.Fatalf("load panicked: %v", r)
					}
				}()
				train, tst, err = load(path)
			}()
			if test.wantErr {
				if err == nil {
					tt.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				tt.Fatalf("unexpected error: %v", err)
			}
			if len(train) != test.wantTrain[0] {
				tt.Fatalf("train rows: got %d, want %d", len(train), test.wantTrain[0])
			}
			if len(train[0]) != test.wantTrain[1] {
				tt.Errorf("train dim: got %d, want %d", len(train[0]), test.wantTrain[1])
			}
			if len(tst) != test.wantTest[0] {
				tt.Fatalf("test rows: got %d, want %d", len(tst), test.wantTest[0])
			}
			if len(tst[0]) != test.wantTest[1] {
				tt.Errorf("test dim: got %d, want %d", len(tst[0]), test.wantTest[1])
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_main(t *testing.T) {
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			main()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_run(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		load    bool
// 		path    string
// 		dim     int
// 		vectors [][]float32
// 		ids     []uint
// 		dur     time.Duration
// 		output  func(header string)
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           load:false,
// 		           path:"",
// 		           dim:0,
// 		           vectors:nil,
// 		           ids:nil,
// 		           dur:nil,
// 		           output:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           load:false,
// 		           path:"",
// 		           dim:0,
// 		           vectors:nil,
// 		           ids:nil,
// 		           dur:nil,
// 		           output:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			run(test.args.ctx, test.args.load, test.args.path, test.args.dim, test.args.vectors, test.args.ids, test.args.dur, test.args.output)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_sleep(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		duration time.Duration
// 		limit    time.Duration
// 		fn       func()
// 		efn      func()
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           duration:nil,
// 		           limit:nil,
// 		           fn:nil,
// 		           efn:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           duration:nil,
// 		           limit:nil,
// 		           fn:nil,
// 		           efn:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			sleep(test.args.ctx, test.args.duration, test.args.limit, test.args.fn, test.args.efn)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_load(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		wantTrain [][]float32
// 		wantTest  [][]float32
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, [][]float32, [][]float32, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotTrain [][]float32, gotTest [][]float32, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotTrain, w.wantTrain) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTrain, w.wantTrain)
// 		}
// 		if !reflect.DeepEqual(gotTest, w.wantTest) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTest, w.wantTest)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotTrain, gotTest, err := load(test.args.path)
// 			if err := checkFunc(test.want, gotTrain, gotTest, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// run used to log an NGT initialization failure and return nil, so main treated a failed phase as success.
func Test_run_ngtInitError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		load bool
	}{
		{name: "new", load: false},
		{name: "load", load: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Dimension 1 is below algorithm.MinimumVectorDimensionSize, so ngt.New/Load fail before touching any index file.
			if err := run(t.Context(), tc.load, 1, nil, nil, 0, func(string) {}); err == nil {
				t.Fatal("run returned nil for a failed NGT initialization")
			}
		})
	}
}
