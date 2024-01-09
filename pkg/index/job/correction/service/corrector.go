// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type contextTimeKey string

const (
	insertMethod                          = "core.v1.Vald/Insert"
	updateMethod                          = "core.v1.Vald/Update"
	deleteMethod                          = "core.v1.Vald/Delete"
	correctionStartTimeKey contextTimeKey = "correctionStartTimeKey"
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
	discoverer                discoverer.Client
	agentAddrs                []string
	sortedByIndexCntAddrs     []string
	uuidsCount                uint32
	uncommittedUUIDsCount     uint32
	checkedID                 bbolt.Bbolt
	checkedIndexCount         atomic.Uint64
	correctedOldIndexCount    atomic.Uint64
	correctedReplicationCount atomic.Uint64

	indexReplica               int
	streamListConcurrency      int
	bboltAsyncWriteConcurrency int
}

const filemode = 0o600

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
	if err := c.bboltInit(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *correct) bboltInit() error {
	dpath := file.Join(os.TempDir(), "bbolt")
	err := file.MkdirAll(dpath, os.ModePerm)
	if err != nil {
		return err
	}

	dbfile := file.Join(dpath, "checkedid.db")
	c.checkedID, err = bbolt.New(dbfile, "", os.FileMode(filemode))
	if err != nil {
		return err
	}
	return nil
}

func (c *correct) StartClient(ctx context.Context) (<-chan error, error) {
	return c.discoverer.Start(ctx)
}

func (c *correct) Start(ctx context.Context) error {
	// set current time to context
	ctx = embedTime(ctx)

	// addrs is sorted by the memory usage of each agent(descending order)
	// this is decending because it's supposed to be used for index manager to decide
	// which pod to make a create index rpc(higher memory, first to commit)
	c.agentAddrs = c.discoverer.GetAddrs(ctx)
	if len(c.agentAddrs) <= 1 {
		log.Warnf("target agent (%v) found, but there must be more than two agents for correction to happen", c.agentAddrs)
		return errors.ErrAgentReplicaOne
	}
	log.Debugf("target agent addrs: %v", c.agentAddrs)

	if err := c.loadAgentIndexInfo(ctx); err != nil {
		return err
	}

	log.Info("starting correction with bbolt disk cache...")
	if err := c.correct(ctx); err != nil {
		return err
	}
	log.Info("correction finished successfully")

	return nil
}

