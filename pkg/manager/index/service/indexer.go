//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package service
package service

import (
	"context"
	"math"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
)

type Indexer interface {
	Start(ctx context.Context) (<-chan error, error)
	NumberOfUUIDs() uint32
	NumberOfUncommittedUUIDs() uint32
	IsIndexing() bool
}

type index struct {
	client                 discoverer.Client
	eg                     errgroup.Group
	creationPoolSize       uint32
	indexDuration          time.Duration
	indexDurationLimit     time.Duration
	saveIndexDurationLimit time.Duration
	saveIndexWaitDuration  time.Duration
	saveIndexTargetAddrCh  chan string
	schMap                 sync.Map
	concurrency            int
	indexInfos             indexInfos
	indexing               atomic.Value // bool
	minUncommitted         uint32
	uuidsCount             uint32
	uncommittedUUIDsCount  uint32
}

func New(opts ...Option) (idx Indexer, err error) {
	i := new(index)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(i); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	i.indexing.Store(false)
	if i.indexDuration+i.indexDurationLimit+i.saveIndexDurationLimit == 0 {
		return nil, errors.ErrInvalidConfig
	}
	return i, nil
}

func (idx *index) Start(ctx context.Context) (<-chan error, error) {
	dech, err := idx.client.Start(ctx)
	if err != nil {
		return nil, err
	}
	err = idx.loadInfos(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 100)
	sech := make(chan error, 10)
	idx.saveIndexTargetAddrCh = make(chan string, len(idx.client.GetAddrs(ctx))*2)
	idx.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if idx.indexDuration <= 0 {
			idx.indexDuration = math.MaxInt64
		}
		if idx.indexDurationLimit <= 0 {
			idx.indexDurationLimit = math.MaxInt64
		}
		if idx.saveIndexDurationLimit <= 0 {
			idx.saveIndexDurationLimit = math.MaxInt64
		}
		it := time.NewTicker(idx.indexDuration)
		itl := time.NewTicker(idx.indexDurationLimit)
		stl := time.NewTicker(idx.saveIndexDurationLimit)
		defer it.Stop()
		defer itl.Stop()
		defer stl.Stop()
		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				return err
			}
			return nil
		}
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case err = <-dech:
				ech <- err
			case err = <-sech:
				ech <- err
			case <-it.C:
				err = idx.execute(grpc.WithGRPCMethod(ctx, "core.v1.Agent/CreateIndex"), true, false)
				if err != nil {
					ech <- err
					log.Error("an error occurred during indexing", err)
					err = nil
				}
				it.Reset(idx.indexDuration)
			case <-itl.C:
				err = idx.execute(grpc.WithGRPCMethod(ctx, "core.v1.Agent/CreateIndex"), false, false)
				if err != nil {
					ech <- err
					log.Error("an error occurred during indexing", err)
					err = nil
				}
				itl.Reset(idx.indexDurationLimit)
			case <-stl.C:
				err = idx.execute(grpc.WithGRPCMethod(ctx, "core.v1.Agent/CreateAndSaveIndex"), false, true)
				if err != nil {
					ech <- err
					log.Error("an error occurred during indexing and saving", err)
					err = nil
				}
				stl.Reset(idx.saveIndexDurationLimit)
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			}
		}
	}))
	idx.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(sech)
		for {
			select {
			case <-ctx.Done():
				return
			case addr := <-idx.saveIndexTargetAddrCh:
				idx.schMap.Delete(addr)
				_, err := idx.client.GetClient().
					Do(grpc.WithGRPCMethod(ctx, "core.v1.Agent/SaveIndex"), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
						return agent.NewAgentClient(conn).SaveIndex(ctx, &payload.Empty{}, copts...)
					})
				if err != nil {
					log.Warnf("an error occurred while calling SaveIndex of %s: %s", addr, err)
					select {
					case <-ctx.Done():
						return nil
					case sech <- err:
					}
				}
			}

			idx.waitForNextSaving(ctx)
		}
	}))
	return ech, nil
}

