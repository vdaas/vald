package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
)

const (
	ValdOrg             = "vald.vdaas.org"
	InstrumentationName = "vdaas/vald"
)

// Meter is type alias of metrics.Meter.
type Meter = metric.Meter

// GetMeter returns the Meter object to record metrics.
func GetMeter() Meter {
	return global.MeterProvider().Meter(InstrumentationName)
}

// Unit is type alias of unit.Unit.
type Unit = unit.Unit

// Units defined by OpenTelemetry.
const (
	Dimensionless = unit.Dimensionless
	Bytes         = unit.Bytes
	Milliseconds  = unit.Milliseconds
)

type (
// String = attribute.Bool
)

// Metric represents an interface for metric.
type Metric interface {
	Register(metric.Meter) error
}
