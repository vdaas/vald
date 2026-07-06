package v1

import (
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
		v1.ResourceMemory: resource.MustParse(mr.Memory),
		v1.ResourceCPU:    resource.MustParse(mr.Cpu),
	}
	if mr.Storage != "" {
		rl[v1.ResourceStorage] = resource.MustParse(mr.Storage)
	}
	return rl
}
