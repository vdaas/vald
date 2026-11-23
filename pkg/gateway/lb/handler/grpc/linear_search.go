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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

func (s *server) LinearSearch(
	ctx context.Context, req *payload.Search_Request,
) (*payload.Search_Response, error) {
	return s.searchVector(ctx, vald.LinearSearchRPCName, req, vald.Client.LinearSearch)
}

func (s *server) LinearSearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (*payload.Search_Response, error) {
	return s.searchByID(ctx, vald.LinearSearchByIDRPCName, vald.LinearSearchRPCName, req, vald.Client.LinearSearchByID, vald.Client.LinearSearch)
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) error {
	return streamSearch(s, stream, vald.StreamLinearSearchRPCName, func(req *payload.Search_Request) string {
		return "requestID-" + req.GetConfig().GetRequestId()
	}, s.LinearSearch)
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) error {
	return streamSearch(s, stream, vald.StreamLinearSearchByIDRPCName, func(req *payload.Search_IDRequest) string {
		return "id-" + req.GetId()
	}, s.LinearSearchByID)
}

func (s *server) MultiLinearSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (*payload.Search_Responses, error) {
	return multiSearch(ctx, s, vald.MultiLinearSearchRPCName, reqs.GetRequests(), s.LinearSearch)
}

func (s *server) MultiLinearSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (*payload.Search_Responses, error) {
	return multiSearch(ctx, s, vald.MultiLinearSearchByIDRPCName, reqs.GetRequests(), s.LinearSearchByID)
}
