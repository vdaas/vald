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

package observer

import (
	"archive/tar"
	"context"
	"io/fs"
	"path/filepath"
	"reflect"
	"syscall"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/file/watch"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type StorageObserver interface {
	Start(ctx context.Context) (<-chan error, error)
	PostStop(ctx context.Context) error
}

type observer struct {
	w               watch.Watcher
	eg              errgroup.Group
	storage         storage.Storage
	ch              chan struct{}
	dir             string
	metadataPath    string
	hooks           []Hook
	checkDuration   time.Duration
	postStopTimeout time.Duration
	watchEnabled    bool
	tickerEnabled   bool
}

func New(opts ...Option) (so StorageObserver, err error) {
	o := new(observer)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(o); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	o.w, err = watch.New(
		watch.WithDirs(o.watchDirs()...),
		watch.WithErrGroup(o.eg),
		watch.WithOnWrite(o.onWrite),
		watch.WithOnCreate(o.onCreate),
	)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// errSaveInProgress reports that the backup walk saw the agent's in-place index
// save still running. The upload is aborted and the backup is skipped, not
// failed: the next metadata write event or ticker tick retries it. Fresh mtimes
// after a restore or a migrate() copy also look like this until the agent's
// next save, so reporting it as an error would fail every tick and the exit.
var errSaveInProgress = errors.New("an index save is in progress")

// originIndexDirName is the subdirectory of the index path that the NGT agent
// saves its index into (pkg/agent/core/ngt/service prepareFolders).
const originIndexDirName = "origin"

// originDir is <dir>/origin. Without copy-on-write the chart sets watch_dir to
// the index path, but the agent saves the index into <index path>/origin, so
// the metadata file at o.metadataPath never appears there.
func (o *observer) originDir() string {
	return file.Join(o.dir, originIndexDirName)
}

// indexDir returns the directory holding the agent's latest index and its
// metadata file: <dir>/origin when it has the metadata (it is where the agent
// saves; a metadata file in <dir> itself is then a stale restored copy that the
// agent no longer uses), otherwise <dir>. inPlace reports the <dir>/origin case,
// where the agent rewrites the index files in place.
func (o *observer) indexDir() (dir, metadataPath string, inPlace bool) {
	od := o.originDir()
	mp := file.Join(od, metadata.AgentMetadataFileName)
	if file.Exists(mp) {
		return od, mp, true
	}
	return o.dir, o.metadataPath, false
}

// isMetadataPath reports whether name is the metadata file of either layout.
func (o *observer) isMetadataPath(name string) bool {
	return name == o.metadataPath || name == file.Join(o.originDir(), metadata.AgentMetadataFileName)
}

// watchDirs returns o.dir and, when it already exists, <dir>/origin: the
// watcher is not recursive, and the metadata file may be written there.
func (o *observer) watchDirs() []string {
	if isDir(o.originDir()) {
		return []string{o.dir, o.originDir()}
	}
	return []string{o.dir}
}

// isDir reports whether path is an existing directory. A path that cannot be
// stat'ed is not watched; it is added when its create event arrives.
func isDir(path string) bool {
	exists, fi, err := file.ExistsWithDetail(path)
	return exists && err == nil && fi != nil && fi.IsDir()
}

// watchOrigin adds <dir>/origin to w when name is it, so a metadata file saved
// into an origin directory created after the watcher started is still seen.
func (o *observer) watchOrigin(w watch.Watcher, name string) error {
	if w == nil || name != o.originDir() || !isDir(name) {
		return nil
	}
	return w.Add(name)
}

func (o *observer) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)

	var wech, tech, sech, bech <-chan error
	var err error

	if o.watchEnabled {
		wech, err = o.w.Start(ctx)
		if err != nil {
			return nil, err
		}
	}

	if o.tickerEnabled {
		tech, err = o.startTicker(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}

	sech, err = o.storage.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	bech, err = o.startBackupLoop(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-wech:
			case err = <-tech:
			case err = <-sech:
			case err = <-bech:
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

func (o *observer) PostStop(ctx context.Context) (err error) {
	defer func() {
		if serr := o.storage.Stop(ctx); serr != nil {
			err = errors.Join(err, serr)
		}
	}()

	finalize := func() (err error) {
		err = ctx.Err()
		if err != nil && err != context.Canceled {
			return err
		}
		return nil
	}

	backup := func() (err error) {
		_, mpath, _ := o.indexDir()
		metadata, err := metadata.Load(mpath)
		if err != nil {
			log.Warn("cannot read metadata of the backup files:", err)
			return err
		}

		if metadata.IsInvalid {
			log.Warn("backup skipped because the files are invalid")
			return nil
		}

		err = o.backup(ctx)
		if errors.Is(err, errSaveInProgress) {
			// The agent has stopped, so this is not a running save but mtimes
			// left by a restore or migrate() copy; the remote already has it.
			log.Info("final backup skipped:", err)
			return nil
		}
		return err
	}

	// f runs on the watcher goroutine while the loop below reads the last write
	// time, so keep it atomic. ch is never closed and only ever holds one pending
	// signal: a metadata store fires several events, and a blocking send (or a
	// send after close) would hang or panic the watcher goroutine at shutdown.
	var last atomic.Int64
	last.Store(time.Now().UnixNano())
	ch := make(chan struct{}, 1)

	var w watch.Watcher
	f := func(ctx context.Context, name string) error {
		if o.isMetadataPath(name) {
			select {
			case ch <- struct{}{}:
			default:
			}
			return nil
		}

		last.Store(time.Now().UnixNano())
		return o.watchOrigin(w, name)
	}

	w, err = watch.New(
		watch.WithDirs(o.watchDirs()...),
		watch.WithErrGroup(o.eg),
		watch.WithOnWrite(f),
		watch.WithOnCreate(f),
	)
	if err != nil {
		return err
	}
	o.w = w

	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err = o.w.Start(wctx)
	if err != nil {
		return err
	}
	defer func() {
		e := o.w.Stop(wctx)
		if e != nil {
			log.Error("an error occurred when watcher stopped:", e)
		}
	}()

	ticker := time.NewTicker(postStopCheckInterval(o.postStopTimeout))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return finalize()
		case <-ch:
			return backup()
		case <-ticker.C:
			if time.Since(time.Unix(0, last.Load())) > o.postStopTimeout {
				return backup()
			}
		}
	}
}

