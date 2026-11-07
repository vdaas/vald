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

package grpc

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc/pool"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// fixedConnPool is a pool.Conn whose Get always returns conn.
type fixedConnPool struct {
	pool.Conn
	conn *ClientConn
}

func (p fixedConnPool) Get(context.Context) (*ClientConn, bool) {
	return p.conn, p.conn != nil
}

// Test_gRPCClient_executeRPC_noPerCallOptions verifies that executeRPC does not
// pass g.copts per call: they are already installed as the dial-level default
// call options of every pooled conn, and repeating them per call made grpc
// allocate a combined slice and apply every option twice for each RPC.
func Test_gRPCClient_executeRPC_noPerCallOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "with call options",
			opts: []Option{WithWaitForReady(true), WithMaxRecvMsgSize(1 << 20)},
		},
		{
			name: "without call options",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g, ok := New("test", tc.opts...).(*gRPCClient)
			if !ok {
				t.Fatal("New did not return *gRPCClient")
			}
			conn, err := ggrpc.NewClient("passthrough:///executeRPC-test",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if cerr := conn.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			got := -1
			_, _, err = g.executeRPC(t.Context(), fixedConnPool{conn: conn}, "addr",
				func(_ context.Context, _ *ClientConn, copts ...CallOption) (any, error) {
					got = len(copts)
					return struct{}{}, nil
				})
			if err != nil {
				t.Fatalf("executeRPC() err = %v", err)
			}
			if got != 0 {
				t.Errorf("f received %d per-call options, want 0 (defaults come from the conn)", got)
			}
		})
	}
}
