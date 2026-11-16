// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"context"
	"encoding/gob"
	"io/fs"
	"math"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm"
	core "github.com/vdaas/vald/internal/core/algorithm/faiss"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/memstore"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

type (
	Faiss interface {
		Start(ctx context.Context) <-chan error
		Search(k, nprobe, nq uint32, xq []float32) (*payload.Search_Response, error)
		Insert(uuid string, vec []float32) (err error)
		InsertWithTime(uuid string, vec []float32, t int64) (err error)
		Update(uuid string, vec []float32) (err error)
		UpdateWithTime(uuid string, vec []float32, t int64) (err error)
		UpdateTimestamp(uuid string, ts int64, force bool) (err error)
		Delete(uuid string) (err error)
		DeleteWithTime(uuid string, t int64) (err error)
		Exists(uuid string) (uint32, bool)
		CreateIndex(ctx context.Context) (err error)
		SaveIndex(ctx context.Context) (err error)
		CreateAndSaveIndex(ctx context.Context) (err error)
		Train(nb int, vec []float32) (err error)
		IsIndexing() bool
		IsSaving() bool
		Len() uint64
		NumberOfCreateIndexExecution() uint64
		NumberOfProactiveGCExecution() uint64
		UUIDs(context.Context) (uuids []string)
		InsertVQueueBufferLen() uint64
		DeleteVQueueBufferLen() uint64
		GetDimensionSize() int
		GetTrainSize() int
		Close(ctx context.Context) error
	}

	faiss struct {
		tmpPath           atomic.Value
		eg                errgroup.Group
		kvs               kvs.BidiMap
		core              core.Faiss
		vq                vqueue.Queue
		saving            atomic.Value
		indexing          atomic.Value
		fmap              map[string]int64
		basePath          string
		oldPath           string
		path              string
		addIds            []int64
		addVecs           []float32
		nlist             int
		alen              int
		nocie             atomic.Uint64
		nogce             atomic.Uint64
		wfci              atomic.Uint64
		kvsdbConcurrency  int
		dim               int
		lastNocie         atomic.Uint64
		m                 int
		icnt              uint64
		dur               time.Duration
		sdur              time.Duration
		lim               time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		idelay            time.Duration
		trainSize         atomic.Int64
		addLen            atomic.Int64
		cimu              sync.Mutex
		smu               sync.Mutex
		fmu               sync.Mutex
		cowmu             sync.Mutex
		lmu               sync.Mutex
		isTrained         atomic.Bool
		enableCopyOnWrite bool
		dcd               bool
		enableProactiveGC bool
		inMem             bool
	}
)

const (
	kvsFileName          = "faiss-meta.kvsdb"
	kvsTimestampFileName = "faiss-timestamp.kvsdb"
	noTimeStampFile      = -1

	// faissIndexFileName is the file inside the index directory that holds the
	// faiss index; faiss write_index/read_index need a file, not a directory.
	faissIndexFileName = "faiss.index"

	oldIndexDirName    = "backup"
	originIndexDirName = "origin"

	// inflightIndexingPollInterval is how often waitForInflightIndexing re-checks the indexing state.
	inflightIndexingPollInterval = 100 * time.Millisecond

	// ref: https://github.com/facebookresearch/faiss/wiki/FAQ#can-i-ignore-warning-clustering-xxx-points-to-yyy-centroids
	// ref: https://github.com/facebookresearch/faiss/blob/main/faiss/Clustering.cpp#L38
	minPointsPerCentroid int = 39
)

func New(cfg *config.Faiss, opts ...Option) (Faiss, error) {
	var (
		f = &faiss{
			fmap:              make(map[string]int64),
			dim:               cfg.Dimension,
			nlist:             cfg.Nlist,
			m:                 cfg.M,
			enableProactiveGC: cfg.EnableProactiveGC,
			enableCopyOnWrite: cfg.EnableCopyOnWrite,
			kvsdbConcurrency:  cfg.KVSDB.Concurrency,
		}
		err error
	)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(f); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if len(f.path) == 0 {
		f.inMem = true
	}

	if !f.inMem {
		if err = moveAsideLegacyIndexFile(f.path); err != nil {
			log.Warn(err)
		}
	}

	if f.enableCopyOnWrite && !f.inMem && len(f.path) != 0 {
		sep := string(os.PathSeparator)
		f.path, err = filepath.Abs(strings.ReplaceAll(f.path, sep+sep, sep))
		if err != nil {
			log.Warn(err)
		}

		f.basePath = f.path
		f.oldPath = file.Join(f.basePath, oldIndexDirName)
		f.path = file.Join(f.basePath, originIndexDirName)
		err = file.MkdirAll(f.oldPath, fs.ModePerm)
		if err != nil {
			log.Warn(err)
		}
		err = file.MkdirAll(f.path, fs.ModePerm)
		if err != nil {
			log.Warn(err)
		}
		err = f.mktmp()
		if err != nil {
			return nil, err
		}
	}

	err = f.initFaiss(
		core.WithDimension(cfg.Dimension),
		core.WithNlist(cfg.Nlist),
		core.WithM(cfg.M),
		core.WithNbitsPerIdx(cfg.NbitsPerIdx),
		core.WithMethodType(cfg.MethodType),
		core.WithMetricType(cfg.MetricType),
	)
	if err != nil {
		return nil, err
	}

	if f.dur == 0 || f.alen == 0 {
		f.dcd = true
	}

	if f.vq == nil {
		f.vq, err = vqueue.New()
		if err != nil {
			return nil, err
		}
	}

	f.indexing.Store(false)
	f.saving.Store(false)

	return f, nil
}

