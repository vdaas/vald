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

// Package sidecar provides functions for sidecar stats
package sidecar

import (
	"context"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/observer"
)

type sidecarMetrics struct {
	uploadTotal   metrics.Int64Measure
	uploadBytes   metrics.Int64Measure
	uploadLatency metrics.Float64Measure

	storageTypeKey metrics.Key
	bucketNameKey  metrics.Key
	filenameKey    metrics.Key

	mu sync.Mutex
	ms []metrics.MeasurementWithTags
}

type MetricsHook interface {
	metrics.Metric
	observer.Hook
}

func New() (MetricsHook, error) {
	var err error
	sm := new(sidecarMetrics)

	sm.uploadTotal = *metrics.Int64(
		metrics.ValdOrg+"/agent/sidecar/completed_upload_total",
		"cumulative count of completed upload execution",
		metrics.UnitDimensionless)

	sm.uploadBytes = *metrics.Int64(
		metrics.ValdOrg+"/agent/sidecar/upload_bytes",
		"uploaded bytes at the last backup execution",
		metrics.UnitBytes)

	sm.uploadLatency = *metrics.Float64(
		metrics.ValdOrg+"/agent/sidecar/upload_latency",
		"upload latency",
		metrics.UnitMilliseconds)

	sm.storageTypeKey, err = metrics.NewKey("agent_sidecar_storage_type")
	if err != nil {
		return nil, err
	}

	sm.bucketNameKey, err = metrics.NewKey("agent_sidecar_bucket_name")
	if err != nil {
		return nil, err
	}

	sm.filenameKey, err = metrics.NewKey("agent_sidecar_filename")
	if err != nil {
		return nil, err
	}

	sm.ms = make([]metrics.MeasurementWithTags, 0)

	return sm, nil
}

func (sm *sidecarMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (sm *sidecarMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	sm.mu.Lock()
	defer func() {
		sm.ms = make([]metrics.MeasurementWithTags, 0)
		sm.mu.Unlock()
	}()

	return sm.ms, nil
}

func (sm *sidecarMetrics) View() []*metrics.View {
	uploadKeys := []metrics.Key{
		sm.storageTypeKey,
		sm.bucketNameKey,
		sm.filenameKey,
	}

	return []*metrics.View{
		{
			Name:        "agent_sidecar_completed_upload_total",
			Description: sm.uploadTotal.Description(),
			TagKeys:     uploadKeys,
			Measure:     &sm.uploadTotal,
			Aggregation: metrics.Count(),
		},
		{
			Name:        "agent_sidecar_upload_bytes",
			Description: sm.uploadBytes.Description(),
			TagKeys:     uploadKeys,
			Measure:     &sm.uploadBytes,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_sidecar_upload_latency",
			Description: sm.uploadLatency.Description(),
			TagKeys:     uploadKeys,
			Measure:     &sm.uploadLatency,
			Aggregation: metrics.RoughMillisecondsDistribution,
		},
	}
}

func (sm *sidecarMetrics) BeforeProcess(ctx context.Context, info *observer.BackupInfo) (context.Context, error) {
	return ctx, nil
}

func (sm *sidecarMetrics) AfterProcess(ctx context.Context, info *observer.BackupInfo) error {
	tags := map[metrics.Key]string{
		sm.storageTypeKey: info.StorageInfo.Type,
		sm.bucketNameKey:  info.BucketName,
		sm.filenameKey:    info.Filename,
	}

	latencyMillis := float64(info.EndTime.Sub(info.StartTime)) / float64(time.Millisecond)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.ms = append(
		sm.ms,
		metrics.MeasurementWithTags{
			Measurement: sm.uploadTotal.M(1),
			Tags:        tags,
		},
		metrics.MeasurementWithTags{
			Measurement: sm.uploadBytes.M(info.Bytes),
			Tags:        tags,
		},
		metrics.MeasurementWithTags{
			Measurement: sm.uploadLatency.M(latencyMillis),
			Tags:        tags,
		},
	)

	return nil
}
