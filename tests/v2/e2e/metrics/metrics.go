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

// Package metrics provides a metrics collection and aggregation system for E2E tests.
// It supports recording latencies, queue waits, error counts, and custom counters,
// aggregated both globally and in time/request-based windows (Scales).
package metrics

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// requestIDCtxKey is the key for storing the request ID in the context.
// It's an unexported type to prevent collisions with context keys from other packages.
type requestIDCtxKey struct{}

// WithRequestID attaches a request ID to the context for RangeScale bucketing.
// This allows tracking metrics for specific requests across different operations.
func WithRequestID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, requestIDCtxKey{}, id)
}

// requestIDFromCtx retrieves the request ID from the context.
// It returns the ID and a boolean indicating whether the ID was found.
func requestIDFromCtx(ctx context.Context) (uint64, bool) {
	id, ok := ctx.Value(requestIDCtxKey{}).(uint64)
	return id, ok
}

// GetRequestResult returns a RequestResult from the pool.
// This helps reduce GC pressure by reusing objects.
func GetRequestResult() *RequestResult {
	return requestResultPool.Get().(*RequestResult)
}

// PutRequestResult returns a RequestResult to the pool.
// It resets all fields to ensure no data leakage or contamination.
func PutRequestResult(rr *RequestResult) {
	rr.Reset()
	requestResultPool.Put(rr)
}

var requestResultPool = sync.Pool{
	New: func() any {
		return new(RequestResult)
	},
}

// RequestResult represents the result of a single request.
type RequestResult struct {
	RequestID string        // request ID
	Status    codes.Code    // gRPC status code
	Err       error         // error content (Status!=OK时)
	Msg       string        // error message
	QueuedAt  time.Time     // time when the request was queued
	StartedAt time.Time     // time when the RPC started
	EndedAt   time.Time     // time when the RPC ended
	QueueWait time.Duration // StartedAt - QueuedAt
	Latency   time.Duration // RPC execution time = EndedAt - StartedAt
}

// Reset resets the RequestResult to its zero value.
func (rr *RequestResult) Reset() {
	rr.RequestID = ""
	rr.Status = 0
	rr.Err = nil
	rr.Msg = ""
	rr.QueuedAt = time.Time{}
	rr.StartedAt = time.Time{}
	rr.EndedAt = time.Time{}
	rr.QueueWait = 0
	rr.Latency = 0
}

// validate ensures the RequestResult has the necessary timing information.
// It calculates Latency and QueueWait if they are zero but timestamps are present.
func (rr *RequestResult) validate() {
	if rr == nil {
		return
	}
	if rr.Latency == 0 && !rr.StartedAt.IsZero() && !rr.EndedAt.IsZero() {
		rr.Latency = rr.EndedAt.Sub(rr.StartedAt)
	}
	if rr.QueueWait == 0 && !rr.QueuedAt.IsZero() && !rr.StartedAt.IsZero() {
		rr.QueueWait = rr.StartedAt.Sub(rr.QueuedAt)
	}
}

// --- Counter ---

// CounterHandle provides a direct, efficient way to increment a registered counter.
// It holds a pointer to an atomic value for thread-safe updates.
type CounterHandle struct {
	value *atomic.Uint64
}

// Inc increments the counter by 1.
func (h *CounterHandle) Inc() {
	if h == nil || h.value == nil {
		return
	}
	h.value.Add(1)
}

// Add increments the counter by a given value.
func (h *CounterHandle) Add(val int64) {
	if h == nil || h.value == nil {
		return
	}
	if val < 0 {
		return // Counters do not support negative increments
	}
	h.value.Add(uint64(val))
}

// --- Collector ---