// moveAsideLegacyIndexFile renames a regular file at the index path to a
// backup name next to it. Older agents wrote the faiss index directly to the
// index path, but the index path is now a directory (see faissIndexFileName)
// and such a file makes every MkdirAll, and so every save, fail. The file can
// never be loaded as an index directory, so it is kept, not deleted.
func moveAsideLegacyIndexFile(path string) error {
	exist, fi, err := file.ExistsWithDetail(path)
	if !exist {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "failed to inspect index path %s", path)
	}
	if fi == nil || fi.IsDir() {
		return nil
	}
	base := path + ".legacy-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	dst := base
	for i := 1; file.Exists(dst); i++ {
		dst = base + "-" + strconv.Itoa(i)
	}
	if err = os.Rename(path, dst); err != nil {
		return errors.Wrapf(err, "failed to move legacy index file %s aside to %s", path, dst)
	}
	log.Warnf("moved legacy index file %s aside to %s; the index path must be a directory", path, dst)
	return nil
}

func (f *faiss) initFaiss(opts ...core.Option) error {
	var err error

	if f.kvs == nil {
		f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
	}

	if f.inMem {
		log.Debug("vald agent starts with in-memory mode")
		f.core, err = core.New(opts...)
		return err
	}

	ctx := context.TODO()
	err = f.load(ctx, f.path, opts...)
	// A timed-out load leaves its goroutines running and still writing f.core,
	// f.kvs and f.fmap, so no fallback may run alongside it; the process is
	// expected to exit on this error.
	if errors.Is(err, errors.ErrIndexLoadTimeout) {
		return err
	}
	var current uint64
	if err != nil {
		if !f.enableCopyOnWrite {
			log.Debugf("failed to load vald index from %s\t error: %v", f.path, err)
			if f.kvs == nil {
				f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
			} else if f.kvs.Len() > 0 {
				f.kvs.Close()
				f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
			}

			if f.core != nil {
				f.core.Close()
				f.core = nil
			}
			f.resetTrainState()
			f.core, err = core.New(append(opts, core.WithIndexPath(file.Join(f.path, faissIndexFileName)))...)
			return err
		}

		if errors.Is(err, errors.ErrIndicesAreTooFewComparedToMetadata) && f.kvs != nil {
			current = f.kvs.Len()
			log.Warnf(
				"load vald primary index success from %s\t error: %v\tbut index data are too few %d compared to metadata count now trying to load from old copied index data from %s and compare them",
				f.path,
				err,
				current,
				f.oldPath,
			)
		} else {
			log.Warnf("failed to load vald primary index from %s\t error: %v\ttrying to load from old copied index data from %s", f.path, err, f.oldPath)
		}
	} else {
		return nil
	}

	err = f.load(ctx, f.oldPath, opts...)
	if errors.Is(err, errors.ErrIndexLoadTimeout) {
		return err
	}
	if err == nil {
		if current != 0 && f.kvs.Len() < current {
			log.Warnf(
				"load vald secondary index success from %s\t error: %v\tbut index data are too few %d compared to primary data now trying to load from primary index data again from %s and start up with them",
				f.oldPath,
				err,
				f.kvs.Len(),
				f.oldPath,
			)

			err = f.load(ctx, f.path, opts...)
			// The primary has the same too-few kvsdb as on the first load, so
			// accept that error as long as its index was loaded; waitLoad marks
			// only an error-free load as trained.
			if errors.Is(err, errors.ErrIndicesAreTooFewComparedToMetadata) && f.core != nil {
				f.isTrained.Store(true)
				return nil
			}
			if err == nil || errors.Is(err, errors.ErrIndexLoadTimeout) {
				return err
			}
		} else {
			return nil
		}
	}

	log.Warnf("failed to load vald secondary index from %s and %s\t error: %v\ttrying to load from non-CoW index data from %s for backwards compatibility", f.path, f.oldPath, err, f.basePath)
	err = f.load(ctx, f.basePath, opts...)
	if errors.Is(err, errors.ErrIndexLoadTimeout) {
		return err
	}
	if err == nil {
		if cerr := file.CopyDir(ctx, f.basePath, f.path); cerr != nil {
			log.Warnf("failed to copy backwards-compatible index from %s to %s: %v", f.basePath, f.path, cerr)
		}
		return nil
	}

	if current != 0 {
		// No other index is usable, so start up with the too-few primary rather
		// than an empty index; its reload reports the same too-few error.
		log.Warnf("failed to load vald backwards index from %s\t error: %v\ttrying to load from primary index data again from %s and start up with them", f.basePath, err, f.path)
		err = f.load(ctx, f.path, opts...)
		if errors.Is(err, errors.ErrIndicesAreTooFewComparedToMetadata) && f.core != nil {
			f.isTrained.Store(true)
			return nil
		}
		if err == nil || errors.Is(err, errors.ErrIndexLoadTimeout) {
			return err
		}
	}

	tpath := f.tmpPath.Load().(string)
	log.Warnf(
		"failed to load vald backwards index from %s and %s and %s\t error: %v\tvald agent couldn't find any index from agent volume in %s trying to start as new index from %s",
		f.path,
		f.oldPath,
		f.basePath,
		err,
		f.basePath,
		tpath,
	)

	if f.core != nil {
		f.core.Close()
		f.core = nil
	}
	f.resetTrainState()
	f.core, err = core.New(append(opts, core.WithIndexPath(file.Join(tpath, faissIndexFileName)))...)
	if err != nil {
		return err
	}

	if f.kvs == nil {
		f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
	} else if f.kvs.Len() > 0 {
		f.kvs.Close()
		f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
	}

	return nil
}

// resetTrainState clears the state restored by a successful load, so that a
// fresh core created by an initFaiss fallback is trained again by CreateIndex
// and its ids start from zero together with the fresh kvs.
func (f *faiss) resetTrainState() {
	f.isTrained.Store(false)
	atomic.StoreUint64(&f.icnt, 0)
}

