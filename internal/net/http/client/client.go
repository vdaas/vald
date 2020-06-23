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

package client

import (
	"net/http"
	"reflect"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	htr "github.com/vdaas/vald/internal/net/http/transport"
	"golang.org/x/net/http2"
)

type transport struct {
	*http.Transport
	backoffOpts []backoff.Option
}

func New(opts ...Option) (*http.Client, error) {
	tr := new(transport)
	tr.Transport = new(http.Transport)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(tr); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	err := http2.ConfigureTransport(tr.Transport)
	if err != nil {
		return nil, err
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
