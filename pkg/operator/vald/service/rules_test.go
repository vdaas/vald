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

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
)

func TestResolveAgentNodePool(t *testing.T) {
	makeInfra := func(generalReplicas int, agentReplicas int, withAgent bool) v1.ValdOperatorReleaseInfra {
		pools := v1.NodePools{
			v1.NodePoolTypeGeneral: v1.NodePool{
				Name:            "general",
				Replicas:        generalReplicas,
				MachineResource: v1.MachineResource{Cpu: "4", Memory: "8Gi"},
			},
		}
		if withAgent {
			pools[v1.NodePoolTypeValdAgent] = v1.NodePool{
				Name:            "agent",
				Replicas:        agentReplicas,
				MachineResource: v1.MachineResource{Cpu: "16", Memory: "32Gi"},
			}
		}
		return v1.ValdOperatorReleaseInfra{NodePools: pools}
	}

	t.Run("agent pool present with replicas: use agent", func(t *testing.T) {
		got := resolveAgentNodePool(makeInfra(3, 2, true))
		assert.Equal(t, 2, got.NodeCount)
		assert.Equal(t, resource.MustParse("16"), got.MachineResource[k8s.ResourceCPU])
	})

	t.Run("agent pool present but replicas == 0: fall back to general", func(t *testing.T) {
		got := resolveAgentNodePool(makeInfra(3, 0, true))
		assert.Equal(t, 3, got.NodeCount)
		assert.Equal(t, resource.MustParse("4"), got.MachineResource[k8s.ResourceCPU])
	})

	t.Run("agent pool absent: fall back to general", func(t *testing.T) {
		got := resolveAgentNodePool(makeInfra(3, 0, false))
		assert.Equal(t, 3, got.NodeCount)
		assert.Equal(t, resource.MustParse("4"), got.MachineResource[k8s.ResourceCPU])
	})

	t.Run("no pools at all: empty spec without panicking", func(t *testing.T) {
		got := resolveAgentNodePool(v1.ValdOperatorReleaseInfra{})
		assert.Zero(t, got.NodeCount)
		assert.Nil(t, got.MachineResource)
	})

	t.Run("agent pool with zero replicas and general pool absent: empty spec", func(t *testing.T) {
		got := resolveAgentNodePool(v1.ValdOperatorReleaseInfra{
			NodePools: v1.NodePools{
				v1.NodePoolTypeValdAgent: v1.NodePool{
					Name:            "agent",
					Replicas:        0,
					MachineResource: v1.MachineResource{Cpu: "16", Memory: "32Gi"},
				},
			},
		})
		assert.Zero(t, got.NodeCount)
		assert.Nil(t, got.MachineResource)
	})
}

func TestAgentPvSize(t *testing.T) {
	const oneGi = int64(1) << 30
	const oneMi = int64(1) << 20

	tests := []struct {
		name         string
		memoryBytes  int64
		bufferRatio  float64
		minSizeBytes int64
		wantBytes    int64
	}{
		{"memory * ratio above min: ratio applies", 4 * oneGi, 1.5, oneGi, 6 * oneGi},
		{"memory * ratio below min: min applies", 100 * oneMi, 1.5, oneGi, oneGi},
		{"configurable min raises the floor", 100 * oneMi, 1.5, 2 * oneGi, 2 * oneGi},
		{"configurable ratio scales the calculation", 4 * oneGi, 2.0, oneGi, 8 * oneGi},
		{"rounds up to whole Gi", 1500 * oneMi, 1.0, oneGi, 2 * oneGi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agentPvSize(tt.memoryBytes, tt.bufferRatio, tt.minSizeBytes)
			gotQty := resource.MustParse(got)
			assert.Equal(t, tt.wantBytes, gotQty.Value(), "PV size in bytes")
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_resolveAgentNodePool(t *testing.T) {
// 	type args struct {
// 		infra v1.ValdOperatorReleaseInfra
// 	}
// 	type want struct {
// 		want agentNodePoolSpec
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, agentNodePoolSpec) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got agentNodePoolSpec) error {
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
// 		           infra:nil,
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
// 		           infra:nil,
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
// 			got := resolveAgentNodePool(test.args.infra)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_agentPvSize(t *testing.T) {
// 	type args struct {
// 		memoryBytes    int64
// 		pvBufferRatio  float64
// 		pvMinSizeBytes int64
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
// 		           memoryBytes:0,
// 		           pvBufferRatio:0,
// 		           pvMinSizeBytes:0,
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
// 		           memoryBytes:0,
// 		           pvBufferRatio:0,
// 		           pvMinSizeBytes:0,
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
// 			got := agentPvSize(test.args.memoryBytes, test.args.pvBufferRatio, test.args.pvMinSizeBytes)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
