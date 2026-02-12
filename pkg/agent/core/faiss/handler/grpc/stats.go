//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	stats "github.com/vdaas/vald/internal/net/grpc/stats"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) ResourceStatsDetail(
	ctx context.Context, _ *payload.Empty,
) (res *payload.Info_Stats_ResourceStatsDetail, err error) {
	_, span := trace.StartSpan(ctx, apiName+".ResourceStatsDetail")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	rs, err := stats.GetResourceStats(ctx)
	if err != nil {
		return nil, err
	}

	return &payload.Info_Stats_ResourceStatsDetail{
		Details: map[string]*payload.Info_Stats_ResourceStats{
			s.name: rs,
		},
	}, nil
}
