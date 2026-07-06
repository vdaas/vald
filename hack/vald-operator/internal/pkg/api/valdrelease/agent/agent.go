// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package agent

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/discoverer"

	v1 "k8s.io/api/core/v1"
)

const (
	DefaultIndexPath           = "/var/ngt/index"
	DefaultAgentMaxSurge       = "1"
	DefaultAgentMaxUnavailable = "1"

	ResourceRatio = 0.6 // fraction of node resources allocated to agent pods on this node
)

type Agent struct {
	Logging                   *defaults.Logging                    `json:"logging,omitempty"`
	Affinity                  *v1.Affinity                         `json:"affinity,omitempty"`
	NodeSelector              map[string]string                    `json:"nodeSelector,omitempty"`
	Tolerations               *[]v1.Toleration                     `json:"tolerations,omitempty"`
	Kind                      common.KindType                      `json:"kind"`
	RollingUpdate             *discoverer.RollingUpdateValdreelase `json:"rollingUpdate,omitempty"`
	MinReplicas               int                                  `json:"minReplicas"`
	MaxReplicas               int                                  `json:"maxReplicas"`
	Resources                 *v1.ResourceRequirements             `json:"resources,omitempty"`
	TopologySpreadConstraints []v1.TopologySpreadConstraint        `json:"topologySpreadConstraints,omitempty"`
	NGT                       NGT                                  `json:"ngt"`
	PersistentVolume          *PersistentVolume                    `json:"persistentVolume,omitempty"`
}

type PersistentVolume struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	AccessMode   string `yaml:"accessMode" json:"accessMode"`
	StorageClass string `yaml:"storageClass" json:"storageClass"`
	Size         string `yaml:"size" json:"size"`
}

type NGT struct {
	Dimension          int    `json:"dimension"`
	DistanceType       string `json:"distance_type"`
	ObjectType         string `json:"object_type"`
	CreationEdgeSize   int    `json:"creation_edge_size"`
	SearchEdgeSize     int    `json:"search_edge_size"`
	DefaultRadius      *int   `json:"default_radius,omitempty"`
	DefaultEpsilon     *int   `json:"default_epsilon,omitempty"`
	IndexPath          string `yaml:"index_path" json:"index_path"`
	EnableInMemoryMode bool   `yaml:"enable_in_memory_mode" json:"enable_in_memory_mode"`
	EnableCopyOnWrite  bool   `json:"enable_copy_on_write,omitempty"`
}
