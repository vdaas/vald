package gateway

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/agent"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (l *Lb) getMaxReplica(ar int) int {
	if r := ar * 2; r >= 1 {
		return r
	}
	return 1
}

func (l *Lb) getMinReplica(ar int) int {
	if r := ar / 2; r >= 1 {
		return r
	}
	return 1
}

func (l *Lb) SetReplica(a agent.Agent) {
	ar := a.MaxReplicas
	l.MinReplicas = l.getMinReplica(ar)
	l.MaxReplicas = l.getMaxReplica(ar)
}

func (l *Lb) SetResources() {
	l.Resources = &v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("200m"),
			v1.ResourceMemory: resource.MustParse("150Mi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("2000m"),
			v1.ResourceMemory: resource.MustParse("700Mi"),
		},
	}
}

func (l *Lb) SetTopologySpreadConstraints() {
	l.TopologySpreadConstraints = []v1.TopologySpreadConstraint{
		common.BuildTopologySpreadConstraint(componentLabelGatewayLb),
	}
}
