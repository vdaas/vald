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

package valdrelease

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestAddToScheme verifies that AddToScheme registers both the item and the
// generic list kind, and that the scheme can construct and recognize the typed
// object for its GVK. Generic instantiations carry mangled reflect type names,
// so the list kind must be registered explicitly (see SchemeBuilder).
func TestAddToScheme(t *testing.T) {
	s := runtime.NewScheme()
	assert.NoError(t, AddToScheme(s))

	// scheme.New must construct a concrete *ValdRelease for the ValdRelease GVK.
	obj, err := s.New(GVK)
	assert.NoError(t, err)
	_, ok := obj.(*ValdRelease)
	assert.True(t, ok, "scheme must construct a *ValdRelease for the ValdRelease GVK")

	// The typed object must resolve back to its GVK.
	gvks, _, err := s.ObjectKinds(&ValdRelease{})
	assert.NoError(t, err)
	assert.Contains(t, gvks, GVK)

	// The list kind must be registered explicitly as ValdReleaseList.
	listGVK := GroupVersion.WithKind("ValdReleaseList")
	listObj, err := s.New(listGVK)
	assert.NoError(t, err, "the list kind must be registered under its explicit name")
	assert.NotNil(t, listObj)
}

// TestValdRelease_RuntimeObject proves ValdRelease satisfies runtime.Object and
// that DeepCopyObject returns a fully independent typed copy: mutating the
// original after the copy must not be observable through the copy.
func TestValdRelease_RuntimeObject(t *testing.T) {
	var _ runtime.Object = &ValdRelease{}

	vr := &ValdRelease{
		Name:      "orig",
		Namespace: "ns",
		Labels:    map[string]string{"k": "v"},
		Spec:      Values{Agent: &Agent{MinReplicas: new(3)}},
	}

	cp, ok := vr.DeepCopyObject().(*ValdRelease)
	assert.True(t, ok)
	assert.Equal(t, "orig", cp.GetName())
	assert.Equal(t, 3, *cp.Spec.Agent.MinReplicas)

	// Mutate every level of the original; the deep copy must stay unchanged.
	vr.SetName("changed")
	vr.GetLabels()["k"] = "mutated"
	*vr.Spec.Agent.MinReplicas = 99

	assert.Equal(t, "orig", cp.GetName(), "metadata must be independent")
	assert.Equal(t, "v", cp.GetLabels()["k"], "label map must be independent")
	assert.Equal(t, 3, *cp.Spec.Agent.MinReplicas, "spec pointers must be independent")
}

// TestValdRelease_JSONRoundTrip proves a populated spec survives a
// marshal→unmarshal cycle with its pointer and enum fields intact.
func TestValdRelease_JSONRoundTrip(t *testing.T) {
	vr := &ValdRelease{
		APIVersion: GroupVersion.String(), Kind: GVK.Kind,
		Name: "rt", Namespace: "ns",
		Spec: Values{
			Agent: &Agent{
				MinReplicas: new(2),
				Ngt: &AgentNgt{
					Dimension:    new(128),
					DistanceType: new(AgentNgtDistanceType("l2")),
				},
			},
			Gateway: &Gateway{Lb: &GatewayLb{MaxReplicas: new(4)}},
		},
	}

	raw, err := json.Marshal(vr)
	assert.NoError(t, err)

	var got ValdRelease
	assert.NoError(t, json.Unmarshal(raw, &got))

	assert.Equal(t, "rt", got.GetName())
	assert.Equal(t, GVK.Kind, got.Kind)
	assert.Equal(t, 2, *got.Spec.Agent.MinReplicas)
	assert.Equal(t, 128, *got.Spec.Agent.Ngt.Dimension)
	assert.Equal(t, AgentNgtDistanceType("l2"), *got.Spec.Agent.Ngt.DistanceType)
	assert.Equal(t, 4, *got.Spec.Gateway.Lb.MaxReplicas)
}

