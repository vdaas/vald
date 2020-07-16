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

// Package service manages the main logic of server.
package service

import (
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
)

type Option func(*client) error

var (
	defaultOpts = []Option{
		WithKVTable("kv"),
		WithVKTable("vk"),
	}
)

func WithCassandraOpts(opts ...cassandra.Option) Option {
	return func(c *client) error {
		if c.cassandraOpts == nil {
			c.cassandraOpts = opts
			return nil
		}

		c.cassandraOpts = append(c.cassandraOpts, opts...)

		return nil
	}
}

func WithKVTable(name string) Option {
	return func(c *client) error {
		if name != "" {
			c.kvTable = name
		}

		return nil
	}
}
func WithVKTable(name string) Option {
	return func(c *client) error {
		if name != "" {
			c.vkTable = name
		}

		return nil
	}
}
