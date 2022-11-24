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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Queue
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Queue, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Queue, err error) error {
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PushInsert(t *testing.T) {
	type args struct {
		uuid   string
		vector []float32
		date   int64
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
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
		           uuid: "",
		           vector: nil,
		           date: 0,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           vector: nil,
		           date: 0,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			err := v.PushInsert(test.args.uuid, test.args.vector, test.args.date)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PushDelete(t *testing.T) {
	type args struct {
		uuid string
		date int64
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
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
		           uuid: "",
		           date: 0,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           date: 0,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			err := v.PushDelete(test.args.uuid, test.args.date)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopInsert(t *testing.T) {
	type args struct {
		ctx context.Context
		now int64
		f   func(uuid string, vector []float32) bool
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           f: nil,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           ctx: nil,
		           f: nil,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.RangePopInsert(test.args.ctx, test.args.now, test.args.f)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		now int64
		f   func(uuid string) bool
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           f: nil,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           ctx: nil,
		           f: nil,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.RangePopDelete(test.args.ctx, test.args.now, test.args.f)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_GetVector(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct {
		want  []float32
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 3000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return (nil, false) when the uiid dose not exist in uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want:  nil,
					want1: false,
				},
			}
		}(),
		func() test {
			uii := []index{
				{
					uuid: "146bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 3000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return (nil, false) when the uuid is empty",
				args: args{
					uuid: "",
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want:  nil,
					want1: false,
				},
			}
		}(),
		func() test {
			return test{
				name: "return (nil, false) when the uiim is empty",
				args: args{
					uuid: "146bbe1a-bc48-11eb-8529-0242ac130003",
				},
				fields: fields{},
				want: want{
					want:  nil,
					want1: false,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "446bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return (1, true) when the uiid dose not exist in udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want:  []float32{1},
					want1: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return (1, true) when the udim is empty",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want:  []float32{1},
					want1: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return (1, true) when the date of uiim is equal the date of udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want:  []float32{1},
					want1: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000001,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return (1, true) when the date of uiim is newer than the date of udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want:  []float32{1},
					want1: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   999999999,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return (nil, false) when the date of uiim is older than the date of udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want:  nil,
					want1: false,
				},
			}
		}(),
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got, got1 := v.GetVector(test.args.uuid)
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_IVExists(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 3000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return false when the uiid dose not exist in uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			uii := []index{
				{
					uuid: "146bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 3000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return false when the uuid is empty",
				args: args{
					uuid: "",
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			return test{
				name: "return false when the uiim is empty",
				args: args{
					uuid: "146bbe1a-bc48-11eb-8529-0242ac130003",
				},
				fields: fields{},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "446bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return true when the uiid dose not exist in udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return true when the udim is empty",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return true when the date of uiim is equal the date of udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000001,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return true when the date of uiim is newer than the date of udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   999999999,
				},
				{
					uuid:   "246bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   2000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return false when the date of uiim is older than the date of udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: false,
				},
			}
		}(),
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got := v.IVExists(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_DVExists(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			udk := []key{
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "346bbe1a-bc48-11eb-8529-0242ac130003",
					date: 3000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return false when the uiid dose not exist in udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					udim: udim,
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			udk := []key{
				{
					uuid: "146bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 3000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return false when the uuid is empty",
				args: args{
					uuid: "",
				},
				fields: fields{
					udim: udim,
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			return test{
				name: "return false when the udim is empty",
				args: args{
					uuid: "146bbe1a-bc48-11eb-8529-0242ac130003",
				},
				fields: fields{},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			uii := []index{
				{
					uuid:   "346bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "446bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   4000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return true when the uiid dose not exist in uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			return test{
				name: "return true when the uiim is empty",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					udim: udim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			udk := []key{
				{
					uuid: uuid,
					date: 1000000000,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "346bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   4000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return true when the date of udim is equal the date of uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			udk := []key{
				{
					uuid: uuid,
					date: 1000000001,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "346bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   4000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return true when the date of udim is newer than the date of uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "146bbe1a-bc48-11eb-8529-0242ac130003"
			udk := []key{
				{
					uuid: uuid,
					date: 999999999,
				},
				{
					uuid: "246bbe1a-bc48-11eb-8529-0242ac130003",
					date: 2000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			uii := []index{
				{
					uuid:   uuid,
					vector: []float32{1},
					date:   1000000000,
				},
				{
					uuid:   "346bbe1a-bc48-11eb-8529-0242ac130003",
					vector: []float32{2},
					date:   4000000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			return test{
				name: "return false when the date of udim is older than the date of uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
					udim: udim,
				},
				want: want{
					want: false,
				},
			}
		}(),
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got := v.DVExists(test.args.uuid)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_addInsert(t *testing.T) {
	type args struct {
		i index
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           i: index{},
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           i: index{},
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.addInsert(test.args.i)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_addDelete(t *testing.T) {
	type args struct {
		d key
	}
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           d: key{},
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           d: key{},
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.addDelete(test.args.d)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_IVQLen(t *testing.T) {
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct {
		wantL int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotL int) error {
		if !reflect.DeepEqual(gotL, w.wantL) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotL, w.wantL)
		}
		return nil
	}
	tests := []test{
		func() test {
			size := 0
			uii := make([]index, size)

			return test{
				name: "return 0 when the capacity and length is 0",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			size := 10
			uii := make([]index, size)

			return test{
				name: "return 10 when the capacity and length is 10",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			c := 10
			l := 5
			uii := make([]index, l, c)

			return test{
				name: "return 5 when the capacity is 10 and the length is 5",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: l,
				},
			}
		}(),
		func() test {
			iniLen := 5
			isrtSize := 2
			uii := make([]index, iniLen, 10)
			for i := 0; i < isrtSize; i++ {
				uii = append(uii, index{})
			}

			return test{
				name: "return 7 when the capacity is 10 and the initial length is 5 but the inserted size is 2",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: iniLen + isrtSize,
				},
			}
		}(),
		func() test {
			size := 10
			uii := make([]index, 0, size)
			for i := 0; i < size; i++ {
				uii = append(uii, index{})
			}

			return test{
				name: "return 10 when the inserted size is 10",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			insertSize := 5
			size := 10
			uii := make([]index, 0, size)
			for i := 0; i < insertSize; i++ {
				uii = append(uii, index{})
			}

			return test{
				name: "return 5 when the capacity is 10 and the inserted size is 5",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: insertSize,
				},
			}
		}(),
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			gotL := v.IVQLen()
			if err := checkFunc(test.want, gotL); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_DVQLen(t *testing.T) {
	type fields struct {
		uii      []index
		uiim     uiim
		udk      []key
		udim     udim
		eg       errgroup.Group
		iBufSize int
		dBufSize int
	}
	type want struct {
		wantL int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotL int) error {
		if !reflect.DeepEqual(gotL, w.wantL) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotL, w.wantL)
		}
		return nil
	}
	tests := []test{
		func() test {
			size := 0
			udk := make([]key, size)

			return test{
				name: "return 0 when the capacity and length is 0",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			size := 10
			udk := make([]key, size)

			return test{
				name: "return 10 when the capacity and length is 10",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			c := 10
			l := 5
			udk := make([]key, l, c)

			return test{
				name: "return 5 when the capacity is 10 and the length is 5",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: l,
				},
			}
		}(),
		func() test {
			iniLen := 5
			isrtSize := 2
			udk := make([]key, iniLen, 10)
			for i := 0; i < isrtSize; i++ {
				udk = append(udk, key{})
			}

			return test{
				name: "return 7 when the capacity is 10 and the initial length is 5 but the inserted size is 2",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: iniLen + isrtSize,
				},
			}
		}(),
		func() test {
			size := 10
			udk := make([]key, 0, size)
			for i := 0; i < size; i++ {
				udk = append(udk, key{})
			}

			return test{
				name: "return 10 when the inserted size is 10",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			insertSize := 5
			size := 10
			udk := make([]key, 0, size)
			for i := 0; i < insertSize; i++ {
				udk = append(udk, key{})
			}

			return test{
				name: "return 5 when the capacity is 10 and the inserted size is 5",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: insertSize,
				},
			}
		}(),
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
			v := &vqueue{
				uii:      test.fields.uii,
				uiim:     test.fields.uiim,
				udk:      test.fields.udk,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			gotL := v.DVQLen()
			if err := checkFunc(test.want, gotL); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
