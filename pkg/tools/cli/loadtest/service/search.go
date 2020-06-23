//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
package service

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

func newSearch() (requestFunc, loaderFunc) {
	return func(dataset assets.Dataset) ([]interface{}, error) {
			vectors := dataset.Query()
			requests := make([]interface{}, len(vectors))
			for j, v := range vectors {
				requests[j] = &payload.Search_Request{
					Vector: v,
				}
			}
			return requests, nil
		},
		func(ctx context.Context, c vald.ValdClient, i interface{}, copts ...grpc.CallOption) error {
			_, err := c.Search(ctx, i.(*payload.Search_Request), copts...)
			return err
		}
}
