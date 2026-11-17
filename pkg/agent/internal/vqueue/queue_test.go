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

package vqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/sync/atomic"
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
		want      want
		timestamp int64
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
// 		})
// 	}
// }
//
// func Test_vqueue_PopInsert(t *testing.T) {
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
// 		wantVector    []float32
// 		wantTimestamp int64
// 		wantOk        bool
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
// 	defaultCheckFunc := func(w want, gotVector []float32, gotTimestamp int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotVector, w.wantVector) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVector, w.wantVector)
// 		}
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
// 			gotVector, gotTimestamp, gotOk := v.PopInsert(test.args.uuid)
// 			if err := checkFunc(test.want, gotVector, gotTimestamp, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
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
// 		})
// 	}
// }
//
// func Test_vqueue_GetVectorWithTimestamp(t *testing.T) {
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
// 		wantVec    []float32
// 		wantIts    int64
// 		wantDts    int64
// 		wantExists bool
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
// 	defaultCheckFunc := func(w want, gotVec []float32, gotIts int64, gotDts int64, gotExists bool) error {
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotIts, w.wantIts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIts, w.wantIts)
// 		}
// 		if !reflect.DeepEqual(gotDts, w.wantDts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDts, w.wantDts)
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
// 			gotVec, gotIts, gotDts, gotExists := v.GetVectorWithTimestamp(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotIts, gotDts, gotExists); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vqueue_getVector(t *testing.T) {
// 	type args struct {
// 		uuid                  string
// 		enableDeleteTimestamp bool
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct {
// 		wantVec []float32
// 		wantIts int64
// 		wantDts int64
// 		wantOk  bool
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
// 	defaultCheckFunc := func(w want, gotVec []float32, gotIts int64, gotDts int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotIts, w.wantIts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIts, w.wantIts)
// 		}
// 		if !reflect.DeepEqual(gotDts, w.wantDts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDts, w.wantDts)
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
// 		           enableDeleteTimestamp:false,
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
// 		           enableDeleteTimestamp:false,
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
// 			gotVec, gotIts, gotDts, gotOk := v.getVector(test.args.uuid, test.args.enableDeleteTimestamp)
// 			if err := checkFunc(test.want, gotVec, gotIts, gotDts, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
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
// 		wantIts int64
// 		wantOk  bool
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
// 	defaultCheckFunc := func(w want, gotIts int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotIts, w.wantIts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIts, w.wantIts)
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
// 			gotIts, gotOk := v.IVExists(test.args.uuid)
// 			if err := checkFunc(test.want, gotIts, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
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
// 		wantDts int64
// 		wantOk  bool
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
// 	defaultCheckFunc := func(w want, gotDts int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotDts, w.wantDts) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDts, w.wantDts)
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
// 			gotDts, gotOk := v.DVExists(test.args.uuid)
// 			if err := checkFunc(test.want, gotDts, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
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
// 	type want struct{}
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
// 	type want struct{}
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
// 		in0 context.Context
// 		f   func(uuid string, vector []float32, ts int64) bool
// 	}
// 	type fields struct {
// 		il sync.Map[string, *index]
// 		dl sync.Map[string, *index]
// 		ic uint64
// 		dc uint64
// 	}
// 	type want struct{}
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
// 		           in0:nil,
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
// 		           in0:nil,
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
// 			v.Range(test.args.in0, test.args.f)
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
// 		})
// 	}
// }
//
// func Test_vqueue_loadIVQ(t *testing.T) {
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
// 		wantVec []float32
// 		wantTs  int64
// 		wantOk  bool
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
// 	defaultCheckFunc := func(w want, gotVec []float32, gotTs int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotTs, w.wantTs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTs, w.wantTs)
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
// 			gotVec, gotTs, gotOk := v.loadIVQ(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotTs, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vqueue_loadDVQ(t *testing.T) {
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
// 		wantTs int64
// 		wantOk bool
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
// 	defaultCheckFunc := func(w want, gotTs int64, gotOk bool) error {
// 		if !reflect.DeepEqual(gotTs, w.wantTs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTs, w.wantTs)
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
// 			gotTs, gotOk := v.loadDVQ(test.args.uuid)
// 			if err := checkFunc(test.want, gotTs, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newer(t *testing.T) {
// 	type args struct {
// 		ts1 int64
// 		ts2 int64
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           ts1:0,
// 		           ts2:0,
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
// 		           ts1:0,
// 		           ts2:0,
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
// 			got := newer(test.args.ts1, test.args.ts2)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// RangePopInsert/RangePopDelete used to pop by key after calling f, dropping an entry
// that was pushed for the same uuid while f ran (e.g. the insert half of an Update during
// CreateIndex). Only the entry that was handed to f may be removed.
func TestRangePop_keepsEntryPushedDuringCallback(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	t.Run("insert", func(t *testing.T) {
		t.Parallel()
		q, err := New()
		require.NoError(t, err)
		require.NoError(t, q.PushInsert(uuid, []float32{1}, 1))
		q.RangePopInsert(t.Context(), 10, func(id string, _ []float32, _ int64) bool {
			require.NoError(t, q.PushInsert(id, []float32{2}, 5))
			return true
		})
		vec, ts, ok := q.GetVector(uuid)
		require.True(t, ok)
		require.Equal(t, []float32{2}, vec)
		require.Equal(t, int64(5), ts)
		require.Equal(t, 1, q.IVQLen())
	})
	t.Run("delete keeps newer insert", func(t *testing.T) {
		t.Parallel()
		q, err := New()
		require.NoError(t, err)
		require.NoError(t, q.PushDelete(uuid, 1))
		q.RangePopDelete(t.Context(), 10, func(id string, _ int64) bool {
			require.NoError(t, q.PushDelete(id, 5))
			require.NoError(t, q.PushInsert(id, []float32{3}, 6))
			return true
		})
		// f removed uuid, which satisfies the delete at ts 5, so no stale delete may
		// survive to remove the insert at ts 6 on the next CreateIndex.
		require.Equal(t, 0, q.DVQLen())
		vec, its, ok := q.GetVector(uuid)
		require.True(t, ok)
		require.Equal(t, []float32{3}, vec)
		require.Equal(t, int64(6), its)
		require.Equal(t, 1, q.IVQLen())
	})
	t.Run("delete cancels insert older than the newest delete", func(t *testing.T) {
		t.Parallel()
		q, err := New()
		require.NoError(t, err)
		require.NoError(t, q.PushDelete(uuid, 1))
		q.RangePopDelete(t.Context(), 10, func(id string, _ int64) bool {
			require.NoError(t, q.PushInsert(id, []float32{3}, 3))
			require.NoError(t, q.PushDelete(id, 5))
			return true
		})
		require.Equal(t, 0, q.DVQLen())
		_, _, ok := q.GetVector(uuid)
		require.False(t, ok, "insert at ts 3 is older than the delete at ts 5")
		require.Equal(t, 0, q.IVQLen())
	})
}

// A delete that f skipped (the indexed object is newer) satisfies nothing else: it must pop
// only the entry handed to f, keep a delete pushed while f ran, and cancel no insert.
func TestRangePopDelete_skippedDelete(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		push    func(t *testing.T, q Queue)
		name    string
		wantDVQ int
		wantDTS int64
	}{
		{
			name: "pops the skipped entry",
			push: func(*testing.T, Queue) {},
		},
		{
			name: "keeps a newer delete pushed during f",
			push: func(t *testing.T, q Queue) {
				t.Helper()
				require.NoError(t, q.PushDelete(uuid, 5))
			},
			wantDVQ: 1,
			wantDTS: 5,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q, err := New()
			require.NoError(t, err)
			require.NoError(t, q.PushInsert(uuid, []float32{1}, 1))
			require.NoError(t, q.PushDelete(uuid, 2))
			var got []int64
			q.RangePopDelete(t.Context(), 10, func(id string, ts int64) bool {
				require.Equal(t, uuid, id)
				got = append(got, ts)
				tc.push(t, q)
				return false
			})
			require.Equal(t, []int64{2}, got, "f must receive the delete timestamp")
			require.Equal(t, tc.wantDVQ, q.DVQLen())
			dts, ok := q.DVExists(uuid)
			require.Equal(t, tc.wantDVQ == 1, ok)
			require.Equal(t, tc.wantDTS, dts)
			vec, its, ok := q.PopInsert(uuid)
			require.True(t, ok, "a skipped delete must not cancel the insert")
			require.Equal(t, []float32{1}, vec)
			require.Equal(t, int64(1), its)
		})
	}
}

