// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"cmp"
	"context"
	"fmt"
	"os"
	"reflect"
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/db/kvs/pogreb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Corrector interface {
	Start(ctx context.Context) error
	StartClient(ctx context.Context) (<-chan error, error)
	PreStop(ctx context.Context) error

	// For metrics
	NumberOfCheckedIndex() uint64
	NumberOfCorrectedOldIndex() uint64
	NumberOfCorrectedReplication() uint64
}

type correct struct {
	eg          errgroup.Group
	discoverer  discoverer.Client
	gateway     vc.Client
	checkedList pogreb.DB

	checkedIndexCount         atomic.Uint64
	correctedOldIndexCount    atomic.Uint64
	correctedReplicationCount atomic.Uint64

	indexReplica                 int
	streamListConcurrency        int
	backgroundSyncInterval       time.Duration
	backgroundCompactionInterval time.Duration
}

func New(opts ...Option) (_ Corrector, err error) {
	c := new(correct)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}

	dir := file.Join(os.TempDir(), "checked")
	err = file.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Errorf("failed to create dir %s", dir)
		return nil, err
	}
	path := file.Join(dir, "checked_id.db")
	db, err := pogreb.New(pogreb.WithPath(path),
		pogreb.WithBackgroundCompactionInterval(c.backgroundCompactionInterval),
		pogreb.WithBackgroundSyncInterval(c.backgroundSyncInterval))
	if err != nil {
		log.Errorf("failed to open checked List kvs DB %s", path)
		return nil, err
	}
	c.checkedList = db
	return c, nil
}

