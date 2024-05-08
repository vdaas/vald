// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

type (
	View = view.View

	// Meter is type alias of metrics.Meter.
	Meter = metric.Meter
	// Metric represents an interface for metric.
	Metric interface {
		View() ([]View, error)
		Register(Meter) error
	}
)

const (
	ValdOrg = "vald.vdaas.org"
	// Units defined by OpenTelemetry.
	// Dimensionless is a type alias of unit.Dimensionless.
	Dimensionless = "1"
	// Bytes is a type alias of unit.Bytes.
	Bytes = "By"
	// Milliseconds is a type alias of unit.Milliseconds.
	Milliseconds = "ms"
)

var (
	// WithUnit returns an metric.WithUnit option.
	WithUnit = metric.WithUnit
	// WithDescription returns an metric.WithDescription option.
	WithDescription = metric.WithDescription
	// WithAttributes returns an metric.WithAttributes option.
	WithAttributes = metric.WithAttributes

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

	DefaultMillisecondsDistribution = []float64{
		0.01,
		0.05,
		0.1,
		0.3,
		0.6,
		0.8,
		1,
		2,
		3,
		4,
		5,
		6,
		8,
		10,
		13,
		16,
		20,
		25,
		30,
		40,
		50,
		65,
		80,
		100,
		130,
		160,
		200,
		250,
		300,
		400,
		500,
		650,
		800,
		1000,
		2000,
		5000,
		10000,
		20000,
		50000,
		100000,
	}
)

// GetMeter returns the Meter object to record metrics.
func GetMeter() Meter {
	return otel.GetMeterProvider().Meter(ValdOrg)
}
