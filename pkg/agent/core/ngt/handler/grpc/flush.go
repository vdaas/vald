// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) Flush(ctx context.Context, req *payload.Flush_Request) (*payload.Info_Index_Count, error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.FlushRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return nil, err
}