func (c *correct) StartClient(ctx context.Context) (_ <-chan error, err error) {
	ech := make(chan error, 2)
	gch, err := c.gateway.Start(ctx)
	if err != nil {
		return nil, err
	}
	dch, err := c.discoverer.Start(ctx)
	if err != nil {
		return nil, err
	}
	c.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-dch:
			case err = <-gch:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

func (c *correct) Start(ctx context.Context) (err error) {
	detail, err := c.gateway.IndexDetail(ctx, new(payload.Empty))
	if err != nil {
		return err
	}
	counts := detail.GetCounts()
	agents := make([]string, 0, len(counts))
	for agent := range counts {
		agents = append(agents, agent)
	}
	slices.SortFunc(agents, func(left, right string) int {
		return cmp.Compare(counts[right].GetStored(), counts[left].GetStored())
	})

	for _, agent := range agents {
		count, ok := counts[agent]
		if ok && count != nil {
			log.Infof("index info: addr(%s), stored(%d), uncommitted(%d), indexing=%t, saving=%t", agent, count.GetStored(), count.GetUncommitted(), count.GetIndexing(), count.GetSaving())
		}
	}
	log.Infof("sorted agents: %v,\tdiscovered agents: %v", agents, c.discoverer.GetAddrs(ctx))

	errs := make([]error, 0, len(agents))

	emptyReq := new(payload.Object_List_Request)

	start := time.Now()

	emptyByte := []byte("")

	corrected := 1

	replicas := slices.Clone(agents)

	log.Infof("processing order of agents: %v", agents)
	if err := c.discoverer.GetClient().OrderedRange(ctx, agents, func(ctx context.Context,
		addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (err error) {
		defer func() {
			corrected++
			if err != nil {
				errs = append(errs, err)
			}
		}()
		if len(replicas) > 0 {
			replicas = replicas[1:]
		}
		var (
			stored      uint32
			uncommitted uint32
			indexing    bool
			saving      bool
			debugMsg    string
		)
		count, ok := counts[addr]
		if ok && count != nil {
			stored = count.GetStored()
			uncommitted = count.GetUncommitted()
			indexing = count.GetIndexing()
			saving = count.GetSaving()
			debugMsg = fmt.Sprintf(
				"agent %s (total index detail = stored: %d, uncommitted: %d, indexing=%t, saving=%t), stream concurrency: %d, processing %d/%d, replicas: size(%d) = addrs%v",
				addr,
				stored,
				uncommitted,
				indexing,
				saving,
				c.streamListConcurrency,
				corrected,
				len(agents),
				len(replicas),
				replicas,
			)
			if stored+uncommitted == 0 {
				// id no indices in agent skip process
				log.Warnf("skipping index correction process due to zero index detected for %s", debugMsg)
				return nil
			}
		}

		eg, egctx := errgroup.WithContext(ctx)
		eg.SetLimit(c.streamListConcurrency)
		ctx, cancel := context.WithCancelCause(egctx)
		stream, err := vc.NewValdClient(conn).StreamListObject(ctx, emptyReq, copts...)
		if err != nil || stream == nil {
			return err
		}
		log.Infof("starting correction for %s", debugMsg)
		// The number of items to be received in advance is not known in advance.
		// This is because there is a possibility of new items being inserted during processing.
		for {
			select {
			case <-ctx.Done():
				if !errors.Is(ctx.Err(), context.Canceled) {
					log.Errorf("context done unexpectedly: %v for %s", ctx.Err(), debugMsg)
				}
				if !errors.Is(context.Cause(ctx), io.EOF) {
					log.Errorf("context canceled due to %v for %s", ctx.Err(), debugMsg)
				}
				err = eg.Wait()
				if err != nil {
					log.Errorf("correction returned error status errgroup returned error: %v for %s", ctx.Err(), debugMsg)
				} else {
					log.Infof("correction finished for %s", debugMsg)
				}
				return nil
			default:
				res, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						cancel(io.EOF)
					} else {
						cancel(errors.ErrStreamListObjectStreamFinishedUnexpectedly(err))
					}
				} else if res != nil && res.GetVector() != nil && res.GetVector().GetId() != "" && res.GetVector().GetTimestamp() < start.UnixNano() {
					eg.Go(safety.RecoverFunc(func() (err error) {
						vec := res.GetVector()
						ts := vec.GetTimestamp()
						id := vec.GetId()

						_, ok, err := c.checkedList.Get(id)
						if err != nil {
							log.Errorf("failed to perform Get from check list but still try to finish processing without cache: %v", err)
						}
						if ok {
							// already checked index
							return nil
						}
						defer func() {
							c.checkedList.Set(id, emptyByte)
							c.checkedIndexCount.Add(1)
						}()

						// If the replicas slice is empty, it means the last agent's correction process, the data for this ID is unique in the current Vald Cluster and only exists for this Agent.
						// The meaning of unique data is that the timestamp of this data is up-to-date and there is no over-replica problem, just a lack of replicas.
						// Therefore, the process is only to correct the missing replicas.
						if len(replicas) <= 0 {
							diff := c.indexReplica - 1
							// correct index replica shortage
							if diff > 0 {
								return c.correctShortage(egctx, id, addr, debugMsg, vec, make(map[string]*payload.Object_Timestamp), diff)
							}
							return nil
						}

						// load index replica from other agents and store it to found map
						found, skipped, latest, latestAgent, err := c.loadReplicaInfo(egctx, addr, id, replicas, counts, ts, start)
						if err != nil {
							return err
						}
						if len(found) != 0 && ((len(replicas) > 0 && len(skipped) == 0) || (len(skipped) > 0 && len(skipped) < len(replicas))) {
							// current object timestamp is not latest get latest object from other agent index replica
							if ts < latest && latestAgent != addr {
								latestObject := c.getLatestObject(egctx, id, addr, latestAgent, latest)
								if latestObject != nil && latestObject.GetVector() != nil && latestObject.GetId() != "" && latestObject.GetTimestamp() >= latest {
									vec = latestObject
								}
							}
							c.correctTimestamp(ctx, id, vec, found)
						} else if len(skipped) > 0 {
							log.Debugf("timestamp correction for index id %s skipped, replica %s, skipped agents: %v", id, addr, skipped)
						}
						diff := c.indexReplica - (len(found) + 1)
						if diff > 0 { // correct index replica shortage
							return c.correctShortage(egctx, id, addr, debugMsg, vec, found, diff)
						} else if diff < 0 { // correct index replica oversupply
							return c.correctOversupply(egctx, id, addr, debugMsg, found, diff)
						}
						return nil
					}))
				}
			}
		}
	}); err != nil {
		// This only happens when ErrGRPCClientConnNotFound is returned.
		// In other cases, OrderedRange continues processing, so error is used to keep track of the error status of correction.
		errs = append(errs, err)
	}
	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c *correct) PreStop(_ context.Context) error {
	log.Info("removing persistent cache files...")
	return c.checkedList.Close(true)
}

