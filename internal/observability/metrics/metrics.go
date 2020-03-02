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

	NewKey = tag.NewKey

	ValdOrg = "vald.vdaas.org"
)

type Measurement = stats.Measurement
type View = view.View

type Int64Measure = stats.Int64Measure
type Float64Measure = stats.Float64Measure

type Key = tag.Key

type MeasurementWithTags struct {
	Measurement Measurement
	Tags        map[Key]string
}

type Metric interface {
	MeasurementsCount() int
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

func IsMeasureType(t reflect.Type) bool {
	switch t.Name() {
	case "Int64Measure":
		return true
	case "Float64Measure":
		return true
	}
	return false
}
