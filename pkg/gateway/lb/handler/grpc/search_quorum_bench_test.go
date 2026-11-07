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
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// newStragglerServer returns an LB server over `agents` fake agents that answer
// Search with a fresh random top-k after fast, except agent-0 which takes slow.
func newStragglerServer(agents, replica int, fast, slow time.Duration) *server {
	gw := &regGateway{
		agents: agents,
		addrs:  []string{loopbackIP},
		broadCast: func(ctx context.Context, f func(context.Context, string, vald.Client, ...grpc.CallOption) error) error {
			eg, ectx := errgroup.New(ctx)
			for i := range agents {
				target := "agent-" + strconv.Itoa(i)
				delay := fast
				if i == 0 {
					delay = slow
				}
				client := &regClient{search: func(ctx context.Context, _ *payload.Search_Request) (*payload.Search_Response, error) {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(delay):
						return newRandomResponse(), nil
					}
				}}
				eg.Go(func() error {
					return f(ectx, target, client)
				})
			}
			return eg.Wait()
		},
	}
	return &server{gateway: gw, replica: replica, timeout: time.Second}
}

// BenchmarkSearch_Straggler measures a Search over 6 agents (index replica 3)
// where one agent answers after 200ms and the others after 1ms, with the
// early termination disabled (default) and with a 10ms quorum grace.
func BenchmarkSearch_Straggler(b *testing.B) {
	const num = 10
	req := &payload.Search_Request{
		Vector: make([]float32, 8),
		Config: &payload.Search_Config{Num: num, Timeout: int64(time.Second)},
	}
	for _, grace := range []string{"", "10ms"} {
		name := "grace=" + grace
		if grace == "" {
			name = "grace=disabled"
		}
		b.Run(name, func(b *testing.B) {
			s := newStragglerServer(6, 3, time.Millisecond, 200*time.Millisecond)
			WithSearchQuorumGrace(grace)(s)
			b.ResetTimer()
			for range b.N {
				res, err := s.Search(b.Context(), req)
				if err != nil {
					b.Fatal(err)
				}
				if got := len(res.GetResults()); got != num {
					b.Fatalf("Search() returned %d results, want %d", got, num)
				}
			}
		})
	}
}
