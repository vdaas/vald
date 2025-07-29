//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package client

import (
	"net/http"
	"reflect"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	htr "github.com/vdaas/vald/internal/net/http/transport"
	"golang.org/x/net/http2"
)

type transport struct {
	*http.Transport
	backoffOpts []backoff.Option
}

// New initializes the HTTP2 transport with exponential backoff and returns the HTTP client for it, or returns any error occurred.
func New(opts ...Option) (*http.Client, error) {
	return NewWithTransport(http.DefaultTransport, opts...)
}

// NewWithTransport initializes the HTTP2 transport with the given RoundTripper and
// exponential backoff, returning the HTTP client for it, or any error that occurred.
func NewWithTransport(rt http.RoundTripper, opts ...Option) (*http.Client, error) {
	tr := new(transport)
	t, ok := rt.(*http.Transport)
	if ok {
		tr.Transport = t.Clone()
	} else {
		// Initialize with default transport if the provided one is not *http.Transport
		tr.Transport = http.DefaultTransport.(*http.Transport).Clone()
	}
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(tr); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}

	var err error
	err = http2.ConfigureTransport(tr.Transport)
	if err != nil {
		log.Warnf("Transport is already configured for HTTP2 error: %v", err)
	}

	return &http.Client{
		Transport: htr.NewExpBackoff(
			htr.WithRoundTripper(tr.Transport),
			htr.WithBackoff(
				backoff.New(tr.backoffOpts...),
			),
		),
	}, nil
}
