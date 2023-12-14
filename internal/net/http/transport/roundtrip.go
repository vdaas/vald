//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package transport provides http transport roundtrip option
package transport

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
)

type ert struct {
	transport http.RoundTripper
	bo        backoff.Backoff
}

// NewExpBackoff returns the backoff roundtripper implementation.
func NewExpBackoff(opts ...Option) http.RoundTripper {
	e := new(ert)
	for _, opt := range append(defaultOptions, opts...) {
		opt(e)
	}

	return e
}

// RoundTrip round trip the HTTP request and return the response.
// If backoff is not set, the default roundTrip implementation will be used.
// It round trip the request and returns the response, and return any error occurred.
// It returns errors.ErrTransportRetryable to indicate if the request is consider as retryable.
func (e *ert) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if req != nil {
		defer closeBody(req.Body)
	}

	if e.bo == nil {
		return e.doRoundTrip(req)
	}
	_, err = e.bo.Do(req.Context(), func(ctx context.Context) (interface{}, bool, error) {
		r, err := e.doRoundTrip(req)
		if err != nil {
			return nil, errors.Is(err, errors.ErrTransportRetryable), err
		}
		res = r
		return r, false, nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *ert) doRoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = e.transport.RoundTrip(req)
	if err != nil {
		log.Error(err)
		if res != nil { // just in case we check the response as it depends on RoundTrip impl.
			closeBody(res.Body)
			if retryableStatusCode(res.StatusCode) {
				return nil, errors.Join(errors.ErrTransportRetryable, err)
			}
		}
		return nil, err
	}

	if res != nil && retryableStatusCode(res.StatusCode) {
		closeBody(res.Body)
		return nil, errors.ErrTransportRetryable
	}
	return res, nil
}

func retryableStatusCode(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
		http.StatusBadGateway,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

func closeBody(rc io.ReadCloser) {
	if rc != nil {
		if _, err := io.Copy(io.Discard, rc); err != nil {
			log.Error(err)
		}
		if err := rc.Close(); err != nil {
			log.Error(err)
		}
	}
}
