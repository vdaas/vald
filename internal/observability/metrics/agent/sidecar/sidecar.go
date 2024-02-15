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
package sidecar

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/observer"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
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

func (*sidecarMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        uploadTotalMetricsName,
				Description: uploadTotalMetricsDescription,
			},
			view.Stream{
<<<<<<< HEAD
				Aggregation: view.AggregationSum{},
=======
				Aggregation: meric.AggregationSum{},
>>>>>>> feature/agent/qbg
			},
		),
		view.NewView(
			view.Instrument{
				Name:        uploadBytesMetricsName,
				Description: uploadBytesMetricsDescription,
			},
			view.Stream{
<<<<<<< HEAD
				Aggregation: view.AggregationLastValue{},
=======
				Aggregation: meric.AggregationLastValue{},
>>>>>>> feature/agent/qbg
			},
		),
		view.NewView(
			view.Instrument{
				Name:        uploadLatencyMetricsName,
				Description: uploadLatencyMetricsDescription,
			},
			view.Stream{
<<<<<<< HEAD
				Aggregation: view.AggregationExplicitBucketHistogram{
=======
				Aggregation: meric.AggregationExplicitBucketHistogram{
>>>>>>> feature/agent/qbg
					Boundaries: metrics.RoughMillisecondsDistribution,
				},
			},
		),
	}, nil
}

func (sm *sidecarMetrics) Register(m metrics.Meter) error {
	uploadTotal, err := m.Int64ObservableCounter(
		uploadTotalMetricsName,
		metrics.WithDescription(uploadTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	uploadBytes, err := m.Int64ObservableGauge(
		uploadBytesMetricsName,
		metrics.WithDescription(uploadBytesMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}
	uploadLatency, err := m.Float64ObservableGauge(
		uploadLatencyMetricsName,
		metrics.WithDescription(uploadLatencyMetricsDescription),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			sm.mu.Lock()
			defer sm.mu.Unlock()

			if sm.info == nil {
				return nil
			}

			attrs := api.WithAttributes(
				attribute.String(sm.storageTypeKey, sm.info.StorageInfo.Type),
				attribute.String(sm.bucketNameKey, sm.info.BucketName),
				attribute.String(sm.filenameKey, sm.info.Filename),
			)

<<<<<<< HEAD
			o.ObserveInt64(uploadTotal, 1, attrs)
			o.ObserveInt64(uploadBytes, sm.info.Bytes, attrs)
			latencyMillis := float64(sm.info.EndTime.Sub(sm.info.StartTime)) / float64(time.Millisecond)
			o.ObserveFloat64(uploadLatency, latencyMillis, attrs)
=======
			o.ObserveInt64(uploadTotal, 1, attrs...)
			o.ObserveInt64(uploadBytes, sm.info.Bytes, attrs...)
			latencyMillis := float64(sm.info.EndTime.Sub(sm.info.StartTime)) / float64(time.Millisecond)
			o.ObserveFloat64(uploadLatency, latencyMillis, attrs...)
>>>>>>> feature/agent/qbg
			sm.info = nil

			return nil
		},
		uploadTotal,
		uploadBytes,
		uploadLatency,
	)
	return err
}

func (*sidecarMetrics) BeforeProcess(ctx context.Context, _ *observer.BackupInfo) (context.Context, error) {
	return ctx, nil
}

func (sm *sidecarMetrics) AfterProcess(_ context.Context, info *observer.BackupInfo) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.info = info
	return nil
}
