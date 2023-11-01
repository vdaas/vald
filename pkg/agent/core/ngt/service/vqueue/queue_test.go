//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetVector(t *testing.T) {
	type want struct {
		vec       []float32
		timestamp int64
		exists    bool
	}
	type test struct {
		name      string
		uuid      string
		vec       []float32
		timestamp int64
		want      want
	}

	now := time.Now().UnixNano()

	tests := []test{
		{
			name:      "success insert and delete",
			uuid:      "test-uuid",
			vec:       []float32{1.0, 2.0, 3.0},
			timestamp: now,
			want: want{
				vec:       []float32{1.0, 2.0, 3.0},
				timestamp: now,
				exists:    true,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			vq, err := New()
			require.NoError(t, err)

			// Insert data into the queue
			vq.PushInsert(test.uuid, test.vec, test.timestamp)

			// Test that the data exists in the queue
			gotVec, gotTimestamp, exists := vq.GetVector(test.uuid)
			require.Equal(t, test.want.vec, gotVec)
			require.Equal(t, test.want.timestamp, gotTimestamp)
			require.Equal(t, test.want.exists, exists)

			// Delete data from the queue
			vq.PushDelete(test.uuid, time.Now().UnixNano())

			// Test that the data no longer exists in the queue
			_, _, exists = vq.GetVector(test.uuid)
			require.False(t, exists)
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Queue
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Queue, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Queue, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_PushInsert(t *testing.T) {
// 	type args struct {
// 		uuid      string
// 		vector    []float32
// 		timestamp int64
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vector:nil,
// 		           timestamp:0,
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           vector:nil,
// 		           timestamp:0,
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			err := v.PushInsert(test.args.uuid, test.args.vector, test.args.timestamp)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_PushDelete(t *testing.T) {
// 	type args struct {
// 		uuid      string
// 		timestamp int64
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           timestamp:0,
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           timestamp:0,
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			err := v.PushDelete(test.args.uuid, test.args.timestamp)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_PopDelete(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantTimestamp int64
// 		wantOk        bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int64, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotTimestamp int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotTimestamp, gotOk := v.PopDelete(test.args.uuid)
// 			if err := checkFunc(test.want, gotTimestamp, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_GetVector(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantVec       []float32
// 		wantTimestamp int64
// 		wantExists    bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, int64, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotTimestamp int64, gotExists bool) error {
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		if !reflect.DeepEqual(gotExists, w.wantExists) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotExists, w.wantExists)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotVec, gotTimestamp, gotExists := v.GetVector(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotTimestamp, gotExists); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_GetVectorWithVQTimestamp(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantVec   []float32
// 		wantIts   int64
// 		wantDts   int64
// 		wantValid bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, int64, int64, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotIts int64, gotDts int64, gotValid bool) error {
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotIts, w.wantIts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIts, w.wantIts)
// 		}
// 		if !reflect.DeepEqual(gotDts, w.wantDts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDts, w.wantDts)
// 		}
// 		if !reflect.DeepEqual(gotValid, w.wantValid) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotValid, w.wantValid)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotVec, gotIts, gotDts, gotValid := v.GetVectorWithVQTimestamp(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotIts, gotDts, gotValid); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_IVExists(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantTimestamp int64
// 		wantOk        bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int64, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotTimestamp int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotTimestamp, gotOk := v.IVExists(test.args.uuid)
// 			if err := checkFunc(test.want, gotTimestamp, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_DVExists(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantTimestamp int64
// 		wantOk        bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int64, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotTimestamp int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotTimestamp, gotOk := v.DVExists(test.args.uuid)
// 			if err := checkFunc(test.want, gotTimestamp, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_RangePopInsert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		now int64
// 		f   func(uuid string, vector []float32, timestamp int64) bool
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           now:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           now:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			v.RangePopInsert(test.args.ctx, test.args.now, test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vqueue_RangePopDelete(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		now int64
// 		f   func(uuid string) bool
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           now:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           now:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			v.RangePopDelete(test.args.ctx, test.args.now, test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vqueue_Range(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(uuid string, vector []float32, ts int64) bool
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			v.Range(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vqueue_IVQLen(t *testing.T) {
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantL int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotL int) error {
// 		if !reflect.DeepEqual(gotL, w.wantL) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotL, w.wantL)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
// 		       },
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
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
// 		           },
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotL := v.IVQLen()
// 			if err := checkFunc(test.want, gotL); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_vqueue_DVQLen(t *testing.T) {
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantL int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotL int) error {
// 		if !reflect.DeepEqual(gotL, w.wantL) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotL, w.wantL)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
// 		       },
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
// 		           fields: fields {
// 		           il:nil,
// 		           dl:nil,
// 		           ic:0,
// 		           dc:0,
// 		           },
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
// 			v := &vqueue{
// 				il: test.fields.il,
// 				dl: test.fields.dl,
// 				ic: test.fields.ic,
// 				dc: test.fields.dc,
// 			}
//
// 			gotL := v.DVQLen()
// 			if err := checkFunc(test.want, gotL); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
