package config

type ReadreplicaRotate struct {
	AgentName string `json:"agent_name" yaml:"agent_name"`
}

func (r *ReadreplicaRotate) Bind() *ReadreplicaRotate {
	r.AgentName = GetActualValue(r.AgentName)
	return r
}
