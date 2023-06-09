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
package grpc

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/test/data/strings"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestPairingHeap(t *testing.T) {
	type test struct {
		name      string
		args      []*DistPayload
		checkFunc func([]*DistPayload) error
	}
	defaultCheckFunc := func(got []*DistPayload) error {
		last := got[0].distance
		for i, current := range got {
			cmp := last.Cmp(current.distance)
			if cmp > 0 {
				gids := make([]string, 0, len(got))
				for i, dp := range got {
					gids = append(gids, fmt.Sprintf("%d:%s:%f", i, dp.raw.GetId(), dp.raw.GetDistance()))
				}
				return errors.Errorf(
					"unsorted return detected cmp: %d, last: %s, current: %s, got: idx=%d,\tid=%s,\traw[\"%#v\"]",
					cmp,
					last.String(),
					current.distance.String(),
					i,
					got[i].raw.GetId(),
					gids,
				)
			}
			last = current.distance
		}
		return nil
	}
	tests := []test{
		func() test {
			dl := 100
			ods := make([]*payload.Object_Distance, 0, dl)
			for i := 0; i < dl; i++ {
				ods = append(ods, &payload.Object_Distance{
					Id:       strings.Random(12),
					Distance: rand.Float32(),
				})
			}
			dps := make([]*DistPayload, 0, dl) // random DistPayload arguments for heap insert data
			for _, dist := range ods {
				dps = append(dps, &DistPayload{raw: dist, distance: big.NewFloat(float64(dist.GetDistance()))})
			}
			return test{
				name:      "check ExtractMin results are sorted",
				args:      dps,
				checkFunc: defaultCheckFunc,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			ph := new(PairingHeap)
			for _, dp := range test.args {
				ph = ph.Insert(dp)
			}
			result := make([]*DistPayload, 0, len(test.args))
			for !ph.IsEmpty() {
				var min *DistPayload
				min, ph = ph.ExtractMin()
				result = append(result, min)
			}
			if err := checkFunc(result); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

func TestPairingHeap_IsEmpty(t *testing.T) {
	type fields struct {
		DistPayload *DistPayload
		Children    []*PairingHeap
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           DistPayload:DistPayload{},
		           Children:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			ph := &PairingHeap{
				DistPayload: test.fields.DistPayload,
				Children:    test.fields.Children,
			}

			got := ph.IsEmpty()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestPairingHeap_Insert(t *testing.T) {
	type args struct {
		dp *DistPayload
	}
	type fields struct {
		DistPayload *DistPayload
		Children    []*PairingHeap
	}
	type want struct {
		want *PairingHeap
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *PairingHeap) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *PairingHeap) error {
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
		           dp:DistPayload{},
		       },
		       fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
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
		           dp:DistPayload{},
		           },
		           fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
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
			ph := &PairingHeap{
				DistPayload: test.fields.DistPayload,
				Children:    test.fields.Children,
			}

			got := ph.Insert(test.args.dp)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestPairingHeap_Merge(t *testing.T) {
	type args struct {
		h2 *PairingHeap
	}
	type fields struct {
		DistPayload *DistPayload
		Children    []*PairingHeap
	}
	type want struct {
		want *PairingHeap
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *PairingHeap) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *PairingHeap) error {
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
		           h2:PairingHeap{},
		       },
		       fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
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
		           h2:PairingHeap{},
		           },
		           fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
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
			ph := &PairingHeap{
				DistPayload: test.fields.DistPayload,
				Children:    test.fields.Children,
			}

			got := ph.Merge(test.args.h2)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestPairingHeap_ExtractMin(t *testing.T) {
	type fields struct {
		DistPayload *DistPayload
		Children    []*PairingHeap
	}
	type want struct {
		want  *DistPayload
		want1 *PairingHeap
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *DistPayload, *PairingHeap) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *DistPayload, got1 *PairingHeap) error {
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
		       fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			ph := &PairingHeap{
				DistPayload: test.fields.DistPayload,
				Children:    test.fields.Children,
			}

			got, got1 := ph.ExtractMin()
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestPairingHeap_mergePairs(t *testing.T) {
	type args struct {
		pairs []*PairingHeap
	}
	type fields struct {
		DistPayload *DistPayload
		Children    []*PairingHeap
	}
	type want struct {
		want *PairingHeap
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *PairingHeap) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *PairingHeap) error {
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
		           pairs:nil,
		       },
		       fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
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
		           pairs:nil,
		           },
		           fields: fields {
		           DistPayload:DistPayload{},
		           Children:nil,
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
			ph := &PairingHeap{
				DistPayload: test.fields.DistPayload,
				Children:    test.fields.Children,
			}

			got := ph.mergePairs(test.args.pairs)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
