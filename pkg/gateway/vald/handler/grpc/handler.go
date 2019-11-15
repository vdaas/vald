//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/apis/grpc/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
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
	filter   service.Filter
	replica  int
}

func New(opts ...Option) vald.ValdServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (*payload.Object_ID, error) {
	uuid, err := s.metadata.GetUUID(ctx, meta.GetId())
	return &payload.Object_ID{
		Id: uuid,
	}, err
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	return s.search(ctx, req.GetConfig(), func(ctx context.Context, ac agent.AgentClient) (*payload.Search_Response, error) {
		return ac.Search(ctx, req)
	})
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	meta := req.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		return nil, err
	}
	req.Id = uuid
	return s.search(ctx, req.GetConfig(),
		func(ctx context.Context, ac agent.AgentClient) (*payload.Search_Response, error) {
			// TODO rewrite ObjectID
			meta := req.GetId()
			uuid, err := s.metadata.GetUUID(ctx, meta)
			if err != nil {
				return nil, err
			}
			req.Id = uuid
			return ac.SearchByID(ctx, req)
		})
}

func (s *server) search(ctx context.Context, cfg *payload.Search_Config,
	f func(ctx context.Context, ac agent.AgentClient) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error) {

	maxDist := uint32(math.MaxUint32)
	num := int(cfg.GetNum())

	res.Results = make([]*payload.Object_Distance, 0, s.gateway.GetAgentCount()*num)
	dch := make(chan *payload.Object_Distance, cap(res.GetResults())/2)
	eg, ectx := errgroup.New(ctx)
	var cancel context.CancelFunc
	var timeout time.Duration
	if to := cfg.GetTimeout(); to != 0 {
		timeout = time.Duration(to)
	} else {
		timeout = s.timeout
	}
	ectx, cancel = context.WithTimeout(ectx, timeout)

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
				id := dist.GetId()
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
			uuids := make([]string, 0, len(res.Results))
			for _, r := range res.Results {
				uuids = append(uuids, r.GetId())
			}
			if s.metadata != nil {
				metas, err := s.metadata.GetMetas(ctx, uuids...)
				if err == nil {
					for i, k := range metas {
						res.Results[i].Id = k
					}
				}
			}
			if s.filter != nil {
				r, err := s.filter.FilterSearch(ctx, res)
				if err != nil {
					return res, err
				}
				res = r
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
	log.Debug(vec)
	uuid := fuid.String()
	meta := vec.GetId()
	err = s.metadata.SetUUIDandMeta(ctx, uuid, meta)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	vec.Id = uuid
	mu := new(sync.Mutex)
	targets := make([]string, 0, s.replica)
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, ac agent.AgentClient) (err error) {
		_, err = ac.Insert(ctx, vec)
		if err != nil {
			log.Error(err)
			return err
		}
		target = strings.SplitN(target, ":", 2)[0]
		mu.Lock()
		targets = append(targets, target)
		mu.Unlock()
		return nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if s.backup != nil {
		err = s.backup.Register(ctx, &payload.Backup_MetaVector{
			Uuid:   uuid,
			Meta:   meta,
			Vector: vec.GetVector(),
			Ips:    targets,
		})
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return new(payload.Empty), nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	metaMap := make(map[string]string)
	for i, vec := range vecs.GetVectors() {
		uuid := fuid.String()
		metaMap[uuid] = vec.GetId()
		vecs.Vectors[i].Id = uuid
	}
	mu := new(sync.Mutex)
	targets := make([]string, 0, s.replica)
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, ac agent.AgentClient) (err error) {
		_, err = ac.MultiInsert(ctx, vecs)
		if err != nil {
			return err
		}
		target = strings.SplitN(target, ":", 2)[0]
		mu.Lock()
		targets = append(targets, target)
		mu.Unlock()
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = s.metadata.SetUUIDandMetas(ctx, metaMap)
	if err != nil {
		return nil, err
	}

	if s.backup != nil {
		mvecs := new(payload.Backup_MetaVectors)
		mvecs.Vectors = make([]*payload.Backup_MetaVector, 0, len(vecs.GetVectors()))
		for _, vec := range vecs.GetVectors() {
			uuid := vec.GetId()
			mvecs.Vectors = append(mvecs.Vectors, &payload.Backup_MetaVector{
				Uuid:   uuid,
				Meta:   metaMap[uuid],
				Vector: vec.GetVector(),
				Ips:    targets,
			})
		}
		err = s.backup.RegisterMultiple(ctx, mvecs)
		if err != nil {
			return nil, err
		}
	}
	return new(payload.Empty), nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (res *payload.Empty, err error) {
	meta := vec.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		return nil, err
	}
	vec.Id = uuid
	locs, err := s.backup.GetLocation(ctx, uuid)
	if err != nil {
		return nil, err
	}
	lmap := make(map[string]struct{}, len(locs))
	for _, loc := range locs {
		lmap[loc] = struct{}{}
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, ac agent.AgentClient) error {
		target = strings.SplitN(target, ":", 2)[0]
		_, ok := lmap[target]
		if ok {
			_, err = ac.Update(ctx, vec)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = s.backup.Register(ctx, &payload.Backup_MetaVector{
		Uuid:   uuid,
		Meta:   meta,
		Vector: vec.GetVector(),
		Ips:    locs,
	})

	return new(payload.Empty), nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	ids := make([]string, 0, len(vecs.GetVectors()))
	for _, vec := range vecs.GetVectors() {
		ids = append(ids, vec.GetId())
	}
	_, err = s.MultiRemove(ctx, &payload.Object_IDs{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	_, err = s.MultiInsert(ctx, vecs)
	if err != nil {
		return nil, err
	}
	return new(payload.Empty), nil
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (*payload.Empty, error) {
	meta := id.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		return nil, err
	}
	locs, err := s.backup.GetLocation(ctx, uuid)
	if err != nil {
		return nil, err
	}
	lmap := make(map[string]struct{}, len(locs))
	for _, loc := range locs {
		lmap[loc] = struct{}{}
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, ac agent.AgentClient) error {
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
	if err != nil {
		return nil, err
	}
	_, err = s.metadata.DeleteMeta(ctx, uuid)
	if err != nil {
		return nil, err
	}
	err = s.backup.Remove(ctx, uuid)
	return new(payload.Empty), nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Empty, err error) {
	uuids, err := s.metadata.GetUUIDs(ctx, ids.GetIds()...)
	if err != nil {
		return nil, err
	}
	lmap := make(map[string][]string, s.gateway.GetAgentCount())
	for _, uuid := range uuids {
		locs, err := s.backup.GetLocation(ctx, uuid)
		if err != nil {
			return nil, err
		}
		for _, loc := range locs {
			lmap[loc] = append(lmap[loc], uuid)
		}
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, ac agent.AgentClient) error {
		uuids, ok := lmap[target]
		if ok {
			_, err := ac.MultiRemove(ctx, &payload.Object_IDs{
				Ids: uuids,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	_, err = s.metadata.DeleteMetas(ctx, uuids...)
	if err != nil {
		return nil, err
	}
	err = s.backup.RemoveMultiple(ctx, uuids...)
	if err != nil {
		return nil, err
	}
	return new(payload.Empty), nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (vec *payload.Backup_MetaVector, err error) {
	meta := id.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		return nil, err
	}
	vec, err = s.backup.GetObject(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.GetObject(ctx, data.(*payload.Object_ID))
		})
}
