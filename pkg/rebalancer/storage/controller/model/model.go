package model

// Pod & PodMetrics
type Pod struct {
	Name        string
	Namespace   string
	MemoryLimit float64
	MemoryUsage float64
}

type StatefulSet struct {
	Name            string
	Namespace       string
	DesiredReplicas *int32 // StatefulSetSpec.Replicas
	Replicas        int32  // StatefulSetStatus.Replicas
}
