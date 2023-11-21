package config

type ReadreplicaRotate struct {
	AgentNamespace      string `json:"agent_namespace" yaml:"agent_namespace"`
	ReadReplicaLabelKey string `json:"read_replica_label_key" yaml:"read_replica_label_key"`
	ReadReplicaId       string `json:"read_replica_id" yaml:"read_replica_id"`
	VolumeName          string `json:"volume_name" yaml:"volume_name"`
}

func (r *ReadreplicaRotate) Bind() *ReadreplicaRotate {
	r.AgentNamespace = GetActualValue(r.AgentNamespace)
	r.ReadReplicaLabelKey = GetActualValue(r.ReadReplicaLabelKey)
	r.ReadReplicaId = GetActualValue(r.ReadReplicaId)
	r.VolumeName = GetActualValue(r.VolumeName)

	return r
}
