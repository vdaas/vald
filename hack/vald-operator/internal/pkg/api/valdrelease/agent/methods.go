package agent

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"
	v1 "k8s.io/api/core/v1"
)

const componentLabelAgent = "agent"

// SetPvEnable configures the agent for persistent storage. The caller resolves
// StorageClass, AccessMode, and the rendered Size (typically via a domain
// rule) and passes them in — keeping policy out of this method preserves
// testability and keeps the api package free of domain concerns.
func (a *Agent) SetPvEnable(sc, am, size string) {
	// Override the NGT settings required for PV-backed operation.
	a.NGT.EnableCopyOnWrite = true
	a.NGT.EnableInMemoryMode = false
	a.NGT.IndexPath = DefaultIndexPath

	a.PersistentVolume = &PersistentVolume{
		Enabled:      true,
		StorageClass: sc,
		Size:         size,
		AccessMode:   am,
	}
}

func (a *Agent) SetReplica(nr int, podsPerNode int) {
	rep := nr * podsPerNode
	a.MinReplicas = rep
	a.MaxReplicas = rep
}

func (a *Agent) SetResources(mc v1.ResourceList, podsPerNode int) {
	div := float64(podsPerNode)

	cpu := mc.Cpu().Value()
	memory := mc.Memory().Value()

	temp := &v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    common.CalcResource(cpu, ResourceRatio, div),
			v1.ResourceMemory: common.CalcResource(memory, ResourceRatio, div),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU: common.CalcResource(cpu, ResourceRatio),
			// Memory limit is intentionally omitted: the NGT index grows after startup,
			// and a hard limit would cause OOM kills as the index size increases.
		},
	}

	a.Resources = &v1.ResourceRequirements{
		Requests: common.NormalizeResourceList(temp.Requests),
		Limits:   common.NormalizeResourceList(temp.Limits),
	}
}

func (a *Agent) SetTopologySpreadConstraints() {
	a.TopologySpreadConstraints = []v1.TopologySpreadConstraint{
		common.BuildTopologySpreadConstraint(componentLabelAgent),
	}
}