// An Update (delete at t, insert at t+1) pushed after RangePopDelete took its snapshot
// used to have its insert applied by RangePopInsert while its delete stayed queued, so the
// next CreateIndex removed the updated object. The insert must wait for the delete.
func TestRangePopInsert_waitsForPendingDelete(t *testing.T) {
	t.Parallel()
	const (
		uuid  = "uuid"
		other = "other"
		now   = 10
	)
	tests := []struct {
		name string
		dts  int64
		its  int64
		// insertQueued reports whether the insert is still queued after the first run.
		insertQueued bool
	}{
		{name: "update pushed after the delete snapshot waits", dts: 2, its: 3, insertQueued: true},
		{name: "equal timestamps wait", dts: 3, its: 3, insertQueued: true},
		{name: "newer delete still cancels the insert", dts: 5, its: 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q, err := New()
			require.NoError(t, err)
			require.NoError(t, q.PushDelete(other, 1))
			// f runs after the Range snapshot, so these pushes are not in it.
			q.RangePopDelete(t.Context(), now, func(string, int64) bool {
				require.NoError(t, q.PushDelete(uuid, tc.dts))
				require.NoError(t, q.PushInsert(uuid, []float32{1}, tc.its))
				return true
			})
			dts, ok := q.PopDelete(other)
			require.False(t, ok, "other's delete must have been popped, got %d", dts)

			var applied []string
			q.RangePopInsert(t.Context(), now, func(id string, _ []float32, _ int64) bool {
				applied = append(applied, id)
				return true
			})
			require.Empty(t, applied)
			if tc.insertQueued {
				require.Equal(t, 1, q.IVQLen())
			} else {
				require.Equal(t, 0, q.IVQLen())
			}
			require.Equal(t, 1, q.DVQLen())

			// The next run applies the delete first and then the insert.
			var order []string
			q.RangePopDelete(t.Context(), now, func(id string, _ int64) bool {
				order = append(order, "delete:"+id)
				return true
			})
			q.RangePopInsert(t.Context(), now, func(id string, _ []float32, ts int64) bool {
				order = append(order, "insert:"+id)
				require.Equal(t, tc.its, ts)
				return true
			})
			want := []string{"delete:" + uuid}
			if tc.insertQueued {
				want = append(want, "insert:"+uuid)
			}
			require.Equal(t, want, order)
			require.Equal(t, 0, q.IVQLen())
			require.Equal(t, 0, q.DVQLen())
		})
	}
}

