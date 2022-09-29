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
package metrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const ValdOrg = "vald.vdaas.org"

var (
	RoughMillisecondsDistribution = []float64{
		1,
		5,
		10,
		30,
		60,
		8,
		100,
		200,
		300,
		400,
		500,
		600,
		800,
		1000,
		1300,
		1600,
		2000,
		2500,
		3000,
		4000,
		5000,
		6500,
		8000,
		10000,
		13000,
		16000,
		20000,
		25000,
		30000,
		40000,
		50000,
		65000,
		80000,
		100000,
		200000,
		500000,
		1000000,
		2000000,
		5000000,
		10000000,
	}
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

type View = view.View

type Viewer interface {
	View() ([]*View, error)
}

// Metric represents an interface for metric.
type Metric interface {
	Viewer
	Register(Meter) error
}