func (c *correct) PreStop(_ context.Context) error {
	log.Info("removing persistent cache files...")
	return c.checkedID.Close(true)
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

// skipcq: GO-R1005
func (c *correct) correct(ctx context.Context) (err error) {
	// Vector with time after this should not be processed
	correctionStartTime, err := correctionStartTime(ctx)
	if err != nil {
		log.Errorf("cannot determine correction start time: %w", err)
		return err
	}

	curTargetAgent := 0
	jobErrs := make([]error, 0, c.streamListConcurrency)
	if err := c.discoverer.GetClient().OrderedRange(ctx, c.sortedByIndexCntAddrs,
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
			defer func() {
				if err != nil {
					// catch the err that happened in this scope using named return err
					jobErrs = append(jobErrs, err)
				}
				curTargetAgent++
			}()

			// context and errgroup for stream.Recv and correction
			sctx, scancel := context.WithCancel(ctx)
			defer scancel()
			seg, sctx := errgroup.WithContext(sctx)
			seg.SetLimit(c.streamListConcurrency)

			// errgroup for bbolt AsyncSet
			bolteg, ctx := errgroup.WithContext(ctx)
			bolteg.SetLimit(c.bboltAsyncWriteConcurrency)

			log.Infof("starting correction for agent %s, stream concurrency: %d, bbolt concurrency: %d", addr, c.streamListConcurrency, c.bboltAsyncWriteConcurrency)

			vc := vald.NewValdClient(conn)
			stream, err := vc.StreamListObject(ctx, &payload.Object_List_Request{})
			if err != nil {
				return err
			}

			var mu sync.Mutex
			// The number of items to be received in advance is not known in advance.
			// This is because there is a possibility of new items being inserted during processing.
			for {
				select {
				case <-sctx.Done():
					if !errors.Is(sctx.Err(), context.Canceled) {
						log.Errorf("context done unexpectedly: %v", sctx.Err())
					}

					// Finalize
					err = seg.Wait()
					if err != nil {
						log.Errorf("err group returned error: %v", err)
					}

					berr := bolteg.Wait()
					if berr != nil {
						log.Errorf("bbolt err group returned error: %v", err)
						err = errors.Join(err, berr)
					} else {
						log.Info("bbolt all batch finished")
					}

					log.Infof("correction finished for agent %s", addr)
					return err

				default:
					seg.Go(safety.RecoverFunc(func() error {
						mu.Lock()
						// As long as we don't stream.Recv() from the stream, we do not consume the memory of the message.
						// So by limiting the number of this errgroup.Go instances, we can limit the memory usage
						// https://github.com/grpc/grpc-go/blob/33f9fa2e6e5bcf4cf8fe45133e23779ae6e43f6c/rpc_util.go#L795
						res, err := stream.Recv()
						mu.Unlock()

						if err != nil {
							if errors.Is(err, io.EOF) {
								scancel()
								return nil
							}
							return errors.ErrStreamListObjectStreamFinishedUnexpectedly(err)
						}

						vec := res.GetVector()
						if vec == nil {
							st := res.GetStatus()
							log.Error(st.GetCode(), st.GetMessage(), st.GetDetails())
							return errors.ErrFailedToReceiveVectorFromStream
						}

						// skip if the vector is inserted after correction start
						if vec.GetTimestamp() > correctionStartTime.UnixNano() {
							log.Debugf("timestamp of vector(id: %s, timestamp: %v) is newer than correction start time(%v). skipping...",
								vec.GetId(),
								vec.GetTimestamp(),
								correctionStartTime.UnixNano(),
							)
							return nil
						}

						// check if the index is already checked
						id := vec.GetId()
						_, ok, err := c.checkedID.Get([]byte(id))
						if err != nil {
							log.Errorf("failed to perform Get from bbolt but still try to finish processing without cache: %v", err)
						}
						if ok {
							// already checked index
							return nil
						}

						if err := c.checkConsistency(
							ctx,
							&vectorReplica{
								addr: addr,
								vec:  vec,
							},
							curTargetAgent,
						); err != nil {
							return errors.ErrFailedToCheckConsistency(err)
						}

						//  now this id is checked so set it to the disk cache
						c.checkedID.AsyncSet(bolteg, []byte(id), nil)
						c.checkedIndexCount.Add(1)

						return nil
					}))
				}
			}
		},
	); err != nil {
		// This only happnes when ErrGRPCClientConnNotFound is returned.
		// In other cases, OrderedRange continues processing, so jobErrrs is used to keep track of the error status of correction.
		return err
	}

	jobErrs = errors.RemoveDuplicates(jobErrs)
	return errors.Join(jobErrs...)
}

type vectorReplica struct {
	addr string
	vec  *payload.Object_Vector
}

