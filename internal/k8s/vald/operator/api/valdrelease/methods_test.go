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

package valdrelease

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// --- Agent -----------------------------------------------------------------

func TestAgent_SetReplica(t *testing.T) {
	const podsPerNode = 2
	tests := []struct {
		name         string
		nodeReplicas int
		want         int
	}{
		{"2 nodes × 2 pods", 2, 4},
		{"3 nodes × 2 pods", 3, 6},
		{"1 node × 2 pods", 1, 2},
		{"0 nodes", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{}
			a.SetReplica(tt.nodeReplicas, podsPerNode)
			assert.Equal(t, tt.want, *a.MinReplicas)
			assert.Equal(t, tt.want, *a.MaxReplicas)
		})
	}
}

// TestAgent_SetResources_Values verifies exact resource values for representative node specs.
//
// Formula:
//
//	CPU request  = nodeCPU_cores × AgentResourceRatio / podsPerNode  [millicores]
//	CPU limit    = nodeCPU_cores × AgentResourceRatio                [millicores]
//	RAM request  = nodeRAM_bytes × AgentResourceRatio / podsPerNode  [rounded up to M]
//	RAM limit    = NOT SET (NGT index grows post-startup; hard limit causes OOM kills)
const testNodeRAM = "10000M"

func TestAgent_SetResources_Values(t *testing.T) {
	tests := []struct {
		name            string
		nodeCPU         string
		nodeRAM         string
		podsPerNode     int
		wantCPUReqMilli int64
		wantCPULimMilli int64
		wantRAMReqBytes int64
	}{
		{
			// 16 × 0.6 / 2 = 4800m req, 9600m lim; 10000M × 0.6 / 2 = 3000M req
			name:        "16CPU 10000M node, 2 pods/node",
			podsPerNode: 2, nodeCPU: "16", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 4800, wantCPULimMilli: 9600,
			wantRAMReqBytes: 3_000_000_000,
		},
		{
			// 8 × 0.6 / 2 = 2400m req, 4800m lim
			name:        "8CPU 10000M node, 2 pods/node",
			podsPerNode: 2, nodeCPU: "8", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 2400, wantCPULimMilli: 4800,
			wantRAMReqBytes: 3_000_000_000,
		},
		{
			// podsPerNode=1: req == lim
			name:        "16CPU node, 1 pod/node",
			podsPerNode: 1, nodeCPU: "16", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 9600, wantCPULimMilli: 9600,
			wantRAMReqBytes: 6_000_000_000,
		},
		{
			// podsPerNode=4: req = lim / 4
			name:        "16CPU node, 4 pods/node",
			podsPerNode: 4, nodeCPU: "16", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 2400, wantCPULimMilli: 9600,
			wantRAMReqBytes: 1_500_000_000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(tt.nodeCPU),
				corev1.ResourceMemory: resource.MustParse(tt.nodeRAM),
			}
			a := &Agent{}
			a.SetResources(mc, tt.podsPerNode)

			assert.NotNil(t, a.Resources)

			gotCPUReq := a.Resources.Requests[corev1.ResourceCPU]
			gotCPULim := a.Resources.Limits[corev1.ResourceCPU]
			gotRAMReq := a.Resources.Requests[corev1.ResourceMemory]

			assert.Equal(t, tt.wantCPUReqMilli, gotCPUReq.MilliValue(), "CPU request (milli)")
			assert.Equal(t, tt.wantCPULimMilli, gotCPULim.MilliValue(), "CPU limit (milli)")
			assert.Equal(t, tt.wantRAMReqBytes, gotRAMReq.Value(), "RAM request (bytes)")

			// Structural invariant: limit covers all pods on the node, request is per-pod.
			assert.Equal(
				t,
				gotCPUReq.MilliValue()*int64(tt.podsPerNode),
				gotCPULim.MilliValue(),
				"CPU limit must equal request × podsPerNode",
			)

			// Memory limit must NOT be set.
			_, hasMemLim := a.Resources.Limits[corev1.ResourceMemory]
			assert.False(t, hasMemLim, "memory limit must not be set (NGT index grows post-startup)")
		})
	}
}

