// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package metrics

type SearchMetrics struct {
	// Recall of the search.
	Recall  float64
	// Queries per second.
	Qps     float64
	// Epsilon value for the search.
	Epsilon float32
}

type Metrics struct {
	// Name of the dataset.
	DatasetName      string
	// Search metrics.
	Search           []*SearchMetrics
	// Build time of the index.
	BuildTime        int64
	// Search edge size.
	SearchEdgeSize   int
	// Creation edge size.
	CreationEdgeSize int
	// Number of nearest neighbors.
	K                int
}
