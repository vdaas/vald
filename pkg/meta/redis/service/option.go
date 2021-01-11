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

// Package service manages the main logic of server.
package service

import (
	"github.com/vdaas/vald/internal/db/kvs/redis"
)

type Option func(*client) error

var defaultOptions = []Option{
	WithKVPrefix("kv"),
	WithVKPrefix("vk"),
	WithPrefixDelimiter("-"),
}

func WithRedisClient(r redis.Redis) Option {
	return func(c *client) error {
		if r != nil {
			c.db = r
		}

		return nil
	}
}

func WithRedisClientConnector(connector redis.Connector) Option {
	return func(c *client) error {
		if connector != nil {
			c.connector = connector
		}

		return nil
	}
}

func WithKVPrefix(name string) Option {
	return func(c *client) error {
		if name != "" {
			c.kvPrefix = name
		}

		return nil
	}
}

func WithVKPrefix(name string) Option {
	return func(c *client) error {
		if name != "" {
			c.vkPrefix = name
		}

		return nil
	}
}

func WithPrefixDelimiter(del string) Option {
	return func(c *client) error {
		if del != "" {
			c.prefixDelimiter = del
		}

		return nil
	}
}
