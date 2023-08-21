// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/slices"
	valdsync "github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/index/job/correction/config"
	stdeg "golang.org/x/sync/errgroup"
)

type Corrector interface {
	Start(ctx context.Context) (<-chan error, error)
}

type correct struct {
	cfg                   *config.Data
	discoverer            discoverer.Client
	indexInfos            valdsync.Map[string, *payload.Info_Index_Count]
	uuidsCount            uint32
	uncommittedUUIDsCount uint32
	checkedId             map[string]struct{} // TODO: use mmap if necessary
	rwmu                  sync.RWMutex
}

func New(cfg *config.Data, discoverer discoverer.Client) (Corrector, error) {
	return &correct{
		cfg:        cfg,
		discoverer: discoverer,
		checkedId:  make(map[string]struct{}),
	}, nil
}

func (c *correct) Start(ctx context.Context) (<-chan error, error) {
	dech, err := c.discoverer.Start(ctx)
	if err != nil {
		return nil, err
	}

	addrs := c.discoverer.GetAddrs(ctx)
	log.Debug("agent addrs found:", addrs)

	if l := len(addrs); l <= 1 {
		log.Warn("only %d agent found, there must be more than two agents for correction to happen", l)
		return nil, err
	}

	err = c.loadInfos(ctx)
	if err != nil {
		return nil, err
	}

	// DEBUG:
	c.indexInfos.Range(func(addr string, info *payload.Info_Index_Count) bool {
		log.Debugf("index info: addr(%s), stored(%d), uncommitted(%d)", addr, info.GetStored(), info.GetUncommitted())
		return true
	})

	// This blocks. Should we run with errorgroup?
	log.Info("starting correction...")
	if err := c.correct(ctx, addrs); err != nil {
		log.Errorf("there's some errors while correction: %v", err)
		return nil, err
	}
	log.Info("correction finished successfully")

	// ech := make(chan error, 100)
	// c.eg.Go(safety.RecoverFunc(func() (err error) {
	// 	defer close(ech)
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			err = ctx.Err()
	// 			if err != nil && err != context.Canceled {
	// 				return err
	// 			}
	// 			return nil
	// 		case err = <-dech:
	// 			ech <- err
	// 		}
	// 	}
	// }))
	return dech, nil
}

func (c *correct) correct(ctx context.Context, addrs []string) (err error) {
	if err := c.discoverer.GetClient().OrderedRange(ctx, addrs,
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			vc := vald.NewValdClient(conn)
			stream, err := vc.StreamListObject(ctx, &payload.Object_List_Request{})
			if err != nil {
				return err
			}

			seg, ctx := stdeg.WithContext(ctx)
			concurrency := c.cfg.Corrector.GetStreamListConcurrency()
			seg.SetLimit(concurrency)

			log.Infof("starting correction for agent %s, concurrency: %d", addr, concurrency)

			finalize := func() error {
				err = seg.Wait()
				if err != nil {
					log.Errorf("err group returned error: %v", err)
					return err
				}
				log.Infof("correction finished for agent %s", addr)
				return nil
			}
			defer finalize()

			streamEnd := make(chan struct{})
			var once sync.Once
			var mu sync.Mutex
			// これをさらにerrgroupで囲みたくなるが、さすがに頭がおかしくなりそう
			// 事前にRecvすべき件数はわかるのだからその回数だけfor文を回すようにする方がいいか
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-streamEnd:
					return nil
				default:
					// TODO: when vald internal errgroup is changed to block when eg limitation is reached,
					//       switch to vald version of errgroup.
					seg.Go(func() error {
						mu.Lock()
						// As long as we don't stream.Recv() from the stream, we do not consume the memory of the message.
						// So by limiting the number of this errgroup.Go instances, we can limit the memory usage
						// https://github.com/grpc/grpc-go/blob/33f9fa2e6e5bcf4cf8fe45133e23779ae6e43f6c/rpc_util.go#L795
						res, err := stream.Recv()
						mu.Unlock()

						if errors.Is(err, io.EOF) {
							log.Debugf("StreamListObject stream finished for agent %s", addr)
							once.Do(func() {
								close(streamEnd)
							})
							return nil
						}
						if err != nil {
							log.Errorf("StreamListObject stream finished unexpectedly: %v", err)
							return err
						}

						if res.GetVector() == nil {
							st := res.GetStatus()
							log.Error(st.GetCode(), st.GetMessage(), st.GetDetails())
							// continue
							return nil
						}

						log.Debugf("received object in StreamListObject: agent(%s), id(%s), timestamp(%v)", addr, res.GetVector().GetId(), res.GetVector().GetTimestamp())

						// check if the index is already checked
						c.rwmu.RLock()
						_, ok := c.checkedId[res.GetVector().GetId()]
						c.rwmu.RUnlock()
						if ok {
							// already checked index
							return nil
						}

						if err := c.checkConsistency(
							ctx,
							&vectorReplica{
								addr: addr,
								vec:  res.GetVector(),
							},
							addrs,
						); err != nil {
							// TODO: valdとstdでerrorの処理が違うので注意
							// （valdはerrが着信するまでにスタートしていた処理は行われる）
							// (stdはerrが着信すると他は全て止まる)
							log.Errorf("failed to check consistency: %v", err)
							return nil // continue other processes
						}

						c.rwmu.Lock()
						c.checkedId[res.GetVector().GetId()] = struct{}{}
						c.rwmu.Unlock()

						return nil
					})
				}
			}
		},
	); err != nil {
		log.Errorf("failed to range over agents(%v): %v", addrs, err)
		return err
	}

	return nil
}

