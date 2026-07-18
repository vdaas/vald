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
package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestCalcResource(t *testing.T) {
	tests := []struct {
		name      string
		divs      []float64
		value     int64
		ratio     float64
		wantMilli int64
	}{
		// ratio=0.5 is exactly representable in float64 (=2^-1), so results are exact.
		{"no div", nil, 4, 0.5, 2000},
		{"div=2", []float64{2}, 4, 0.5, 1000},
		{"div=4", []float64{4}, 8, 0.5, 1000},
		// ratio=0.25 is also exact.
		{"0.25 no div", nil, 8, 0.25, 2000},
		{"0.25 div=2", []float64{2}, 8, 0.25, 1000},
		// Scale-up: large CPU values typical of real nodes (int64 safe range).
		{"16 cores × 0.5", nil, 16, 0.5, 8000},
		{"16 cores × 0.5 / 2", []float64{2}, 16, 0.5, 4000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcResource(tt.value, tt.ratio, tt.divs...)
			assert.Equal(t, tt.wantMilli, got.MilliValue())
		})
	}
}

func TestNormalizeResourceList_CPU(t *testing.T) {
	tests := []struct {
		name       string
		wantStr    string
		inputMilli int64
	}{
		{"whole core: 3000m → 3", "3", 3000},
		{"whole core: 1000m → 1", "1", 1000},
		{"fractional: 1500m stays", "1500m", 1500},
		{"fractional: 600m stays", "600m", 600},
		{"fractional: 4800m → 4800m (not whole)", "4800m", 4800},
		{"whole core: 6000m → 6", "6", 6000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := resource.Quantity{}
			q.SetMilli(tt.inputMilli)
			rl := NormalizeResourceList(v1.ResourceList{v1.ResourceCPU: q})
			got := rl[v1.ResourceCPU]
			assert.Equal(t, tt.wantStr, got.String())
		})
	}
}

func TestNormalizeResourceList_Memory(t *testing.T) {
	tests := []struct {
		name       string
		inputBytes int64 // passed via SetMilli(inputBytes * 1000) to simulate CalcResource output
		wantBytes  int64
	}{
		// Values should be rounded up to nearest Mega (10^6 bytes).
		{"3000M exact", 3_000_000_000, 3_000_000_000},
		{"rounds up to next M", 2_999_999_001, 3_000_000_000},
		{"1500M exact", 1_500_000_000, 1_500_000_000},
		{"1500M with sub-M remainder rounds up", 1_499_999_001, 1_500_000_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := resource.Quantity{}
			q.SetMilli(tt.inputBytes * 1000)
			rl := NormalizeResourceList(v1.ResourceList{v1.ResourceMemory: q})
			got := rl[v1.ResourceMemory]
			assert.Equal(t, tt.wantBytes, got.Value())
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNormalizeResourceList(t *testing.T) {
// 	type args struct {
// 		rl v1.ResourceList
// 	}
// 	type want struct {
// 		want v1.ResourceList
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, v1.ResourceList) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got v1.ResourceList) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           rl:nil,
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
// 		           rl:nil,
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
// 			got := NormalizeResourceList(test.args.rl)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_getScale(t *testing.T) {
// 	type args struct {
// 		quantity resource.Quantity
// 	}
// 	type want struct {
// 		want resource.Scale
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, resource.Scale) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got resource.Scale) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           quantity:nil,
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
// 		           quantity:nil,
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
// 			got := getScale(test.args.quantity)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBuildTopologySpreadConstraint(t *testing.T) {
// 	type args struct {
// 		componentLabel string
// 	}
// 	type want struct {
// 		want v1.TopologySpreadConstraint
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, v1.TopologySpreadConstraint) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got v1.TopologySpreadConstraint) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           componentLabel:"",
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
// 		           componentLabel:"",
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
// 			got := BuildTopologySpreadConstraint(test.args.componentLabel)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
