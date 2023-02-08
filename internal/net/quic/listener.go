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

	quic "github.com/lucas-clemente/quic-go"
)

type Listener struct {
	quic.Listener

	ctx context.Context
}

func Listen(ctx context.Context, addr string, tcfg *tls.Config) (net.Listener, error) {
	ql, err := quic.ListenAddr(addr, tcfg, &quic.Config{
		// Versions:             nil,
		// ConnectionIDLength:   0,
		// HandshakeIdleTimeout: 0,
		// MaxIdleTimeout:       0,
		// AcceptToken: func(clientAddr net.Addr, token *quic.Token) bool {
		// 	return true
		// },
		// TokenStore:                     quic.NewLRUTokenStore(clientAddr),
		// InitialStreamReceiveWindow:     0,
		// InitialConnectionReceiveWindow: 0,
		// MaxStreamReceiveWindow:         0,
		// MaxConnectionReceiveWindow:     0,
		// MaxIncomingStreams:             0,
		// MaxIncomingUniStreams:          0,
		// StatelessResetKey:              nil,
		// KeepAlive:                      true,
		// DisablePathMTUDiscovery:        false,
		EnableDatagrams: true,
		// Tracer:                         logging.NewMultiplexedTracer(),
	})
	if err != nil {
		return nil, err
	}
	return &Listener{
		Listener: ql,
		ctx:      ctx,
	}, nil
}

func (l *Listener) Accept() (net.Conn, error) {
	sess, err := l.Listener.Accept(l.ctx)
	if err != nil {
		return nil, err
	}

	stream, err := sess.AcceptStream(l.ctx)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Session: sess,
		Stream:  stream,
	}, nil
}
