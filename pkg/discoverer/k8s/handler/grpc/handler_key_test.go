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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// TestSingleflightKey_DistinctRequestsDoNotCollide pins that requests whose
// fields differ only in where a '-' falls (valid in Kubernetes names) get
// distinct singleflight keys, so one caller is never handed the result of a
// different namespace/name query.
func TestSingleflightKey_DistinctRequestsDoNotCollide(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key  func(string, *payload.Discoverer_Request) string
		a, b *payload.Discoverer_Request
		name string
	}{
		{
			name: "pods: namespace/name boundary",
			key:  singleflightKey,
			a:    &payload.Discoverer_Request{Namespace: "vald", Name: "agent-ngt"},
			b:    &payload.Discoverer_Request{Namespace: "vald-agent", Name: "ngt"},
		},
		{
			name: "pods: node/namespace boundary",
			key:  singleflightKey,
			a:    &payload.Discoverer_Request{Node: "node-a", Namespace: "ns", Name: "agent"},
			b:    &payload.Discoverer_Request{Node: "node", Namespace: "a-ns", Name: "agent"},
		},
		{
			name: "svcs: namespace/name boundary",
			key:  singleflightKeyForSvcs,
			a:    &payload.Discoverer_Request{Namespace: "vald", Name: "lb-gateway"},
			b:    &payload.Discoverer_Request{Namespace: "vald-lb", Name: "gateway"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ka, kb := tc.key(podPrefix, tc.a), tc.key(podPrefix, tc.b)
			if ka == kb {
				t.Errorf("distinct requests %v and %v share the singleflight key %q", tc.a, tc.b, ka)
			}
		})
	}
}
