//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/apis/grpc/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/vald/service"
)

type Server vald.ValdServer

type server struct {
	eg      errgroup.Group
	gateway service.ValdProxy
	timeout time.Duration
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, oid *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	return s.search(ctx, int(req.GetConfig().GetNum()), req.GetConfig().GetTimeout(), func(ac agent.AgentClient) (*payload.Search_Response, error) {
		return ac.Search(ctx, req)
	})
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	return s.search(ctx, int(req.GetConfig().GetNum()), req.GetConfig().GetTimeout(), func(ac agent.AgentClient) (*payload.Search_Response, error) {
		// TODO rewrite ObjectID
		return ac.SearchByID(ctx, req)
	})
}

func (s *server) search(ctx context.Context, num int, to int64, f func(ac agent.AgentClient) (*payload.Search_Response, error)) (res *payload.Search_Response, err error) {
	maxDist := uint32(math.MaxUint32)
	res.Results = make([]*payload.Object_Distance, 0, len(s.gateway.GetIPs())*num)
	dch := make(chan *payload.Object_Distance, cap(res.GetResults())/2)
	eg, ctx := errgroup.New(ctx)
	var cancel context.CancelFunc
	if to != 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(to))
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	eg.Go(safety.RecoverFunc(func() error {
		defer cancel()
		cl := new(checkList)
		return s.gateway.BroadCast(ctx, eg, func(ac agent.AgentClient) error {
			r, err := f(ac)
			if err != nil {
				return err
			}
			for _, dist := range r.GetResults() {
				id := dist.GetId().GetId()
				if dist.GetDistance() < math.Float32frombits(atomic.LoadUint32(&maxDist)) {
					if !cl.Exists(id) {
						dch <- dist
						cl.Check(id)
					}
					return nil
				}
			}
			return nil
		})
	}))
	for {
		select {
		case <-ctx.Done():
			err = eg.Wait()
			if len(res.GetResults()) > num && num != 0 {
				res.Results = res.Results[:num]
			}
			return res, err
		case dist := <-dch:
			pos := len(res.GetResults())
			if pos >= num &&
				dist.GetDistance() < math.Float32frombits(atomic.LoadUint32(&maxDist)) {
				atomic.StoreUint32(&maxDist, math.Float32bits(dist.GetDistance()))
			}
			for idx := pos; idx >= 0; idx-- {
				if res.GetResults()[idx].GetDistance() > dist.GetDistance() {
					pos = idx
					break
				}
			}
			if pos != 0 {
				res.Results = append(res.GetResults()[:pos+1], res.GetResults()[pos:]...)
				if len(res.GetResults()) > num && num != 0 {
					res.Results = res.GetResults()[:num]
				}
			}
			if pos <= num {
				res.Results[pos] = dist
			}
		}
	}
}

func (s *server) StreamSearch(stream vald.Vald_StreamSearchServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.Search(ctx, data.(*payload.Search_Request))
	})
}

func (s *server) StreamSearchByID(stream vald.Vald_StreamSearchByIDServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
	})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (*payload.Common_Error, error) {
	return nil, nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Common_Errors, err error) {
	return nil, nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (*payload.Common_Error, error) {
	return nil, nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Common_Errors, err error) {
	return nil, nil
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (*payload.Common_Error, error) {
	return nil, nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Common_Errors, err error) {
	return nil, nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (*payload.Object_Vector, error) {
	return nil, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}
