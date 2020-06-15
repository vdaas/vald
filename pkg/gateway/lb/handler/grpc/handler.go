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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway
	timeout           time.Duration
	strict            bool
	replica           int
	streamConcurrency int
}

func New(opts ...Option) vald.ValdServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	var once sync.Once
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Exists/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		oid, err := vc.Exists(ctx, meta, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil
		}
		if oid != nil && oid.Id != "" {
			once.Do(func() {
				id = new(payload.Object_ID)
				id.Id = oid.Id
				cancel()
			})
		}
		return nil
	})
	if err != nil || id == nil || id.Id == "" {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err, meta.GetId(), info.Get())
	}
	return id, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(req.Vector) < 2 {
		err = errors.ErrInvalidDimensionSize(len(req.Vector), 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Search API invalid vector argument", err, req, info.Get())
	}
	return s.search(ctx, req.GetConfig(),
		func(ctx context.Context, vc vald.ValdClient, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.Search(ctx, req, copts...)
		})
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(req.GetId()) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("SearchByID API invalid uuid", err, req, info.Get())
	}
	vec, err := s.GetObject(ctx, &payload.Object_ID{
		Id: req.GetId(),
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's object not found", req.GetId()), err, info.Get())
	}
	return s.search(ctx, req.GetConfig(),
		func(ctx context.Context, vc vald.ValdClient, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.Search(ctx, &payload.Search_Request{
				Vector: vec.GetVector(),
				Config: req.GetConfig(),
			}, copts...)
		})
}

func (s *server) search(ctx context.Context, cfg *payload.Search_Config,
	f func(ctx context.Context, vc vald.ValdClient, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	maxDist := uint32(math.MaxUint32)
	num := int(cfg.GetNum())
	res = new(payload.Search_Response)
	res.Results = make([]*payload.Object_Distance, 0, s.gateway.GetAgentCount(ctx)*num)
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
		// cl := new(checkList)
		visited := make(map[string]bool, len(res.Results))
		mu := sync.RWMutex{}
		return s.gateway.BroadCast(ectx, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.search/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			r, err := f(ctx, vc, copts...)
			if err != nil {
				log.Debug(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil
			}
			for _, dist := range r.GetResults() {
				if dist.GetDistance() > math.Float32frombits(atomic.LoadUint32(&maxDist)) {
					return nil
				}
				id := dist.GetId()
				mu.RLock()
				if !visited[id] {
					mu.RUnlock()
					mu.Lock()
					visited[id] = true
					mu.Unlock()
					dch <- dist
				} else {
					mu.RUnlock()
				}
			}
			return nil
		})
	}))
	for {
		select {
		case <-ectx.Done():
			err = eg.Wait()
			if err != nil {
				log.Error(err)
			}
			close(dch)
			if num != 0 && len(res.GetResults()) > num {
				res.Results = res.Results[:num]
			}
			return res, nil
		case dist := <-dch:
			if len(res.GetResults()) >= num &&
				dist.GetDistance() < math.Float32frombits(atomic.LoadUint32(&maxDist)) {
				atomic.StoreUint32(&maxDist, math.Float32bits(dist.GetDistance()))
			}
			switch len(res.GetResults()) {
			case 0:
				res.Results = append(res.Results, dist)
				continue
			case 1:
				if res.GetResults()[0].GetDistance() <= dist.GetDistance() {
					res.Results = append(res.Results, dist)
				} else {
					res.Results = append([]*payload.Object_Distance{dist}, res.Results[0])
				}
				continue
			}

			pos := len(res.GetResults())
			for idx := pos; idx >= 1; idx-- {
				if res.GetResults()[idx-1].GetDistance() <= dist.GetDistance() {
					pos = idx - 1
					break
				}
			}
			switch {
			case pos == len(res.GetResults()):
				res.Results = append([]*payload.Object_Distance{dist}, res.Results...)
			case pos == len(res.GetResults())-1:
				res.Results = append(res.GetResults(), dist)
			case pos >= 0:
				res.Results = append(res.GetResults()[:pos+1], res.GetResults()[pos:]...)
				res.Results[pos+1] = dist
			}
			if len(res.GetResults()) > num && num != 0 {
				res.Results = res.GetResults()[:num]
			}
		}
	}
}

