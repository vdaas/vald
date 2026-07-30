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
	"github.com/vdaas/vald/internal/k8s/resource"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NodePoolTypeGeneral   NodePoolType = "general"
	NodePoolTypeValdAgent NodePoolType = "agent"
)

// JSON is the apiextensions JSON type used by the Overlay field.
type JSON = apiextensionsv1.JSON

type ValdOperatorReleaseSpec struct {
	Infrastructure []ValdOperatorReleaseInfra `json:"infrastructure" yaml:"infrastructure"`
	VectorEngine   VectorEngine               `json:"vectorEngine"   yaml:"vectorEngine"`
}

type ValdOperatorReleaseInfra struct {
	Role      RoleType       `json:"role"      yaml:"role"`
	Type      ClusterType    `json:"type"      yaml:"type"`
	Active    bool           `json:"active"    yaml:"active"`
	Clusters  []DestClusters `json:"clusters"  yaml:"clusters"`
	NodePools NodePools      `json:"nodePools" yaml:"nodePools"`
}

type RoleType string

type ClusterType string

type DestClusters struct {
	ID   string `json:"id"   yaml:"id"`
	Name string `json:"name"`
}

type NodePools map[NodePoolType]NodePool

type NodePool struct {
	Name            string          `json:"name"            yaml:"name"`
	MachineResource MachineResource `json:"machineResource" yaml:"machineResource"`
	Replicas        int             `json:"replicas"        yaml:"replicas"`
}

type NodePoolType string

type MachineResource struct {
	Name    string `json:"name"              yaml:"name"`
	Cpu     string `json:"cpu"               yaml:"cpu"`
	Memory  string `json:"memory"            yaml:"memory"`
	Storage string `json:"storage,omitempty" yaml:"storage,omitempty"`
}

type VectorEngine struct {
	Name string `json:"name" yaml:"name"`
	Vald Vald   `json:"vald" yaml:"vald"`
}

type Vald struct {
	Defaults   ValdDefaults         `json:"defaults"   yaml:"defaults"`
	Agent      Agent                `json:"agent"      yaml:"agent"`
	Indexer    Indexer              `json:"indexer"    yaml:"indexer"`
	Gateway    Gateway              `json:"gateway"    yaml:"gateway"`
	Discoverer Discoverer           `json:"discoverer" yaml:"discoverer"`
	Overlay    apiextensionsv1.JSON `json:"overlay"    yaml:"overlay,omitempty"`
}

type ValdDefaults struct {
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
}

type Agent struct {
	LogLevel         string                 `json:"logLevel,omitempty"         yaml:"logLevel,omitempty"`
	Ngt              Ngt                    `json:"ngt"                        yaml:"ngt,omitempty"`
	PersistentVolume *AgentPersistentVolume `json:"persistentVolume,omitempty" yaml:"persistentVolume,omitempty"`
}

type AgentPersistentVolume struct {
	Enabled      bool   `json:"enabled"                yaml:"enabled"`
	StorageClass string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
	AccessMode   string `json:"accessMode,omitempty"   yaml:"accessMode,omitempty"`
}
type Ngt struct {
	CreationEdgeSize int    `json:"creationEdgeSize" yaml:"creationEdgeSize"`
	SearchEdgeSize   int    `json:"searchEdgeSize"   yaml:"searchEdgeSize"`
	Dimension        int    `json:"dimension"        yaml:"dimension"`
	DistanceType     string `json:"distanceType"     yaml:"distanceType"`
	ObjectType       string `json:"objectType"       yaml:"objectType"`
}

type Indexer struct {
	LogLevel      string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	IndexSchedule string `json:"indexSchedule"      yaml:"indexSchedule,omitempty"`
	IndexSuspend  bool   `json:"indexSuspend"       yaml:"indexSuspend"`
	SaveSuspend   bool   `json:"saveSuspend"        yaml:"saveSuspend"`
	SaveSchedule  string `json:"saveSchedule"       yaml:"saveSchedule,omitempty"`
	Concurrency   int    `json:"concurrency"        yaml:"concurrency"`
	Manager       bool   `json:"manager"            yaml:"manager"`
	IndexDuration string `json:"indexDuration"      yaml:"indexDuration"`
	SaveDuration  string `json:"saveDuration"       yaml:"saveDuration"`
}

type Gateway struct {
	LogLevel     string          `json:"logLevel,omitempty"    yaml:"logLevel,omitempty"`
	IndexReplica int             `json:"indexReplica"          yaml:"indexReplica"`
	ServiceType  string          `json:"serviceType,omitempty" yaml:"serviceType,omitempty"`
	Ingress      *GatewayIngress `json:"ingress,omitempty"     yaml:"ingress,omitempty"`
}

type GatewayIngress struct {
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Host    string `json:"host,omitempty"    yaml:"host,omitempty"`
}

type Discoverer struct {
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	Kind     string `json:"kind"               yaml:"kind"`
}

type ValdOperatorReleaseStatus struct {
	resource.Base[ValdOperatorReleaseStatus, *ValdOperatorReleaseStatus] `json:"-" yaml:"-"`

	Conditions []metav1.Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Phase      string             `json:"phase,omitempty"      yaml:"phase,omitempty"`
	Progress   Progress           `json:"progress"             yaml:"progress,omitempty"`
}

type ValdOperatorRelease struct {
	resource.Base[ValdOperatorRelease, *ValdOperatorRelease] `json:"-" yaml:"-"`

	metav1.TypeMeta   `json:",inline"  yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata,omitempty"`

	Spec   ValdOperatorReleaseSpec   `json:"spec"   yaml:"spec,omitempty"`
	Status ValdOperatorReleaseStatus `json:"status" yaml:"status,omitempty"`
}

type ValdOperatorReleaseList = resource.List[ValdOperatorRelease, *ValdOperatorRelease]
