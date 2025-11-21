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
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
)

type (
	// Option represents a functional option for configuring the metrics collector.
	Option func(*collector) error

	// HistogramOption represents a functional option for configuring a histogram.
	HistogramOption func(*histogram) error

	// TDigestOption configures a TDigest.
	TDigestOption func(*tdigest) error

	// ExemplarOption represents a functional option for configuring an Exemplar.
	ExemplarOption func(*exemplar)
)

var (
	defaultTDigestOpts = []TDigestOption{
		WithTDigestCompression(100),
		WithTDigestCompressionTriggerFactor(1.5),
		WithTDigestQuantiles([]float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99}...),
	}

	defaultHistogramOpts = []HistogramOption{
		WithHistogramMin(1000),        // 1us
		WithHistogramMax(60000000000), // 60s (soft limit)
		WithHistogramGrowth(1.2),
		WithHistogramNumBuckets(100), // Covers range from us to min
		WithHistogramNumShards(16),
	}

	defaultExemplarOpts = []ExemplarOption{
		WithExemplarCapacity(10),
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
func WithTimeScale(name string, width, capacity uint64) Option {
	return func(c *collector) error {
		s, err := newScale(name, width, capacity, len(c.counters), TimeScale, c.latencies, c.queueWaits, c.exemplars)
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

// WithHistogramMin sets the minimum value for the histogram.
func WithHistogramMin(min float64) HistogramOption {
	return func(c *histogram) error {
		c.min = min
		if c.min <= 0 {
			return errors.New("histogram min must be > 0 for geometric buckets")
		}
		return nil
	}
}

// WithHistogramMax sets the maximum value for the histogram.
func WithHistogramMax(max float64) HistogramOption {
	return func(c *histogram) error {
		c.max = max
		return nil
	}
}

// WithHistogramGrowth sets the growth factor for the histogram.
func WithHistogramGrowth(growth float64) HistogramOption {
	return func(c *histogram) error {
		c.growth = growth
		if c.growth <= 1 {
			return errors.New("histogram growth must be > 1 for geometric buckets")
		}
		return nil
	}
}

// WithHistogramNumBuckets sets the number of buckets for the histogram.
func WithHistogramNumBuckets(n int) HistogramOption {
	return func(c *histogram) error {
		c.numBuckets = n
		if c.numBuckets < 2 {
			return errors.New("numBuckets must be at least 2")
		}
		return nil
	}
}

// WithHistogramNumShards sets the number of shards for the histogram.
func WithHistogramNumShards(n int) HistogramOption {
	return func(c *histogram) error {
		c.numShards = n
		if c.numShards <= 0 {
			return errors.New("numShards must be positive")
		}
		return nil
	}
}

// WithTDigestCompression sets the compression for the t-digest.
func WithTDigestCompression(c float64) TDigestOption {
	return func(t *tdigest) error {
		t.compression = c
		if t.compression <= 0 {
			return errors.New("tdigest compression must be > 0")
		}
		return nil
	}
}

// WithTDigestCompressionTriggerFactor sets the compression trigger factor for the t-digest.
func WithTDigestCompressionTriggerFactor(f float64) TDigestOption {
	return func(t *tdigest) error {
		t.compressionTriggerFactor = f
		if t.compressionTriggerFactor <= 0 {
			return errors.New("tdigest compressionTriggerFactor must be > 0")
		}
		return nil
	}
}

// WithQuantiles sets the quantiles to be used in the String() method.
func WithTDigestQuantiles(quantiles ...float64) TDigestOption {
	return func(t *tdigest) error {
		if len(quantiles) > 0 {
			t.quantiles = quantiles
		}
		return nil
	}
}

// WithExemplarCapacity sets the capacity for the exemplar.
func WithExemplarCapacity(k int) ExemplarOption {
	return func(e *exemplar) {
		if k >= 1 {
			e.k = k
		}
	}
}