type vectorReplica struct {
	addr string
	vec  *payload.Object_Vector
}

// Validate len(addrs) >= 2 before calling this function
func (c *correct) checkConsistency(ctx context.Context, targetReplica *vectorReplica, addrs []string) error {
	// copy the addrs slice but delete the curAgentAddr
	otherAddrs := make([]string, 0, len(addrs)-1)
	availableAddrs := make(map[string]struct{})
	for _, addr := range addrs {
		if addr != targetReplica.addr {
			otherAddrs = append(otherAddrs, addr)
			availableAddrs[addr] = struct{}{}
		}
	}

	foundReplicas := make([]*vectorReplica, 0, len(otherAddrs))
	var mu sync.Mutex
	if err := c.discoverer.GetClient().OrderedRangeConcurrent(ctx, otherAddrs, len(otherAddrs),
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			vc := vald.NewValdClient(conn)
			v, err := vc.GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: targetReplica.vec.GetId(),
				},
			})
			if err != nil {
				if st, ok := status.FromError(err); !ok {
					log.Errorf("gRPC call returned not a gRPC status error: %v", err)
					return err
				} else if st.Code() == codes.NotFound {
					// when replica of agent > index replica, this happens
					return nil
				} else {
					log.Errorf("failed to GetObject with unexpected error. code: %v, message: %s", st.Code(), st.Message())
					return err
				}
			}

			log.Debugf("object found: agent(%s), id(%v), timestamp(%v)", addr, v.GetId(), v.GetTimestamp())
			mu.Lock()
			foundReplicas = append(foundReplicas, &vectorReplica{
				addr: addr,
				vec:  v,
			})
			delete(availableAddrs, addr)
			mu.Unlock()

			return nil
		},
	); err != nil {
		return err
	}

	// check timestamps
	if err := c.correctTimestamp(ctx, targetReplica, foundReplicas); err != nil {
		return fmt.Errorf("failed to fix timestamp: %w", err)
	}

	// check replica number
	replica := len(foundReplicas) + 1
	if err := c.correctReplica(ctx, targetReplica, foundReplicas, replica, availableAddrs); err != nil {
		return fmt.Errorf("failed to fix index replica: %w", err)
	}

	return nil
}

func (c *correct) correctTimestamp(ctx context.Context, targetReplica *vectorReplica, foundReplicas []*vectorReplica) error {
	if len(foundReplicas) == 0 {
		// no replica found. nothing to do about timestamp
		return nil
	}

	allReplicas := append(foundReplicas, targetReplica)

	// sort by timestamp
	slices.SortFunc(allReplicas, func(i, j *vectorReplica) bool {
		// largest timestamp means the latest
		return i.vec.GetTimestamp() > j.vec.GetTimestamp()
	})

	latest := allReplicas[0]
	latestTs := latest.vec.GetTimestamp()
	for _, replica := range allReplicas {
		if replica.vec.GetTimestamp() == latestTs {
			// no inconsistency
			continue
		}

		// udate the vector with the new one
		log.Infof("timestamp inconsistency detected with vector(id: %s, timestamp: %v). updating with the latest vector(id: %s, timestamp: %v)",
			replica.vec.GetId(),
			replica.vec.GetTimestamp(),
			latest.vec.GetId(),
			latest.vec.GetTimestamp(),
		)
		if err := c.updateObject(ctx, replica.addr, latest.vec); err != nil {
			return err
		}
	}

	return nil
}

