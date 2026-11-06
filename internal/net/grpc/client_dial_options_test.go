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
	"testing"

	"google.golang.org/grpc"
)

// TestGRPCClient_mergedDialOptions verifies that per-call dial options never
// leak into the shared g.dopts backing array, which would let concurrent
// Connect calls overwrite each other's options.
func TestGRPCClient_mergedDialOptions(t *testing.T) {
	t.Parallel()
	base := grpc.WithUserAgent("base")
	first := grpc.WithUserAgent("first")
	second := grpc.WithUserAgent("second")

	tests := []struct {
		name string
		// spare is the spare capacity of g.dopts.
		spare int
	}{
		{name: "g.dopts without spare capacity", spare: 0},
		{name: "g.dopts with spare capacity", spare: 4},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dopts := make([]DialOption, 1, 1+test.spare)
			dopts[0] = base
			g := &gRPCClient{dopts: dopts}

			got1 := g.mergedDialOptions([]DialOption{first})
			got2 := g.mergedDialOptions([]DialOption{second})

			if len(got1) != 2 || got1[0] != base || got1[1] != first {
				t.Errorf("first call result was modified or wrong: %v", got1)
			}
			if len(got2) != 2 || got2[0] != base || got2[1] != second {
				t.Errorf("second call result is wrong: %v", got2)
			}
			if len(g.dopts) != 1 {
				t.Errorf("g.dopts length changed: %d", len(g.dopts))
			}
			if full := g.dopts[:cap(g.dopts)]; test.spare > 0 && (full[1] != nil || full[2] != nil) {
				t.Errorf("g.dopts spare capacity was written: %v", full)
			}
		})
	}
}
