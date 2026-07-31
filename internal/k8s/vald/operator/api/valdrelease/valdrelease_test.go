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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orig",
			Namespace: "ns",
			Labels:    map[string]string{"k": "v"},
		},
		Spec: Values{Agent: &Agent{MinReplicas: new(3)}},
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
		TypeMeta:   metav1.TypeMeta{APIVersion: GroupVersion.String(), Kind: GVK.Kind},
		ObjectMeta: metav1.ObjectMeta{Name: "rt", Namespace: "ns"},
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
