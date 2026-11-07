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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func Test_skipSearchByIDFallback(t *testing.T) {
	t.Parallel()
	tests := []struct {
		err      error
		name     string
		canceled bool
		want     bool
	}{
		{name: "DeadlineExceeded status skips the fallback", err: status.Error(codes.DeadlineExceeded, "timeout"), want: true},
		{name: "Canceled status skips the fallback", err: status.Error(codes.Canceled, "canceled"), want: true},
		{name: "done request ctx skips the fallback", canceled: true, err: status.Error(codes.Internal, "internal"), want: true},
		{name: "NotFound still falls back", err: status.Error(codes.NotFound, "not found")},
		{name: "Internal still falls back", err: status.Error(codes.Internal, "internal")},
		{name: "Unavailable still falls back", err: status.Error(codes.Unavailable, "unavailable")},
		{name: "non-status error still falls back", err: errors.New("plain")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			if tc.canceled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			if got := skipSearchByIDFallback(ctx, tc.err); got != tc.want {
				t.Errorf("skipSearchByIDFallback() = %v, want %v", got, tc.want)
			}
		})
	}
}
