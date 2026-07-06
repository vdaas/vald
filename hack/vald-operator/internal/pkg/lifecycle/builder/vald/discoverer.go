package vald

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/discoverer"
)

func (b *VrsBuilder) buildDiscoverer() discoverer.Discoverer {
	input := b.CR.Spec.VectorEngine.Vald.Discoverer

	d := discoverer.Discoverer{
		Logging: b.buildLogging(input.LogLevel),
		ClusterRole: discoverer.ClusterRole{
			Name: b.CR.Namespace,
		},
		ClusterRoleBinding: discoverer.ClusterRoleBinding{
			Name: b.CR.Namespace,
		},
		Kind: common.KindType(input.Kind),
	}
	return d
}
