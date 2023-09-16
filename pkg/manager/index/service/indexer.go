//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"sync/atomic"
	"time"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Indexer interface {
	Start(ctx context.Context) (<-chan error, error)
	NumberOfUUIDs() uint32
	NumberOfUncommittedUUIDs() uint32
	IsIndexing() bool
	IsSaving() bool
	LoadIndexDetail() *payload.Info_Index_Detail
}

type index struct {
	client                 discoverer.Client
	eg                     errgroup.Group
	creationPoolSize       uint32
	indexDuration          time.Duration
	indexDurationLimit     time.Duration
	saveIndexDuration      time.Duration
	saveIndexDurationLimit time.Duration
	shouldSaveList         sync.Map[string, struct{}]
	createIndexConcurrency int
	saveIndexConcurrency   int
	indexInfos             sync.Map[string, *payload.Info_Index_Count]
	indexing               atomic.Bool
	saving                 atomic.Bool
	minUncommitted         uint32
	uuidsCount             uint32
	uncommittedUUIDsCount  uint32
}

var empty = struct{}{}

func New(opts ...Option) (idx Indexer, err error) {
	i := new(index)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(i); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	i.indexing.Store(false)
	i.saving.Store(false)
	if i.indexDuration+i.indexDurationLimit+i.saveIndexDurationLimit <= 0 {
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
	idx.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if idx.indexDuration <= 0 {
			idx.indexDuration = math.MaxInt64
		}
		if idx.indexDurationLimit <= 0 {
			idx.indexDurationLimit = math.MaxInt64
		}
		if idx.saveIndexDuration <= 0 {
			idx.saveIndexDuration = math.MaxInt64
		}
		if idx.saveIndexDurationLimit <= 0 {
			idx.saveIndexDurationLimit = math.MaxInt64
		}
		it := time.NewTicker(idx.indexDuration)
		itl := time.NewTicker(idx.indexDurationLimit)
		st := time.NewTicker(idx.saveIndexDuration)
		stl := time.NewTicker(idx.saveIndexDurationLimit)
		defer it.Stop()
		defer itl.Stop()
		defer st.Stop()
		defer stl.Stop()
		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			return nil
		}
		var mu sync.Mutex
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case err = <-dech:
			case <-it.C: // index duration ticker
				// execute CreateIndex. This execution ignores low index agent.
				err = idx.createIndex(grpc.WithGRPCMethod(ctx, "core.v1.Agent/CreateIndex"), true)
				if err != nil &&
					!errors.Is(err, context.Canceled) &&
					!errors.Is(err, context.DeadlineExceeded) {
					err = errors.Wrap(err, "an error occurred during create indexing")
				}
				it.Reset(idx.indexDuration)
			case <-itl.C: // index duration limit ticker
				// execute CreateIndex. This execution always executes CreateIndex regardless of the state of the uncommitted index.
				err = idx.createIndex(grpc.WithGRPCMethod(ctx, "core.v1.Agent/CreateIndex"), false)
				if err != nil &&
					!errors.Is(err, context.Canceled) &&
					!errors.Is(err, context.DeadlineExceeded) {
					err = errors.Wrap(err, "an error occurred during force create indexing")
				}
				itl.Reset(idx.indexDurationLimit)
			case <-st.C: // save index duration ticker
				//  execute SaveIndex in concurrent.
				idx.eg.Go(safety.RecoverFunc(func() (err error) {
					if !mu.TryLock() {
						return
					}
					defer mu.Unlock()
					defer st.Reset(idx.saveIndexDuration)
					err = idx.saveIndex(grpc.WithGRPCMethod(ctx, "core.v1.Agent/SaveIndex"), false)
					if err != nil &&
						!errors.Is(err, context.Canceled) &&
						!errors.Is(err, context.DeadlineExceeded) {
						err = errors.Wrap(err, "an error occurred during save indexing")
						log.Error(err)
						select {
						case <-ctx.Done():
							return nil
						case ech <- err:
						}
					}
					return nil
				}))
			case <-stl.C: // save index duration limit ticker
				//  execute SaveIndex in concurrent.
				idx.eg.Go(safety.RecoverFunc(func() (err error) {
					if !mu.TryLock() {
						return
					}
					defer mu.Unlock()
					defer stl.Reset(idx.saveIndexDurationLimit)
					err = idx.saveIndex(grpc.WithGRPCMethod(ctx, "core.v1.Agent/SaveIndex"), true)
					if err != nil &&
						!errors.Is(err, context.Canceled) &&
						!errors.Is(err, context.DeadlineExceeded) {
						err = errors.Wrap(err, "an error occurred during force save indexing")
						log.Error(err)
						select {
						case <-ctx.Done():
							return nil
						case ech <- err:
						}
					}
					return nil
				}))
			}
			if err != nil &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

