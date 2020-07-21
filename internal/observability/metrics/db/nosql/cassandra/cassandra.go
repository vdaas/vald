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

// Package cassandra provides cassandra metrics functions
package cassandra

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type cassandraMetrics struct {
	queryTotal   metrics.Int64Measure
	queryLatency metrics.Float64Measure
}

type Observer interface {
	metrics.Metric
	cassandra.QueryObserver
}

func New() Observer {
	return &cassandraMetrics{
		queryTotal:   *metrics.Int64(metrics.ValdOrg+"/db/nosql/cassandra/completed_query_total", "cumulative count of completed queries", metrics.UnitDimensionless),
		queryLatency: *metrics.Float64(metrics.ValdOrg+"/db/nosql/cassandra/query_latency", "query latency", metrics.UnitMilliseconds),
	}
}

func (cm *cassandraMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (cm *cassandraMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (cm *cassandraMetrics) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "db_nosql_cassandra_completed_query_total",
			Description: "cumulative count of completed queries",
			Measure:     &cm.queryTotal,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "db_nosql_cassandra_query_latency",
			Description: "query latency",
			Measure:     &cm.queryLatency,
			Aggregation: metrics.DefaultMillisecondsDistribution,
		},
	}
}

// ObserveQuery updates the member variables by passed ObservedQuery instance.
func (cm *cassandraMetrics) ObserveQuery(ctx context.Context, q cassandra.ObservedQuery) {
	latencyMillis := float64(q.End.Sub(q.Start)) / float64(time.Millisecond)
	cm.queryTotal.M(1)
	cm.queryLatency.M(latencyMillis)
}
