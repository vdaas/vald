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

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type cassandraMetrics struct {
	queryObserver cassandra.QueryObserver
	queryTotal    metrics.Int64Measure
	queryMsTotal  metrics.Float64Measure
}

func New(qo cassandra.QueryObserver) metrics.Metric {
	return &cassandraMetrics{
		queryObserver: qo,
		queryTotal:    *metrics.Int64(metrics.ValdOrg+"/db/nosql/cassandra/completed_query_total", "cumulative count of completed queries", metrics.UnitDimensionless),
		queryMsTotal:  *metrics.Float64(metrics.ValdOrg+"/db/nosql/cassandra/query_milliseconds_total", "cumulative count of query time in milliseconds", metrics.UnitMilliseconds),
	}
}

func (cm *cassandraMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		cm.queryTotal.M(int64(cm.queryObserver.CompletedQueryTotal())),
		cm.queryMsTotal.M(float64(cm.queryObserver.QueryNsTotal()) / 1000000.0),
	}, nil
}

func (cm *cassandraMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (cm *cassandraMetrics) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "completed_query_total",
			Description: "cumulative count of completed queries",
			Measure:     &cm.queryTotal,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "query_milliseconds_total",
			Description: "cumulative count of query time in milliseconds",
			Measure:     &cm.queryMsTotal,
			Aggregation: metrics.Count(),
		},
	}
}
