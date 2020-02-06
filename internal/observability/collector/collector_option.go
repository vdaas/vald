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
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/timeutil"
)

type CollectorOption func(*collector) error

var (
	collectorDefaultOpts = []CollectorOption{
		WithDuration("5s"),
	}
)

func WithDuration(dur string) CollectorOption {
	return func(c *collector) error {
		if dur != "" {
			d, err := timeutil.Parse(dur)
			if err != nil {
				return err
			}

			c.duration = d
		}
		return nil
	}
}

func WithMetrics(metrics ...metrics.Metric) CollectorOption {
	return func(c *collector) error {
		if c.metrics != nil && len(c.metrics) > 0 {
			c.metrics = append(c.metrics, metrics...)
		} else {
			c.metrics = metrics
		}
		return nil
	}
}
