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
	"unsafe"

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
	InsertWithTime(uuid string, vec []float32, t int64) (err error)
	InsertMultiple(vecs map[string][]float32) (err error)
	InsertMultipleWithTime(vecs map[string][]float32, t int64) (err error)
	Update(uuid string, vec []float32) (err error)
	UpdateWithTime(uuid string, vec []float32, t int64) (err error)
	UpdateMultiple(vecs map[string][]float32) (err error)
	UpdateMultipleWithTime(vecs map[string][]float32, t int64) (err error)
	Delete(uuid string) (err error)
	DeleteWithTime(uuid string, t int64) (err error)
	DeleteMultiple(uuids ...string) (err error)
	DeleteMultipleWithTime(uuids []string, t int64) (err error)
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
	cimu      sync.Mutex // create index mutex
	lastNoice uint64     // last number of create index execution this value prevent unnecessary saveindex.

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

	n.dim = cfg.Dimension

	n.kvs = kvs.New(kvs.WithConcurrency(cfg.KVSDB.Concurrency))

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

	if exist, _, err := file.ExistsWithDetail(n.path); !exist {
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
		if exist, fi, err := file.ExistsWithDetail(filepath.Join(n.path, kvsFileName)); !exist || fi != nil && fi.Size() == 0 {
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
	n.eg.Go(safety.RecoverFunc(func() error {
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
		return nil
	}))

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
		defer close(ech)
		if n.dur <= 0 {
			n.dur = math.MaxInt64
		}
		if n.sdur <= 0 {
			n.sdur = math.MaxInt64
		}
		if n.lim <= 0 {
			n.lim = math.MaxInt64
		}

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
		log.Errorf("cgo error detected: ngt code api returned error %v", err)
		if n.IsIndexing() {
			return nil, errors.ErrCreateIndexingIsInProgress
		}
		return nil, err
	}

	if len(sr) == 0 {
		return nil, errors.ErrEmptySearchResult
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

func (n *ngt) InsertWithTime(uuid string, vec []float32, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.insert(uuid, vec, t, true)
}

func (n *ngt) insert(uuid string, vec []float32, t int64, validation bool) (err error) {
	if len(uuid) == 0 {
		err = errors.ErrUUIDNotFound(0)
		return err
	}
	if validation {
		_, ok := n.Exists(uuid)
		if ok {
			return errors.ErrUUIDAlreadyExists(uuid)
		}
	}
	return n.vq.PushInsert(uuid, vec, t)
}

func (n *ngt) InsertMultiple(vecs map[string][]float32) (err error) {
	return n.insertMultiple(vecs, time.Now().UnixNano())
}

func (n *ngt) InsertMultipleWithTime(vecs map[string][]float32, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.insertMultiple(vecs, t)
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
	return n.update(uuid, vec, time.Now().UnixNano())
}

func (n *ngt) UpdateWithTime(uuid string, vec []float32, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.update(uuid, vec, t)
}

func (n *ngt) update(uuid string, vec []float32, t int64) (err error) {
	if err = n.readyForUpdate(uuid, vec); err != nil {
		return err
	}
	err = n.delete(uuid, t)
	if err != nil {
		return err
	}
	t++
	return n.insert(uuid, vec, t, false)
}

func (n *ngt) UpdateMultiple(vecs map[string][]float32) (err error) {
	return n.updateMultiple(vecs, time.Now().UnixNano())
}

func (n *ngt) UpdateMultipleWithTime(vecs map[string][]float32, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.updateMultiple(vecs, t)
}

func (n *ngt) updateMultiple(vecs map[string][]float32, t int64) (err error) {
	uuids := make([]string, 0, len(vecs))
	for uuid, vec := range vecs {
		if err = n.readyForUpdate(uuid, vec); err != nil {
			delete(vecs, uuid)
		} else {
			uuids = append(uuids, uuid)
		}
	}
	err = n.deleteMultiple(uuids, t)
	if err != nil {
		return err
	}
	t++
	return n.insertMultiple(vecs, t)
}

func (n *ngt) Delete(uuid string) (err error) {
	return n.delete(uuid, time.Now().UnixNano())
}

func (n *ngt) DeleteWithTime(uuid string, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.delete(uuid, t)
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

func (n *ngt) DeleteMultipleWithTime(uuids []string, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.deleteMultiple(uuids, t)
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

func (n *ngt) CreateIndex(ctx context.Context, poolSize uint32) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt/service/NGT.CreateIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ic := n.vq.IVQLen() + n.vq.DVQLen()
	if ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}
	err = func() error {
		ticker := time.NewTicker(time.Millisecond * 100)
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
		return nil
	}()
	if err != nil {
		return err
	}
	n.cimu.Lock()
	defer n.cimu.Unlock()
	n.indexing.Store(true)
	now := time.Now().UnixNano()
	defer n.indexing.Store(false)
	defer n.gc()
	ic = n.vq.IVQLen() + n.vq.DVQLen()
	if ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}
	log.Infof("create index operation started, uncommitted indexes = %d", ic)
	log.Debug("create index delete phase started")
	n.vq.RangePopDelete(ctx, now, func(uuid string) bool {
		oid, ok := n.kvs.Delete(uuid)
		if !ok {
			log.Warn(errors.ErrObjectIDNotFound(uuid))
			return true
		}
		if err := n.core.Remove(uint(oid)); err != nil {
			log.Error(err)
		}
		return true
	})
	log.Debug("create index delete phase finished")
	n.gc()
	log.Debug("create index insert phase started")
	n.vq.RangePopInsert(ctx, now, func(uuid string, vector []float32) bool {
		oid, err := n.core.Insert(vector)
		if err != nil {
			log.Error(err)
		} else {
			n.kvs.Set(uuid, uint32(oid))
		}
		return true
	})
	log.Debug("create index insert phase finished")
	log.Debug("create graph and tree phase started")
	log.Debugf("pool size = %d", poolSize)
	err = n.core.CreateIndex(poolSize)
	if err != nil {
		log.Error("an error occurred on creating graph and tree phase:", err)
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
		return n.saveIndex(ctx)
	}
	return nil
}

func (n *ngt) saveIndex(ctx context.Context) (err error) {
	noice := atomic.LoadUint64(&n.nocie)
	if atomic.LoadUint64(&n.lastNoice) == noice {
		return
	}
	atomic.SwapUint64(&n.lastNoice, noice)
	err = func() error {
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
		return nil
	}()
	if err != nil {
		return err
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
	if err != nil &&
		!errors.Is(err, errors.ErrUncommittedIndexNotFound) &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return n.SaveIndex(ctx)
}

func (n *ngt) Exists(uuid string) (oid uint32, ok bool) {
	ok = n.vq.IVExists(uuid)
	if !ok {
		oid, ok = n.kvs.Get(uuid)
		if !ok {
			log.Debugf("Exists\tuuid: %s's data not found in kvsdb and insert vqueue\terror: %v", uuid, errors.ErrObjectIDNotFound(uuid))
			return 0, false
		}
		if n.vq.DVExists(uuid) {
			log.Debugf("Exists\tuuid: %s's data found in kvsdb and not found in insert vqueue, but delete vqueue data exists. the object will be delete soon\terror: %v",
				uuid, errors.ErrObjectIDNotFound(uuid))
			return 0, false
		}
	}
	return oid, ok
}

func (n *ngt) GetObject(uuid string) (vec []float32, err error) {
	vec, ok := n.vq.GetVector(uuid)
	if !ok {
		oid, ok := n.kvs.Get(uuid)
		if !ok {
			log.Debugf("GetObject\tuuid: %s's data not found in kvsdb and insert vqueue", uuid)
			return nil, errors.ErrObjectIDNotFound(uuid)
		}
		if n.vq.DVExists(uuid) {
			log.Debugf("GetObject\tuuid: %s's data found in kvsdb and not found in insert vqueue, but delete vqueue data exists. the object will be delete soon", uuid)
			return nil, errors.ErrObjectIDNotFound(uuid)
		}
		vec, err = n.core.GetVector(uint(oid))
		if err != nil {
			log.Debugf("GetObject\tuuid: %s oid: %d's vector not found in ngt index", uuid, oid)
			return nil, errors.ErrObjectNotFound(err, uuid)
		}
	}
	return vec, nil
}

func (n *ngt) readyForUpdate(uuid string, vec []float32) (err error) {
	if len(uuid) == 0 {
		return errors.ErrUUIDNotFound(0)
	}
	if len(vec) != n.GetDimensionSize() {
		return errors.ErrInvalidDimensionSize(len(vec), n.GetDimensionSize())
	}
	ovec, err := n.GetObject(uuid)
	if err != nil ||
		len(vec) != len(ovec) ||
		f32stos(vec) != f32stos(ovec) {
		// if error (GetObject cannot find vector) or vector length is not equal or if difference exists let's try update
		return nil
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

func (n *ngt) GetDimensionSize() int {
	return n.dim
}

func (n *ngt) Close(ctx context.Context) (err error) {
	err = n.kvs.Close()
	if len(n.path) != 0 {
		cerr := n.CreateIndex(ctx, n.poolSize)
		if cerr != nil &&
			!errors.Is(err, errors.ErrUncommittedIndexNotFound) &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			if err != nil {
				err = errors.Wrap(cerr, err.Error())
			} else {
				err = cerr
			}
		}
		serr := n.SaveIndex(ctx)
		if serr != nil &&
			!errors.Is(err, errors.ErrUncommittedIndexNotFound) &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			if err != nil {
				err = errors.Wrap(serr, err.Error())
			} else {
				err = serr
			}
		}
	}
	n.core.Close()
	return
}

func f32stos(fs []float32) string {
	lf := 4 * len(fs)
	buf := (*(*[1]byte)(unsafe.Pointer(&(fs[0]))))[:]
	addr := unsafe.Pointer(&buf)
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(8)))) = lf
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(16)))) = lf
	return *(*string)(unsafe.Pointer(&buf))
}
