// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vald

import (
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/connectivity"
	"github.com/vdaas/vald/internal/sync"
)

const connCachePruneMinSize = 8

// ConnCache caches per-*grpc.ClientConn vald.Client wrappers lock-freely and
// opportunistically prunes Shutdown connections when new entries are inserted.
type ConnCache struct {
	m sync.Map[*grpc.ClientConn, vald.Client]
}

// Get returns a cached vald.Client for conn, creating one on a miss.
// A nil conn returns the shared nil-safe Client.
func (c *ConnCache) Get(conn *grpc.ClientConn) vald.Client {
	if conn == nil {
		return nilSingleClient
	}
	if cli, ok := c.m.Load(conn); ok {
		return cli
	}
	return c.getOrInsert(conn)
}

// Load returns the cached vald.Client for conn if present.
func (c *ConnCache) Load(conn *grpc.ClientConn) (vald.Client, bool) {
	if conn == nil {
		return nil, false
	}
	return c.m.Load(conn)
}

func (c *ConnCache) getOrInsert(conn *grpc.ClientConn) vald.Client {
	cli, loaded := c.m.LoadOrCompute(conn, func() vald.Client {
		return vald.NewValdClient(conn)
	})
	if !loaded && c.m.Len() >= connCachePruneMinSize {
		c.m.DeleteFunc(func(k *grpc.ClientConn, _ vald.Client) bool {
			return k != conn && k.Target() != "" && k.GetState() == connectivity.Shutdown
		})
	}
	return cli
}
