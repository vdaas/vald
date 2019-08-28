// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package transport provides http transport roundtrip option
package transport

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
)

type ert struct {
	transport http.RoundTripper
	bo        backoff.Backoff
}

func NewExpBackoff(opts ...Option) http.RoundTripper {
	e := new(ert)
	for _, opt := range append(defaultOpts, opts...) {
		opt(e)
	}

	return e
}

func (e *ert) RoundTrip(req *http.Request) (res *http.Response, err error) {
	r, err := e.bo.Do(req.Context(), func() (interface{}, error) {
		res, err := e.transport.RoundTrip(req)
		if err == nil {
			return res, nil
		}
		if res != nil {
			io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()
		}
		switch res.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusServiceUnavailable,
			http.StatusMovedPermanently,
			http.StatusBadGateway,
			http.StatusGatewayTimeout:
			return nil, errors.ErrTransportRetryable
		}
		return res, nil
	})

	if err != nil {
		return nil, err
	}

	res, ok := r.(*http.Response)
	if !ok {
		return nil, errors.ErrInvalidTypeConversion(r, res)
	}

	return res, nil
}
