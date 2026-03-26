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
package config

// NOT IMPLEMENTED BELOW
//
// func TestFaiss_Bind(t *testing.T) {
// 	type fields struct {
// 		VQueue                  *VQueue
// 		KVSDB                   *KVSDB
// 		AutoIndexDurationLimit  string
// 		InitialDelayMaxDuration string
// 		LoadIndexTimeoutFactor  string
// 		MethodType              string
// 		MetricType              string
// 		MaxLoadIndexTimeout     string
// 		AutoIndexCheckDuration  string
// 		AutoSaveIndexDuration   string
// 		IndexPath               string
// 		MinLoadIndexTimeout     string
// 		AutoIndexLength         int
// 		M                       int
// 		NbitsPerIdx             int
// 		Nlist                   int
// 		Dimension               int
// 		EnableInMemoryMode      bool
// 		EnableProactiveGC       bool
// 		EnableCopyOnWrite       bool
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
// 		           VQueue:VQueue{},
// 		           KVSDB:KVSDB{},
// 		           AutoIndexDurationLimit:"",
// 		           InitialDelayMaxDuration:"",
// 		           LoadIndexTimeoutFactor:"",
// 		           MethodType:"",
// 		           MetricType:"",
// 		           MaxLoadIndexTimeout:"",
// 		           AutoIndexCheckDuration:"",
// 		           AutoSaveIndexDuration:"",
// 		           IndexPath:"",
// 		           MinLoadIndexTimeout:"",
// 		           AutoIndexLength:0,
// 		           M:0,
// 		           NbitsPerIdx:0,
// 		           Nlist:0,
// 		           Dimension:0,
// 		           EnableInMemoryMode:false,
// 		           EnableProactiveGC:false,
// 		           EnableCopyOnWrite:false,
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
// 		           VQueue:VQueue{},
// 		           KVSDB:KVSDB{},
// 		           AutoIndexDurationLimit:"",
// 		           InitialDelayMaxDuration:"",
// 		           LoadIndexTimeoutFactor:"",
// 		           MethodType:"",
// 		           MetricType:"",
// 		           MaxLoadIndexTimeout:"",
// 		           AutoIndexCheckDuration:"",
// 		           AutoSaveIndexDuration:"",
// 		           IndexPath:"",
// 		           MinLoadIndexTimeout:"",
// 		           AutoIndexLength:0,
// 		           M:0,
// 		           NbitsPerIdx:0,
// 		           Nlist:0,
// 		           Dimension:0,
// 		           EnableInMemoryMode:false,
// 		           EnableProactiveGC:false,
// 		           EnableCopyOnWrite:false,
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
// 				VQueue:                  test.fields.VQueue,
// 				KVSDB:                   test.fields.KVSDB,
// 				AutoIndexDurationLimit:  test.fields.AutoIndexDurationLimit,
// 				InitialDelayMaxDuration: test.fields.InitialDelayMaxDuration,
// 				LoadIndexTimeoutFactor:  test.fields.LoadIndexTimeoutFactor,
// 				MethodType:              test.fields.MethodType,
// 				MetricType:              test.fields.MetricType,
// 				MaxLoadIndexTimeout:     test.fields.MaxLoadIndexTimeout,
// 				AutoIndexCheckDuration:  test.fields.AutoIndexCheckDuration,
// 				AutoSaveIndexDuration:   test.fields.AutoSaveIndexDuration,
// 				IndexPath:               test.fields.IndexPath,
// 				MinLoadIndexTimeout:     test.fields.MinLoadIndexTimeout,
// 				AutoIndexLength:         test.fields.AutoIndexLength,
// 				M:                       test.fields.M,
// 				NbitsPerIdx:             test.fields.NbitsPerIdx,
// 				Nlist:                   test.fields.Nlist,
// 				Dimension:               test.fields.Dimension,
// 				EnableInMemoryMode:      test.fields.EnableInMemoryMode,
// 				EnableProactiveGC:       test.fields.EnableProactiveGC,
// 				EnableCopyOnWrite:       test.fields.EnableCopyOnWrite,
// 			}
//
// 			got := f.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
