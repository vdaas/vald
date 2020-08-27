package cassandra

import (
	"context"
	"net"
)

type DialerMock struct {
	DialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error)
}

func (dm *DialerMock) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return dm.DialContextFunc(ctx, network, addr)
}
