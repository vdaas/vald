package observability

// TODO: Fix observability-v2 to observability
import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability-v2/exporter"
	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"github.com/vdaas/vald/internal/observability-v2/trace"
)

type Observability interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	Stop(ctx context.Context) error
}

type observability struct {
	eg        errgroup.Group
	exporters []exporter.Exporter
	tracer    trace.Tracer
	metrics   []metrics.Metric
}

func NewObservability(opts ...Option) (Observability, error) {
	o := &observability{}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(o); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return o, nil
}

func (o *observability) PreStart(ctx context.Context) error {
	for i, ex := range o.exporters {
		if err := ex.Start(ctx); err != nil {
			for _, ex := range o.exporters[:i] {
				if err := ex.Stop(ctx); err != nil {
					log.Error(err)
				}
			}
			return err
		}
	}

	meter := metrics.GetMeter()
	for _, m := range o.metrics {
		if err := m.Register(meter); err != nil {
			return err
		}
	}

	if err := o.tracer.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (o *observability) Start(ctx context.Context) <-chan error {
	return nil
}

func (o *observability) Stop(ctx context.Context) (werr error) {
	for _, ex := range o.exporters {
		if err := ex.Stop(ctx); err != nil {
			log.Error(err)
			werr = errors.Wrap(werr, err.Error())
		}
	}
	return
}
