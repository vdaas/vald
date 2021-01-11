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

// Package mysql provides mysql metrics functions
package mysql

import (
	"context"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/trace"
)

type mysqlMetrics struct {
	queryTotal   metrics.Int64Measure
	queryLatency metrics.Float64Measure

	mu sync.Mutex
	ms []metrics.Measurement

	mysql.NullEventReceiver
}

type EventReceiver interface {
	metrics.Metric
	mysql.EventReceiver
	mysql.TracingEventReceiver
}

type startTimeKey struct{}

func New() (e EventReceiver, err error) {
	ms := new(mysqlMetrics)

	ms.queryTotal = *metrics.Int64(
		metrics.ValdOrg+"/db/rdb/mysql/completed_query_total",
		"cumulative count of completed queries",
		metrics.UnitDimensionless,
	)

	ms.queryLatency = *metrics.Float64(
		metrics.ValdOrg+"/db/rdb/mysql/query_latency",
		"query latency",
		metrics.UnitMilliseconds,
	)

	ms.ms = make([]metrics.Measurement, 0)

	return ms, nil
}

func (mm *mysqlMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	mm.mu.Lock()
	defer func() {
		mm.ms = make([]metrics.Measurement, 0)
		mm.mu.Unlock()
	}()

	return mm.ms, nil
}

func (mm *mysqlMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (mm *mysqlMetrics) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "db_rdb_mysql_completed_query_total",
			Description: mm.queryTotal.Description(),
			Measure:     &mm.queryTotal,
			Aggregation: metrics.Count(),
		},
		{
			Name:        "db_rdb_mysql_query_latency",
			Description: mm.queryLatency.Description(),
			Measure:     &mm.queryLatency,
			Aggregation: metrics.DefaultMillisecondsDistribution,
		},
	}
}

func (mm *mysqlMetrics) SpanStart(ctx context.Context, eventName, query string) context.Context {
	ctx, span := trace.StartSpan(ctx, "vald/internal/db/rdb/mysql")
	if span != nil {
		span.AddAttributes(
			trace.StringAttribute("event_name", eventName),
			trace.StringAttribute("query", query),
		)
	}

	return context.WithValue(ctx, startTimeKey{}, time.Now())
}

func (mm *mysqlMetrics) SpanError(ctx context.Context, err error) {
	span := trace.FromContext(ctx)
	if span != nil {
		span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		span.End()
	}
}

func (mm *mysqlMetrics) SpanFinish(ctx context.Context) {
	span := trace.FromContext(ctx)
	if span != nil {
		span.End()
	}

	startTime, _ := ctx.Value(startTimeKey{}).(time.Time)
	latencyMillis := float64(time.Since(startTime)) / float64(time.Millisecond)

	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.ms = append(
		mm.ms,
		mm.queryTotal.M(1),
		mm.queryLatency.M(latencyMillis),
	)
}
