// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package retry

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
)

type Retry interface {
	Out(
		fn func(vals ...any) error,
		vals ...any,
	)
	Outf(
		fn func(format string, vals ...any) error,
		format string, vals ...any,
	)
}

type retry struct {
	warnFn  func(vals ...any)
	errorFn func(vals ...any)
}

// failure marks a logging-failure report handed to warnFn/errorFn. Those are
// usually the very logger that is failing and call back into Out, so a report
// that cannot be written is dropped instead of being reported again, which
// would recurse without bound.
type failure struct{ error }

func (f failure) Unwrap() error { return f.error }

func isFailure(vals []any) bool {
	if len(vals) != 1 {
		return false
	}
	_, ok := vals[0].(failure)
	return ok
}

func New(opts ...Option) Retry {
	r := new(retry)
	for _, opt := range append(defaultOption, opts...) {
		opt(r)
	}
	return r
}

func (r *retry) Out(fn func(vals ...any) error, vals ...any) {
	if fn != nil {
		if err := fn(vals...); err != nil && !isFailure(vals) {
			rv := reflect.ValueOf(fn)
			r.report(r.warnFn, errors.ErrLoggingRetry(err, rv))
			err = fn(vals...)
			if err != nil {
				r.report(r.errorFn, errors.ErrLoggingFailed(err, rv))
				err = fn(vals...)
				if err != nil {
					// do not panic: a failing logger must not crash the caller
					r.report(r.errorFn, errors.ErrLoggingFailed(err, rv))
				}
			}
		}
	}
}

func (r *retry) Outf(fn func(format string, vals ...any) error, format string, vals ...any) {
	if fn != nil {
		if err := fn(format, vals...); err != nil {
			rv := reflect.ValueOf(fn)
			r.report(r.warnFn, errors.ErrLoggingRetry(err, rv))
			err = fn(format, vals...)
			if err != nil {
				r.report(r.errorFn, errors.ErrLoggingFailed(err, rv))

				err = fn(format, vals...)
				if err != nil {
					// do not panic: a failing logger must not crash the caller
					r.report(r.errorFn, errors.ErrLoggingFailed(err, rv))
				}
			}
		}
	}
}

func (*retry) report(fn func(vals ...any), err error) {
	fn(failure{err})
}
