//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package redis provides implementation of Go API for redis interface
package cassandra

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*client) error

var (
	defaultOpts = []Option{}
)

func WithHosts(hosts ...string) Option {
	return func(c *client) error {
		if len(hosts) == 0 {
			return nil
		}
		if c.hosts == nil {
			c.hosts = hosts
		} else {
			c.hosts = append(c.hosts, hosts...)
		}
		return nil
	}
}

func WithCQLVersion(version string) Option {
	return func(c *client) error {
		c.cqlVersion = version
		return nil
	}
}

func WithTimeout(dur string) Option {
	return func(c *client) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute // FIXME
		}
		c.timeout = d
		return nil
	}
}

func WithConnectTimeout(dur string) Option {
	return func(c *client) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		c.connectTimeout = d
		return nil
	}
}

func WithPort(port int) Option {
	return func(c *client) error {
		c.port = port
		return nil
	}
}

func WithNumConns(numConns int) Option {
	return func(c *client) error {
		c.numConns = numConns
		return nil
	}
}

var consistenciesMap = map[string]gocql.Consistency{
	"any":         gocql.Any,
	"one":         gocql.One,
	"two":         gocql.Two,
	"three":       gocql.Three,
	"quorum":      gocql.Quorum,
	"all":         gocql.All,
	"localQuorum": gocql.LocalQuorum,
	"eachQuorum":  gocql.EachQuorum,
	"localOne":    gocql.LocalOne,
}

func WithConsistency(consistency string) Option {
	return func(c *client) error {
		actual, ok := consistenciesMap[consistency]
		if !ok {
			return errors.ErrCassandraInvalidConsistencyType(consistency)
		}
		c.consistency = actual
		return nil
	}
}

func WithMaxPreparedStmts(maxPreparedStmts int) Option {
	return func(c *client) error {
		c.maxPreparedStmts = maxPreparedStmts
		return nil
	}
}

func WithMaxRoutingKeyInfo(maxRoutingKeyInfo int) Option {
	return func(c *client) error {
		c.maxRoutingKeyInfo = maxRoutingKeyInfo
		return nil
	}
}

func WithPageSize(pageSize int) Option {
	return func(c *client) error {
		c.pageSize = pageSize
		return nil
	}
}

func WithDefaultTimestamp(defaultTimestamp bool) Option {
	return func(c *client) error {
		c.defaultTimestamp = defaultTimestamp
		return nil
	}
}

func WithMaxWaitSchemaAgreement(maxWaitSchemaAgreement string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(maxWaitSchemaAgreement)
		if err != nil {
			return err
		}
		c.maxWaitSchemaAgreement = d
		return nil
	}
}

func WithReconnectInterval(reconnectInterval string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(reconnectInterval)
		if err != nil {
			return err
		}
		c.reconnectInterval = d
		return nil
	}
}

func WithReconnectionPolicyInitialInterval(initialInterval string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(initialInterval)
		if err != nil {
			return err
		}
		c.reconnectionPolicy.initialInterval = d
		return nil
	}
}

func WithReconnectionPolicyMaxRetries(maxRetries int) Option {
	return func(c *client) error {
		c.reconnectionPolicy.maxRetries = maxRetries
		return nil
	}
}

func WithWriteCoalesceWaitTime(writeCoalesceWaitTime string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(writeCoalesceWaitTime)
		if err != nil {
			return err
		}
		c.writeCoalesceWaitTime = d
		return nil
	}
}

func WithKeyspace(keyspace string) Option {
	return func(c *client) error {
		c.keyspace = keyspace
		return nil
	}
}

func WithKVTable(kvTable string) Option {
	return func(c *client) error {
		c.kvTable = kvTable
		return nil
	}
}

func WithVKTable(vkTable string) Option {
	return func(c *client) error {
		c.vkTable = vkTable
		return nil
	}
}

func WithUsername(username string) Option {
	return func(c *client) error {
		c.username = username
		return nil
	}
}

func WithPassword(password string) Option {
	return func(c *client) error {
		c.password = password
		return nil
	}
}
