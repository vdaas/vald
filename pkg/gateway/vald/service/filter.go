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
	"reflect"

	"github.com/vdaas/vald/apis/grpc/filter/egress"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Filter interface {
	Start(ctx context.Context) (<-chan error, error)
	FilterSearch(context.Context, *payload.Search_Response) (*payload.Search_Response, error)
}

type filter struct {
	client grpc.Client
}

func NewFilter(opts ...FilterOption) (ef Filter, err error) {
	f := new(filter)
	for _, opt := range append(defaultFilterOpts, opts...) {
		if err := opt(f); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return f, nil
}

func (f *filter) Start(ctx context.Context) (<-chan error, error) {
	return f.client.StartConnectionMonitor(ctx)
}

func (f *filter) FilterSearch(ctx context.Context, res *payload.Search_Response) (*payload.Search_Response, error) {
	var rerr error
	f.client.Range(ctx,
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			r, err := egress.NewEgressFilterClient(conn).Filter(ctx, res, copts...)
			if err != nil {
				rerr = errors.Wrap(rerr, err.Error())
			} else {
				res = r
			}
			return nil
		})

	return res, rerr
}
