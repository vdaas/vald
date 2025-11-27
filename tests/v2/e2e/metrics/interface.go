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
	"unsafe"

	"github.com/zeebo/xxh3"
)

// Collector is the interface for the metrics collector.
type Collector interface {
	Record(ctx context.Context, rr *RequestResult)
	MergeInto(dest Collector) error
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
	BoundsHash() uint64
	Clone() Histogram
	merge(other *histogram) error
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
	Quantiles() []float64
	Merge(other TDigest) error
	Clone() TDigest
	fmt.Stringer
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

// sliceHeader is a stripped-down version of reflect.SliceHeader used for unsafe pointer conversions.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// computeHash computes a xxh3 checksum for the given slice of float64 value.
func computeHash(vals ...float64) uint64 {
	if len(vals) == 0 {
		return 0
	}
	// Reinterpret the float64 slice as a byte slice to avoid allocation.
	// This is safe because both slices point to the same underlying data.
	header := (*sliceHeader)(unsafe.Pointer(&vals)) //nolint:gosec
	header.Len *= 8
	header.Cap *= 8
	return xxh3.Hash(*(*[]byte)(unsafe.Pointer(header))) //nolint:gosec
}
