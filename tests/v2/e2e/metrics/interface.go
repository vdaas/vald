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
	"time"
	"unsafe"

	"github.com/zeebo/xxh3"
)

// sizeOfFloat64 is the size of float64 in bytes.
const sizeOfFloat64 = 8

// Recorder is the interface for recording metrics.
type Recorder interface {
	Record(ctx context.Context, rr *RequestResult)
}

// Reporter is the interface for reporting and exporting data.
type Reporter interface {
	GlobalSnapshot() *GlobalSnapshot
	RangeScalesSnapshot() map[string]*ScaleSnapshot
	TimeScalesSnapshot() map[string]*ScaleSnapshot
}

// CounterManager is the interface for managing custom counters.
type CounterManager interface {
	CounterHandle(name string) (*CounterHandle, error)
	IncCounter(name string, val int64)
}

// Collector is the interface for the metrics collector.
type Collector interface {
	Recorder
	Reporter
	CounterManager
	MergeInto(dest Collector) error
	Clone() (Collector, error)
	merge(other *collector) error
	Reset()
}

// Histogram is the interface for a histogram.
type Histogram interface {
	Record(val float64)
	Merge(other Histogram) error
	Snapshot() *HistogramSnapshot
	BoundsHash() uint64
	Clone() Histogram
	Reset()
}

// Exemplar is the interface for an exemplar.
type Exemplar interface {
	Offer(latency time.Duration, requestID string, err error, msg string)
	Snapshot() []*ExemplarItem
	DetailedSnapshot() (*ExemplarDetails, error)
	Merge(other Exemplar) error
	Clone() Exemplar
	Reset()
}

// TDigest is the interface for approximate percentile estimators like t-digest.
type TDigest interface {
	Add(value float64)
	Quantile(q float64) float64
	CDF(value, min, max float64) float64
	ForEachCentroid(f func(mean, weight float64) bool)
	Quantiles() []float64
	Merge(other TDigest) error
	Clone() TDigest
	Reset()
}

// Slot is the interface for a time/range slot in a Scale.
type Slot interface {
	Record(rr *RequestResult, windowIdx uint64)
	Merge(other Slot) error
	Snapshot() *SlotSnapshot
	Reset()
	Clone() Slot
}

// Scale is the interface for a metrics scale (ring buffer of slots).
type Scale interface {
	Record(ctx context.Context, rr *RequestResult)
	Merge(other Scale) error
	Snapshot() *ScaleSnapshot
	Reset()
	Clone() Scale
	Type() ScaleType
	Name() string
}

// ScaleType represents the type of scale (Range or Time).
type ScaleType uint8

const (
	// RangeScale buckets metrics by request ID.
	RangeScale ScaleType = iota
	// TimeScale buckets metrics by time.
	TimeScale
)

// computeHash computes a xxh3 checksum for the given slice of float64 value.
func computeHash(vals ...float64) uint64 {
	if len(vals) == 0 {
		return 0
	}
	// Reinterpret the float64 slice as a byte slice to avoid allocation using unsafe.Slice.
	// This is safe because both slices point to the same underlying data.
	data := unsafe.Slice((*byte)(unsafe.Pointer(&vals[0])), len(vals)*sizeOfFloat64) //nolint:gosec // using unsafe for zero-copy slice conversion for performance
	return xxh3.Hash(data)
}

// shardIndex calculates the shard index for a given hash and number of shards.
func shardIndex(hash uint64, n int) int {
	if n <= 1 {
		return 0
	}
	// Optimized for power of 2
	if (n & (n - 1)) == 0 {
		// Fix G115: hash & uint64(n-1) will be in [0, n-1].
		return int(hash & uint64(n-1)) //nolint:gosec // bitwise and with positive int fits in int
	}
	// Fix G115: hash % uint64(n) will be in [0, n-1]. Since n is int and n > 0, the result fits in int.
	return int(hash % uint64(n)) //nolint:gosec // hash modulo length is always within int bounds
}