// Concurrent pushes for one uuid used to race between LoadOrStore and Store, so an
// older timestamp could overwrite a newer one.
func TestPushInsert_concurrentKeepsNewest(t *testing.T) {
	t.Parallel()
	q, err := New()
	require.NoError(t, err)
	const (
		uuid    = "uuid"
		writers = 64
	)
	start := make(chan struct{})
	done := make(chan struct{}, writers)
	for i := range writers {
		go func() {
			defer func() { done <- struct{}{} }()
			<-start
			if err := q.PushInsert(uuid, []float32{float32(i)}, int64(i+1)); err != nil {
				t.Error(err)
			}
		}()
	}
	close(start)
	for range writers {
		<-done
	}
	_, ts, ok := q.GetVector(uuid)
	require.True(t, ok)
	require.Equal(t, int64(writers), ts)
	require.Equal(t, 1, q.IVQLen())
}

// push used to publish an entry before counting it, so a concurrent pop could
// decrement the counter first and wrap it below zero (IVQLen/DVQLen < 0, and a
// makeslice panic in RangePop*). The counter must stay non-negative under any
// interleaving and match the queue at quiescence; a push onto an existing key must
// leave it unchanged.
func TestPush_counterNeverUnderflows(t *testing.T) {
	t.Parallel()
	t.Run("existing key keeps the count", func(t *testing.T) {
		t.Parallel()
		q, err := New()
		require.NoError(t, err)
		require.NoError(t, q.PushInsert("uuid", []float32{1}, 2))
		require.NoError(t, q.PushInsert("uuid", []float32{2}, 1)) // older: ignored
		require.NoError(t, q.PushInsert("uuid", []float32{3}, 3)) // newer: replaces
		require.Equal(t, 1, q.IVQLen())
		vec, ts, ok := q.GetVector("uuid")
		require.True(t, ok)
		require.Equal(t, []float32{3}, vec)
		require.Equal(t, int64(3), ts)
		require.NoError(t, q.PushDelete("uuid", 2))
		require.NoError(t, q.PushDelete("uuid", 1))
		require.NoError(t, q.PushDelete("uuid", 4))
		require.Equal(t, 1, q.DVQLen())
		_, _, ok = q.PopInsert("uuid")
		require.True(t, ok)
		_, ok = q.PopDelete("uuid")
		require.True(t, ok)
		require.Equal(t, 0, q.IVQLen())
		require.Equal(t, 0, q.DVQLen())
	})
	t.Run("concurrent push and pop", func(t *testing.T) {
		t.Parallel()
		qi, err := New()
		require.NoError(t, err)
		q, ok := qi.(*vqueue)
		require.True(t, ok)
		const (
			workers = 8
			iters   = 2000
		)
		uuids := []string{"a", "b", "c", "d"}
		var negative atomic.Bool
		check := func() {
			if q.IVQLen() < 0 || q.DVQLen() < 0 {
				negative.Store(true)
			}
		}
		start := make(chan struct{})
		done := make(chan struct{}, 2*workers)
		for w := range workers {
			go func() {
				defer func() { done <- struct{}{} }()
				<-start
				for i := range iters {
					id := uuids[(w+i)%len(uuids)]
					ts := int64(w*iters + i + 1)
					if err := q.PushInsert(id, []float32{1}, ts); err != nil {
						t.Error(err)
						return
					}
					if err := q.PushDelete(id, ts); err != nil {
						t.Error(err)
						return
					}
					check()
				}
			}()
			go func() {
				defer func() { done <- struct{}{} }()
				<-start
				for i := range iters {
					id := uuids[(w+i)%len(uuids)]
					q.PopInsert(id)
					check()
					q.PopDelete(id)
					check()
				}
			}()
		}
		close(start)
		timeout := time.NewTimer(time.Minute)
		defer timeout.Stop()
		for range 2 * workers {
			select {
			case <-done:
			case <-timeout.C:
				t.Fatal("workers did not finish")
			}
		}
		require.False(t, negative.Load(), "queue length went below zero")
		var ic, dc int
		for _, id := range uuids {
			if _, _, ok := q.loadIVQ(id); ok {
				ic++
			}
			if _, ok := q.loadDVQ(id); ok {
				dc++
			}
		}
		require.Equal(t, ic, q.IVQLen())
		require.Equal(t, dc, q.DVQLen())
	})
}
