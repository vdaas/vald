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
	bo          backoff.Backoff
	backoffOpts []backoff.Option
}

// New initializes the HTTP2 transport with exponential backoff and returns the HTTP client for it, or returns any error occurred.
func New(opts ...Option) (c *http.Client, err error) {
	tr := new(transport)
	tr.Transport = new(http.Transport)

	e := new(errors.ErrCriticalOption)
	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(tr); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}

	if _, err = http2.ConfigureTransports(tr.Transport); err != nil {
		// ConfigureTransports configures a net/http HTTP/1 Transport to use HTTP/2.
		// It returns an error if transport has already been HTTP/2-enabled.
		// so we can ignore this error
		log.Warnf("failed to configure http2 transport: %v", err)
	}

	return &http.Client{
		Transport: htr.NewExpBackoff(
			htr.WithRoundTripper(tr.Transport),
			htr.WithBackoff(func() backoff.Backoff {
				if tr.bo != nil {
					return tr.bo
				}
				return backoff.New(tr.backoffOpts...)
			}()),
		),
	}, nil
}
