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
package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodePools_GetNodePool(t *testing.T) {
	pool := NodePool{Name: "general-pool", Replicas: 2}
	pools := NodePools{
		NodePoolTypeGeneral: pool,
	}

	t.Run("existing type returns pool", func(t *testing.T) {
		got := pools.GetNodePool(NodePoolTypeGeneral)
		assert.NotNil(t, got)
		assert.Equal(t, pool.Name, got.Name)
	})

	t.Run("non-existent type returns nil", func(t *testing.T) {
		got := pools.GetNodePool(NodePoolTypeValdAgent)
		assert.Nil(t, got)
	})

	t.Run("empty NodePools returns nil", func(t *testing.T) {
		var empty NodePools
		got := empty.GetNodePool(NodePoolTypeGeneral)
		assert.Nil(t, got)
	})
}

func TestMachineResource_GetResourceList(t *testing.T) {
	t.Run("with storage", func(t *testing.T) {
		mr := &MachineResource{
			Cpu:     "500m",
			Memory:  "1Gi",
			Storage: "10Gi",
		}
		rl := mr.GetResourceList()
		assert.Equal(t, resource.MustParse("500m"), rl[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("1Gi"), rl[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("10Gi"), rl[corev1.ResourceStorage])
	})

	t.Run("without storage omits storage key", func(t *testing.T) {
		mr := &MachineResource{
			Cpu:    "1",
			Memory: "2Gi",
		}
		rl := mr.GetResourceList()
		assert.Equal(t, resource.MustParse("1"), rl[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("2Gi"), rl[corev1.ResourceMemory])
		_, hasStorage := rl[corev1.ResourceStorage]
		assert.False(t, hasStorage)
	})
}

func newValdOperatorReleaseFixture() *ValdOperatorRelease {
	return &ValdOperatorRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fixture",
			Namespace: "default",
			Labels:    map[string]string{"role": "green"},
		},
		Spec: ValdOperatorReleaseSpec{
			Infrastructure: []ValdOperatorReleaseInfra{
				{
					Role:     "green",
					Active:   true,
					Clusters: []DestClusters{{ID: "id-1", Name: "cluster-a"}},
					NodePools: NodePools{
						NodePoolTypeGeneral: NodePool{
							Name:            "general-pool",
							MachineResource: MachineResource{Cpu: "1", Memory: "2Gi"},
							Replicas:        1,
						},
					},
				},
			},
			VectorEngine: VectorEngine{
				Name: "vald",
				Vald: Vald{
					Agent: Agent{
						Ngt:              Ngt{Dimension: 128},
						PersistentVolume: &AgentPersistentVolume{Enabled: true},
					},
					Gateway: Gateway{Ingress: &GatewayIngress{Enabled: true, Host: "vald.example"}},
				},
			},
		},
		Status: ValdOperatorReleaseStatus{
			Phase: "WaitingCreateVrs",
			Conditions: []metav1.Condition{
				{Type: "WaitingCreateVrs", Status: metav1.ConditionTrue, Reason: "Progressing"},
			},
		},
	}
}

func TestValdOperatorRelease_DeepCopy(t *testing.T) {
	orig := newValdOperatorReleaseFixture()
	cp := orig.DeepCopy()
	if !assert.NotNil(t, cp) {
		return
	}

	// mutate every reference field of the copy
	cp.Labels["role"] = "mutated"
	cp.Spec.Infrastructure[0].Clusters[0].Name = "mutated"
	np := cp.Spec.Infrastructure[0].NodePools[NodePoolTypeGeneral]
	np.Name = "mutated"
	cp.Spec.Infrastructure[0].NodePools[NodePoolTypeGeneral] = np
	cp.Spec.VectorEngine.Vald.Agent.PersistentVolume.Enabled = false
	cp.Spec.VectorEngine.Vald.Gateway.Ingress.Host = "mutated"
	cp.Status.Conditions[0].Status = metav1.ConditionFalse

	// the original must be unaffected
	assert.Equal(t, "green", orig.Labels["role"])
	assert.Equal(t, "cluster-a", orig.Spec.Infrastructure[0].Clusters[0].Name)
	assert.Equal(t, "general-pool", orig.Spec.Infrastructure[0].NodePools[NodePoolTypeGeneral].Name)
	assert.True(t, orig.Spec.VectorEngine.Vald.Agent.PersistentVolume.Enabled)
	assert.Equal(t, "vald.example", orig.Spec.VectorEngine.Vald.Gateway.Ingress.Host)
	assert.Equal(t, metav1.ConditionTrue, orig.Status.Conditions[0].Status)
}

func TestValdOperatorRelease_DeepCopyObject(t *testing.T) {
	orig := newValdOperatorReleaseFixture()

	obj := orig.DeepCopyObject()
	if !assert.NotNil(t, obj) {
		return
	}
	cp, ok := obj.(*ValdOperatorRelease)
	assert.True(t, ok, "DeepCopyObject() = %T, want *ValdOperatorRelease", obj)
	assert.Equal(t, orig.Name, cp.GetName())

	list := &ValdOperatorReleaseList{Items: []ValdOperatorRelease{*orig}}
	lobj := list.DeepCopyObject()
	lcp, ok := lobj.(*ValdOperatorReleaseList)
	assert.True(t, ok, "DeepCopyObject() = %T, want *ValdOperatorReleaseList", lobj)
	if assert.Len(t, lcp.Items, 1) {
		lcp.Items[0].Labels["role"] = "mutated"
		assert.Equal(t, "green", orig.Labels["role"])
	}
}

func TestValdOperatorReleaseStatus_DeepCopy(t *testing.T) {
	orig := newValdOperatorReleaseFixture()
	st := orig.Status.DeepCopy()
	if !assert.NotNil(t, st) {
		return
	}
	st.Conditions[0].Status = metav1.ConditionFalse
	assert.Equal(t, metav1.ConditionTrue, orig.Status.Conditions[0].Status)
}

// NOT IMPLEMENTED BELOW
//
// func Test_parseQuantity(t *testing.T) {
// 	type args struct {
// 		s string
// 	}
// 	type want struct {
// 		want resource.Quantity
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, resource.Quantity) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got resource.Quantity) error {
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
// 		           s:"",
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
// 		           s:"",
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
// 			got := parseQuantity(test.args.s)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
