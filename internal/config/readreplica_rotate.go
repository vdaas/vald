// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

type ReadreplicaRotate struct {
	AgentNamespace      string `json:"agent_namespace"        yaml:"agent_namespace"`
	ReadReplicaLabelKey string `json:"read_replica_label_key" yaml:"read_replica_label_key"`
	ReadReplicaId       string `json:"read_replica_id"        yaml:"read_replica_id"`
	VolumeName          string `json:"volume_name"            yaml:"volume_name"`
}

func (r *ReadreplicaRotate) Bind() *ReadreplicaRotate {
	r.AgentNamespace = GetActualValue(r.AgentNamespace)
	r.ReadReplicaLabelKey = GetActualValue(r.ReadReplicaLabelKey)
	r.ReadReplicaId = GetActualValue(r.ReadReplicaId)
	r.VolumeName = GetActualValue(r.VolumeName)

	return r
}
