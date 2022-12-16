//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/conv"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/vqueue"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
)

type NGT interface {
	Start(ctx context.Context) <-chan error
	Search(vec []float32, size uint32, epsilon, radius float32) ([]model.Distance, error)
	SearchByID(uuid string, size uint32, epsilon, radius float32) ([]float32, []model.Distance, error)
	LinearSearch(vec []float32, size uint32) ([]model.Distance, error)
	LinearSearchByID(uuid string, size uint32) ([]float32, []model.Distance, error)
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
	RegenerateIndex(ctx context.Context) (err error)
	GetObject(uuid string) (vec []float32, err error)
	CreateIndex(ctx context.Context, poolSize uint32) (err error)
	SaveIndex(ctx context.Context) (err error)
	Exists(string) (uint32, bool)
	CreateAndSaveIndex(ctx context.Context, poolSize uint32) (err error)
	IsIndexing() bool
	IsFlushing() bool
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
	fmu  sync.Mutex
	fmap map[string]uint32 // failure map for index
	vq   vqueue.Queue

	// statuses
	indexing  atomic.Value
	flushing  atomic.Value
	saving    atomic.Value
	cimu      sync.Mutex // create index mutex
	lastNocie uint64     // last number of create index execution this value prevent unnecessary saveindex.

	// counters
	nocie uint64 // number of create index execution
	nogce uint64 // number of proactive GC execution

	// parameters
	cfg  *config.NGT
	opts []Option

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
	enableCopyOnWrite bool // if this value is true, agent component will write backup file using Copy on Write and saves old files to the old directory

	path      string       // index path
	smu       sync.Mutex   // save index lock
	tmpPath   atomic.Value // temporary index path for Copy on Write
	oldPath   string       // old volume path
	basePath  string       // index base directory for CoW
	cowmu     sync.Mutex   // copy on write move lock
	backupGen uint64       // number of backup generation

	poolSize uint32  // default pool size
	radius   float32 // default radius
	epsilon  float32 // default epsilon

	idelay time.Duration // initial delay duration
	dcd    bool          // disable commit daemon

	kvsdbConcurrency int // kvsdb concurrency
}

const (
	kvsFileName = "ngt-meta.kvsdb"

	oldIndexDirName    = "backup"
	originIndexDirName = "origin"
)

func New(cfg *config.NGT, opts ...Option) (nn NGT, err error) {
	return newNGT(cfg, opts...)
}

