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
package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
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
func GetRequestResult() *RequestResult {
	return requestResultPool.Get()
}

// PutRequestResult returns a RequestResult to the pool.
func PutRequestResult(rr *RequestResult) {
	rr.RequestID = ""
	rr.Status = 0
	rr.Err = nil
	rr.Msg = ""
	rr.QueuedAt = time.Time{}
	rr.StartedAt = time.Time{}
	rr.EndedAt = time.Time{}
	rr.QueueWait = 0
	rr.Latency = 0
	requestResultPool.Put(rr)
}

// RequestResult represents the result of a single request.
type RequestResult struct {
	RequestID string        // request ID
	Status    codes.Code    // gRPC status code
	Err       error         // error content (Status!=OK时)
	Msg       string        // status message
	QueuedAt  time.Time     // time when the request was queued
	StartedAt time.Time     // time when the RPC started
	EndedAt   time.Time     // time when the RPC ended
	QueueWait time.Duration // StartedAt - QueuedAt
	Latency   time.Duration // RPC execution time = EndedAt - StartedAt
}

// validate ensures the RequestResult has the necessary timing information.
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
	if rr.Msg == "" {
		if rr.Err != nil {
			rr.Msg = fmt.Sprintf("status code: %s, error: %s", rr.Status.String(), rr.Err.Error())
		} else {
			rr.Msg = rr.Status.String()
		}
	}
}

// --- Counter ---

// CounterHandle provides a direct, efficient way to increment a registered counter.
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

// --- Slot ---

// Slot holds the metrics for a single window in a scale.
type Slot struct {
	Total     atomic.Uint64   // total number of requests in this slot
	Errors    atomic.Uint64   // number of errored requests in this slot
	updatedNS atomic.Int64    // UnixNano timestamp of the last update
	Latency   Histogram       // latency histogram
	QueueWait Histogram       // queue wait histogram
	Counters  []atomic.Uint64 // custom counters
	Exemplars Exemplar        // exemplar heap
}

// newSlot creates a new Slot.
func newSlot(numCounters int, latencies, queueWaits Histogram, exemplars Exemplar) *Slot {
	return &Slot{
		Latency:   latencies,
		QueueWait: queueWaits,
		Counters:  make([]atomic.Uint64, numCounters),
		Exemplars: exemplars,
	}
}

// --- Scales ---

// Scale is a ring buffer of slots for windowed metrics.
type Scale struct {
	mu        sync.RWMutex
	slots     []*Slot // ring buffer of slots
	width     uint64  // width of each slot
	capacity  uint64  // number of slots in the ring buffer
	name      string  // name of the scale
	scaleType scaleType
}

type scaleType uint8

const (
	rangeScale scaleType = iota
	timeScale
)

// record adds a request result to the appropriate slot in the scale.
func (s *Scale) record(ctx context.Context, rr *RequestResult) {
	var idx uint64
	var ok bool

	switch s.scaleType {
	case rangeScale:
		idx, ok = requestIDFromCtx(ctx)
		if !ok {
			return
		}
	case timeScale:
		idx = uint64(rr.EndedAt.Unix())
	}

	slot := s.getSlot(idx)
	slot.Total.Add(1)
	if rr.Err != nil {
		slot.Errors.Add(1)
	}
	slot.updatedNS.Store(rr.EndedAt.UnixNano())
	slot.Latency.Record(float64(rr.Latency.Nanoseconds()))
	slot.QueueWait.Record(float64(rr.QueueWait.Nanoseconds()))
	slot.Exemplars.Offer(rr.Latency, rr.RequestID)
}

// newScale creates a new scale with the given configuration.
func newScale(name string, width, capacity uint64, numCounters int, st scaleType) (*Scale, error) {
	if width == 0 {
		return nil, errors.New("scale width must be > 0")
	}
	if capacity == 0 {
		return nil, errors.New("scale capacity must be > 0")
	}
	slots := make([]*Slot, capacity)
	for i := range slots {
		slots[i] = newSlot(
			numCounters,
			histogramPool.Get(),
			histogramPool.Get(),
			exemplarPool.Get(),
		)
	}
	s := &Scale{
		name:      name,
		width:     width,
		capacity:  capacity,
		slots:     slots,
		scaleType: st,
	}

	return s, nil
}

