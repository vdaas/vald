// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package prometheus

import (
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type Option func(e *exp) error

var defaultOpts = []Option{
	WithEndpoint("/metrics"),
	WithNamespace("vald"),
	WithCollectInterval("500ms"),
	WithCollectTimeout("10s"),
	WithInMemoty(true),
	WithHistogramDistribution(
		exporter.DefaultMillisecondsHistogramDistribution,
	),
}

func WithEndpoint(ep string) Option {
	return func(e *exp) error {
		if ep == "" {
			return errors.NewErrInvalidOption("endpoint", ep)
		}
		e.endpoint = ep
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(e *exp) error {
		if ns == "" {
			return errors.NewErrInvalidOption("namespace", ns)
		}
		e.namespace = ns
		return nil
	}
}

func WithCollectInterval(period string) Option {
	return func(e *exp) error {
		if len(period) == 0 {
			return errors.NewErrInvalidOption("collectInterval", period)
		}

		dur, err := time.ParseDuration(period)
		if err != nil {
			return errors.NewErrInvalidOption("collectInterval", period, err)
		}
		e.collectInterval = dur
		return nil
	}
}

func WithCollectTimeout(timeout string) Option {
	return func(e *exp) error {
		if len(timeout) == 0 {
			return errors.NewErrInvalidOption("collectTimeout", timeout)
		}

		dur, err := time.ParseDuration(timeout)
		if err != nil {
			return errors.NewErrInvalidOption("collectTimeout", timeout, err)
		}
		e.collectTimeout = dur
		return nil
	}
}

func WithInMemoty(ok bool) Option {
	return func(e *exp) error {
		e.inmemoryEnabled = ok
		return nil
	}
}

func WithView(views ...metrics.Viewer) Option {
	return func(e *exp) error {
		return nil
	}
}

func WithHistogramDistribution(fs []float64) Option {
	return func(e *exp) error {
		if len(fs) == 0 {
			return errors.NewErrInvalidOption("histogramBoundarie", fs)
		}
		e.histogramBoundarie = fs
		return nil
	}
}
