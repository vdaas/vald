package vald

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAlwaysAvailable(t *testing.T) {
	c := AlwaysAvailable()
	assert.True(t, c.HasGeneralPool)
	assert.True(t, c.HasAgentPool)
}

func TestResolveNodePoolCapability_NilClient(t *testing.T) {
	_, err := ResolveNodePoolCapability(context.Background(), nil, "default", "")
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
					labelKey(prefix, nodePoolLabelType):      string(controllerv1.NodePoolTypeGeneral),
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
					labelKey(prefix, nodePoolLabelType):      string(controllerv1.NodePoolTypeGeneral),
				}),
				makeNode("n2", map[string]string{
					labelKey(prefix, nodePoolLabelNamespace): namespace,
					labelKey(prefix, nodePoolLabelType):      string(controllerv1.NodePoolTypeValdAgent),
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
					labelKey(prefix, nodePoolLabelType):      string(controllerv1.NodePoolTypeGeneral),
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
			got, err := ResolveNodePoolCapability(context.Background(), c, namespace, prefix)
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
