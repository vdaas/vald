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
package quic

import (
	"context"
	"net"

	quic "github.com/quic-go/quic-go"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/tls"
)

type Conn struct {
	*quic.Conn
	*quic.Stream
}

type qconn struct {
	connectionCache sync.Map[string, *quic.Conn]
}

var defaultQconn = new(qconn)

func NewConn(ctx context.Context, conn *quic.Conn) (net.Conn, error) {
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn,
		stream,
	}, nil
}

func (c *Conn) Close() (err error) {
	return c.Stream.Close()
}

func DialContext(ctx context.Context, addr string, tcfg *tls.Config) (net.Conn, error) {
	if tcfg == nil {
		return nil, errors.ErrCertificationFailed
	}
	if len(tcfg.NextProtos) == 0 {
		return nil, errors.ErrEmptyALPNs
	}
	return defaultQconn.dialQuicContext(ctx, addr, tcfg)
}

func (q *qconn) dialQuicContext(
	ctx context.Context, addr string, tcfg *tls.Config,
) (net.Conn, error) {
	conn, ok := q.connectionCache.Load(addr)
	if ok {
		return NewConn(ctx, conn)
	}
	conn, err := quic.DialAddr(ctx, addr, tcfg, &quic.Config{
		/*
			// GetConfigForClient is called for incoming connections.
			// If the error is not nil, the connection attempt is refused.
			GetConfigForClient func(info *ClientHelloInfo) (*Config, error)
			// The QUIC versions that can be negotiated.
			// If not set, it uses all versions available.
			Versions []Version
			// HandshakeIdleTimeout is the idle timeout before completion of the handshake.
			// If we don't receive any packet from the peer within this time, the connection attempt is aborted.
			// Additionally, if the handshake doesn't complete in twice this time, the connection attempt is also aborted.
			// If this value is zero, the timeout is set to 5 seconds.
			HandshakeIdleTimeout time.Duration
			// MaxIdleTimeout is the maximum duration that may pass without any incoming network activity.
			// The actual value for the idle timeout is the minimum of this value and the peer's.
			// This value only applies after the handshake has completed.
			// If the timeout is exceeded, the connection is closed.
			// If this value is zero, the timeout is set to 30 seconds.
			MaxIdleTimeout time.Duration
			// The TokenStore stores tokens received from the server.
			// Tokens are used to skip address validation on future connection attempts.
			// The key used to store tokens is the ServerName from the tls.Config, if set
			// otherwise the token is associated with the server's IP address.
			TokenStore TokenStore
			// InitialStreamReceiveWindow is the initial size of the stream-level flow control window for receiving data.
			// If the application is consuming data quickly enough, the flow control auto-tuning algorithm
			// will increase the window up to MaxStreamReceiveWindow.
			// If this value is zero, it will default to 512 KB.
			// Values larger than the maximum varint (quicvarint.Max) will be clipped to that value.
			InitialStreamReceiveWindow uint64
			// MaxStreamReceiveWindow is the maximum stream-level flow control window for receiving data.
			// If this value is zero, it will default to 6 MB.
			// Values larger than the maximum varint (quicvarint.Max) will be clipped to that value.
			MaxStreamReceiveWindow uint64
			// InitialConnectionReceiveWindow is the initial size of the stream-level flow control window for receiving data.
			// If the application is consuming data quickly enough, the flow control auto-tuning algorithm
			// will increase the window up to MaxConnectionReceiveWindow.
			// If this value is zero, it will default to 512 KB.
			// Values larger than the maximum varint (quicvarint.Max) will be clipped to that value.
			InitialConnectionReceiveWindow uint64
			// MaxConnectionReceiveWindow is the connection-level flow control window for receiving data.
			// If this value is zero, it will default to 15 MB.
			// Values larger than the maximum varint (quicvarint.Max) will be clipped to that value.
			MaxConnectionReceiveWindow uint64
			// AllowConnectionWindowIncrease is called every time the connection flow controller attempts
			// to increase the connection flow control window.
			// If set, the caller can prevent an increase of the window. Typically, it would do so to
			// limit the memory usage.
			// To avoid deadlocks, it is not valid to call other functions on the connection or on streams
			// in this callback.
			AllowConnectionWindowIncrease func(conn Connection, delta uint64) bool
			// MaxIncomingStreams is the maximum number of concurrent bidirectional streams that a peer is allowed to open.
			// If not set, it will default to 100.
			// If set to a negative value, it doesn't allow any bidirectional streams.
			// Values larger than 2^60 will be clipped to that value.
			MaxIncomingStreams int64
			// MaxIncomingUniStreams is the maximum number of concurrent unidirectional streams that a peer is allowed to open.
			// If not set, it will default to 100.
			// If set to a negative value, it doesn't allow any unidirectional streams.
			// Values larger than 2^60 will be clipped to that value.
			MaxIncomingUniStreams int64
			// KeepAlivePeriod defines whether this peer will periodically send a packet to keep the connection alive.
			// If set to 0, then no keep alive is sent. Otherwise, the keep alive is sent on that period (or at most
			// every half of MaxIdleTimeout, whichever is smaller).
			KeepAlivePeriod time.Duration
			// InitialPacketSize is the initial size of packets sent.
			// It is usually not necessary to manually set this value,
			// since Path MTU discovery very quickly finds the path's MTU.
			// If set too high, the path might not support packets that large, leading to a timeout of the QUIC handshake.
			// Values below 1200 are invalid.
			InitialPacketSize uint16
			// DisablePathMTUDiscovery disables Path MTU Discovery (RFC 8899).
			// This allows the sending of QUIC packets that fully utilize the available MTU of the path.
			// Path MTU discovery is only available on systems that allow setting of the Don't Fragment (DF) bit.
			DisablePathMTUDiscovery bool
			// Allow0RTT allows the application to decide if a 0-RTT connection attempt should be accepted.
			// Only valid for the server.
			Allow0RTT bool
			// Enable QUIC datagram support (RFC 9221).
			EnableDatagrams bool
			Tracer          func(context.Context, logging.Perspective, ConnectionID) *logging.ConnectionTracer
		*/
	})
	if err != nil {
		return nil, err
	}
	q.connectionCache.Store(addr, conn)
	return NewConn(ctx, conn)
}

func (q *qconn) Close() (err error) {
	q.connectionCache.Range(func(addr string, conn *quic.Conn) bool {
		e := conn.CloseWithError(0, addr)
		if e != nil {
			err = errors.Wrap(err, e.Error())
		}
		return true
	})
	return nil
}
