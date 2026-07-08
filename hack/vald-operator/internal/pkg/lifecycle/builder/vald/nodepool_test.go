package vald

import (
	"testing"

	"github.com/stretchr/testify/assert"
	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/manager"
)

func newAffinityBuilder(capability NodePoolCapability) *VrsBuilder {
	cr := &controllerv1.ValdOperatorRelease{}
	cr.SetNamespace("ns-affinity")
	return &VrsBuilder{
		CR:         cr,
		Config:     &config.Config{NodePoolLabelPrefix: "vald.vdaas.org"},
		Capability: capability,
	}
}

func TestVrsBuilder_ApplyNodeAffinities_AllComponents(t *testing.T) {
	b := newAffinityBuilder(AlwaysAvailable())
	row := &valdrelease.ValdRelease{}
	row.Spec.Manager.Index.Saver = &manager.Saver{}
	row.Spec.Manager.Index.Creator = &manager.Creator{}

	b.applyNodeAffinities(row)

	const (
		nsKey   = "vald.vdaas.org/namespace"
		typeKey = "vald.vdaas.org/type"
	)

	assertGeneral := func(name string, ns map[string]string, tols *[]any) {
		t.Helper()
		assert.Equal(t, "ns-affinity", ns[nsKey], name+" NodeSelector namespace")
		assert.Equal(t, "general", ns[typeKey], name+" NodeSelector type")
	}

	// Agent: agent pool (because Capability says HasAgentPool=true)
	assert.Equal(t, "agent", row.Spec.Agent.NodeSelector[typeKey], "agent uses agent pool")
	assert.NotNil(t, row.Spec.Agent.Tolerations, "agent tolerations set")

	// All others: general pool
	assertGeneral("gateway.Lb", row.Spec.Gateway.Lb.NodeSelector, nil)
	assertGeneral("discoverer", row.Spec.Discoverer.NodeSelector, nil)
	assertGeneral("manager.Index", row.Spec.Manager.Index.NodeSelector, nil)
	assertGeneral("manager.Index.Saver", row.Spec.Manager.Index.Saver.NodeSelector, nil)
	assertGeneral("manager.Index.Creator", row.Spec.Manager.Index.Creator.NodeSelector, nil)

	assert.NotNil(t, row.Spec.Gateway.Lb.Tolerations)
	assert.NotNil(t, row.Spec.Discoverer.Tolerations)
	assert.NotNil(t, row.Spec.Manager.Index.Tolerations)
	assert.NotNil(t, row.Spec.Manager.Index.Saver.Tolerations)
	assert.NotNil(t, row.Spec.Manager.Index.Creator.Tolerations)
}

func TestVrsBuilder_ApplyNodeAffinities_AgentFallsBackToGeneral_WhenNoAgentPool(t *testing.T) {
	// Capability declares the cluster has no dedicated agent pool.
	b := newAffinityBuilder(NodePoolCapability{HasGeneralPool: true, HasAgentPool: false})
	row := &valdrelease.ValdRelease{}

	b.applyNodeAffinities(row)

	assert.Equal(t, "general", row.Spec.Agent.NodeSelector["vald.vdaas.org/type"],
		"agent must fall back to general pool when no agent pool exists")
}

func TestVrsBuilder_ApplyNodeAffinities_NilSaverCreator(t *testing.T) {
	b := newAffinityBuilder(AlwaysAvailable())
	row := &valdrelease.ValdRelease{}
	// Saver and Creator are nil (Index.Enabled=true path uses Indexer instead).

	// Must not panic; non-nil components still get values.
	b.applyNodeAffinities(row)

	assert.Nil(t, row.Spec.Manager.Index.Saver)
	assert.Nil(t, row.Spec.Manager.Index.Creator)
	assert.NotNil(t, row.Spec.Agent.Tolerations)
}