func (o *observer) startTicker(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		ct := time.NewTicker(o.checkDuration)
		defer ct.Stop()

		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				return err
			}
			return nil
		}

		for {
			select {
			case <-ctx.Done():
				return finalize()
			case <-ct.C:
				_, mpath, _ := o.indexDir()
				metadata, err := metadata.Load(mpath)
				if err != nil {
					log.Warn("cannot read metadata of the backup files:", err)
					select {
					case <-ctx.Done():
						return finalize()
					case ech <- err:
					}
					continue
				}

				if metadata.IsInvalid {
					log.Warn("backup skipped because the files are invalid")
					continue
				}

				if err = o.requestBackup(ctx); err != nil {
					log.Error("failed to request backup:", err)
					select {
					case <-ctx.Done():
						return finalize()
					case ech <- err:
					}
				}
			}
		}
	}))

	return ech, nil
}

func (o *observer) startBackupLoop(ctx context.Context) (<-chan error, error) {
	o.ch = make(chan struct{}, 1)

	ech := make(chan error, 100)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)

		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				return err
			}
			return nil
		}

		for {
			select {
			case <-ctx.Done():
				return finalize()
			case <-o.ch:
				if err := o.backup(ctx); err != nil {
					if errors.Is(err, errSaveInProgress) {
						log.Info("backup skipped and will be retried:", err)
						continue
					}
					log.Error("an error occurred during backup:", err)
					select {
					case <-ctx.Done():
						return finalize()
					case ech <- err:
					}
				}
			}
		}
	}))

	return ech, nil
}

