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
package service

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
)

type registerClientStub struct {
	MirrorClient
	res *payload.Mirror_Targets
}

func (c *registerClientStub) Register(
	_ context.Context, _ *payload.Mirror_Targets, _ ...grpc.CallOption,
) (*payload.Mirror_Targets, error) {
	return c.res, nil
}

// Test_mirr_registers_dedup verifies that every distinct target returned by the
// peers is collected exactly once, even if several peers report the same list.
func Test_mirr_registers_dedup(t *testing.T) {
	t.Parallel()
	peers := []string{"peer-1:8081", "peer-2:8081"}
	res := &payload.Mirror_Targets{
		Targets: []*payload.Mirror_Target{
			{Host: "192.168.0.1", Port: 8081},
			{Host: "192.168.0.2", Port: 8081},
		},
	}
	m := &mirr{
		gateway: &GatewayMock{
			DoMultiFunc: func(ctx context.Context, _ []string, f func(context.Context, string, MirrorClient, ...grpc.CallOption) error) error {
				for _, peer := range peers {
					if err := f(ctx, peer, &registerClientStub{res: res}); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	got, err := m.registers(t.Context(), &payload.Mirror_Targets{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(res.GetTargets()) {
		t.Fatalf("got %d targets, want %d: %v", len(got), len(res.GetTargets()), got)
	}
	seen := make(map[string]int, len(got))
	for _, tgt := range got {
		seen[tgt.GetHost()]++
	}
	for _, want := range res.GetTargets() {
		if n := seen[want.GetHost()]; n != 1 {
			t.Errorf("target %s appeared %d times, want 1", want.GetHost(), n)
		}
	}
}
