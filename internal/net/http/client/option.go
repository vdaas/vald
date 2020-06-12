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
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*transport) error

var (
	defaultOptions = []Option{
		WithProxy(http.ProxyFromEnvironment),
		WithEnableKeepAlives(true),
		WithEnableCompression(true),
	}
)

func WithProxy(px func(*http.Request) (*url.URL, error)) Option {
	return func(tr *transport) error {
		if px != nil {
			tr.Proxy = px
		}

		return nil
	}
}

func WithDialContext(dx func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(tr *transport) error {
		if dx != nil {
			tr.DialContext = dx
		}

		return nil
	}

}

func WithTLSHandshakeTimeout(dur string) Option {
	return func(tr *transport) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		tr.TLSHandshakeTimeout = d

		return nil
	}
}

func WithEnableKeepAlives(enable bool) Option {
	return func(tr *transport) error {
		tr.DisableKeepAlives = !enable

		return nil
	}
}

func WithEnableCompression(enable bool) Option {
	return func(tr *transport) error {
		tr.DisableCompression = !enable

		return nil
	}
}

func WithMaxIdleConns(cn int) Option {
	return func(tr *transport) error {
		tr.MaxIdleConns = cn

		return nil
	}
}

func WithMaxIdleConnsPerHost(cn int) Option {
	return func(tr *transport) error {
		tr.MaxIdleConnsPerHost = cn

		return nil
	}
}

func WithMaxConnsPerHost(cn int) Option {
	return func(tr *transport) error {
		tr.MaxConnsPerHost = cn

		return nil
	}
}

func WithIdleConnTimeout(dur string) Option {
	return func(tr *transport) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		tr.IdleConnTimeout = d

		return nil
	}
}

func WithResponseHeaderTimeout(dur string) Option {
	return func(tr *transport) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		tr.ResponseHeaderTimeout = d

		return nil
	}
}

func WithExpectContinueTimeout(dur string) Option {
	return func(tr *transport) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		tr.ExpectContinueTimeout = d

		return nil
	}
}

func WithProxyConnectHeader(header http.Header) Option {
	return func(tr *transport) error {
		if header != nil {
			tr.ProxyConnectHeader = header
		}

		return nil
	}
}

func WithMaxResponseHeaderBytes(bs int64) Option {
	return func(tr *transport) error {
		tr.MaxResponseHeaderBytes = bs

		return nil
	}
}

func WithWriteBufferSize(bs int64) Option {
	return func(tr *transport) error {
		tr.WriteBufferSize = int(bs)

		return nil
	}
}

func WithReadBufferSize(bs int64) Option {
	return func(tr *transport) error {
		tr.ReadBufferSize = int(bs)

		return nil
	}
}

func WithForceAttemptHTTP2(force bool) Option {
	return func(tr *transport) error {
		tr.ForceAttemptHTTP2 = force

		return nil
	}
}

func WithBackoffOpts(opts ...backoff.Option) Option {
	return func(tr *transport) error {
		if tr.backoffOpts == nil {
			tr.backoffOpts = opts
		}

		tr.backoffOpts = append(tr.backoffOpts, opts...)

		return nil
	}
}
