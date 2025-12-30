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
package tikv

import (
	"bytes"
	"context"
	"encoding/hex"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/tikv"
	"github.com/vdaas/vald/internal/client/v1/client/meta"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	apiName = "vald/internal/client/meta/v1/client/meta/tikv"
	waitInterval = time.Second * 2
)

type Client interface {
	meta.MetadataClient
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	StartPD(context.Context) (<-chan error, error)
	Stop(context.Context) error
	StopPD(context.Context) error
}

type client struct {
	addrs []string
	c     grpc.Client

	pd *pdClient

	regionErrorRetryLimit int

	clusterId uint64

	// range cache
	rmu     sync.RWMutex
	ranges  []*rangeEntry
	regions map[uint64]*rangeEntry
	storeIdToAddr map[uint64]string
}

// rangeEntry maps [start,end) to addr.
type rangeEntry struct {
	start []byte
	end   []byte // nil or empty means +Inf
	addr  string
	ctx   *tikv.Context
}

func (c *client) lookupRange(key []byte) *rangeEntry {
	idx := sort.Search(len(c.ranges), func(i int) bool {
		return bytes.Compare(c.ranges[i].start, key) > 0
	}) - 1
	if idx >= 0 {
		re := c.ranges[idx]
		if len(re.end) == 0 || bytes.Compare(key, re.end) < 0 {
			return re
		}
	}
	return nil
}

type lookupResult struct {
	keys [][]byte
	re   *rangeEntry
}

func (c *client) lookupAddrs(ctx context.Context, keys [][]byte) (map[uint64]*lookupResult, error) {
	unknownKeys := make([][]byte, 0, len(keys))
	res := make(map[uint64]*lookupResult)
	func () {
		c.rmu.RLock()
		defer c.rmu.RUnlock()
		for _, key := range keys {
			re := c.lookupRange(key)
			if re == nil {
				unknownKeys = append(unknownKeys, key)
				continue
			}
			if res[re.ctx.RegionId] == nil {
				res[re.ctx.RegionId] = &lookupResult{
					keys: [][]byte{key},
					re:   re,
				}
				continue
			}
			res[re.ctx.RegionId].keys = append(res[re.ctx.RegionId].keys, key)
		}
	}()
	err := c.refresh(ctx, unknownKeys)
	if err != nil {
		return nil, err
	}
	c.rmu.RLock()
	defer c.rmu.RUnlock()
	for _, key := range unknownKeys {
		re := c.lookupRange(key)
		if re.addr == "" {
			return nil, errors.Errorf("address not found for key: %s", hex.EncodeToString(key))
		}
		if res[re.ctx.RegionId] == nil {
			res[re.ctx.RegionId] = &lookupResult{
				keys: [][]byte{key},
				re:   re,
			}
			continue
		}
		res[re.ctx.RegionId].keys = append(res[re.ctx.RegionId].keys, key)
	}
	return res, nil
}

