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
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
)

// requestIDKey is the key for storing the request ID in the context.
type requestIDKey struct{}

// WithRequestID attaches a request ID to the context for RangeScale bucketing.
func WithRequestID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// requestIDFromCtx retrieves the request ID from the context.
func requestIDFromCtx(ctx context.Context) (uint64, bool) {
	id, ok := ctx.Value(requestIDKey{}).(uint64)
	return id, ok
}

var requestResultPool = sync.Pool{
	New: func() any {
		return new(RequestResult)
	},
}

// GetRequestResult returns a RequestResult from the pool.
func GetRequestResult() *RequestResult {
	return requestResultPool.Get().(*RequestResult)
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
	h.value.Add(uint64(val))
}

// --- Slot ---

// Slot holds the metrics for a single window in a scale.
type Slot struct {
	Total     atomic.Uint64   // total number of requests in this slot
	Errors    atomic.Uint64   // number of errored requests in this slot
	updatedNS atomic.Int64    // UnixNano timestamp of the last update
	Latency   *Histogram      // latency histogram
	QueueWait *Histogram      // queue wait histogram
	Counters  []atomic.Uint64 // custom counters
	Exemplars *Exemplar       // exemplar heap
}

// newSlot creates a new Slot.
func newSlot(numCounters int, latencies, queueWaits *Histogram, exemplars *Exemplar) *Slot {
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
	mu       sync.RWMutex
	slots    []*Slot // ring buffer of slots
	width    uint64  // width of each slot
	capacity uint64  // number of slots in the ring buffer
	name     string  // name of the scale
}

