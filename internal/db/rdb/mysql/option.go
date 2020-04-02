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

package mysql

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*mySQLClient) error

var (
	defaultOpts = []Option{
		WithCharset("utf8mb4"),
		WithTimezone("Local"),
		WithInitialPingDuration("30ms"),
		WithInitialPingTimeLimit("5m"),
		// WithConnectionLifeTimeLimit("2m"),
		// WithMaxOpenConns(40),
		// WithMaxIdleConns(50),
	}
)

func WithTimezone(tz string) Option {
	return func(m *mySQLClient) error {
		if tz != "" {
			m.timezone = tz
		}
		return nil
	}
}

func WithCharset(cs string) Option {
	return func(m *mySQLClient) error {
		if cs != "" {
			m.charset = cs
		}
		return nil
	}
}

func WithDB(db string) Option {
	return func(m *mySQLClient) error {
		if db != "" {
			m.db = db
		}
		return nil
	}
}

func WithHost(host string) Option {
	return func(m *mySQLClient) error {
		if host != "" {
			m.host = host
		}
		return nil
	}
}

func WithPort(port int) Option {
	return func(m *mySQLClient) error {
		m.port = port
		return nil
	}
}

func WithUser(user string) Option {
	return func(m *mySQLClient) error {
		if user != "" {
			m.user = user
		}
		return nil
	}
}

func WithPass(pass string) Option {
	return func(m *mySQLClient) error {
		if pass != "" {
			m.pass = pass
		}
		return nil
	}
}

func WithName(name string) Option {
	return func(m *mySQLClient) error {
		if name != "" {
			m.name = name
		}
		return nil
	}
}

func WithInitialPingTimeLimit(lim string) Option {
	return func(m *mySQLClient) error {
		if lim == "" {
			return nil
		}
		pd, err := timeutil.Parse(lim)
		if err != nil {
			pd = time.Second * 30
		}
		m.initialPingTimeLimit = pd
		return nil
	}
}

func WithInitialPingDuration(dur string) Option {
	return func(m *mySQLClient) error {
		if dur == "" {
			return nil
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Millisecond * 50
		}
		m.initialPingDuration = pd
		return nil
	}
}

func WithConnectionLifeTimeLimit(dur string) Option {
	return func(m *mySQLClient) error {
		if dur == "" {
			return nil
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Second * 30
		}
		m.connMaxLifeTime = pd
		return nil
	}
}

func WithMaxIdleConns(conns int) Option {
	return func(m *mySQLClient) error {
		if conns != 0 {
			m.maxIdleConns = conns
		}
		return nil
	}
}

func WithMaxOpenConns(conns int) Option {
	return func(m *mySQLClient) error {
		if conns != 0 {
			m.maxOpenConns = conns
		}
		return nil
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(m *mySQLClient) error {
		if cfg != nil {
			m.tlsConfig = cfg
		}
		return nil
	}
}

func WithDialer(der func(ctx context.Context, addr, port string) (net.Conn, error)) Option {
	return func(m *mySQLClient) error {
		if der != nil {
			m.dialer = der
		}
		return nil
	}
}
