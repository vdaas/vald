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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/atomic"
)

type (
	// Option represents a functional option for configuring the metrics collector.
	Option func(*collector) error

	// HistogramOption represents a functional option for configuring a histogram.
	HistogramOption func(*histogramConfig) error

	// TDigestOption configures a TDigest.
	TDigestOption func(*tdigestConfig) error

	// ExemplarOption represents a functional option for configuring an Exemplar.
	ExemplarOption func(*exemplarConfig)
)

var (
	defaultTDigestOpts = []TDigestOption{
		WithTDigestCompression(100),
		WithTDigestCompressionTriggerFactor(1.5),
		WithTDigestQuantiles([]float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99}...),
	}

	defaultHistogramOpts = []HistogramOption{
		WithHistogramNumShards(16),
	}

	defaultExemplarOpts = []ExemplarOption{
		WithExemplarCapacity(10),
		WithExemplarNumShards(8),
	}

	defaultOptions = []Option{
		WithLatencyHistogram(defaultHistogramOpts...),
		WithQueueWaitHistogram(defaultHistogramOpts...),
		WithExemplar(defaultExemplarOpts...),
		WithLatencyTDigest(defaultTDigestOpts...),
		WithQueueWaitTDigest(defaultTDigestOpts...),
	}
)

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
func WithTimeScale(name string, width time.Duration, capacity uint64) Option {
	return func(c *collector) error {
		if width <= 0 {
			return errors.New("time scale width must be positive")
		}
		// width is checked to be positive, so casting to uint64 is safe
		s, err := newScale(name, uint64(width), capacity, len(c.counters), TimeScale, c.latencies, c.queueWaits, c.exemplars) //nolint:gosec // width is checked to be positive
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
		s, err := newScale(name, width, capacity, len(c.counters), RangeScale, c.latencies, c.queueWaits, c.exemplars)
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
		h, err := NewHistogram(opts...)
		if err != nil {
			return err
		}
		c.latencies = h
		return nil
	}
}

// WithQueueWaitHistogram sets the histogram for queue wait metrics.
func WithQueueWaitHistogram(opts ...HistogramOption) Option {
	return func(c *collector) error {
		h, err := NewHistogram(opts...)
		if err != nil {
			return err
		}
		c.queueWaits = h
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
		c.latPercentiles = t
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
		c.qwPercentiles = t
		return nil
	}
}

// WithExemplar sets the exemplar for the collector.
func WithExemplar(opts ...ExemplarOption) Option {
	return func(c *collector) error {
		e := NewExemplar(opts...)
		c.exemplars = e
		return nil
	}
}

// WithBucketInterval is deprecated and ignored.
func WithBucketInterval(interval time.Duration) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithTailSegments is deprecated and ignored.
func WithTailSegments(count int) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithHistogramMaxBuckets is deprecated and ignored.
func WithHistogramMaxBuckets(count int) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithHistogramNumShards sets the number of shards for the histogram.
func WithHistogramNumShards(n int) HistogramOption {
	return func(cfg *histogramConfig) error {
		cfg.NumShards = n
		if cfg.NumShards <= 0 {
			return errors.New("numShards must be positive")
		}
		return nil
	}
}

// Deprecated options maintained for compatibility but no-op or mapped if possible.

// WithHistogramMin is deprecated.
func WithHistogramMin(minVal float64) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithHistogramMax is deprecated.
func WithHistogramMax(maxVal float64) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithHistogramGrowth is deprecated.
func WithHistogramGrowth(growth float64) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithHistogramNumBuckets is deprecated.
func WithHistogramNumBuckets(n int) HistogramOption {
	return func(cfg *histogramConfig) error {
		return nil
	}
}

// WithTDigestCompression sets the compression for the t-digest.
func WithTDigestCompression(c float64) TDigestOption {
	return func(cfg *tdigestConfig) error {
		cfg.Compression = c
		if cfg.Compression <= 0 {
			return errors.New("tdigest compression must be > 0")
		}
		return nil
	}
}

// WithTDigestCompressionTriggerFactor sets the compression trigger factor for the t-digest.
func WithTDigestCompressionTriggerFactor(f float64) TDigestOption {
	return func(cfg *tdigestConfig) error {
		cfg.CompressionTriggerFactor = f
		if cfg.CompressionTriggerFactor <= 0 {
			return errors.New("tdigest compressionTriggerFactor must be > 0")
		}
		return nil
	}
}

// WithTDigestQuantiles sets the quantiles to be used in the String() method.
func WithTDigestQuantiles(quantiles ...float64) TDigestOption {
	return func(cfg *tdigestConfig) error {
		if len(quantiles) > 0 {
			cfg.Quantiles = quantiles
		}
		return nil
	}
}

// WithTDigestNumShards sets the number of shards for the t-digest.
func WithTDigestNumShards(n int) TDigestOption {
	return func(cfg *tdigestConfig) error {
		if n < 1 {
			return errors.New("tdigest num shards must be >= 1")
		}
		cfg.NumShards = n
		return nil
	}
}

// WithExemplarCapacity sets the capacity for the exemplar.
func WithExemplarCapacity(k int) ExemplarOption {
	return func(cfg *exemplarConfig) {
		if k >= 1 {
			cfg.Capacity = k
		}
	}
}

// WithExemplarNumShards sets the number of shards for the exemplar.
func WithExemplarNumShards(n int) ExemplarOption {
	return func(cfg *exemplarConfig) {
		if n >= 1 {
			cfg.NumShards = n
		}
	}
}