// getSlot returns the slot for the given index.
func (s *Scale) getSlot(idx uint64) *Slot {
	slotIdx := (idx / s.width) % s.capacity
	return s.slots[slotIdx]
}

// Merge merges another scale into this one.
func (s *Scale) Merge(other *Scale) error {
	if s.width != other.width || s.capacity != other.capacity {
		return errors.New("incompatible scales")
	}
	if s == other {
		return nil
	}
	// To prevent deadlocks, always lock in a consistent order.
	if uintptr(unsafe.Pointer(s)) < uintptr(unsafe.Pointer(other)) {
		s.mu.Lock()
		other.mu.Lock()
	} else {
		other.mu.Lock()
		s.mu.Lock()
	}
	defer s.mu.Unlock()
	defer other.mu.Unlock()

	for i := range s.slots {
		ss := s.slots[i]
		os := other.slots[i]
		ss.Total.Add(os.Total.Load())
		ss.Errors.Add(os.Errors.Load())
		if os.updatedNS.Load() > ss.updatedNS.Load() {
			ss.updatedNS.Store(os.updatedNS.Load())
		}
		if err := ss.Latency.Merge(os.Latency); err != nil {
			return err
		}
		if err := ss.QueueWait.Merge(os.QueueWait); err != nil {
			return err
		}
		for _, ex := range os.Exemplars.Snapshot() {
			ss.Exemplars.Offer(ex.latency, ex.requestID)
		}
	}
	return nil
}

// Snapshot returns a snapshot of the scale.
func (s *Scale) Snapshot() *ScaleSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	slots := make([]*SlotSnapshot, len(s.slots))
	for i, slot := range s.slots {
		counters := make([]uint64, len(slot.Counters))
		for j := range counters {
			counters[j] = slot.Counters[j].Load()
		}
		slots[i] = &SlotSnapshot{
			Total:       slot.Total.Load(),
			Errors:      slot.Errors.Load(),
			LastUpdated: slot.updatedNS.Load(),
			Latencies:   slot.Latency.Snapshot(),
			QueueWaits:  slot.QueueWait.Snapshot(),
			Counters:    counters,
			Exemplars:   slot.Exemplars.Snapshot(),
		}
	}

	return &ScaleSnapshot{
		Name:     s.name,
		Width:    s.width,
		Capacity: s.capacity,
		Slots:    slots,
	}
}

// --- Collector ---

// collector is the main entry point for metrics aggregation. It is thread-safe.
type collector struct {
	mu sync.RWMutex

	// Atomic counters for total and errored requests.
	total  atomic.Uint64
	errors atomic.Uint64

	// Global metrics.
	latencies      Histogram
	queueWaits     Histogram
	latPercentiles QuantileSketch
	qwPercentiles  QuantileSketch
	exemplars      Exemplar

	// Custom counters, stored in a map for thread-safe access.
	counters map[string]*CounterHandle

	// Thread-safe slice of Scale for metrics.
	scales []*Scale

	// gRPC status code counts, stored in a map for thread-safe access.
	codes map[codes.Code]*atomic.Uint64
}