func (c *client) refresh(ctx context.Context, keys [][]byte) error {
	if c.clusterId == 0 {
		res, err := c.pd.GetClusterInfo(ctx, &tikv.GetClusterInfoRequest{
			Header: &tikv.ResponseHeader{},
		})
		if err != nil {
			return errors.Errorf("PD GetClusterInfo failed, err: %w", err)
		}
		if res.Header.Error != nil {
			return errors.Errorf("PD GetClusterInfo failed, message: %s", res.Header.Error.Message)
		}
		c.clusterId = res.Header.ClusterId
	}
	c.storeIdToAddr = make(map[uint64]string)
	var regions []*tikv.Region
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		res, err := c.pd.GetAllStores(ctx, &tikv.GetAllStoresRequest{
			Header:                 &tikv.RequestHeader{ClusterId: c.clusterId},
			ExcludeTombstoneStores: true,
		})
		if err != nil {
			return errors.Errorf("PD GetAllStores failed, err: %w", err)
		}
		for _, store := range res.Stores {
			c.storeIdToAddr[store.Id] = store.Address
		}
		return nil
	})
	g.Go(func() error {
		req := make([]*tikv.KeyRange, len(keys))
		sort.Slice(keys, func(i, j int) bool {
    		return bytes.Compare(keys[i], keys[j]) < 0
		})
		for i, key := range keys {
			req[i] = &tikv.KeyRange{
				StartKey: key,
				EndKey:   nil,
			}
		}
		res, err := c.pd.BatchScanRegions(ctx, &tikv.BatchScanRegionsRequest{
			Header: &tikv.RequestHeader{ClusterId: c.clusterId},
			Ranges: req,
			ContainAllKeyRange: true,
		})
		if err != nil {
			return errors.Errorf("PD BatchScanRegions failed, err: %w", err)
		}
		regions = res.Regions
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	c.rmu.Lock()
	defer c.rmu.Unlock()
	for _, r := range regions {
		if r == nil || r.Region == nil || r.Leader == nil {
			continue
		}
		if _, ok := c.storeIdToAddr[r.Leader.StoreId]; !ok {
			return errors.Errorf("store id %d not found", r.Leader.StoreId)
		}
		start := slices.Clone(r.Region.StartKey)
		end := slices.Clone(r.Region.EndKey)
		re := &rangeEntry{
			start: start,
			end:   end,
			addr:  c.storeIdToAddr[r.Leader.StoreId],
			ctx: &tikv.Context{
				RegionId: r.Region.Id,
				RegionEpoch: &tikv.RegionEpoch{
					Version: r.Region.RegionEpoch.Version,
					ConfVer: r.Region.RegionEpoch.ConfVer,
				},
				Peer: &tikv.Peer{
					Id:      r.Leader.Id,
					StoreId: r.Leader.StoreId,
				},
			},
		}
		c.regions[r.Region.Id] = re
		c.ranges = mergeByNewerVersion(c.ranges, re)
	}
	return nil
}

func epochCmp(a, b *tikv.RegionEpoch) int {
    // return 1 if a newer than b, 0 if equal, -1 if older
    if a.GetConfVer() != b.GetConfVer() {
        if a.GetConfVer() > b.GetConfVer() { return 1 }
        return -1
    }
    if a.GetVersion() != b.GetVersion() {
        if a.GetVersion() > b.GetVersion() { return 1 }
        return -1
    }
    return 0
}

func mergeByNewerVersion(old []*rangeEntry, newE *rangeEntry) []*rangeEntry {
		// log.Errorf("merging %#v into existing ranges %#v", newE, old)
    out := make([]*rangeEntry, 0, len(old)+1)
    inserted := false

    for _, e := range old {
        if !overlap(e.start, e.end, newE.start, newE.end) {
            out = append(out, e)
            continue
        }

        cmp := epochCmp(newE.ctx.RegionEpoch, e.ctx.RegionEpoch)
        switch {
        case cmp > 0:
            // new is newer: drop old overlapping entry
            continue
        case cmp == 0:
            // same epoch: replace entire entry (start/end/ctx/addr)
            out = append(out, newE)
            inserted = true
        default:
            // new is older: keep old set
            return old
        }
    }

    if !inserted {
        out = append(out, newE)
    }

    sort.Slice(out, func(i, j int) bool {
        return bytes.Compare(out[i].start, out[j].start) < 0
    })
		// log.Errorf("out: %#v", out)
    return out
}

func overlap(aStart, aEnd, bStart, bEnd []byte) bool {
	// Does [aStart, aEnd) and [bStart, bEnd) overlap?
	if len(aEnd) != 0 && bytes.Compare(aEnd, bStart) <= 0 {
		return false
	}
	if len(bEnd) != 0 && bytes.Compare(bEnd, aStart) <= 0 {
		return false
	}
	return true
}

var errNotFound = errors.New("tikv: key not found")

func New(opts ...Option) (Client, error) {
	c := new(client)
	c.pd = new(pdClient)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.regions = make(map[uint64]*rangeEntry)

	if c.c == nil {
		c.c = grpc.New("TiKV Client")
	}
	if c.pd.c == nil {
		if len(c.pd.addrs) == 0 {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		c.pd.c = grpc.New("PD Client", grpc.WithAddrs(c.pd.addrs...))
	}
	return c, nil
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	return c.c.StartConnectionMonitor(ctx)
}

func (c *client) StartPD(ctx context.Context) (<-chan error, error) {
	return c.pd.c.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	return c.c.Close(ctx)
}

func (c *client) StopPD(ctx context.Context) error {
	return c.pd.c.Close(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

func (c *client) GRPCClientPD() grpc.Client {
	return c.pd.c
}

func (c *client) handleRegionError(ctx context.Context, regionErr *tikv.Error, keys [][]byte, refresh bool) error {
	log.Errorf("region error: %+v", regionErr)
	c.rmu.Lock()
	if err := regionErr.ServerIsBusy; err != nil {
		c.rmu.Unlock()
		time.Sleep(waitInterval)
		return nil
	}
	if err := regionErr.NotLeader; regionErr.NotLeader != nil {
		if err.Leader == nil {
			// sleep while electing new leader
			c.rmu.Unlock()
			time.Sleep(waitInterval)
			return nil
		}
		if _, ok := c.regions[err.RegionId]; ok {
			c.regions[err.RegionId].addr = c.storeIdToAddr[err.Leader.StoreId]
			c.regions[err.RegionId].ctx.Peer.Id = err.Leader.Id
			c.regions[err.RegionId].ctx.Peer.StoreId = err.Leader.StoreId
			c.rmu.Unlock()
			return nil
		}
	}
	if err := regionErr.EpochNotMatch; err != nil {
		for _, r := range err.CurrentRegions {
			re := &rangeEntry{
				start: slices.Clone(r.StartKey),
				end:   slices.Clone(r.EndKey),
				addr:  c.storeIdToAddr[r.Peers[0].StoreId],
				ctx: &tikv.Context{
					RegionId: r.Id,
					RegionEpoch: &tikv.RegionEpoch{
						Version: r.RegionEpoch.Version,
						ConfVer: r.RegionEpoch.ConfVer,
					},
					Peer: &tikv.Peer{
						Id:      r.Peers[0].Id,
						StoreId: r.Peers[0].StoreId,
					},
				},
			}
			c.ranges = mergeByNewerVersion(c.ranges, re)
		}
		c.rmu.Unlock()
		return nil
	}
	c.rmu.Unlock()
	if !refresh {
		return nil
	}
	if err := c.refresh(ctx, keys); err != nil {
		return err
	}
	return nil
}

func (c *client) Get(ctx context.Context, key []byte) ([]byte, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Get"), apiName+"/Get")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	for range c.regionErrorRetryLimit {
		lookups, err := c.lookupAddrs(ctx, [][]byte{key})
		if err != nil {
			return nil, err
		}
		for _, lookup := range lookups {
			c.c.Connect(ctx, lookup.re.addr)
			res, err := grpc.Do(ctx, lookup.re.addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawGetResponse, error) {
				return tikv.NewTikvClient(conn).RawGet(ctx, &tikv.RawGetRequest{
					Context: lookup.re.ctx,
					Key:     key,
					Cf:      "default",
				}, copts...)
			})
			if err != nil {
				return nil, errors.Errorf("Get failed, err: %w", err)
			}
			if res.Error != "" {
				return nil, errors.Errorf("Get failed, message: %s", res.Error)
			}
			if res.NotFound {
				return nil, errNotFound
			}
			if res.RegionError != nil {
				if err = c.handleRegionError(ctx, res.RegionError, [][]byte{key}, true); err != nil {
					return nil, err
				}
				// After refresh, retry with updated region info.
				continue
			}
			return res.Value, nil
		}
	}
	return nil, errors.Errorf("exceeded region error retry limit for key: %s", hex.EncodeToString(key))
}

func (c *client) BatchGet(ctx context.Context, keys [][]byte) ([][]byte, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/BatchGet"), apiName+"/BatchGet")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	// remaining keys to fetch
	remaining := make([][]byte, len(keys))
	copy(remaining, keys)

	// map for final results
	resultKV := make(map[string][]byte, len(keys))

	for range c.regionErrorRetryLimit {

		lookups, err := c.lookupAddrs(ctx, keys)
		if err != nil {
			return nil, err
		}

		// failed keys in this iteration
		var failedKeys [][]byte
		var fkMu sync.Mutex

		// mutex for result map
		var kvMu sync.Mutex

		g, gctx := errgroup.WithContext(ctx)

		for _, lookup := range lookups {
			addr := lookup.re.addr
			keys := lookup.keys
			g.Go(func() error {
				c.c.Connect(gctx, addr)
				res, err := grpc.Do(gctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchGetResponse, error) {
					return tikv.NewTikvClient(conn).RawBatchGet(ctx, &tikv.RawBatchGetRequest{
						Context: lookup.re.ctx,
						Keys:    keys,
						Cf:      "default",
					}, copts...)
				})
				if err != nil {
					return errors.Errorf("BatchGet failed: %w", err)
				}
				if res.RegionError != nil {
					if err = c.handleRegionError(gctx, res.RegionError, keys, false); err != nil {
						return err
					}
					fkMu.Lock()
					failedKeys = append(failedKeys, keys...)
					fkMu.Unlock()
					return nil
				}
				for _, pair := range res.Pairs {
					if pair.Error != nil {
						return errors.Errorf("KeyError happened %+v", pair.Error)
					}
					kvMu.Lock()
					resultKV[hex.EncodeToString(pair.Key)] = pair.Value
					kvMu.Unlock()
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, err
		}

		if len(failedKeys) == 0 {
			break // all done
		}

		if err := c.refresh(ctx, failedKeys); err != nil {
			return nil, err
		}

		remaining = failedKeys
	}

	if len(resultKV) != len(keys) {
		return nil, errors.New("exceeded region error retry limit for BatchGet")
	}

	// order results in input order
	out := make([][]byte, len(keys))
	for i, key := range keys {
		out[i] = resultKV[hex.EncodeToString(key)]
	}
	return out, nil
}

func (c *client) Put(ctx context.Context, key, val []byte) error {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Put"), apiName+"/Put")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	for range c.regionErrorRetryLimit {
		lookups, err := c.lookupAddrs(ctx, [][]byte{key})
		if err != nil {
			return err
		}

		for _, lookup := range lookups {
			c.c.Connect(ctx, lookup.re.addr)
			res, err := grpc.Do(ctx, lookup.re.addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawPutResponse, error) {
				return tikv.NewTikvClient(conn).RawPut(ctx, &tikv.RawPutRequest{
					Context: lookup.re.ctx,
					Key:     key,
					Value:   val,
					Cf:      "default",
				}, copts...)
			})
			if err != nil {
				return errors.Errorf("Put failed, err: %w", err)
			}
			if res.Error != "" {
				return errors.Errorf("Put failed, message: %s", res.Error)
			}
			if res.RegionError != nil {
				if err = c.handleRegionError(ctx, res.RegionError, [][]byte{key}, true); err != nil {
					return err
				}
				// retry with refreshed region info
				goto RETRY
			}
			return nil
		}
	RETRY:
	}
	return errors.Errorf("exceeded region error retry limit for key: %s", hex.EncodeToString(key))
}

func (c *client) BatchPut(ctx context.Context, keys, vals [][]byte) error {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/BatchPut"), apiName+"/BatchPut")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	// map key(hex) -> value for quick access
	valMap := make(map[string][]byte, len(keys))
	for i, k := range keys {
		valMap[hex.EncodeToString(k)] = vals[i]
	}

	remaining := make([][]byte, len(keys))
	copy(remaining, keys)

	for range c.regionErrorRetryLimit {
		lookups, err := c.lookupAddrs(ctx, remaining)
		if err != nil {
			return err
		}

		var failedKeys [][]byte
		var fkMu sync.Mutex

		g, gctx := errgroup.WithContext(ctx)

		for _, lookup := range lookups {
			lookup := lookup
			g.Go(func() error {
				// build pairs for this addr
				pairs := make([]*tikv.KvPair, len(lookup.keys))
				for i, k := range lookup.keys {
					pairs[i] = &tikv.KvPair{
						Key:   k,
						Value: valMap[hex.EncodeToString(k)],
					}
				}

				c.c.Connect(gctx, lookup.re.addr)
				res, err := grpc.Do(gctx, lookup.re.addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchPutResponse, error) {
					return tikv.NewTikvClient(conn).RawBatchPut(ctx, &tikv.RawBatchPutRequest{
						Context: lookup.re.ctx,
						Pairs:   pairs,
						Cf:      "default",
					}, copts...)
				})
				if err != nil {
					return errors.Errorf("BatchPut failed, err: %w", err)
				}
				if res.Error != "" {
					return errors.Errorf("BatchPut failed, message: %s", res.Error)
				}
				if res.RegionError != nil {
					if err = c.handleRegionError(gctx, res.RegionError, lookup.keys, false); err != nil {
						return err
					}
					fkMu.Lock()
					failedKeys = append(failedKeys, lookup.keys...)
					fkMu.Unlock()
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		if len(failedKeys) == 0 {
			return nil // success
		}

		if err := c.refresh(ctx, failedKeys); err != nil {
			return err
		}

		remaining = failedKeys
	}

	if len(remaining) != 0 {
		return errors.New("exceeded region error retry limit for BatchPut")
	}
	return nil
}

func (c *client) Delete(ctx context.Context, key []byte) error {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Delete"), apiName+"/Delete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	for range c.regionErrorRetryLimit {
		lookups, err := c.lookupAddrs(ctx, [][]byte{key})
		if err != nil {
			return err
		}

		for _, lookup := range lookups {
			c.c.Connect(ctx, lookup.re.addr)
			res, err := grpc.Do(ctx, lookup.re.addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawDeleteResponse, error) {
				return tikv.NewTikvClient(conn).RawDelete(ctx, &tikv.RawDeleteRequest{
					Context: lookup.re.ctx,
					Key:     key,
					Cf:      "default",
				}, copts...)
			})
			if err != nil {
				return errors.Errorf("Delete failed, err: %w", err)
			}
			if res.Error != "" {
				return errors.Errorf("Delete failed, message: %s", res.Error)
			}
			if res.RegionError != nil {
				if err = c.handleRegionError(ctx, res.RegionError, [][]byte{key}, true); err != nil {
					return err
				}
				// retry
				goto RETRY
			}
			return nil
		}
	RETRY:
	}
	return errors.Errorf("exceeded region error retry limit for key: %s", hex.EncodeToString(key))
}

func (c *client) BatchDelete(ctx context.Context, keys [][]byte) error {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/BatchDelete"), apiName+"/BatchDelete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	remaining := make([][]byte, len(keys))
	copy(remaining, keys)

	for range c.regionErrorRetryLimit {
		lookups, err := c.lookupAddrs(ctx, remaining)
		if err != nil {
			return err
		}

		var failedKeys [][]byte
		var fkMu sync.Mutex

		g, gctx := errgroup.WithContext(ctx)

		for _, lookup := range lookups {
			addr := lookup.re.addr
			lookup := lookup
			g.Go(func() error {
				c.c.Connect(gctx, addr)
				res, err := grpc.Do(gctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchDeleteResponse, error) {
					return tikv.NewTikvClient(conn).RawBatchDelete(ctx, &tikv.RawBatchDeleteRequest{
						Context: lookup.re.ctx,
						Keys:    lookup.keys,
						Cf:      "default",
					}, copts...)
				})
				if err != nil {
					return errors.Errorf("BatchDelete failed, err: %w", err)
				}
				if res.Error != "" {
					return errors.Errorf("BatchDelete failed, message: %s", res.Error)
				}
				if res.RegionError != nil {
					if err = c.handleRegionError(gctx, res.RegionError, lookup.keys, false); err != nil {
						return err
					}
					fkMu.Lock()
					failedKeys = append(failedKeys, lookup.keys...)
					fkMu.Unlock()
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		if len(failedKeys) == 0 {
			return nil // success
		}

		if err := c.refresh(ctx, failedKeys); err != nil {
			return err
		}

		remaining = failedKeys
	}

	if len(remaining) != 0 {
		return errors.New("exceeded region error retry limit for BatchDelete")
	}
	return nil
}