func (c *correct) correctReplica(
	ctx context.Context,
	targetReplica *vectorReplica,
	foundReplicas []*vectorReplica,
	replica int,
	availableAddrs map[string]struct{},
) error {
	// diff < 0 means there is less replica than the correct number
	diff := replica - c.cfg.Gateway.IndexReplica
	if diff == 0 {
		// replica number is correct
		return nil
	}

	// when there are less replicas than the correct number, add the extra replicas
	// TODO: refine this logic. pretty complicated
	if diff < 0 {
		log.Infof("replica shortage of vector %s. inserting to other agents...",
			targetReplica.vec.GetId())
		if len(availableAddrs) == 0 {
			// TODO: define errors in errors pkg
			return fmt.Errorf("no available agent to insert replica")
		}

		// availableAddrsからdiff個選んでinsert処理する
		// TODO:　どのagentにinsertするのが最適化のロジックを考える
		//       とりあえずはランダムに入れとく
		for addr := range availableAddrs {
			if diff == 0 {
				break
			}
			log.Infof("inserting replica to %s", addr)
			if err := c.insertObject(ctx, addr, targetReplica.vec); err != nil {
				log.Errorf("failed to insert object to agent(%s): %v", addr, err)
				continue
			}
			diff++
		}

		if diff < 0 {
			return fmt.Errorf("failed to insert the sufficient amount of index to meet the replica setting")
		}

		return nil
	}

	// when there are more replicas than the correct number, delete the extra replicas
	log.Infof("replica oversupply of vector %s. deleting...",
		targetReplica.vec.GetId())
	// delete from myself
	if err := c.deleteObject(ctx, targetReplica.addr, targetReplica.vec); err != nil {
		log.Errorf("failed to delete object from agent(%s): %v", targetReplica.addr, err)
	} else {
		diff--
	}

	// delte from others
	for _, replica := range foundReplicas {
		if diff == 0 {
			break
		}
		if err := c.deleteObject(ctx, replica.addr, replica.vec); err != nil {
			log.Errorf("failed to delete object from agent(%s): %v", replica.addr, err)
			continue
		}
		diff--
	}

	if diff > 0 {
		return fmt.Errorf("failed to delete the sufficient amount of index to meet the replica setting")
	}

	return nil
}

func (c *correct) updateObject(ctx context.Context, addr string, vector *payload.Object_Vector) error {
	res, err := c.discoverer.GetClient().
		Do(grpc.WithGRPCMethod(ctx, "core.v1.Vald/Update"), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewUpdateClient(conn).Update(ctx, &payload.Update_Request{
				Vector: vector,
				// FIXME: this should be deleted after Config.Timestamp deprecation
				Config: &payload.Update_Config{
					// TODO: Decrementing because it's gonna be incremented befor being pushed
					// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
					// so we should consider refactoring vqueue.
					Timestamp: vector.GetTimestamp() - 1,
				},
			}, copts...)
		})
	if err != nil {
		return err
	}

	if v, ok := res.(*payload.Object_Location); ok {
		log.Infof("vector successfully updated. address: %s, uuid: %v", addr, v.GetUuid())
	}

	return nil
}

func (c *correct) insertObject(ctx context.Context, addr string, vector *payload.Object_Vector) error {
	res, err := c.discoverer.GetClient().
		Do(grpc.WithGRPCMethod(ctx, "core.v1.Vald/Insert"), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewInsertClient(conn).Insert(ctx, &payload.Insert_Request{
				Vector: vector,
				// FIXME: this should be deleted after Config.Timestamp deprecation
				Config: &payload.Insert_Config{
					Timestamp: vector.GetTimestamp(),
				},
			}, copts...)
		})
	if err != nil {
		return err
	}

	if v, ok := res.(*payload.Object_Location); ok {
		log.Infof("vector successfully inserted. address: %s, uuid: %v", addr, v.GetUuid())
	}

	return nil
}

func (c *correct) deleteObject(ctx context.Context, addr string, vector *payload.Object_Vector) error {
	res, err := c.discoverer.GetClient().
		Do(grpc.WithGRPCMethod(ctx, "core.v1.Vald/Delete"), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewRemoveClient(conn).Remove(ctx, &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: vector.GetId(),
				},
			}, copts...)
		})
	if err != nil {
		return err
	}

	if v, ok := res.(*payload.Object_Location); ok {
		log.Infof("vector successfully deleted. address: %s, uuid: %v", addr, v.GetUuid())
	}

	return nil
}

func (c *correct) loadInfos(ctx context.Context) (err error) {
	// FIXME: o11yは最後に整える
	// ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "core.v1.Agent/IndexInfo"), "vald/manager-index/service/Indexer.loadInfos")
	// defer func() {
	// 	if span != nil {
	// 		span.End()
	// 	}
	// }()

	var u, ucu uint32
	var infoMap valdsync.Map[string, *payload.Info_Index_Count]
	err = c.discoverer.GetClient().RangeConcurrent(ctx, len(c.discoverer.GetAddrs(ctx)),
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
	atomic.StoreUint32(&c.uuidsCount, atomic.LoadUint32(&u))
	atomic.StoreUint32(&c.uncommittedUUIDsCount, atomic.LoadUint32(&ucu))
	c.indexInfos.Range(func(addr string, _ *payload.Info_Index_Count) bool {
		info, ok := infoMap.Load(addr)
		if !ok {
			c.indexInfos.Delete(addr)
		}
		c.indexInfos.Store(addr, info)
		infoMap.Delete(addr)
		return true
	})
	infoMap.Range(func(addr string, info *payload.Info_Index_Count) bool {
		c.indexInfos.Store(addr, info)
		return true
	})
	return nil
}
