//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
//

package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	return s.UnimplementedValdServer.UnimplementedSearchServer.LinearSearch(ctx, req)
}

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	return s.UnimplementedValdServer.UnimplementedSearchServer.LinearSearchByID(ctx, req)
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	return s.UnimplementedValdServer.UnimplementedSearchServer.StreamLinearSearch(stream)
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) (err error) {
	return s.UnimplementedValdServer.UnimplementedSearchServer.StreamLinearSearchByID(stream)
}

func (s *server) MultiLinearSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	return s.UnimplementedValdServer.UnimplementedSearchServer.MultiLinearSearch(ctx, reqs)
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	return s.UnimplementedValdServer.UnimplementedSearchServer.MultiLinearSearchByID(ctx, reqs)
}