// TestAgent_SetPvEnable verifies that calling SetPvEnable populates NGT
// settings and PersistentVolume from the caller-provided values.
func TestAgent_SetPvEnable(t *testing.T) {
	a := &Agent{}
	a.SetPvEnable("fast-ssd", "ReadWriteOnce", "6Gi")

	assert.True(t, *a.Ngt.EnableCopyOnWrite)
	assert.False(t, *a.Ngt.EnableInMemoryMode)
	assert.Equal(t, DefaultIndexPath, *a.Ngt.IndexPath)

	assert.NotNil(t, a.PersistentVolume)
	assert.True(t, *a.PersistentVolume.Enabled)
	assert.Equal(t, "fast-ssd", *a.PersistentVolume.StorageClass)
	assert.Equal(t, "ReadWriteOnce", *a.PersistentVolume.AccessMode)
	assert.Equal(t, "6Gi", *a.PersistentVolume.Size)
}

// --- Gateway (LB) ----------------------------------------------------------

func TestGatewayLb_SetReplica(t *testing.T) {
	tests := []struct {
		name     string
		agentMax int
		wantMin  int
		wantMax  int
	}{
		{"6 agent replicas", 6, 3, 12},
		{"1 agent replica → min stays 1", 1, 1, 2},
		{"2 agent replicas", 2, 1, 4},
		// ar=0: both getMaxReplica(0)=1, getMinReplica(0)=1
		{"0 agent replicas → floor to 1", 0, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := &GatewayLb{}
			a := &Agent{MaxReplicas: new(tt.agentMax)}
			lb.SetReplica(a)
			assert.Equal(t, tt.wantMin, *lb.MinReplicas)
			assert.Equal(t, tt.wantMax, *lb.MaxReplicas)
		})
	}
}

// --- Discoverer ------------------------------------------------------------

func TestDiscoverer_ApplyDefaultsByKind(t *testing.T) {
	const (
		maxSurge       = "30%"
		maxUnavailable = "0%"
	)

	// Empty string in the existing/want columns represents a nil pointer on the
	// generated type (the field is omitted).
	tests := []struct {
		name                  string
		kind                  string
		existingServiceType   string
		existingTrafficPolicy string
		wantServiceType       string
		wantTrafficPolicy     string
		wantRollingUpdate     bool
	}{
		{
			name:              "DaemonSet sets NodePort and Local defaults",
			kind:              string(DiscovererKindDaemonSet),
			wantServiceType:   string(corev1.ServiceTypeNodePort),
			wantTrafficPolicy: string(corev1.ServiceExternalTrafficPolicyTypeLocal),
			wantRollingUpdate: true,
		},
		{
			name:                "DaemonSet does not override existing ServiceType",
			kind:                string(DiscovererKindDaemonSet),
			existingServiceType: string(corev1.ServiceTypeClusterIP),
			wantServiceType:     string(corev1.ServiceTypeClusterIP),
			wantTrafficPolicy:   string(corev1.ServiceExternalTrafficPolicyTypeLocal),
			wantRollingUpdate:   true,
		},
		{
			name:                  "DaemonSet does not override existing ExternalTrafficPolicy",
			kind:                  string(DiscovererKindDaemonSet),
			existingTrafficPolicy: string(corev1.ServiceExternalTrafficPolicyTypeCluster),
			wantServiceType:       string(corev1.ServiceTypeNodePort),
			wantTrafficPolicy:     string(corev1.ServiceExternalTrafficPolicyTypeCluster),
			wantRollingUpdate:     true,
		},
		{
			name:              "Deployment sets ClusterIP",
			kind:              string(DiscovererKindDeployment),
			wantServiceType:   string(corev1.ServiceTypeClusterIP),
			wantRollingUpdate: false,
		},
		{
			name:              "unknown kind makes no change",
			kind:              "Unknown",
			wantServiceType:   "",
			wantTrafficPolicy: "",
			wantRollingUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{Kind: new(DiscovererKind(tt.kind))}
			if tt.existingServiceType != "" {
				d.ServiceType = new(DiscovererServiceType(tt.existingServiceType))
			}
			if tt.existingTrafficPolicy != "" {
				d.ExternalTrafficPolicy = new(tt.existingTrafficPolicy)
			}
			d.ApplyDefaultsByKind(maxSurge, maxUnavailable)

			if tt.wantServiceType == "" {
				assert.Nil(t, d.ServiceType)
			} else {
				assert.NotNil(t, d.ServiceType)
				assert.Equal(t, tt.wantServiceType, string(*d.ServiceType))
			}

			if tt.wantTrafficPolicy == "" {
				assert.Nil(t, d.ExternalTrafficPolicy)
			} else {
				assert.NotNil(t, d.ExternalTrafficPolicy)
				assert.Equal(t, tt.wantTrafficPolicy, *d.ExternalTrafficPolicy)
			}

			if tt.wantRollingUpdate {
				assert.NotNil(t, d.RollingUpdate)
				assert.Equal(t, maxSurge, *d.RollingUpdate.MaxSurge)
				assert.Equal(t, maxUnavailable, *d.RollingUpdate.MaxUnavailable)
			} else {
				assert.Nil(t, d.RollingUpdate)
			}
		})
	}
}