func (c *correct) NumberOfCheckedIndex() uint64 {
	return c.checkedIndexCount.Load()
}

func (c *correct) NumberOfCorrectedOldIndex() uint64 {
	return c.correctedOldIndexCount.Load()
}

func (c *correct) NumberOfCorrectedReplication() uint64 {
	return c.correctedReplicationCount.Load()
}

func (c *correct) loadReplicaInfo(
	ctx context.Context,
	originAddr, id string,
	replicas []string,
	counts map[string]*payload.Info_Index_Count,
	ts int64,
	start time.Time,
) (
	found map[string]*payload.Object_Timestamp,
	skipped []string,
	latest int64,
	latestAgent string,
	err error,
) {
	var mu sync.Mutex
	latestAgent = originAddr
	skipped = make([]string, 0, len(replicas))
	found = make(map[string]*payload.Object_Timestamp, c.indexReplica-1)
	tss := time.Unix(0, start.UnixNano()).Format(time.RFC3339Nano)
	err = c.discoverer.GetClient().OrderedRangeConcurrent(ctx, replicas, len(replicas),
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			if originAddr == addr {
				return nil
			}
			count, ok := counts[addr] // counts is read-only we don't need to lock.
			if ok && count != nil && count.GetStored() == 0 && count.GetUncommitted() == 0 {
				mu.Lock()
				skipped = append(skipped, addr)
				mu.Unlock()
				return nil
			}

			ots, err := vc.NewValdClient(conn).GetTimestamp(ctx, &payload.Object_TimestampRequest{
				Id: &payload.Object_ID{
					Id: id,
				},
			})
			if err != nil {
				if st, ok := status.FromError(err); !ok || st == nil {
					log.Errorf("gRPC call GetTimestamp to agent: %s, id: %s returned not a gRPC status error: %v", addr, id, err)
					return err
				} else if st.Code() == codes.NotFound {
					// when replica of agent > index replica, this happens
					return nil
				} else if st.Code() == codes.Canceled {
					return nil
				} else {
					log.Errorf("failed to GetTimestamp with unexpected error. agent: %s, id: %s, code: %v, message: %s", addr, id, st.Code(), st.Message())
					return err
				}
			}

			if ots == nil {
				// not found
				return nil
			}

			// skip if the vector is inserted after correction start
			if ots.GetTimestamp() > start.UnixNano() {
				log.Debugf("timestamp of vector(id: %s, timestamp: %s) is newer than correction start time(%s). skipping...",
					ots.GetId(),
					time.Unix(0, ots.GetTimestamp()).Format(time.RFC3339Nano),
					tss,
				)
				return nil
			}
			mu.Lock()
			found[addr] = ots
			if latest < ots.GetTimestamp() {
				latest = ots.GetTimestamp()
				if latest > ts {
					latestAgent = addr
				}
			}
			mu.Unlock()
			return nil
		},
	)
	return found, skipped, latest, latestAgent, err
}

