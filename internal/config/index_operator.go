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
package config

import "github.com/vdaas/vald/internal/k8s"

// IndexOperator represents the configurations for index k8s operator.
type IndexOperator struct {
	// Namespace represent the namespace of this pod
	Namespace string `json:"namespace" yaml:"namespace"`

	// AgentName represent agents meta_name for service discovery
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// AgentNamespace represent agent namespace location
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`

	// RotatorName represent rotator name for service discovery
	RotatorName string `json:"rotator_name" yaml:"rotator_name"`

	// TargetReadReplicaIDAnnotationsKey represents the environment variable name for target read replica id.
	TargetReadReplicaIDAnnotationsKey string `json:"target_read_replica_id_annotations_key" yaml:"target_read_replica_id_annotations_key"`

	// RotationJobConcurrency represents indexing concurrency.
	RotationJobConcurrency uint `json:"rotation_job_concurrency" yaml:"rotation_job_concurrency"`

	// ReadReplicaEnabled represents whether read replica is enabled or not.
	ReadReplicaEnabled bool `json:"read_replica_enabled" yaml:"read_replica_enabled"`

	// ReadReplicaLabelKey represents the label key for read replica.
	ReadReplicaLabelKey string `json:"read_replica_label_key" yaml:"read_replica_label_key"`

	// JobTemplates represents the job templates for indexing.
	JobTemplates IndexJobTemplates `json:"job_templates" yaml:"job_templates"`
}

type IndexJobTemplates struct {
	Rotate     *k8s.Job `json:"rotate"     yaml:"rotate"`
	Creation   *k8s.Job `json:"creation"   yaml:"creation"`
	Save       *k8s.Job `json:"save"       yaml:"save"`
	Correction *k8s.Job `json:"correction" yaml:"correction"`
}

// Bind binds the actual data from the IndexJobTemplates receiver fields.
func (ijt *IndexJobTemplates) Bind() *IndexJobTemplates {
	return ijt
}

func (ic *IndexOperator) Bind() *IndexOperator {
	ic.Namespace = GetActualValue(ic.Namespace)
	ic.AgentName = GetActualValue(ic.AgentName)
	ic.AgentNamespace = GetActualValue(ic.AgentNamespace)
	ic.RotatorName = GetActualValue(ic.RotatorName)
	ic.TargetReadReplicaIDAnnotationsKey = GetActualValue(ic.TargetReadReplicaIDAnnotationsKey)
	ic.ReadReplicaLabelKey = GetActualValue(ic.ReadReplicaLabelKey)

	// Bind for the value struct field. Its Bind method has a pointer receiver.
	ic.JobTemplates.Bind()

	return ic
}
