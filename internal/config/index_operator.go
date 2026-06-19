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
	// JobTemplates represents the job templates.
	JobTemplates IndexJobTemplates `json:"job_templates" yaml:"job_templates"`
	// Namespace represents the namespace.
	Namespace string `json:"namespace" yaml:"namespace"`
	// AgentName represents the agent name.
	AgentName string `json:"agent_name" yaml:"agent_name"`
	// AgentNamespace represents the agent namespace.
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`
	// RotatorName represents the rotator name.
	RotatorName string `json:"rotator_name" yaml:"rotator_name"`
	// TargetReadReplicaIDAnnotationsKey represents the target read replica ID annotations key.
	TargetReadReplicaIDAnnotationsKey string `json:"target_read_replica_id_annotations_key" yaml:"target_read_replica_id_annotations_key"`
	// ReadReplicaLabelKey represents the read replica label key.
	ReadReplicaLabelKey string `json:"read_replica_label_key" yaml:"read_replica_label_key"`
	// RotationJobConcurrency represents the rotation job concurrency.
	RotationJobConcurrency uint `json:"rotation_job_concurrency" yaml:"rotation_job_concurrency"`
	// ReadReplicaEnabled enables read replica.
	ReadReplicaEnabled bool `json:"read_replica_enabled" yaml:"read_replica_enabled"`
}

type IndexJobTemplates struct {
	// Rotate represents the rotate job template.
	Rotate *k8s.Job `json:"rotate" yaml:"rotate"`
	// Creation represents the creation job template.
	Creation *k8s.Job `json:"creation" yaml:"creation"`
	// Save represents the save job template.
	Save *k8s.Job `json:"save" yaml:"save"`
	// Correction represents the correction job template.
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
