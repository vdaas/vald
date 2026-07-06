package vald

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/agent"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/discoverer"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *VrsBuilder) buildAgent() agent.Agent {
	input := b.CR.Spec.VectorEngine.Vald.Agent

	mu := intstr.FromString(agent.DefaultAgentMaxUnavailable)
	ms := intstr.FromString(agent.DefaultAgentMaxSurge)

	a := agent.Agent{
		Logging: b.buildLogging(input.LogLevel),
		Kind:    common.KindTypeStatefulSet,
		RollingUpdate: &discoverer.RollingUpdateValdreelase{
			MaxUnavailable: &mu,
			MaxSurge:       &ms,
		},
		NGT: *b.buildAgentNgt(),
	}
	// TODO: Discuss with the OSS team whether alternative agent backends (e.g., QBG, Faiss)
	// should be supported alongside NGT.
	return a
}

func (b *VrsBuilder) buildAgentNgt() *agent.NGT {
	input := b.CR.Spec.VectorEngine.Vald.Agent.Ngt
	return &agent.NGT{
		Dimension:          input.Dimension,
		DistanceType:       input.DistanceType,
		ObjectType:         input.ObjectType,
		SearchEdgeSize:     input.SearchEdgeSize,
		CreationEdgeSize:   input.CreationEdgeSize,
		EnableInMemoryMode: true,
	}
}
