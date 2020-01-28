package retry

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
)

type Retry interface {
	Out(
		fn func(vals ...interface{}) error,
		vals ...interface{},
	)

	Outf(
		fn func(format string, vals ...interface{}) error,
		format string, vals ...interface{},
	)
}

type retry struct {
	warnfn  func(vals ...interface{})
	errorfn func(vals ...interface{})
}

func New(opts ...Option) Retry {
	r := new(retry)
	for _, opt := range append(defaultOption, opts...) {
		opt(r)
	}
	return r
}

func (r *retry) Out(
	fn func(vals ...interface{}) error,
	vals ...interface{},
) {
	r.Outf(func(format string, vals ...interface{}) error {
		return fn(vals...)
	}, "", vals...)
}

func (r *retry) Outf(
	fn func(format string, vals ...interface{}) error,
	format string, vals ...interface{},
) {
	if err := fn(format, vals...); err != nil {
		rv := reflect.ValueOf(fn)

		r.warnfn(errors.ErrLoggingRetry(err, rv))

		err = fn(format, vals...)
		if err != nil {
			r.errorfn(errors.ErrLoggingFailed(err, rv))

			err = fn(format, vals...)
			if err != nil {
				panic(err)
			}
		}
	}
}
