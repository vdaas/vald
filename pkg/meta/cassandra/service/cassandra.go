//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

// Package service manages the main logic of server.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/config"
)

type Cassandra interface {
	Connect(context.Context) error
	Disconnect() error
	Get(string) (string, error)
	GetMultiple(...string) ([]string, error)
	GetInverse(string) (string, error)
	GetInverseMultiple(...string) ([]string, error)
	Set(string, string) error
	SetMultiple(map[string]string) error
	Delete(string) (string, error)
	DeleteMultiple(...string) ([]string, error)
	DeleteInverse(string) (string, error)
	DeleteInverseMultiple(...string) ([]string, error)
}

type client struct {
	// db              cassandra.Cassandra
	// opts            []cassandra.Option
	// topts           []tcp.DialerOption
	kvPrefix        string // TODO cassandraの場合はprefixではなくtable 変えてもいいかもね
	vkPrefix        string
	prefixDelimiter string
}

func New(cfg *config.Cassandra) (Cassandra, error) {
	c := new(client)
	return c, nil
}

func (c *client) Disconnect() error {
	return nil
}

func (c *client) Connect(ctx context.Context) error {
	return nil
}

func (c *client) Get(key string) (string, error) {
	return c.get(c.kvPrefix, key)
}

func (c *client) GetMultiple(keys ...string) (vals []string, err error) {
	return c.getMulti(c.kvPrefix, keys...)
}

func (c *client) GetInverse(val string) (string, error) {
	return c.get(c.vkPrefix, val)
}

func (c *client) GetInverseMultiple(vals ...string) ([]string, error) {
	return c.getMulti(c.vkPrefix, vals...)
}

func (c *client) appendPrefix(prefix, key string) string {
	return prefix + c.prefixDelimiter + key
}

func (c *client) get(prefix, key string) (string, error) {
	return "", nil
}

func (c *client) getMulti(prefix string, keys ...string) (vals []string, err error) {
	return nil, nil
}

func (c *client) Set(key, val string) error {
	return nil
}

func (c *client) SetMultiple(kvs map[string]string) (err error) {
	return nil
}

func (c *client) Delete(key string) (string, error) {
	return c.delete(c.kvPrefix, c.vkPrefix, key)
}

func (c *client) DeleteMultiple(keys ...string) ([]string, error) {
	return c.deleteMulti(c.kvPrefix, c.vkPrefix, keys...)
}

func (c *client) DeleteInverse(val string) (string, error) {
	return c.delete(c.vkPrefix, c.kvPrefix, val)
}

func (c *client) DeleteInverseMultiple(vals ...string) ([]string, error) {
	return c.deleteMulti(c.vkPrefix, c.kvPrefix, vals...)
}

func (c *client) delete(pfx, pfxInv, key string) (val string, err error) {
	return "", nil
}

func (c *client) deleteMulti(pfx, pfxInv string, keys ...string) (vals []string, err error) {
	return nil, nil
}