// NewCollector creates and initializes a new Collector with the provided options.
func NewCollector(opts ...Option) (Collector, error) {
	c := &collector{
		counters: make(map[string]*CounterHandle),
		codes:    make(map[codes.Code]*atomic.Uint64),
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

// Merge merges the metrics from another collector into this one.
func (c *collector) Merge(other Collector) error {
	return other.merge(c)
}

func (c *collector) merge(other *collector) error {
	if c == other || other == nil {
		return nil
	}

	// To prevent deadlocks, always lock in a consistent order.
	if uintptr(unsafe.Pointer(c)) < uintptr(unsafe.Pointer(other)) {
		other.mu.Lock()
		c.mu.Lock()
	} else {
		c.mu.Lock()
		other.mu.Lock()
	}
	defer c.mu.Unlock()
	defer other.mu.Unlock()

	other.total.Add(c.total.Load())
	other.errors.Add(c.errors.Load())

	if err := other.latencies.Merge(c.latencies); err != nil {
		return err
	}
	if err := other.queueWaits.Merge(c.queueWaits); err != nil {
		return err
	}
	if err := other.latPercentiles.Merge(c.latPercentiles); err != nil {
		return err
	}
	if err := other.qwPercentiles.Merge(c.qwPercentiles); err != nil {
		return err
	}
	for _, ex := range c.exemplars.Snapshot() {
		other.exemplars.Offer(ex.latency, ex.requestID)
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

	return nil
}

// Record processes a single RequestResult, updating all relevant metrics.
// This method is optimized for high-throughput, low-latency execution.
func (c *collector) Record(ctx context.Context, rr *RequestResult) {
	rr.validate()
	c.total.Add(1)
	if rr.Err != nil {
		c.errors.Add(1)
	}

	c.latencies.Record(float64(rr.Latency.Nanoseconds()))
	c.queueWaits.Record(float64(rr.QueueWait.Nanoseconds()))
	c.latPercentiles.Add(float64(rr.Latency.Nanoseconds()))
	c.qwPercentiles.Add(float64(rr.QueueWait.Nanoseconds()))
	c.exemplars.Offer(rr.Latency, rr.RequestID)
	c.mu.RLock()
	counter, ok := c.codes[rr.Status]
	c.mu.RUnlock()
	if !ok {
		c.mu.Lock()
		// re-check if another goroutine already inserted the key
		counter, ok = c.codes[rr.Status]
		if !ok {
			counter = new(atomic.Uint64)
			c.codes[rr.Status] = counter
		}
		c.mu.Unlock()
	}
	counter.Add(1)

	for _, s := range c.scales {
		s.record(ctx, rr)
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

func (c *collector) deepCopy() (Collector, error) {
	return NewCollector(defaultOptions...)
}

// --- Snapshots ---

// GlobalSnapshot returns a snapshot of all aggregated metrics since initialization.
func (c *collector) GlobalSnapshot() *GlobalSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	codes := make(map[codes.Code]uint64, len(c.codes))
	for code, val := range c.codes {
		codes[code] = val.Load()
	}
	return &GlobalSnapshot{
		Total:          c.total.Load(),
		Errors:         c.errors.Load(),
		Latencies:      c.latencies.Snapshot(),
		QueueWaits:     c.queueWaits.Snapshot(),
		LatPercentiles: c.latPercentiles,
		QWPercentiles:  c.qwPercentiles,
		Exemplars:      c.exemplars.Snapshot(),
		Codes:          codes,
	}
}

// RangeScalesSnapshot returns snapshots for all configured range-based windows.
func (c *collector) RangeScalesSnapshot() map[string]*ScaleSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshots := make(map[string]*ScaleSnapshot, len(c.scales))
	for _, s := range c.scales {
		if s.scaleType == rangeScale {
			snapshots[s.name] = s.Snapshot()
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
		if s.scaleType == timeScale {
			snapshots[s.name] = s.Snapshot()
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
	merged, err := base.deepCopy()
	if err != nil {
		return nil, err
	}

	for _, c := range collectors[1:] {
		if err := merged.Merge(c); err != nil {
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

	// Validate compatibility
	for _, s := range snapshots[1:] {
		if s.BoundsCRC32 != base.BoundsCRC32 || s.SketchKind != base.SketchKind {
			return nil, errors.New("incompatible snapshots")
		}
	}

	// Merge all fields from the snapshots.
	for _, s := range snapshots {
		merged.Total += s.Total
		merged.Errors += s.Errors
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
	Total          uint64                `json:"total"`
	Errors         uint64                `json:"errors"`
	Latencies      *HistogramSnapshot    `json:"latencies"`
	QueueWaits     *HistogramSnapshot    `json:"queue_waits"`
	LatPercentiles QuantileSketch        `json:"lat_percentiles"`
	QWPercentiles  QuantileSketch        `json:"qw_percentiles"`
	Exemplars      []*item               `json:"exemplars"`
	Codes          map[codes.Code]uint64 `json:"codes"`
	SchemaVersion  string                `json:"schema_version"`
	BoundsCRC32    uint32                `json:"bounds_crc32"`
	SketchKind     string                `json:"sketch_kind"`
	InvariantsOK   bool                  `json:"invariants_ok"`
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
	Exemplars   []*item            `json:"exemplars"`
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
		fmt.Fprintf(&sb, "\t- RequestID:\t%s,\tLatency:\t%s\n", ex.requestID, ex.latency)
	}

	return sb.String()
}

// --- Interfaces and supporting types ---

// QuantileSketch defines the interface for approximate percentile estimators like t-digest.
type QuantileSketch interface {
	Add(value float64)
	Quantile(q float64) float64
	Merge(other QuantileSketch) error
	fmt.Stringer
}
