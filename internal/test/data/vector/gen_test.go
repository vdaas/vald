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
package vector

import (
	"math"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestFloat32VectorGenerator(t *testing.T) {
	type args struct {
		d   Distribution
		n   int
		dim int
	}
	type want struct {
		n   int
		dim int
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(args, want, Float32VectorGeneratorFunc, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(a args, w want, got Float32VectorGeneratorFunc, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if got != nil {
			vectors := got(a.n, a.dim)
			if len(vectors) != w.n && len(vectors[0]) != w.dim {
				return errors.Errorf("got: \"%d\",\"%d\"\n\t\t\t\twant: \"%d\",\"%d\"", len(vectors), len(vectors[0]), w.n, w.dim)
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "success generating gaussian distributed random vector",
			args: args{
				d:   Gaussian,
				n:   20,
				dim: 10,
			},
			want: want{
				n:   20,
				dim: 10,
				err: nil,
			},
		},
		{
			name: "success generating uniform distributed random vector",
			args: args{
				d:   Uniform,
				n:   20,
				dim: 10,
			},
			want: want{
				n:   20,
				dim: 10,
				err: nil,
			},
		},
		{
			name: "fail generating random vector with unknown distribution",
			args: args{
				d: -1,
			},
			want: want{
				err: ErrUnknownDistribution,
			},
		},
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

			got, err := Float32VectorGenerator(test.args.d)
			if err := checkFunc(test.args, test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestUint8VectorGenerator(t *testing.T) {
	type args struct {
		d   Distribution
		n   int
		dim int
	}
	type want struct {
		n   int
		dim int
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(args, want, Uint8VectorGeneratorFunc, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(a args, w want, got Uint8VectorGeneratorFunc, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if got != nil {
			vectors := got(a.n, a.dim)
			if len(vectors) != w.n && len(vectors[0]) != w.dim {
				return errors.Errorf("got: \"%d\",\"%d\"\n\t\t\t\twant: \"%d\",\"%d\"", len(vectors), len(vectors[0]), w.n, w.dim)
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "success generating gaussian distributed random vector",
			args: args{
				d:   Gaussian,
				n:   200,
				dim: 100,
			},
			want: want{
				n:   200,
				dim: 100,
				err: nil,
			},
		},
		{
			name: "success generating uniform distributed random vector",
			args: args{
				d:   Uniform,
				n:   200,
				dim: 100,
			},
			want: want{
				n:   200,
				dim: 100,
				err: nil,
			},
		},
		{
			name: "fail generating random vector with unknown distribution",
			args: args{
				d: -1,
			},
			want: want{
				err: ErrUnknownDistribution,
			},
		},
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

			got, err := Uint8VectorGenerator(test.args.d)
			if err := checkFunc(test.args, test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_float32VectorGenerator(t *testing.T) {
	type args struct {
		n   int
		dim int
		gen func() float32
	}
	type want struct {
		wantRet [][]float32
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRet [][]float32) error {
		if !reflect.DeepEqual(gotRet, w.wantRet) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRet, w.wantRet)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           n: 0,
		           dim: 0,
		           gen: nil,
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
		           n: 0,
		           dim: 0,
		           gen: nil,
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

			gotRet := float32VectorGenerator(test.args.n, test.args.dim, test.args.gen)
			if err := checkFunc(test.want, gotRet); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestUniformDistributedFloat32VectorGenerator(t *testing.T) {
	type args struct {
		n   int
		dim int
	}
	type want struct {
		want [][]float32
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]float32) error {
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
		           n: 0,
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
		           n: 0,
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

			got := UniformDistributedFloat32VectorGenerator(test.args.n, test.args.dim)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGaussianDistributedFloat32VectorGenerator(t *testing.T) {
	type args struct {
		n   int
		dim int
	}
	type want struct {
		want [][]float32
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]float32) error {
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
		           n: 0,
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
		           n: 0,
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

			got := GaussianDistributedFloat32VectorGenerator(test.args.n, test.args.dim)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_uint8VectorGenerator(t *testing.T) {
	type args struct {
		n   int
		dim int
		gen func() uint8
	}
	type want struct {
		wantRet [][]uint8
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]uint8) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRet [][]uint8) error {
		if !reflect.DeepEqual(gotRet, w.wantRet) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRet, w.wantRet)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           n: 0,
		           dim: 0,
		           gen: nil,
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
		           n: 0,
		           dim: 0,
		           gen: nil,
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

			gotRet := uint8VectorGenerator(test.args.n, test.args.dim, test.args.gen)
			if err := checkFunc(test.want, gotRet); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestUniformDistributedUint8VectorGenerator(t *testing.T) {
	type args struct {
		n   int
		dim int
	}
	type want struct {
		want [][]uint8
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]uint8) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]uint8) error {
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
		           n: 0,
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
		           n: 0,
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

			got := UniformDistributedUint8VectorGenerator(test.args.n, test.args.dim)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGaussianDistributedUint8VectorGenerator(t *testing.T) {
	type args struct {
		n   int
		dim int
	}
	type want struct {
		want [][]uint8
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]uint8) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]uint8) error {
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
		           n: 0,
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
		           n: 0,
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

			got := GaussianDistributedUint8VectorGenerator(test.args.n, test.args.dim)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gaussianDistributedUint8VectorGenerator(t *testing.T) {
	type args struct {
		n     int
		dim   int
		mean  float64
		sigma float64
	}
	type want struct {
		want [][]uint8
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]uint8) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]uint8) error {
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
		           n: 0,
		           dim: 0,
		           mean: 0,
		           sigma: 0,
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
		           n: 0,
		           dim: 0,
		           mean: 0,
		           sigma: 0,
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

			got := gaussianDistributedUint8VectorGenerator(test.args.n, test.args.dim, test.args.mean, test.args.sigma)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGenF32Vec(t *testing.T) {
	type args struct {
		dist Distribution
		num  int
		dim  int
	}
	type want struct {
		wantLen int
		wantDim int
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if len(got) != w.wantLen {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twantLen: \"%#v\"", got, w.wantLen)
		}
		for _, vec := range got {
			if len(vec) != w.wantDim {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twantDim: \"%#v\"", got, w.wantDim)
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "return 1 generated float32 vector",
			args: args{
				dist: Gaussian,
				num:  1,
				dim:  5,
			},
			want: want{
				wantLen: 1,
				wantDim: 5,
			},
		},
		{
			name: "return 5 generated float32 vector",
			args: args{
				dist: Gaussian,
				num:  5,
				dim:  5,
			},
			want: want{
				wantLen: 5,
				wantDim: 5,
			},
		},
		{
			name: "return no generated float32 vector with num 0",
			args: args{
				dist: Gaussian,
				num:  0,
				dim:  5,
			},
			want: want{
				wantLen: 0,
			},
		},
		{
			name: "return generated float32 vector with dim 0",
			args: args{
				dist: Gaussian,
				num:  1,
				dim:  0,
			},
			want: want{
				wantLen: 1,
				wantDim: 0,
			},
		},
		// max dim and max num test is ignored due to test timeout
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

			got, err := GenF32Vec(test.args.dist, test.args.num, test.args.dim)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGenUint8Vec(t *testing.T) {
	type args struct {
		dist Distribution
		num  int
		dim  int
	}
	type want struct {
		wantLen int
		wantDim int
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if len(got) != w.wantLen {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twantLen: \"%#v\"", got, w.wantLen)
		}
		for _, vec := range got {
			if len(vec) != w.wantDim {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twantDim: \"%#v\"", got, w.wantDim)
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "return 1 generated uint8 vector",
			args: args{
				dist: Gaussian,
				num:  1,
				dim:  5,
			},
			want: want{
				wantLen: 1,
				wantDim: 5,
			},
		},
		{
			name: "return 5 generated uint8 vector",
			args: args{
				dist: Gaussian,
				num:  5,
				dim:  5,
			},
			want: want{
				wantLen: 5,
				wantDim: 5,
			},
		},
		{
			name: "return no generated uint8 vector with num 0",
			args: args{
				dist: Gaussian,
				num:  0,
				dim:  5,
			},
			want: want{
				wantLen: 0,
			},
		},
		{
			name: "return generated uint8 vector with dim 0",
			args: args{
				dist: Gaussian,
				num:  1,
				dim:  0,
			},
			want: want{
				wantLen: 1,
				wantDim: 0,
			},
		},
		// max dim and max num test is ignored due to test timeout
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

			got, err := GenUint8Vec(test.args.dist, test.args.num, test.args.dim)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGenSameValueVec(t *testing.T) {
	type args struct {
		size int
		val  float32
	}
	type want struct {
		want []float32
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []float32) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			val := float32(1)
			return test{
				name: "return same value vector with size 1",
				args: args{
					size: 1,
					val:  val,
				},
				want: want{
					want: []float32{
						val,
					},
				},
			}
		}(),
		func() test {
			val := float32(1)
			return test{
				name: "return same value vector with size 5",
				args: args{
					size: 5,
					val:  val,
				},
				want: want{
					want: []float32{
						val, val, val, val, val,
					},
				},
			}
		}(),
		func() test {
			val := float32(1)
			return test{
				name: "return same value vector with size 0",
				args: args{
					size: 0,
					val:  val,
				},
				want: want{
					want: []float32{},
				},
			}
		}(),
		func() test {
			val := float32(math.SmallestNonzeroFloat32)
			return test{
				name: "return same value vector with min value",
				args: args{
					size: 1,
					val:  val,
				},
				want: want{
					want: []float32{
						val,
					},
				},
			}
		}(),
		func() test {
			val := float32(math.MaxFloat32)
			return test{
				name: "return same value vector with max value",
				args: args{
					size: 1,
					val:  val,
				},
				want: want{
					want: []float32{
						val,
					},
				},
			}
		}(),
		// max size test is ignored due to test timeout
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

			got := GenSameValueVec(test.args.size, test.args.val)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestConvertVectorsUint8ToFloat32(t *testing.T) {
	type args struct {
		vectors [][]uint8
	}
	type want struct {
		wantRet [][]float32
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, [][]float32) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRet [][]float32) error {
		if !reflect.DeepEqual(gotRet, w.wantRet) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRet, w.wantRet)
		}
		return nil
	}
	tests := []test{
		{
			name: "convert 1 float32 vector from uint8 vector",
			args: args{
				vectors: [][]uint8{
					{
						1, 2, 3,
					},
				},
			},
			want: want{
				wantRet: [][]float32{
					{
						1, 2, 3,
					},
				},
			},
		},
		{
			name: "convert 3 float32 vector from uint8 vector",
			args: args{
				vectors: [][]uint8{
					{
						1, 2, 3,
					},
					{
						4, 5, 6,
					},
					{
						7, 8, 9,
					},
				},
			},
			want: want{
				wantRet: [][]float32{
					{
						1, 2, 3,
					},
					{
						4, 5, 6,
					},
					{
						7, 8, 9,
					},
				},
			},
		},
		{
			name: "convert 0 value vector",
			args: args{
				vectors: [][]uint8{
					{
						0, 0, 0,
					},
				},
			},
			want: want{
				wantRet: [][]float32{
					{
						0, 0, 0,
					},
				},
			},
		},
		{
			name: "convert max value vector",
			args: args{
				vectors: [][]uint8{
					{
						math.MaxUint8, math.MaxUint8, math.MaxUint8,
					},
				},
			},
			want: want{
				wantRet: [][]float32{
					{
						math.MaxUint8, math.MaxUint8, math.MaxUint8,
					},
				},
			},
		},
		{
			name: "return empty slice when vectors is empty",
			args: args{
				vectors: [][]uint8{},
			},
			want: want{
				wantRet: [][]float32{},
			},
		},
		{
			name: "return empty slice when vectors is nil",
			args: args{
				vectors: nil,
			},
			want: want{
				wantRet: [][]float32{},
			},
		},
		{
			name: "return empty vector when vector is empty",
			args: args{
				vectors: [][]uint8{
					{},
				},
			},
			want: want{
				wantRet: [][]float32{
					{},
				},
			},
		},
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

			gotRet := ConvertVectorsUint8ToFloat32(test.args.vectors)
			if err := checkFunc(test.want, gotRet); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