func (s *server) StreamSearch(stream vald.Vald_StreamSearchServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Search(ctx, data.(*payload.Search_Request))
		})
}

func (s *server) StreamSearchByID(stream vald.Vald_StreamSearchByIDServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
		})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(vec.Vector) < 2 {
		err = errors.ErrInvalidDimensionSize(len(vec.Vector), 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Search API invalid vector argument", err, vec, info.Get())
	}
	if s.strict {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: vec.GetId(),
		})
		if err == nil && id != nil && len(id.GetId()) != 0 {
			err = errors.ErrMetaDataAlreadyExists(vec.GetId())
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(
				fmt.Sprintf("Insert API ID = %v already exists", vec.GetId()), err, info.Get())
		}
	}

	mu := new(sync.Mutex)
	targets := make([]string, 0, s.replica)
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Insert/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		_, err = vc.Insert(ctx, vec, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			if err == errors.ErrRPCCallFailed(target, context.Canceled) {
				return nil
			}
			return err
		}
		target = strings.SplitN(target, ":", 2)[0]
		mu.Lock()
		targets = append(targets, target)
		mu.Unlock()
		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "Insert API (do multiple) failed to Insert uuid = %s\t info = %#v", vec.GetId(), info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	log.Debugf("Insert API insert succeeded to %v", targets)
	return &payload.Object_Location{
		Uuid: vec.GetId(),
		Ips:  targets,
	}, nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if s.strict {
		for _, vec := range vecs.GetVectors() {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: vec.GetId(),
			})
			if err == nil && id != nil && len(id.GetId()) != 0 {
				err = errors.ErrMetaDataAlreadyExists(vec.GetId())
				log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, status.WrapWithAlreadyExists(
					fmt.Sprintf("MultiInsert API ID = %v already exists", vec.GetId()), err, info.Get())
			}
		}
	}

	mu := new(sync.Mutex)
	targets := make([]string, 0, s.replica)
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.MultiInsert/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		_, err = vc.MultiInsert(ctx, vecs, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return err
		}
		target = strings.SplitN(target, ":", 2)[0]
		mu.Lock()
		targets = append(targets, target)
		mu.Unlock()
		return nil
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed request %#v", vecs), err, info.Get())
	}
	return new(payload.Object_Locations), nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err := s.Remove(ctx, &payload.Object_ID{
		Id: vec.GetId(),
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
		}
		return nil, err
	}

	return s.Insert(ctx, vec)
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ids := make([]string, 0, len(vecs.GetVectors()))
	for _, vec := range vecs.GetVectors() {
		ids = append(ids, vec.GetId())
	}
	_, err = s.MultiRemove(ctx, &payload.Object_IDs{
		Ids: ids,
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Remove request %#v", ids), err, info.Get())
	}
	_, err = s.MultiInsert(ctx, vecs)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Insert request %#v", vecs), err, info.Get())
	}
	return new(payload.Object_Locations), nil
}

func (s *server) Upsert(ctx context.Context, vec *payload.Object_Vector) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	meta := vec.GetId()
	exists, errs := s.metadata.Exists(ctx, meta)
	if errs != nil {
		log.Error(errs)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(errs.Error()))
		}
	}

	if exists {
		_, err := s.Update(ctx, vec)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	} else {
		_, err := s.Insert(ctx, vec)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}

	return new(payload.Object_Location), errs
}

