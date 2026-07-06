package vald

import (
	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	v1 "k8s.io/api/core/v1"
)

// AgentNodePoolSpec describes the resolved node-pool inputs for the agent
// of a single Infrastructure entry. The "resolved" part is the
// general-pool fallback: when no dedicated agent pool exists, the spec
// reflects the general pool instead.
type AgentNodePoolSpec struct {
	NodeCount       int
	MachineResource v1.ResourceList
}

// DomainRules captures the domain-owned policy decisions the Builder needs
// at Build time. The Builder does not import the domain package directly
// (which would create a cycle); instead the domain type satisfies this
// interface and is injected at construction time.
type DomainRules interface {
	ResolveAgentNodePool(infra controllerv1.MvaldreleaseInfra) AgentNodePoolSpec
	AgentPvSize(memoryBytes int64, pvBufferRatio float64, pvMinSizeBytes int64) string
}
