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

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
)

func TestAlwaysAvailable(t *testing.T) {
	c := alwaysAvailable()
	assert.True(t, c.HasGeneralPool)
	assert.True(t, c.HasAgentPool)
}

func TestResolveNodePoolCapability_NilClient(t *testing.T) {
	_, err := resolveNodePoolCapability(context.Background(), nil, "default", "")
	assert.Error(t, err)
}

func TestResolveNodePoolCapability(t *testing.T) {
	const (
		namespace = "test-ns"
		prefix    = "vald.vdaas.org"
	)

	tests := []struct {
		name           string
		nodes          []k8s.RuntimeObject
		wantHasGeneral bool
		wantHasAgent   bool
	}{
		{
			name:           "no nodes",
			nodes:          nil,
			wantHasGeneral: false,
			wantHasAgent:   false,
		},
		{
			name: "only general pool node",
			nodes: []k8s.RuntimeObject{
				makeNode("n1", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): namespace,
					labelKey(prefix, nodePoolLabelType):      string(v1.NodePoolTypeGeneral),
				}),
			},
			wantHasGeneral: true,
			wantHasAgent:   false,
		},
		{
			name: "both pools present",
			nodes: []k8s.RuntimeObject{
				makeNode("n1", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): namespace,
					labelKey(prefix, nodePoolLabelType):      string(v1.NodePoolTypeGeneral),
				}),
				makeNode("n2", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): namespace,
					labelKey(prefix, nodePoolLabelType):      string(v1.NodePoolTypeValdAgent),
				}),
			},
			wantHasGeneral: true,
			wantHasAgent:   true,
		},
		{
			name: "node in another namespace does not count",
			nodes: []k8s.RuntimeObject{
				makeNode("n1", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): "other-ns",
					labelKey(prefix, nodePoolLabelType):      string(v1.NodePoolTypeGeneral),
				}),
			},
			wantHasGeneral: false,
			wantHasAgent:   false,
		},
		{
			// The single-List implementation classifies the type label in Go:
			// unknown pool types in the namespace must not count as any pool.
			name: "node with an unknown pool type does not count",
			nodes: []k8s.RuntimeObject{
				makeNode("n1", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): namespace,
					labelKey(prefix, nodePoolLabelType):      "gpu",
				}),
			},
			wantHasGeneral: false,
			wantHasAgent:   false,
		},
	}

	scheme := k8s.NewScheme()
	assert.NoError(t, k8s.AddCoreToScheme(scheme))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mock.NewFakeClientBuilder().WithScheme(scheme).WithRuntimeObjects(tt.nodes...).Build()
			got, err := resolveNodePoolCapability(context.Background(), c, namespace, prefix)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantHasGeneral, got.HasGeneralPool, "HasGeneralPool")
			assert.Equal(t, tt.wantHasAgent, got.HasAgentPool, "HasAgentPool")
		})
	}
}

func makeNode(name string, labels map[string]string) *k8s.Node {
	return &k8s.Node{
		ObjectMeta: k8s.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_labelKey(t *testing.T) {
// 	type args struct {
// 		prefix string
// 		suffix string
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
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
// 		           prefix:"",
// 		           suffix:"",
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
// 		           prefix:"",
// 		           suffix:"",
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
// 			got := labelKey(test.args.prefix, test.args.suffix)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_alwaysAvailable(t *testing.T) {
// 	type want struct {
// 		want nodePoolCapability
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, nodePoolCapability) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got nodePoolCapability) error {
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
//
// 			got := alwaysAvailable()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resolveNodePoolCapability(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		c           k8s.Client
// 		namespace   string
// 		labelPrefix string
// 	}
// 	type want struct {
// 		want nodePoolCapability
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, nodePoolCapability, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got nodePoolCapability, err error) error {
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
// 		           ctx:nil,
// 		           c:nil,
// 		           namespace:"",
// 		           labelPrefix:"",
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
// 		           c:nil,
// 		           namespace:"",
// 		           labelPrefix:"",
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
// 			got, err := resolveNodePoolCapability(test.args.ctx, test.args.c, test.args.namespace, test.args.labelPrefix)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
