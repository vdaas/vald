package ngt

type AccuracyEntry struct {
	Epsilon  float64 `json:"epsilon"`
	Accuracy float64 `json:"accuracy"`
}

type Result struct {
	AccuracyTable                         []AccuracyEntry `json:"AccuracyTable"`
	BatchSizeForCreation                  string          `json:"BatchSizeForCreation"`
	BuildTimeLimit                        string          `json:"BuildTimeLimit"`
	DatabaseType                          string          `json:"DatabaseType"`
	Dimension                             string          `json:"Dimension"`
	DistanceType                          string          `json:"DistanceType"`
	DynamicEdgeSizeBase                   string          `json:"DynamicEdgeSizeBase"`
	DynamicEdgeSizeRate                   string          `json:"DynamicEdgeSizeRate"`
	EdgeSizeForCreation                   string          `json:"EdgeSizeForCreation"`
	EdgeSizeForSearch                     string          `json:"EdgeSizeForSearch"`
	EdgeSizeLimitForCreation              string          `json:"EdgeSizeLimitForCreation"`
	EpsilonForCreation                    string          `json:"EpsilonForCreation"`
	EpsilonForInsertionOrder              string          `json:"EpsilonForInsertionOrder"`
	EpsilonType                           string          `json:"EpsilonType"`
	GraphType                             string          `json:"GraphType"`
	IncomingEdge                          string          `json:"IncomingEdge"`
	IncrimentalEdgeSizeLimitForTruncation string          `json:"IncrimentalEdgeSizeLimitForTruncation"`
	IndexType                             string          `json:"IndexType"`
	MaxMagnitude                          string          `json:"MaxMagnitude"`
	NumberOfNeighborsForInsertionOrder    string          `json:"NumberOfNeighborsForInsertionOrder"`
	ObjectAlignment                       string          `json:"ObjectAlignment"`
	ObjectType                            string          `json:"ObjectType"`
	OutgoingEdge                          string          `json:"OutgoingEdge"`
	PathAdjustmentInterval                string          `json:"PathAdjustmentInterval"`
	PrefetchOffset                        string          `json:"PrefetchOffset"`
	PrefetchSize                          string          `json:"PrefetchSize"`
	QuantizationClippingRate              string          `json:"QuantizationClippingRate"`
	QuantizationOffset                    string          `json:"QuantizationOffset"`
	QuantizationScale                     string          `json:"QuantizationScale"`
	RefinementObjectType                  string          `json:"RefinementObjectType"`
	SeedSize                              string          `json:"SeedSize"`
	SeedType                              string          `json:"SeedType"`
	ThreadPoolSize                        string          `json:"ThreadPoolSize"`
	TruncationThreadPoolSize              string          `json:"TruncationThreadPoolSize"`
}

type Metadata struct {
	DataSize int    `json:"data_size"`
	Error    string `json:"error"`
}

type Ngt struct {
	Result   Result    `json:"result"`
	Metadata *Metadata `json:"metadata,omitempty"`
}
