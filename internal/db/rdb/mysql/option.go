//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Option represents the functional option for mySQLClient.
type Option func(*mySQLClient) error

var defaultOptions = []Option{
	WithCharset("utf8mb4"),
	WithTimezone("Local"),
	WithInitialPingDuration("30ms"),
	WithInitialPingTimeLimit("5m"),
	// WithConnectionLifeTimeLimit("2m"),
	// WithMaxOpenConns(40),
	// WithMaxIdleConns(50),
}

// WithTimezone returns the option to set the timezone.
func WithTimezone(tz string) Option {
	return func(m *mySQLClient) error {
		if tz != "" {
			m.timezone = tz
		}
		return nil
	}
}

// WithCharset returns the option to set the charset.
func WithCharset(cs string) Option {
	return func(m *mySQLClient) error {
		if cs != "" {
			m.charset = cs
		}
		return nil
	}
}

// WithDB returns the option to set the db.
func WithDB(db string) Option {
	return func(m *mySQLClient) error {
		if db != "" {
			m.db = db
		}
		return nil
	}
}

// WithNetwork returns the option to set the network type (tcp, unix).
func WithNetwork(network string) Option {
	return func(m *mySQLClient) error {
		if network != "" {
			m.network = network
		}
		return nil
	}
}

// WithSocketPath returns the option to set the socketPath for unix domain socket connection.
func WithSocketPath(socketPath string) Option {
	return func(m *mySQLClient) error {
		if socketPath != "" {
			m.socketPath = socketPath
		}
		return nil
	}
}

// WithHost returns the option to set the host.
func WithHost(host string) Option {
	return func(m *mySQLClient) error {
		if host != "" {
			m.host = host
		}
		return nil
	}
}

// WithPort returns the option to set the port.
func WithPort(port uint16) Option {
	return func(m *mySQLClient) error {
		m.port = port
		return nil
	}
}

// WithUser returns the option to set the user.
func WithUser(user string) Option {
	return func(m *mySQLClient) error {
		if user != "" {
			m.user = user
		}
		return nil
	}
}

// WithPass returns the option to set the pass.
func WithPass(pass string) Option {
	return func(m *mySQLClient) error {
		if pass != "" {
			m.pass = pass
		}
		return nil
	}
}

// WithName returns the option to sst the name.
func WithName(name string) Option {
	return func(m *mySQLClient) error {
		if name != "" {
			m.name = name
		}
		return nil
	}
}

// WithInitialPingTimeLimit returns the option to set the initialPingTimeLimit.
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

// WithInitialPingDuration returns the option to set the initialPingDuration.
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

// WithConnectionLifeTimeLimit returns the option to set the connMaxLifeTime.
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

// WithMaxIdleConns returns the option to set the maxIdleConns.
// If conns is negative numner, no idle connections are retained.
// ref: https://golang.org/src/database/sql/sql.go?s=24983:25019#L879
func WithMaxIdleConns(conns int) Option {
	return func(m *mySQLClient) error {
		if conns != 0 {
			m.maxIdleConns = conns
		}
		return nil
	}
}

// WithMaxOpenConns returns the option to set the maxOpenConns.
// If conns is negative numner, no limit on the number of open connections.
// ref: https://golang.org/src/database/sql/sql.go?s=24983:25019#L923
func WithMaxOpenConns(conns int) Option {
	return func(m *mySQLClient) error {
		if conns != 0 {
			m.maxOpenConns = conns
		}
		return nil
	}
}

// WithTLSConfig returns the option to set the tlsConfig.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(m *mySQLClient) error {
		if cfg != nil {
			m.tlsConfig = cfg
		}
		return nil
	}
}

// WithDialer returns the option to set the dialer.
func WithDialer(der net.Dialer) Option {
	return func(m *mySQLClient) error {
		if der != nil {
			m.dialer = der
		}
		return nil
	}
}

// WithDialerFunc returns the option to set the dialer function.
func WithDialerFunc(der func(ctx context.Context, addr, port string) (net.Conn, error)) Option {
	return func(m *mySQLClient) error {
		if der != nil {
			m.dialerFunc = der
		}
		return nil
	}
}

// WithEventReceiver returns the option to set the eventReceiver.
func WithEventReceiver(er EventReceiver) Option {
	return func(m *mySQLClient) error {
		if er != nil {
			m.eventReceiver = er
		}

		return nil
	}
}
