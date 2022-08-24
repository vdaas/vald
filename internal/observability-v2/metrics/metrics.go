package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
)

const (
	ValdOrg = "vald.vdaas.org"
)

// Meter is type alias of metrics.Meter.
type Meter = metric.Meter

// GetMeter returns the Meter object to record metrics.
func GetMeter() Meter {
	return global.MeterProvider().Meter(ValdOrg)
}

// Unit is type alias of unit.Unit.
type Unit = unit.Unit

// Units defined by OpenTelemetry.
const (
	// Dimensionless is a type alias of unit.Dimensionless.
	Dimensionless = unit.Dimensionless
	// Bytes is a type alias of unit.Bytes.
	Bytes = unit.Bytes
	// Milliseconds is a type alias of unit.Milliseconds.
	Milliseconds = unit.Milliseconds
)

type (
	// AsynchronousInstrument is type alias of instrument.Asynchronous.
	AsynchronousInstrument = instrument.Asynchronous
	// SynchronousInstrument is type alias of instrument.Synchronous.
	SynchronousInstrument = instrument.Synchronous
)

// WithUnit returns an instrument.WithUnit option.
func WithUnit(u Unit) instrument.Option {
	return instrument.WithUnit(u)
}

// WithDescription returns an instrument.WithDescription option.
func WithDescription(desc string) instrument.Option {
	return instrument.WithDescription(desc)
}

// Metric represents an interface for metric.
type Metric interface {
	Register(metric.Meter) error
}