func (c *correct) getLatestObject(
	ctx context.Context, id, addr, latestAgent string, latest int64,
) (latestObject *payload.Object_Vector) {
	_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), latestAgent, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (any, error) {
		obj, err := vc.NewValdClient(conn).GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: id,
			},
		}, copts...)
		if err != nil {
			if st, ok := status.FromError(err); !ok || st == nil {
				log.Errorf("gRPC call returned not a gRPC status error: %v", err)
				return nil, err
			} else if st.Code() == codes.NotFound {
				return nil, nil
			} else if st.Code() == codes.Canceled {
				return nil, nil
			}
			return nil, err
		}
		if obj == nil {
			// not found
			return nil, nil
		}
		if obj.GetTimestamp() >= latest && obj.GetId() != "" && obj.GetVector() != nil {
			latestObject = obj
		}
		return obj, nil
	})
	if err != nil {
		log.Errorf("failed to load latest object id: %s, agent: %s, timestamp: %d, error: %v", id, addr, latest, err)
	}
	if latestObject != nil && latestObject.GetTimestamp() < latest {
		latestObject.Timestamp = latest
	}
	return latestObject
}

func (c *correct) correctTimestamp(
	ctx context.Context,
	id string,
	latestObject *payload.Object_Vector,
	found map[string]*payload.Object_Timestamp,
) {
	tss := time.Unix(0, latestObject.GetTimestamp()).Format(time.RFC3339Nano) // timestamp string
	for addr, ots := range found {                                            // correct timestamp inconsistency
		if latestObject.GetTimestamp() > ots.GetTimestamp() {
			log.Infof("timestamp inconsistency detected with vector(id: %s, timestamp: %s). updating with the latest vector(id: %s, timestamp: %s)",
				ots.GetId(),
				time.Unix(0, ots.GetTimestamp()).Format(time.RFC3339Nano),
				latestObject.GetId(),
				tss,
			)
			_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), addr, func(ctx context.Context,
				conn *grpc.ClientConn,
				copts ...grpc.CallOption,
			) (any, error) {
				client := vc.NewValdClient(conn)
				_, err := client.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{
					Id:        latestObject.GetId(),
					Timestamp: latestObject.GetTimestamp(),
				}, copts...)
				if err != nil {
					if st, ok := status.FromError(err); !ok || st == nil {
						log.Errorf("gRPC call returned not a gRPC status error: %v", err)
						return nil, err
					} else if st.Code() == codes.Canceled ||
						st.Code() == codes.AlreadyExists ||
						st.Code() == codes.InvalidArgument ||
						st.Code() == codes.NotFound {
						return nil, nil
					}
					return nil, err
				}
				log.Infof("vector successfully updated. address: %s, uuid: %s, timestamp: %s", addr, latestObject.GetId(), tss)
				c.correctedOldIndexCount.Add(1)
				return nil, nil
			})
			if err != nil {
				log.Errorf("failed to fix timestamp to %s for id %s agent %s error: %w", tss, id, addr, err)
			}
		}
	}
}

func (c *correct) correctOversupply(
	ctx context.Context,
	id, selfAddr, debugMsg string,
	found map[string]*payload.Object_Timestamp,
	diff int,
) (err error) {
	addrs := c.discoverer.GetAddrs(ctx)
	log.Infof("replica oversupply(configured: %d, stored: %d, diff: %d) of vector id: %s detected for %s. deleting from agents = %v", c.indexReplica, len(found)+1, diff, id, debugMsg, found)
	if len(addrs) == 0 {
		return errors.ErrNoAvailableAgentToRemove
	}
	req := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: id,
		},
	}
	for _, daddr := range addrs {
		if diff < 0 {
			_, ok := found[daddr]
			if ok || daddr == selfAddr {
				_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), daddr, func(ctx context.Context,
					conn *grpc.ClientConn,
					copts ...grpc.CallOption,
				) (any, error) {
					_, err := vc.NewValdClient(conn).Remove(ctx, req, copts...)
					if err != nil {
						if st, ok := status.FromError(err); !ok || st == nil {
							log.Errorf("gRPC call returned not a gRPC status error: %v", err)
							return nil, err
						} else if st.Code() == codes.NotFound {
							diff++
							c.correctedReplicationCount.Add(1)
							return nil, nil
						} else if st.Code() == codes.Canceled {
							return nil, nil
						}
						return nil, err
					}
					diff++
					c.correctedReplicationCount.Add(1)
					return nil, nil
				})
				if err != nil {
					log.Errorf("failed to delete object from agent(%s): %w", daddr, err)
				}
			}
		}
	}
	return nil
}

