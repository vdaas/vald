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
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/db/kvs/bbolt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
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
	PreStop(ctx context.Context) error
}

type correct struct {
	cfg                   *config.Data
	discoverer            discoverer.Client
	agentAddrs            []string
	indexInfos            valdsync.Map[string, *payload.Info_Index_Count]
	uuidsCount            uint32
	uncommittedUUIDsCount uint32
	checkedID             bbolt.Bbolt
}

func New(cfg *config.Data, discoverer discoverer.Client) (Corrector, error) {
	d := filepath.Join(os.TempDir(), "bbolt")
	file.MkdirAll(d, os.ModePerm)
	dbfile := filepath.Join(d, "checkedid.db")
	bolt, err := bbolt.New(dbfile, "", os.FileMode(0o600))
	if err != nil {
		return nil, err
	}

	return &correct{
		cfg:        cfg,
		discoverer: discoverer,
		checkedID:  bolt,
	}, nil
}

func (c *correct) Start(ctx context.Context) (<-chan error, error) {
	dech, err := c.discoverer.Start(ctx)
	if err != nil {
		return nil, err
	}

	// addrs is sorted by the memory usage of each agent(descending order)
	// this is decending because it's supposed to be used for index manager to decide
	// which pod to make a create index rpc(higher memory, first to commit)
	c.agentAddrs = c.discoverer.GetAddrs(ctx)
	log.Debug("agent addrs found:", c.agentAddrs)

	if l := len(c.agentAddrs); l <= 1 {
		log.Warn("only %d agent found, there must be more than two agents for correction to happen", l)
		return nil, err
	}

	err = c.loadInfos(ctx)
	if err != nil {
		return nil, err
	}

	// For debugging
	c.indexInfos.Range(func(addr string, info *payload.Info_Index_Count) bool {
		log.Debugf("index info: addr(%s), stored(%d), uncommitted(%d)", addr, info.GetStored(), info.GetUncommitted())
		return true
	})

	log.Info("starting correction...")
	if c.cfg.Corrector.UseCache {
		log.Info("with bbolt disk cache...")
		if err := c.correctWithCache(ctx); err != nil {
			log.Errorf("there's some errors while correction: %v", err)
			return nil, err
		}
	} else {
		log.Info("without cache...")
		if err := c.correct(ctx); err != nil {
			log.Errorf("there's some errors while correction: %v", err)
			return nil, err
		}
	}
	log.Info("correction finished successfully")

	return dech, nil
}

func (c *correct) PreStop(_ context.Context) error {
	log.Info("removing persistent cache files...")
	if err := c.checkedID.Close(true); err != nil {
		return err
	}
	return nil
}

func (c *correct) correct(ctx context.Context) (err error) {
	if err := c.discoverer.GetClient().OrderedRange(ctx, c.agentAddrs,
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

			streamEnd := make(chan struct{})
			var once sync.Once
			var mu sync.Mutex
			// maybe just iterate through the number of indexes is ok?
			// that way, we don't have to use this `streamEnd` channel
			for {
				select {
				case <-ctx.Done():
					return finalize()
				case <-streamEnd:
					return finalize()
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
						if err := c.checkConsistency(
							ctx,
							&vectorReplica{
								addr: addr,
								vec:  res.GetVector(),
							},
							c.agentAddrs, // FIXME: no cache pattern always have to check all the agents
						); err != nil {
							// TODO: valdとstdでerrorの処理が違うので注意
							// （valdはerrが着信するまでにスタートしていた処理は行われる）
							// (stdはerrが着信すると他は全て止まる)
							log.Errorf("failed to check consistency: %v", err)
							return nil // continue other processes
						}

						return nil
					})
				}
			}
		},
	); err != nil {
		log.Errorf("failed to range over agents(%v): %v", c.agentAddrs, err)
		return err
	}

	return nil
}