// Validate len(addrs) >= 2 before calling this function
func (c *correct) checkConsistency(ctx context.Context, targetReplica *vectorReplica, targetAgentIdx int) error {
	// leftAgentAddrs is the agents' addr that hasn't been corrected yet.
	leftAgentAddrs := c.sortedByIndexCntAddrs[targetAgentIdx+1:]

	// Vector with time after this should not be processed
	correctionStartTime, err := correctionStartTime(ctx)
	if err != nil {
		log.Errorf("cannot determine correction start time: %w", err)
		return err
	}

	foundReplicas := make([]*vectorReplica, 0, len(c.sortedByIndexCntAddrs))
	var mu sync.Mutex
	if err := c.discoverer.GetClient().OrderedRangeConcurrent(ctx, leftAgentAddrs, len(leftAgentAddrs),
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			vecMeta, err := agent.NewAgentClient(conn).GetTimestamp(ctx, &payload.Object_GetTimestampRequest{
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

			// skip if the vector is inserted after correction start
			if vecMeta.GetTimestamp() > correctionStartTime.UnixNano() {
				log.Debugf("timestamp of vector(id: %s, timestamp: %v) is newer than correction start time(%v). skipping...",
					vecMeta.GetId(),
					vecMeta.GetTimestamp(),
					correctionStartTime.UnixNano(),
				)
				return nil
			}

			mu.Lock()
			foundReplicas = append(foundReplicas, &vectorReplica{
				addr: addr,
				// the vector itself will be fetched when it's needed
				vec: &payload.Object_Vector{
					Id:        vecMeta.GetId(),
					Timestamp: vecMeta.GetTimestamp(),
				},
			})
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
	if err := c.correctReplica(ctx, targetReplica, foundReplicas); err != nil {
		return fmt.Errorf("failed to fix index replica: %w", err)
	}

	return nil
}

func (c *correct) correctTimestamp(ctx context.Context, targetReplica *vectorReplica, foundReplicas []*vectorReplica) error {
	if len(foundReplicas) == 0 {
		// no replica found. nothing to do about timestamp
		return nil
	}

	// skipcq: CRT-D0001
	allReplicas := append(foundReplicas, targetReplica)

	// sort by timestamp
	slices.SortFunc(allReplicas, func(i, j *vectorReplica) int {
		// largest timestamp means the latest
		return cmp.Compare(j.vec.GetTimestamp(), i.vec.GetTimestamp())
	})

	latest := allReplicas[0]
	latestTS := latest.vec.GetTimestamp()
	for _, replica := range allReplicas {
		if replica.vec.GetTimestamp() == latestTS {
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
		c.correctedOldIndexCount.Add(1)
		if err := c.updateObject(ctx, replica, latest); err != nil {
			return err
		}
	}

	return nil
}

// correctReplica corrects the number of replicas of the target vector.
// skipcq: GO-R1005
func (c *correct) correctReplica(
	ctx context.Context,
	targetReplica *vectorReplica,
	foundReplicas []*vectorReplica,
) error {
	// diff < 0 means there is less replica than the correct number
	existReplica := len(foundReplicas) + 1
	diff := existReplica - c.indexReplica
	if diff == 0 {
		// replica number is correct
		return nil
	}

	// availableAddrs = c.agentAddrs - foundReplicas - targetReplica.addr
	// here we use c.agentAddrs because we want to decide by memory usage order
	// not the number of indexes
	availableAddrs := make([]string, 0, len(c.agentAddrs))
	for _, addr := range c.agentAddrs {
		if addr == targetReplica.addr {
			continue
		}
		if slices.ContainsFunc(foundReplicas, func(replica *vectorReplica) bool {
			return replica.addr == addr
		}) {
			continue
		}
		availableAddrs = append(availableAddrs, addr)
	}

	// when there are less replicas than the correct number, add the extra replicas
	if diff < 0 {
		log.Infof("replica shortage of vector %s. inserting to other agents...", targetReplica.vec.GetId())
		c.correctedReplicationCount.Add(1)
		if len(availableAddrs) == 0 {
			return errors.ErrNoAvailableAgentToInsert
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
			return errors.ErrFailedToCorrectReplicaNum
		}

		return nil
	}

	// when there are more replicas than the correct number, delete the extra replicas
	log.Infof("replica oversupply of vector %s. deleting...",
		targetReplica.vec.GetId())
	c.correctedReplicationCount.Add(1)
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
		return errors.ErrFailedToCorrectReplicaNum
	}

	return nil
}

func (c *correct) updateObject(ctx context.Context, dest, src *vectorReplica) error {
	// check if the src vector has content not just timestamp
	if vec := src.vec.GetVector(); len(vec) == 0 {
		if err := c.fillVectorField(ctx, src); err != nil {
			return err
		}
	}

	res, err := c.discoverer.GetClient().
		Do(grpc.WithGRPCMethod(ctx, updateMethod), dest.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			// TODO: use UpdateTimestamp when it's implemented because here we just want to update only the timestamp but not the vector
			return vald.NewUpdateClient(conn).Update(ctx, &payload.Update_Request{
				Vector: src.vec,
				// TODO: this should be deleted after Config.Timestamp deprecation
				Config: &payload.Update_Config{
					// TODO: Decrementing because it's gonna be incremented befor being pushed
					// to vqueue in the agent. This is a not ideal workaround for the current vqueue implementation
					// so we should consider refactoring vqueue.
					Timestamp: src.vec.GetTimestamp() - 1,
				},
			}, copts...)
		})
	if err != nil {
		return err
	}

	if v, ok := res.(*payload.Object_Location); ok {
		log.Infof("vector successfully updated. address: %s, uuid: %v", dest.addr, v.GetUuid())
	}

	return nil
}

func (c *correct) fillVectorField(ctx context.Context, replica *vectorReplica) error {
	res, err := c.discoverer.GetClient().
		Do(grpc.WithGRPCMethod(ctx, "core.v1.Vald/GetObject"), replica.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: replica.vec.GetId(),
				},
			}, copts...)
		})
	if err != nil {
		return err
	}

	if v, ok := res.(*payload.Object_Vector); ok {
		vec := v.GetVector()
		if len(vec) == 0 {
			return err
		}
		replica.vec.Vector = v.GetVector()
	}

	return nil
}

