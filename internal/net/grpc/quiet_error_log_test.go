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
)

// The quiet flag set by a fan-out caller must survive the contexts the client
// derives for each target (method wrapping, cancellation), and must not be set
// on an unmarked context.
func TestWithQuietErrorLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ctx  func() context.Context
		name string
		want bool
	}{
		{
			name: "unmarked context is not quiet",
			ctx:  context.Background,
			want: false,
		},
		{
			name: "marked context is quiet",
			ctx:  func() context.Context { return WithQuietErrorLog(context.Background()) },
			want: true,
		},
		{
			name: "derived contexts keep the flag",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(WithQuietErrorLog(context.Background()))
				t.Cleanup(cancel)
				return WrapGRPCMethod(ctx, "BroadCast/agent-0")
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := isQuietErrorLog(tc.ctx()); got != tc.want {
				t.Errorf("isQuietErrorLog() = %v, want %v", got, tc.want)
			}
		})
	}
}
