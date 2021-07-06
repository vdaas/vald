//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package service manages the main logic of server.
package service

import (
	"context"
	"encoding/gob"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/vqueue"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
)

type NGT interface {
	Start(ctx context.Context) <-chan error
	Search(vec []float32, size uint32, epsilon, radius float32) ([]model.Distance, error)
	SearchByID(uuid string, size uint32, epsilon, radius float32) ([]model.Distance, error)
	Insert(uuid string, vec []float32) (err error)
	InsertMultiple(vecs map[string][]float32) (err error)
	Update(uuid string, vec []float32) (err error)
	UpdateMultiple(vecs map[string][]float32) (err error)
	Delete(uuid string) (err error)
	DeleteMultiple(uuids ...string) (err error)
	GetObject(uuid string) (vec []float32, err error)
	CreateIndex(ctx context.Context, poolSize uint32) (err error)
	SaveIndex(ctx context.Context) (err error)
	Exists(string) (uint32, bool)
	CreateAndSaveIndex(ctx context.Context, poolSize uint32) (err error)
	IsIndexing() bool
	IsSaving() bool
	Len() uint64
	NumberOfCreateIndexExecution() uint64
	NumberOfProactiveGCExecution() uint64
	UUIDs(context.Context) (uuids []string)
	DeleteVQueueBufferLen() uint64
	InsertVQueueBufferLen() uint64
	DeleteVQueueChannelLen() uint64
	InsertVQueueChannelLen() uint64
	GetDimensionSize() int
	Close(ctx context.Context) error
}

type ngt struct {
	// instances
	core core.NGT
	eg   errgroup.Group
	kvs  kvs.BidiMap
	vq   vqueue.Queue

	// statuses
	indexing  atomic.Value
	saving    atomic.Value
	lastNoice uint64 // last number of create index execution this value prevent unnecessary saveindex.

	// counters
	nocie uint64 // number of create index execution
	nogce uint64 // number of proactive GC execution

	// configurations
	inMem bool // in-memory mode
	dim   int  // dimension size
	alen  int  // auto indexing length

	lim  time.Duration // auto indexing time limit
	dur  time.Duration // auto indexing check duration
	sdur time.Duration // auto save index check duration

	minLit    time.Duration // minimum load index timeout
	maxLit    time.Duration // maximum load index timeout
	litFactor time.Duration // load index timeout factor

	enableProactiveGC bool // if this value is true, agent component will purge GC memory more proactive

	path string // index path

	poolSize uint32  // default pool size
	radius   float32 // default radius
	epsilon  float32 // default epsilon

	idelay time.Duration // initial delay duration
	dcd    bool          // disable commit daemon
}

const (
	kvsFileName = "ngt-meta.kvsdb"
)

func New(cfg *config.NGT, opts ...Option) (nn NGT, err error) {
	n := new(ngt)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(n); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	n.kvs = kvs.New()
	n.dim = cfg.Dimension

	err = n.initNGT(
		core.WithInMemoryMode(n.inMem),
		core.WithIndexPath(n.path),
		core.WithDefaultPoolSize(n.poolSize),
		core.WithDefaultRadius(n.radius),
		core.WithDefaultEpsilon(n.epsilon),
		core.WithDimension(cfg.Dimension),
		core.WithDistanceTypeByString(cfg.DistanceType),
		core.WithObjectTypeByString(cfg.ObjectType),
		core.WithBulkInsertChunkSize(cfg.BulkInsertChunkSize),
		core.WithCreationEdgeSize(cfg.CreationEdgeSize),
		core.WithSearchEdgeSize(cfg.SearchEdgeSize),
	)
	if err != nil {
		return nil, err
	}

	if n.dur == 0 || n.alen == 0 {
		n.dcd = true
	}
	if n.vq == nil {
		n.vq, err = vqueue.New(
			vqueue.WithErrGroup(n.eg),
			vqueue.WithInsertBufferSize(cfg.VQueue.InsertBufferSize),
			vqueue.WithDeleteBufferSize(cfg.VQueue.DeleteBufferSize),
			vqueue.WithInsertBufferPoolSize(cfg.VQueue.InsertBufferPoolSize),
			vqueue.WithDeleteBufferPoolSize(cfg.VQueue.DeleteBufferPoolSize),
		)
		if err != nil {
			return nil, err
		}
	}

	n.indexing.Store(false)
	n.saving.Store(false)

	return n, nil
}

