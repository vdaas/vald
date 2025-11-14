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

import "sync/atomic"

// Option represents a functional option for configuring the metrics collector.
type Option func(*Collector) error

// HistogramOption represents a functional option for configuring a histogram.
type HistogramOption func(*histogramConfig)

// TDigestOption represents a functional option for configuring a TDigest.
type TDigestOption func(*tdigestConfig)

// ExemplarOption represents a functional option for configuring an exemplar.
type ExemplarOption func(*exemplarConfig)

type (
	histogramConfig struct {
		min        float64
		max        float64
		growth     float64
		numBuckets int
		numShards  int
	}
	tdigestConfig struct {
		compression              float64
		compressionTriggerFactor float64
	}
	exemplarConfig struct {
		capacity int
	}
)

var (
	defaultHistogramConfig = histogramConfig{
		min:        1, // >0 to satisfy NewHistogram’s geometric bucket requirement
		max:        1e9,
		growth:     1.6,
		numBuckets: 30,
		numShards:  16,
	}
	defaultTDigestConfig = tdigestConfig{
		compression:              100,
		compressionTriggerFactor: 10,
	}
	defaultExemplarConfig = exemplarConfig{
		capacity: 10,
	}
)

// WithTDigestCompressionTriggerFactor sets the compression trigger factor for the TDigest.
func WithTDigestCompressionTriggerFactor(factor float64) TDigestOption {
	return func(t *tdigestConfig) {
		t.compressionTriggerFactor = factor
	}
}

// WithHistogramMin sets the minimum value for the histogram.
func WithHistogramMin(min float64) HistogramOption {
	return func(h *histogramConfig) {
		h.min = min
	}
}

// WithHistogramMax sets the maximum value for the histogram.
func WithHistogramMax(max float64) HistogramOption {
	return func(h *histogramConfig) {
		h.max = max
	}
}

// WithHistogramGrowth sets the growth factor for the histogram.
func WithHistogramGrowth(growth float64) HistogramOption {
	return func(h *histogramConfig) {
		h.growth = growth
	}
}

// WithHistogramNumBuckets sets the number of buckets for the histogram.
func WithHistogramNumBuckets(numBuckets int) HistogramOption {
	return func(h *histogramConfig) {
		h.numBuckets = numBuckets
	}
}

// WithHistogramNumShards sets the number of shards for the histogram.
func WithHistogramNumShards(numShards int) HistogramOption {
	return func(h *histogramConfig) {
		h.numShards = numShards
	}
}

// WithTDigestCompression sets the compression for the TDigest.
func WithTDigestCompression(compression float64) TDigestOption {
	return func(t *tdigestConfig) {
		t.compression = compression
	}
}

// WithExemplarCapacity sets the capacity for the exemplar.
func WithExemplarCapacity(capacity int) ExemplarOption {
	return func(e *exemplarConfig) {
		e.capacity = capacity
	}
}

// WithLatencyHistogram returns an option to set the latency histogram configuration.
func WithLatencyHistogram(opts ...HistogramOption) Option {
	return func(c *Collector) error {
		h, err := NewHistogram(opts...)
		if err != nil {
			return err
		}
		c.latencies = h
		for _, opt := range opts {
			opt(&c.hcfg)
		}
		return nil
	}
}

// WithQueueWaitHistogram returns an option to set the queue wait histogram configuration.
func WithQueueWaitHistogram(opts ...HistogramOption) Option {
	return func(c *Collector) error {
		h, err := NewHistogram(opts...)
		if err != nil {
			return err
		}
		c.queueWaits = h
		for _, opt := range opts {
			opt(&c.hcfg)
		}
		return nil
	}
}

// WithLatencyTDigest returns an option to set the latency TDigest configuration.
func WithLatencyTDigest(opts ...TDigestOption) Option {
	return func(c *Collector) error {
		cfg := defaultTDigestConfig
		for _, opt := range opts {
			opt(&cfg)
		}
		var err error
		c.latPercentiles, err = NewTDigest(cfg.compression, cfg.compressionTriggerFactor)
		return err
	}
}

// WithQueueWaitTDigest returns an option to set the queue wait TDigest configuration.
func WithQueueWaitTDigest(opts ...TDigestOption) Option {
	return func(c *Collector) error {
		cfg := defaultTDigestConfig
		for _, opt := range opts {
			opt(&cfg)
		}
		var err error
		c.qwPercentiles, err = NewTDigest(cfg.compression, cfg.compressionTriggerFactor)
		return err
	}
}

// WithExemplar returns an option to set the exemplar configuration.
func WithExemplar(opts ...ExemplarOption) Option {
	return func(c *Collector) error {
		cfg := c.ecfg
		for _, opt := range opts {
			opt(&cfg)
		}
		c.ecfg = cfg
		c.exemplars = NewExemplar(cfg.capacity)
		return nil
	}
}

// WithRangeScale is an option to add a range scale.
// It is important to register all custom counters via WithCustomCounters *before* adding any scales.
func WithRangeScale(name string, width, capacity uint64) Option {
	return func(c *Collector) error {
		rs, err := NewRangeScale(name, width, capacity, len(c.counters), c.hcfg, c.ecfg)
		if err != nil {
			return err
		}
		c.rangeScales = append(c.rangeScales, rs)
		return nil
	}
}

// WithTimeScale is an option to add a time scale.
// It is important to register all custom counters via WithCustomCounters *before* adding any scales.
func WithTimeScale(name string, widthSec, capacity uint64) Option {
	return func(c *Collector) error {
		ts, err := NewTimeScale(name, widthSec, capacity, len(c.counters), c.hcfg, c.ecfg)
		if err != nil {
			return err
		}
		c.timeScales = append(c.timeScales, ts)
		return nil
	}
}

// WithCustomCounters is an option to add custom counters.
// This option should be used *before* any WithRangeScale or WithTimeScale options
// to ensure that the scales are initialized with the correct number of counters.
func WithCustomCounters(names ...string) Option {
	return func(c *Collector) error {
		for _, name := range names {
			c.counters[name] = &CounterHandle{
				value: new(atomic.Uint64),
			}
		}
		return nil
	}
}
