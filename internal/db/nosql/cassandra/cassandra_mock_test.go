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

type CompressorMock struct {
	NameFunc   func() string
	EncodeFunc func(data []byte) ([]byte, error)
	DecodeFunc func(data []byte) ([]byte, error)
}

func (cm *CompressorMock) Name() string {
	return cm.NameFunc()
}
func (cm *CompressorMock) Encode(data []byte) ([]byte, error) {
	return cm.EncodeFunc(data)
}
func (cm *CompressorMock) Decode(data []byte) ([]byte, error) {
	return cm.DecodeFunc(data)
}
