//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package transport provides http transport roundtrip option
package transport

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type ert struct {
	transport http.RoundTripper
	bo        backoff.Backoff
}

// NewExpBackoff returns the backoff roundtripper implementation
func NewExpBackoff(opts ...Option) http.RoundTripper {
	e := new(ert)
	for _, opt := range append(defaultOpts, opts...) {
		opt(e)
	}

	return e
}

// RoundTrip round trip the HTTP request and return the response
func (e *ert) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if e.bo == nil {
		return e.roundTrip(req)
	}
	var r interface{}
	r, err = e.bo.Do(req.Context(), func() (interface{}, error) {
		return e.roundTrip(req)
	})
	if err != nil {
		return nil, err
	}

	var ok bool
	res, ok = r.(*http.Response)
	if !ok {
		return nil, errors.ErrInvalidTypeConversion(r, res)
	}
	return res, nil
}

func (e *ert) roundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = e.transport.RoundTrip(req)
	if err == nil {
		return res, nil
	}
	if res == nil {
		return nil, err
	}
	_, err = io.Copy(ioutil.Discard, res.Body)
	if err != nil {
		log.Error(err)
	}
	err = res.Body.Close()
	if err != nil {
		log.Error(err)
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
}