func (n *ngt) initNGT(opts ...core.Option) (err error) {
	if n.inMem {
		log.Debug("vald agent starts with in-memory mode")
		n.core, err = core.New(opts...)
		return err
	}

	_, err = os.Stat(n.path)
	if os.IsNotExist(err) {
		log.Debugf("index file not exists,\tpath: %s,\terr: %v", n.path, err)
		n.core, err = core.New(opts...)
		return err
	}
	if os.IsPermission(err) {
		log.Errorf("no permission for index path,\tpath: %s,\terr: %v", n.path, err)
		return err
	}

	log.Debugf("load index from %s", n.path)

	agentMetadata, err := metadata.Load(filepath.Join(n.path, metadata.AgentMetadataFileName))
	if err != nil {
		log.Warnf("cannot read metadata from %s: %s", metadata.AgentMetadataFileName, err)
	}
	if os.IsNotExist(err) || agentMetadata == nil || agentMetadata.NGT == nil || agentMetadata.NGT.IndexCount == 0 {
		log.Warnf("cannot read metadata from %s: %v", metadata.AgentMetadataFileName, err)

		if fi, err := os.Stat(filepath.Join(n.path, kvsFileName)); os.IsNotExist(err) || fi.Size() == 0 {
			log.Warn("kvsdb file is not exist")
			n.core, err = core.New(opts...)
			return err
		}

		if os.IsPermission(err) {
			log.Errorf("no permission for kvsdb file,\tpath: %s,\terr: %v", filepath.Join(n.path, kvsFileName), err)
			return err
		}
	}

	var timeout time.Duration
	if agentMetadata != nil && agentMetadata.NGT != nil {
		log.Debugf("the backup index size is %d. starting to load...", agentMetadata.NGT.IndexCount)
		timeout = time.Duration(
			math.Min(
				math.Max(
					float64(agentMetadata.NGT.IndexCount)*float64(n.litFactor),
					float64(n.minLit),
				),
				float64(n.maxLit),
			),
		)
	} else {
		log.Debugf("cannot inspect the backup index size. starting to load default value.")
		timeout = time.Duration(math.Min(float64(n.minLit), float64(n.maxLit)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	eg, _ := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() (err error) {
		n.core, err = core.Load(opts...)
		return err
	}))

	eg.Go(safety.RecoverFunc(n.loadKVS))

	ech := make(chan error, 1)

	// NOTE: when it exceeds the timeout while loading,
	// it should exit this function and leave this goroutine running.
	go func() {
		defer close(ech)
		err = safety.RecoverFunc(func() (err error) {
			err = eg.Wait()
			if err != nil {
				return err
			}
			cancel()
			return nil
		})()
		if err != nil {
			ech <- err
		}
	}()

	select {
	case err := <-ech:
		return err
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Errorf("cannot load index backup data within the timeout %s. the process is going to be killed.", timeout)

			err := metadata.Store(
				filepath.Join(n.path, metadata.AgentMetadataFileName),
				&metadata.Metadata{
					IsInvalid: true,
					NGT: &metadata.NGT{
						IndexCount: 0,
					},
				},
			)
			if err != nil {
				return err
			}

			return errors.ErrIndexLoadTimeout
		}
	}

	return nil
}

func (n *ngt) loadKVS() error {
	gob.Register(map[string]uint32{})

	f, err := file.Open(
		filepath.Join(n.path, kvsFileName),
		os.O_RDONLY|os.O_SYNC,
		os.ModePerm,
	)
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Wrap(err, derr.Error())
			}
		}
	}()

	m := make(map[string]uint32)
	err = gob.NewDecoder(f).Decode(&m)
	if err != nil {
		log.Errorf("error decoding kvsdb file,\terr: %v", err)
		return err
	}

	for k, id := range m {
		n.kvs.Set(k, id)
	}

	return nil
}

