package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestIndex_SetResources(t *testing.T) {
	idx := &Index{}
	idx.SetResources()

	assert.NotNil(t, idx.Resources)
	assert.NotEmpty(t, idx.Resources.Requests[v1.ResourceCPU])
	assert.NotEmpty(t, idx.Resources.Requests[v1.ResourceMemory])
	assert.NotEmpty(t, idx.Resources.Limits[v1.ResourceCPU])
	assert.NotEmpty(t, idx.Resources.Limits[v1.ResourceMemory])
}

func TestIndex_SetTopologySpreadConstraints(t *testing.T) {
	idx := &Index{}
	idx.SetTopologySpreadConstraints()

	assert.Len(t, idx.TopologySpreadConstraints, 1)
	assert.Equal(t, "manager-index", idx.TopologySpreadConstraints[0].LabelSelector.MatchLabels["app.kubernetes.io/component"])
}
