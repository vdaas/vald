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

// Package middleware provides rest.Func Middleware
package middleware

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/safety"
)

type timeout struct {
	dur time.Duration
	eg  errgroup.Group
}

func NewTimeout(opts ...TimeoutOption) Wrapper {
	t := new(timeout)
	for _, opt := range append(defaultTimeoutOpts, opts...) {
		opt(t)
	}
	return t
}

func (t *timeout) Wrap(h rest.Func) rest.Func {
	return func(w http.ResponseWriter, r *http.Request) (code int, err error) {
		ctx, cancel := context.WithTimeout(r.Context(), t.dur)
		defer cancel()
		start := fastime.UnixNanoNow()
		// run the custom handler logic in go routine,
		// report error to error channel
		ech := make(chan error, 1)
		sch := make(chan int, 1)
		t.eg.Go(safety.RecoverFunc(func() (err error) {
			defer close(ech)
			defer close(sch)
			// it is the responsibility for handler to close the request
			var code int
			code, err = h(w, r.WithContext(ctx))
			if err != nil {
				sch <- code
				runtime.Gosched()
			}
			ech <- err
			return nil
		}))

		select {
		case err = <-ech:
			// handler finished first, may have error returned
			if err != nil {
				code = <-sch
				err = errors.ErrHandler(err)
				return code, err
			}
			return http.StatusOK, nil
		case <-ctx.Done():
			// timeout passed or parent context canceled first,
			// it is the responsibility for handler to response to the user
			return http.StatusRequestTimeout, errors.ErrHandlerTimeout(ctx.Err(), time.Duration(fastime.UnixNanoNow()-start))
		}
	}
}