func (idx *index) createIndex(ctx context.Context, enableLowIndexSkip bool) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-index/service/Indexer.execute")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if idx.indexing.Load() {
		return nil
	}
	idx.indexing.Store(true)
	defer idx.indexing.Store(false)
	return errors.Join(idx.client.GetClient().OrderedRangeConcurrent(ctx, idx.client.GetAddrs(ctx),
		idx.createIndexConcurrency,
		func(ctx context.Context,
			addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
		) (err error) {
			info, ok := idx.indexInfos.Load(addr)
			if ok && (info.GetUncommitted() == 0 || (enableLowIndexSkip && info.GetUncommitted() < idx.minUncommitted)) {
				return nil
			}
			_, err = agent.NewAgentClient(conn).CreateIndex(ctx, &payload.Control_CreateIndexRequest{
				PoolSize: idx.creationPoolSize,
			}, copts...)
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil && st.Code() == codes.FailedPrecondition {
					log.Debugf("CreateIndex of %s skipped, message: %s, err: %v", addr, st.Message(), errors.Join(st.Err(), err))
					return nil
				}
				log.Warnf("an error occurred while calling CreateIndex of %s: %s", addr, err)
				return err
			}
			_, ok = idx.shouldSaveList.LoadOrStore(addr, empty)
			if ok {
				log.Debugf("addr %s already queued for saveIndex", addr)
				return nil
			}
			return nil
		}), idx.loadInfos(ctx))
}

func (idx *index) saveIndex(ctx context.Context, force bool) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-index/service/Indexer.saveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if idx.saving.Load() {
		return nil
	}
	idx.saving.Store(true)
	defer idx.saving.Store(false)
	return idx.client.GetClient().OrderedRangeConcurrent(ctx, idx.client.GetAddrs(ctx),
		idx.saveIndexConcurrency,
		func(ctx context.Context,
			addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
		) (err error) {
			_, ok := idx.shouldSaveList.LoadAndDelete(addr)
			if !ok && !force {
				return nil
			}
			_, err = agent.NewAgentClient(conn).SaveIndex(ctx, new(payload.Empty), copts...)
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil && st.Code() == codes.FailedPrecondition {
					log.Debugf("CreateIndex of %s skipped, message: %s, err: %v", addr, st.Message(), errors.Join(st.Err(), err))
					return nil
				}
				log.Warnf("an error occurred while calling CreateIndex of %s: %s", addr, err)
				return err
			}
			return nil
		})
}

func (idx *index) loadInfos(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "core.v1.Agent/IndexInfo"), "vald/manager-index/service/Indexer.loadInfos")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var u, ucu uint32
	var infoMap sync.Map[string, *payload.Info_Index_Count]
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
	return idx.indexing.Load()
}

func (idx *index) IsSaving() bool {
	return idx.saving.Load()
}

func (idx *index) NumberOfUUIDs() uint32 {
	return atomic.LoadUint32(&idx.uuidsCount)
}

func (idx *index) NumberOfUncommittedUUIDs() uint32 {
	return atomic.LoadUint32(&idx.uncommittedUUIDsCount)
}

func (idx *index) LoadIndexDetail() (detail *payload.Info_Index_Detail) {
	detail = new(payload.Info_Index_Detail)
	idx.indexInfos.Range(func(addr string, info *payload.Info_Index_Count) bool {
		detail.Counts[addr] = info
		return true
	})
	return detail
}
