package otlp

import (
	"time"

	"github.com/vdaas/vald/internal/errors"
)

type Option func(*exp) error

var defaultOpts = []Option{
	WithServiceName("vald"),
	WithTraceBatchTimeout("1s"),
	WithTraceExportTimeout("1m"),
	WithTraceMaxExportBatchSize(1024),
	WithTraceMaxQueueSize(256),
	WithMetricsExportInterval("1s"),
	WithMetricsExportTimeout("1m"),
}

func WithServiceName(s string) Option {
	return func(e *exp) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("serviceName", s)
		}
		e.serviceName = s
		return nil
	}
}

func WithCollectorEndpoint(ep string) Option {
	return func(e *exp) error {
		if len(ep) == 0 {
			return errors.NewErrCriticalOption("collectorEndpoint", ep)
		}
		e.collectorEndpoint = ep
		return nil
	}
}

func WithTraceBatchTimeout(s string) Option {
	return func(e *exp) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("traceBatchTimeout", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("traceBatchTimeout", s, err)
		}
		e.tBatchTimeout = dur
		return nil
	}
}

func WithTraceExportTimeout(s string) Option {
	return func(e *exp) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("traceExportTimeout", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("traceExportTimeout", s, err)
		}
		e.tExportTimeout = dur
		return nil
	}
}

func WithTraceMaxExportBatchSize(size int) Option {
	return func(e *exp) error {
		if size < 0 {
			return errors.NewErrInvalidOption("traceMaxExportBatchSize", size)
		}
		e.tMaxExportBatchSize = size
		return nil
	}
}

func WithTraceMaxQueueSize(size int) Option {
	return func(e *exp) error {
		if size < 0 {
			return errors.NewErrInvalidOption("traceMaxQueueSize", size)
		}
		e.tMaxQueueSize = size
		return nil
	}
}

func WithMetricsExportInterval(s string) Option {
	return func(e *exp) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("metricsExportInterval", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("metricsExportInterval", s, err)
		}
		e.mExportInterval = dur
		return nil
	}
}

func WithMetricsExportTimeout(s string) Option {
	return func(e *exp) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("metricsExportTimeout", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("metricsExportTimeout", s, err)
		}
		e.mExportTimeout = dur
		return nil
	}
}
