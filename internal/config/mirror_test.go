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
// func TestMirror_Bind(t *testing.T) {
// 	type fields struct {
// 		Net               *Net
// 		Client            *GRPCClient
// 		SelfMirrorAddr    string
// 		GatewayAddr       string
// 		PodName           string
// 		RegisterDuration  string
// 		Namespace         string
// 		DiscoveryDuration string
// 		Colocation        string
// 		Group             string
// 	}
// 	type want struct {
// 		want *Mirror
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *Mirror) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *Mirror) error {
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
// 		           Net:Net{},
// 		           Client:GRPCClient{},
// 		           SelfMirrorAddr:"",
// 		           GatewayAddr:"",
// 		           PodName:"",
// 		           RegisterDuration:"",
// 		           Namespace:"",
// 		           DiscoveryDuration:"",
// 		           Colocation:"",
// 		           Group:"",
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
// 		           Net:Net{},
// 		           Client:GRPCClient{},
// 		           SelfMirrorAddr:"",
// 		           GatewayAddr:"",
// 		           PodName:"",
// 		           RegisterDuration:"",
// 		           Namespace:"",
// 		           DiscoveryDuration:"",
// 		           Colocation:"",
// 		           Group:"",
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
// 			m := &Mirror{
// 				Net:               test.fields.Net,
// 				Client:            test.fields.Client,
// 				SelfMirrorAddr:    test.fields.SelfMirrorAddr,
// 				GatewayAddr:       test.fields.GatewayAddr,
// 				PodName:           test.fields.PodName,
// 				RegisterDuration:  test.fields.RegisterDuration,
// 				Namespace:         test.fields.Namespace,
// 				DiscoveryDuration: test.fields.DiscoveryDuration,
// 				Colocation:        test.fields.Colocation,
// 				Group:             test.fields.Group,
// 			}
//
// 			got := m.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