func (f *faiss) load(ctx context.Context, path string, opts ...core.Option) error {
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
	if err != nil && errors.Is(err, fs.ErrNotExist) || agentMetadata == nil || agentMetadata.Faiss == nil || agentMetadata.Faiss.IndexCount == 0 {
		err = errors.Wrapf(err, "cannot read metadata from path: %s\tmetadata: %v", path, agentMetadata)
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

	kvsTimestampFilePath := file.Join(path, kvsTimestampFileName)
	log.Debugf("now starting to load kvs timestamp data from %s", kvsTimestampFilePath)
	exist, fi, err = file.ExistsWithDetail(kvsTimestampFilePath)
	switch {
	case !exist, fi == nil, fi != nil && fi.Size() == 0, err != nil && errors.Is(err, fs.ErrNotExist):
		log.Warnf("timestamp kvsdb file does not exists,\tpath: %s,\terr: %v", kvsTimestampFilePath, err)
	case err != nil && errors.Is(err, fs.ErrPermission):
		if fi != nil {
			err = errors.ErrFailedToOpenFile(err, kvsTimestampFilePath, 0, fi.Mode())
		}
		log.Warnf("invalid permission for loading timestamp kvsdb file from %s", kvsTimestampFilePath)
	}

	indexFilePath := file.Join(path, faissIndexFileName)
	exist, fi, err = file.ExistsWithDetail(indexFilePath)
	switch {
	case !exist, fi == nil, fi.Size() == 0, err != nil && errors.Is(err, fs.ErrNotExist):
		err = errors.Wrapf(errors.ErrIndexFileNotFound, "faiss index file does not exists,\tpath: %s,\terr: %v", indexFilePath, err)
		return err
	case err != nil && errors.Is(err, fs.ErrPermission):
		// fi is non-nil here: the case above already returned for fi == nil.
		err = errors.ErrFailedToOpenFile(err, indexFilePath, 0, fi.Mode())
		err = errors.Wrapf(err, "invalid permission for loading faiss index file from %s", indexFilePath)
		return err
	}

	var timeout time.Duration
	if agentMetadata != nil && agentMetadata.Faiss != nil {
		log.Debugf("the backup index size is %d. starting to load...", agentMetadata.Faiss.IndexCount)
		timeout = time.Duration(
			math.Min(
				math.Max(
					float64(agentMetadata.Faiss.IndexCount)*float64(f.litFactor),
					float64(f.minLit),
				),
				float64(f.maxLit),
			),
		)
	} else {
		log.Debugf("cannot inspect the backup index size. starting to load default value.")
		timeout = time.Duration(math.Min(float64(f.minLit), float64(f.maxLit)))
	}

	log.Debugf(
		"index path: %s and metadata: %s and kvsdb file: %s and timestamp kvsdb file: %s exists and successfully load metadata, now starting to load full index and kvs data in concurrent",
		path,
		metadataPath,
		kvsFilePath,
		kvsTimestampFilePath,
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	eg, _ := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() (err error) {
		// Serialize the Close/reassign of f.core. On a load timeout load()
		// returns but leaves this goroutine running (the process is expected to
		// be killed) while the caller falls back to another load(); without
		// this lock the abandoned goroutine and the new one would concurrently
		// Close/reassign the same cgo handle — a data race and a double-Close
		// (use-after-free) of the FAISS C++ index.
		f.lmu.Lock()
		defer f.lmu.Unlock()
		if f.core != nil {
			f.core.Close()
			f.core = nil
		}
		f.core, err = core.Load(append(opts, core.WithIndexPath(indexFilePath))...)
		if err != nil {
			err = errors.Wrapf(err, "failed to load faiss index from path: %s", indexFilePath)
			return err
		}
		return nil
	}))

	eg.Go(safety.RecoverFunc(func() (err error) {
		err = f.loadKVS(ctx, path, timeout)
		if err != nil {
			err = errors.Wrapf(err, "failed to load kvsdb data from path: %s, %s", kvsFilePath, kvsTimestampFilePath)
			return err
		}
		if f.kvs != nil && float64(agentMetadata.Faiss.IndexCount/2) > float64(f.kvs.Len()) {
			return errors.ErrIndicesAreTooFewComparedToMetadata
		}
		return nil
	}))

	return f.waitLoad(ctx, eg.Wait, path, metadataPath, timeout)
}

// waitLoad waits for wait (the concurrent index and kvsdb load) to finish
// within ctx and marks the loaded index as trained on success.
// NOTE: when it exceeds the timeout while loading,
// it returns ErrIndexLoadTimeout and leaves the load goroutine running.
func (f *faiss) waitLoad(
	ctx context.Context, wait func() error, path, metadataPath string, timeout time.Duration,
) error {
	ech := make(chan error, 1)
	// wait must not cancel ctx before sending its result: a parked select would
	// then be woken by ctx.Done() and skip the success path below.
	f.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- safety.RecoverFunc(func() (err error) {
			err = wait()
			if err != nil {
				log.Error(err)
				return err
			}
			return nil
		})()
		return nil
	}))

	select {
	case err := <-ech:
		if err == nil {
			// Only a trained index is ever saved (saveIndex skips otherwise), so a
			// loaded one is trained; without this the restored index is never saved
			// again and CreateIndex would buffer new vectors for a second Train.
			f.isTrained.Store(true)
		}
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Errorf("cannot load index backup data from %s within the timeout %s. the process is going to be killed.", path, timeout)
			err := metadata.Store(metadataPath,
				&metadata.Metadata{
					IsInvalid: true,
					Faiss: &metadata.Faiss{
						IndexCount: 0,
					},
				},
			)
			if err != nil {
				// keep the timeout identity: initFaiss must not start another
				// load while the abandoned one is still running.
				return errors.Join(errors.ErrIndexLoadTimeout, err)
			}
			return errors.ErrIndexLoadTimeout
		}
	}

	return nil
}

