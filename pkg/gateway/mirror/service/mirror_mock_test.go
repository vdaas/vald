package service

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// MirrorMock represents mock struct for Gateway.
type MirrorMock struct {
	Mirror
	ConnectFunc         func(ctx context.Context, targets ...*payload.Mirror_Target) error
	DisconnectFunc      func(ctx context.Context, targets ...*payload.Mirror_Target) error
	IsConnectedFunc     func(ctx context.Context, addr string) bool
	RangeMirrorAddrFunc func(f func(addr string, _ any) bool)
}

// Connect calls ConnectFunc object.
func (mm *MirrorMock) Connect(ctx context.Context, targets ...*payload.Mirror_Target) error {
	return mm.ConnectFunc(ctx, targets...)
}

// Disconnect calls DisconnectFunc object.
func (mm *MirrorMock) Disconnect(ctx context.Context, targets ...*payload.Mirror_Target) error {
	return mm.DisconnectFunc(ctx, targets...)
}

// IsConnected calls IsConnectedFunc object.
func (mm *MirrorMock) IsConnected(ctx context.Context, addr string) bool {
	return mm.IsConnectedFunc(ctx, addr)
}

// RangeMirrorAddr calls RangeMirrorAddrFunc object.
func (mm *MirrorMock) RangeMirrorAddr(f func(addr string, _ any) bool) {
	mm.RangeMirrorAddrFunc(f)
}
