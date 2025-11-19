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
type Option func(*collector) error

// HistogramOption represents a functional option for configuring a histogram.
type HistogramOption func(*histogramConfig)

type histogramConfig struct {
	min        float64
	max        float64
	growth     float64
	numBuckets int
	numShards  int
}

var defaultHistogramConfig = histogramConfig{
	min:        1,
	max:        5000,
	growth:     1.2,
	numBuckets: 50,
	numShards:  16,
}

var defaultOptions = []Option{
	WithLatencyHistogram(
		WithHistogramMin(1),
		WithHistogramMax(5000),
		WithHistogramGrowth(1.2),
		WithHistogramNumBuckets(50),
		WithHistogramNumShards(16),
	),
	WithQueueWaitHistogram(
		WithHistogramMin(1),
		WithHistogramMax(5000),
		WithHistogramGrowth(1.2),
		WithHistogramNumBuckets(50),
		WithHistogramNumShards(16),
	),
	WithExemplar(
		WithCapacity(10),
	),
	WithLatencyTDigest(
		WithTDigestCompression(100),
		WithTDigestCompressionTriggerFactor(1.5),
		WithQuantiles(0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99),
	),
	WithQueueWaitTDigest(
		WithTDigestCompression(100),
		WithTDigestCompressionTriggerFactor(1.5),
		WithQuantiles(0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99),
	),
}
type exemplarConfig struct {
	capacity int
}

var defaultExemplarConfig = exemplarConfig{
	capacity: 10,
}

const (
	defaultTDigestCompression              = 100
	defaultTDigestCompressionTriggerFactor = 1.5
)

type tdigestConfig struct {
	compression              float64
	compressionTriggerFactor float64
	quantiles                []float64
}

var defaultTDigestConfig = tdigestConfig{
	compression:              defaultTDigestCompression,
	compressionTriggerFactor: defaultTDigestCompressionTriggerFactor,
	quantiles:                []float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99},
}

// WithCustomCounters registers custom counters with the collector.
func WithCustomCounters(names ...string) Option {
	return func(c *collector) error {
		for _, name := range names {
			c.counters[name] = new(CounterHandle)
			c.counters[name].value = new(atomic.Uint64)
		}
		return nil
	}
}

// WithTimeScale adds a time-based scale to the collector.
func WithTimeScale(name string, width, capacity uint64) Option {
	return func(c *collector) error {
		s, err := newScale(name, width, capacity, len(c.counters), c.hcfg, c.ecfg, timeScale, c.histogramPool)
		if err != nil {
			return err
		}
		c.scales = append(c.scales, s)
		return nil
	}
}

// WithRangeScale adds a range-based scale to the collector.
func WithRangeScale(name string, width, capacity uint64) Option {
	return func(c *collector) error {
		s, err := newScale(name, width, capacity, len(c.counters), c.hcfg, c.ecfg, rangeScale, c.histogramPool)
		if err != nil {
			return err
		}
		c.scales = append(c.scales, s)
		return nil
	}
}

// WithLatencyHistogram sets the histogram for latency metrics.
func WithLatencyHistogram(opts ...HistogramOption) Option {
	return func(c *collector) error {
		for _, opt := range opts {
			opt(&c.hcfg)
		}
		h, err := NewHistogram(opts...)
		if err != nil {
			return err
		}
		c.global.latencies = h
		return nil
	}
}

// WithQueueWaitHistogram sets the histogram for queue wait metrics.
func WithQueueWaitHistogram(opts ...HistogramOption) Option {
	return func(c *collector) error {
		for _, opt := range opts {
			opt(&c.hcfg)
		}
		h, err := NewHistogram(opts...)
		if err != nil {
			return err
		}
		c.global.queueWaits = h
		return nil
	}
}

// WithLatencyTDigest sets the t-digest for latency metrics.
func WithLatencyTDigest(opts ...TDigestOption) Option {
	return func(c *collector) error {
		t, err := NewTDigest(opts...)
		if err != nil {
			return err
		}
		c.global.latPercentiles = t
		return nil
	}
}

// WithQueueWaitTDigest sets the t-digest for queue wait metrics.
func WithQueueWaitTDigest(opts ...TDigestOption) Option {
	return func(c *collector) error {
		t, err := NewTDigest(opts...)
		if err != nil {
			return err
		}
		c.global.qwPercentiles = t
		return nil
	}
}

// WithExemplar sets the exemplar for the collector.
func WithExemplar(opts ...ExemplarOption) Option {
	return func(c *collector) error {
		e := NewExemplar(opts...)
		c.global.exemplars = e
		return nil
	}
}

// WithHistogramMin sets the minimum value for the histogram.
func WithHistogramMin(min float64) HistogramOption {
	return func(c *histogramConfig) {
		c.min = min
	}
}

// WithHistogramMax sets the maximum value for the histogram.
func WithHistogramMax(max float64) HistogramOption {
	return func(c *histogramConfig) {
		c.max = max
	}
}

// WithHistogramGrowth sets the growth factor for the histogram.
func WithHistogramGrowth(growth float64) HistogramOption {
	return func(c *histogramConfig) {
		c.growth = growth
	}
}

// WithHistogramNumBuckets sets the number of buckets for the histogram.
func WithHistogramNumBuckets(n int) HistogramOption {
	return func(c *histogramConfig) {
		c.numBuckets = n
	}
}

// WithHistogramNumShards sets the number of shards for the histogram.
func WithHistogramNumShards(n int) HistogramOption {
	return func(c *histogramConfig) {
		c.numShards = n
	}
}

// WithExemplarCapacity sets the capacity for the exemplar.
func WithExemplarCapacity(k int) ExemplarOption {
	return func(e *exemplar) {
		e.k = k
	}
}

// WithTDigestCompression sets the compression for the t-digest.
func WithTDigestCompression(c float64) TDigestOption {
	return func(t *TDigest) {
		t.compression = c
	}
}

// WithTDigestCompressionTriggerFactor sets the compression trigger factor for the t-digest.
func WithTDigestCompressionTriggerFactor(f float64) TDigestOption {
	return func(t *TDigest) {
		t.compressionTriggerFactor = f
	}
}

// WithTDigestQuantiles sets the quantiles for the t-digest.
func WithTDigestQuantiles(q ...float64) TDigestOption {
	return func(t *TDigest) {
		t.quantiles = q
	}
}