func (n *ngt) Start(ctx context.Context) <-chan error {
	if n.dcd {
		return nil
	}
	ech := make(chan error, 2)
	n.eg.Go(safety.RecoverFunc(func() (err error) {
		vqech, err := n.vq.Start(ctx)
		if err != nil {
			return err
		}
		if n.sdur == 0 {
			n.sdur = n.dur + time.Second
		}
		if n.lim == 0 {
			n.lim = n.dur * 2
		}
		defer close(ech)

		timer := time.NewTimer(n.idelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		timer.Stop()

		tick := time.NewTicker(n.dur)
		sTick := time.NewTicker(n.sdur)
		limit := time.NewTicker(n.lim)
		defer tick.Stop()
		defer sTick.Stop()
		defer limit.Stop()
		for {
			err = nil
			select {
			case <-ctx.Done():
				err = n.CreateIndex(ctx, n.poolSize)
				if err != nil && !errors.Is(err, errors.ErrUncommittedIndexNotFound) {
					ech <- err
					return errors.Wrap(ctx.Err(), err.Error())
				}
				return ctx.Err()
			case <-tick.C:
				if n.vq.IVQLen() >= n.alen {
					err = n.CreateIndex(ctx, n.poolSize)
				}
			case <-limit.C:
				err = n.CreateAndSaveIndex(ctx, n.poolSize)
			case <-sTick.C:
				err = n.SaveIndex(ctx)
			case err := <-vqech:
				if err != nil {
					ech <- err
					err = nil
				}
			}
			if err != nil && err != errors.ErrUncommittedIndexNotFound {
				ech <- err
				runtime.Gosched()
				err = nil
			}
		}
	}))
	return ech
}

func (n *ngt) Search(vec []float32, size uint32, epsilon, radius float32) ([]model.Distance, error) {
	if n.IsIndexing() {
		return nil, errors.ErrCreateIndexingIsInProgress
	}
	sr, err := n.core.Search(vec, int(size), epsilon, radius)
	if err != nil {
		return nil, err
	}

	ds := make([]model.Distance, 0, len(sr))
	for _, d := range sr {
		if err = d.Error; d.ID == 0 && err != nil {
			log.Warnf("an error occurred while searching: %s", err)
			continue
		}
		key, ok := n.kvs.GetInverse(d.ID)
		if ok {
			ds = append(ds, model.Distance{
				ID:       key,
				Distance: d.Distance,
			})
		} else {
			log.Warn("not found", d.ID, d.Distance)
		}
	}

	return ds, nil
}

func (n *ngt) SearchByID(uuid string, size uint32, epsilon, radius float32) (dst []model.Distance, err error) {
	if n.IsIndexing() {
		return nil, errors.ErrCreateIndexingIsInProgress
	}
	log.Debugf("SearchByID\tuuid: %s size: %d epsilon: %f radius: %f", uuid, size, epsilon, radius)
	vec, err := n.GetObject(uuid)
	if err != nil {
		log.Debugf("SearchByID\tuuid: %s's vector not found", uuid)
		return nil, err
	}
	return n.Search(vec, size, epsilon, radius)
}

func (n *ngt) Insert(uuid string, vec []float32) (err error) {
	return n.insert(uuid, vec, time.Now().UnixNano(), true)
}

func (n *ngt) insert(uuid string, vec []float32, t int64, validation bool) (err error) {
	if len(uuid) == 0 {
		err = errors.ErrUUIDNotFound(0)
		return err
	}
	if validation && !n.vq.DVExists(uuid) {
		// if delete schedule exists we can insert new vector
		_, ok := n.kvs.Get(uuid)
		if ok || n.vq.IVExists(uuid) {
			return errors.ErrUUIDAlreadyExists(uuid)
		}
	}
	return n.vq.PushInsert(uuid, vec, t)
}

func (n *ngt) InsertMultiple(vecs map[string][]float32) (err error) {
	return n.insertMultiple(vecs, time.Now().UnixNano())
}

func (n *ngt) insertMultiple(vecs map[string][]float32, now int64) (err error) {
	for uuid, vec := range vecs {
		ierr := n.insert(uuid, vec, now, true)
		if ierr != nil {
			if err != nil {
				err = errors.Wrap(ierr, err.Error())
			} else {
				err = ierr
			}
		}
	}
	return err
}

func (n *ngt) Update(uuid string, vec []float32) (err error) {
	now := time.Now().UnixNano()
	if err = n.readyForUpdate(uuid, vec); err != nil {
		return err
	}
	err = n.delete(uuid, now)
	if err != nil {
		return err
	}
	now++
	return n.insert(uuid, vec, now, false)
}

func (n *ngt) UpdateMultiple(vecs map[string][]float32) (err error) {
	uuids := make([]string, 0, len(vecs))
	for uuid, vec := range vecs {
		if err = n.readyForUpdate(uuid, vec); err != nil {
			delete(vecs, uuid)
		} else {
			uuids = append(uuids, uuid)
		}
	}
	now := time.Now().UnixNano()
	err = n.deleteMultiple(uuids, now)
	if err != nil {
		return err
	}
	now++
	return n.insertMultiple(vecs, now)
}

func (n *ngt) Delete(uuid string) (err error) {
	return n.delete(uuid, time.Now().UnixNano())
}

func (n *ngt) delete(uuid string, t int64) (err error) {
	if len(uuid) == 0 {
		err = errors.ErrUUIDNotFound(0)
		return err
	}
	_, ok := n.kvs.Get(uuid)
	if !ok && !n.vq.IVExists(uuid) {
		return errors.ErrObjectIDNotFound(uuid)
	}
	return n.vq.PushDelete(uuid, t)
}

func (n *ngt) DeleteMultiple(uuids ...string) (err error) {
	return n.deleteMultiple(uuids, time.Now().UnixNano())
}

func (n *ngt) deleteMultiple(uuids []string, now int64) (err error) {
	for _, uuid := range uuids {
		ierr := n.delete(uuid, now)
		if ierr != nil {
			if err != nil {
				err = errors.Wrap(ierr, err.Error())
			} else {
				err = ierr
			}
		}
	}
	return err
}

func (n *ngt) GetObject(uuid string) (vec []float32, err error) {
	oid, ok := n.kvs.Get(uuid)
	if !ok {
		log.Debugf("GetObject\tuuid: %s's kvs data not found, trying to read from vqueue", uuid)
		vec, ok := n.vq.GetVector(uuid)
		if !ok {
			log.Debugf("GetObject\tuuid: %s's vqueue data not found", uuid)
			return nil, errors.ErrObjectIDNotFound(uuid)
		}
		return vec, nil
	}
	log.Debugf("GetObject\tGetVector oid: %d", oid)
	vec, err = n.core.GetVector(uint(oid))
	if err != nil {
		log.Debugf("GetObject\tuuid: %s oid: %d's vector not found", uuid, oid)
		return nil, errors.ErrObjectNotFound(err, uuid)
	}
	return vec, nil
}

func (n *ngt) CreateIndex(ctx context.Context, poolSize uint32) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt/service/NGT.CreateIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if n.IsIndexing() || n.IsSaving() {
		return nil
	}
	ic := n.vq.IVQLen() + n.vq.DVQLen()
	if ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}
	n.indexing.Store(true)
	defer n.indexing.Store(false)
	defer n.gc()

	log.Infof("create index operation started, uncommitted indexes = %d", ic)
	log.Debug("create index delete phase started")
	n.vq.RangePopDelete(ctx, func(uuid string) bool {
		var ierr error
		oid, ok := n.kvs.Delete(uuid)
		if ok {
			ierr = n.core.Remove(uint(oid))
		} else {
			ierr = errors.ErrObjectIDNotFound(uuid)
		}
		if ierr != nil {
			log.Error(ierr)
			err = errors.Wrap(err, ierr.Error())
		}
		return true
	})
	log.Debug("create index delete phase finished")
	n.gc()
	log.Debug("create index insert phase started")
	n.vq.RangePopInsert(ctx, func(uuid string, vector []float32) bool {
		oid, ierr := n.core.Insert(vector)
		if ierr != nil {
			log.Error(ierr)
			err = errors.Wrap(err, ierr.Error())
		} else {
			n.kvs.Set(uuid, uint32(oid))
		}
		return true
	})
	log.Debug("create index insert phase finished")
	log.Debug("create graph and tree phase started")
	log.Debugf("pool size = %d", poolSize)
	cierr := n.core.CreateIndex(poolSize)
	if cierr != nil {
		log.Error("an error occurred on creating graph and tree phase:", cierr)
		err = errors.Wrap(err, cierr.Error())
	}
	log.Debug("create graph and tree phase finished")

	log.Info("create index operation finished")
	atomic.AddUint64(&n.nocie, 1)
	return err
}

