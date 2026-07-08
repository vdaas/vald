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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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
		nodes          []runtime.Object
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
			nodes: []runtime.Object{
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
			nodes: []runtime.Object{
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
			nodes: []runtime.Object{
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
			nodes: []runtime.Object{
				makeNode("n1", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): namespace,
					labelKey(prefix, nodePoolLabelType):      "gpu",
				}),
			},
			wantHasGeneral: false,
			wantHasAgent:   false,
		},
	}

	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(tt.nodes...).Build()
			got, err := resolveNodePoolCapability(context.Background(), c, namespace, prefix)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantHasGeneral, got.HasGeneralPool, "HasGeneralPool")
			assert.Equal(t, tt.wantHasAgent, got.HasAgentPool, "HasAgentPool")
		})
	}
}

func makeNode(name string, labels map[string]string) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
}

// NOT IMPLEMENTED BELOW
