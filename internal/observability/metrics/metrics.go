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

// Package metrics provides metrics.
package metrics

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	Int64   = stats.Int64
	Float64 = stats.Float64

	UnitDimensionless = stats.UnitDimensionless
	UnitBytes         = stats.UnitBytes
	UnitMilliseconds  = stats.UnitMilliseconds

	Count        = view.Count
	Distribution = view.Distribution
	LastValue    = view.LastValue
	Sum          = view.Sum

	DefaultMillisecondsDistribution = Distribution(
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
	)
	RoughMillisecondsDistribution = Distribution(
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
	)

	NewKey = tag.NewKey

	ValdOrg = "vald.vdaas.org"
)

type (
	Measurement = stats.Measurement
	View        = view.View
)

type (
	Int64Measure   = stats.Int64Measure
	Float64Measure = stats.Float64Measure
)

type Key = tag.Key

type MeasurementWithTags struct {
	Measurement Measurement
	Tags        map[Key]string
}

type Metric interface {
	Measurement(ctx context.Context) ([]Measurement, error)
	MeasurementWithTags(ctx context.Context) ([]MeasurementWithTags, error)
	View() []*View
}

func RegisterView(views ...*View) error {
	return view.Register(views...)
}

func Record(ctx context.Context, ms ...Measurement) {
	stats.Record(ctx, ms...)
}

func RecordWithTags(ctx context.Context, mwts ...MeasurementWithTags) (errs error) {
	for _, mwt := range mwts {
		mutators := make([]tag.Mutator, 0, len(mwt.Tags))
		for k, v := range mwt.Tags {
			mutators = append(mutators, tag.Upsert(k, v))
		}
		err := stats.RecordWithTags(ctx, mutators, mwt.Measurement)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}
	return errs
}

func MeasurementsCount(m Metric) int {
	cnt := 0
	rv := reflect.Indirect(reflect.ValueOf(m))
	for i := 0; i < rv.NumField(); i++ {
		switch rv.Field(i).Type() {
		case reflect.TypeOf(Int64Measure{}),
			reflect.TypeOf(&Int64Measure{}),
			reflect.TypeOf(Float64Measure{}),
			reflect.TypeOf(&Float64Measure{}):
			cnt++
		}
	}
	return cnt
}
