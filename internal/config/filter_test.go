//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestEgressFilter_Bind(t *testing.T) {
	type fields struct {
		Client          *GRPCClient
		DistanceFilters []*DistanceFilterConfig
		ObjectFilters   []*ObjectFilterConfig
	}
	type want struct {
		want *EgressFilter
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *EgressFilter) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *EgressFilter) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return EgressFilter when the bind successes",
				fields: fields{
					DistanceFilters: []*DistanceFilterConfig{
						{
							Addr:  "192.168.1.2",
							Query: "distQuery",
						},
					},
					ObjectFilters: []*ObjectFilterConfig{
						{
							Addr:  "192.168.1.3",
							Query: "objQuery",
						},
					},
				},
				want: want{
					want: &EgressFilter{
						DistanceFilters: []*DistanceFilterConfig{
							{
								Addr:  "192.168.1.2",
								Query: "distQuery",
							},
						},
						ObjectFilters: []*ObjectFilterConfig{
							{
								Addr:  "192.168.1.3",
								Query: "objQuery",
							},
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return EgressFilter when the bind successes and the Client is not nil",
				fields: fields{
					DistanceFilters: []*DistanceFilterConfig{
						{
							Addr:  "192.168.1.2",
							Query: "distQuery",
						},
					},
					ObjectFilters: []*ObjectFilterConfig{
						{
							Addr:  "192.168.1.3",
							Query: "objQuery",
						},
					},
					Client: new(GRPCClient),
				},
				want: want{
					want: &EgressFilter{
						DistanceFilters: []*DistanceFilterConfig{
							{
								Addr:  "192.168.1.2",
								Query: "distQuery",
							},
						},
						ObjectFilters: []*ObjectFilterConfig{
							{
								Addr:  "192.168.1.3",
								Query: "objQuery",
							},
						},
						Client: &GRPCClient{
							ConnectionPool: new(ConnectionPool),
							DialOption: &DialOption{
								Insecure: true,
							},
							TLS: new(TLS),
						},
					},
				},
			}
		}(),
		func() test {
			suffix := "_FOR_TEST_EGRESS_FILTER_BIND"
			m := map[string]string{
				"DISTANCE_FILTERS" + suffix:       "192.168.1.2",
				"OBJECT_FILTERS" + suffix:         "192.168.1.3",
				"DISTANCE_FILTERS_QUERY" + suffix: "distQuery",
				"OBJECT_FILTERS_QUERY" + suffix:   "objQuery",
			}
			return test{
				name: "return EgressFilter when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					DistanceFilters: []*DistanceFilterConfig{
						{
							Addr:  "_DISTANCE_FILTERS" + suffix + "_",
							Query: "_DISTANCE_FILTERS_QUERY" + suffix + "_",
						},
					},
					ObjectFilters: []*ObjectFilterConfig{
						{
							Addr:  "_OBJECT_FILTERS" + suffix + "_",
							Query: "_OBJECT_FILTERS_QUERY" + suffix + "_",
						},
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &EgressFilter{
						DistanceFilters: []*DistanceFilterConfig{
							{
								Addr:  "192.168.1.2",
								Query: "distQuery",
							},
						},
						ObjectFilters: []*ObjectFilterConfig{
							{
								Addr:  "192.168.1.3",
								Query: "objQuery",
							},
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			e := &EgressFilter{
				Client:          test.fields.Client,
				DistanceFilters: test.fields.DistanceFilters,
				ObjectFilters:   test.fields.ObjectFilters,
			}

			got := e.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIngressFilter_Bind(t *testing.T) {
	type fields struct {
		Client        *GRPCClient
		Vectorizer    string
		SearchFilters []string
		InsertFilters []string
		UpdateFilters []string
		UpsertFilters []string
	}
	type want struct {
		want *IngressFilter
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *IngressFilter) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *IngressFilter) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return IngressFilter when the bind successes",
				fields: fields{
					Vectorizer: "192.168.1.2",
					SearchFilters: []string{
						"192.168.1.3",
					},
					InsertFilters: []string{
						"192.168.1.4",
					},
					UpdateFilters: []string{
						"192.168.1.5",
					},
					UpsertFilters: []string{
						"192.168.1.6",
					},
				},
				want: want{
					want: &IngressFilter{
						Vectorizer: "192.168.1.2",
						SearchFilters: []string{
							"192.168.1.3",
						},
						InsertFilters: []string{
							"192.168.1.4",
						},
						UpdateFilters: []string{
							"192.168.1.5",
						},
						UpsertFilters: []string{
							"192.168.1.6",
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return IngressFilter when the bind successes when the bind successes and the Client is not nil",
				fields: fields{
					Vectorizer: "192.168.1.2",
					SearchFilters: []string{
						"192.168.1.3",
					},
					InsertFilters: []string{
						"192.168.1.4",
					},
					UpdateFilters: []string{
						"192.168.1.5",
					},
					UpsertFilters: []string{
						"192.168.1.6",
					},
					Client: new(GRPCClient),
				},
				want: want{
					want: &IngressFilter{
						Vectorizer: "192.168.1.2",
						SearchFilters: []string{
							"192.168.1.3",
						},
						InsertFilters: []string{
							"192.168.1.4",
						},
						UpdateFilters: []string{
							"192.168.1.5",
						},
						UpsertFilters: []string{
							"192.168.1.6",
						},
						Client: &GRPCClient{
							ConnectionPool: new(ConnectionPool),
							DialOption: &DialOption{
								Insecure: true,
							},
							TLS: new(TLS),
						},
					},
				},
			}
		}(),
		func() test {
			sufix := "_FOR_TEST_INGRESS_FILTER_BIND"
			m := map[string]string{
				"VECTORIZER" + sufix:     "192.168.1.2",
				"SEARCH_FILTERS" + sufix: "192.168.1.3",
				"INSERT_FILTERS" + sufix: "192.168.1.4",
				"UPDATE_FILTERS" + sufix: "192.168.1.5",
				"UPSERT_FILTERS" + sufix: "192.168.1.6",
			}

			return test{
				name: "return IngressFilter when the bind successes",
				fields: fields{
					Vectorizer: "_VECTORIZER" + sufix + "_",
					SearchFilters: []string{
						"_SEARCH_FILTERS" + sufix + "_",
					},
					InsertFilters: []string{
						"_INSERT_FILTERS" + sufix + "_",
					},
					UpdateFilters: []string{
						"_UPDATE_FILTERS" + sufix + "_",
					},
					UpsertFilters: []string{
						"_UPSERT_FILTERS" + sufix + "_",
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &IngressFilter{
						Vectorizer: "192.168.1.2",
						SearchFilters: []string{
							"192.168.1.3",
						},
						InsertFilters: []string{
							"192.168.1.4",
						},
						UpdateFilters: []string{
							"192.168.1.5",
						},
						UpsertFilters: []string{
							"192.168.1.6",
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			i := &IngressFilter{
				Client:        test.fields.Client,
				Vectorizer:    test.fields.Vectorizer,
				SearchFilters: test.fields.SearchFilters,
				InsertFilters: test.fields.InsertFilters,
				UpdateFilters: test.fields.UpdateFilters,
				UpsertFilters: test.fields.UpsertFilters,
			}

			got := i.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestDistanceFilterConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		Addr  string
// 		Query string
// 	}
// 	type want struct {
// 		want *DistanceFilterConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *DistanceFilterConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *DistanceFilterConfig) error {
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
// 		           Addr:"",
// 		           Query:"",
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
// 		           Addr:"",
// 		           Query:"",
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
// 			c := &DistanceFilterConfig{
// 				Addr:  test.fields.Addr,
// 				Query: test.fields.Query,
// 			}
//
// 			got := c.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func TestObjectFilterConfig_Bind(t *testing.T) {
// 	type fields struct {
// 		Addr  string
// 		Query string
// 	}
// 	type want struct {
// 		want *ObjectFilterConfig
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ObjectFilterConfig) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *ObjectFilterConfig) error {
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
// 		           Addr:"",
// 		           Query:"",
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
// 		           Addr:"",
// 		           Query:"",
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
// 			c := &ObjectFilterConfig{
// 				Addr:  test.fields.Addr,
// 				Query: test.fields.Query,
// 			}
//
// 			got := c.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
