//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package collector provides metrics collector
package collector

import (
	"context"
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/safety"
)

type Collector interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	Stop(ctx context.Context)
}

type collector struct {
	duration time.Duration
	metrics  []metrics.Metric
	eg       errgroup.Group
}

func New(opts ...CollectorOption) (Collector, error) {
	co := new(collector)

	for _, opt := range append(collectorDefaultOpts, opts...) {
		err := opt(co)
		if err != nil {
			return nil, err
		}
	}

	return co, nil
}

func (c *collector) PreStart(ctx context.Context) error {
	views := make([]*metrics.View, 0, len(c.metrics))
	for _, metric := range c.metrics {
		views = append(views, metric.View()...)
	}

	return metrics.RegisterView(views...)
}

func (c *collector) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 2)

	c.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		tick := time.NewTicker(c.duration)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tick.C:
				err := c.collect(ctx)
				if err != nil {
					runtime.Gosched()
					ech <- err
				}
			}
		}
	}))

	return ech
}

func (c *collector) Stop(ctx context.Context) {
}

func (c *collector) collect(ctx context.Context) (err error) {
	cnt := 0
	for _, metric := range c.metrics {
		measurementsCount := metrics.MeasurementsCount(metric)
		if cnt < measurementsCount {
			cnt = measurementsCount
		}
	}
	measurements := make([]metrics.Measurement, 0, len(c.metrics)*cnt)
	measurementsWithTags := make([]metrics.MeasurementWithTags, 0, len(c.metrics)*cnt)
	for _, metric := range c.metrics {
		ms, err := metric.Measurement(ctx)
		if err != nil {
			return err
		}
		measurements = append(measurements, ms...)
		mwts, err := metric.MeasurementWithTags(ctx)
		if err != nil {
			return err
		}
		measurementsWithTags = append(measurementsWithTags, mwts...)
	}

	metrics.Record(ctx, measurements...)
	return metrics.RecordWithTags(ctx, measurementsWithTags...)
}