func newNGT(cfg *config.NGT, opts ...Option) (n *ngt, err error) {
	n = &ngt{
		fmap:              make(map[string]uint32),
		dim:               cfg.Dimension,
		enableProactiveGC: cfg.EnableProactiveGC,
		enableCopyOnWrite: cfg.EnableCopyOnWrite,
		kvsdbConcurrency:  cfg.KVSDB.Concurrency,
		cfg:               cfg,
		opts:              opts,
	}

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(n); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if len(n.path) == 0 {
		n.inMem = true
	}

	if n.enableCopyOnWrite && !n.inMem && len(n.path) != 0 {
		sep := string(os.PathSeparator)
		n.path, err = filepath.Abs(strings.ReplaceAll(n.path, sep+sep, sep))
		if err != nil {
			log.Warn(err)
		}
		n.basePath = n.path
		n.oldPath = file.Join(n.basePath, oldIndexDirName)
		n.path = file.Join(n.basePath, originIndexDirName)
		err = file.MkdirAll(n.oldPath, fs.ModePerm)
		if err != nil {
			log.Warn(err)
		}
		err = file.MkdirAll(n.path, fs.ModePerm)
		if err != nil {
			log.Warn(err)
		}
		err = n.mktmp()
		if err != nil {
			return nil, err
		}
	}

	err = n.initNGT(
		core.WithInMemoryMode(n.inMem),
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

func (n *ngt) load(ctx context.Context, path string, opts ...core.Option) (err error) {
	exist, fi, err := file.ExistsWithDetail(path)
	switch {
	case !exist, fi == nil, fi != nil && fi.Size() == 0, err != nil && errors.Is(err, fs.ErrNotExist):
		err = errors.Wrapf(errors.ErrIndexFileNotFound, "index file does not exists,\tpath: %s,\terr: %v", path, err)
		return err
	case err != nil && errors.Is(err, fs.ErrPermission):
		if fi != nil {
			err = errors.Wrap(errors.ErrFailedToOpenFile(err, path, 0, fi.Mode()), "invalid permission for loading index path")
		}
		return err
	case exist && fi != nil && fi.IsDir():
		if fi.Mode().IsDir() && !strings.HasSuffix(path, string(os.PathSeparator)) {
			path += string(os.PathSeparator)
		}
		files, err := filepath.Glob(file.Join(filepath.Dir(path), "*"))
		if err != nil || len(files) == 0 {
			err = errors.Wrapf(errors.ErrIndexFileNotFound, "index path exists but no file does not exists in the directory,\tpath: %s,\tfiles: %v\terr: %v", path, files, err)
			return err
		}
		if strings.HasSuffix(path, string(os.PathSeparator)) {
			path = strings.TrimSuffix(path, string(os.PathSeparator))
		}
	}
	metadataPath := file.Join(path, metadata.AgentMetadataFileName)
	log.Debugf("index path: %s exists, now starting to check metadata from %s", path, metadataPath)
	exist, fi, err = file.ExistsWithDetail(metadataPath)
	switch {
	case !exist, fi == nil, fi != nil && fi.Size() == 0, err != nil && errors.Is(err, fs.ErrNotExist):
		err = errors.Wrapf(errors.ErrIndexFileNotFound, "metadata file does not exists,\tpath: %s,\terr: %v", metadataPath, err)
		return err
	case err != nil && errors.Is(err, fs.ErrPermission):
		if fi != nil {
			err = errors.Wrap(errors.ErrFailedToOpenFile(err, metadataPath, 0, fi.Mode()), "invalid permission for loading metadata")
		}
		return err
	}

	log.Debugf("index path: %s and metadata: %s exists, now starting to load metadata", path, metadataPath)
	agentMetadata, err := metadata.Load(metadataPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) || agentMetadata == nil || agentMetadata.NGT == nil || agentMetadata.NGT.IndexCount == 0 {
		err = errors.Wrapf(err, "cannot read metadata from path: %s\tmetadata: %s", path, agentMetadata)
		return err
	}
	kvsFilePath := file.Join(path, kvsFileName)
	log.Debugf("index path: %s and metadata: %s exists and successfully load metadata, now starting to load kvs data from %s", path, metadataPath, kvsFilePath)
	exist, fi, err = file.ExistsWithDetail(kvsFilePath)
	switch {
	case !exist, fi == nil, fi != nil && fi.Size() == 0, err != nil && errors.Is(err, fs.ErrNotExist):
		err = errors.Wrapf(errors.ErrIndexFileNotFound, "kvsdb file does not exists,\tpath: %s,\terr: %v", kvsFilePath, err)
		return err
	case err != nil && errors.Is(err, fs.ErrPermission):
		if fi != nil {
			err = errors.ErrFailedToOpenFile(err, kvsFilePath, 0, fi.Mode())
		}
		err = errors.Wrapf(err, "invalid permission for loading kvsdb file from %s", kvsFilePath)
		return err
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

	log.Debugf("index path: %s and metadata: %s and kvsdb file: %s exists and successfully load metadata, now starting to load full index and kvs data in concurrent", path, metadataPath, kvsFilePath)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	eg, _ := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() (err error) {
		if n.core != nil {
			n.core.Close()
			n.core = nil
		}
		n.core, err = core.Load(append(opts, core.WithIndexPath(path))...)
		if err != nil {
			err = errors.Wrapf(err, "failed to load ngt index from path: %s", path)
			return err
		}
		return nil
	}))

	eg.Go(safety.RecoverFunc(func() (err error) {
		err = n.loadKVS(kvsFilePath)
		if err != nil {
			err = errors.Wrapf(err, "failed to load kvsdb data from path: %s", kvsFilePath)
			return err
		}
		if n.kvs != nil && float64(agentMetadata.NGT.IndexCount/2) > float64(n.kvs.Len()) {
			return errors.ErrIndicesAreTooFewComparedToMetadata
		}
		return nil
	}))

	ech := make(chan error, 1)
	// NOTE: when it exceeds the timeout while loading,
	// it should exit this function and leave this goroutine running.
	n.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- safety.RecoverFunc(func() (err error) {
			err = eg.Wait()
			if err != nil {
				log.Error(err)
				return err
			}
			cancel()
			return nil
		})()
		return nil
	}))

	select {
	case err := <-ech:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Errorf("cannot load index backup data from %s within the timeout %s. the process is going to be killed.", path, timeout)
			err := metadata.Store(metadataPath,
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

func (n *ngt) initNGT(opts ...core.Option) (err error) {
	if n.kvs == nil {
		n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
	}
	if n.inMem {
		log.Debug("vald agent starts with in-memory mode")
		n.core, err = core.New(opts...)
		return err
	}

	ctx := context.Background()
	err = n.load(ctx, n.path, opts...)
	var current uint64
	if err != nil {
		if !n.enableCopyOnWrite {
			log.Debug("failed to load vald index from %s\t error: %v", n.path, err)
			if n.kvs == nil {
				n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
			} else if n.kvs.Len() > 0 {
				n.kvs.Close()
				n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
			}
			if n.core != nil {
				n.core.Close()
				n.core = nil
			}
			n.core, err = core.New(append(opts, core.WithIndexPath(n.path))...)
			return err
		}
		if errors.Is(err, errors.ErrIndicesAreTooFewComparedToMetadata) && n.kvs != nil {
			current = n.kvs.Len()
			log.Warnf(
				"load vald primary index success from %s\t error: %v\tbut index data are too few %d compared to metadata count now trying to load from old copied index data from %s and compare them",
				n.path,
				err,
				current,
				n.oldPath,
			)
		} else {
			log.Warnf("failed to load vald primary index from %s\t error: %v\ttrying to load from old copied index data from %s", n.path, err, n.oldPath)
		}
	} else {
		return nil
	}
	err = n.load(ctx, n.oldPath, opts...)
	if err == nil {
		if current != 0 && n.kvs.Len() < current {
			log.Warnf(
				"load vald secondary index success from %s\t error: %v\tbut index data are too few %d compared to primary data now trying to load from primary index data again from %s and start up with them",
				n.oldPath,
				err,
				n.kvs.Len(),
				n.oldPath,
			)
			err = n.load(ctx, n.path, opts...)
			if err == nil {
				return nil
			}
		} else {
			return nil
		}
	}
	log.Warnf("failed to load vald secondary index from %s and %s\t error: %v\ttrying to load from non-CoW index data from %s for backwards compatibility", n.path, n.oldPath, err, n.basePath)
	err = n.load(ctx, n.basePath, opts...)
	if err == nil {
		file.CopyDir(ctx, n.basePath, n.path)
		return nil
	}
	tpath := n.tmpPath.Load().(string)
	log.Warnf(
		"failed to load vald backwards index from %s and %s and %s\t error: %v\tvald agent couldn't find any index from agent volume in %s trying to start as new index from %s",
		n.path,
		n.oldPath,
		n.basePath,
		err,
		n.basePath,
		tpath,
	)
	if n.core != nil {
		n.core.Close()
		n.core = nil
	}
	n.core, err = core.New(append(opts, core.WithIndexPath(tpath))...)
	if err != nil {
		return err
	}
	if n.kvs == nil {
		n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
	} else if n.kvs.Len() > 0 {
		n.kvs.Close()
		n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
	}
	return nil
}

func (n *ngt) loadKVS(path string) (err error) {
	gob.Register(map[string]uint32{})

	var f *os.File
	f, err = file.Open(
		path,
		os.O_RDONLY|os.O_SYNC,
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

	m := make(map[string]uint32)
	err = gob.NewDecoder(f).Decode(&m)
	if err != nil {
		log.Errorf("error decoding kvsdb file,\terr: %v", err)
		return err
	}

	if n.kvs == nil {
		n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
	} else if n.kvs.Len() > 0 {
		n.kvs.Close()
		n.kvs = kvs.New(kvs.WithConcurrency(n.kvsdbConcurrency))
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
	n.removeInvalidIndex(ctx)
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

		if n.idelay > 0 {
			timer := time.NewTimer(n.idelay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			timer.Stop()
		}

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
	if n.IsFlushing() {
		return nil, errors.ErrFlushingIsInProgress
	}
	if n.IsIndexing() {
		return nil, errors.ErrCreateIndexingIsInProgress
	}
	sr, err := n.core.Search(vec, int(size), epsilon, radius)
	if err != nil {
		if n.IsIndexing() {
			return nil, errors.ErrCreateIndexingIsInProgress
		}
		log.Errorf("cgo error detected: ngt api returned error %v", err)
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

func (n *ngt) SearchByID(uuid string, size uint32, epsilon, radius float32) (vec []float32, dst []model.Distance, err error) {
	if n.IsFlushing() {
		return nil, nil, errors.ErrFlushingIsInProgress
	}
	if n.IsIndexing() {
		return nil, nil, errors.ErrCreateIndexingIsInProgress
	}
	vec, err = n.GetObject(uuid)
	if err != nil {
		return nil, nil, err
	}
	dst, err = n.Search(vec, size, epsilon, radius)
	if err != nil {
		return vec, nil, err
	}
	return vec, dst, nil
}

func (n *ngt) LinearSearch(vec []float32, size uint32) ([]model.Distance, error) {
	if n.IsFlushing() {
		return nil, errors.ErrFlushingIsInProgress
	}
	if n.IsIndexing() {
		return nil, errors.ErrCreateIndexingIsInProgress
	}
	sr, err := n.core.LinearSearch(vec, int(size))
	if err != nil {
		if n.IsIndexing() {
			return nil, errors.ErrCreateIndexingIsInProgress
		}
		log.Errorf("cgo error detected: ngt api returned error %v", err)
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

func (n *ngt) LinearSearchByID(uuid string, size uint32) (vec []float32, dst []model.Distance, err error) {
	if n.IsFlushing() {
		return nil, nil, errors.ErrFlushingIsInProgress
	}
	if n.IsIndexing() {
		return nil, nil, errors.ErrCreateIndexingIsInProgress
	}
	vec, err = n.GetObject(uuid)
	if err != nil {
		return nil, nil, err
	}
	dst, err = n.LinearSearch(vec, size)
	if err != nil {
		return vec, nil, err
	}
	return vec, dst, nil
}

func (n *ngt) Insert(uuid string, vec []float32) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	return n.insert(uuid, vec, time.Now().UnixNano(), true)
}

func (n *ngt) InsertWithTime(uuid string, vec []float32, t int64) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
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
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	return n.insertMultiple(vecs, time.Now().UnixNano(), true)
}

func (n *ngt) InsertMultipleWithTime(vecs map[string][]float32, t int64) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.insertMultiple(vecs, t, true)
}

func (n *ngt) insertMultiple(vecs map[string][]float32, now int64, validation bool) (err error) {
	for uuid, vec := range vecs {
		ierr := n.insert(uuid, vec, now, validation)
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
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	return n.update(uuid, vec, time.Now().UnixNano())
}

func (n *ngt) UpdateWithTime(uuid string, vec []float32, t int64) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.update(uuid, vec, t)
}

func (n *ngt) update(uuid string, vec []float32, t int64) (err error) {
	if err = n.readyForUpdate(uuid, vec); err != nil {
		return err
	}
	err = n.delete(uuid, t, true) // `true` is to return NotFound error with non-existent ID
	if err != nil {
		return err
	}
	t++
	return n.insert(uuid, vec, t, false)
}

func (n *ngt) UpdateMultiple(vecs map[string][]float32) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	return n.updateMultiple(vecs, time.Now().UnixNano())
}

func (n *ngt) UpdateMultipleWithTime(vecs map[string][]float32, t int64) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
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
	err = n.deleteMultiple(uuids, t, true) // `true` is to return NotFound error with non-existent ID
	if err != nil {
		return err
	}
	t++
	return n.insertMultiple(vecs, t, false)
}

func (n *ngt) Delete(uuid string) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	return n.delete(uuid, time.Now().UnixNano(), true)
}

func (n *ngt) DeleteWithTime(uuid string, t int64) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.delete(uuid, t, true)
}

func (n *ngt) delete(uuid string, t int64, validation bool) (err error) {
	if len(uuid) == 0 {
		err = errors.ErrUUIDNotFound(0)
		return err
	}
	if validation {
		_, ok := n.kvs.Get(uuid)
		if !ok && !n.vq.IVExists(uuid) {
			return errors.ErrObjectIDNotFound(uuid)
		}
	}
	return n.vq.PushDelete(uuid, t)
}

func (n *ngt) DeleteMultiple(uuids ...string) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	return n.deleteMultiple(uuids, time.Now().UnixNano(), true)
}

func (n *ngt) DeleteMultipleWithTime(uuids []string, t int64) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
	}
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return n.deleteMultiple(uuids, t, true)
}

