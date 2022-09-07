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

// Package errors
package errors

import (
	"math"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestErrCreateProperty(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrCreateProperty error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to create property: ngt error"),
			},
		},
		{
			name: "return an ErrCreateProperty error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to create property"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCreateProperty(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIndexNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrIndexNotFound error",
			want: want{
				want: New("index not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrIndexNotFound
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIndexLoadTimeout(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrIndexLoadTimeout error",
			want: want{
				want: New("index load timeout"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrIndexLoadTimeout
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidDimensionSize(t *testing.T) {
	type args struct {
		current int
		limit   int
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrInvalidDimensionSize error when current is 10 and limit is 5",
			args: args{
				current: 10,
				limit:   5,
			},
			want: want{
				want: New("dimension size 10 is invalid, the supporting dimension size must be between 2 ~ 5"),
			},
		},
		{
			name: "return an ErrInvalidDimensionSize error when current is 0 and limit is 5",
			args: args{
				current: 0,
				limit:   5,
			},
			want: want{
				want: New("dimension size 0 is invalid, the supporting dimension size must be between 2 ~ 5"),
			},
		},
		{
			name: "return an ErrInvalidDimensionSize error when current is 10 and limit is 0",
			args: args{
				current: 10,
				limit:   0,
			},
			want: want{
				want: New("dimension size 10 is invalid, the supporting dimension size must be bigger than 2"),
			},
		},
		{
			name: "return an ErrInvalidDimensionSize error when current is 0 and limit is 0",
			args: args{
				current: 0,
				limit:   0,
			},
			want: want{
				want: New("dimension size 0 is invalid, the supporting dimension size must be bigger than 2"),
			},
		},
		{
			name: "return an ErrInvalidDimensionSize error when current and limit are the minimum value of int",
			args: args{
				current: int(math.MinInt64),
				limit:   int(math.MinInt64),
			},
			want: want{
				want: Errorf("dimension size %d is invalid, the supporting dimension size must be between 2 ~ %d", int(math.MinInt64), int(math.MinInt64)),
			},
		},
		{
			name: "return an ErrInvalidDimensionSize error when current and limit are the minimum value of int",
			args: args{
				current: int(math.MaxInt64),
				limit:   int(math.MaxInt64),
			},
			want: want{
				want: Errorf("dimension size %d is invalid, the supporting dimension size must be between 2 ~ %d", int(math.MaxInt64), int(math.MaxInt64)),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrInvalidDimensionSize(test.args.current, test.args.limit)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrInvalidUUID error when uuid is empty string",
			args: args{
				uuid: "",
			},
			want: want{
				want: New("uuid \"\" is invalid"),
			},
		},
		{
			name: "return an ErrInvalidUUID error when uuid is foo",
			args: args{
				uuid: "foo",
			},
			want: want{
				want: New("uuid \"foo\" is invalid"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrInvalidUUID(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrDimensionLimitExceed(t *testing.T) {
	type args struct {
		current int
		limit   int
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrDimensionLimitExceed error when current is 10 and limit is 5",
			args: args{
				current: 10,
				limit:   5,
			},
			want: want{
				want: New("supported dimension limit exceed:\trequired = 10,\tlimit = 5"),
			},
		},

		{
			name: "return an ErrDimensionLimitExceed error when current is 0 and limit is 0",
			args: args{
				current: 0,
				limit:   0,
			},
			want: want{
				want: New("supported dimension limit exceed:\trequired = 0,\tlimit = 0"),
			},
		},
		{
			name: "return an ErrDimensionLimitExceed error when current and limit are the minimum value of int",
			args: args{
				current: int(math.MinInt64),
				limit:   int(math.MinInt64),
			},
			want: want{
				want: Errorf("supported dimension limit exceed:\trequired = %d,\tlimit = %d", int(math.MinInt64), int(math.MinInt64)),
			},
		},
		{
			name: "return an ErrDimensionLimitExceed error when current and limit are the maximum value of int",
			args: args{
				current: int(math.MaxInt64),
				limit:   int(math.MaxInt64),
			},
			want: want{
				want: Errorf("supported dimension limit exceed:\trequired = %d,\tlimit = %d", int(math.MaxInt64), int(math.MaxInt64)),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrDimensionLimitExceed(test.args.current, test.args.limit)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIncompatibleDimensionSize(t *testing.T) {
	type args struct {
		req int
		dim int
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrIncompatibleDimensionSize error when req is 640 and dim is 720",
			args: args{
				req: 640,
				dim: 720,
			},
			want: want{
				want: New("incompatible dimension size detected\trequested: 640,\tconfigured: 720"),
			},
		},
		{
			name: "return an ErrIncompatibleDimensionSize error when req is empty and dim is 720",
			args: args{
				dim: 720,
			},
			want: want{
				want: New("incompatible dimension size detected\trequested: 0,\tconfigured: 720"),
			},
		},
		{
			name: "return an ErrIncompatibleDimensionSize error when req is 640",
			args: args{
				req: 640,
			},
			want: want{
				want: New("incompatible dimension size detected\trequested: 640,\tconfigured: 0"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrIncompatibleDimensionSize(test.args.req, test.args.dim)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUnsupportedObjectType(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrUnsupportedObjectType error",
			want: want{
				want: New("unsupported ObjectType"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrUnsupportedObjectType
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUnsupportedDistanceType(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrUnsupportedDistanceType error",
			want: want{
				want: New("unsupported DistanceType"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrUnsupportedDistanceType
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetDistanceType(t *testing.T) {
	type args struct {
		err      error
		distance string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrFailedToSetDistanceType error when err is ngt error and distance is l2",
			args: args{
				err:      New("ngt error"),
				distance: "l2",
			},
			want: want{
				want: New("failed to set distance type l2: ngt error"),
			},
		},
		{
			name: "return a wrapped ErrFailedToSetDistanceType error when err is ngt error and distance is empty",
			args: args{
				err:      New("ngt error"),
				distance: "",
			},
			want: want{
				want: New("failed to set distance type : ngt error"),
			},
		},
		{
			name: "return an ErrFailedToSetDistanceType error when err is nil and distance is cos",
			args: args{
				err:      nil,
				distance: "cos",
			},
			want: want{
				want: New("failed to set distance type cos"),
			},
		},
		{
			name: "return an ErrFailedToSetDistanceType error when err is nil and distance is empty",
			args: args{
				err:      nil,
				distance: "",
			},
			want: want{
				want: New("failed to set distance type "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrFailedToSetDistanceType(test.args.err, test.args.distance)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetObjectType(t *testing.T) {
	type args struct {
		err error
		t   string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrFailedToSetObjectType error when err is ngt error and t is Float",
			args: args{
				err: New("ngt error"),
				t:   "Float",
			},
			want: want{
				want: New("failed to set object type Float: ngt error"),
			},
		},
		{
			name: "return a wrapped ErrFailedToSetObjectType error when err is ngt error and t is empty",
			args: args{
				err: New("ngt error"),
				t:   "",
			},
			want: want{
				want: New("failed to set object type : ngt error"),
			},
		},
		{
			name: "return an ErrFailedToSetObjectType error when err is nil and t is Int",
			args: args{
				err: nil,
				t:   "Int",
			},
			want: want{
				want: New("failed to set object type Int"),
			},
		},
		{
			name: "return an ErrFailedToSetObjectType error when err is nil and t is empty",
			args: args{
				err: nil,
				t:   "",
			},
			want: want{
				want: New("failed to set object type "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrFailedToSetObjectType(test.args.err, test.args.t)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetDimension(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrFailedToSetDimension error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to set dimension: ngt error"),
			},
		},
		{
			name: "return an ErrFailedToSetDimension error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to set dimension"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrFailedToSetDimension(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetCreationEdgeSize(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrFailedToSetCreationEdgeSize error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to set creation edge size: ngt error"),
			},
		},
		{
			name: "return an ErrFailedToSetCreationEdgeSize error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to set creation edge size"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrFailedToSetCreationEdgeSize(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrFailedToSetSearchEdgeSize(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrFailedToSetSearchEdgeSize error when err is ngt error",
			args: args{
				err: New("ngt error"),
			},
			want: want{
				want: New("failed to set search edge size: ngt error"),
			},
		},
		{
			name: "return an ErrFailedToSetSearchEdgeSize error when err is nil",
			args: args{
				err: nil,
			},
			want: want{
				want: New("failed to set search edge size"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrFailedToSetSearchEdgeSize(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUncommittedIndexExists(t *testing.T) {
	type args struct {
		num uint64
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrUncommittedIndexExists error when num is 100",
			args: args{
				num: 100,
			},
			want: want{
				want: New("100 indexes are not committed"),
			},
		},

		{
			name: "return an ErrUncommittedIndexExists error when num is 0",
			args: args{
				num: 0,
			},
			want: want{
				want: New("0 indexes are not committed"),
			},
		},
		{
			name: "return an ErrUncommittedIndexExists error when num is the maximum value of uint64",
			args: args{
				num: math.MaxUint64,
			},
			want: want{
				want: Errorf("%d indexes are not committed", uint(math.MaxUint64)),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrUncommittedIndexExists(test.args.num)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUncommittedIndexNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrUncommittedIndexNotFound error",
			want: want{
				want: New("uncommitted indexes are not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrUncommittedIndexNotFound
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCAPINotImplemented(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrCAPINotImplemented error",
			want: want{
				want: New("not implemented in C API"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrCAPINotImplemented
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUUIDAlreadyExists(t *testing.T) {
	type args struct {
		uuid string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrUUIDAlreadyExists error when uuid is 550e8400-e29b-41d4",
			args: args{
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4 index already exists"),
			},
		},
		{
			name: "return an ErrUUIDAlreadyExists error when uuid is empty",
			args: args{
				uuid: "",
			},
			want: want{
				want: New("ngt uuid  index already exists"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrUUIDAlreadyExists(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrUUIDNotFound(t *testing.T) {
	type args struct {
		id uint32
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrUUIDNotFound error when id is 1234",
			args: args{
				id: 1234,
			},
			want: want{
				want: New("ngt object uuid 1234's metadata not found"),
			},
		},
		{
			name: "return an ErrUUIDNotFound error when id is the maximum value of uint32",
			args: args{
				id: math.MaxUint32,
			},
			want: want{
				want: Errorf("ngt object uuid %d's metadata not found", math.MaxUint32),
			},
		},
		{
			name: "return an ErrUUIDNotFound error when id is 0",
			args: args{
				id: 0,
			},
			want: want{
				want: New("ngt object uuid not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrUUIDNotFound(test.args.id)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrObjectIDNotFound(t *testing.T) {
	type args struct {
		uuid string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrObjectIDNotFound error when uuid is 550e8400-e29b-41d4.",
			args: args{
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4's object id not found"),
			},
		},
		{
			name: "return an ErrObjectIDNotFound error when uuid is empty.",
			args: args{
				uuid: "",
			},
			want: want{
				want: New("ngt uuid 's object id not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrObjectIDNotFound(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrObjectNotFound(t *testing.T) {
	type args struct {
		err  error
		uuid string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return a wrapped ErrObjectNotFound error when err is ngt error and uuid is 550e8400-e29b-41d4",
			args: args{
				err:  New("ngt error"),
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4's object not found: ngt error"),
			},
		},
		{
			name: "return a wrapped ErrObjectNotFound error when err is ngt error and uuid is empty",
			args: args{
				err:  New("ngt error"),
				uuid: "",
			},
			want: want{
				want: New("ngt uuid 's object not found: ngt error"),
			},
		},
		{
			name: "return an ErrObjectNotFound error when err is nil and uuid is 550e8400-e29b-41d4",
			args: args{
				err:  nil,
				uuid: "550e8400-e29b-41d4",
			},
			want: want{
				want: New("ngt uuid 550e8400-e29b-41d4's object not found"),
			},
		},
		{
			name: "return an ErrObjectNotFound error when err is nil and uuid is empty",
			args: args{
				err:  nil,
				uuid: "",
			},
			want: want{
				want: New("ngt uuid 's object not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrObjectNotFound(test.args.err, test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRemoveRequestedBeforeIndexing(t *testing.T) {
	type args struct {
		oid uint
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrRemoveRequestedBeforeIndexing error when oid is 100",
			args: args{
				oid: 100,
			},
			want: want{
				want: New("object id 100 is not indexed we cannot remove it"),
			},
		},
		{
			name: "return an ErrRemoveRequestedBeforeIndexing error when oid is 0",
			args: args{
				oid: 0,
			},
			want: want{
				want: New("object id 0 is not indexed we cannot remove it"),
			},
		},
		{
			name: "return an ErrRemoveRequestedBeforeIndexing error when oid is maximum value of uint",
			args: args{
				oid: uint(math.MaxUint64),
			},
			want: want{
				want: Errorf("object id %d is not indexed we cannot remove it", uint(math.MaxUint64)),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRemoveRequestedBeforeIndexing(test.args.oid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewNGTError(t *testing.T) {
	type args struct {
		msg string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !Is(err, w.err) {
			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           msg: "",
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
		           msg: "",
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := NewNGTError(test.args.msg)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNGTError_Error(t *testing.T) {
	type fields struct {
		Msg string
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
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Msg: "",
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
		           Msg: "",
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := NGTError{
				Msg: test.fields.Msg,
			}

			got := n.Error()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
