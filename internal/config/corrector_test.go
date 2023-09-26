// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
// func TestCorrector_Bind(t *testing.T) {
// 	type fields struct {
// 		AgentPort                  int
// 		AgentName                  string
// 		AgentNamespace             string
// 		AgentDNS                   string
// 		CreationPoolSize           uint32
// 		NodeName                   string
// 		StreamListConcurrency      int
// 		BboltAsyncWriteConcurrency int
// 		Discoverer                 *DiscovererClient
// 	}
// 	type want struct {
// 		want *Corrector
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *Corrector) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *Corrector) error {
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
// 		           AgentPort:0,
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           AgentDNS:"",
// 		           CreationPoolSize:0,
// 		           NodeName:"",
// 		           StreamListConcurrency:0,
// 		           BboltAsyncWriteConcurrency:0,
// 		           Discoverer:DiscovererClient{},
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
// 		           AgentPort:0,
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           AgentDNS:"",
// 		           CreationPoolSize:0,
// 		           NodeName:"",
// 		           StreamListConcurrency:0,
// 		           BboltAsyncWriteConcurrency:0,
// 		           Discoverer:DiscovererClient{},
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
// 			c := &Corrector{
// 				AgentPort:                  test.fields.AgentPort,
// 				AgentName:                  test.fields.AgentName,
// 				AgentNamespace:             test.fields.AgentNamespace,
// 				AgentDNS:                   test.fields.AgentDNS,
// 				CreationPoolSize:           test.fields.CreationPoolSize,
// 				NodeName:                   test.fields.NodeName,
// 				StreamListConcurrency:      test.fields.StreamListConcurrency,
// 				BboltAsyncWriteConcurrency: test.fields.BboltAsyncWriteConcurrency,
// 				Discoverer:                 test.fields.Discoverer,
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
// func TestCorrector_GetStreamListConcurrency(t *testing.T) {
// 	type fields struct {
// 		AgentPort                  int
// 		AgentName                  string
// 		AgentNamespace             string
// 		AgentDNS                   string
// 		CreationPoolSize           uint32
// 		NodeName                   string
// 		StreamListConcurrency      int
// 		BboltAsyncWriteConcurrency int
// 		Discoverer                 *DiscovererClient
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		           AgentPort:0,
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           AgentDNS:"",
// 		           CreationPoolSize:0,
// 		           NodeName:"",
// 		           StreamListConcurrency:0,
// 		           BboltAsyncWriteConcurrency:0,
// 		           Discoverer:DiscovererClient{},
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
// 		           AgentPort:0,
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           AgentDNS:"",
// 		           CreationPoolSize:0,
// 		           NodeName:"",
// 		           StreamListConcurrency:0,
// 		           BboltAsyncWriteConcurrency:0,
// 		           Discoverer:DiscovererClient{},
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
// 			c := &Corrector{
// 				AgentPort:                  test.fields.AgentPort,
// 				AgentName:                  test.fields.AgentName,
// 				AgentNamespace:             test.fields.AgentNamespace,
// 				AgentDNS:                   test.fields.AgentDNS,
// 				CreationPoolSize:           test.fields.CreationPoolSize,
// 				NodeName:                   test.fields.NodeName,
// 				StreamListConcurrency:      test.fields.StreamListConcurrency,
// 				BboltAsyncWriteConcurrency: test.fields.BboltAsyncWriteConcurrency,
// 				Discoverer:                 test.fields.Discoverer,
// 			}
//
// 			got := c.GetStreamListConcurrency()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func TestCorrector_GetBboltAsyncWriteConcurrency(t *testing.T) {
// 	type fields struct {
// 		AgentPort                  int
// 		AgentName                  string
// 		AgentNamespace             string
// 		AgentDNS                   string
// 		CreationPoolSize           uint32
// 		NodeName                   string
// 		StreamListConcurrency      int
// 		BboltAsyncWriteConcurrency int
// 		Discoverer                 *DiscovererClient
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		           AgentPort:0,
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           AgentDNS:"",
// 		           CreationPoolSize:0,
// 		           NodeName:"",
// 		           StreamListConcurrency:0,
// 		           BboltAsyncWriteConcurrency:0,
// 		           Discoverer:DiscovererClient{},
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
// 		           AgentPort:0,
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           AgentDNS:"",
// 		           CreationPoolSize:0,
// 		           NodeName:"",
// 		           StreamListConcurrency:0,
// 		           BboltAsyncWriteConcurrency:0,
// 		           Discoverer:DiscovererClient{},
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
// 			c := &Corrector{
// 				AgentPort:                  test.fields.AgentPort,
// 				AgentName:                  test.fields.AgentName,
// 				AgentNamespace:             test.fields.AgentNamespace,
// 				AgentDNS:                   test.fields.AgentDNS,
// 				CreationPoolSize:           test.fields.CreationPoolSize,
// 				NodeName:                   test.fields.NodeName,
// 				StreamListConcurrency:      test.fields.StreamListConcurrency,
// 				BboltAsyncWriteConcurrency: test.fields.BboltAsyncWriteConcurrency,
// 				Discoverer:                 test.fields.Discoverer,
// 			}
//
// 			got := c.GetBboltAsyncWriteConcurrency()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