// --- Shared component behavior (SetTopologySpreadConstraints / SetResources) -
//
// Agent, GatewayLb, Discoverer, and ManagerIndex all delegate
// SetTopologySpreadConstraints to the shared setTopologySpreadConstraints
// helper, and GatewayLb/Discoverer/ManagerIndex delegate their fixed
// SetResources to the shared fixedResources helper (Agent's SetResources is
// derived from node sizing and is covered separately by
// TestAgent_SetResources_Values above). These table-driven tests cover all
// four/three components through that shared machinery in one place.

func TestSetTopologySpreadConstraints(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *TopologySpreadConstraints
		wantLabel string
	}{
		{
			name: "Agent",
			setup: func() *TopologySpreadConstraints {
				a := &Agent{}
				a.SetTopologySpreadConstraints()
				return a.TopologySpreadConstraints
			},
			wantLabel: "agent",
		},
		{
			name: "GatewayLb",
			setup: func() *TopologySpreadConstraints {
				l := &GatewayLb{}
				l.SetTopologySpreadConstraints()
				return l.TopologySpreadConstraints
			},
			wantLabel: "gateway-lb",
		},
		{
			name: "Discoverer",
			setup: func() *TopologySpreadConstraints {
				d := &Discoverer{}
				d.SetTopologySpreadConstraints()
				return d.TopologySpreadConstraints
			},
			wantLabel: "discoverer",
		},
		{
			name: "ManagerIndex",
			setup: func() *TopologySpreadConstraints {
				i := &ManagerIndex{}
				i.SetTopologySpreadConstraints()
				return i.TopologySpreadConstraints
			},
			wantLabel: "manager-index",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.setup()

			assert.NotNil(t, got)
			assert.Len(t, *got, 1)
			tsc := (*got)[0]
			assert.Equal(t, int32(1), tsc.MaxSkew)
			assert.Equal(t, "kubernetes.io/hostname", tsc.TopologyKey)
			assert.Equal(t, corev1.DoNotSchedule, tsc.WhenUnsatisfiable)
			assert.NotNil(t, tsc.LabelSelector)
			assert.Equal(t, tt.wantLabel, tsc.LabelSelector.MatchLabels["app.kubernetes.io/component"])
		})
	}
}

