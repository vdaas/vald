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

// Package cassandra provides cassandra metrics functions
package cassandra

import (
	"context"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type cassandraMetrics struct {
	queryTotal         metrics.Int64Measure
	queryAttemptsTotal metrics.Int64Measure
	queryLatency       metrics.Float64Measure

	keyspaceKey    metrics.Key
	clusterNameKey metrics.Key
	dataCenterKey  metrics.Key
	hostIDKey      metrics.Key
	hostPortKey    metrics.Key
	rackKey        metrics.Key
	versionKey     metrics.Key

	mu sync.Mutex
	ms []metrics.MeasurementWithTags
}

type Observer interface {
	metrics.Metric
	cassandra.QueryObserver
}

func New() (o Observer, err error) {
	cms := new(cassandraMetrics)

	cms.queryTotal = *metrics.Int64(
		metrics.ValdOrg+"/db/nosql/cassandra/completed_query_total",
		"cumulative count of completed queries",
		metrics.UnitDimensionless,
	)

	cms.queryAttemptsTotal = *metrics.Int64(
		metrics.ValdOrg+"/db/nosql/cassandra/completed_query_attempts_total",
		"cumulative count of query attempts (number of retry or fetching next page)",
		metrics.UnitDimensionless,
	)

	cms.queryLatency = *metrics.Float64(
		metrics.ValdOrg+"/db/nosql/cassandra/query_latency",
		"query latency",
		metrics.UnitMilliseconds,
	)

	cms.keyspaceKey, err = metrics.NewKey("cassandra_keyspace")
	if err != nil {
		return nil, err
	}

	cms.clusterNameKey, err = metrics.NewKey("cassandra_cluster_name")
	if err != nil {
		return nil, err
	}

	cms.dataCenterKey, err = metrics.NewKey("cassandra_data_center")
	if err != nil {
		return nil, err
	}

	cms.hostIDKey, err = metrics.NewKey("cassandra_host_id")
	if err != nil {
		return nil, err
	}

	cms.hostPortKey, err = metrics.NewKey("cassandra_host_port")
	if err != nil {
		return nil, err
	}

	cms.rackKey, err = metrics.NewKey("cassandra_rack")
	if err != nil {
		return nil, err
	}

	cms.versionKey, err = metrics.NewKey("cassandra_version")
	if err != nil {
		return nil, err
	}

	cms.ms = make([]metrics.MeasurementWithTags, 0)

	return cms, nil
}

func (cm *cassandraMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (cm *cassandraMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	cm.mu.Lock()
	defer func() {
		cm.ms = make([]metrics.MeasurementWithTags, 0)
		cm.mu.Unlock()
	}()

	return cm.ms, nil
}

func (cm *cassandraMetrics) View() []*metrics.View {
	keys := []metrics.Key{
		cm.keyspaceKey,
		cm.clusterNameKey,
		cm.dataCenterKey,
		cm.hostIDKey,
		cm.hostPortKey,
		cm.rackKey,
		cm.versionKey,
	}

	return []*metrics.View{
		{
			Name:        "db_nosql_cassandra_completed_query_total",
			Description: cm.queryTotal.Description(),
			TagKeys:     keys,
			Measure:     &cm.queryTotal,
			Aggregation: metrics.Count(),
		},
		{
			Name:        "db_nosql_cassandra_query_attempts_total",
			Description: cm.queryAttemptsTotal.Description(),
			TagKeys:     keys,
			Measure:     &cm.queryAttemptsTotal,
			Aggregation: metrics.Count(),
		},
		{
			Name:        "db_nosql_cassandra_query_latency",
			Description: cm.queryLatency.Description(),
			TagKeys:     keys,
			Measure:     &cm.queryLatency,
			Aggregation: metrics.DefaultMillisecondsDistribution,
		},
	}
}

// ObserveQuery updates the member variables by passed ObservedQuery instance.
func (cm *cassandraMetrics) ObserveQuery(ctx context.Context, q cassandra.ObservedQuery) {
	latencyMillis := float64(q.End.Sub(q.Start)) / float64(time.Millisecond)
	tags := map[metrics.Key]string{
		cm.keyspaceKey:    q.Keyspace,
		cm.clusterNameKey: q.Host.ClusterName(),
		cm.dataCenterKey:  q.Host.DataCenter(),
		cm.hostIDKey:      q.Host.HostID(),
		cm.hostPortKey:    q.Host.HostnameAndPort(),
		cm.rackKey:        q.Host.Rack(),
		cm.versionKey:     q.Host.Version().String(),
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.ms = append(
		cm.ms,
		metrics.MeasurementWithTags{
			Measurement: cm.queryTotal.M(1),
			Tags:        tags,
		},
		metrics.MeasurementWithTags{
			Measurement: cm.queryAttemptsTotal.M(1 + int64(q.Attempt)),
			Tags:        tags,
		},
		metrics.MeasurementWithTags{
			Measurement: cm.queryLatency.M(latencyMillis),
			Tags:        tags,
		},
	)
}
