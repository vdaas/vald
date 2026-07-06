package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/agent"
	v1 "k8s.io/api/core/v1"
)

func TestLb_SetReplica(t *testing.T) {
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
			lb := &Lb{}
			a := agent.Agent{MaxReplicas: tt.agentMax}
			lb.SetReplica(a)
			assert.Equal(t, tt.wantMin, lb.MinReplicas)
			assert.Equal(t, tt.wantMax, lb.MaxReplicas)
		})
	}
}

func TestLb_SetResources(t *testing.T) {
	lb := &Lb{}
	lb.SetResources()

	assert.NotNil(t, lb.Resources)
	assert.NotEmpty(t, lb.Resources.Requests[v1.ResourceCPU])
	assert.NotEmpty(t, lb.Resources.Requests[v1.ResourceMemory])
	assert.NotEmpty(t, lb.Resources.Limits[v1.ResourceCPU])
	assert.NotEmpty(t, lb.Resources.Limits[v1.ResourceMemory])
}

func TestLb_SetTopologySpreadConstraints(t *testing.T) {
	lb := &Lb{}
	lb.SetTopologySpreadConstraints()

	assert.Len(t, lb.TopologySpreadConstraints, 1)
	assert.Equal(t, "gateway-lb", lb.TopologySpreadConstraints[0].LabelSelector.MatchLabels["app.kubernetes.io/component"])
}
