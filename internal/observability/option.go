package observability

// TODO: Fix observability-v2 to observability
import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability-v2/exporter"
	"github.com/vdaas/vald/internal/observability-v2/trace"
)

type Option func(*observability) error

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
	}
)

// WithErrGroup returns an option that sets the errgroup.
func WithErrGroup(eg errgroup.Group) Option {
	return func(o *observability) error {
		if eg != nil {
			o.eg = eg
		}
		return nil
	}
}

// WithExporters returns an option that sets the exporters.
func WithExporters(exps ...exporter.Exporter) Option {
	return func(o *observability) error {
		if len(exps) != 0 {
			if o.exporters == nil {
				o.exporters = exps
			} else {
				o.exporters = append(o.exporters, exps...)
			}
		}
		return nil
	}
}

// WithTracer returns an option that sets the tracer.
func WithTracer(tr trace.Tracer) Option {
	return func(o *observability) error {
		if tr != nil {
			o.tracer = tr
		}
		return nil
	}
}
