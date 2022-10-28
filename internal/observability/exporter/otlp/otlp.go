package otlp

import (
	"context"
	"os"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type exp struct {
	serviceName       string
	collectorEndpoint string

	traceExporter *otlptrace.Exporter
	traceProvider *trace.TracerProvider

	tBatchTimeout       time.Duration
	tExportTimeout      time.Duration
	tMaxExportBatchSize int
	tMaxQueueSize       int

	metricsExporter metric.Exporter
	meterProvider   *metric.MeterProvider
	metricsViews    []metrics.View

	mExportInterval time.Duration
	mExportTimeout  time.Duration
}

func New(opts ...Option) (exporter.Exporter, error) {
	e := new(exp)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(e); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return e, nil
}

func (e *exp) initTracer(ctx context.Context) (err error) {
	e.traceExporter, err = otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(e.collectorEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	e.traceProvider = trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(e.traceExporter,
			trace.WithBatchTimeout(e.tBatchTimeout),
			trace.WithExportTimeout(e.tExportTimeout),
			trace.WithMaxExportBatchSize(e.tMaxExportBatchSize),
			trace.WithMaxQueueSize(e.tMaxQueueSize),
		),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(e.serviceName),
		)),
	)
	return nil
}

func (e *exp) initMeter(ctx context.Context) (err error) {
	e.metricsExporter, err = otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(e.collectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	e.meterProvider = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(
			e.metricsExporter,
			metric.WithInterval(e.mExportInterval),
			metric.WithTimeout(e.mExportTimeout),
		), e.metricsViews...),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(e.serviceName),
			attribute.String("target_pod", os.Getenv("MY_POD_NAME")),                // TODO: fix it later
			attribute.String("target_node", os.Getenv("MY_NODE_NAME")),              // TODO: fix it later
			attribute.String("kubernetes_name", e.serviceName),                      // TODO: fix it later
			attribute.String("kubernetes_namespace", os.Getenv("MY_POD_NAMESPACE")), // TODO: fix it later
		)),
	)
	return nil
}

func (e *exp) Start(ctx context.Context) error {
	if err := e.initTracer(ctx); err != nil {
		return err
	}
	if err := e.initMeter(ctx); err != nil {
		return err
	}

	otel.SetTracerProvider(e.traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	global.SetMeterProvider(e.meterProvider)
	return nil
}

func (e *exp) Stop(ctx context.Context) error {
	if err := e.traceProvider.ForceFlush(ctx); err != nil {
		log.Errorf("failed to flush trace data: %v", err)
	}
	if err := e.traceProvider.Shutdown(ctx); err != nil {
		log.Errorf("failed to shutdown trace provider: %v", err)
	}
	if err := e.traceExporter.Shutdown(ctx); err != nil {
		log.Warn("failed to shutdown trace exporter: %v", err)
	}
	if err := e.meterProvider.ForceFlush(ctx); err != nil {
		log.Errorf("failed to flush metrics data: %v", err)
	}
	if err := e.metricsExporter.Shutdown(ctx); err != nil {
		log.Errorf("failed to shutdown metrics exporter: %v", err)
	}
	return nil
}
