package manager

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (i *Index) SetResources() {
	i.Resources = &v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("200m"),
			v1.ResourceMemory: resource.MustParse("80Mi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1000m"),
			v1.ResourceMemory: resource.MustParse("500Mi"),
		},
	}
}

func (i *Index) SetTopologySpreadConstraints() {
	i.TopologySpreadConstraints = []v1.TopologySpreadConstraint{
		common.BuildTopologySpreadConstraint(componentLabelManagerIndex),
	}
}
