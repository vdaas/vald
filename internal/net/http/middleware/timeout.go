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

package middleware

import (
	"bytes"
	"context"
	"maps"
	"net/http"
	"time"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type timeout struct {
	eg  errgroup.Group
	dur time.Duration
}

func NewTimeout(opts ...TimeoutOption) Wrapper {
	t := new(timeout)
	for _, opt := range append(defaultTimeoutOpts, opts...) {
		opt(t)
	}
	return t
}

// timeoutWriter buffers the handler's response so that the handler goroutine
// never touches the real http.ResponseWriter, which net/http recycles once
// ServeHTTP returns on timeout (same pattern as net/http.TimeoutHandler).
type timeoutWriter struct {
	h           http.Header
	buf         bytes.Buffer
	mu          sync.Mutex
	code        int
	wroteHeader bool
	timedOut    bool
}

func (tw *timeoutWriter) Header() http.Header { return tw.h }

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if !tw.wroteHeader {
		tw.wroteHeader = true
		tw.code = http.StatusOK
	}
	return tw.buf.Write(p)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.wroteHeader = true
	tw.code = code
}

// timeout discards the buffered response and makes later writes fail.
func (tw *timeoutWriter) timeout() {
	tw.mu.Lock()
	tw.timedOut = true
	tw.mu.Unlock()
}

// flush copies the buffered response to w. It must be called only after the handler returned.
func (tw *timeoutWriter) flush(w http.ResponseWriter) error {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	dst := w.Header()
	maps.Copy(dst, tw.h)
	if tw.wroteHeader {
		w.WriteHeader(tw.code)
	}
	if tw.buf.Len() == 0 {
		return nil
	}
	_, err := w.Write(tw.buf.Bytes())
	return err
}

type timeoutResult struct {
	err  error
	code int
}

func (t *timeout) Wrap(h rest.Func) rest.Func {
	return func(w http.ResponseWriter, r *http.Request) (code int, err error) {
		ctx, cancel := context.WithTimeout(r.Context(), t.dur)
		defer cancel()
		start := fastime.UnixNanoNow()
		tw := &timeoutWriter{h: make(http.Header)}
		// run the custom handler logic in go routine and report exactly one result
		// (a recovered panic included) to the buffered result channel.
		rch := make(chan timeoutResult, 1)
		hr := r.WithContext(ctx)
		t.eg.Go(func() error {
			var res timeoutResult
			// it is the responsibility for handler to close the request
			perr := safety.RecoverFunc(func() (err error) {
				res.code, res.err = h(tw, hr)
				return nil
			})()
			if perr != nil {
				res = timeoutResult{code: http.StatusInternalServerError, err: perr}
			}
			rch <- res
			return perr
		})

		select {
		case <-ctx.Done():
			// timeout passed or parent context canceled first,
			// the late output of the handler is discarded.
			tw.timeout()
			return http.StatusRequestTimeout, errors.ErrHandlerTimeout(ctx.Err(), time.Duration(fastime.UnixNanoNow()-start))
		case res := <-rch:
			if cerr := ctx.Err(); cerr != nil {
				tw.timeout()
				return http.StatusRequestTimeout, errors.ErrHandlerTimeout(errors.Join(cerr, res.err), time.Duration(fastime.UnixNanoNow()-start))
			}
			// handler finished first, may have error returned
			if ferr := tw.flush(w); ferr != nil && res.err == nil {
				return http.StatusServiceUnavailable, errors.ErrHandler(ferr)
			}
			if res.err != nil {
				return res.code, errors.ErrHandler(res.err)
			}
			return http.StatusOK, nil
		}
	}
}