func (c *correct) insertObject(ctx context.Context, addr string, vector *payload.Object_Vector) error {
	res, err := c.discoverer.GetClient().
		Do(grpc.WithGRPCMethod(ctx, insertMethod), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewInsertClient(conn).Insert(ctx, &payload.Insert_Request{
				Vector: vector,
				// TODO: this should be deleted after Config.Timestamp deprecation
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
		Do(grpc.WithGRPCMethod(ctx, deleteMethod), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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

// loadAgentIndexInfo loads the index info of each agent and sort them by the number of indexes
// then append the result to c.sortedByIndexCntAddrs.
// This sort is required because we want to process the agents with the least number of indexes first
// for performance to filter out the agent as early as possible from broadcast in checkConsistency function.
func (c *correct) loadAgentIndexInfo(ctx context.Context) (err error) {
	var u, ucu uint32
	var infoMap sync.Map[string, *payload.Info_Index_Count]
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

	type indexInfo struct {
		stored int
		addr   string
	}

	var infos []indexInfo
	infoMap.Range(func(addr string, info *payload.Info_Index_Count) bool {
		log.Infof("index info: addr(%s), stored(%d), uncommitted(%d)", addr, info.GetStored(), info.GetUncommitted())

		infos = append(infos, indexInfo{
			addr:   addr,
			stored: int(info.GetStored() + info.GetUncommitted()),
		})
		return true
	})

	slices.SortFunc(infos, func(i, j indexInfo) int {
		return cmp.Compare(i.stored, j.stored)
	})
	for _, info := range infos {
		c.sortedByIndexCntAddrs = append(c.sortedByIndexCntAddrs, info.addr)
	}
	log.Infof("processing order of agents: %v", c.sortedByIndexCntAddrs)
	return nil
}

func embedTime(ctx context.Context) context.Context {
	v := ctx.Value(correctionStartTimeKey)
	if _, ok := v.(time.Time); ok {
		return ctx
	}
	return context.WithValue(ctx, correctionStartTimeKey, time.Now())
}

func correctionStartTime(ctx context.Context) (time.Time, error) {
	v := ctx.Value(correctionStartTimeKey)
	if t, ok := v.(time.Time); ok {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("timeKey is not embedded in context")
}
