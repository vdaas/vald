//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package metrics

import (
	"context"
	"fmt"
	"time"
)

// Collector is the interface for the metrics collector.
type Collector interface {
	Record(ctx context.Context, rr *RequestResult)
	Merge(other Collector) error
	GlobalSnapshot() *GlobalSnapshot
	RangeScalesSnapshot() map[string]*ScaleSnapshot
	TimeScalesSnapshot() map[string]*ScaleSnapshot
	CounterHandle(name string) (*CounterHandle, error)
	IncCounter(name string, val int64)
	Clone() (Collector, error)
	merge(other *collector) error
	Reset()
}

// Histogram is the interface for a histogram.
type Histogram interface {
	Record(val float64)
	Merge(other Histogram) error
	Snapshot() *HistogramSnapshot
	BoundsCRC32() uint32
	Clone() Histogram
	merge(other *histogram) error
	Reset()
}

// Exemplar is the interface for an exemplar.
type Exemplar interface {
	Offer(latency time.Duration, requestID string)
	Snapshot() []*item
	Clone() Exemplar
	Reset()
}

// TDigest is the interface for approximate percentile estimators like t-digest.
type TDigest interface {
	Add(value float64)
	Quantile(q float64) float64
	Quantiles() []float64
	Merge(other TDigest) error
	Clone() TDigest
	fmt.Stringer
	Reset()
}
