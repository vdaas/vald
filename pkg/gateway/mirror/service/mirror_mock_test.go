// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
