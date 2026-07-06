package valdrelease

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/agent"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/discoverer"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/gateway"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/manager"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func newTestValdRelease() *ValdRelease {
	return &ValdRelease{
		Spec: Spec{
			Agent:      agent.Agent{},
			Gateway:    gateway.Gateway{},
			Discoverer: discoverer.Discoverer{Kind: "Deployment"},
			Manager: manager.Manager{
				Index: manager.Index{
					Enabled: true,
					Creator: &manager.Creator{Enabled: true},
					Saver:   &manager.Saver{Enabled: true},
				},
			},
		},
	}
}

func TestSetRelationalResources(t *testing.T) {
	vr := newTestValdRelease()
	mc := v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse("8"),
		v1.ResourceMemory: resource.MustParse("16Gi"),
	}

	vr.SetRelationalResources(3, mc, ResourceParams{
		AgentPodsPerNode:           2,
		DiscovererDSMaxSurge:       "30%",
		DiscovererDSMaxUnavailable: "0%",
	})

	assert.NotZero(t, vr.Spec.Agent.MinReplicas)
	assert.NotZero(t, vr.Spec.Agent.MaxReplicas)
	assert.NotNil(t, vr.Spec.Agent.Resources)
	assert.NotZero(t, vr.Spec.Gateway.Lb.MinReplicas)
	assert.NotZero(t, vr.Spec.Gateway.Lb.MaxReplicas)
	assert.NotNil(t, vr.Spec.Gateway.Lb.Resources)
	assert.NotNil(t, vr.Spec.Discoverer.Resources)
	assert.NotNil(t, vr.Spec.Manager.Index.Resources)
}
