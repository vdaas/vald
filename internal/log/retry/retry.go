// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
		fn func(vals ...interface{}) error,
		vals ...interface{},
	)
	Outf(
		fn func(format string, vals ...interface{}) error,
		format string, vals ...interface{},
	)
}

type retry struct {
	warnFn  func(vals ...interface{})
	errorFn func(vals ...interface{})
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
	if fn != nil {
		if err := fn(vals...); err != nil {
			rv := reflect.ValueOf(fn)
			r.warnFn(errors.ErrLoggingRetry(err, rv))
			err = fn(vals...)
			if err != nil {
				r.errorFn(errors.ErrLoggingFailed(err, rv))
				err = fn(vals...)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (r *retry) Outf(
	fn func(format string, vals ...interface{}) error,
	format string, vals ...interface{},
) {
	if fn != nil {
		if err := fn(format, vals...); err != nil {
			rv := reflect.ValueOf(fn)
			r.warnFn(errors.ErrLoggingRetry(err, rv))
			err = fn(format, vals...)
			if err != nil {
				r.errorFn(errors.ErrLoggingFailed(err, rv))

				err = fn(format, vals...)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
