//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestNGT_Bind(t *testing.T) {
	type fields struct {
		VQueue                  *VQueue
		KVSDB                   *KVSDB
		AutoSaveIndexDuration   string
		DistanceType            string
		ObjectType              string
		AutoIndexDurationLimit  string
		AutoIndexCheckDuration  string
		IndexPath               string
		InitialDelayMaxDuration string
		CreationEdgeSize        int
		SearchEdgeSize          int
		AutoIndexLength         int
		BulkInsertChunkSize     int
		Dimension               int
		EnableInMemoryMode      bool
	}
	type want struct {
		want *NGT
	}
	type test struct {
		want       want
		checkFunc  func(want, *NGT) error
		beforeFunc func(*testing.T)
		afterFunc  func()
		name       string
		fields     fields
	}
	defaultCheckFunc := func(w want, got *NGT) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return NGT when all fields contain no prefix/suffix symbol",
			fields: fields{
				IndexPath:               "config/ngt",
				Dimension:               1000,
				BulkInsertChunkSize:     100,
				DistanceType:            "l2",
				ObjectType:              "float",
				CreationEdgeSize:        3,
				SearchEdgeSize:          5,
				AutoIndexDurationLimit:  "1h",
				AutoIndexCheckDuration:  "30m",
				AutoSaveIndexDuration:   "30m",
				AutoIndexLength:         100,
				InitialDelayMaxDuration: "1h",
				EnableInMemoryMode:      false,
				VQueue:                  new(VQueue),
				KVSDB:                   new(KVSDB),
			},
			want: want{
				want: &NGT{
					IndexPath:               "config/ngt",
					Dimension:               1000,
					BulkInsertChunkSize:     100,
					DistanceType:            "l2",
					ObjectType:              "float",
					CreationEdgeSize:        3,
					SearchEdgeSize:          5,
					AutoIndexDurationLimit:  "1h",
					AutoIndexCheckDuration:  "30m",
					AutoSaveIndexDuration:   "30m",
					AutoIndexLength:         100,
					InitialDelayMaxDuration: "1h",
					EnableInMemoryMode:      false,
					VQueue:                  new(VQueue),
					KVSDB:                   new(KVSDB),
				},
			},
		},
		{
			name: "return NGT with environment variable when it contains `_` as prefix and suffix",
			fields: fields{
				IndexPath:               "_NGT_BIND_INDEX_PATH_",
				Dimension:               1000,
				BulkInsertChunkSize:     100,
				DistanceType:            "_NGT_BIND_DISTANCE_TYPE_",
				ObjectType:              "_NGT_BIND_OBJECT_TYPE_",
				CreationEdgeSize:        3,
				SearchEdgeSize:          5,
				AutoIndexDurationLimit:  "_NGT_BIND_AUTO_INDEX_DURATION_LIMIT_",
				AutoIndexCheckDuration:  "_NGT_BIND_AUTO_INDEX_CHECK_DURATION_",
				AutoSaveIndexDuration:   "_NGT_BIND_AUTO_SAVE_INDEX_DURATION_",
				AutoIndexLength:         100,
				InitialDelayMaxDuration: "_NGT_BIND_INITIAL_DELAY_MAX_DURATION_",
				EnableInMemoryMode:      false,
				VQueue:                  new(VQueue),
				KVSDB:                   new(KVSDB),
			},
			beforeFunc: func(t *testing.T) {
				t.Helper()
				t.Setenv("NGT_BIND_INDEX_PATH", "config/ngt")
				t.Setenv("NGT_BIND_DISTANCE_TYPE", "l2")
				t.Setenv("NGT_BIND_OBJECT_TYPE", "float")
				t.Setenv("NGT_BIND_AUTO_INDEX_DURATION_LIMIT", "1h")
				t.Setenv("NGT_BIND_AUTO_INDEX_CHECK_DURATION", "30m")
				t.Setenv("NGT_BIND_AUTO_SAVE_INDEX_DURATION", "30m")
				t.Setenv("NGT_BIND_INITIAL_DELAY_MAX_DURATION", "1h")
			},
			want: want{
				want: &NGT{
					IndexPath:               "config/ngt",
					Dimension:               1000,
					BulkInsertChunkSize:     100,
					DistanceType:            "l2",
					ObjectType:              "float",
					CreationEdgeSize:        3,
					SearchEdgeSize:          5,
					AutoIndexDurationLimit:  "1h",
					AutoIndexCheckDuration:  "30m",
					AutoSaveIndexDuration:   "30m",
					AutoIndexLength:         100,
					InitialDelayMaxDuration: "1h",
					EnableInMemoryMode:      false,
					VQueue:                  new(VQueue),
					KVSDB:                   new(KVSDB),
				},
			},
		},
		{
			name: "returns NGT when all fields are empty",
			want: want{
				want: &NGT{
					VQueue: new(VQueue),
					KVSDB:  new(KVSDB),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			n := &NGT{
				IndexPath:               test.fields.IndexPath,
				Dimension:               test.fields.Dimension,
				BulkInsertChunkSize:     test.fields.BulkInsertChunkSize,
				DistanceType:            test.fields.DistanceType,
				ObjectType:              test.fields.ObjectType,
				CreationEdgeSize:        test.fields.CreationEdgeSize,
				SearchEdgeSize:          test.fields.SearchEdgeSize,
				AutoIndexDurationLimit:  test.fields.AutoIndexDurationLimit,
				AutoIndexCheckDuration:  test.fields.AutoIndexCheckDuration,
				AutoSaveIndexDuration:   test.fields.AutoSaveIndexDuration,
				AutoIndexLength:         test.fields.AutoIndexLength,
				InitialDelayMaxDuration: test.fields.InitialDelayMaxDuration,
				EnableInMemoryMode:      test.fields.EnableInMemoryMode,
			}

			got := n.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestKVSDB_Bind(t *testing.T) {
// 	type fields struct {
// 		Concurrency int
// 	}
// 	type want struct {
// 		want *KVSDB
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *KVSDB) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *KVSDB) error {
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
// 		       fields: fields {
// 		           Concurrency:0,
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
// 		           Concurrency:0,
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
// 			k := &KVSDB{
// 				Concurrency: test.fields.Concurrency,
// 			}
//
// 			got := k.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestVQueue_Bind(t *testing.T) {
// 	type fields struct {
// 		InsertBufferPoolSize int
// 		DeleteBufferPoolSize int
// 	}
// 	type want struct {
// 		want *VQueue
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *VQueue) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *VQueue) error {
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
// 		       fields: fields {
// 		           InsertBufferPoolSize:0,
// 		           DeleteBufferPoolSize:0,
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
// 		           InsertBufferPoolSize:0,
// 		           DeleteBufferPoolSize:0,
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
// 			vq := &VQueue{
// 				InsertBufferPoolSize: test.fields.InsertBufferPoolSize,
// 				DeleteBufferPoolSize: test.fields.DeleteBufferPoolSize,
// 			}
//
// 			got := vq.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