func (c *correct) correctWithCache(ctx context.Context) (err error) {
	// leftAgentAddrs is the agents' addr that hasn't been corrected yet.
	// This is used to know which agents possibly have the same index as the target replica.
	// We can say this because, thanks to caching, there is no way that the target replica is
	// in the agent that has already been corrected.
	leftAgentAddrs := make([]string, len(c.agentAddrs))
	n := copy(leftAgentAddrs, c.agentAddrs)
	if n != len(c.agentAddrs) {
		return fmt.Errorf("failed to copy agentAddrs")
	}

	if err := c.discoverer.GetClient().OrderedRange(ctx, c.agentAddrs,
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			// current address is the leftAgentAddrs[0] because this is OrderedRange and
			// leftAgentAddrs is copied from c.agentAddrs
			leftAgentAddrs = leftAgentAddrs[1:]

			vc := vald.NewValdClient(conn)
			stream, err := vc.StreamListObject(ctx, &payload.Object_List_Request{})
			if err != nil {
				return err
			}

			seg, ctx := stdeg.WithContext(ctx)
			concurrency := c.cfg.Corrector.GetStreamListConcurrency()
			seg.SetLimit(concurrency)

			bolteg, ctx := stdeg.WithContext(ctx)
			bolteg.SetLimit(2048)

			finalize := func() error {
				err = seg.Wait()
				if err != nil {
					log.Errorf("err group returned error: %v", err)
					return err
				}

				err = bolteg.Wait()
				if err != nil {
					log.Errorf("bolt err group returned error: %v", err)
					return err
				}
				log.Info("bbolt all batch finished")

				log.Infof("correction finished for agent %s", addr)
				return nil
			}
			defer finalize()

			streamEnd := make(chan struct{})
			var once sync.Once
			var mu sync.Mutex

			log.Infof("starting correction for agent %s, concurrency: %d", addr, concurrency)

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
						id := res.GetVector().GetId()

						ok := false
						_, ok, err = c.checkedID.Get([]byte(id))
						if err != nil {
							log.Errorf("failed to perform Get from bbolt: %v", err)
						}

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
							leftAgentAddrs,
						); err != nil {
							// TODO: valdとstdでerrorの処理が違うので注意
							// （valdはerrが着信するまでにスタートしていた処理は行われる）
							// (stdはerrが着信すると他は全て止まる)
							log.Errorf("failed to check consistency: %v", err)
							return nil // continue other processes
						}

						// TODO: define error group
						c.checkedID.AsyncSet(bolteg, []byte(id), nil)

						return nil
					})
				}
			}
		},
	); err != nil {
		log.Errorf("failed to range over agents(%v): %v", c.agentAddrs, err)
		return err
	}

	return nil
}

type vectorReplica struct {
	addr string
	vec  *payload.Object_Vector
}

// Validate len(addrs) >= 2 before calling this function
func (c *correct) checkConsistency(ctx context.Context, targetReplica *vectorReplica, leftAgentAddrs []string) error {
	// availableAddrs is the agents' addr that doesn't have the target replica thus is available to insert the replica
	// to fix the index replica number if required.
	availableAddrs := make([]string, 0, len(c.agentAddrs)-1)
	for _, addr := range c.agentAddrs {
		if addr != targetReplica.addr {
			availableAddrs = append(availableAddrs, addr)
		}
	}

	foundReplicas := make([]*vectorReplica, 0, len(availableAddrs))
	var mu sync.Mutex
	if err := c.discoverer.GetClient().OrderedRangeConcurrent(ctx, leftAgentAddrs, len(leftAgentAddrs),
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			// To avoid GetObject to myself. To maintain backward compatibility for withoug cache operation
			if addr == targetReplica.addr {
				return nil
			}

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
			for i, a := range availableAddrs {
				if a == addr {
					availableAddrs = availableAddrs[:i+copy(availableAddrs[i:], availableAddrs[i+1:])]
					break
				}
			}
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
	if err := c.correctReplica(ctx, targetReplica, foundReplicas, availableAddrs); err != nil {
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
	availableAddrs []string,
) error {
	// diff < 0 means there is less replica than the correct number
	existReplica := len(foundReplicas) + 1
	diff := existReplica - c.cfg.Gateway.IndexReplica
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

		// inserting with the reverse order of availableAddrs since the last agent has the lowest memory usage
		for i := len(availableAddrs) - 1; i >= 0 && diff < 0; i-- {
			addr := availableAddrs[i]
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

	// delte from others if there's more to delete
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
