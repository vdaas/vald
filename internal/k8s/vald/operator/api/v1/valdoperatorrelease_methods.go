// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package v1

import (
	"github.com/vdaas/vald/internal/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (np NodePools) GetNodePool(nt NodePoolType) *NodePool {
	if len(np) == 0 {
		return nil
	}
	node, ok := np[nt]
	if !ok {
		return nil
	}
	return &node
}

func (mr *MachineResource) GetResourceList() v1.ResourceList {
	rl := v1.ResourceList{
		v1.ResourceMemory: parseQuantity(mr.Memory),
		v1.ResourceCPU:    parseQuantity(mr.Cpu),
	}
	if mr.Storage != "" {
		rl[v1.ResourceStorage] = parseQuantity(mr.Storage)
	}
	return rl
}

// parseQuantity parses a user-supplied quantity string from the CRD.
// Invalid values fall back to the zero Quantity instead of panicking,
// since CRD fields are external input and must not crash the controller.
func parseQuantity(s string) resource.Quantity {
	q, err := resource.ParseQuantity(s)
	if err != nil {
		log.Warnf("failed to parse resource quantity %q from ValdOperatorRelease, falling back to zero Quantity: %v", s, err)
		return resource.Quantity{}
	}
	return q
}
