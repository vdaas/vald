//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/metrics/cpu"
	"github.com/vdaas/vald/internal/observability/metrics/mem"
	"github.com/vdaas/vald/internal/observability/metrics/runtime"
	"github.com/vdaas/vald/internal/observability/metrics/version"
	"github.com/vdaas/vald/internal/timeutil"
)

type CollectorOption func(*collector) error

var (
	collectorDefaultOpts = []CollectorOption{
		WithErrGroup(errgroup.Get()),
		WithDuration("5s"),
	}
)

func WithErrGroup(eg errgroup.Group) CollectorOption {
	return func(c *collector) error {
		if eg != nil {
			c.eg = eg
		}
		return nil
	}
}

func WithDuration(dur string) CollectorOption {
	return func(c *collector) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		c.duration = d
		return nil
	}
}

func WithMetrics(metrics ...metrics.Metric) CollectorOption {
	return func(c *collector) error {
		if metrics == nil {
			return nil
		}
		if c.metrics != nil && len(c.metrics) > 0 {
			c.metrics = append(c.metrics, metrics...)
		} else {
			c.metrics = metrics
		}
		return nil
	}
}

func WithVersionInfo(enabled bool) CollectorOption {
	return func(c *collector) error {
		if !enabled {
			return nil
		}
		versionInfo, err := version.NewMetric()
		if err != nil {
			return err
		}
		return WithMetrics(versionInfo)(c)
	}
}

func WithCPUMetrics(enabled bool) CollectorOption {
	return func(c *collector) error {
		if !enabled {
			return nil
		}
		cpu, err := cpu.NewMetric()
		if err != nil {
			return err
		}
		return WithMetrics(cpu)(c)
	}
}

func WithMemoryMetrics(enabled bool) CollectorOption {
	return func(c *collector) error {
		if !enabled {
			return nil
		}
		return WithMetrics(mem.NewMetric())(c)
	}
}

func WithGoroutineMetrics(enabled bool) CollectorOption {
	return func(c *collector) error {
		if !enabled {
			return nil
		}
		return WithMetrics(runtime.NewGoroutineMetrics())(c)
	}
}

func WithCGOMetrics(enabled bool) CollectorOption {
	return func(c *collector) error {
		if !enabled {
			return nil
		}
		return WithMetrics(runtime.NewCGOMetrics())(c)
	}
}
