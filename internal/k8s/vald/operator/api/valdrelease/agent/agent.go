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

// Package agent defines the agent component model of the generated ValdRelease.
package agent

import (
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/defaults"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/discoverer"
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
	RollingUpdate             *discoverer.RollingUpdateValdRelease `json:"rollingUpdate,omitempty"`
	MinReplicas               int                                  `json:"minReplicas"`
	MaxReplicas               int                                  `json:"maxReplicas"`
	Resources                 *v1.ResourceRequirements             `json:"resources,omitempty"`
	TopologySpreadConstraints []v1.TopologySpreadConstraint        `json:"topologySpreadConstraints,omitempty"`
	NGT                       NGT                                  `json:"ngt"`
	PersistentVolume          *PersistentVolume                    `json:"persistentVolume,omitempty"`
}

type PersistentVolume struct {
	Enabled      bool   `yaml:"enabled"      json:"enabled"`
	AccessMode   string `yaml:"accessMode"   json:"accessMode"`
	StorageClass string `yaml:"storageClass" json:"storageClass"`
	Size         string `yaml:"size"         json:"size"`
}

type NGT struct {
	Dimension          int      `json:"dimension"`
	DistanceType       string   `json:"distance_type"`
	ObjectType         string   `json:"object_type"`
	CreationEdgeSize   int      `json:"creation_edge_size"`
	SearchEdgeSize     int      `json:"search_edge_size"`
	DefaultRadius      *float32 `json:"default_radius,omitempty"`
	DefaultEpsilon     *float32 `json:"default_epsilon,omitempty"`
	IndexPath          string   `json:"index_path"                     yaml:"index_path"`
	EnableInMemoryMode bool     `json:"enable_in_memory_mode"          yaml:"enable_in_memory_mode"`
	EnableCopyOnWrite  bool     `json:"enable_copy_on_write,omitempty"`
}
