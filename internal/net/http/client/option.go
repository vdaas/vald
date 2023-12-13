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

package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/timeutil"
)

// Option represent the functional option for transport.
type Option func(*transport) error

var defaultOptions = []Option{
	WithProxy(http.ProxyFromEnvironment),
	WithEnableKeepalives(true),
	WithEnableCompression(true),
}

// WithProxy returns the option to set the transport proxy.
func WithProxy(px func(*http.Request) (*url.URL, error)) Option {
	return func(tr *transport) error {
		if px == nil {
			return errors.NewErrInvalidOption("proxy", px)
		}
		tr.Proxy = px

		return nil
	}
}

// WithDialContext returns the option to set the dial context.
func WithDialContext(dx func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(tr *transport) error {
		if dx == nil {
			return errors.NewErrInvalidOption("dialContext", dx)
		}
		tr.DialContext = dx

		return nil
	}
}

// WithTLSHandshakeTimeout returns the option to set the TLS handshake timeout.
func WithTLSHandshakeTimeout(dur string) Option {
	return func(tr *transport) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("TLSHandshakeTimeout", dur)
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return errors.NewErrCriticalOption("TLSHandshakeTimeout", dur, err)
		}

		tr.TLSHandshakeTimeout = d

		return nil
	}
}

// WithEnableKeepalives returns the option to enable keep alive.
func WithEnableKeepalives(enable bool) Option {
	return func(tr *transport) error {
		tr.DisableKeepAlives = !enable

		return nil
	}
}

// WithEnableCompression returns the option to enable compression.
func WithEnableCompression(enable bool) Option {
	return func(tr *transport) error {
		tr.DisableCompression = !enable

		return nil
	}
}

// WithMaxIdleConns returns the option to set the max idle connection.
func WithMaxIdleConns(cn int) Option {
	return func(tr *transport) error {
		tr.MaxIdleConns = cn

		return nil
	}
}

// WithMaxIdleConnsPerHost returns the option to set the max idle connection per host.
func WithMaxIdleConnsPerHost(cn int) Option {
	return func(tr *transport) error {
		tr.MaxIdleConnsPerHost = cn

		return nil
	}
}

// WithMaxConnsPerHost returns the option to set the max connections per host.
func WithMaxConnsPerHost(cn int) Option {
	return func(tr *transport) error {
		tr.MaxConnsPerHost = cn

		return nil
	}
}

// WithIdleConnTimeout returns the option to set the idle connection timeout.
func WithIdleConnTimeout(dur string) Option {
	return func(tr *transport) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("idleConnTimeout", dur)
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return errors.NewErrCriticalOption("idleConnTimeout", dur, err)
		}

		tr.IdleConnTimeout = d

		return nil
	}
}

// WithResponseHeaderTimeout returns the option to set the response header timeout.
func WithResponseHeaderTimeout(dur string) Option {
	return func(tr *transport) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("responseHeaderTimeout", dur)
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return errors.NewErrCriticalOption("responseHeaderTimeout", dur, err)
		}

		tr.ResponseHeaderTimeout = d

		return nil
	}
}

// WithExpectContinueTimeout returns the option to set the expect continue timeout.
func WithExpectContinueTimeout(dur string) Option {
	return func(tr *transport) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("expectContinueTimeout", dur)
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return errors.NewErrCriticalOption("expectContinueTimeout", dur, err)
		}

		tr.ExpectContinueTimeout = d

		return nil
	}
}

// WithProxyConnectHeader returns the option to set the proxy connect header.
func WithProxyConnectHeader(header http.Header) Option {
	return func(tr *transport) error {
		if header == nil {
			return errors.NewErrInvalidOption("proxyConnectHeader", header)
		}
		tr.ProxyConnectHeader = header

		return nil
	}
}

// WithMaxResponseHeaderBytes returns the option to set the max response header bytes.
func WithMaxResponseHeaderBytes(bs int64) Option {
	return func(tr *transport) error {
		tr.MaxResponseHeaderBytes = bs
		return nil
	}
}

// WithWriteBufferSize returns the option to set the write buffer size.
func WithWriteBufferSize(bs int64) Option {
	return func(tr *transport) error {
		tr.WriteBufferSize = int(bs)
		return nil
	}
}

// WithReadBufferSize returns the option to set the read buffer size.
func WithReadBufferSize(bs int64) Option {
	return func(tr *transport) error {
		tr.ReadBufferSize = int(bs)
		return nil
	}
}

// WithForceAttemptHTTP2 returns the option to force attempt HTTP2 for the HTTP transport.
func WithForceAttemptHTTP2(force bool) Option {
	return func(tr *transport) error {
		tr.ForceAttemptHTTP2 = force
		return nil
	}
}

// WithBackoffOpts returns the option to set the options to initialize backoff.
func WithBackoffOpts(opts ...backoff.Option) Option {
	return func(tr *transport) error {
		if len(opts) == 0 {
			return errors.NewErrInvalidOption("backoffOpts", opts)
		}
		if tr.backoffOpts == nil {
			tr.backoffOpts = opts
			return nil
		}

		tr.backoffOpts = append(tr.backoffOpts, opts...)
		return nil
	}
}
