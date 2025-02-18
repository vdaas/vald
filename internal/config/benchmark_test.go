// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
// func TestBenchmarkTarget_Bind(t *testing.T) {
// 	type fields struct {
// 		Host string
// 		Port int
// 		Meta map[string]string
// 	}
// 	type want struct {
// 		want *BenchmarkTarget
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkTarget) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkTarget) error {
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
// 		           Host:"",
// 		           Port:0,
// 		           Meta:nil,
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
// 		           Host:"",
// 		           Port:0,
// 		           Meta:nil,
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
// 			tr := &BenchmarkTarget{
// 				Host: test.fields.Host,
// 				Port: test.fields.Port,
// 				Meta: test.fields.Meta,
// 			}
//
// 			got := tr.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkDataset_Bind(t *testing.T) {
// 	type fields struct {
// 		Name    string
// 		Group   string
// 		Indexes int
// 		Range   *BenchmarkDatasetRange
// 		URL     string
// 	}
// 	type want struct {
// 		want *BenchmarkDataset
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkDataset) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkDataset) error {
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
// 		           Name:"",
// 		           Group:"",
// 		           Indexes:0,
// 		           Range:BenchmarkDatasetRange{},
// 		           URL:"",
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
// 		           Name:"",
// 		           Group:"",
// 		           Indexes:0,
// 		           Range:BenchmarkDatasetRange{},
// 		           URL:"",
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
// 			d := &BenchmarkDataset{
// 				Name:    test.fields.Name,
// 				Group:   test.fields.Group,
// 				Indexes: test.fields.Indexes,
// 				Range:   test.fields.Range,
// 				URL:     test.fields.URL,
// 			}
//
// 			got := d.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkJobRule_Bind(t *testing.T) {
// 	type fields struct {
// 		Name string
// 		Type string
// 	}
// 	type want struct {
// 		want *BenchmarkJobRule
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkJobRule) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkJobRule) error {
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
// 		           Name:"",
// 		           Type:"",
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
// 		           Name:"",
// 		           Type:"",
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
// 			r := &BenchmarkJobRule{
// 				Name: test.fields.Name,
// 				Type: test.fields.Type,
// 			}
//
// 			got := r.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestInsertConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		SkipStrictExistCheck bool
// 		Timestamp            string
// 	}
// 	type want struct {
// 		want *InsertConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *InsertConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *InsertConfig) error {
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
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
// 			cfg := &InsertConfig{
// 				SkipStrictExistCheck: test.fields.SkipStrictExistCheck,
// 				Timestamp:            test.fields.Timestamp,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestUpdateConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		SkipStrictExistCheck  bool
// 		Timestamp             string
// 		DisableBalancedUpdate bool
// 	}
// 	type want struct {
// 		want *UpdateConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *UpdateConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *UpdateConfig) error {
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
// 		           DisableBalancedUpdate:false,
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
// 		           DisableBalancedUpdate:false,
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
// 			cfg := &UpdateConfig{
// 				SkipStrictExistCheck:  test.fields.SkipStrictExistCheck,
// 				Timestamp:             test.fields.Timestamp,
// 				DisableBalancedUpdate: test.fields.DisableBalancedUpdate,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestUpsertConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		SkipStrictExistCheck  bool
// 		Timestamp             string
// 		DisableBalancedUpdate bool
// 	}
// 	type want struct {
// 		want *UpsertConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *UpsertConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *UpsertConfig) error {
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
// 		           DisableBalancedUpdate:false,
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
// 		           DisableBalancedUpdate:false,
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
// 			cfg := &UpsertConfig{
// 				SkipStrictExistCheck:  test.fields.SkipStrictExistCheck,
// 				Timestamp:             test.fields.Timestamp,
// 				DisableBalancedUpdate: test.fields.DisableBalancedUpdate,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestSearchConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		Epsilon              float32
// 		Radius               float32
// 		Num                  int32
// 		MinNum               int32
// 		Timeout              string
// 		EnableLinearSearch   bool
// 		AggregationAlgorithm string
// 	}
// 	type want struct {
// 		want *SearchConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *SearchConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *SearchConfig) error {
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
// 		           Epsilon:0,
// 		           Radius:0,
// 		           Num:0,
// 		           MinNum:0,
// 		           Timeout:"",
// 		           EnableLinearSearch:false,
// 		           AggregationAlgorithm:"",
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
// 		           Epsilon:0,
// 		           Radius:0,
// 		           Num:0,
// 		           MinNum:0,
// 		           Timeout:"",
// 		           EnableLinearSearch:false,
// 		           AggregationAlgorithm:"",
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
// 			cfg := &SearchConfig{
// 				Epsilon:              test.fields.Epsilon,
// 				Radius:               test.fields.Radius,
// 				Num:                  test.fields.Num,
// 				MinNum:               test.fields.MinNum,
// 				Timeout:              test.fields.Timeout,
// 				EnableLinearSearch:   test.fields.EnableLinearSearch,
// 				AggregationAlgorithm: test.fields.AggregationAlgorithm,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestRemoveConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		SkipStrictExistCheck bool
// 		Timestamp            string
// 	}
// 	type want struct {
// 		want *RemoveConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *RemoveConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *RemoveConfig) error {
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
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
// 		           SkipStrictExistCheck:false,
// 		           Timestamp:"",
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
// 			cfg := &RemoveConfig{
// 				SkipStrictExistCheck: test.fields.SkipStrictExistCheck,
// 				Timestamp:            test.fields.Timestamp,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestObjectConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		FilterConfig FilterConfig
// 	}
// 	type want struct {
// 		want *ObjectConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ObjectConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *ObjectConfig) error {
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
// 		           FilterConfig:FilterConfig{},
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
// 		           FilterConfig:FilterConfig{},
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
// 			cfg := &ObjectConfig{
// 				FilterConfig: test.fields.FilterConfig,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestFilterTarget_Bind(t *testing.T) {
// 	type fields struct {
// 		Host string
// 		Port int32
// 	}
// 	type want struct {
// 		want *FilterTarget
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *FilterTarget) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *FilterTarget) error {
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
// 		           Host:"",
// 		           Port:0,
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
// 		           Host:"",
// 		           Port:0,
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
// 			cfg := &FilterTarget{
// 				Host: test.fields.Host,
// 				Port: test.fields.Port,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestFilterConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		Targets []*FilterTarget
// 	}
// 	type want struct {
// 		want *FilterConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *FilterConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *FilterConfig) error {
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
// 		           Targets:nil,
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
// 		           Targets:nil,
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
// 			cfg := &FilterConfig{
// 				Targets: test.fields.Targets,
// 			}
//
// 			got := cfg.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkJob_Bind(t *testing.T) {
// 	type fields struct {
// 		Target             *BenchmarkTarget
// 		Dataset            *BenchmarkDataset
// 		Replica            int
// 		Repetition         int
// 		JobType            string
// 		InsertConfig       *InsertConfig
// 		UpdateConfig       *UpdateConfig
// 		UpsertConfig       *UpsertConfig
// 		SearchConfig       *SearchConfig
// 		RemoveConfig       *RemoveConfig
// 		ObjectConfig       *ObjectConfig
// 		ClientConfig       *GRPCClient
// 		Rules              []*BenchmarkJobRule
// 		BeforeJobName      string
// 		BeforeJobNamespace string
// 		RPS                int
// 		ConcurrencyLimit   int
// 	}
// 	type want struct {
// 		want *BenchmarkJob
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkJob) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkJob) error {
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
// 		           Target:BenchmarkTarget{},
// 		           Dataset:BenchmarkDataset{},
// 		           Replica:0,
// 		           Repetition:0,
// 		           JobType:"",
// 		           InsertConfig:InsertConfig{},
// 		           UpdateConfig:UpdateConfig{},
// 		           UpsertConfig:UpsertConfig{},
// 		           SearchConfig:SearchConfig{},
// 		           RemoveConfig:RemoveConfig{},
// 		           ObjectConfig:ObjectConfig{},
// 		           ClientConfig:GRPCClient{},
// 		           Rules:nil,
// 		           BeforeJobName:"",
// 		           BeforeJobNamespace:"",
// 		           RPS:0,
// 		           ConcurrencyLimit:0,
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
// 		           Target:BenchmarkTarget{},
// 		           Dataset:BenchmarkDataset{},
// 		           Replica:0,
// 		           Repetition:0,
// 		           JobType:"",
// 		           InsertConfig:InsertConfig{},
// 		           UpdateConfig:UpdateConfig{},
// 		           UpsertConfig:UpsertConfig{},
// 		           SearchConfig:SearchConfig{},
// 		           RemoveConfig:RemoveConfig{},
// 		           ObjectConfig:ObjectConfig{},
// 		           ClientConfig:GRPCClient{},
// 		           Rules:nil,
// 		           BeforeJobName:"",
// 		           BeforeJobNamespace:"",
// 		           RPS:0,
// 		           ConcurrencyLimit:0,
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
// 			b := &BenchmarkJob{
// 				Target:             test.fields.Target,
// 				Dataset:            test.fields.Dataset,
// 				Replica:            test.fields.Replica,
// 				Repetition:         test.fields.Repetition,
// 				JobType:            test.fields.JobType,
// 				InsertConfig:       test.fields.InsertConfig,
// 				UpdateConfig:       test.fields.UpdateConfig,
// 				UpsertConfig:       test.fields.UpsertConfig,
// 				SearchConfig:       test.fields.SearchConfig,
// 				RemoveConfig:       test.fields.RemoveConfig,
// 				ObjectConfig:       test.fields.ObjectConfig,
// 				ClientConfig:       test.fields.ClientConfig,
// 				Rules:              test.fields.Rules,
// 				BeforeJobName:      test.fields.BeforeJobName,
// 				BeforeJobNamespace: test.fields.BeforeJobNamespace,
// 				RPS:                test.fields.RPS,
// 				ConcurrencyLimit:   test.fields.ConcurrencyLimit,
// 			}
//
// 			got := b.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkScenario_Bind(t *testing.T) {
// 	type fields struct {
// 		Target  *BenchmarkTarget
// 		Dataset *BenchmarkDataset
// 		Jobs    []*BenchmarkJob
// 	}
// 	type want struct {
// 		want *BenchmarkScenario
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkScenario) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkScenario) error {
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
// 		           Target:BenchmarkTarget{},
// 		           Dataset:BenchmarkDataset{},
// 		           Jobs:nil,
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
// 		           Target:BenchmarkTarget{},
// 		           Dataset:BenchmarkDataset{},
// 		           Jobs:nil,
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
// 			b := &BenchmarkScenario{
// 				Target:  test.fields.Target,
// 				Dataset: test.fields.Dataset,
// 				Jobs:    test.fields.Jobs,
// 			}
//
// 			got := b.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkJobImageInfo_Bind(t *testing.T) {
// 	type fields struct {
// 		Repository string
// 		Tag        string
// 		PullPolicy string
// 	}
// 	type want struct {
// 		want *BenchmarkJobImageInfo
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkJobImageInfo) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkJobImageInfo) error {
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
// 		           Repository:"",
// 		           Tag:"",
// 		           PullPolicy:"",
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
// 		           Repository:"",
// 		           Tag:"",
// 		           PullPolicy:"",
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
// 			b := &BenchmarkJobImageInfo{
// 				Repository: test.fields.Repository,
// 				Tag:        test.fields.Tag,
// 				PullPolicy: test.fields.PullPolicy,
// 			}
//
// 			got := b.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