func (n *ngt) SaveIndex(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt/service/NGT.SaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(n.path) != 0 && !n.inMem {
		err = n.saveIndex(ctx)
	}

	return
}

func (n *ngt) saveIndex(ctx context.Context) (err error) {
	noice := atomic.LoadUint64(&n.nocie)
	if atomic.LoadUint64(&n.lastNoice) == noice {
		return
	}
	atomic.SwapUint64(&n.lastNoice, noice)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	// wait for not indexing & not saving
	for n.IsIndexing() || n.IsSaving() {
		runtime.Gosched()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
	n.saving.Store(true)
	defer n.gc()
	defer n.saving.Store(false)

	eg, ctx := errgroup.New(ctx)

	// we want to ensure the acutal kvs size between kvsdb and metadata,
	// so we create thie counter to count the actual kvs size instead of using kvs.Len()
	var kvsLen uint64

	eg.Go(safety.RecoverFunc(func() (err error) {
		if n.path != "" {
			m := make(map[string]uint32, n.Len())
			var mu sync.Mutex
			n.kvs.Range(ctx, func(key string, id uint32) bool {
				mu.Lock()
				m[key] = id
				kvsLen++
				mu.Unlock()
				return true
			})
			var f *os.File
			f, err = file.Open(
				filepath.Join(n.path, kvsFileName),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				fs.ModePerm,
			)
			if err != nil {
				return err
			}
			defer func() {
				if f != nil {
					derr := f.Close()
					if derr != nil {
						err = errors.Wrap(err, derr.Error())
						log.Errorf("[rebalance controller] failed to close kvsdb file: %v", err)
					}
				}
			}()
			gob.Register(map[string]uint32{})
			err = gob.NewEncoder(f).Encode(&m)
			if err != nil {
				return err
			}
			err = f.Sync()
			if err != nil {
				return err
			}
		}
		return nil
	}))

	eg.Go(safety.RecoverFunc(func() error {
		return n.core.SaveIndex()
	}))

	err = eg.Wait()
	if err != nil {
		return err
	}

	// TODO: delete this code
	time.Sleep(10 * time.Second)

	log.Infof("[rebalance controller] metadata area. kvs length: %d", n.kvs.Len())

	return metadata.Store(
		filepath.Join(n.path, metadata.AgentMetadataFileName),
		&metadata.Metadata{
			IsInvalid: false,
			NGT: &metadata.NGT{
				IndexCount: kvsLen,
			},
		},
	)
}

func (n *ngt) CreateAndSaveIndex(ctx context.Context, poolSize uint32) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt/service/NGT.CreateAndSaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = n.CreateIndex(ctx, poolSize)
	if err != nil {
		return err
	}
	return n.SaveIndex(ctx)
}