func (f *faiss) loadKVS(ctx context.Context, path string, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	eg, _ := errgroup.New(ctx)

	m := make(map[string]uint32)
	mt := make(map[string]int64)

	eg.Go(safety.RecoverFunc(func() (err error) {
		gob.Register(map[string]uint32{})
		var fi *os.File
		fi, err = file.Open(
			file.Join(path, kvsFileName),
			os.O_RDONLY|os.O_SYNC,
			fs.ModePerm,
		)
		if err != nil {
			return err
		}
		defer func() {
			if fi != nil {
				derr := fi.Close()
				if derr != nil {
					err = errors.Wrap(err, derr.Error())
				}
			}
		}()
		err = gob.NewDecoder(fi).Decode(&m)
		if err != nil {
			log.Errorf("error decoding kvsdb file,\terr: %v", err)
			return err
		}
		return nil
	}))

	eg.Go(safety.RecoverFunc(func() (err error) {
		gob.Register(map[string]int64{})
		var ft *os.File
		ft, err = file.Open(
			file.Join(path, kvsTimestampFileName),
			os.O_RDONLY|os.O_SYNC,
			fs.ModePerm,
		)
		if err != nil {
			log.Warnf("error opening timestamp kvsdb file,\terr: %v", err)
		}
		defer func() {
			if ft != nil {
				derr := ft.Close()
				if derr != nil {
					err = errors.Wrap(err, derr.Error())
				}
			}
		}()
		err = gob.NewDecoder(ft).Decode(&mt)
		if err != nil {
			log.Warnf("error decoding timestamp kvsdb file,\terr: %v", err)
		}
		return nil
	}))

	err = eg.Wait()
	if err != nil {
		return err
	}

	if f.kvs == nil {
		f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
	} else if f.kvs.Len() > 0 {
		f.kvs.Close()
		f.kvs = kvs.New(kvs.WithConcurrency(f.kvsdbConcurrency))
	}
	// Resume id allocation after the largest loaded id; restarting from 0 would
	// hand out ids that already label vectors in the loaded index.
	var next uint64
	for k, id := range m {
		if uint64(id) >= next {
			next = uint64(id) + 1
		}
		if ts, ok := mt[k]; ok {
			f.kvs.Set(k, id, ts)
		} else {
			// NOTE: SaveIndex do not write ngt-timestamp.kvsdb with timestamp 0.
			f.kvs.Set(k, id, 0)
			f.fmap[k] = int64(id)
		}
	}
	for k := range mt {
		if _, ok := m[k]; !ok {
			f.fmap[k] = noTimeStampFile
		}
	}
	atomic.StoreUint64(&f.icnt, next)

	return nil
}

func (f *faiss) mktmp() error {
	if !f.enableCopyOnWrite {
		return nil
	}

	path, err := file.MkdirTemp(file.Join(os.TempDir(), "vald"))
	if err != nil {
		log.Warnf("failed to create temporary index file path directory %s:\terr: %v", path, err)
		return err
	}

	f.tmpPath.Store(path)

	return nil
}

func (f *faiss) Start(ctx context.Context) <-chan error {
	if f.dcd {
		return nil
	}

	ech := make(chan error, 2)
	f.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if f.dur <= 0 {
			f.dur = math.MaxInt64
		}
		if f.sdur <= 0 {
			f.sdur = math.MaxInt64
		}
		if f.lim <= 0 {
			f.lim = math.MaxInt64
		}

		if f.idelay > 0 {
			timer := time.NewTimer(f.idelay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			timer.Stop()
		}

		tick := time.NewTicker(f.dur)
		sTick := time.NewTicker(f.sdur)
		limit := time.NewTicker(f.lim)
		defer tick.Stop()
		defer sTick.Stop()
		defer limit.Stop()
		for {
			err = nil
			select {
			case <-ctx.Done():
				// No final CreateIndex here: ctx is already done, so CreateIndex would refuse to
				// touch the vqueue. Close (PostStop) runs the final CreateIndex/SaveIndex with a
				// live ctx.
				return ctx.Err()
			case <-tick.C:
				if f.vq.IVQLen() >= f.alen {
					err = f.CreateIndex(ctx)
				}
			case <-limit.C:
				err = f.CreateAndSaveIndex(ctx)
			case <-sTick.C:
				err = f.SaveIndex(ctx)
			}
			if err != nil && !errors.Is(err, errors.ErrUncommittedIndexNotFound) {
				ech <- err
				runtime.Gosched()
				err = nil
			}
		}
	}))

	return ech
}

func (f *faiss) Train(nb int, xb []float32) error {
	err := f.core.Train(nb, xb)
	if err != nil {
		log.Errorf("failed to faiss train: %v", err)
		return err
	}

	return nil
}

func (f *faiss) Insert(uuid string, vec []float32) error {
	return f.insert(uuid, vec, time.Now().UnixNano(), true)
}

func (f *faiss) InsertWithTime(uuid string, vec []float32, t int64) error {
	if t <= 0 {
		t = time.Now().UnixNano()
	}

	return f.insert(uuid, vec, t, true)
}

func (f *faiss) insert(uuid string, xb []float32, t int64, validation bool) error {
	if len(uuid) == 0 {
		err := errors.ErrUUIDNotFound(0)
		return err
	}

	if validation {
		_, ok := f.Exists(uuid)
		if ok {
			return errors.ErrUUIDAlreadyExists(uuid)
		}
	}

	return f.vq.PushInsert(uuid, xb, t)
}

func (f *faiss) Update(uuid string, vec []float32) error {
	return f.update(uuid, vec, time.Now().UnixNano())
}

func (f *faiss) UpdateWithTime(uuid string, vec []float32, t int64) error {
	if t <= 0 {
		t = time.Now().UnixNano()
	}
	return f.update(uuid, vec, t)
}

