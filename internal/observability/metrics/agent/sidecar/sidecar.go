//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package sidecar

import (
	"context"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/observer"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	uploadTotalMetricsName        = "agent_sidecar_completed_upload_total"
	uploadTotalMetricsDescription = "Cumulative count of completed upload execution"

	uploadBytesMetricsName        = "agent_sidecar_upload_bytes"
	uploadBytesMetricsDescription = "Uploaded bytes at the last backup execution"

	uploadLatencyMetricsName        = "agent_sidecar_upload_latency"
	uploadLatencyMetricsDescription = "Upload latency"
)

type MetricsHook interface {
	metrics.Metric
	observer.Hook
}

type sidecarMetrics struct {
	storageTypeKey string
	bucketNameKey  string
	filenameKey    string

	mu   sync.Mutex
	info *observer.BackupInfo
}

func New() MetricsHook {
	return &sidecarMetrics{
		storageTypeKey: "agent_sidecar_storage_type",
		bucketNameKey:  "agent_sidecar_bucket_name",
		filenameKey:    "agent_sidecar_filename",
	}
}

func (sm *sidecarMetrics) View() ([]*metrics.View, error) {
	uploadTotal, err := view.New(
		view.MatchInstrumentName(uploadTotalMetricsName),
		view.WithSetDescription(uploadTotalMetricsDescription),
		view.WithSetAggregation(aggregation.Sum{}),
	)
	if err != nil {
		return nil, err
	}

	uploadBytes, err := view.New(
		view.MatchInstrumentName(uploadBytesMetricsName),
		view.WithSetDescription(uploadBytesMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	uploadLatency, err := view.New(
		view.MatchInstrumentName(uploadLatencyMetricsName),
		view.WithSetDescription(uploadLatencyMetricsDescription),
		view.WithSetAggregation(aggregation.ExplicitBucketHistogram{
			Boundaries: metrics.RoughMillisecondsDistribution,
		}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&uploadTotal,
		&uploadBytes,
		&uploadLatency,
	}, nil
}

func (sm *sidecarMetrics) Register(m metrics.Meter) error {
	uploadTotal, err := m.AsyncInt64().Counter(
		uploadTotalMetricsName,
		metrics.WithDescription(uploadTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	uploadBytes, err := m.AsyncInt64().Gauge(
		uploadBytesMetricsName,
		metrics.WithDescription(uploadBytesMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}
	uploadLatency, err := m.AsyncFloat64().Gauge(
		uploadLatencyMetricsName,
		metrics.WithDescription(uploadLatencyMetricsDescription),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			uploadTotal,
			uploadBytes,
			uploadLatency,
		},
		func(ctx context.Context) {
			sm.mu.Lock()
			defer sm.mu.Unlock()

			if sm.info == nil {
				return
			}

			attrs := []attribute.KeyValue{
				attribute.String(sm.storageTypeKey, sm.info.StorageInfo.Type),
				attribute.String(sm.bucketNameKey, sm.info.BucketName),
				attribute.String(sm.filenameKey, sm.info.Filename),
			}

			latencyMillis := float64(sm.info.EndTime.Sub(sm.info.StartTime)) / float64(time.Millisecond)

			uploadTotal.Observe(ctx, 1, attrs...)
			uploadBytes.Observe(ctx, sm.info.Bytes, attrs...)
			uploadLatency.Observe(ctx, latencyMillis, attrs...)

			sm.info = nil
		},
	)
}

func (sm *sidecarMetrics) BeforeProcess(ctx context.Context, info *observer.BackupInfo) (context.Context, error) {
	return ctx, nil
}

func (sm *sidecarMetrics) AfterProcess(ctx context.Context, info *observer.BackupInfo) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.info = info
	return nil
}