// collector is the main entry point for metrics aggregation. It is thread-safe.
// It manages global metrics, per-window scales, and custom counters.
type collector struct {
	mu sync.RWMutex
	id uint64 // Unique ID for lock ordering

	// Atomic counters for total and errored requests.
	total  atomic.Uint64
	errors atomic.Uint64

	// Time range of the collection.
	startTime   atomic.Int64
	lastUpdated atomic.Int64

	// Global metrics.
	latencies      Histogram
	queueWaits     Histogram
	latPercentiles TDigest
	qwPercentiles  TDigest
	exemplars      Exemplar

	// Custom counters, stored in a map for thread-safe access.
	counters map[string]*CounterHandle

	// Thread-safe slice of Scale for metrics (ring buffers).
	scales []Scale

	// gRPC status code counts, using a lock-free array for performance.
	// Array indices correspond to gRPC status codes (0-16).
	// Index 17 is for unknown codes, 18 for non-gRPC errors.
	codes [20]atomic.Uint64
}

var collectorIDCounter atomic.Uint64

// NewCollector creates and initializes a new Collector with the provided options.
func NewCollector(opts ...Option) (Collector, error) {
	c := &collector{
		id:       collectorIDCounter.Add(1),
		counters: make(map[string]*CounterHandle),
	}

	// Prepend default options and apply all options.
	opts = append(defaultOptions, opts...)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	for name := range c.counters {
		c.counters[name] = &CounterHandle{
			value: new(atomic.Uint64),
		}
	}
	return c, nil
}

// Reset resets the collector and all its components.
// It clears all recorded data but maintains the configuration and capacity.
func (c *collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total.Store(0)
	c.errors.Store(0)
	if c.latencies != nil {
		c.latencies.Reset()
	}
	if c.queueWaits != nil {
		c.queueWaits.Reset()
	}
	if c.latPercentiles != nil {
		c.latPercentiles.Reset()
	}
	if c.qwPercentiles != nil {
		c.qwPercentiles.Reset()
	}
	if c.exemplars != nil {
		c.exemplars.Reset()
	}
	for _, h := range c.counters {
		h.value.Store(0)
	}
	for i := range c.codes {
		c.codes[i].Store(0)
	}
	for _, s := range c.scales {
		s.Reset()
	}
}

// MergeInto merges this collector's metrics into the destination collector.
// It uses the internal `merge` method where the receiver is the source and the argument is the destination.
func (c *collector) MergeInto(dest Collector) error {
	return c.merge(dest.(*collector))
}

// merge performs the actual merging logic.
// It acquires locks on both collectors (ordered by unique ID) to ensure safety.
func (c *collector) merge(other *collector) error {
	if c == other || other == nil {
		return nil
	}

	// To prevent deadlocks, always lock in a consistent order.
	if c.id < other.id {
		c.mu.Lock()
		other.mu.Lock()
	} else {
		other.mu.Lock()
		c.mu.Lock()
	}
	defer c.mu.Unlock()
	defer other.mu.Unlock()

	other.total.Add(c.total.Load())
	other.errors.Add(c.errors.Load())

	if t := c.startTime.Load(); t > 0 {
		for {
			curr := other.startTime.Load()
			if curr != 0 && t >= curr {
				break
			}
			if other.startTime.CompareAndSwap(curr, t) {
				break
			}
		}
	}
	if t := c.lastUpdated.Load(); t > 0 {
		for {
			curr := other.lastUpdated.Load()
			if t <= curr {
				break
			}
			if other.lastUpdated.CompareAndSwap(curr, t) {
				break
			}
		}
	}

	if c.latencies != nil && other.latencies != nil {
		if err := other.latencies.Merge(c.latencies); err != nil {
			return err
		}
	}
	if c.queueWaits != nil && other.queueWaits != nil {
		if err := other.queueWaits.Merge(c.queueWaits); err != nil {
			return err
		}
	}
	if c.latPercentiles != nil && other.latPercentiles != nil {
		if err := other.latPercentiles.Merge(c.latPercentiles); err != nil {
			return err
		}
	}
	if c.qwPercentiles != nil && other.qwPercentiles != nil {
		if err := other.qwPercentiles.Merge(c.qwPercentiles); err != nil {
			return err
		}
	}
	if c.exemplars != nil && other.exemplars != nil {
		if err := other.exemplars.Merge(c.exemplars); err != nil {
			return err
		}
	}
	for name, h := range c.counters {
		if mh, ok := other.counters[name]; ok {
			mh.value.Add(h.value.Load())
		} else {
			other.counters[name] = &CounterHandle{
				value: new(atomic.Uint64),
			}
			other.counters[name].value.Store(h.value.Load())
		}
	}
	for i, s := range c.scales {
		if i < len(other.scales) {
			if err := other.scales[i].Merge(s); err != nil {
				return err
			}
		} else {
			other.scales = append(other.scales, s)
		}
	}

	// Merge codes
	for i := range c.codes {
		other.codes[i].Add(c.codes[i].Load())
	}

	return nil
}

