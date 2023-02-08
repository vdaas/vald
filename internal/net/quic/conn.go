// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package quic

import (
	"context"
	"crypto/tls"
	"net"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/vdaas/vald/internal/errors"
)

type Conn struct {
	quic.Connection
	quic.Stream
}

type qconn struct {
	connectionCache sync.Map
}

var defaultQconn = new(qconn)

func NewConn(ctx context.Context, conn quic.Connection) (net.Conn, error) {
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Connection: conn,
		Stream:     stream,
	}, nil
}

func (c *Conn) Close() (err error) {
	return c.Stream.Close()
}

func DialQuicContext(ctx context.Context, addr string, tcfg *tls.Config) (net.Conn, error) {
	return defaultQconn.dialQuicContext(ctx, addr, tcfg)
}

func (q *qconn) dialQuicContext(ctx context.Context, addr string, tcfg *tls.Config) (net.Conn, error) {
	si, ok := q.connectionCache.Load(addr)
	if ok {
		if conn, ok := si.(quic.Connection); ok {
			return NewConn(ctx, conn)
		}
	}
	conn, err := quic.DialAddrContext(ctx, addr, tcfg, nil)
	if err != nil {
		return nil, err
	}
	q.connectionCache.Store(addr, conn)
	return NewConn(ctx, conn)
}

func (q *qconn) Close() (err error) {
	q.connectionCache.Range(func(addr, si interface{}) bool {
		if conn, ok := si.(quic.Connection); ok {
			e := conn.CloseWithError(0, addr.(string))
			if e != nil {
				err = errors.Wrap(err, e.Error())
			}
		}
		return true
	})
	return nil
}