// newScale creates a new scale with the given configuration.
func newScale(
	name string, width, capacity uint64, numCounters int, hcfg histogramConfig, ecfg exemplarConfig,
) *Scale {
	slots := make([]*Slot, capacity)
	for i := range slots {
		slots[i] = newSlot(
			numCounters,
			NewHistogram(hcfg.min, hcfg.max, hcfg.growth, hcfg.numBuckets, hcfg.numShards),
			NewHistogram(hcfg.min, hcfg.max, hcfg.growth, hcfg.numBuckets, hcfg.numShards),
			NewExemplar(ecfg.capacity),
		)
	}
	return &Scale{
		name:     name,
		width:    width,
		capacity: capacity,
		slots:    slots,
	}
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

// RangeScale aggregates metrics over a range of request IDs.
type RangeScale struct {
	*Scale
}

// NewRangeScale creates a new RangeScale.
func NewRangeScale(
	name string, width, capacity uint64, numCounters int, hcfg histogramConfig, ecfg exemplarConfig,
) *RangeScale {
	return &RangeScale{
		Scale: newScale(name, width, capacity, numCounters, hcfg, ecfg),
	}
}

// Record updates the appropriate slot based on the request ID in the context.
func (rs *RangeScale) Record(ctx context.Context, rr *RequestResult) {
	if reqID, ok := requestIDFromCtx(ctx); ok {
		slot := rs.getSlot(reqID)
		slot.Total.Add(1)
		if rr.Err != nil {
			slot.Errors.Add(1)
		}
		slot.updatedNS.Store(rr.EndedAt.UnixNano())
		slot.Latency.Record(float64(rr.Latency.Nanoseconds()))
		slot.QueueWait.Record(float64(rr.QueueWait.Nanoseconds()))
		slot.Exemplars.Offer(rr.Latency, rr.RequestID)
	}
}

// TimeScale aggregates metrics over a time window.
type TimeScale struct {
	*Scale
}

// NewTimeScale creates a new TimeScale.
func NewTimeScale(
	name string,
	widthSec, capacity uint64,
	numCounters int,
	hcfg histogramConfig,
	ecfg exemplarConfig,
) *TimeScale {
	return &TimeScale{
		Scale: newScale(name, widthSec, capacity, numCounters, hcfg, ecfg),
	}
}

// Record updates the appropriate slot based on the request's end time.
func (ts *TimeScale) Record(rr *RequestResult) {
	idx := uint64(rr.EndedAt.Unix())
	slot := ts.getSlot(idx)
	slot.Total.Add(1)
	if rr.Err != nil {
		slot.Errors.Add(1)
	}
	slot.updatedNS.Store(rr.EndedAt.UnixNano())
	slot.Latency.Record(float64(rr.Latency.Nanoseconds()))
	slot.QueueWait.Record(float64(rr.QueueWait.Nanoseconds()))
	slot.Exemplars.Offer(rr.Latency, rr.RequestID)
}

// --- Collector ---

// Collector is the main entry point for metrics aggregation. It is thread-safe.
type Collector struct {
	mu             sync.RWMutex
	total          atomic.Uint64
	errors         atomic.Uint64
	latencies      *Histogram
	queueWaits     *Histogram
	latPercentiles QuantileSketch
	qwPercentiles  QuantileSketch
	exemplars      *Exemplar
	counters       map[string]*CounterHandle
	rangeScales    []*RangeScale
	timeScales     []*TimeScale

	hcfg histogramConfig
	ecfg exemplarConfig
}

// NewCollector creates and initializes a new Collector with the provided options.
func NewCollector(opts ...Option) *Collector {
	c := &Collector{
		counters: make(map[string]*CounterHandle),
		hcfg:     defaultHistogramConfig,
		ecfg:     defaultExemplarConfig,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.latencies == nil {
		c.latencies = NewHistogram(c.hcfg.min, c.hcfg.max, c.hcfg.growth, c.hcfg.numBuckets, c.hcfg.numShards)
	}
	if c.queueWaits == nil {
		c.queueWaits = NewHistogram(c.hcfg.min, c.hcfg.max, c.hcfg.growth, c.hcfg.numBuckets, c.hcfg.numShards)
	}
	if c.exemplars == nil {
		c.exemplars = NewExemplar(c.ecfg.capacity)
	}

	if c.latPercentiles == nil {
		c.latPercentiles, _ = NewTDigest(defaultTDigestConfig.compression, defaultTDigestConfig.compressionTriggerFactor)
	}
	if c.qwPercentiles == nil {
		c.qwPercentiles, _ = NewTDigest(defaultTDigestConfig.compression, defaultTDigestConfig.compressionTriggerFactor)
	}
	for name := range c.counters {
		c.counters[name] = &CounterHandle{
			value: new(atomic.Uint64),
		}
	}
	return c
}

// Merge merges the metrics from another collector into this one.
func (c *Collector) Merge(other *Collector) error {
	if c == other {
		return nil
	}

	// To prevent deadlocks, always lock in a consistent order.
	if uintptr(unsafe.Pointer(c)) < uintptr(unsafe.Pointer(other)) {
		c.mu.Lock()
		other.mu.Lock()
	} else {
		other.mu.Lock()
		c.mu.Lock()
	}
	defer c.mu.Unlock()
	defer other.mu.Unlock()

	c.total.Add(other.total.Load())
	c.errors.Add(other.errors.Load())

	if err := c.latencies.Merge(other.latencies); err != nil {
		return err
	}
	if err := c.queueWaits.Merge(other.queueWaits); err != nil {
		return err
	}
	if err := c.latPercentiles.Merge(other.latPercentiles); err != nil {
		return err
	}
	if err := c.qwPercentiles.Merge(other.qwPercentiles); err != nil {
		return err
	}
	for _, ex := range other.exemplars.Snapshot() {
		c.exemplars.Offer(ex.latency, ex.requestID)
	}
	for name, h := range other.counters {
		if mh, ok := c.counters[name]; ok {
			mh.value.Add(h.value.Load())
		} else {
			c.counters[name] = h
		}
	}
	for i, rs := range other.rangeScales {
		if i < len(c.rangeScales) {
			if err := c.rangeScales[i].Scale.Merge(rs.Scale); err != nil {
				return err
			}
		} else {
			c.rangeScales = append(c.rangeScales, rs)
		}
	}
	for i, ts := range other.timeScales {
		if i < len(c.timeScales) {
			if err := c.timeScales[i].Scale.Merge(ts.Scale); err != nil {
				return err
			}
		} else {
			c.timeScales = append(c.timeScales, ts)
		}
	}

	return nil
}

// Record processes a single RequestResult, updating all relevant metrics.
// This method is optimized for high-throughput, low-latency execution.
func (c *Collector) Record(ctx context.Context, rr *RequestResult) {
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

	for _, rs := range c.rangeScales {
		rs.Record(ctx, rr)
	}
	for _, ts := range c.timeScales {
		ts.Record(rr)
	}
}

// CounterHandle returns a handle for a pre-registered custom counter.
func (c *Collector) CounterHandle(name string) (*CounterHandle, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if h, ok := c.counters[name]; ok {
		return h, nil
	}
	return nil, errors.New("counter not found")
}

// IncCounter increments a custom counter by a given value.
// It is a convenience wrapper and may be slightly slower than using a CounterHandle.
func (c *Collector) IncCounter(name string, val int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if h, ok := c.counters[name]; ok {
		h.Add(val)
	}
}

// --- Snapshots ---

// GlobalSnapshot returns a snapshot of all aggregated metrics since initialization.
func (c *Collector) GlobalSnapshot() *GlobalSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return &GlobalSnapshot{
		Total:          c.total.Load(),
		Errors:         c.errors.Load(),
		Latencies:      c.latencies.Snapshot(),
		QueueWaits:     c.queueWaits.Snapshot(),
		LatPercentiles: c.latPercentiles,
		QWPercentiles:  c.qwPercentiles,
		Exemplars:      c.exemplars.Snapshot(),
	}
}

// RangeScalesSnapshot returns snapshots for all configured range-based windows.
func (c *Collector) RangeScalesSnapshot() map[string]*ScaleSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshots := make(map[string]*ScaleSnapshot, len(c.rangeScales))
	for _, rs := range c.rangeScales {
		snapshots[rs.name] = rs.Snapshot()
	}
	return snapshots
}

