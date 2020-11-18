package model

import "time"

// Job -
type Job struct {
	Name      string
	Namespace string
	Active    int32
	StartTime time.Time
}

// Pod & PodMetrics
type Pod struct {
	Name        string
	Namespace   string
	MemoryLimit float64
	MemoryUsage float64
}
