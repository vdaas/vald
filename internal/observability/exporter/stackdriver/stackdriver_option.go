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

// Package stackdriver provides a stackdriver exporter.
package stackdriver

import (
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/client/google"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(e *exp) error

var (
	defaultOpts = []Option{
		WithOnErrorFunc(func(err error) {
			if err != nil {
				log.Warnf("Error when uploading stats or spans to Stackdriver: %v", err)
			}
		}),
		WithMonitoredResource(monitoredresource.Autodetect()),
		WithMetricPrefix("vald.vdaas.org/"),
		WithTimeout("5s"),
		WithReportingInterval("0"),
		WithNumberOfWorkers(1),
	}
)

func WithMonitoring(enabled bool) Option {
	return func(e *exp) error {
		e.monitoringEnabled = enabled

		return nil
	}
}

func WithTracing(enabled bool) Option {
	return func(e *exp) error {
		e.tracingEnabled = enabled

		return nil
	}
}

func WithProjectID(pid string) Option {
	return func(e *exp) error {
		if pid == "" {
			return nil
		}

		e.ProjectID = pid

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithProjectID(pid),
			}
		} else {
			e.topts = append(e.topts, trace.WithProjectID(pid))
		}

		return nil
	}
}

func WithLocation(loc string) Option {
	return func(e *exp) error {
		if loc != "" {
			e.Location = loc
		}

		return nil
	}
}

func WithOnErrorFunc(f func(error)) Option {
	return func(e *exp) error {
		if f == nil {
			return nil
		}

		e.OnError = f

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithOnError(f),
			}
		} else {
			e.topts = append(e.topts, trace.WithOnError(f))
		}

		return nil
	}
}

func WithMonitoringClientOptions(copts ...google.Option) Option {
	return func(e *exp) error {
		opts := make([]google.Option, 0, len(copts))
		for _, opt := range copts {
			if opt != nil {
				opts = append(opts, opt)
			}
		}

		if e.MonitoringClientOptions == nil {
			e.MonitoringClientOptions = opts
			return nil
		}

		e.MonitoringClientOptions = append(e.MonitoringClientOptions, opts...)

		return nil
	}
}

func WithTraceClientOptions(copts ...google.Option) Option {
	return func(e *exp) error {
		opts := make([]google.Option, 0, len(copts))
		for _, opt := range copts {
			if opt != nil {
				opts = append(opts, opt)
			}
		}

		if e.TraceClientOptions == nil {
			e.TraceClientOptions = opts
			return nil
		}

		e.TraceClientOptions = append(e.TraceClientOptions, opts...)

		return nil
	}
}

func WithBundleDelayThreshold(dur string) Option {
	return func(e *exp) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		e.BundleDelayThreshold = d

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithBundleDelayThreshold(d),
			}
		} else {
			e.topts = append(e.topts, trace.WithBundleDelayThreshold(d))
		}

		return nil
	}
}

func WithBundleCountThreshold(cnt int) Option {
	return func(e *exp) error {
		e.BundleCountThreshold = cnt

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithBundleCountThreshold(cnt),
			}
		} else {
			e.topts = append(e.topts, trace.WithBundleCountThreshold(cnt))
		}

		return nil
	}
}

func WithTraceSpansBufferMaxBytes(bs int) Option {
	return func(e *exp) error {
		e.TraceSpansBufferMaxBytes = bs

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithBufferMaxBytes(bs),
			}
		} else {
			e.topts = append(e.topts, trace.WithBufferMaxBytes(bs))
		}

		return nil
	}
}

func WithMonitoredResource(mr monitoredresource.Interface) Option {
	return func(e *exp) error {
		if mr != nil {
			e.MonitoredResource = mr
		}

		return nil
	}
}

func WithMetricPrefix(prefix string) Option {
	return func(e *exp) error {
		if prefix != "" {
			e.MetricPrefix = prefix
		}

		return nil
	}
}

func WithGetMetricDisplayName(f func(view *metrics.View) string) Option {
	return func(e *exp) error {
		if f != nil {
			e.GetMetricDisplayName = f
		}

		return nil
	}
}

func WithGetMetricPrefix(f func(name string) string) Option {
	return func(e *exp) error {
		if f != nil {
			e.GetMetricPrefix = f
		}

		return nil
	}
}

func WithDefaultMonitoringLabels(lbs *stackdriver.Labels) Option {
	return func(e *exp) error {
		if lbs != nil {
			e.DefaultMonitoringLabels = lbs
		}
		return nil
	}
}

func WithSkipCMD(skip bool) Option {
	return func(e *exp) error {
		e.SkipCMD = skip

		return nil
	}
}

func WithTimeout(dur string) Option {
	return func(e *exp) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		e.Timeout = d

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithTimeout(d),
			}
		} else {
			e.topts = append(e.topts, trace.WithTimeout(d))
		}

		return nil
	}
}

func WithReportingInterval(dur string) Option {
	return func(e *exp) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		e.ReportingInterval = d

		return nil
	}
}

func WithNumberOfWorkers(n int) Option {
	return func(e *exp) error {
		e.NumberOfWorkers = n

		if e.topts == nil {
			e.topts = []trace.Option{
				trace.WithMaxNumberOfWorkers(n),
			}
		} else {
			e.topts = append(e.topts, trace.WithMaxNumberOfWorkers(n))
		}

		return nil
	}
}
