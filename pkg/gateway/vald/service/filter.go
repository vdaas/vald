//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"runtime"
	"sync"
	"time"

	egress "github.com/vdaas/vald/apis/grpc/filter/egress"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Filter interface {
	Start(ctx context.Context) <-chan error
	FilterSearch(context.Context, *payload.Search_Response) (*payload.Search_Response, error)
}

type filters []*grpc.ClientConn

type filter struct {
	hcDur   time.Duration
	addrs   []string
	eg      errgroup.Group
	filters sync.Map
	bo      backoff.Backoff
	gopts   []grpc.DialOption
	copts   []grpc.CallOption
}

func NewFilter(opts ...FilterOption) (ef Filter, err error) {
	f := new(filter)
	for _, opt := range append(defaultFilterOpts, opts...) {
		err = opt(f)

		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *filter) Start(ctx context.Context) <-chan error {
	ech := make(chan error, len(f.addrs))
	if f.addrs == nil || len(f.addrs) == 0 {
		defer close(ech)
		return ech
	}
	invalidFilters := make([]int, 0, len(f.addrs))
	for i, addr := range f.addrs {
		conn, err := grpc.DialContext(ctx, addr,
			append(f.gopts, grpc.WithBlock())...)
		if err != nil {
			invalidFilters = append(invalidFilters, i)
			ech <- err
		} else {
			f.filters.Store(addr, conn)
		}
	}

	f.eg.Go(safety.RecoverFunc(func() (err error) {
		tick := time.NewTicker(f.hcDur)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				close(ech)
				return ctx.Err()
			case <-tick.C:
				f.filters.Range(func(iaddr, iconn interface{}) bool {
					conn, ok := iconn.(*grpc.ClientConn)
					if !ok {
						return true
					}
					addr, ok := iaddr.(string)

					if conn == nil ||
						conn.GetState() == connectivity.Shutdown ||
						conn.GetState() == connectivity.TransientFailure {
						if conn != nil {
							err = conn.Close()
							if err != nil {
								ech <- err
							}
						}
						conn, err = grpc.DialContext(ctx, addr, f.gopts...)
						if err != nil {
							ech <- err
							runtime.Gosched()
						} else {
							f.filters.Store(addr, conn)
						}
					}
					return true
				})
			}
		}
		return nil
	}))
	return ech
}

func (f *filter) FilterSearch(ctx context.Context, res *payload.Search_Response) (*payload.Search_Response, error) {

	var rerr error

	f.filters.Range(func(iaddr, iconn interface{}) bool {
		conn, ok := iconn.(*grpc.ClientConn)
		if !ok {
			return true
		}
		addr, ok := iaddr.(string)
		if !ok {
			return true
		}
		r, err := egress.NewEgressFilterClient(conn).Filter(ctx, res, f.copts...)
		if err != nil {
			rerr = errors.Wrap(rerr, fmt.Sprintf("addr: %s\terror: %s", addr, err.Error()))
		} else {
			res = r
		}
		return true
	})

	return res, rerr
}
