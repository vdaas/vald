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
	UncommittedUUIDs() (uuids []string)
	DeleteVCacheLen() uint64
	InsertVCacheLen() uint64
	Close(ctx context.Context) error
}

type ngt struct {
	// instances
	core core.NGT
	eg   errgroup.Group
	kvs  kvs.BidiMap
	ivc  *vcaches // insertion vector cache
	dvc  *vcaches // deletion vector cache

	// statuses
	indexing  atomic.Value
	saving    atomic.Value
	lastNoice uint64 // last number of create index execution this value prevent unnecessary saveindex.

	// counters
	ic    uint64 // insert count
	nocie uint64 // number of create index execution
	nogce uint64 // number of proactive GC execution

	// configurations
	inMem bool // in-memory mode

	alen int // auto indexing length

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

type vcache struct {
	vector []float32
	date   int64
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
	if n.ivc == nil {
		n.ivc = new(vcaches)
	}
	if n.dvc == nil {
		n.dvc = new(vcaches)
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
		log.Debugf("no permission for index path,\tpath: %s,\terr: %v", n.path, err)
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
			log.Debugf("no permission for kvsdb file,\tpath: %s,\terr: %v", filepath.Join(n.path, kvsFileName), err)
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

	defer f.Close()

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
				if int(atomic.LoadUint64(&n.ic)) >= n.alen {
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
		return make([]model.Distance, 0), nil
	}
	sr, err := n.core.Search(vec, int(size), epsilon, radius)
	if err != nil {
		return nil, err
	}

	ds := make([]model.Distance, 0, len(sr))
	for _, d := range sr {
		if err = d.Error; d.ID == 0 && err != nil {
			log.Debug("an error occurred while searching:", err)
			continue
		}
		key, ok := n.kvs.GetInverse(d.ID)
		if ok {
			ds = append(ds, model.Distance{
				ID:       key,
				Distance: d.Distance,
			})
		}
	}

	return ds, nil
}

func (n *ngt) SearchByID(uuid string, size uint32, epsilon, radius float32) (dst []model.Distance, err error) {
	if n.IsIndexing() {
		log.Debug("SearchByID\t now indexing...")
		return make([]model.Distance, 0), nil
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
	if validation {
		id, ok := n.kvs.Get(uuid)
		if ok {
			err = errors.ErrUUIDAlreadyExists(uuid, uint(id))
			return err
		}
		_, ok = n.insertCache(uuid)
		if ok {
			err = errors.ErrUUIDAlreadyExists(uuid, uint(id))
			return err
		}
	}
	n.ivc.Store(uuid, vcache{
		vector: vec,
		date:   t,
	})
	atomic.AddUint64(&n.ic, 1)
	return nil
}

func (n *ngt) InsertMultiple(vecs map[string][]float32) (err error) {
	t := time.Now().UnixNano()
	for uuid, vec := range vecs {
		ierr := n.insert(uuid, vec, t, true)
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
	if !n.readyForUpdate(uuid, vec) {
		return nil
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
		if n.readyForUpdate(uuid, vec) {
			uuids = append(uuids, uuid)
		} else {
			delete(vecs, uuid)
		}
	}
	err = n.DeleteMultiple(uuids...)
	if err != nil {
		for _, uuid := range uuids {
			n.dvc.Delete(uuid)
		}
		return err
	}
	t := time.Now().UnixNano()
	for uuid, vec := range vecs {
		ierr := n.insert(uuid, vec, t, false)
		if ierr != nil {
			n.dvc.Delete(uuid)
			n.ivc.Delete(uuid)
			atomic.AddUint64(&n.ic, ^uint64(0))
			if err != nil {
				err = errors.Wrap(ierr, err.Error())
			} else {
				err = ierr
			}
		}
	}
	return err
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
	if !ok {
		_, ok := n.insertCache(uuid)
		if !ok {
			err = errors.ErrObjectIDNotFound(uuid)
			return err
		}
	}
	if vc, ok := n.ivc.Load(uuid); ok && vc.date < t {
		n.ivc.Delete(uuid)
		atomic.AddUint64(&n.ic, ^uint64(0))
	}
	n.dvc.Store(uuid, vcache{
		date: t,
	})
	return nil
}

func (n *ngt) DeleteMultiple(uuids ...string) (err error) {
	t := time.Now().UnixNano()
	for _, uuid := range uuids {
		ierr := n.delete(uuid, t)
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
		log.Debugf("GetObject\tuuid: %s's kvs data not found, trying to read from vcache", uuid)
		ivc, ok := n.insertCache(uuid)
		if !ok {
			log.Debugf("GetObject\tuuid: %s's vcache data not found", uuid)
			return nil, errors.ErrObjectIDNotFound(uuid)
		}
		return ivc.vector, nil
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
	ic := atomic.LoadUint64(&n.ic)
	if ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}
	n.indexing.Store(true)
	atomic.StoreUint64(&n.ic, 0)
	t := time.Now().UnixNano()
	defer n.indexing.Store(false)
	defer n.gc()

	log.Infof("create index operation started, uncommitted indexes = %d", ic)
	delList := make([]string, 0, ic)
	n.dvc.Range(func(uuid string, dvc vcache) bool {
		if dvc.date > t {
			return true
		}
		if ivc, ok := n.ivc.Load(uuid); ok && ivc.date < t && ivc.date < dvc.date {
			n.ivc.Delete(uuid)
		}
		delList = append(delList, uuid)
		return true
	})
	log.Info("create index delete kvs phase started")
	log.Debugf("deleting kvs: %#v", delList)
	doids := make([]uint, 0, ic)
	for _, duuid := range delList {
		n.dvc.Delete(duuid)
		id, ok := n.kvs.Delete(duuid)
		if !ok {
			log.Error(errors.ErrObjectIDNotFound(duuid).Error())
			err = errors.Wrap(err, errors.ErrObjectIDNotFound(duuid).Error())
		} else {
			doids = append(doids, uint(id))
		}
	}
	log.Info("create index delete kvs phase finished")

	log.Info("create index delete index phase started")
	log.Debugf("deleting index: %#v", doids)
	brerr := n.core.BulkRemove(doids...)
	log.Info("create index delete index phase finished")
	if brerr != nil {
		log.Error("an error occurred on deleting index phase:", brerr)
		err = errors.Wrap(err, brerr.Error())
	}
	uuids := make([]string, 0, ic)
	vecs := make([][]float32, 0, ic)
	n.ivc.Range(func(uuid string, ivc vcache) bool {
		if ivc.date <= t {
			uuids = append(uuids, uuid)
			vecs = append(vecs, ivc.vector)
		}
		return true
	})
	for _, uuid := range uuids {
		n.ivc.Delete(uuid)
	}
	n.gc()
	log.Info("create index insert index phase started")
	log.Debugf("inserting index: %#v", vecs)
	oids, errs := n.core.BulkInsert(vecs)
	log.Info("create index insert index phase finished")
	if errs != nil && len(errs) != 0 {
		for _, bierr := range errs {
			if bierr != nil {
				log.Error("an error occurred on inserting index phase:", bierr)
				err = errors.Wrap(err, bierr.Error())
			}
		}
	}

	log.Info("create index insert kvs phase started")
	log.Debugf("uuids = %#v\t\toids = %#v", uuids, oids)
	for i, uuid := range uuids {
		if len(oids) > i {
			oid := uint32(oids[i])
			if oid != 0 {
				n.kvs.Set(uuid, oid)
			}
		}
	}
	log.Info("create index insert kvs phase finished")

	log.Info("create graph and tree phase started")
	log.Debugf("pool size = %d", poolSize)
	cierr := n.core.CreateIndex(poolSize)
	if cierr != nil {
		log.Error("an error occurred on creating graph and tree phase:", cierr)
		err = errors.Wrap(err, cierr.Error())
	}
	log.Info("create graph and tree phase finished")

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

	eg.Go(safety.RecoverFunc(func() error {
		if n.path != "" {
			m := make(map[string]uint32, n.kvs.Len())
			var mu sync.Mutex
			n.kvs.Range(ctx, func(key string, id uint32) bool {
				mu.Lock()
				m[key] = id
				mu.Unlock()
				return true
			})
			f, err := file.Open(
				filepath.Join(n.path, kvsFileName),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				os.ModePerm,
			)
			if err != nil {
				return err
			}
			defer f.Close()
			gob.Register(map[string]uint32{})
			return gob.NewEncoder(f).Encode(&m)
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
				IndexCount: n.Len(),
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
		_, ok = n.insertCache(uuid)
	}
	return oid, ok
}

func (n *ngt) readyForUpdate(uuid string, vec []float32) (ready bool) {
	if len(uuid) == 0 || len(vec) == 0 {
		return false
	}
	ovec, err := n.GetObject(uuid)
	if err != nil || len(vec) != len(ovec) {
		// if error (GetObject cannot find vector) or vector length is not equal let's try update
		return true
	}
	for i, v := range vec {
		if v != ovec[i] {
			// if difference exists return true for update
			return true
		}
	}
	// if no difference exists (same vector already exists) return false for skip update
	return false
}

func (n *ngt) insertCache(uuid string) (*vcache, bool) {
	iv, ok := n.ivc.Load(uuid)
	if ok {
		dv, ok := n.dvc.Load(uuid)
		if !ok {
			return &iv, true
		}
		if ok && dv.date <= iv.date {
			return &iv, true
		}
		n.ivc.Delete(uuid)
		atomic.AddUint64(&n.ic, ^uint64(0))
	}
	return nil, false
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

func (n *ngt) UncommittedUUIDs() (uuids []string) {
	var mu sync.Mutex
	uuids = make([]string, 0, atomic.LoadUint64(&n.ic))
	n.ivc.Range(func(uuid string, vc vcache) bool {
		mu.Lock()
		uuids = append(uuids, uuid)
		mu.Unlock()
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

func (n *ngt) InsertVCacheLen() uint64 {
	return n.ivc.Len()
}

func (n *ngt) DeleteVCacheLen() uint64 {
	return n.dvc.Len()
}

func (n *ngt) Close(ctx context.Context) (err error) {
	if len(n.path) != 0 {
		err = n.CreateAndSaveIndex(ctx, n.poolSize)
	}
	n.core.Close()
	return
}