func (s *server) StreamUpsert(stream vald.Vald_StreamUpsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Upsert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpsert(ctx context.Context, vecs *payload.Object_Vectors) (*payload.Object_Locations, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))
	updateVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))

	var errs error
	for _, vec := range vecs.GetVectors() {
		exists, err := s.metadata.Exists(ctx, vec.GetId())
		if err != nil {
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			errs = errors.Wrap(errs, err.Error())
		}

		if exists {
			updateVecs = append(updateVecs, vec)
		} else {
			insertVecs = append(insertVecs, vec)
		}
	}

	eg, ectx := errgroup.New(ctx)

	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(updateVecs) > 0 {
			_, err = s.MultiUpdate(ectx, &payload.Object_Vectors{
				Vectors: updateVecs,
			})
		}
		return err
	}))

	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(insertVecs) > 0 {
			_, err = s.MultiInsert(ectx, &payload.Object_Vectors{
				Vectors: insertVecs,
			})
		}
		return err
	}))

	err := eg.Wait()
	if err != nil {
		errs = errors.Wrap(errs, err.Error())
		return nil, status.WrapWithInternal("MultiUpsert API failed", errs, info.Get())
	}

	return new(payload.Object_Locations), errs
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (locs *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if s.strict {
		sid, err := s.Exists(ctx, id)
		if err != nil || sid == nil || len(sid.GetId()) == 0 {
			err = errors.ErrObjectNotFound(err, id.GetId())
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, status.WrapWithNotFound(
				fmt.Sprintf("Remove API ID = %v not found", id.GetId()), err, info.Get())
		}
	}
	var mu sync.Mutex
	locs = &payload.Object_Location{
		Uuid: id.GetId(),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.Remove/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.Remove(ctx, id, copts...)
		if err != nil {
			log.Debug(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil
		}
		mu.Lock()
		locs.Ips = append(locs.Ips, loc.GetIps()...)
		locs.Name = loc.GetName()
		mu.Unlock()
		return nil
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed request uuid %s", id.GetId()), err, info.Get())
	}
	return locs, nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if s.strict {
		for _, id := range ids.GetIds() {
			sid, err := s.Exists(ctx, &payload.Object_ID{
				Id: id,
			})
			if err != nil || sid == nil || len(sid.GetId()) == 0 {
				err = errors.ErrObjectNotFound(err, id)
				log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil, status.WrapWithNotFound(
					fmt.Sprintf("MultiRemove API ID = %v not found", id), err, info.Get())
			}
		}
	}
	var mu sync.Mutex
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(ids.GetIds())),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.MultiRemove/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.MultiRemove(ctx, ids, copts...)
		if err != nil {
			log.Debug(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil
		}
		mu.Lock()
		locs.Locations = append(locs.Locations, loc.Locations...)
		mu.Unlock()
		return nil
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to request uuids %v ", ids.GetIds()), err, info.Get())
	}
	lmap := make(map[string]*payload.Object_Location, len(locs.Locations))
	for _, loc := range locs.Locations {
		_, ok := lmap[loc.Uuid]
		if !ok {
			lmap[loc.Uuid] = new(payload.Object_Location)
		}
		lmap[loc.Uuid].Ips = append(lmap[loc.Uuid].GetIps(), loc.GetIps()...)
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(lmap)),
	}
	for _, id := range ids.GetIds() {
		loc, ok := lmap[id]
		if !ok {
			loc = new(payload.Object_Location)
		}
		locs.Locations = append(locs.Locations, loc)
	}
	return locs, nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	var once sync.Once
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ValdClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, "vald/gateway-lb.GetObject/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		ovec, err := vc.GetObject(ctx, id, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil
		}
		if ovec != nil && ovec.GetId() != "" && ovec.GetVector() != nil {
			once.Do(func() {
				vec = ovec
				cancel()
			})
		}
		return nil
	})
	if err != nil || vec == nil || vec.GetId() != "" || vec.GetVector() != nil {
		err = errors.ErrObjectNotFound(err, id.GetId())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API uuid %s's object not found", vec.GetId()), err, info.Get())
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-lb.StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.GetObject(ctx, data.(*payload.Object_ID))
		})
}