func (f *faiss) update(uuid string, vec []float32, t int64) (err error) {
	// the insert runs at t+1, which wraps negative for MaxInt64 and would be
	// dropped after the delete; reject before anything is deleted.
	if t == math.MaxInt64 {
		return errors.ErrInvalidTimestamp(t)
	}
	if err = f.readyForUpdate(uuid, vec); err != nil {
		return err
	}

	err = f.delete(uuid, t, true) // `true` is to return NotFound error with non-existent ID
	if err != nil {
		return err
	}

	t++
	return f.insert(uuid, vec, t, false)
}

func (f *faiss) UpdateTimestamp(uuid string, ts int64, force bool) (err error) {
	return memstore.UpdateTimestamp(f.kvs, f.vq, uuid, ts, force)
}

func (f *faiss) readyForUpdate(uuid string, vec []float32) (err error) {
	if len(uuid) == 0 {
		return errors.ErrUUIDNotFound(0)
	}

	if len(vec) != f.GetDimensionSize() {
		return errors.ErrInvalidDimensionSize(len(vec), f.GetDimensionSize())
	}

	// not impl GetObject()

	return nil
}

func (f *faiss) CreateIndex(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-faiss/service/Faiss.CreateIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	// addVecs is guarded by cimu; addLen mirrors its vector count for this lock-free check.
	ic := f.vq.IVQLen() + f.vq.DVQLen() + int(f.addLen.Load())
	if ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}

	wf := f.wfci.Add(1)
	if wf > 1 {
		f.wfci.Add(^uint64(0))
		log.Debugf("concurrent create index waiting detected this request will be ignored, concurrent: %d", wf)
		return nil
	}

	err := func() error {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		// wait for not indexing & not saving
		for f.IsIndexing() || f.IsSaving() {
			runtime.Gosched()
			select {
			case <-ctx.Done():
				f.wfci.Add(^uint64(0))
				return ctx.Err()
			case <-ticker.C:
			}
		}
		f.wfci.Add(^uint64(0))
		return nil
	}()
	if err != nil {
		return err
	}

	f.cimu.Lock()
	defer f.cimu.Unlock()
	f.indexing.Store(true)
	defer f.indexing.Store(false)
	defer f.gc()
	now := time.Now().UnixNano()
	ic = f.vq.IVQLen() + f.vq.DVQLen() + (len(f.addVecs) / f.dim)
	if ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}
	// A done ctx makes the vqueue ranges stop after their first entry, so do not start a run
	// that could apply only part of the queue; the entries stay queued for the next run.
	if err := ctx.Err(); err != nil {
		return err
	}

	log.Infof("create index operation started, uncommitted indexes = %d", ic)
	log.Debug("create index delete phase started")
	var dcnt uint64
	f.vq.RangePopDelete(ctx, now, func(uuid string, dts int64) bool {
		log.Debugf("start delete operation for kvsdb id: %s", uuid)
		// Newest wins: a delete older than the indexed object (e.g. a stale Remove, or the
		// delete half of a stale Update whose insert is skipped below) must not remove it.
		if _, kts, ok := f.kvs.Get(uuid); ok && kts > dts {
			log.Debugf("skipped delete operation for kvsdb id: %s, delete timestamp %d is older than the indexed timestamp %d", uuid, dts, kts)
			return false
		}
		oid, ok := f.kvs.Delete(uuid)
		if !ok {
			log.Warn(errors.ErrObjectIDNotFound(uuid))
			return true
		}
		dcnt++
		if f.removeBuffered(int64(oid)) {
			log.Debugf("removed untrained buffered vector and kvsdb id: %s, oid: %d", uuid, oid)
			return true
		}
		log.Debugf("start remove operation for faiss index id: %s, oid: %d", uuid, oid)
		ntotal, err := f.core.Remove(1, []int64{int64(oid)})
		if err != nil {
			log.Errorf("failed to remove oid: %d from faiss index. error: %v", oid, err)
			f.fmu.Lock()
			f.fmap[uuid] = int64(oid)
			f.fmu.Unlock()
		}
		log.Debugf("removed from faiss index and kvsdb id: %s, oid: %d, index size: %d", uuid, oid, ntotal)
		return true
	})
	log.Debug("create index delete phase finished")

	f.gc()

	// RangePopDelete stops early only when ctx is done. Applying inserts after a partial delete
	// phase could apply the insert half of an Update whose delete is still queued: kvs.Set would
	// repoint the uuid to a new id and leave the old id orphaned in the index or the add buffer.
	// Keep the inserts queued, but still finish the run (train/add of what is already buffered)
	// when deletes were applied so that it is counted (nocie) and the deletes get saved.
	skipped := ctx.Err()
	if skipped != nil && dcnt == 0 {
		return skipped
	}
	if skipped != nil {
		log.Warnf("create index insert phase skipped after %d deletes because the context is done: %v", dcnt, skipped)
	} else {
		log.Debug("create index insert phase started")
		f.vq.RangePopInsert(ctx, now, func(uuid string, vector []float32, timestamp int64) bool {
			// Newest wins: an insert not newer than the indexed object (e.g. the insert half
			// of a stale Update whose delete was skipped above) is dropped.
			ooid, kts, exists := f.kvs.Get(uuid)
			if exists && kts >= timestamp {
				log.Debugf("skipped insert operation for kvsdb id: %s, insert timestamp %d is not newer than the indexed timestamp %d", uuid, timestamp, kts)
				return true
			}
			// A NaN/Inf vector would make faiss Train/Add fail and stall every later
			// indexing run, so drop it here like NGT does.
			if algorithm.HasNaNOrInf(vector) {
				log.Warnf("skipped insert operation for kvsdb id: %s: %v", uuid, errors.ErrNaNOrInfVector)
				return true
			}
			log.Debugf("start stack operation for faiss index id: %s, icnt: %d", uuid, f.icnt)
			f.addVecs = append(f.addVecs, vector...)
			f.addIds = append(f.addIds, int64(f.icnt)) //nolint:gosec // icnt counts inserts one at a time and cannot reach 2^63
			f.addLen.Add(1)

			log.Debugf("start insert operation for kvsdb id: %s, icnt: %d", uuid, f.icnt)
			f.kvs.Set(uuid, uint32(f.icnt), timestamp) //nolint:gosec // the kvs oid is uint32 by contract; this narrowing predates and is unchanged by the newest-wins rework
			atomic.AddUint64(&f.icnt, 1)
			// The uuid still pointed to an older object whose delete was skipped or never
			// queued; Set repointed it, so remove the old object instead of orphaning it.
			if exists && !f.removeBuffered(int64(ooid)) {
				if _, rerr := f.core.Remove(1, []int64{int64(ooid)}); rerr != nil {
					log.Errorf("failed to remove replaced oid: %d of uuid: %s from faiss index. error: %v", ooid, uuid, rerr)
				}
			}

			f.fmu.Lock()
			_, ok := f.fmap[uuid]
			if ok {
				delete(f.fmap, uuid)
			}
			f.fmu.Unlock()
			log.Debugf("finished to insert index and kvsdb id: %s, icnt: %d", uuid, f.icnt)
			return true
		})
	}

	var max int
	if f.nlist > int(math.Pow(2, float64(f.m))) {
		max = f.nlist
	} else {
		max = int(math.Pow(2, float64(f.m)))
	}
	if !f.isTrained.Load() && len(f.addVecs)/f.dim >= max*minPointsPerCentroid {
		log.Debug("faiss train phase started")
		log.Debugf("max * minPointsPerCentroid: %d", max*minPointsPerCentroid)
		err := f.core.Train(len(f.addVecs)/f.dim, f.addVecs)
		if err != nil {
			log.Errorf("failed to faiss train: %v", err)
			return err
		}
		f.isTrained.Store(true)
		f.trainSize.Store(int64(len(f.addVecs) / f.dim))
		log.Debug("faiss train phase finished")
	}
	if f.isTrained.Load() && len(f.addVecs) > 0 {
		log.Debug("faiss add phase started")
		ntotal, err := f.core.Add(len(f.addVecs)/f.dim, f.addVecs, f.addIds)
		if err != nil {
			log.Errorf("failed to faiss add: %v", err)
			return err
		}
		f.addVecs = nil
		f.addIds = nil
		f.addLen.Store(0)
		log.Debugf("is trained: %v, index size: %d", f.isTrained.Load(), ntotal)
		log.Debug("faiss add phase finished")
	}
	log.Debug("create index insert phase finished")

	f.nocie.Add(1)
	log.Info("create index operation finished")

	return skipped
}

