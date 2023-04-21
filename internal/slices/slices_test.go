// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package slices

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestSort(t *testing.T) {
	type args struct {
		x []int
	}
	type want struct {
		want []int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func([]int, want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(got []int, w want) error {
		if len(got) != len(w.want) {
			return errors.New("len not match")
		}
		for i := 0; i < len(got); i++ {
			if got[i] != w.want[i] {
				return errors.New("slice not sorted")
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "success to sort 10 elements",
			args: args{
				x: []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
			},
			want: want{
				want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		{
			name: "success to sort 1 elements",
			args: args{
				x: []int{0},
			},
			want: want{
				want: []int{0},
			},
		},
		{
			name: "success to sort 0 elements",
			args: args{
				x: []int{},
			},
			want: want{
				want: []int{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			Sort(test.args.x)
			if err := checkFunc(test.args.x, test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSortFunc(t *testing.T) {
	type args struct {
		x    []int
		less func(left, right int) bool
	}
	type want struct {
		want []int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func([]int, want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}

	defaultLessFn := func(l, r int) bool {
		return l < r
	}
	defaultCheckFunc := func(got []int, w want) error {
		if len(got) != len(w.want) {
			return errors.New("len not match")
		}
		for i := 0; i < len(got); i++ {
			if got[i] != w.want[i] {
				return errors.New("slice not sorted")
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "success to sort 10 elements",
			args: args{
				x:    []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
				less: defaultLessFn,
			},
			want: want{
				want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		{
			name: "success to sort 1 elements",
			args: args{
				x:    []int{0},
				less: defaultLessFn,
			},
			want: want{
				want: []int{0},
			},
		},
		{
			name: "success to sort 0 elements",
			args: args{
				x:    []int{},
				less: defaultLessFn,
			},
			want: want{
				want: []int{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			SortFunc(test.args.x, test.args.less)
			if err := checkFunc(test.args.x, test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSortStableFunc(t *testing.T) {
	type args struct {
		x    []int
		less func(left, right int) bool
	}
	type want struct {
		want []int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func([]int, want) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}

	defaultLessFn := func(l, r int) bool {
		return l < r
	}
	defaultCheckFunc := func(got []int, w want) error {
		if len(got) != len(w.want) {
			return errors.New("len not match")
		}
		for i := 0; i < len(got); i++ {
			if got[i] != w.want[i] {
				return errors.New("slice not sorted")
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "success to sort 10 elements",
			args: args{
				x:    []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
				less: defaultLessFn,
			},
			want: want{
				want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		{
			name: "success to sort 1 elements",
			args: args{
				x:    []int{0},
				less: defaultLessFn,
			},
			want: want{
				want: []int{0},
			},
		},
		{
			name: "success to sort 0 elements",
			args: args{
				x:    []int{},
				less: defaultLessFn,
			},
			want: want{
				want: []int{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			SortStableFunc(test.args.x, test.args.less)
			if err := checkFunc(test.args.x, test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		x    []int
		less func(left, right int) bool
	}
	type want struct {
		want []int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []int) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got []int) error {
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
		           x:nil,
		           less:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           x:nil,
		           less:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := RemoveDuplicates(test.args.x, test.args.less)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
