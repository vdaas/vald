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

// Package service
package service

import (
	"context"
	"fmt"
	"math"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
)

type Indexer interface {
	Start(ctx context.Context) (<-chan error, error)
	UUIDs(context.Context) []string
	UncommittedUUIDs() []string
	IsIndexing() bool
}

type index struct {
	uuids              atomic.Value // []string uuid
	uncommittedUUIDs   atomic.Value // []string uncommitted uuid
	indexing           atomic.Value // bool
	agentName          string
	namespace          string
	nodeName           string
	agentPort          int
	concurrency        int
	agentARecord       string
	indexInfos         infoMap
	minUncommitted     int
	indexDuration      time.Duration
	indexDurationLimit time.Duration
	agents             atomic.Value // []string ips
	dscAddr            string
	dscDur             time.Duration
	dscClient          grpc.Client
	acClient           grpc.Client
	agentOpts          []grpc.Option
	eg                 errgroup.Group
}

func New(opts ...Option) (idx Indexer, err error) {
	i := new(index)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(i); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	i.uuids.Store(make([]string, 10))
	i.uncommittedUUIDs.Store(make([]string, 10))
	i.indexing.Store(false)
	return i, nil
}

func (idx *index) Start(ctx context.Context) (<-chan error, error) {
	dech, err := idx.dscClient.StartConnectionMonitor(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 100)
	err = idx.discover(ctx, ech)
	if err != nil {
		close(ech)
		idx.dscClient.Close()
		return nil, err
	}
	idx.acClient = grpc.New(
		append(
			idx.agentOpts,
			grpc.WithAddrs(idx.agents.Load().([]string)...),
			grpc.WithErrGroup(idx.eg),
		)...,
	)

	aech, err := idx.acClient.StartConnectionMonitor(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	idx.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		fch := make(chan struct{}, 1)
		defer close(fch)
		it := time.NewTicker(idx.indexDuration)
		defer it.Stop()
		itl := time.NewTicker(idx.indexDurationLimit)
		defer itl.Stop()
		dt := time.NewTicker(idx.dscDur)
		defer dt.Stop()
		finalize := func() (err error) {
			var errs error
			err = idx.dscClient.Close()
			if err != nil {
				errs = errors.Wrap(errs, err.Error())
			}
			err = idx.acClient.Close()
			if err != nil {
				errs = errors.Wrap(errs, err.Error())
			}
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				errs = errors.Wrap(errs, err.Error())
			}
			return errs
		}
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case err = <-dech:
			case err = <-aech:
			case <-fch:
				err = idx.discover(ctx, ech)
				if err != nil {
					ech <- err
					err = nil
				}
			case <-dt.C:
				err = idx.discover(ctx, ech)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
					time.Sleep(idx.dscDur / 5)
					fch <- struct{}{}
				}
			case <-it.C:
				err = idx.execute(ctx, true)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
				}
			case <-itl.C:
				err = idx.execute(ctx, false)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
				}
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			} else {
				log.Debug(idx.acClient.GetAddrs())
			}
		}
	}))
	return ech, nil
}