// removeBuffered drops oid and its vector from the not-yet-added buffer
// (addIds/addVecs) by swap-remove and reports whether it was found. Otherwise a
// deleted vector would later be added to the index as a label with no kvs entry.
// It must be called with f.cimu held, as every other addIds/addVecs access is.
func (f *faiss) removeBuffered(oid int64) bool {
	if f.dim <= 0 || len(f.addVecs) != len(f.addIds)*f.dim {
		return false
	}
	for i, id := range f.addIds {
		if id != oid {
			continue
		}
		last := len(f.addIds) - 1
		f.addIds[i] = f.addIds[last]
		f.addIds = f.addIds[:last]
		copy(f.addVecs[i*f.dim:(i+1)*f.dim], f.addVecs[last*f.dim:])
		f.addVecs = f.addVecs[:last*f.dim]
		f.addLen.Add(-1)
		return true
	}
	return false
}

func (f *faiss) SaveIndex(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-faiss/service/Faiss.SaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if !f.inMem {
		return f.saveIndex(ctx)
	}

	return nil
}

func (f *faiss) saveIndex(ctx context.Context) (err error) {
	nocie := f.nocie.Load()
	if f.lastNocie.Load() == nocie || !f.isTrained.Load() {
		return nil
	}
	// Mark this generation as saved up front so concurrent saves are deduplicated,
	// but roll it back if the save fails; otherwise every later SaveIndex would
	// skip until the next CreateIndex and the new index would never be persisted.
	prev := f.lastNocie.Swap(nocie)
	defer func() {
		if err != nil {
			f.lastNocie.CompareAndSwap(nocie, prev)
		}
	}()

	err = func() error {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		// wait for not indexing & not saving
		for f.IsIndexing() || f.IsSaving() {
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

	f.saving.Store(true)
	defer f.gc()
	defer f.saving.Store(false)

	// no cleanup invalid index

	eg, ectx := errgroup.New(ctx)
	// we want to ensure the actual kvs size between kvsdb and metadata,
	// so we create this counter to count the actual kvs size instead of using kvs.Len()
	var (
		kvsLen uint64
		path   string
	)

	if f.enableCopyOnWrite {
		path = f.tmpPath.Load().(string)
	} else {
		path = f.path
	}

	f.smu.Lock()
	defer f.smu.Unlock()

	// The kvsdb writers below create the directory only as a side effect, while
	// the index writer does not create it at all, so create it before they race.
	if path != "" {
		err = file.MkdirAll(path, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	if f.kvs.Len() > 0 && path != "" {
		eg.Go(safety.RecoverFunc(func() error {
			return f.saveKVS(ectx, path, &kvsLen)
		}))
	}

	eg.Go(safety.RecoverFunc(func() (err error) {
		f.fmu.Lock()
		fl := len(f.fmap)
		f.fmu.Unlock()

		if fl > 0 && path != "" {
			var fi *os.File
			fi, err = file.Open(
				file.Join(path, "invalid-"+kvsFileName),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				fs.ModePerm,
			)
			if err != nil {
				return err
			}
			defer func() {
				if fi != nil {
					derr := fi.Close()
					if derr != nil {
						err = errors.Join(err, derr)
					}
				}
			}()

			gob.Register(map[string]int64{})
			f.fmu.Lock()
			err = gob.NewEncoder(fi).Encode(&f.fmap)
			f.fmu.Unlock()
			if err != nil {
				return err
			}
			err = fi.Sync()
			if err != nil {
				return err
			}
		}

		return nil
	}))

	eg.Go(safety.RecoverFunc(func() error {
		return f.core.SaveIndexWithPath(file.Join(path, faissIndexFileName))
	}))

	err = eg.Wait()
	if err != nil {
		return err
	}

	err = metadata.Store(
		file.Join(path, metadata.AgentMetadataFileName),
		&metadata.Metadata{
			IsInvalid: false,
			Faiss: &metadata.Faiss{
				IndexCount: kvsLen,
			},
		},
	)
	if err != nil {
		return err
	}

	return f.moveAndSwitchSavedData(ctx)
}

// saveKVS writes the uuid->id kvsdb and the uuid->timestamp kvsdb files under
// path from one snapshot of f.kvs, adding the snapshot size to kvsLen.
func (f *faiss) saveKVS(ctx context.Context, path string, kvsLen *uint64) (err error) {
	m := make(map[string]uint32, f.Len())
	mt := make(map[string]int64, f.Len())
	defer func() {
		m = nil
		mt = nil
	}()
	var mu sync.Mutex

	// Take the snapshot with a non-cancelable context: kvs.Range silently
	// stops early on cancellation, and the kvsdb/timestamp files below are
	// truncated and rewritten from this snapshot (and kvsLen is recorded as
	// the metadata IndexCount), so a partial range would persist a truncated
	// uuid->id map. The file writes below and the index save in saveIndex do
	// not observe ctx either, so the snapshot now completes like they do.
	f.kvs.Range(context.WithoutCancel(ctx), func(key string, id uint32, ts int64) bool {
		mu.Lock()
		m[key] = id
		mt[key] = ts
		mu.Unlock()
		atomic.AddUint64(kvsLen, 1)
		return true
	})

	var fi *os.File
	fi, err = file.Open(
		file.Join(path, kvsFileName),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		fs.ModePerm,
	)
	if err != nil {
		return err
	}
	defer func() {
		if fi != nil {
			derr := fi.Close()
			if derr != nil {
				err = errors.Wrap(err, derr.Error())
			}
		}
	}()

	gob.Register(map[string]uint32{})
	err = gob.NewEncoder(fi).Encode(&m)
	if err != nil {
		return err
	}

	err = fi.Sync()
	if err != nil {
		return err
	}

	m = make(map[string]uint32)

	var ft *os.File
	ft, err = file.Open(
		file.Join(path, kvsTimestampFileName),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		fs.ModePerm,
	)
	if err != nil {
		return err
	}
	defer func() {
		if ft != nil {
			derr := ft.Close()
			if derr != nil {
				err = errors.Wrap(err, derr.Error())
			}
		}
	}()

	gob.Register(map[string]int64{})
	err = gob.NewEncoder(ft).Encode(&mt)
	if err != nil {
		return err
	}

	err = ft.Sync()
	if err != nil {
		return err
	}

	mt = make(map[string]int64)

	return nil
}

func (f *faiss) moveAndSwitchSavedData(ctx context.Context) error {
	if !f.enableCopyOnWrite {
		return nil
	}

	var err error
	f.cowmu.Lock()
	defer f.cowmu.Unlock()

	err = file.MoveDir(ctx, f.path, f.oldPath)
	if err != nil {
		log.Warnf("failed to backup data from %s to %s error: %v", f.path, f.oldPath, err)
	}

	path := f.tmpPath.Load().(string)
	err = file.MoveDir(ctx, path, f.path)
	if err != nil {
		log.Warnf("failed to move temporary index data from %s to %s error: %v, trying to rollback secondary backup data from %s to %s", path, f.path, err, f.oldPath, f.path)
		return file.MoveDir(ctx, f.oldPath, f.path)
	}
	defer log.Warnf("finished to copy index from %s => %s => %s", path, f.path, f.oldPath)

	return f.mktmp()
}

func (f *faiss) CreateAndSaveIndex(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-faiss/service/Faiss.CreateAndSaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := f.CreateIndex(ctx)
	if errors.IsNot(err, errors.ErrUncommittedIndexNotFound, context.Canceled, context.DeadlineExceeded) {
		return err
	}

	return f.SaveIndex(ctx)
}

func (f *faiss) Search(
	k, nprobe, nq uint32, xq []float32,
) (res *payload.Search_Response, err error) {
	if f.IsIndexing() {
		return nil, errors.ErrCreateIndexingIsInProgress
	}
	if nprobe == 0 {
		nprobe = 1
	}

	sr, err := f.core.Search(searchK(k, max(f.Len(), atomic.LoadUint64(&f.icnt))), int(nprobe), int(nq), xq)
	if err != nil {
		if f.IsIndexing() {
			return nil, errors.ErrCreateIndexingIsInProgress
		}
		if errors.Is(err, errors.ErrSearchResultEmptyButNoDataStored) && f.Len() == 0 {
			return nil, nil
		}
		if errors.Is(err, errors.ErrNaNOrInfVector) {
			return nil, err
		}
		log.Errorf("cgo error detected during search: faiss api returned error %v", err)
		return nil, err
	}

	return f.toSearchResponse(sr)
}

// searchK bounds the client-supplied k (unbounded uint32) before the core
// search, which allocates k*nq labels and distances up front. FAISS cannot
// return more neighbours than the index holds and pads the rest with -1,
// which toSearchResult drops, so capping k at an upper bound of the index size
// (n) leaves the results unchanged. It is also capped at MaxInt32 so the C.int
// conversion in core cannot wrap.
func searchK(k uint32, n uint64) int {
	return int(min(uint64(k), max(n, 1), math.MaxInt32))
}

func (f *faiss) Delete(uuid string) (err error) {
	return f.delete(uuid, time.Now().UnixNano(), true)
}

func (f *faiss) DeleteWithTime(uuid string, t int64) (err error) {
	if t <= 0 {
		t = time.Now().UnixNano()
	}

	return f.delete(uuid, t, true)
}

func (f *faiss) delete(uuid string, t int64, validation bool) error {
	if len(uuid) == 0 {
		return errors.ErrUUIDNotFound(0)
	}

	if validation {
		_, _, ok := f.kvs.Get(uuid)
		_, ivqok := f.vq.IVExists(uuid)
		if !ok && !ivqok {
			return errors.ErrObjectIDNotFound(uuid)
		}
	}

	return f.vq.PushDelete(uuid, t)
}

func (f *faiss) Exists(uuid string) (oid uint32, ok bool) {
	return memstore.Exists(f.kvs, f.vq, uuid)
}

func (f *faiss) GetObject(uuid string) (vec []float32, timestamp int64, err error) {
	return memstore.GetObject(f.kvs, f.vq, uuid, nil)
}

func (f *faiss) IsIndexing() bool {
	i, ok := f.indexing.Load().(bool)
	return i && ok
}

func (f *faiss) IsSaving() bool {
	s, ok := f.saving.Load().(bool)
	return s && ok
}

func (f *faiss) UUIDs(ctx context.Context) (uuids []string) {
	return memstore.UUIDs(ctx, f.kvs, f.vq)
}

func (f *faiss) NumberOfCreateIndexExecution() uint64 {
	return f.nocie.Load()
}

func (f *faiss) NumberOfProactiveGCExecution() uint64 {
	return f.nogce.Load()
}

func (f *faiss) gc() {
	if f.enableProactiveGC {
		runtime.GC()
		f.nogce.Add(1)
	}
}

func (f *faiss) Len() uint64 {
	return f.kvs.Len()
}

func (f *faiss) InsertVQueueBufferLen() uint64 {
	return uint64(f.vq.IVQLen())
}

func (f *faiss) DeleteVQueueBufferLen() uint64 {
	return uint64(f.vq.DVQLen())
}

func (f *faiss) GetDimensionSize() int {
	return f.dim
}

func (f *faiss) GetTrainSize() int {
	return int(f.trainSize.Load())
}

func (f *faiss) Close(ctx context.Context) (err error) {
	defer f.core.Close()
	defer func() {
		if !errors.IsNot(err, context.Canceled, context.DeadlineExceeded) {
			err = nil
		}
		kerr := f.kvs.Close()
		if kerr != nil {
			if err != nil {
				err = errors.Join(kerr, err)
			} else {
				err = kerr
			}
		}
	}()
	if len(f.path) != 0 {
		// An auto-index CreateIndex/SaveIndex started by Start may still be running: its batch is
		// already popped from the vqueue (so the CreateIndex below sees nothing to do) while nocie
		// is not bumped yet (so the SaveIndex below is deduplicated). Wait for it, bounded by ctx,
		// so the final SaveIndex persists that batch. Both flags are cleared by defers, and a
		// CreateIndex waiter (wfci) leaves its wait loop on its own ctx or once indexing ends.
		if werr := f.waitForInflightIndexing(ctx); werr != nil {
			log.Warnf("stopped waiting for in-flight create/save index operation on close: %v", werr)
		}
		cerr := f.CreateIndex(ctx)
		if errors.IsNot(cerr, errors.ErrUncommittedIndexNotFound, context.Canceled, context.DeadlineExceeded) {
			if err != nil {
				err = errors.Join(cerr, err)
			} else {
				err = cerr
			}
		}
		serr := f.SaveIndex(ctx)
		if errors.IsNot(serr, errors.ErrUncommittedIndexNotFound, context.Canceled, context.DeadlineExceeded) {
			if err != nil {
				err = errors.Join(serr, err)
			} else {
				err = serr
			}
		}
	}
	return err
}

// waitForInflightIndexing blocks until no CreateIndex/SaveIndex is running or waiting, or ctx is done.
func (f *faiss) waitForInflightIndexing(ctx context.Context) error {
	ticker := time.NewTicker(inflightIndexingPollInterval)
	defer ticker.Stop()
	for f.IsIndexing() || f.IsSaving() || f.wfci.Load() > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
	return nil
}

// ListObjectFunc applies the input function on each index stored in the kvs and vqueue.
// Use this function for performing something on each object with caring about the memory usage.
// If the vector exists in the vqueue, this vector is not indexed so the oid(object ID) is processed as 0.
func (f *faiss) ListObjectFunc(
	ctx context.Context, fn func(uuid string, oid uint32, ts int64) bool,
) {
	memstore.ListObjectFunc(ctx, f.kvs, f.vq, fn)
}

func (f *faiss) toSearchResponse(
	sr []algorithm.SearchResult,
) (res *payload.Search_Response, err error) {
	if len(sr) == 0 {
		if f.Len() == 0 {
			return nil, nil
		}
		return nil, errors.ErrEmptySearchResult
	}

	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(sr)),
	}
	// One backing array for all results instead of one allocation per result.
	dists := make([]payload.Object_Distance, len(sr))
	for _, d := range sr {
		if err = d.Error; d.ID == 0 && err != nil {
			log.Warnf("an error occurred while searching: %v", err)
			continue
		}
		if math.IsNaN(float64(d.Distance)) {
			// same policy as the LB aggregation rankable: a NaN distance is not rankable
			continue
		}
		key, _, ok := f.kvs.GetInverse(d.ID)
		if ok {
			dist := &dists[len(res.Results)]
			dist.Id = key
			dist.Distance = d.Distance
			res.Results = append(res.Results, dist)
		} else {
			log.Warn("not found", d.ID, d.Distance)
		}
	}
	if len(res.GetResults()) == 0 {
		if f.Len() == 0 {
			return nil, nil
		}
		return nil, errors.ErrEmptySearchResult
	}
	return res, nil
}
