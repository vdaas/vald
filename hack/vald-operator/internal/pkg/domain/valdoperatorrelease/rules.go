package valdoperatorrelease

import (
	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	builder "github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/builder/vald"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Compile-time check: Domain satisfies builder.DomainRules.
var _ builder.DomainRules = (*Domain)(nil)

// ResolveAgentNodePool returns the resolved node-pool inputs for the agent
// of a single Infrastructure entry, applying the general-pool fallback:
// when no dedicated agent pool exists (Replicas == 0), the spec falls
// back to the general pool.
func (d *Domain) ResolveAgentNodePool(infra v1.ValdOperatorReleaseInfra) builder.AgentNodePoolSpec {
	gn := infra.NodePools.GetNodePool(v1.NodePoolTypeGeneral)
	an := infra.NodePools.GetNodePool(v1.NodePoolTypeValdAgent)
	if an == nil || an.Replicas == 0 {
		return builder.AgentNodePoolSpec{
			NodeCount:       gn.Replicas,
			MachineResource: gn.MachineResource.GetResourceList(),
		}
	}
	return builder.AgentNodePoolSpec{
		NodeCount:       an.Replicas,
		MachineResource: an.MachineResource.GetResourceList(),
	}
}

// AgentPvSize returns the PV size for the agent, applying the buffer ratio
// and minimum floor:
//
//	size = max(memoryBytes * pvBufferRatio, pvMinSizeBytes)
//
// The result is rounded up to a whole Gi so the rendered Quantity is
// k8s-friendly ("15Gi" rather than raw bytes). Caller supplies the formula
// inputs (typically from *config.Config).
func (d *Domain) AgentPvSize(memoryBytes int64, pvBufferRatio float64, pvMinSizeBytes int64) string {
	const gi = int64(1) << 30
	size := max(int64(float64(memoryBytes)*pvBufferRatio), pvMinSizeBytes)
	sizeGi := (size + gi - 1) / gi
	return resource.NewQuantity(sizeGi*gi, resource.BinarySI).String()
}