// securityContexts returns fresh, non-nil pod/container security contexts.
func securityContexts() (*corev1.PodSecurityContext, *corev1.SecurityContext) {
	return &corev1.PodSecurityContext{RunAsUser: new(int64(1000))},
		&corev1.SecurityContext{RunAsNonRoot: new(true)}
}

// TestDeepCopyInto_SecurityContexts guards that DeepCopyInto does not share the
// corev1 security-context pointers (nor their nested pointers) with the source.
func TestDeepCopyInto_SecurityContexts(t *testing.T) {
	t.Parallel()
	type pair struct {
		pod *corev1.PodSecurityContext
		sc  *corev1.SecurityContext
	}
	tests := []struct {
		// copyFn builds a source with both contexts set, deep-copies it, and
		// returns the source and copy pointers.
		copyFn func() (src, dst pair)
		name   string
	}{
		{name: "*Agent", copyFn: func() (pair, pair) {
			in, out := new(Agent), new(Agent)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*Discoverer", copyFn: func() (pair, pair) {
			in, out := new(Discoverer), new(Discoverer)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*GatewayFilter", copyFn: func() (pair, pair) {
			in, out := new(GatewayFilter), new(GatewayFilter)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*GatewayLb", copyFn: func() (pair, pair) {
			in, out := new(GatewayLb), new(GatewayLb)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*GatewayMirror", copyFn: func() (pair, pair) {
			in, out := new(GatewayMirror), new(GatewayMirror)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*ManagerIndex", copyFn: func() (pair, pair) {
			in, out := new(ManagerIndex), new(ManagerIndex)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*ManagerIndexOperator", copyFn: func() (pair, pair) {
			in, out := new(ManagerIndexOperator), new(ManagerIndexOperator)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*ManagerIndexReadreplicaRotator", copyFn: func() (pair, pair) {
			in, out := new(ManagerIndexReadreplicaRotator), new(ManagerIndexReadreplicaRotator)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			src, dst := tc.copyFn()
			if dst.pod == nil || dst.sc == nil {
				t.Fatalf("copy dropped security contexts: pod=%v sc=%v", dst.pod, dst.sc)
			}
			if dst.pod == src.pod || dst.pod.RunAsUser == src.pod.RunAsUser {
				t.Error("PodSecurityContext is shared with the source")
			}
			if dst.sc == src.sc || dst.sc.RunAsNonRoot == src.sc.RunAsNonRoot {
				t.Error("SecurityContext is shared with the source")
			}
			if *dst.pod.RunAsUser != 1000 || !*dst.sc.RunAsNonRoot {
				t.Errorf("copied values differ: RunAsUser=%d RunAsNonRoot=%t", *dst.pod.RunAsUser, *dst.sc.RunAsNonRoot)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestVrsStatus_DeepCopyInto(t *testing.T) {
// 	type args struct {
// 		out *VrsStatus
// 	}
// 	type fields struct {
// 		Status    Status
// 		Condition metav1.Condition
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
// 		           out:VrsStatus{},
// 		       },
// 		       fields: fields {
// 		           Status:nil,
// 		           Condition:nil,
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
// 		           out:VrsStatus{},
// 		           },
// 		           fields: fields {
// 		           Status:nil,
// 		           Condition:nil,
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
// 			in := &VrsStatus{
// 				Status:    test.fields.Status,
// 				Condition: test.fields.Condition,
// 			}
//
// 			in.DeepCopyInto(test.args.out)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestValdRelease_DeepCopyInto(t *testing.T) {
// 	type args struct {
// 		out *ValdRelease
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
// 		           out:ValdRelease{},
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
// 		           out:ValdRelease{},
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
// 			in := &ValdRelease{
// 				Base:       test.fields.Base,
// 				TypeMeta:   test.fields.TypeMeta,
// 				ObjectMeta: test.fields.ObjectMeta,
// 				Spec:       test.fields.Spec,
// 				Status:     test.fields.Status,
// 			}
//
// 			in.DeepCopyInto(test.args.out)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// TestGrpcClient_CallOption guards that the call_option values (defaults.grpc.client.call_option.* in
// charts/vald/values.yaml) are typed under call_option, survive the typed round trip the
// vald-operator makes, and that DeepCopyInto does not share their pointers.
func TestGrpcClient_CallOption(t *testing.T) {
	t.Parallel()
	var in GrpcClient
	assert.NoError(t, json.Unmarshal([]byte(`{"call_option":{"wait_for_ready":true,"max_recv_msg_size":4096,"content_subtype":"proto"}}`), &in))
	if !assert.NotNil(t, in.CallOption, "call_option was dropped while unmarshalling") {
		return
	}
	if assert.NotNil(t, in.CallOption.WaitForReady) && assert.NotNil(t, in.CallOption.MaxRecvMsgSize) && assert.NotNil(t, in.CallOption.ContentSubtype) {
		assert.True(t, *in.CallOption.WaitForReady)
		assert.Equal(t, 4096, *in.CallOption.MaxRecvMsgSize)
		assert.Equal(t, "proto", *in.CallOption.ContentSubtype)
	}

	raw, err := json.Marshal(in)
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"call_option":{`)
	assert.Contains(t, string(raw), `"wait_for_ready":true`)

	var out GrpcClient
	in.DeepCopyInto(&out)
	if assert.NotNil(t, out.CallOption, "DeepCopyInto dropped call_option") {
		assert.NotSame(t, in.CallOption, out.CallOption)
		if assert.NotNil(t, out.CallOption.WaitForReady) {
			assert.True(t, *out.CallOption.WaitForReady)
			assert.NotSame(t, in.CallOption.WaitForReady, out.CallOption.WaitForReady)
		}
	}
}

// TestImage_Registry guards that image.registry (gateway.lb.image.registry in charts/vald/values.yaml,
// read by the chart templates) survives the typed round trip the vald-operator makes and that
// DeepCopyInto does not share the pointer.
func TestImage_Registry(t *testing.T) {
	t.Parallel()
	var in Image
	assert.NoError(t, json.Unmarshal([]byte(`{"registry":"registry.example","repository":"vdaas/vald-lb-gateway"}`), &in))
	if assert.NotNil(t, in.Registry, "registry was dropped while unmarshalling") {
		assert.Equal(t, "registry.example", *in.Registry)
	}

	raw, err := json.Marshal(in)
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"registry":"registry.example"`)

	var out Image
	in.DeepCopyInto(&out)
	if assert.NotNil(t, out.Registry, "DeepCopyInto dropped registry") {
		assert.Equal(t, "registry.example", *out.Registry)
		assert.NotSame(t, in.Registry, out.Registry)
	}
}

// TestGatewayLbGatewayConfig_SearchQuorumGrace guards that gateway.lb.gateway_config.search_quorum_grace
// survives the typed round trip the vald-operator makes (a map of Helm values is marshalled and
// unmarshalled into ValdRelease) and that DeepCopyInto does not share the pointer. A field that is in
// charts/vald/values.schema.json but missing from the generated type is dropped silently there.
func TestGatewayLbGatewayConfig_SearchQuorumGrace(t *testing.T) {
	t.Parallel()
	var in GatewayLbGatewayConfig
	assert.NoError(t, json.Unmarshal([]byte(`{"search_quorum_grace":"10ms","multi_operation_concurrency":0}`), &in))
	if assert.NotNil(t, in.SearchQuorumGrace, "search_quorum_grace was dropped while unmarshalling") {
		assert.Equal(t, "10ms", *in.SearchQuorumGrace)
	}

	raw, err := json.Marshal(in)
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"search_quorum_grace":"10ms"`)

	var out GatewayLbGatewayConfig
	in.DeepCopyInto(&out)
	if assert.NotNil(t, out.SearchQuorumGrace, "DeepCopyInto dropped search_quorum_grace") {
		assert.Equal(t, "10ms", *out.SearchQuorumGrace)
		assert.NotSame(t, in.SearchQuorumGrace, out.SearchQuorumGrace)
	}
}
