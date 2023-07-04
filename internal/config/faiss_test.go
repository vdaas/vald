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
package config

// NOT IMPLEMENTED BELOW
//
// func TestFaiss_Bind(t *testing.T) {
// 	type fields struct {
// 		IndexPath               string
// 		Dimension               int
// 		Nlist                   int
// 		M                       int
// 		NbitsPerIdx             int
// 		MetricType              string
// 		EnableInMemoryMode      bool
// 		AutoIndexCheckDuration  string
// 		AutoSaveIndexDuration   string
// 		AutoIndexDurationLimit  string
// 		AutoIndexLength         int
// 		InitialDelayMaxDuration string
// 		MinLoadIndexTimeout     string
// 		MaxLoadIndexTimeout     string
// 		LoadIndexTimeoutFactor  string
// 		EnableProactiveGC       bool
// 		EnableCopyOnWrite       bool
// 		VQueue                  *VQueue
// 		KVSDB                   *KVSDB
// 	}
// 	type want struct {
// 		want *Faiss
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *Faiss) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *Faiss) error {
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
// 		           IndexPath:"",
// 		           Dimension:0,
// 		           Nlist:0,
// 		           M:0,
// 		           NbitsPerIdx:0,
// 		           MetricType:"",
// 		           EnableInMemoryMode:false,
// 		           AutoIndexCheckDuration:"",
// 		           AutoSaveIndexDuration:"",
// 		           AutoIndexDurationLimit:"",
// 		           AutoIndexLength:0,
// 		           InitialDelayMaxDuration:"",
// 		           MinLoadIndexTimeout:"",
// 		           MaxLoadIndexTimeout:"",
// 		           LoadIndexTimeoutFactor:"",
// 		           EnableProactiveGC:false,
// 		           EnableCopyOnWrite:false,
// 		           VQueue:VQueue{},
// 		           KVSDB:KVSDB{},
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
// 		           IndexPath:"",
// 		           Dimension:0,
// 		           Nlist:0,
// 		           M:0,
// 		           NbitsPerIdx:0,
// 		           MetricType:"",
// 		           EnableInMemoryMode:false,
// 		           AutoIndexCheckDuration:"",
// 		           AutoSaveIndexDuration:"",
// 		           AutoIndexDurationLimit:"",
// 		           AutoIndexLength:0,
// 		           InitialDelayMaxDuration:"",
// 		           MinLoadIndexTimeout:"",
// 		           MaxLoadIndexTimeout:"",
// 		           LoadIndexTimeoutFactor:"",
// 		           EnableProactiveGC:false,
// 		           EnableCopyOnWrite:false,
// 		           VQueue:VQueue{},
// 		           KVSDB:KVSDB{},
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
// 			f := &Faiss{
// 				IndexPath:               test.fields.IndexPath,
// 				Dimension:               test.fields.Dimension,
// 				Nlist:                   test.fields.Nlist,
// 				M:                       test.fields.M,
// 				NbitsPerIdx:             test.fields.NbitsPerIdx,
// 				MetricType:              test.fields.MetricType,
// 				EnableInMemoryMode:      test.fields.EnableInMemoryMode,
// 				AutoIndexCheckDuration:  test.fields.AutoIndexCheckDuration,
// 				AutoSaveIndexDuration:   test.fields.AutoSaveIndexDuration,
// 				AutoIndexDurationLimit:  test.fields.AutoIndexDurationLimit,
// 				AutoIndexLength:         test.fields.AutoIndexLength,
// 				InitialDelayMaxDuration: test.fields.InitialDelayMaxDuration,
// 				MinLoadIndexTimeout:     test.fields.MinLoadIndexTimeout,
// 				MaxLoadIndexTimeout:     test.fields.MaxLoadIndexTimeout,
// 				LoadIndexTimeoutFactor:  test.fields.LoadIndexTimeoutFactor,
// 				EnableProactiveGC:       test.fields.EnableProactiveGC,
// 				EnableCopyOnWrite:       test.fields.EnableCopyOnWrite,
// 				VQueue:                  test.fields.VQueue,
// 				KVSDB:                   test.fields.KVSDB,
// 			}
//
// 			got := f.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
