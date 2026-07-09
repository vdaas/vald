//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package service

import (
	"github.com/vdaas/vald/internal/k8s"
	resource "github.com/vdaas/vald/internal/k8s/resource"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
)

// agentNodePoolSpec describes the resolved node-pool inputs for the agent of
// a single Infrastructure entry. The "resolved" part is the general-pool
// fallback: when no dedicated agent pool exists, the spec reflects the
// general pool instead.
type agentNodePoolSpec struct {
	MachineResource k8s.ResourceList
	NodeCount       int
}

// resolveAgentNodePool returns the resolved node-pool inputs for the agent
// of a single Infrastructure entry, applying the general-pool fallback:
// when no dedicated agent pool exists (Replicas == 0), the spec falls
// back to the general pool.
func resolveAgentNodePool(infra v1.ValdOperatorReleaseInfra) agentNodePoolSpec {
	gn := infra.NodePools.GetNodePool(v1.NodePoolTypeGeneral)
	an := infra.NodePools.GetNodePool(v1.NodePoolTypeValdAgent)
	if an == nil || an.Replicas == 0 {
		if gn == nil {
			// Neither a dedicated agent pool nor a general pool to fall back
			// to: the rule cannot be satisfied, so return an empty spec
			// (NodeCount 0 -> no replicas) instead of panicking on
			// user-supplied input.
			return agentNodePoolSpec{}
		}
		return agentNodePoolSpec{
			NodeCount:       gn.Replicas,
			MachineResource: gn.MachineResource.GetResourceList(),
		}
	}
	return agentNodePoolSpec{
		NodeCount:       an.Replicas,
		MachineResource: an.MachineResource.GetResourceList(),
	}
}

// agentPvSize returns the PV size for the agent, applying the buffer ratio
// and minimum floor:
//
//	size = max(memoryBytes * pvBufferRatio, pvMinSizeBytes)
//
// The result is rounded up to a whole Gi so the rendered Quantity is
// k8s-friendly ("15Gi" rather than raw bytes). Caller supplies the formula
// inputs (typically from *config.Config).
func agentPvSize(memoryBytes int64, pvBufferRatio float64, pvMinSizeBytes int64) string {
	const gi = int64(1) << 30
	size := max(int64(float64(memoryBytes)*pvBufferRatio), pvMinSizeBytes)
	sizeGi := (size + gi - 1) / gi
	return resource.NewQuantity(sizeGi*gi, resource.BinarySI).String()
}
