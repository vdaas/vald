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

	"github.com/vdaas/vald/apis/grpc/v1/tikv"
	"github.com/vdaas/vald/internal/client/v1/client/meta"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	apiName = "vald/internal/client/meta/v1/client/meta/tikv"
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
}

// rangeEntry maps [start,end) to addr.
type rangeEntry struct {
	start []byte
	end   []byte // nil or empty means +Inf
	addr  string
	ver   uint64
}

func (c *client) lookupAddr(key []byte) string {
	idx := sort.Search(len(c.ranges), func(i int) bool {
		return bytes.Compare(c.ranges[i].start, key) > 0
	}) - 1
	if idx >= 0 {
		re := c.ranges[idx]
		if len(re.end) == 0 || bytes.Compare(key, re.end) < 0 {
			return re.addr
		}
	}
	return ""
}

func (c *client) lookupAddrs(ctx context.Context, keys [][]byte) (map[string][][]byte, error) {
	unknownKeys := make([][]byte, 0, len(keys))
	res := make(map[string][][]byte)
	c.rmu.RLock()
	for _, key := range keys {
		addr := c.lookupAddr(key)
		if addr == "" {
			unknownKeys = append(unknownKeys, key)
		}
		res[addr] = append(res[addr], key)
	}
	c.rmu.RUnlock()
	err := c.refresh(ctx, unknownKeys)
	if err != nil {
		return nil, err
	}
	for _, key := range unknownKeys {
		addr := c.lookupAddr(key)
		if addr == "" {
			return nil, errors.Errorf("address not found for key: %s", hex.EncodeToString(key))
		}
		res[addr] = append(res[addr], key)
	}
	return res, nil
}