func (o *observer) onWrite(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.onWrite")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if !o.isMetadataPath(name) {
		return nil
	}

	ok, err := o.isValidMetadata(name)
	if err != nil {
		log.Warn("cannot read metadata of the backup files:", err)
		return err
	}

	if ok {
		return o.requestBackup(ctx)
	}

	return o.terminate()
}

func (o *observer) onCreate(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.onCreate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if err := o.watchOrigin(o.w, name); err != nil {
		return err
	}

	if !o.isMetadataPath(name) {
		return nil
	}

	ok, err := o.isValidMetadata(name)
	if err != nil {
		log.Warn("cannot read metadata of the backup files:", err)
		return err
	}

	if ok {
		return o.requestBackup(ctx)
	}

	return o.terminate()
}

func (*observer) isValidMetadata(path string) (bool, error) {
	metadata, err := metadata.Load(path)
	if err != nil {
		return false, err
	}

	return !metadata.IsInvalid, nil
}

func (*observer) terminate() error {
	log.Error("the process will be terminated because the files are invalid")

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}

	return p.Signal(syscall.SIGTERM)
}

func (o *observer) requestBackup(context.Context) error {
	select {
	case o.ch <- struct{}{}:
	default:
		log.Debug("cannot request backup: channel is full")
	}

	return nil
}