// Record processes a single RequestResult, updating all relevant metrics.
// This method is optimized for high-throughput, low-latency execution.
// It updates global histograms, t-digests, exemplars, and all configured scales.
func (c *collector) Record(ctx context.Context, rr *RequestResult) {
	rr.validate()

	if s := rr.StartedAt.UnixNano(); s > 0 {
		for {
			curr := c.startTime.Load()
			if curr != 0 && s >= curr {
				break
			}
			if c.startTime.CompareAndSwap(curr, s) {
				break
			}
		}
	}
	if e := rr.EndedAt.UnixNano(); e > 0 {
		for {
			curr := c.lastUpdated.Load()
			if e <= curr {
				break
			}
			if c.lastUpdated.CompareAndSwap(curr, e) {
				break
			}
		}
	}

	c.total.Add(1)
	if rr.Err != nil {
		c.errors.Add(1)
	}

	if c.latencies != nil {
		c.latencies.Record(float64(rr.Latency.Nanoseconds()))
	}
	if c.queueWaits != nil {
		c.queueWaits.Record(float64(rr.QueueWait.Nanoseconds()))
	}
	if c.latPercentiles != nil {
		c.latPercentiles.Add(float64(rr.Latency.Nanoseconds()))
	}
	if c.qwPercentiles != nil {
		c.qwPercentiles.Add(float64(rr.QueueWait.Nanoseconds()))
	}
	if c.exemplars != nil {
		c.exemplars.Offer(rr.Latency, rr.RequestID, rr.Err, rr.Msg)
	}

	code := resolveStatusCode(rr.Status, rr.Err)
	if int(code) < len(c.codes) {
		c.codes[code].Add(1)
	} else {
		c.codes[codes.Unknown].Add(1) // Fallback for out-of-range codes
	}

	for _, s := range c.scales {
		s.Record(ctx, rr)
	}
}

// CounterHandle returns a handle for a pre-registered custom counter.
func (c *collector) CounterHandle(name string) (*CounterHandle, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if h, ok := c.counters[name]; ok {
		return h, nil
	}
	return nil, errors.New("counter not found")
}

// IncCounter increments a custom counter by a given value.
// It is a convenience wrapper and may be slightly slower than using a CounterHandle.
func (c *collector) IncCounter(name string, val int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if h, ok := c.counters[name]; ok {
		h.Add(val)
	}
}

