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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/apis/grpc/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/vald/service"
)

type server struct {
	eg       errgroup.Group
	gateway  service.Gateway
	metadata service.Meta
	backup   service.Backup
	timeout  time.Duration
	filters  []service.Filter
	replica  int
}

func New(opts ...Option) vald.ValdServer {
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
	return s.search(ctx, req.GetConfig(), func(ctx context.Context, ac agent.AgentClient) (*payload.Search_Response, error) {
		return ac.Search(ctx, req)
	})
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	val, err := s.metadata.GetMetaInverse(ctx, req.GetId().GetId())
	if err != nil {
		return nil, err
	}
	req.Id.Id = val.GetKey()
	return s.search(ctx, req.GetConfig(),
		func(ctx context.Context, ac agent.AgentClient) (*payload.Search_Response, error) {
			// TODO rewrite ObjectID
			return ac.SearchByID(ctx, req)
		})
}

func (s *server) search(ctx context.Context, cfg *payload.Search_Config,
	f func(ctx context.Context, ac agent.AgentClient) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error) {

	maxDist := uint32(math.MaxUint32)
	num := int(cfg.GetNum())
	to := cfg.GetTimeout()
	res.Results = make([]*payload.Object_Distance, 0, s.gateway.GetAgentCount()*num)
	dch := make(chan *payload.Object_Distance, cap(res.GetResults())/2)
	eg, ectx := errgroup.New(ctx)
	var cancel context.CancelFunc
	if to != 0 {
		ectx, cancel = context.WithTimeout(ectx, time.Duration(to))
	} else {
		ectx, cancel = context.WithCancel(ectx)
	}
	eg.Go(safety.RecoverFunc(func() error {
		defer cancel()
		cl := new(checkList)
		return s.gateway.BroadCast(ectx, func(ctx context.Context, target string, ac agent.AgentClient) error {
			r, err := f(ctx, ac)
			if err != nil {
				return err
			}
			for _, dist := range r.GetResults() {
				if dist.GetDistance() > math.Float32frombits(atomic.LoadUint32(&maxDist)) {
					return nil
				}
				id := dist.GetId().GetId()
				if !cl.Exists(id) {
					dch <- dist
					cl.Check(id)
				}
			}
			return nil
		})
	}))
	for {
		select {
		case <-ectx.Done():
			err = eg.Wait()
			close(dch)
			if len(res.GetResults()) > num && num != 0 {
				res.Results = res.Results[:num]
			}
			keys := make([]string, 0, len(res.Results))
			for _, r := range res.Results {
				keys = append(keys, r.GetId().GetId())
			}
			if s.metadata != nil {
				metas, err := s.metadata.GetMetas(ctx, keys...)
				if err == nil {
					for i, k := range metas {
						res.Results[i].Id = &payload.Object_ID{
							Id: k,
						}
					}
				}
			}
			if s.filters != nil {
				for _, filter := range s.filters {
					res, err = filter.FilterSearch(res)
					if err != nil {
						return res, err
					}
				}
			}
			// TODO metadataを引いてFilterに転送
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
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Search(ctx, data.(*payload.Search_Request))
		})
}

func (s *server) StreamSearchByID(stream vald.Vald_StreamSearchByIDServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
		})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (ce *payload.Empty, err error) {
	uuid := fuid.String()
	err = s.metadata.SetMeta(ctx, vec.Id.GetId(), uuid)
	if err != nil {
		return nil, err
	}
	vec.Id = &payload.Object_ID{
		Id: uuid,
	}
	mu := new(sync.Mutex)
	targets := make([]string, 0, s.replica)
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, ac agent.AgentClient) (err error) {
		_, err = ac.Insert(ctx, vec)
		if err != nil {
			return err
		}
		mu.Lock()
		targets = append(targets, target)
		mu.Unlock()
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = s.backup.Register(vec.GetId().GetId(), targets...)
	return nil, err
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	for _, vec := range vecs.GetVectors() {
		vec.Id = &payload.Object_ID{
			Id: s.metadata.SetMeta(vec.Id.GetId()),
		}
		_, err := s.Insert(ctx, vec)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (*payload.Empty, error) {
	vec.Id = &payload.Object_ID{
		Id: s.metadata.GetMeta(vec.Id.GetId()),
	}
	locs, err := s.backup.GetLocation(vec.GetId().GetId())
	if err != nil {
		return nil, err
	}
	lmap := make(map[string]struct{}, len(locs))
	for _, loc := range locs {
		lmap[loc] = struct{}{}
	}
	s.gateway.BroadCast(ctx, func(ctx context.Context, target string, ac agent.AgentClient) error {
		_, ok := lmap[target]
		if ok {
			_, err = ac.Update(ctx, vec)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil, nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	for _, vec := range vecs.GetVectors() {
		_, err := s.Update(ctx, vec)
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (*payload.Empty, error) {
	uuid := s.metadata.GetMeta(id.GetId())
	locs, err := s.backup.GetLocation(uuid)
	if err != nil {
		return nil, err
	}
	lmap := make(map[string]struct{}, len(locs))
	for _, loc := range locs {
		lmap[loc] = struct{}{}
	}
	s.gateway.BroadCast(ctx, func(ctx context.Context, target string, ac agent.AgentClient) error {
		_, ok := lmap[target]
		if ok {
			_, err = ac.Remove(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil, nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Empty, err error) {
	for _, id := range ids.GetIds() {
		_, err := s.Remove(ctx, &payload.Object_ID{
			Id: s.metadata.GetMeta(id.GetId()),
		})
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (*payload.Object_Vector, error) {
	// TODO get Object from backup
	return nil, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.GetObject(ctx, data.(*payload.Object_ID))
		})
}