func TestSetFixedResources(t *testing.T) {
	// commonReqCPU is the CPU request shared by all three fixed-resource
	// components below (GatewayLb, Discoverer, ManagerIndex).
	const commonReqCPU = "200m"

	tests := []struct {
		name          string
		setup         func() *Resources
		wantReqCPU    string
		wantReqMemory string
		wantLimCPU    string
		wantLimMemory string
	}{
		{
			name: "GatewayLb",
			setup: func() *Resources {
				l := &GatewayLb{}
				l.SetResources()
				return l.Resources
			},
			wantReqCPU: commonReqCPU, wantReqMemory: "150Mi",
			wantLimCPU: "2000m", wantLimMemory: "700Mi",
		},
		{
			name: "Discoverer",
			setup: func() *Resources {
				d := &Discoverer{}
				d.SetResources()
				return d.Resources
			},
			wantReqCPU: commonReqCPU, wantReqMemory: "65Mi",
			wantLimCPU: "600m", wantLimMemory: "200Mi",
		},
		{
			name: "ManagerIndex",
			setup: func() *Resources {
				i := &ManagerIndex{}
				i.SetResources()
				return i.Resources
			},
			wantReqCPU: commonReqCPU, wantReqMemory: "80Mi",
			wantLimCPU: "1000m", wantLimMemory: "500Mi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.setup()

			assert.NotNil(t, got)
			assert.Equal(t, resource.MustParse(tt.wantReqCPU), got.Requests[corev1.ResourceCPU])
			assert.Equal(t, resource.MustParse(tt.wantReqMemory), got.Requests[corev1.ResourceMemory])
			assert.Equal(t, resource.MustParse(tt.wantLimCPU), got.Limits[corev1.ResourceCPU])
			assert.Equal(t, resource.MustParse(tt.wantLimMemory), got.Limits[corev1.ResourceMemory])
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_setTopologySpreadConstraints(t *testing.T) {
// 	type args struct {
// 		componentLabel string
// 	}
// 	type want struct {
// 		want *TopologySpreadConstraints
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *TopologySpreadConstraints) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *TopologySpreadConstraints) error {
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
// 		           componentLabel:"",
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
// 		           componentLabel:"",
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
// 			got := setTopologySpreadConstraints(test.args.componentLabel)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_fixedResources(t *testing.T) {
// 	type args struct {
// 		reqCPU    string
// 		reqMemory string
// 		limCPU    string
// 		limMemory string
// 	}
// 	type want struct {
// 		want *Resources
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *Resources) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *Resources) error {
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
// 		           reqCPU:"",
// 		           reqMemory:"",
// 		           limCPU:"",
// 		           limMemory:"",
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
// 		           reqCPU:"",
// 		           reqMemory:"",
// 		           limCPU:"",
// 		           limMemory:"",
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
// 			got := fixedResources(test.args.reqCPU, test.args.reqMemory, test.args.limCPU, test.args.limMemory)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestAgent_SetResources(t *testing.T) {
// 	type args struct {
// 		mc          corev1.ResourceList
// 		podsPerNode int
// 	}
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Algorithm                     *AgentAlgorithm
// 		Annotations                   *map[string]any
// 		ClusterRole                   *AgentClusterRole
// 		ClusterRoleBinding            *AgentClusterRoleBinding
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		Faiss                         *AgentFaiss
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		InitContainers                *InitContainers
// 		Kind                          *AgentKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		Ngt                           *AgentNgt
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PersistentVolume              *AgentPersistentVolume
// 		PodAnnotations                *map[string]any
// 		PodManagementPolicy           *AgentPodManagementPolicy
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Readreplica                   *AgentReadreplica
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *AgentRollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *AgentServiceType
// 		Sidecar                       *AgentSidecar
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *AgentUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           mc:nil,
// 		           podsPerNode:0,
// 		       },
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Algorithm:nil,
// 		           Annotations:nil,
// 		           ClusterRole:AgentClusterRole{},
// 		           ClusterRoleBinding:AgentClusterRoleBinding{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Faiss:AgentFaiss{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           Ngt:AgentNgt{},
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PersistentVolume:AgentPersistentVolume{},
// 		           PodAnnotations:nil,
// 		           PodManagementPolicy:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:AgentReadreplica{},
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:AgentRollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           Sidecar:AgentSidecar{},
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           mc:nil,
// 		           podsPerNode:0,
// 		           },
// 		           fields: fields {
// 		           Affinity:nil,
// 		           Algorithm:nil,
// 		           Annotations:nil,
// 		           ClusterRole:AgentClusterRole{},
// 		           ClusterRoleBinding:AgentClusterRoleBinding{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Faiss:AgentFaiss{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           Ngt:AgentNgt{},
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PersistentVolume:AgentPersistentVolume{},
// 		           PodAnnotations:nil,
// 		           PodManagementPolicy:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:AgentReadreplica{},
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:AgentRollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           Sidecar:AgentSidecar{},
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			a := &Agent{
// 				Affinity:                      test.fields.Affinity,
// 				Algorithm:                     test.fields.Algorithm,
// 				Annotations:                   test.fields.Annotations,
// 				ClusterRole:                   test.fields.ClusterRole,
// 				ClusterRoleBinding:            test.fields.ClusterRoleBinding,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				Faiss:                         test.fields.Faiss,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				InitContainers:                test.fields.InitContainers,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				Ngt:                           test.fields.Ngt,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PersistentVolume:              test.fields.PersistentVolume,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodManagementPolicy:           test.fields.PodManagementPolicy,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Readreplica:                   test.fields.Readreplica,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				Sidecar:                       test.fields.Sidecar,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			a.SetResources(test.args.mc, test.args.podsPerNode)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestAgent_SetTopologySpreadConstraints(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Algorithm                     *AgentAlgorithm
// 		Annotations                   *map[string]any
// 		ClusterRole                   *AgentClusterRole
// 		ClusterRoleBinding            *AgentClusterRoleBinding
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		Faiss                         *AgentFaiss
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		InitContainers                *InitContainers
// 		Kind                          *AgentKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		Ngt                           *AgentNgt
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PersistentVolume              *AgentPersistentVolume
// 		PodAnnotations                *map[string]any
// 		PodManagementPolicy           *AgentPodManagementPolicy
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Readreplica                   *AgentReadreplica
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *AgentRollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *AgentServiceType
// 		Sidecar                       *AgentSidecar
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *AgentUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Algorithm:nil,
// 		           Annotations:nil,
// 		           ClusterRole:AgentClusterRole{},
// 		           ClusterRoleBinding:AgentClusterRoleBinding{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Faiss:AgentFaiss{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           Ngt:AgentNgt{},
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PersistentVolume:AgentPersistentVolume{},
// 		           PodAnnotations:nil,
// 		           PodManagementPolicy:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:AgentReadreplica{},
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:AgentRollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           Sidecar:AgentSidecar{},
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Algorithm:nil,
// 		           Annotations:nil,
// 		           ClusterRole:AgentClusterRole{},
// 		           ClusterRoleBinding:AgentClusterRoleBinding{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Faiss:AgentFaiss{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           Ngt:AgentNgt{},
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PersistentVolume:AgentPersistentVolume{},
// 		           PodAnnotations:nil,
// 		           PodManagementPolicy:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:AgentReadreplica{},
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:AgentRollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           Sidecar:AgentSidecar{},
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			a := &Agent{
// 				Affinity:                      test.fields.Affinity,
// 				Algorithm:                     test.fields.Algorithm,
// 				Annotations:                   test.fields.Annotations,
// 				ClusterRole:                   test.fields.ClusterRole,
// 				ClusterRoleBinding:            test.fields.ClusterRoleBinding,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				Faiss:                         test.fields.Faiss,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				InitContainers:                test.fields.InitContainers,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				Ngt:                           test.fields.Ngt,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PersistentVolume:              test.fields.PersistentVolume,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodManagementPolicy:           test.fields.PodManagementPolicy,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Readreplica:                   test.fields.Readreplica,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				Sidecar:                       test.fields.Sidecar,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			a.SetTopologySpreadConstraints()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestGatewayLb_getMaxReplica(t *testing.T) {
// 	type args struct {
// 		ar int
// 	}
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		GatewayConfig                 *GatewayLbGatewayConfig
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		Ingress                       *GatewayLbIngress
// 		InitContainers                *InitContainers
// 		InternalTrafficPolicy         *string
// 		Kind                          *GatewayLbKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *GatewayLbServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *GatewayLbUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
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
// 		       args: args {
// 		           ar:0,
// 		       },
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           ar:0,
// 		           },
// 		           fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			l := &GatewayLb{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				GatewayConfig:                 test.fields.GatewayConfig,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				Ingress:                       test.fields.Ingress,
// 				InitContainers:                test.fields.InitContainers,
// 				InternalTrafficPolicy:         test.fields.InternalTrafficPolicy,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			got := l.getMaxReplica(test.args.ar)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestGatewayLb_getMinReplica(t *testing.T) {
// 	type args struct {
// 		ar int
// 	}
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		GatewayConfig                 *GatewayLbGatewayConfig
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		Ingress                       *GatewayLbIngress
// 		InitContainers                *InitContainers
// 		InternalTrafficPolicy         *string
// 		Kind                          *GatewayLbKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *GatewayLbServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *GatewayLbUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
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
// 		       args: args {
// 		           ar:0,
// 		       },
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           ar:0,
// 		           },
// 		           fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			l := &GatewayLb{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				GatewayConfig:                 test.fields.GatewayConfig,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				Ingress:                       test.fields.Ingress,
// 				InitContainers:                test.fields.InitContainers,
// 				InternalTrafficPolicy:         test.fields.InternalTrafficPolicy,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			got := l.getMinReplica(test.args.ar)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestGatewayLb_SetResources(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		GatewayConfig                 *GatewayLbGatewayConfig
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		Ingress                       *GatewayLbIngress
// 		InitContainers                *InitContainers
// 		InternalTrafficPolicy         *string
// 		Kind                          *GatewayLbKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *GatewayLbServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *GatewayLbUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			l := &GatewayLb{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				GatewayConfig:                 test.fields.GatewayConfig,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				Ingress:                       test.fields.Ingress,
// 				InitContainers:                test.fields.InitContainers,
// 				InternalTrafficPolicy:         test.fields.InternalTrafficPolicy,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			l.SetResources()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestGatewayLb_SetTopologySpreadConstraints(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		GatewayConfig                 *GatewayLbGatewayConfig
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		Ingress                       *GatewayLbIngress
// 		InitContainers                *InitContainers
// 		InternalTrafficPolicy         *string
// 		Kind                          *GatewayLbKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *GatewayLbServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *GatewayLbUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           GatewayConfig:GatewayLbGatewayConfig{},
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           Ingress:GatewayLbIngress{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			l := &GatewayLb{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				GatewayConfig:                 test.fields.GatewayConfig,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				Ingress:                       test.fields.Ingress,
// 				InitContainers:                test.fields.InitContainers,
// 				InternalTrafficPolicy:         test.fields.InternalTrafficPolicy,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			l.SetTopologySpreadConstraints()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestDiscoverer_SetResources(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		ClusterRole                   *DiscovererClusterRole
// 		ClusterRoleBinding            *DiscovererClusterRoleBinding
// 		Discoverer                    *DiscovererDiscoverer
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		InitContainers                *InitContainers
// 		InternalTrafficPolicy         *string
// 		Kind                          *DiscovererKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *DiscovererServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *DiscovererUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           ClusterRole:DiscovererClusterRole{},
// 		           ClusterRoleBinding:DiscovererClusterRoleBinding{},
// 		           Discoverer:DiscovererDiscoverer{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           ClusterRole:DiscovererClusterRole{},
// 		           ClusterRoleBinding:DiscovererClusterRoleBinding{},
// 		           Discoverer:DiscovererDiscoverer{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			d := &Discoverer{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				ClusterRole:                   test.fields.ClusterRole,
// 				ClusterRoleBinding:            test.fields.ClusterRoleBinding,
// 				Discoverer:                    test.fields.Discoverer,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				InitContainers:                test.fields.InitContainers,
// 				InternalTrafficPolicy:         test.fields.InternalTrafficPolicy,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			d.SetResources()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestDiscoverer_SetTopologySpreadConstraints(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		ClusterRole                   *DiscovererClusterRole
// 		ClusterRoleBinding            *DiscovererClusterRoleBinding
// 		Discoverer                    *DiscovererDiscoverer
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		Hpa                           *Hpa
// 		Image                         *Image
// 		InitContainers                *InitContainers
// 		InternalTrafficPolicy         *string
// 		Kind                          *DiscovererKind
// 		Logging                       *Logging
// 		MaxReplicas                   *int
// 		MaxUnavailable                *string
// 		MinReplicas                   *int
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *DiscovererServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *DiscovererUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           ClusterRole:DiscovererClusterRole{},
// 		           ClusterRoleBinding:DiscovererClusterRoleBinding{},
// 		           Discoverer:DiscovererDiscoverer{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           ClusterRole:DiscovererClusterRole{},
// 		           ClusterRoleBinding:DiscovererClusterRoleBinding{},
// 		           Discoverer:DiscovererDiscoverer{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Hpa:Hpa{},
// 		           Image:Image{},
// 		           InitContainers:nil,
// 		           InternalTrafficPolicy:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxReplicas:nil,
// 		           MaxUnavailable:nil,
// 		           MinReplicas:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			d := &Discoverer{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				ClusterRole:                   test.fields.ClusterRole,
// 				ClusterRoleBinding:            test.fields.ClusterRoleBinding,
// 				Discoverer:                    test.fields.Discoverer,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				Hpa:                           test.fields.Hpa,
// 				Image:                         test.fields.Image,
// 				InitContainers:                test.fields.InitContainers,
// 				InternalTrafficPolicy:         test.fields.InternalTrafficPolicy,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxReplicas:                   test.fields.MaxReplicas,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				MinReplicas:                   test.fields.MinReplicas,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			d.SetTopologySpreadConstraints()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestManagerIndex_SetResources(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		Corrector                     *ManagerIndexCorrector
// 		Creator                       *ManagerIndexCreator
// 		Deleter                       *ManagerIndexDeleter
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		Image                         *Image
// 		Indexer                       *ManagerIndexIndexer
// 		InitContainers                *InitContainers
// 		Kind                          *ManagerIndexKind
// 		Logging                       *Logging
// 		MaxUnavailable                *string
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		Operator                      *ManagerIndexOperator
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Readreplica                   *ManagerIndexReadreplica
// 		Replicas                      *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		Saver                         *ManagerIndexSaver
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *ManagerIndexServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *ManagerIndexUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Corrector:ManagerIndexCorrector{},
// 		           Creator:ManagerIndexCreator{},
// 		           Deleter:ManagerIndexDeleter{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Image:Image{},
// 		           Indexer:ManagerIndexIndexer{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxUnavailable:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           Operator:ManagerIndexOperator{},
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:ManagerIndexReadreplica{},
// 		           Replicas:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           Saver:ManagerIndexSaver{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Corrector:ManagerIndexCorrector{},
// 		           Creator:ManagerIndexCreator{},
// 		           Deleter:ManagerIndexDeleter{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Image:Image{},
// 		           Indexer:ManagerIndexIndexer{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxUnavailable:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           Operator:ManagerIndexOperator{},
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:ManagerIndexReadreplica{},
// 		           Replicas:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           Saver:ManagerIndexSaver{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			i := &ManagerIndex{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				Corrector:                     test.fields.Corrector,
// 				Creator:                       test.fields.Creator,
// 				Deleter:                       test.fields.Deleter,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				Image:                         test.fields.Image,
// 				Indexer:                       test.fields.Indexer,
// 				InitContainers:                test.fields.InitContainers,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				Operator:                      test.fields.Operator,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Readreplica:                   test.fields.Readreplica,
// 				Replicas:                      test.fields.Replicas,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				Saver:                         test.fields.Saver,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			i.SetResources()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestManagerIndex_SetTopologySpreadConstraints(t *testing.T) {
// 	type fields struct {
// 		Affinity                      *Affinity
// 		Annotations                   *map[string]any
// 		Corrector                     *ManagerIndexCorrector
// 		Creator                       *ManagerIndexCreator
// 		Deleter                       *ManagerIndexDeleter
// 		Enabled                       *bool
// 		Env                           *Env
// 		ExternalTrafficPolicy         *string
// 		Image                         *Image
// 		Indexer                       *ManagerIndexIndexer
// 		InitContainers                *InitContainers
// 		Kind                          *ManagerIndexKind
// 		Logging                       *Logging
// 		MaxUnavailable                *string
// 		Name                          *string
// 		NodeName                      *string
// 		NodeSelector                  *NodeSelector
// 		Observability                 *Observability
// 		Operator                      *ManagerIndexOperator
// 		PodAnnotations                *map[string]any
// 		PodPriority                   *PodPriority
// 		PodSecurityContext            *corev1.PodSecurityContext
// 		ProgressDeadlineSeconds       *int
// 		Readreplica                   *ManagerIndexReadreplica
// 		Replicas                      *int
// 		Resources                     *Resources
// 		RevisionHistoryLimit          *int
// 		RollingUpdate                 *RollingUpdate
// 		Saver                         *ManagerIndexSaver
// 		SecurityContext               *corev1.SecurityContext
// 		ServerConfig                  *ServerConfig
// 		Service                       *Service
// 		ServiceAccountName            *string
// 		ServiceType                   *ManagerIndexServiceType
// 		TerminationGracePeriodSeconds *int
// 		TimeZone                      *string
// 		Tolerations                   *Tolerations
// 		TopologySpreadConstraints     *TopologySpreadConstraints
// 		UnhealthyPodEvictionPolicy    *ManagerIndexUnhealthyPodEvictionPolicy
// 		Version                       *Version
// 		VolumeMounts                  *VolumeMounts
// 		Volumes                       *Volumes
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Corrector:ManagerIndexCorrector{},
// 		           Creator:ManagerIndexCreator{},
// 		           Deleter:ManagerIndexDeleter{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Image:Image{},
// 		           Indexer:ManagerIndexIndexer{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxUnavailable:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           Operator:ManagerIndexOperator{},
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:ManagerIndexReadreplica{},
// 		           Replicas:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           Saver:ManagerIndexSaver{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 		           Affinity:nil,
// 		           Annotations:nil,
// 		           Corrector:ManagerIndexCorrector{},
// 		           Creator:ManagerIndexCreator{},
// 		           Deleter:ManagerIndexDeleter{},
// 		           Enabled:nil,
// 		           Env:nil,
// 		           ExternalTrafficPolicy:nil,
// 		           Image:Image{},
// 		           Indexer:ManagerIndexIndexer{},
// 		           InitContainers:nil,
// 		           Kind:nil,
// 		           Logging:nil,
// 		           MaxUnavailable:nil,
// 		           Name:nil,
// 		           NodeName:nil,
// 		           NodeSelector:nil,
// 		           Observability:nil,
// 		           Operator:ManagerIndexOperator{},
// 		           PodAnnotations:nil,
// 		           PodPriority:PodPriority{},
// 		           PodSecurityContext:nil,
// 		           ProgressDeadlineSeconds:nil,
// 		           Readreplica:ManagerIndexReadreplica{},
// 		           Replicas:nil,
// 		           Resources:nil,
// 		           RevisionHistoryLimit:nil,
// 		           RollingUpdate:RollingUpdate{},
// 		           Saver:ManagerIndexSaver{},
// 		           SecurityContext:nil,
// 		           ServerConfig:ServerConfig{},
// 		           Service:Service{},
// 		           ServiceAccountName:nil,
// 		           ServiceType:nil,
// 		           TerminationGracePeriodSeconds:nil,
// 		           TimeZone:nil,
// 		           Tolerations:nil,
// 		           TopologySpreadConstraints:nil,
// 		           UnhealthyPodEvictionPolicy:nil,
// 		           Version:nil,
// 		           VolumeMounts:nil,
// 		           Volumes:nil,
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
// 			i := &ManagerIndex{
// 				Affinity:                      test.fields.Affinity,
// 				Annotations:                   test.fields.Annotations,
// 				Corrector:                     test.fields.Corrector,
// 				Creator:                       test.fields.Creator,
// 				Deleter:                       test.fields.Deleter,
// 				Enabled:                       test.fields.Enabled,
// 				Env:                           test.fields.Env,
// 				ExternalTrafficPolicy:         test.fields.ExternalTrafficPolicy,
// 				Image:                         test.fields.Image,
// 				Indexer:                       test.fields.Indexer,
// 				InitContainers:                test.fields.InitContainers,
// 				Kind:                          test.fields.Kind,
// 				Logging:                       test.fields.Logging,
// 				MaxUnavailable:                test.fields.MaxUnavailable,
// 				Name:                          test.fields.Name,
// 				NodeName:                      test.fields.NodeName,
// 				NodeSelector:                  test.fields.NodeSelector,
// 				Observability:                 test.fields.Observability,
// 				Operator:                      test.fields.Operator,
// 				PodAnnotations:                test.fields.PodAnnotations,
// 				PodPriority:                   test.fields.PodPriority,
// 				PodSecurityContext:            test.fields.PodSecurityContext,
// 				ProgressDeadlineSeconds:       test.fields.ProgressDeadlineSeconds,
// 				Readreplica:                   test.fields.Readreplica,
// 				Replicas:                      test.fields.Replicas,
// 				Resources:                     test.fields.Resources,
// 				RevisionHistoryLimit:          test.fields.RevisionHistoryLimit,
// 				RollingUpdate:                 test.fields.RollingUpdate,
// 				Saver:                         test.fields.Saver,
// 				SecurityContext:               test.fields.SecurityContext,
// 				ServerConfig:                  test.fields.ServerConfig,
// 				Service:                       test.fields.Service,
// 				ServiceAccountName:            test.fields.ServiceAccountName,
// 				ServiceType:                   test.fields.ServiceType,
// 				TerminationGracePeriodSeconds: test.fields.TerminationGracePeriodSeconds,
// 				TimeZone:                      test.fields.TimeZone,
// 				Tolerations:                   test.fields.Tolerations,
// 				TopologySpreadConstraints:     test.fields.TopologySpreadConstraints,
// 				UnhealthyPodEvictionPolicy:    test.fields.UnhealthyPodEvictionPolicy,
// 				Version:                       test.fields.Version,
// 				VolumeMounts:                  test.fields.VolumeMounts,
// 				Volumes:                       test.fields.Volumes,
// 			}
//
// 			i.SetTopologySpreadConstraints()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestValdRelease_SetScaledResources(t *testing.T) {
// 	type args struct {
// 		ar int
// 		am corev1.ResourceList
// 		p  ResourceParams
// 	}
// 	type fields struct {
// 		Base       resource.Base[ValdRelease, *ValdRelease]
// 		TypeMeta   metav1.TypeMeta
// 		ObjectMeta metav1.ObjectMeta
// 		Spec       Values
// 		Status     VrsStatus
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ar:0,
// 		           am:nil,
// 		           p:ResourceParams{},
// 		       },
// 		       fields: fields {
// 		           Base:nil,
// 		           TypeMeta:nil,
// 		           ObjectMeta:nil,
// 		           Spec:Values{},
// 		           Status:VrsStatus{},
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
// 		           ar:0,
// 		           am:nil,
// 		           p:ResourceParams{},
// 		           },
// 		           fields: fields {
// 		           Base:nil,
// 		           TypeMeta:nil,
// 		           ObjectMeta:nil,
// 		           Spec:Values{},
// 		           Status:VrsStatus{},
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
// 			v := &ValdRelease{
// 				Base:       test.fields.Base,
// 				TypeMeta:   test.fields.TypeMeta,
// 				ObjectMeta: test.fields.ObjectMeta,
// 				Spec:       test.fields.Spec,
// 				Status:     test.fields.Status,
// 			}
//
// 			v.SetScaledResources(test.args.ar, test.args.am, test.args.p)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