func (o *observer) backup(ctx context.Context) (err error) {
	bi := &BackupInfo{
		StorageInfo: o.storage.StorageInfo(),
	}
	bi.StartTime = time.Now()

	for _, hook := range o.hooks {
		ctx, err = hook.BeforeProcess(ctx, bi)
		if err != nil {
			return err
		}
	}

	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/observer/StorageObserver.backup")
	if span != nil {
		span.SetAttributes(
			attribute.String("storage_type", bi.Type),
			attribute.String("bucket_name", bi.BucketName),
			attribute.String("filename", bi.Filename),
		)
	}
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	dir, metadataPath, inPlace := o.indexDir()
	log.Infof("started to backup directory %s", dir)

	// Taken before the walk and compared after it (see below).
	pmfi, err := os.Stat(metadataPath)
	if err != nil {
		return errors.Wrap(err, "failed to stat metadata before backup walk")
	}

	pr, pw := io.Pipe()
	defer func() {
		e := pr.Close()
		if e != nil {
			log.Errorf("error on closing pipe reader: %s", e)
		}
	}()

	// The blob-storage writer gets its own cancellable context. Close commits
	// the upload, so on a failure path the context is cancelled first: the S3
	// uploader and the gocloud writer then fail their final request instead of
	// committing a partial object over the last good backup.
	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sw, err := o.storage.Writer(wctx)
	if err != nil {
		return err
	}
	// swClosed is set once the success path has closed sw itself and checked
	// the result, so the deferred close below does not close it twice.
	swClosed := false
	defer func() {
		if swClosed {
			return
		}
		cancel()
		e := sw.Close()
		if e != nil {
			log.Errorf("error on closing blob-storage writer: %s", e)
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	walk := safety.RecoverFunc(func() (err error) {
		tw := tar.NewWriter(pw)
		defer func() {
			// tar reports a short final entry (a file that shrank after its
			// header was written) only from Close, so the error must reach the
			// pipe reader via err; logging alone would commit a corrupt archive.
			e := tw.Close()
			if e != nil {
				log.Errorf("error on closing tar writer: %s", e)
				err = errors.Join(err, e)
			}
		}()

		return filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err != nil {
				return err
			}

			// In the origin layout the agent rewrites the index files in place
			// and writes the metadata last, so a file newer than the metadata
			// taken before the walk belongs to a save that is still running; the
			// metadata check after the walk cannot see a save that has not
			// finished yet. With copy-on-write a save swaps in a whole directory
			// whose mtimes the copy fallback of MoveDir does not keep, so the
			// check is limited to this layout and the swap is caught by the
			// metadata check below.
			if inPlace && fi.Mode().IsRegular() && path != metadataPath && fi.ModTime().After(pmfi.ModTime()) {
				return errors.Errorf("index file %s is newer than metadata %s: %w", path, metadataPath, errSaveInProgress)
			}

			header, err := tar.FileInfoHeader(fi, path)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			header.Name = filepath.ToSlash(rel)

			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}

			log.Debug("writing: ", path)

			if fi.IsDir() {
				return nil
			}

			data, err := file.Open(path, os.O_RDONLY, fs.ModePerm)
			if err != nil {
				return err
			}
			defer func() {
				e := data.Close()
				if e != nil {
					log.Errorf("failed to close %s: %s", path, e)
				}
			}()

			d, err := io.NewReaderWithContext(ctx, data)
			if err != nil {
				return err
			}

			_, err = io.Copy(tw, d)
			if err != nil {
				return err
			}
			return nil
		})
	})
	o.eg.Go(func() error {
		defer wg.Done()
		// Propagate the walk/tar error (including a recovered panic) to the pipe
		// reader (the blob-storage uploader). A plain Close() surfaces a clean
		// io.EOF, so a truncated archive from a mid-walk failure would be
		// uploaded and reported as a successful backup. On success err is nil
		// and this behaves like Close(). PipeWriter.CloseWithError never returns
		// a non-nil error, but keep the guard so the result is not silently
		// discarded.
		if e := pw.CloseWithError(walk()); e != nil {
			log.Errorf("error on closing pipe writer: %s", e)
		}
		// The error reaches backup()'s caller through the pipe only. o.eg is the
		// process-wide errgroup: returning the error there would cancel the
		// daemon context and fail the shutdown for a backup that is retried.
		return nil
	})

	prr, err := io.NewReaderWithContext(ctx, pr)
	if err != nil {
		return err
	}

	bi.Bytes, err = io.Copy(sw, prr)
	if err != nil {
		// When another task on the errgroup fails, the errgroup cancels ctx with
		// that error as the cause, and the reader may see the cancellation before
		// the pipe error. Report the cause so the failure is not masked as
		// "context canceled".
		if cause := context.Cause(ctx); cause != nil && errors.Is(err, context.Canceled) {
			return cause
		}
		return err
	}

	wg.Wait()

	// The agent writes metadata last when it saves and replaces it when the
	// copy-on-write swap moves a new generation in, so a metadata file that
	// changed during the walk means the archive may mix index generations.
	// Fail before Close so the upload is aborted; the next save event retries.
	mfi, err := os.Stat(metadataPath)
	if err != nil {
		return errors.Wrap(err, "failed to stat metadata after backup walk")
	}
	if !os.SameFile(pmfi, mfi) || !pmfi.ModTime().Equal(mfi.ModTime()) || pmfi.Size() != mfi.Size() {
		return errors.Errorf("index metadata %s changed during backup; the index was saved while it was archived", metadataPath)
	}

	// Closing the blob-storage writer is what finalises the upload (e.g. an S3
	// multipart upload), so its error must be checked before any success hook runs.
	swClosed = true
	err = sw.Close()
	if err != nil {
		return err
	}

	bi.EndTime = time.Now()
	for _, hook := range o.hooks {
		err = hook.AfterProcess(ctx, bi)
		if err != nil {
			return err
		}
	}

	log.Infof("finished to backup directory %s", o.dir)

	return nil
}

// postStopCheckDivisor and maxPostStopCheckInterval shape postStopCheckInterval.
const (
	postStopCheckDivisor     = 5
	maxPostStopCheckInterval = 2 * time.Second
)

// postStopCheckInterval returns the polling interval used while waiting for the backup metadata
// after the agent stopped: a fifth of timeout, capped at 2s. It is clamped to at least 1ns because
// time.NewTicker panics for a non-positive interval (timeout < 5ns would otherwise yield 0).
func postStopCheckInterval(timeout time.Duration) time.Duration {
	return max(min(timeout/postStopCheckDivisor, maxPostStopCheckInterval), time.Nanosecond)
}
