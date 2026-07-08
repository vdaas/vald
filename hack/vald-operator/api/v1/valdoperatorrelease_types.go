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

package v1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NodePoolTypeGeneral   NodePoolType = "general"
	NodePoolTypeValdAgent NodePoolType = "agent"
)

// ValdOperatorReleaseSpec defines the desired state of ValdOperatorRelease
type ValdOperatorReleaseSpec struct {
	Infrastructure []ValdOperatorReleaseInfra `json:"infrastructure" yaml:"infrastructure"`
	VectorEngine   VectorEngine        `json:"vectorEngine" yaml:"vectorEngine"`
}

// +kubebuilder:validation:Required
type ValdOperatorReleaseInfra struct {
	//+kubebuilder:validation:Required
	Role RoleType `json:"role" yaml:"role"`
	//+kubebuilder:validation:Required
	Type ClusterType `json:"type" yaml:"type"`
	// +kubebuilder:validation:Required
	Active bool `json:"active" yaml:"active"`
	// +kubebuilder:validation:Required
	Clusters []DestClusters `json:"clusters" yaml:"clusters"`
	// +kubebuilder:validation:Required
	NodePools NodePools `json:"nodePools" yaml:"nodePools"`
}

type RoleType string

type ClusterType string

type DestClusters struct {
	ID string `json:"id" yaml:"id"`
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// +kubebuilder:validation:Required
// +kubebuilder:validation:MinProperties=1
type NodePools map[NodePoolType]NodePool

// +kubebuilder:validation:Required
type NodePool struct {
	// +kubebuilder:validation:Required
	Name string `json:"name" yaml:"name"`
	// +kubebuilder:validation:Required
	MachineResource MachineResource `json:"machineResource" yaml:"machineResource"`
	// +kubebuilder:validation:Required
	Replicas int `json:"replicas" yaml:"replicas"`
}

type NodePoolType string

type MachineResource struct {
	Name string `json:"name" yaml:"name"`
	// +kubebuilder:validation:MinLength=1
	Cpu string `json:"cpu" yaml:"cpu"`
	// +kubebuilder:validation:MinLength=1
	Memory  string `json:"memory" yaml:"memory"`
	Storage string `json:"storage,omitempty" yaml:"storage,omitempty"`
}

type VectorEngine struct {
	Name string `json:"name" yaml:"name"`
	Vald Vald   `json:"vald" yaml:"vald"`
}

type Vald struct {
	Defaults   ValdDefaults         `json:"defaults" yaml:"defaults"`
	Agent      Agent                `json:"agent" yaml:"agent"`
	Indexer    Indexer              `json:"indexer" yaml:"indexer"`
	Gateway    Gateway              `json:"gateway" yaml:"gateway"`
	Discoverer Discoverer           `json:"discoverer" yaml:"discoverer"`
	Overlay    apiextensionsv1.JSON `json:"overlay,omitempty" yaml:"overlay,omitempty"`
}

type ValdDefaults struct {
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
}

type Agent struct {
	LogLevel         string                 `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	Ngt              Ngt                    `json:"ngt" yaml:"ngt,omitempty"`
	PersistentVolume *AgentPersistentVolume `json:"persistentVolume,omitempty" yaml:"persistentVolume,omitempty"`
}

type AgentPersistentVolume struct {
	Enabled      bool   `json:"enabled" yaml:"enabled"`
	StorageClass string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
	AccessMode   string `json:"accessMode,omitempty" yaml:"accessMode,omitempty"`
}
type Ngt struct {
	CreationEdgeSize int `json:"creationEdgeSize" yaml:"creationEdgeSize"`
	SearchEdgeSize   int `json:"searchEdgeSize" yaml:"searchEdgeSize"`
	// +kubebuilder:validation:Minimum=2
	Dimension    int    `json:"dimension" yaml:"dimension"`
	DistanceType string `json:"distanceType" yaml:"distanceType"`
	ObjectType   string `json:"objectType" yaml:"objectType"`
}

type Indexer struct {
	LogLevel      string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	IndexSchedule string `json:"indexSchedule" yaml:"indexSchedule,omitempty"`
	IndexSuspend  bool   `json:"indexSuspend" yaml:"indexSuspend"`
	SaveSuspend   bool   `json:"saveSuspend" yaml:"saveSuspend"`
	SaveSchedule  string `json:"saveSchedule" yaml:"saveSchedule,omitempty"`
	Concurrency   int    `json:"concurrency" yaml:"concurrency"`
	Manager       bool   `json:"manager" yaml:"manager"`
	IndexDuration string `json:"indexDuration" yaml:"indexDuration"`
	SaveDuration  string `json:"saveDuration" yaml:"saveDuration"`
}

type Gateway struct {
	LogLevel     string          `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	IndexReplica int             `json:"indexReplica" yaml:"indexReplica"`
	ServiceType  string          `json:"serviceType,omitempty" yaml:"serviceType,omitempty"`
	Ingress      *GatewayIngress `json:"ingress,omitempty" yaml:"ingress,omitempty"`
}

type GatewayIngress struct {
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Host    string `json:"host,omitempty" yaml:"host,omitempty"`
}

type Discoverer struct {
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	// +kubebuilder:validation:Enum=DaemonSet;Deployment
	Kind string `json:"kind" yaml:"kind"`
}

// ValdOperatorReleaseStatus defines the observed state of ValdOperatorRelease.
type ValdOperatorReleaseStatus struct {
	// Define observed state of cluster
	Conditions []metav1.Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Phase      string             `json:"phase,omitempty" yaml:"phase,omitempty"`
	Progress   Progress           `json:"progress,omitempty" yaml:"progress,omitempty"`
}

type AnalyzedLog struct {
	Date    metav1.Time `json:"date"`
	Details []string    `json:"details"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vor,scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Progress",type=integer,JSONPath=".status.progress.total"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=".status.conditions[0].status"
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=".metadata.creationTimestamp"
// ValdOperatorRelease is the Schema for the valdoperatorreleases API.
type ValdOperatorRelease struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec   ValdOperatorReleaseSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status ValdOperatorReleaseStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ValdOperatorReleaseList contains a list of ValdOperatorRelease.
type ValdOperatorReleaseList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Items           []ValdOperatorRelease `json:"items" yaml:"items"`
}

func init() {
	SchemeBuilder.Register(&ValdOperatorRelease{}, &ValdOperatorReleaseList{})
}
