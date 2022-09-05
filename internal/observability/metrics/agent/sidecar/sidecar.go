package sidecar

import (
	"context"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/observer"
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

func (sm *sidecarMetrics) Register(m metrics.Meter) error {
	uploadTotal, err := m.AsyncInt64().Counter(
		"agent_sidecar_completed_upload_total",
		metrics.WithDescription("cumulative count of completed upload execution"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	uploadBytes, err := m.AsyncInt64().Gauge(
		"agent_sidecar_upload_bytes",
		metrics.WithDescription("uploaded bytes at the last backup execution"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}
	uploadLatency, err := m.AsyncFloat64().UpDownCounter( // TODO:
		"agent_sidecar_upload_latency",
		metrics.WithDescription("upload latency"),
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