func (n *ngt) deleteMultiple(uuids []string, now int64, validation bool) (err error) {
	for _, uuid := range uuids {
		ierr := n.delete(uuid, now, validation)
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

func (n *ngt) RegenerateIndex(ctx context.Context) (err error) {
	if n.IsFlushing() {
		return errors.ErrFlushingIsInProgress
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
	n.flushing.Store(true)
	n.indexing.Store(true)
	defer n.flushing.Store(false)
	defer n.indexing.Store(false)

	// delete kvs
	err = n.kvs.Close()
	if err != nil {
		log.Errorf("failed to flushing vector to ngt index in delete kvs. error: %v", err)
	}
	n.kvs = nil

	// delete ngt
	n.core.CloseWithoutSaveIndex()
	n.core = nil

	// delete vq
	n.vq = nil

	// gc
	runtime.GC()
	atomic.AddUint64(&n.nogce, 1)

	// delete file
	err = file.DeleteDir(ctx, n.path)
	if err != nil {
		// ここはエラーにすべき？
		log.Errorf("failed to flushing vector to ngt index in delete file. error: %v", err)
	}

	// delete cow
	if n.enableCopyOnWrite {
		err := file.DeleteDir(ctx, n.oldPath)
		if err != nil {
			// ここはエラーにすべき？
			log.Errorf("failed to flushing vector to ngt index in delete file. error: %v", err)
		}
	}

	// renew instance
	nn, err := newNGT(n.cfg, n.opts...)
	if err != nil {
		return err
	}
	// Regenerate with flags set
	nn.flushing.Store(true)
	nn.indexing.Store(true)
	n = nn

	return nil
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
		for n.IsIndexing() || n.IsSaving() || n.IsFlushing() {
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
			log.Errorf("failed to remove oid: %d from ngt index. error: %v", oid, err)
			n.fmu.Lock()
			n.fmap[uuid] = oid
			n.fmu.Unlock()
		}
		return true
	})
	log.Debug("create index delete phase finished")
	n.gc()
	log.Debug("create index insert phase started")
	var icnt uint32
	n.vq.RangePopInsert(ctx, now, func(uuid string, vector []float32) bool {
		oid, err := n.core.Insert(vector)
		if err != nil {
			log.Warnf("failed to insert vector uuid: %s vec: %v to ngt index. error: %v", uuid, vector, err)
			if !errors.Is(err, errors.ErrIncompatibleDimensionSize(len(vector), n.dim)) {
				oid, err = n.core.Insert(vector)
				if err != nil {
					log.Errorf("failed to retry insert vector uuid: %s vec: %v to ngt index. error: %v", uuid, vector, err)
					return true
				}
				n.kvs.Set(uuid, uint32(oid))
				atomic.AddUint32(&icnt, 1)
			}
		} else {
			n.kvs.Set(uuid, uint32(oid))
			atomic.AddUint32(&icnt, 1)
		}
		n.fmu.Lock()
		_, ok := n.fmap[uuid]
		if ok {
			delete(n.fmap, uuid)
		}
		n.fmu.Unlock()
		return true
	})
	if poolSize <= 0 {
		if n.poolSize > 0 && n.poolSize < atomic.LoadUint32(&icnt) {
			poolSize = n.poolSize
		} else {
			poolSize = atomic.LoadUint32(&icnt)
		}
	}
	log.Debug("create index insert phase finished")
	log.Debug("create graph and tree phase started")
	log.Debugf("pool size = %d", poolSize)
	err = n.core.CreateIndex(poolSize)
	if err != nil {
		log.Error("an error occurred on creating graph and tree phase:", err)
	} else {
		atomic.AddUint64(&n.nocie, 1)
	}
	log.Debug("create graph and tree phase finished")
	log.Info("create index operation finished")
	return err
}

func (n *ngt) removeInvalidIndex(ctx context.Context) {
	if n.kvs.Len() == 0 {
		return
	}
	var dcnt uint32
	n.kvs.Range(ctx, func(uuid string, oid uint32) bool {
		vec, err := n.core.GetVector(uint(oid))
		if err != nil || vec == nil || len(vec) != n.dim {
			log.Debugf("invalid index detected err: %v\tuuid: %s\toid: %d will remove", err, uuid, oid)
			n.kvs.Delete(uuid)
			n.fmu.Lock()
			err = n.core.Remove(uint(oid))
			n.fmap[uuid] = oid
			n.fmu.Unlock()
			atomic.AddUint32(&dcnt, 1)
			if err != nil {
				log.Debugf("invalid index remove operation returned error: %v", err)
			}
		}
		return true
	})
	if atomic.LoadUint32(&dcnt) <= 0 {
		return
	}
	var poolSize uint32
	if n.poolSize > 0 && n.poolSize < atomic.LoadUint32(&dcnt) {
		poolSize = n.poolSize
	} else {
		poolSize = atomic.LoadUint32(&dcnt)
	}
	n.cimu.Lock()
	defer n.cimu.Unlock()
	n.indexing.Store(true)
	defer n.indexing.Store(false)
	log.Debug("create graph and tree phase for removing invalid index started")
	err := n.core.CreateIndex(poolSize)
	if err != nil {
		log.Error("an error occurred on creating graph and tree phase:", err)
	} else {
		atomic.AddUint64(&n.nocie, 1)
	}
	log.Debug("create graph and tree phase for removing invalid index finished")
}

func (n *ngt) SaveIndex(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt/service/NGT.SaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if !n.inMem {
		return n.saveIndex(ctx)
	}
	return nil
}

func (n *ngt) saveIndex(ctx context.Context) (err error) {
	nocie := atomic.LoadUint64(&n.nocie)
	if atomic.LoadUint64(&n.lastNocie) == nocie {
		return
	}
	atomic.SwapUint64(&n.lastNocie, nocie)
	err = func() error {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		// wait for not indexing & not saving
		for n.IsIndexing() || n.IsSaving() || n.IsFlushing() {
			if n.IsFlushing() {
				return errors.ErrFlushingIsInProgress
			}
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

	log.Debug("cleanup invalid index started")
	n.removeInvalidIndex(ctx)
	log.Debug("cleanup invalid index finished")

	eg, ectx := errgroup.New(ctx)
	// we want to ensure the acutal kvs size between kvsdb and metadata,
	// so we create this counter to count the actual kvs size instead of using kvs.Len()
	var (
		kvsLen uint64
		path   string
	)

	if n.enableCopyOnWrite {
		path = n.tmpPath.Load().(string)
	} else {
		path = n.path
	}
	n.smu.Lock()
	defer n.smu.Unlock()
	eg.Go(safety.RecoverFunc(func() (err error) {
		if n.kvs.Len() > 0 && path != "" {
			m := make(map[string]uint32, n.Len())
			var mu sync.Mutex
			n.kvs.Range(ectx, func(key string, id uint32) bool {
				mu.Lock()
				m[key] = id
				mu.Unlock()
				atomic.AddUint64(&kvsLen, 1)
				return true
			})
			var f *os.File
			f, err = file.Open(
				file.Join(path, kvsFileName),
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
			m = make(map[string]uint32)
		}
		return nil
	}))

	eg.Go(safety.RecoverFunc(func() (err error) {
		n.fmu.Lock()
		fl := len(n.fmap)
		n.fmu.Unlock()
		if fl > 0 && path != "" {
			var f *os.File
			f, err = file.Open(
				file.Join(path, "invalid-"+kvsFileName),
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
			n.fmu.Lock()
			err = gob.NewEncoder(f).Encode(&n.fmap)
			n.fmu.Unlock()
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
		return n.core.SaveIndexWithPath(path)
	}))

	err = eg.Wait()
	if err != nil {
		return err
	}

	err = metadata.Store(
		file.Join(path, metadata.AgentMetadataFileName),
		&metadata.Metadata{
			IsInvalid: false,
			NGT: &metadata.NGT{
				IndexCount: kvsLen,
			},
		},
	)
	if err != nil {
		return err
	}

	return n.moveAndSwitchSavedData(ctx)
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

func (n *ngt) moveAndSwitchSavedData(ctx context.Context) (err error) {
	if !n.enableCopyOnWrite {
		return nil
	}
	n.cowmu.Lock()
	defer n.cowmu.Unlock()
	err = file.MoveDir(ctx, n.path, n.oldPath)
	if err != nil {
		log.Warnf("failed to backup backup data from %s to %s error: %v", n.path, n.oldPath, err)
	}
	path := n.tmpPath.Load().(string)
	err = file.MoveDir(ctx, path, n.path)
	if err != nil {
		log.Warnf("failed to move temporary index data from %s to %s error: %v, trying to rollback secondary backup data from %s to %s", path, n.path, n.oldPath, n.path, err)
		return file.MoveDir(ctx, n.oldPath, n.path)
	}
	defer log.Warnf("finished to copy index from %s => %s => %s", path, n.path, n.oldPath)
	return n.mktmp()
}

func (n *ngt) mktmp() (err error) {
	if !n.enableCopyOnWrite {
		return nil
	}
	path, err := file.MkdirTemp(file.Join(os.TempDir(), "vald"))
	if err != nil {
		log.Warnf("failed to create temporary index file path directory %s:\terr: %v", path, err)
		return err
	}
	n.tmpPath.Store(path)
	return nil
}

func (n *ngt) Exists(uuid string) (oid uint32, ok bool) {
	if n.IsFlushing() {
		log.Debugf("Exists\tuuid: %s's data will be delete soon\terror: %v",
			uuid, errors.ErrFlushingIsInProgress())
		return 0, false
	}
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
	// if error (GetObject cannot find vector) return error
	if err != nil {
		return err
	}
	// if vector length is not equal or if some difference exists let's try update
	if len(vec) != len(ovec) || conv.F32stos(vec) != conv.F32stos(ovec) {
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

func (n *ngt) IsFlushing() bool {
	i, ok := n.flushing.Load().(bool)
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
