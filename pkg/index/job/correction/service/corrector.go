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
	"io"
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
	path := file.Join(dir, "checkedid.db")
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
	agents := make([]string, 0, detail.GetLiveAgents())
	for agent, count := range counts {
		log.Infof("index info: addr(%s), stored(%d), uncommitted(%d), indexing=%t, saving=%t", agent, count.GetStored(), count.GetUncommitted(), count.GetIndexing(), count.GetSaving())
		agents = append(agents, agent)
	}
	slices.SortFunc(agents, func(left, right string) int {
		return cmp.Compare(counts[left].GetStored(), counts[right].GetStored())
	})

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
		)
		count, ok := counts[addr]
		if ok && count != nil {
			stored = count.GetStored()
			uncommitted = count.GetUncommitted()
			if stored+uncommitted == 0 {
				// id no indices in agent skip process
				return nil
			}
			indexing = count.GetIndexing()
			saving = count.GetSaving()
		}
		debugMsg := fmt.Sprintf("agent %s (stored: %d, uncommitted: %d, indexing=%t, saving=%t), stream concurrency: %d, processing %d/%d, replicas: size(%d) = addrs%v", addr, stored, uncommitted, indexing, saving, c.streamListConcurrency, corrected, len(agents), len(replicas), replicas)

		eg, egctx := errgroup.WithContext(ctx)
		eg.SetLimit(c.streamListConcurrency)
		ctx, cancel := context.WithCancelCause(egctx)
		stream, err := vald.NewObjectClient(conn).StreamListObject(ctx, emptyReq, copts...)
		if err != nil {
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
				if context.Cause(ctx) != io.EOF {
					log.Errorf("context canceled due to %v for %s", ctx.Err(), debugMsg)
				}
				err = eg.Wait()
				if err != nil {
					log.Errorf("errgroup returned error: %v for %s", ctx.Err(), debugMsg)
					return err
				}
				log.Infof("correction finished for %s", debugMsg)
				return nil
			default:
				res, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						cancel(io.EOF)
					} else {
						cancel(errors.ErrStreamListObjectStreamFinishedUnexpectedly(err))
					}
				} else {
					eg.Go(safety.RecoverFunc(func() (err error) {
						vec := res.GetVector()
						if vec == nil || vec.GetId() == "" {
							st := res.GetStatus()
							if st != nil {
								log.Errorf("invalid vector id: %s detected and returned status code: %d, message: %s, details: %v, debug: %s", vec.GetId(), st.GetCode(), st.GetMessage(), st.GetDetails(), debugMsg)
							}
							return errors.ErrFailedToReceiveVectorFromStream
						}

						// skip if the vector is inserted after correction start
						if vec.GetTimestamp() > start.UnixNano() {
							log.Debugf("index correction process for ID: %s skipped due to newer timestamp detected. job started at %s but object timestamp is %s",
								vec.GetId(),
								start.Format(time.RFC3339Nano),
								time.Unix(0, vec.GetTimestamp()).Format(time.RFC3339Nano))
							return nil
						}

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
							addrs := c.discoverer.GetAddrs(egctx)
							// correct index replica shortage
							if diff > 0 {
								log.Infof("replica shortage(diff=%d) of vector id: %s detected from last %s. inserting to other agents = %v", diff, id, debugMsg, addrs)
								if len(addrs) == 0 {
									return errors.ErrNoAvailableAgentToInsert
								}
								req := &payload.Insert_Request{
									Vector: vec,
									// TODO: this should be deleted after Config.Timestamp deprecation
									Config: &payload.Insert_Config{
										Timestamp: vec.GetTimestamp(),
									},
								}
								for _, daddr := range addrs {
									if diff > 0 && daddr != addr {
										_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(egctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), daddr, func(ctx context.Context,
											conn *grpc.ClientConn,
											copts ...grpc.CallOption,
										) (any, error) {
											client := vald.NewValdClient(conn)
											_, err := client.Insert(ctx, req, copts...)
											if err != nil {
												if st, ok := status.FromError(err); !ok {
													log.Errorf("gRPC call returned not a gRPC status error: %v", err)
													return nil, err
												} else if st.Code() == codes.AlreadyExists {
													obj, err := client.GetObject(ctx, &payload.Object_VectorRequest{
														Id: &payload.Object_ID{
															Id: id,
														},
													}, copts...)
													if err != nil {
														if st, ok = status.FromError(err); !ok {
															log.Errorf("gRPC call returned not a gRPC status error: %v", err)
															return nil, err
														} else if st.Code() == codes.NotFound {
															return nil, nil
														} else if st.Code() == codes.Canceled {
															return nil, nil
														}
														return nil, err
													}
													if obj.GetTimestamp() < vec.GetTimestamp() {
														_, err := client.Update(ctx, &payload.Update_Request{
															Vector: vec,
															// TODO: this should be deleted after Config.Timestamp deprecation
															Config: &payload.Update_Config{
																// TODO: Decrementing because it's gonna be incremented befor being pushed
																// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
																// so we should consider refactoring vqueue.
																Timestamp: vec.GetTimestamp() - 1,
															},
														}, copts...)
														if err != nil {
															if st, ok = status.FromError(err); !ok {
																log.Errorf("gRPC call returned not a gRPC status error: %v", err)
																return nil, err
															} else if st.Code() == codes.NotFound {
																return nil, nil
															} else if st.Code() == codes.Canceled {
																return nil, nil
															}
															return nil, err
														}
														c.correctedOldIndexCount.Add(1)
													}
													diff--
													c.correctedReplicationCount.Add(1)
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
											log.Error(fmt.Errorf("failed to insert object to agent(%s): %w", daddr, err))
										}
									}
								}
							}
							return nil
						}

						var (
							latest      int64
							mu          sync.Mutex
							found       = make(map[string]*payload.Object_Timestamp, len(addr))
							latestAgent = addr
						)
						// load index replica from other agents and store it to found map
						if err := c.discoverer.GetClient().OrderedRangeConcurrent(egctx, replicas, len(replicas),
							func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
								ots, err := vald.NewObjectClient(conn).GetTimestamp(ctx, &payload.Object_TimestampRequest{
									Id: &payload.Object_ID{
										Id: id,
									},
								})
								if err != nil {
									if st, ok := status.FromError(err); !ok {
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

								// skip if the vector is inserted after correction start
								if ots.GetTimestamp() > start.UnixNano() {
									log.Debugf("timestamp of vector(id: %s, timestamp: %v) is newer than correction start time(%v). skipping...",
										ots.GetId(),
										ots.GetTimestamp(),
										start.UnixNano(),
									)
									return nil
								}
								mu.Lock()
								found[addr] = ots
								if latest < ots.GetTimestamp() {
									latest = ots.GetTimestamp()
									if latest > vec.GetTimestamp() {
										latestAgent = addr
									}
								}
								mu.Unlock()
								return nil
							},
						); err != nil {
							return err
						}
						latestObject := vec

						// current object timestamp is not latest get latest object from other agent index replica
						if vec.GetTimestamp() < latest && latestAgent != addr {
							_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(egctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), latestAgent, func(ctx context.Context,
								conn *grpc.ClientConn,
								copts ...grpc.CallOption,
							) (any, error) {
								obj, err := vald.NewObjectClient(conn).GetObject(ctx, &payload.Object_VectorRequest{
									Id: &payload.Object_ID{
										Id: id,
									},
								}, copts...)
								if err != nil {
									if st, ok := status.FromError(err); !ok {
										log.Errorf("gRPC call returned not a gRPC status error: %v", err)
										return nil, err
									} else if st.Code() == codes.NotFound {
										return nil, nil
									} else if st.Code() == codes.Canceled {
										return nil, nil
									}
									return nil, err
								}
								if obj.GetTimestamp() >= latest && obj.GetId() != "" && obj.GetVector() != nil {
									latestObject = obj
								}
								return obj, nil
							})
							if err != nil {
								log.Error(fmt.Errorf("failed to load latest object id: %s, agent: %s, timestamp: %d, error: %w", id, addr, latest, err))
							}
						}
						if latestObject.Timestamp < latest {
							latestObject.Timestamp = latest
						}
						tss := time.Unix(0, latestObject.GetTimestamp()).Format(time.RFC3339Nano) // timestamp string
						for addr, ots := range found {                                            // correct timestamp inconsistency
							if latestObject.GetTimestamp() > ots.GetTimestamp() {
								log.Infof("timestamp inconsistency detected with vector(id: %s, timestamp: %s). updating with the latest vector(id: %s, timestamp: %s)",
									ots.GetId(),
									time.Unix(0, ots.GetTimestamp()).Format(time.RFC3339Nano),
									latestObject.GetId(),
									tss,
								)
								_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(egctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), addr, func(ctx context.Context,
									conn *grpc.ClientConn,
									copts ...grpc.CallOption,
								) (any, error) {
									client := vald.NewValdClient(conn)
									// TODO: use UpdateTimestamp when it's implemented because here we just want to update only the timestamp but not the vector
									_, err := client.Update(ctx, &payload.Update_Request{
										Vector: latestObject,
										// TODO: this should be deleted after Config.Timestamp deprecation
										Config: &payload.Update_Config{
											// TODO: Decrementing because it's gonna be incremented befor being pushed
											// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
											// so we should consider refactoring vqueue.
											Timestamp: latestObject.GetTimestamp() - 1,
										},
									}, copts...)
									if err != nil {
										if st, ok := status.FromError(err); !ok {
											log.Errorf("gRPC call returned not a gRPC status error: %v", err)
											return nil, err
										} else if st.Code() == codes.NotFound {
											_, err = client.Insert(ctx, &payload.Insert_Request{
												Vector: latestObject,
												// TODO: this should be deleted after Config.Timestamp deprecation
												Config: &payload.Insert_Config{
													// TODO: Decrementing because it's gonna be incremented befor being pushed
													// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
													// so we should consider refactoring vqueue.
													Timestamp: latestObject.GetTimestamp(),
												},
											}, copts...)
											if err != nil {
												if st, ok = status.FromError(err); !ok {
													log.Errorf("gRPC call returned not a gRPC status error: %v", err)
													return nil, err
												} else if st.Code() == codes.AlreadyExists {
													obj, err := client.GetObject(ctx, &payload.Object_VectorRequest{
														Id: &payload.Object_ID{
															Id: id,
														},
													}, copts...)
													if err != nil {
														if st, ok = status.FromError(err); !ok {
															log.Errorf("gRPC call returned not a gRPC status error: %v", err)
															return nil, err
														} else if st.Code() == codes.NotFound {
															return nil, nil
														} else if st.Code() == codes.Canceled {
															return nil, nil
														}
														return nil, err
													}
													if obj.GetTimestamp() < latestObject.GetTimestamp() {
														_, err = client.Update(ctx, &payload.Update_Request{
															Vector: latestObject,
															// TODO: this should be deleted after Config.Timestamp deprecation
															Config: &payload.Update_Config{
																// TODO: Decrementing because it's gonna be incremented befor being pushed
																// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
																// so we should consider refactoring vqueue.
																Timestamp: latestObject.GetTimestamp() - 1,
															},
														}, copts...)
														if err != nil {
															if st, ok = status.FromError(err); !ok {
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
													return nil, nil
												} else if st.Code() == codes.Canceled {
													return nil, nil
												}
												return nil, err
											}
											c.correctedOldIndexCount.Add(1)
											return nil, nil
										} else if st.Code() == codes.Canceled {
											return nil, nil
										}
										return nil, err
									}
									log.Infof("vector successfully updated. address: %s, uuid: %s, timestamp: %s", addr, latestObject.GetId(), tss)
									c.correctedOldIndexCount.Add(1)
									return nil, nil
								})
								if err != nil {
									log.Error(fmt.Errorf("failed to fix timestamp to %s for id %s agent %s error: %w", tss, id, addr, err))
								}
							}
						}
						currentNumberOfIndexReplica := len(found) + 1
						diff := c.indexReplica - currentNumberOfIndexReplica
						addrs := c.discoverer.GetAddrs(egctx)
						if diff > 0 { // correct index replica shortage
							log.Infof("replica shortage(diff=%d) of vector id: %s detected for %s. inserting to other agents = %v", diff, id, debugMsg, addrs)
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
								if diff > 0 && daddr != addr {
									_, ok := found[daddr]
									if !ok {
										_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(egctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), daddr, func(ctx context.Context,
											conn *grpc.ClientConn,
											copts ...grpc.CallOption,
										) (any, error) {
											client := vald.NewValdClient(conn)
											_, err := client.Insert(ctx, req, copts...)
											if err != nil {
												if st, ok := status.FromError(err); !ok {
													log.Errorf("gRPC call returned not a gRPC status error: %v", err)
													return nil, err
												} else if st.Code() == codes.AlreadyExists {
													obj, err := client.GetObject(ctx, &payload.Object_VectorRequest{
														Id: &payload.Object_ID{
															Id: id,
														},
													}, copts...)
													if err != nil {
														if st, ok = status.FromError(err); !ok {
															log.Errorf("gRPC call returned not a gRPC status error: %v", err)
															return nil, err
														} else if st.Code() == codes.NotFound {
															return nil, nil
														} else if st.Code() == codes.Canceled {
															return nil, nil
														}
														return nil, err
													}
													if obj.GetTimestamp() < latestObject.GetTimestamp() {
														_, err = client.Update(ctx, &payload.Update_Request{
															Vector: latestObject,
															// TODO: this should be deleted after Config.Timestamp deprecation
															Config: &payload.Update_Config{
																// TODO: Decrementing because it's gonna be incremented befor being pushed
																// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
																// so we should consider refactoring vqueue.
																Timestamp: latestObject.GetTimestamp() - 1,
															},
														}, copts...)
														if err != nil {
															if st, ok = status.FromError(err); !ok {
																log.Errorf("gRPC call returned not a gRPC status error: %v", err)
																return nil, err
															} else if st.Code() == codes.NotFound {
																return nil, nil
															} else if st.Code() == codes.Canceled {
																return nil, nil
															}
															return nil, err
														}
														c.correctedOldIndexCount.Add(1)
													}
													diff--
													c.correctedReplicationCount.Add(1)
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
											log.Error(fmt.Errorf("failed to insert object to agent(%s): %w", daddr, err))
										}
									}
								}
							}
						} else if diff < 0 { // correct index replica oversupply
							log.Infof("replica oversupply of vector %s. deleting...", id)
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
									if ok || daddr == addr {
										_, err := c.discoverer.GetClient().Do(grpc.WithGRPCMethod(egctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), daddr, func(ctx context.Context,
											conn *grpc.ClientConn,
											copts ...grpc.CallOption,
										) (any, error) {
											_, err := vald.NewRemoveClient(conn).Remove(ctx, req, copts...)
											if err != nil {
												if st, ok := status.FromError(err); !ok {
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
											log.Error(fmt.Errorf("failed to delete object from agent(%s): %w", daddr, err))
										}
									}
								}
							}
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