func (c *client) refresh(ctx context.Context, keys [][]byte) error {
	if c.clusterId == 0 {
		res, err := c.pd.GetClusterInfo(ctx, &tikv.GetClusterInfoRequest{
			Header: &tikv.ResponseHeader{},
		})
		if err != nil {
			return err
		}
		if res.Header.Error != nil {
			return errors.New(res.Header.Error.Message)
		}
		c.clusterId = res.Header.ClusterId
	}
	storeIdToAddr := make(map[uint64]string)
	var regions []*tikv.Region
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		res, err := c.pd.GetAllStores(ctx, &tikv.GetAllStoresRequest{
			Header: &tikv.RequestHeader{ClusterId: c.clusterId},
			ExcludeTombstoneStores: true,
		})
		if err != nil {
			return err
		}
		for _, store := range res.Stores {
			storeIdToAddr[store.Id] = store.Address
		}
		return nil
	})
	g.Go(func() error {
		req := make([]*tikv.KeyRange, len(keys))
		for i, key := range keys {
			req[i] = &tikv.KeyRange{
				StartKey: key,
				EndKey:   nil,
			}
		}
		res, err := c.pd.BatchScanRegions(ctx, &tikv.BatchScanRegionsRequest{
			Header: &tikv.RequestHeader{ClusterId: c.clusterId},
			Ranges: req,
		})
		if err != nil {
			return err
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
		if _, ok := storeIdToAddr[r.Leader.StoreId]; !ok {
			return errors.Errorf("store id %d not found", r.Leader.StoreId)
		}
		start := slices.Clone(r.Region.StartKey)
		end := slices.Clone(r.Region.EndKey)
		re := &rangeEntry{
			start: start,
			end:   end,
			addr:  storeIdToAddr[r.Leader.StoreId],
			ver:   r.Region.RegionEpoch.Version,
		}
		c.regions[r.Region.Id] = re
		c.ranges = mergeByNewerVersion(c.ranges, re)
	}
	return nil
}

func mergeByNewerVersion(old []*rangeEntry, newE *rangeEntry) []*rangeEntry {
	out := make([]*rangeEntry, 0, len(old)+1)
	inserted := false

	for _, e := range old {
		if !overlap(e.start, e.end, newE.start, newE.end) {
			out = append(out, e)
			continue
		}
		switch {
		case newE.ver > e.ver:
			continue
		case newE.ver == e.ver:
			// Only addr is updated for the same version
			e.addr = newE.addr
			out = append(out, e)
			inserted = true
		default:
			return old
		}
	}
	if !inserted {
		out = append(out, newE)
	}
	sort.Slice(out, func(i, j int) bool {
		return bytes.Compare(out[i].start, out[j].start) < 0
	})
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

func (c *client) Get(ctx context.Context, key []byte) (val []byte, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Get"), apiName+"/Get")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	addrs, err := c.lookupAddrs(ctx, [][]byte{key})
	if err != nil {
		return nil, err
	}
	for range c.regionErrorRetryLimit {
		for addr := range addrs {
			c.c.Connect(ctx, addr)
			res, err := grpc.Do(ctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawGetResponse, error) {
				return tikv.NewTikvClient(conn).RawGet(ctx, &tikv.RawGetRequest{
					Key: key,
				}, copts...)
			})
			if err != nil {
				return nil, err
			}
			if res.Error != "" {
				return nil, errors.New(res.Error)
			}
			if res.NotFound {
				return nil, errNotFound
			}
			if res.RegionError != nil {
				log.Errorf("region error: %+v", res.RegionError)
				continue
			}
			return val, nil
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
	addrs, err := c.lookupAddrs(ctx, keys)
	if err != nil {
		return nil, err
	}
	g, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	kv := make(map[string][]byte, len(keys))
	for addr, keys := range addrs {
		g.Go(func() error {
			for range c.regionErrorRetryLimit {
				c.c.Connect(ctx, addr)
				res, err := grpc.Do(ctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchGetResponse, error) {
					return tikv.NewTikvClient(conn).RawBatchGet(ctx, &tikv.RawBatchGetRequest{
						Keys: keys,
					}, copts...)
				})
				if err != nil {
					return err
				}
				if res.RegionError != nil {
					log.Errorf("region error: %+v", res.RegionError)
					continue
				}
				for _, pair := range res.Pairs {
					if pair.Error != nil {
						return errors.Errorf("KeyError happened %+v", pair.Error)
					}
					mu.Lock()
					kv[hex.EncodeToString(pair.Key)] = pair.Value
					mu.Unlock()
				}
				return nil
			}
			return errors.Errorf("exceeded region error retry limit for address: %s", addr)
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	re := make([][]byte, len(keys))
	for i, key := range keys {
		re[i] = kv[hex.EncodeToString(key)]
	}
	return re, nil
}

func (c *client) Put(ctx context.Context, key, val []byte) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Put"), apiName+"/Put")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	addrs, err := c.lookupAddrs(ctx, [][]byte{key})
	if err != nil {
		return err
	}
	for addr := range addrs {
		for range c.regionErrorRetryLimit {
			c.c.Connect(ctx, addr)
			res, err := grpc.Do(ctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchPutResponse, error) {
				return tikv.NewTikvClient(conn).RawBatchPut(ctx, &tikv.RawBatchPutRequest{
					Pairs: []*tikv.KvPair{{
						Key:   key,
						Value: val,
					}},
				}, copts...)
			})
			if err != nil {
				return err
			}
			if res.Error != "" {
				return errors.New(res.Error)
			}
			if res.RegionError != nil {
				log.Errorf("region error: %+v", res.RegionError)
				continue
			}
			return nil
		}
		return errors.Errorf("exceeded region error retry limit for address: %s", addr)
	}
	return nil
}

func (c *client) BatchPut(ctx context.Context, keys, vals [][]byte) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/BatchPut"), apiName+"/BatchPut")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	kv := make(map[string][]byte, len(keys))
	for i, key := range keys {
		kv[hex.EncodeToString(key)] = vals[i]
	}
	addrs, err := c.lookupAddrs(ctx, keys)
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	for addr, keys := range addrs {
		g.Go(func() error {
			pairs := make([]*tikv.KvPair, len(keys))
			for i := range keys {
				pairs[i] = &tikv.KvPair{
					Key:   keys[i],
					Value: kv[hex.EncodeToString(keys[i])],
				}
			}
			for range c.regionErrorRetryLimit {
				c.c.Connect(ctx, addr)
				res, err := grpc.Do(ctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchPutResponse, error) {
					return tikv.NewTikvClient(conn).RawBatchPut(ctx, &tikv.RawBatchPutRequest{
						Pairs: pairs,
					}, copts...)
				})
				if err != nil {
					return err
				}
				if res.Error != "" {
					return errors.New(res.Error)
				}
				if res.RegionError != nil {
					log.Errorf("region error: %+v", res.RegionError)
					continue
				}
				return nil
			}
			return errors.Errorf("exceeded region error retry limit for address: %s", addr)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (c *client) Delete(ctx context.Context, key []byte) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Delete"), apiName+"/Delete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	addrs, err := c.lookupAddrs(ctx, [][]byte{key})
	if err != nil {
		return err
	}
	for addr := range addrs {
		for range c.regionErrorRetryLimit {
			c.c.Connect(ctx, addr)
			res, err := grpc.Do(ctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchDeleteResponse, error) {
				return tikv.NewTikvClient(conn).RawBatchDelete(ctx, &tikv.RawBatchDeleteRequest{
					Keys: [][]byte{key},
				}, copts...)
			})
			if err != nil {
				return err
			}
			if res.Error != "" {
				return errors.New(res.Error)
			}
			if res.RegionError != nil {
				log.Errorf("region error: %+v", res.RegionError)
				continue
			}
			return nil
		}
		return errors.Errorf("exceeded region error retry limit for address: %s", addr)
	}
	return nil
}

func (c *client) BatchDelete(ctx context.Context, keys [][]byte) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/BatchDelete"), apiName+"/BatchDelete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	addrs, err := c.lookupAddrs(ctx, keys)
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	for addr, keys := range addrs {
		g.Go(func() error {
			for range c.regionErrorRetryLimit {
				c.c.Connect(ctx, addr)
				res, err := grpc.Do(ctx, addr, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchDeleteResponse, error) {
					return tikv.NewTikvClient(conn).RawBatchDelete(ctx, &tikv.RawBatchDeleteRequest{
						Keys: keys,
					}, copts...)
				})
				if err != nil {
					return err
				}
				if res.Error != "" {
					return errors.New(res.Error)
				}
				if res.RegionError != nil {
					log.Errorf("region error: %+v", res.RegionError)
					continue
				}
				return nil
			}
			return errors.Errorf("exceeded region error retry limit for address: %s", addr)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
