package metrics

type SearchMetrics struct {
	Recall float64
	Qps float64
	Epsilon float32
}

type Metrics struct {
	BuildTime        int64
	DatasetName      string
	SearchEdgeSize   int
	CreationEdgeSize int
	K                int
	Search           []*SearchMetrics
}

