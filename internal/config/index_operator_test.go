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
// func TestIndexJobTemplates_Bind(t *testing.T) {
// 	type fields struct {
// 		Rotate     *k8s.Job
// 		Creation   *k8s.Job
// 		Save       *k8s.Job
// 		Correction *k8s.Job
// 	}
// 	type want struct {
// 		want *IndexJobTemplates
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *IndexJobTemplates) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *IndexJobTemplates) error {
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
// 		           Rotate:nil,
// 		           Creation:nil,
// 		           Save:nil,
// 		           Correction:nil,
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
// 		           Rotate:nil,
// 		           Creation:nil,
// 		           Save:nil,
// 		           Correction:nil,
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
// 			ijt := &IndexJobTemplates{
// 				Rotate:     test.fields.Rotate,
// 				Creation:   test.fields.Creation,
// 				Save:       test.fields.Save,
// 				Correction: test.fields.Correction,
// 			}
//
// 			got := ijt.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestIndexOperator_Bind(t *testing.T) {
// 	type fields struct {
// 		Namespace                         string
// 		AgentName                         string
// 		AgentNamespace                    string
// 		RotatorName                       string
// 		TargetReadReplicaIDAnnotationsKey string
// 		RotationJobConcurrency            uint
// 		ReadReplicaEnabled                bool
// 		ReadReplicaLabelKey               string
// 		JobTemplates                      IndexJobTemplates
// 	}
// 	type want struct {
// 		want *IndexOperator
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *IndexOperator) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *IndexOperator) error {
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
// 		           Namespace:"",
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           RotatorName:"",
// 		           TargetReadReplicaIDAnnotationsKey:"",
// 		           RotationJobConcurrency:0,
// 		           ReadReplicaEnabled:false,
// 		           ReadReplicaLabelKey:"",
// 		           JobTemplates:IndexJobTemplates{},
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
// 		           Namespace:"",
// 		           AgentName:"",
// 		           AgentNamespace:"",
// 		           RotatorName:"",
// 		           TargetReadReplicaIDAnnotationsKey:"",
// 		           RotationJobConcurrency:0,
// 		           ReadReplicaEnabled:false,
// 		           ReadReplicaLabelKey:"",
// 		           JobTemplates:IndexJobTemplates{},
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
// 			ic := &IndexOperator{
// 				Namespace:                         test.fields.Namespace,
// 				AgentName:                         test.fields.AgentName,
// 				AgentNamespace:                    test.fields.AgentNamespace,
// 				RotatorName:                       test.fields.RotatorName,
// 				TargetReadReplicaIDAnnotationsKey: test.fields.TargetReadReplicaIDAnnotationsKey,
// 				RotationJobConcurrency:            test.fields.RotationJobConcurrency,
// 				ReadReplicaEnabled:                test.fields.ReadReplicaEnabled,
// 				ReadReplicaLabelKey:               test.fields.ReadReplicaLabelKey,
// 				JobTemplates:                      test.fields.JobTemplates,
// 			}
//
// 			got := ic.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