func (n *ngt) Exists(uuid string) (oid uint32, ok bool) {
	oid, ok = n.kvs.Get(uuid)
	if !ok {
		ok = n.vq.IVExists(uuid)
	} else {
		ok = !n.vq.DVExists(uuid)
	}
	return oid, ok
}

func (n *ngt) readyForUpdate(uuid string, vec []float32) (err error) {
	if len(uuid) == 0 {
		return errors.ErrUUIDNotFound(0)
	}
	if len(vec) == 0 {
		return errors.ErrInvalidDimensionSize(len(vec), 0)
	}
	ovec, err := n.GetObject(uuid)
	if err != nil || len(vec) != len(ovec) {
		// if error (GetObject cannot find vector) or vector length is not equal let's try update
		return nil
	}
	for i, v := range vec {
		if v != ovec[i] {
			// if difference exists return nil for update
			return nil
		}
	}
	// if no difference exists (same vector already exists) return error for skip update
	return errors.ErrUUIDAlreadyExists(uuid)
}

func (n *ngt) IsSaving() bool {
	s, ok := n.saving.Load().(bool)
	return s && ok
}

func (n *ngt) IsIndexing() bool {
	i, ok := n.indexing.Load().(bool)
	return i && ok
}

func (n *ngt) UUIDs(ctx context.Context) (uuids []string) {
	uuids = make([]string, 0, n.kvs.Len())
	n.kvs.Range(ctx, func(uuid string, oid uint32) bool {
		uuids = append(uuids, uuid)
		return true
	})
	return uuids
}

func (n *ngt) NumberOfCreateIndexExecution() uint64 {
	return atomic.LoadUint64(&n.nocie)
}

func (n *ngt) NumberOfProactiveGCExecution() uint64 {
	return atomic.LoadUint64(&n.nogce)
}

func (n *ngt) gc() {
	if n.enableProactiveGC {
		runtime.GC()
		atomic.AddUint64(&n.nogce, 1)
	}
}

func (n *ngt) Len() uint64 {
	return n.kvs.Len()
}

func (n *ngt) InsertVQueueBufferLen() uint64 {
	return uint64(n.vq.IVQLen())
}

func (n *ngt) DeleteVQueueBufferLen() uint64 {
	return uint64(n.vq.DVQLen())
}

func (n *ngt) InsertVQueueChannelLen() uint64 {
	return uint64(n.vq.IVCLen())
}

func (n *ngt) DeleteVQueueChannelLen() uint64 {
	return uint64(n.vq.DVCLen())
}

func (n *ngt) GetDimensionSize() int {
	return n.dim
}

func (n *ngt) Close(ctx context.Context) (err error) {
	if len(n.path) != 0 {
		err = n.CreateAndSaveIndex(ctx, n.poolSize)
	}
	n.core.Close()
	return
}