func (idx *index) execute(ctx context.Context, enableLowIndexSkip, immediateSaving bool) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-index/service/Indexer.execute")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if idx.indexing.Load().(bool) {
		return nil
	}
	idx.indexing.Store(true)
	defer idx.indexing.Store(false)
	addrs := idx.client.GetAddrs(ctx)
	err = idx.client.GetClient().OrderedRangeConcurrent(ctx, addrs,
		idx.concurrency,
		func(ctx context.Context,
			addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
		) (err error) {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			info, ok := idx.indexInfos.Load(addr)
			if ok && (info.GetUncommitted() == 0 || (enableLowIndexSkip && info.GetUncommitted() < idx.minUncommitted)) {
				return nil
			}
			ac := agent.NewAgentClient(conn)
			req := &payload.Control_CreateIndexRequest{
				PoolSize: idx.creationPoolSize,
			}
			if !immediateSaving {
				_, err = ac.CreateIndex(ctx, req, copts...)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil && st.Code() == codes.FailedPrecondition {
						log.Debugf("CreateIndex of %s skipped, message: %s, err: %v", addr, st.Message(), errors.Wrap(st.Err(), err.Error()))
						return nil
					}
					log.Warnf("an error occurred while calling CreateIndex of %s: %s", addr, err)
					return err
				}
				_, ok := idx.schMap.Load(addr)
				if !ok {
					select {
					case <-ctx.Done():
					case idx.saveIndexTargetAddrCh <- addr:
						idx.schMap.Store(addr, struct{}{})
					}
				}
				return nil
			}
			_, err = ac.CreateAndSaveIndex(ctx, req, copts...)
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil && st.Code() == codes.FailedPrecondition {
					log.Debugf("CreateIndex of %s skipped, message: %s, err: %v", addr, st.Message(), errors.Wrap(st.Err(), err.Error()))
					return nil
				}
				log.Warnf("an error occurred while calling CreateAndSaveIndex of %s: %s", addr, err)
				return err
			}
			idx.waitForNextSaving(ctx)
			return nil
		})
	if err != nil {
		return err
	}
	return idx.loadInfos(ctx)
}

func (idx *index) waitForNextSaving(ctx context.Context) {
	if idx.saveIndexWaitDuration > 0 {
		timer := time.NewTimer(idx.saveIndexWaitDuration)
		select {
		case <-ctx.Done():
		case <-timer.C:
		}
		timer.Stop()
	}
}

func (idx *index) loadInfos(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "core.v1.Agent/IndexInfo"), "vald/manager-index/service/Indexer.loadInfos")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var u, ucu uint32
	var infoMap indexInfos
	err = idx.client.GetClient().RangeConcurrent(ctx, len(idx.client.GetAddrs(ctx)),
		func(ctx context.Context,
			addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
		) (err error) {
			select {
			case <-ctx.Done():
				return nil
			default:
				info, err := agent.NewAgentClient(conn).IndexInfo(ctx, new(payload.Empty), copts...)
				if err != nil {
					log.Warnf("an error occurred while calling IndexInfo of %s: %s", addr, err)
					return nil
				}
				infoMap.Store(addr, info)
				atomic.AddUint32(&u, info.GetStored())
				atomic.AddUint32(&ucu, info.GetUncommitted())
			}
			return nil
		})
	if err != nil {
		return err
	}
	atomic.StoreUint32(&idx.uuidsCount, atomic.LoadUint32(&u))
	atomic.StoreUint32(&idx.uncommittedUUIDsCount, atomic.LoadUint32(&ucu))
	idx.indexInfos.Range(func(addr string, _ *payload.Info_Index_Count) bool {
		info, ok := infoMap.Load(addr)
		if !ok {
			idx.indexInfos.Delete(addr)
		}
		idx.indexInfos.Store(addr, info)
		infoMap.Delete(addr)
		return true
	})
	infoMap.Range(func(addr string, info *payload.Info_Index_Count) bool {
		idx.indexInfos.Store(addr, info)
		return true
	})
	return nil
}

func (idx *index) IsIndexing() bool {
	return idx.indexing.Load().(bool)
}

func (idx *index) NumberOfUUIDs() uint32 {
	return atomic.LoadUint32(&idx.uuidsCount)
}

func (idx *index) NumberOfUncommittedUUIDs() uint32 {
	return atomic.LoadUint32(&idx.uncommittedUUIDsCount)
}