// Clone returns a deep copy of the collector.
// It creates a new independent collector with copied data, suitable for snapshotting or isolation.
func (c *collector) Clone() (Collector, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	newC := &collector{
		id:       collectorIDCounter.Add(1),
		counters: make(map[string]*CounterHandle),
	}

	// Copy atomics
	newC.total.Store(c.total.Load())
	newC.errors.Store(c.errors.Load())
	newC.startTime.Store(c.startTime.Load())
	newC.lastUpdated.Store(c.lastUpdated.Load())

	// Copy global metrics using Clone
	if c.latencies != nil {
		newC.latencies = c.latencies.Clone()
	}
	if c.queueWaits != nil {
		newC.queueWaits = c.queueWaits.Clone()
	}
	if c.latPercentiles != nil {
		newC.latPercentiles = c.latPercentiles.Clone()
	}
	if c.qwPercentiles != nil {
		newC.qwPercentiles = c.qwPercentiles.Clone()
	}
	if c.exemplars != nil {
		newC.exemplars = c.exemplars.Clone()
	}

	// Copy counters
	for name, h := range c.counters {
		newC.counters[name] = &CounterHandle{
			value: new(atomic.Uint64),
		}
		newC.counters[name].value.Store(h.value.Load())
	}

	// Copy codes
	for i := range c.codes {
		newC.codes[i].Store(c.codes[i].Load())
	}

	// Copy scales
	if len(c.scales) > 0 {
		newC.scales = make([]Scale, len(c.scales))
		for i, s := range c.scales {
			newC.scales[i] = s.Clone()
		}
	}

	return newC, nil
}

// --- Snapshots ---

// GlobalSnapshot returns a snapshot of all aggregated metrics since initialization.
func (c *collector) GlobalSnapshot() *GlobalSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	codesMap := make(map[codes.Code]uint64)
	for i := range c.codes {
		if v := c.codes[i].Load(); v > 0 {
			codesMap[codes.Code(i)] = v
		}
	}
	var latSnap, qwSnap *HistogramSnapshot
	var exSnap []*ExemplarItem
	var exDetails *ExemplarDetails
	if c.latencies != nil {
		latSnap = c.latencies.Snapshot()
	}
	if c.queueWaits != nil {
		qwSnap = c.queueWaits.Snapshot()
	}
	if c.exemplars != nil {
		exSnap = c.exemplars.Snapshot()
		exDetails, _ = c.exemplars.DetailedSnapshot()
		// Sort Average exemplars by distance to mean latency
		if exDetails != nil && len(exDetails.Average) > 0 && latSnap != nil && latSnap.Total > 0 {
			mean := latSnap.Mean
			slices.SortFunc(exDetails.Average, func(a, b *ExemplarItem) int {
				return cmp.Compare(
					math.Abs(float64(a.Latency)-mean),
					math.Abs(float64(b.Latency)-mean))
			})
		}
	}
	var startTime, lastUpdated time.Time
	if t := c.startTime.Load(); t > 0 {
		startTime = time.Unix(0, t)
	}
	if t := c.lastUpdated.Load(); t > 0 {
		lastUpdated = time.Unix(0, t)
	}

	return &GlobalSnapshot{
		Total:           c.total.Load(),
		Errors:          c.errors.Load(),
		StartTime:       startTime,
		LastUpdated:     lastUpdated,
		Latencies:       latSnap,
		QueueWaits:      qwSnap,
		LatPercentiles:  c.latPercentiles,
		QWPercentiles:   c.qwPercentiles,
		Exemplars:       exSnap,
		ExemplarDetails: exDetails,
		Codes:           codesMap,
	}
}

// RangeScalesSnapshot returns snapshots for all configured range-based windows.
func (c *collector) RangeScalesSnapshot() map[string]*ScaleSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshots := make(map[string]*ScaleSnapshot, len(c.scales))
	for _, s := range c.scales {
		if s.Type() == RangeScale {
			snapshots[s.Name()] = s.Snapshot()
		}
	}
	return snapshots
}

// TimeScalesSnapshot returns snapshots for all configured time-based windows.
func (c *collector) TimeScalesSnapshot() map[string]*ScaleSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshots := make(map[string]*ScaleSnapshot, len(c.scales))
	for _, s := range c.scales {
		if s.Type() == TimeScale {
			snapshots[s.Name()] = s.Snapshot()
		}
	}
	return snapshots
}

// --- Merging ---