// TimeScalesSnapshot returns snapshots for all configured time-based windows.
func (c *Collector) TimeScalesSnapshot() map[string]*ScaleSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshots := make(map[string]*ScaleSnapshot, len(c.timeScales))
	for _, ts := range c.timeScales {
		snapshots[ts.name] = ts.Snapshot()
	}
	return snapshots
}

// --- Merging ---

// MergeCollectors combines multiple collectors into a new one.
// It returns an error if the collectors have incompatible configurations.
func MergeCollectors(collectors ...*Collector) (*Collector, error) {
	if len(collectors) == 0 {
		return nil, nil
	}
	if len(collectors) == 1 {
		return collectors[0], nil
	}

	base := collectors[0]
	for _, c := range collectors[1:] {
		if c.hcfg != base.hcfg || c.ecfg != base.ecfg {
			return nil, errors.New("incompatible collectors")
		}
	}

	// Use the configuration of the first collector as the base.
	baseCfg := collectors[0]
	merged := NewCollector(
		WithLatencyHistogram(
			WithHistogramMin(baseCfg.hcfg.min),
			WithHistogramMax(baseCfg.hcfg.max),
			WithHistogramGrowth(baseCfg.hcfg.growth),
			WithHistogramNumBuckets(baseCfg.hcfg.numBuckets),
			WithHistogramNumShards(baseCfg.hcfg.numShards),
		),
		WithQueueWaitHistogram(
			WithHistogramMin(baseCfg.hcfg.min),
			WithHistogramMax(baseCfg.hcfg.max),
			WithHistogramGrowth(baseCfg.hcfg.growth),
			WithHistogramNumBuckets(baseCfg.hcfg.numBuckets),
			WithHistogramNumShards(baseCfg.hcfg.numShards),
		),
		WithExemplar(
			WithExemplarCapacity(baseCfg.ecfg.capacity),
		),
	)

	for _, c := range collectors {
		merged.total.Add(c.total.Load())
		merged.errors.Add(c.errors.Load())
		if err := merged.latencies.Merge(c.latencies); err != nil {
			return nil, err
		}
		if err := merged.queueWaits.Merge(c.queueWaits); err != nil {
			return nil, err
		}
		if err := merged.latPercentiles.Merge(c.latPercentiles); err != nil {
			return nil, err
		}
		if err := merged.qwPercentiles.Merge(c.qwPercentiles); err != nil {
			return nil, err
		}
		for _, ex := range c.exemplars.Snapshot() {
			merged.exemplars.Offer(ex.latency, ex.requestID)
		}
		for name, h := range c.counters {
			if mh, ok := merged.counters[name]; ok {
				mh.value.Add(h.value.Load())
			} else {
				merged.counters[name] = h
			}
		}
		for i, rs := range c.rangeScales {
			if i < len(merged.rangeScales) {
				if err := merged.rangeScales[i].Scale.Merge(rs.Scale); err != nil {
					return nil, err
				}
			} else {
				merged.rangeScales = append(merged.rangeScales, rs)
			}
		}
		for i, ts := range c.timeScales {
			if i < len(merged.timeScales) {
				if err := merged.timeScales[i].Scale.Merge(ts.Scale); err != nil {
					return nil, err
				}
			} else {
				merged.timeScales = append(merged.timeScales, ts)
			}
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
		merged.Latencies.Merge(s.Latencies)
		merged.QueueWaits.Merge(s.QueueWaits)
		if err := merged.LatPercentiles.Merge(s.LatPercentiles); err != nil {
			return nil, err
		}
		if err := merged.QWPercentiles.Merge(s.QWPercentiles); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

// --- Data Structures for Snapshots ---

// GlobalSnapshot contains the aggregated metrics for all requests.
type GlobalSnapshot struct {
	Total          uint64             `json:"total"`
	Errors         uint64             `json:"errors"`
	Latencies      *HistogramSnapshot `json:"latencies"`
	QueueWaits     *HistogramSnapshot `json:"queue_waits"`
	LatPercentiles QuantileSketch     `json:"lat_percentiles"`
	QWPercentiles  QuantileSketch     `json:"qw_percentiles"`
	Exemplars      []*item            `json:"exemplars"`
	SchemaVersion  string             `json:"schema_version"`
	BoundsCRC32    uint32             `json:"bounds_crc32"`
	SketchKind     string             `json:"sketch_kind"`
	InvariantsOK   bool               `json:"invariants_ok"`
}

// MarshalJSON implements the json.Marshaler interface.
func (s *GlobalSnapshot) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
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

// --- Interfaces and supporting types ---

// QuantileSketch defines the interface for approximate percentile estimators like t-digest.
type QuantileSketch interface {
	Add(value float64)
	Quantile(q float64) float64
	Merge(other QuantileSketch) error
}
