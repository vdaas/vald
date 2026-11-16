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

// hangingSearchClient answers GetObject at once but holds Search/SearchByID
// until the request deadline, i.e. agents too slow for the search timeout.
type hangingSearchClient struct {
	*regClient
}

func (hangingSearchClient) SearchByID(
	ctx context.Context, _ *payload.Search_IDRequest, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

// BenchmarkSearchByID_SearchTimeout measures SearchByID when every agent's
// search exceeds the per-request timeout (benchSearchTimeout): the fallback
// agent-side SearchByID broadcast must not wait out a second timeout.
func BenchmarkSearchByID_SearchTimeout(b *testing.B) {
	const (
		agents             = 3
		dim                = 8
		benchSearchTimeout = 50 * time.Millisecond
	)
	vec := make([]float32, dim)
	hang := func(ctx context.Context, _ *payload.Search_Request) (*payload.Search_Response, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	client := hangingSearchClient{&regClient{
		search: hang,
		getObject: func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{Id: fixtureUUID, Vector: vec}, nil
		},
	}}
	gw := &regGateway{
		agents: agents,
		addrs:  []string{loopbackIP},
		broadCast: func(ctx context.Context, f func(context.Context, string, vald.Client, ...grpc.CallOption) error) error {
			eg, ectx := errgroup.New(ctx)
			for i := range agents {
				target := "agent-" + strconv.Itoa(i)
				eg.Go(func() error {
					return f(ectx, target, client)
				})
			}
			return eg.Wait()
		},
	}
	eg, _ := errgroup.New(b.Context())
	s := &server{gateway: gw, eg: eg, replica: agents, timeout: benchSearchTimeout}
	req := &payload.Search_IDRequest{
		Id: fixtureUUID,
		Config: &payload.Search_Config{
			Num:     10,
			Timeout: int64(benchSearchTimeout),
		},
	}
	b.ResetTimer()
	for range b.N {
		if _, err := s.SearchByID(b.Context(), req); err == nil {
			b.Fatal("SearchByID() must fail when every agent times out")
		}
	}
}