func (idx *index) discover(ctx context.Context, ech chan<- error) (err error) {
	log.Info("starting discoverer discovery")
	addrs := make([]string, 0, 100)
	_, err = idx.dscClient.Do(ctx, idx.dscAddr,
		func(ctx context.Context,
			conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			nodes, err := discoverer.NewDiscovererClient(conn).
				Nodes(ctx, &payload.Discoverer_Request{
					Namespace: idx.namespace,
					Name:      idx.agentName,
					Node:      idx.nodeName,
				}, copts...)
			if err != nil {
				return nil, err
			}
			for i := 0; i < (math.MaxInt32); i++ {
				visited := false
				for i, node := range nodes.GetNodes() {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					default:
						if node != nil && node.GetPods() != nil {
							pods := node.GetPods().GetPods()
							if i < len(pods) {
								addrs = append(addrs, fmt.Sprintf("%s:%d", pods[i].GetIp(), idx.agentPort))
								if !visited {
									visited = true
								}
								break
							}
						}
					}
				}
				if !visited {
					return nil, nil
				}
			}
			// TODO sort addrs here
			return nil, nil
		})
	if err != nil {
		log.Warn("failed to discover agents from discoverer API, trying to discover from dns...")
		ips, err := net.DefaultResolver.LookupIPAddr(ctx, idx.agentARecord)
		if err != nil {
			return errors.ErrAgentAddrCouldNotDiscover(err, idx.agentARecord)
		}
		addrs = make([]string, 0, len(ips))
		for _, ip := range ips {
			addrs = append(addrs, fmt.Sprintf("%s:%d", ip.String(), idx.agentPort))
		}
	}
	if idx.acClient != nil {
		cur := make(map[string]struct{}, len(addrs))
		i := len(addrs)
		connected := make([]string, i)
		for _, addr := range addrs {
			err = idx.acClient.Connect(ctx, addr)
			if err != nil {
				ech <- err
				err = nil
			} else {
				cur[addr] = struct{}{}
				connected[i-1] = addr
				i--
			}
		}
		idx.agents.Store(connected[i:])
		err = idx.acClient.Range(ctx,
			func(ctx context.Context,
				addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
				_, ok := cur[addr]
				if !ok {
					idx.indexInfos.Delete(addr)
					return idx.acClient.Disconnect(addr)
				}
				delete(cur, addr)
				return nil
			})
		if err != nil {
			ech <- err
			err = nil
		}
	}
	log.Info("finished discoverer discovery")
	return idx.loadInfos(ctx)
}

func (idx *index) execute(ctx context.Context, enableLowIndexSkip bool) (err error) {
	if idx.indexing.Load().(bool) {
		return nil
	}
	idx.indexing.Store(true)
	defer idx.indexing.Store(false)
	err = idx.acClient.OrderedRangeConcurrent(ctx, idx.agents.Load().([]string),
		idx.concurrency,
		func(ctx context.Context,
			addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
			select {
			case <-ctx.Done():
				return nil
			default:
				if enableLowIndexSkip {
					info, ok := idx.indexInfos.Load(addr)
					if ok && len(info.GetUncommittedUuids()) < idx.minUncommitted {
						return nil
					}
				}
				_, err := agent.NewAgentClient(conn).CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: 10000,
				}, copts...)
				if err != nil {
					log.Debug(addr, err)
					return err
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return idx.loadInfos(ctx)
}

func (idx *index) loadInfos(ctx context.Context) (err error) {
	var uuids sync.Map
	var ucuuids sync.Map
	err = idx.acClient.RangeConcurrent(ctx, len(idx.agents.Load().([]string)),
		func(ctx context.Context,
			addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
			select {
			case <-ctx.Done():
				return nil
			default:
				info, err := agent.NewAgentClient(conn).IndexInfo(ctx, new(payload.Empty), copts...)
				if err != nil {
					log.Debug(addr, err)
					return err
				}
				idx.indexInfos.Store(addr, info)
				for _, uuid := range info.GetUuids() {
					uuids.Store(uuid, struct{}{})
				}
				for _, ucuuid := range info.GetUncommittedUuids() {
					ucuuids.Store(ucuuid, struct{}{})
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	us := make([]string, 0, len(idx.uuids.Load().([]string)))
	uuids.Range(func(uuid, _ interface{}) bool {
		us = append(us, uuid.(string))
		return true
	})
	ucus := make([]string, 0, len(idx.uncommittedUUIDs.Load().([]string)))
	ucuuids.Range(func(uuid, _ interface{}) bool {
		ucus = append(ucus, uuid.(string))
		return true
	})
	idx.uuids.Store(us)
	idx.uncommittedUUIDs.Store(ucus)
	return nil
}

func (idx *index) IsIndexing() bool {
	return idx.indexing.Load().(bool)
}

func (idx *index) UUIDs(context.Context) []string {
	return idx.uuids.Load().([]string)
}
func (idx *index) UncommittedUUIDs() []string {
	return idx.uncommittedUUIDs.Load().([]string)
}