func (c *correct) correctShortage(
	ctx context.Context,
	id, selfAddr, debugMsg string,
	latestObject *payload.Object_Vector,
	found map[string]*payload.Object_Timestamp,
	diff int,
) (err error) {
	addrs := c.discoverer.GetAddrs(ctx)
	log.Infof("replica shortage(configured: %d, stored: %d, diff: %d) of vector id: %s detected for %s. inserting to other agents = %v", c.indexReplica, len(found)+1, diff, id, debugMsg, addrs)
	if len(addrs) == 0 {
		return errors.ErrNoAvailableAgentToInsert
	}
	req := &payload.Insert_Request{
		Vector: latestObject,
		// TODO: this should be deleted after Config.Timestamp deprecation
		Config: &payload.Insert_Config{
			Timestamp: latestObject.GetTimestamp(),
		},
	}
	for _, daddr := range addrs {
		if diff > 0 && daddr != selfAddr {
			_, ok := found[daddr]
			if !ok {
				_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), daddr, func(ctx context.Context,
					conn *grpc.ClientConn,
					copts ...grpc.CallOption,
				) (any, error) {
					client := vc.NewValdClient(conn)
					_, err := client.Insert(ctx, req, copts...)
					if err != nil {
						if st, ok := status.FromError(err); !ok || st == nil {
							log.Errorf("gRPC call returned not a gRPC status error: %v", err)
							return nil, err
						} else if st.Code() == codes.AlreadyExists {
							var obj *payload.Object_Vector
							obj, err = client.GetObject(ctx, &payload.Object_VectorRequest{
								Id: &payload.Object_ID{
									Id: id,
								},
							}, copts...)
							if err != nil {
								if st, ok = status.FromError(err); !ok || st == nil {
									log.Errorf("gRPC call returned not a gRPC status error: %v", err)
									return nil, err
								} else if st.Code() == codes.NotFound {
									return nil, nil
								} else if st.Code() == codes.Canceled {
									return nil, nil
								}
								return nil, err
							}
							if obj != nil {
								if obj.GetTimestamp() < latestObject.GetTimestamp() {
									_, err = client.Update(ctx, &payload.Update_Request{
										Vector: latestObject,
										// TODO: this should be deleted after Config.Timestamp deprecation
										Config: &payload.Update_Config{
											// TODO: Decrementing because it's gonna be incremented before being pushed
											// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
											// so we should consider refactoring vqueue.
											Timestamp: latestObject.GetTimestamp() - 1,
										},
									}, copts...)
									if err != nil {
										if st, ok = status.FromError(err); !ok || st == nil {
											log.Errorf("gRPC call returned not a gRPC status error: %v", err)
											return nil, err
										} else if st.Code() == codes.NotFound {
											return nil, nil
										} else if st.Code() == codes.Canceled {
											return nil, nil
										}
										return nil, err
									}
								}
								diff--
								c.correctedReplicationCount.Add(1)
							}
							return nil, nil
						} else if st.Code() == codes.Canceled {
							return nil, nil
						}
						return nil, err
					}
					diff--
					c.correctedReplicationCount.Add(1)
					return nil, nil
				})
				if err != nil {
					log.Errorf("failed to insert object to agent(%s): %w", daddr, err)
				}
			}
		}
	}
	return nil
}
