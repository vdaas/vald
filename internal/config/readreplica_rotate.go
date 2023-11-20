package config

type ReadreplicaRotate struct {
	AgentNamespace   string `json:"agent_namespace" yaml:"agent_namespace"`
	DeploymentPrefix string `json:"deployment_prefix" yaml:"deployment_prefix"`
	SnapshotPrefix   string `json:"snapshot_prefix" yaml:"snapshot_prefix"`
	PvcPrefix        string `json:"pvc_prefix" yaml:"pvc_prefix"`
	VolumeName       string `json:"volume_name" yaml:"volume_name"`
}

func (r *ReadreplicaRotate) Bind() *ReadreplicaRotate {
	r.AgentNamespace = GetActualValue(r.AgentNamespace)
	r.DeploymentPrefix = GetActualValue(r.DeploymentPrefix)
	r.SnapshotPrefix = GetActualValue(r.SnapshotPrefix)
	r.PvcPrefix = GetActualValue(r.PvcPrefix)
	r.VolumeName = GetActualValue(r.VolumeName)
	
	return r
}
