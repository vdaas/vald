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
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// Test_filterAddrs pins that filterAddrs lists the configured addrs first and
// formats each request target exactly as fmt.Sprintf("%s:%d") did (an IPv6
// host is not bracketed), with no spare capacity left to regrow into.
func Test_filterAddrs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		configured []string
		targets    []*payload.Filter_Target
	}{
		{name: "nothing"},
		{name: "configured only", configured: []string{"a:1", "b:2"}},
		{
			name: "targets only",
			targets: []*payload.Filter_Target{
				{Host: "filter", Port: 8081},
				{Host: "::1", Port: 0},
				{Host: "", Port: math.MaxUint32},
				nil,
			},
		},
		{
			name:       "configured then targets",
			configured: []string{"configured:8081"},
			targets:    []*payload.Filter_Target{{Host: "10.0.0.1", Port: 80}, {Host: "fe80::1%eth0", Port: 65535}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			want := slices.Clone(tc.configured)
			for _, target := range tc.targets {
				want = append(want, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
			}
			got := filterAddrs(tc.configured, tc.targets)
			if !slices.Equal(got, want) {
				tt.Errorf("filterAddrs() = %q, want %q", got, want)
			}
			if cap(got) != len(want) {
				tt.Errorf("cap(filterAddrs()) = %d, want %d", cap(got), len(want))
			}
		})
	}
}