// MergeCollectors combines multiple collectors into a new one.
// It returns an error if the collectors have incompatible configurations.
func MergeCollectors(collectors ...Collector) (Collector, error) {
	if len(collectors) == 0 {
		return nil, nil
	}
	if len(collectors) == 1 {
		return collectors[0], nil
	}

	base := collectors[0]
	merged, err := base.Clone()
	if err != nil {
		return nil, err
	}

	for _, c := range collectors[1:] {
		if err := c.MergeInto(merged); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

// MergeSnapshots combines multiple snapshots into a new one.
// It returns an error if the snapshots are incompatible.
func MergeSnapshots(snapshots ...*GlobalSnapshot) (*GlobalSnapshot, error) {
	if len(snapshots) == 0 {
		return nil, nil
	}
	if len(snapshots) == 1 {
		return snapshots[0], nil
	}

	merged := &GlobalSnapshot{
		Latencies:      &HistogramSnapshot{},
		QueueWaits:     &HistogramSnapshot{},
		LatPercentiles: snapshots[0].LatPercentiles,
		QWPercentiles:  snapshots[0].QWPercentiles,
		Codes:          make(map[codes.Code]uint64),
	}
	base := snapshots[0]
	merged.StartTime = base.StartTime
	merged.LastUpdated = base.LastUpdated

	// Validate compatibility
	for _, s := range snapshots[1:] {
		if s.BoundsHash != base.BoundsHash || s.SketchKind != base.SketchKind {
			return nil, errors.New("incompatible snapshots")
		}
	}

	// Merge all fields from the snapshots.
	for _, s := range snapshots {
		merged.Total += s.Total
		merged.Errors += s.Errors

		if !s.StartTime.IsZero() && (merged.StartTime.IsZero() || s.StartTime.Before(merged.StartTime)) {
			merged.StartTime = s.StartTime
		}
		if !s.LastUpdated.IsZero() && s.LastUpdated.After(merged.LastUpdated) {
			merged.LastUpdated = s.LastUpdated
		}

		if s.Latencies != nil {
			if merged.Latencies == nil {
				merged.Latencies = &HistogramSnapshot{}
			}
			if err := merged.Latencies.Merge(s.Latencies); err != nil {
				return nil, err
			}
		}
		if s.QueueWaits != nil {
			if merged.QueueWaits == nil {
				merged.QueueWaits = &HistogramSnapshot{}
			}
			if err := merged.QueueWaits.Merge(s.QueueWaits); err != nil {
				return nil, err
			}
		}
		if s.LatPercentiles != nil {
			if merged.LatPercentiles == nil {
				merged.LatPercentiles, _ = NewTDigest(defaultTDigestOpts...)
			}
			if err := merged.LatPercentiles.Merge(s.LatPercentiles); err != nil {
				return nil, err
			}
		}
		if s.QWPercentiles != nil {
			if merged.QWPercentiles == nil {
				merged.QWPercentiles, _ = NewTDigest(defaultTDigestOpts...)
			}
			if err := merged.QWPercentiles.Merge(s.QWPercentiles); err != nil {
				return nil, err
			}
		}
		for code, val := range s.Codes {
			merged.Codes[code] += val
		}
	}

	return merged, nil
}

// --- Data Structures for Snapshots ---

// GlobalSnapshot contains the aggregated metrics for all requests.
type GlobalSnapshot struct {
	Total           uint64                `json:"total"`
	Errors          uint64                `json:"errors"`
	StartTime       time.Time             `json:"start_time"`
	LastUpdated     time.Time             `json:"last_updated"`
	Latencies       *HistogramSnapshot    `json:"latencies"`
	QueueWaits      *HistogramSnapshot    `json:"queue_waits"`
	LatPercentiles  TDigest               `json:"lat_percentiles"`
	QWPercentiles   TDigest               `json:"qw_percentiles"`
	Exemplars       []*ExemplarItem       `json:"exemplars"`
	ExemplarDetails *ExemplarDetails      `json:"exemplar_details,omitempty"`
	Codes           map[codes.Code]uint64 `json:"codes"`
	SchemaVersion   string                `json:"schema_version"`
	BoundsHash      uint64                `json:"bounds_hash"`
	SketchKind      string                `json:"sketch_kind"`
	InvariantsOK    bool                  `json:"invariants_ok"`
}

// MarshalJSON implements the json.Marshaler interface.
func (s *GlobalSnapshot) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

// String implements the fmt.Stringer interface.
func (s *GlobalSnapshot) String() string {
	return NewSnapshotPresenter(s).AsString()
}

// ScaleSnapshot contains the aggregated metrics for a set of windows (slots).
type ScaleSnapshot struct {
	Name     string          `json:"name"`
	Width    uint64          `json:"width"`
	Capacity uint64          `json:"capacity"`
	Slots    []*SlotSnapshot `json:"slots"`
}

// MarshalJSON implements the json.Marshaler interface.
func (s *ScaleSnapshot) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

// String implements the fmt.Stringer interface.
func (s *ScaleSnapshot) String() string {
	if s == nil {
		return ""
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "\n--- Scale: %s (Width: %d, Capacity: %d) ---\n", s.Name, s.Width, s.Capacity)

	totalRequests := uint64(0)
	totalErrors := uint64(0)
	for _, slot := range s.Slots {
		totalRequests += slot.Total
		totalErrors += slot.Errors
	}

	if totalRequests == 0 {
		fmt.Fprint(&sb, "No data collected in this scale.\n")
		return sb.String()
	}

	for i, slot := range s.Slots {
		if slot.Total > 0 {
			fmt.Fprintf(&sb, "\n--- Slot %d ---\n", i)
			fmt.Fprint(&sb, slot.String())
		}
	}
	return sb.String()
}

// SlotSnapshot contains the aggregated metrics for a single window.
type SlotSnapshot struct {
	Total       uint64             `json:"total"`
	Errors      uint64             `json:"errors"`
	LastUpdated int64              `json:"last_updated"`
	Latencies   *HistogramSnapshot `json:"latencies"`
	QueueWaits  *HistogramSnapshot `json:"queue_waits"`
	Counters    []uint64           `json:"counters"`
	Exemplars   []*ExemplarItem    `json:"exemplars"`
}

// String implements the fmt.Stringer interface.
func (s *SlotSnapshot) String() string {
	if s == nil || s.Total == 0 {
		return "No data collected in this slot.\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Total Requests:\t%d\n", s.Total)
	fmt.Fprintf(&sb, "Errors:\t%d (%.2f%%)\n", s.Errors, float64(s.Errors)/float64(s.Total)*100)
	fmt.Fprintf(&sb, "Last Updated:\t%s\n", time.Unix(0, s.LastUpdated))

	fmt.Fprint(&sb, "\nLatency:\n")
	if s.Latencies != nil {
		fmt.Fprint(&sb, s.Latencies.String())
	}

	fmt.Fprint(&sb, "\nQueue Wait:\n")
	if s.QueueWaits != nil {
		fmt.Fprint(&sb, s.QueueWaits.String())
	}

	fmt.Fprintf(&sb, "\nExemplars (Top %d slowest requests):\n", len(s.Exemplars))
	for _, ex := range s.Exemplars {
		fmt.Fprintf(&sb, "\t- RequestID:\t%s,\tLatency:\t%s\n", ex.RequestID, ex.Latency)
	}

	return sb.String()
}

// resolveStatusCode determines the final gRPC status code for a request.
// It prioritizes specific error types (e.g., Canceled, DeadlineExceeded)
// when the status code is initially OK but an error is present.
func resolveStatusCode(code codes.Code, err error) codes.Code {
	if code == codes.OK && err != nil {
		if st, ok := status.FromError(err); ok {
			return st.Code()
		} else if errors.Is(err, context.Canceled) {
			return codes.Canceled
		} else if errors.Is(err, context.DeadlineExceeded) {
			return codes.DeadlineExceeded
		} else {
			return codes.Unknown
		}
	}
	return code
}
