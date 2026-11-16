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
	"bytes"
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

// abortWriter records whether the context handed to the blob-storage writer
// was already cancelled when Close was called. Close commits the upload (S3
// PutObject/CompleteMultipartUpload, gocloud blob.Writer), and a cancelled
// context is what makes those final requests fail instead of committing.
type abortWriter struct {
	writeErr    error
	closeCtxErr error
	// ctxErr is the Err method of the context handed to Writer.
	ctxErr func() error
	// onWrite, when set, is called with the cumulative byte count after each Write.
	onWrite func(total int)
	total   int
	closed  int
}

func (w *abortWriter) Write(p []byte) (int, error) {
	if w.writeErr != nil {
		return 0, w.writeErr
	}
	w.total += len(p)
	if w.onWrite != nil {
		w.onWrite(w.total)
	}
	return len(p), nil
}

func (w *abortWriter) Close() error {
	w.closed++
	w.closeCtxErr = w.ctxErr()
	return nil
}

type abortWriterStorage struct {
	writerStorage
	w *abortWriter
}

func (s abortWriterStorage) Writer(ctx context.Context) (io.WriteCloser, error) {
	s.w.ctxErr = ctx.Err
	return s.w, nil
}

func runBackup(t *testing.T, dir string, w *abortWriter, h *recordHook) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)

	o := &observer{
		eg:           eg,
		dir:          dir,
		metadataPath: writeMetadata(t, dir),
		storage:      abortWriterStorage{w: w},
		hooks:        []Hook{h},
	}

	// backup is bounded by ctx: the pipe reads and the tar walk honour it.
	err := o.backup(egctx)
	if ctx.Err() != nil {
		t.Fatalf("backup did not finish in time: %v", ctx.Err())
	}

	cancel()
	if werr := eg.Wait(); werr != nil && !errors.Is(werr, context.Canceled) {
		t.Logf("eg.Wait() = %v", werr)
	}
	return err
}

// TestObserver_backup_abortOnFailure pins that a failed backup cancels the
// writer context before closing the writer, so the partial stream is aborted
// instead of committed over the last good backup, while a successful backup
// closes the writer with a live context.
func TestObserver_backup_abortOnFailure(t *testing.T) {
	t.Parallel()

	errWrite := errors.New("upload stream broken")

	tests := []struct {
		writeErr  error
		name      string
		wantAbort bool
	}{
		{name: "copy failure aborts the upload", writeErr: errWrite, wantAbort: true},
		{name: "success commits with a live context", writeErr: nil, wantAbort: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("hello"), 0o600); err != nil {
				t.Fatal(err)
			}
			w := &abortWriter{writeErr: tt.writeErr}
			h := new(recordHook)

			err := runBackup(t, dir, w, h)
			if tt.writeErr != nil {
				if !errors.Is(err, tt.writeErr) {
					t.Fatalf("backup() error = %v, want %v", err, tt.writeErr)
				}
			} else if err != nil {
				t.Fatalf("backup() error = %v, want nil", err)
			}
			if w.closed != 1 {
				t.Fatalf("writer closed %d times, want exactly 1", w.closed)
			}
			if gotAbort := errors.Is(w.closeCtxErr, context.Canceled); gotAbort != tt.wantAbort {
				t.Errorf("writer context cancelled at Close = %v (ctx err %v), want %v", gotAbort, w.closeCtxErr, tt.wantAbort)
			}
		})
	}
}

// TestObserver_backup_tarCloseError pins that an error reported only by
// tar.Writer.Close (a file that shrank after its header was written) fails the
// backup instead of committing a corrupt archive and running the success hooks.
func TestObserver_backup_tarCloseError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "f.bin")
	// Far larger than the few 32 KiB chunks the walker can read ahead of the
	// consumer, so the truncation below always lands before the file is done.
	if err := os.WriteFile(path, make([]byte, 1<<20), 0o600); err != nil {
		t.Fatal(err)
	}

	truncated := false
	var truncErr error
	w := &abortWriter{
		onWrite: func(total int) {
			// Two 512-byte headers (root dir, f.bin) come first; any byte
			// beyond them is file content, so the header with the full size
			// has already been written.
			if truncated || total <= 1024 {
				return
			}
			truncated = true
			truncErr = os.Truncate(path, 0)
		},
	}
	h := new(recordHook)

	err := runBackup(t, dir, w, h)
	if !truncated || truncErr != nil {
		t.Fatalf("truncation did not happen: truncated=%v err=%v", truncated, truncErr)
	}
	if err == nil || !strings.Contains(err.Error(), "missed writing") {
		t.Fatalf("backup() error = %v, want the tar close error", err)
	}
	if h.after != 0 {
		t.Errorf("AfterProcess called %d times, want 0", h.after)
	}
	if !errors.Is(w.closeCtxErr, context.Canceled) {
		t.Errorf("writer context not cancelled at Close (ctx err %v): the corrupt archive would be committed", w.closeCtxErr)
	}
}

// TestObserver_backup_metadataChangedAborts pins that a backup whose walk
// overlapped an index save fails and aborts the upload. The agent writes the
// metadata file last when it saves (and swaps it in with the copy-on-write
// move), so a change to it during the walk means the archive may mix index
// generations and must not be committed over the last good backup.
func TestObserver_backup_metadataChangedAborts(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	mpath := filepath.Join(dir, metadata.AgentMetadataFileName)
	staged := filepath.Join(t.TempDir(), metadata.AgentMetadataFileName)

	saved := false
	var saveErr error
	w := &abortWriter{
		onWrite: func(int) {
			if saved {
				return
			}
			saved = true
			// Simulates the agent swapping in a newly saved generation while
			// the walk is running. The replacement is a new inode with the same
			// content, so the change is detected independently of the mtime
			// resolution, and the walk never sees a size mismatch (a tar error
			// would fail the backup for a different reason). It is staged
			// outside dir so the walk cannot visit it.
			if saveErr = os.WriteFile(staged, []byte(`{"is_invalid":false}`), 0o600); saveErr != nil {
				return
			}
			saveErr = os.Rename(staged, mpath)
		},
	}
	h := new(recordHook)

	err := runBackup(t, dir, w, h)
	if !saved || saveErr != nil {
		t.Fatalf("metadata was not rewritten during the backup: saved=%v err=%v", saved, saveErr)
	}
	if err == nil || !strings.Contains(err.Error(), "changed during backup") {
		t.Fatalf("backup() error = %v, want the metadata change error", err)
	}
	if h.after != 0 {
		t.Errorf("AfterProcess called %d times, want 0", h.after)
	}
	if w.closed != 1 {
		t.Fatalf("writer closed %d times, want exactly 1", w.closed)
	}
	if !errors.Is(w.closeCtxErr, context.Canceled) {
		t.Errorf("writer context not cancelled at Close (ctx err %v): the mixed archive would be committed", w.closeCtxErr)
	}
}

// TestObserver_backup_failureNotGroupError pins that a failed backup is
// reported only by backup() (and the observer's error channel), not also
// returned to the errgroup. In the sidecar that is the process-wide group,
// where the walk goroutine's io.ErrClosedPipe cancelled the daemon context and
// was returned at shutdown, failing it even after a successful final backup.
func TestObserver_backup_failureNotGroupError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	errWrite := errors.New("upload stream broken")
	w := &abortWriter{writeErr: errWrite}
	o := &observer{
		eg:           eg,
		dir:          dir,
		metadataPath: writeMetadata(t, dir),
		storage:      abortWriterStorage{w: w},
		hooks:        []Hook{new(recordHook)},
	}

	if err := o.backup(egctx); !errors.Is(err, errWrite) {
		t.Fatalf("backup() error = %v, want %v", err, errWrite)
	}

	waited := make(chan error, 1)
	go func() { waited <- eg.Wait() }()
	select {
	case err := <-waited:
		if err != nil {
			t.Errorf("eg.Wait() = %v, want nil: the backup error belongs to backup()", err)
		}
	case <-ctx.Done():
		t.Fatalf("eg.Wait did not return: %v", ctx.Err())
	}
}

// closeErrWriter is a blob-storage writer whose Close reports closeErr, like an
// S3 multipart upload that fails when it is finalised.
type closeErrWriter struct {
	closeErr error
	closed   int
}

func (*closeErrWriter) Write(p []byte) (int, error) { return len(p), nil }

func (w *closeErrWriter) Close() error {
	w.closed++
	return w.closeErr
}

type writerStorage struct {
	storage.Storage
	w io.WriteCloser
}

func (writerStorage) StorageInfo() *storage.StorageInfo { return new(storage.StorageInfo) }

func (s writerStorage) Writer(context.Context) (io.WriteCloser, error) { return s.w, nil }

type recordHook struct {
	after int
}

func (*recordHook) BeforeProcess(ctx context.Context, _ *BackupInfo) (context.Context, error) {
	return ctx, nil
}

func (h *recordHook) AfterProcess(context.Context, *BackupInfo) error {
	h.after++
	return nil
}

// writeMetadata creates the agent metadata file in dir, which backup requires
// to exist, and returns its path.
func writeMetadata(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, metadata.AgentMetadataFileName)
	if err := os.WriteFile(path, []byte(`{"is_invalid":false}`), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestObserver_backup_closeError pins that the success hooks run only after the
// blob-storage writer was closed successfully. Closing the writer is where the
// upload is finalised, so a Close error means the backup failed.
func TestObserver_backup_closeError(t *testing.T) {
	t.Parallel()

	errClose := errors.New("finalising upload failed")

	tests := []struct {
		closeErr  error
		wantErr   error
		name      string
		wantAfter int
	}{
		{name: "writer close error is returned and success hooks are skipped", closeErr: errClose, wantErr: errClose, wantAfter: 0},
		{name: "writer close success runs the hooks", closeErr: nil, wantErr: nil, wantAfter: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)

			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("hello"), 0o600); err != nil {
				t.Fatal(err)
			}

			w := &closeErrWriter{closeErr: tt.closeErr}
			h := new(recordHook)
			o := &observer{
				eg:           eg,
				dir:          dir,
				metadataPath: writeMetadata(t, dir),
				storage:      writerStorage{w: w},
				hooks:        []Hook{h},
			}

			// backup is bounded by ctx: the pipe reads and the tar walk honour it.
			err := o.backup(egctx)
			if ctx.Err() != nil {
				t.Fatalf("backup did not finish in time: %v", ctx.Err())
			}
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("backup() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("backup() error = %v, want %v", err, tt.wantErr)
			}
			if h.after != tt.wantAfter {
				t.Errorf("AfterProcess called %d times, want %d", h.after, tt.wantAfter)
			}
			if w.closed != 1 {
				t.Errorf("writer closed %d times, want exactly 1", w.closed)
			}

			cancel()
			if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
				t.Logf("eg.Wait() = %v", err)
			}
		})
	}
}

type bufWriteCloser struct {
	bytes.Buffer
}

func (*bufWriteCloser) Close() error { return nil }

func tarNames(t *testing.T, b []byte) []string {
	t.Helper()
	var names []string
	tr := tar.NewReader(bytes.NewReader(b))
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return names
		}
		if err != nil {
			t.Fatalf("reading the archive: %v", err)
		}
		names = append(names, h.Name)
	}
}

type addRecordWatcher struct {
	added []string
}

func (*addRecordWatcher) Start(context.Context) (<-chan error, error) {
	return make(chan error), nil
}

func (w *addRecordWatcher) Add(dirs ...string) error {
	w.added = append(w.added, dirs...)
	return nil
}

func (*addRecordWatcher) Remove(...string) error { return nil }

func (*addRecordWatcher) Stop(context.Context) error { return nil }

// TestObserver_backup_originLayout pins that the backup finds the index the
// agent saves into <watch_dir>/origin. Without copy-on-write the chart sets
// watch_dir to the index path while the agent saves into <index path>/origin,
// so <watch_dir>/metadata.json never exists and every backup failed. The
// archive is rooted at origin, matching what the restorer extracts into
// watch_dir and the agent then migrates into origin.
func TestObserver_backup_originLayout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rootStale bool
	}{
		{name: "only origin holds the metadata"},
		{name: "origin wins over a stale restored root copy", rootStale: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)

			root := t.TempDir()
			origin := filepath.Join(root, originIndexDirName)
			if err := os.Mkdir(origin, 0o700); err != nil {
				t.Fatal(err)
			}
			writeMetadata(t, origin)
			kvs := filepath.Join(origin, "ngt-meta.kvsdb")
			if err := os.WriteFile(kvs, []byte("new"), 0o600); err != nil {
				t.Fatal(err)
			}
			// A finished save wrote the index files before the metadata.
			past := time.Now().Add(-time.Hour)
			if err := os.Chtimes(kvs, past, past); err != nil {
				t.Fatal(err)
			}
			if tt.rootStale {
				writeMetadata(t, root)
			}

			w := new(bufWriteCloser)
			h := new(recordHook)
			o := &observer{
				eg:           eg,
				dir:          root,
				metadataPath: filepath.Join(root, metadata.AgentMetadataFileName),
				storage:      writerStorage{w: w},
				hooks:        []Hook{h},
			}
			err := o.backup(egctx)
			if ctx.Err() != nil {
				t.Fatalf("backup did not finish in time: %v", ctx.Err())
			}
			cancel()
			if werr := eg.Wait(); werr != nil && !errors.Is(werr, context.Canceled) {
				t.Logf("eg.Wait() = %v", werr)
			}
			if err != nil {
				t.Fatalf("backup() error = %v, want nil", err)
			}
			if h.after != 1 {
				t.Fatalf("AfterProcess called %d times, want 1", h.after)
			}
			names := tarNames(t, w.Bytes())
			for _, want := range []string{metadata.AgentMetadataFileName, "ngt-meta.kvsdb"} {
				if !slices.Contains(names, want) {
					t.Errorf("archive entries %v miss %q", names, want)
				}
			}
			for _, n := range names {
				if n == originIndexDirName || filepath.Dir(n) == originIndexDirName {
					t.Errorf("archive entry %q is not rooted at origin: %v", n, names)
				}
			}
		})
	}
}

// TestObserver_onWrite_originMetadata pins that a metadata save into
// <watch_dir>/origin requests a backup, and that onCreate starts watching an
// origin directory created after the watcher (the watcher is not recursive).
func TestObserver_onWrite_originMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	origin := filepath.Join(root, originIndexDirName)
	if err := os.Mkdir(origin, 0o700); err != nil {
		t.Fatal(err)
	}
	w := new(addRecordWatcher)
	o := &observer{
		w:            w,
		ch:           make(chan struct{}, 1),
		dir:          root,
		metadataPath: filepath.Join(root, metadata.AgentMetadataFileName),
	}

	if err := o.onCreate(t.Context(), origin); err != nil {
		t.Fatalf("onCreate(origin) error = %v", err)
	}
	if !slices.Equal(w.added, []string{origin}) {
		t.Fatalf("watcher added %v, want [%s]", w.added, origin)
	}
	if got := o.watchDirs(); !slices.Equal(got, []string{root, origin}) {
		t.Fatalf("watchDirs() = %v, want [%s %s]", got, root, origin)
	}

	if err := o.onWrite(t.Context(), writeMetadata(t, origin)); err != nil {
		t.Fatalf("onWrite(origin metadata) error = %v", err)
	}
	select {
	case <-o.ch:
	default:
		t.Fatal("onWrite of <watch_dir>/origin/metadata.json did not request a backup")
	}
}

// TestObserver_backup_saveInProgress pins that, in the origin layout where the
// agent rewrites the index in place and writes the metadata last, a backup
// that sees an index file newer than the metadata fails and aborts the upload:
// the save is still running, so the archive would mix generations, and the
// metadata check after the walk cannot see a save that has not finished. A
// newer file in the copy-on-write layout (watch_dir holds the metadata) is not
// a save in progress, since the copy fallback of MoveDir does not keep mtimes.
func TestObserver_backup_saveInProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		inPlace bool
		wantErr bool
	}{
		{name: "origin layout with a newer index file fails", inPlace: true, wantErr: true},
		{name: "copy-on-write layout with a newer index file succeeds", inPlace: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)

			root := t.TempDir()
			idx := root
			if tt.inPlace {
				idx = filepath.Join(root, originIndexDirName)
				if err := os.Mkdir(idx, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			mpath := writeMetadata(t, idx)
			past := time.Now().Add(-time.Hour)
			if err := os.Chtimes(mpath, past, past); err != nil {
				t.Fatal(err)
			}
			// Rewritten by a save that has not written the metadata yet.
			if err := os.WriteFile(filepath.Join(idx, "ngt-meta.kvsdb"), []byte("new"), 0o600); err != nil {
				t.Fatal(err)
			}

			w := new(abortWriter)
			h := new(recordHook)
			o := &observer{
				eg:           eg,
				dir:          root,
				metadataPath: filepath.Join(root, metadata.AgentMetadataFileName),
				storage:      abortWriterStorage{w: w},
				hooks:        []Hook{h},
			}
			err := o.backup(egctx)
			if ctx.Err() != nil {
				t.Fatalf("backup did not finish in time: %v", ctx.Err())
			}
			cancel()
			if werr := eg.Wait(); werr != nil && !errors.Is(werr, context.Canceled) {
				t.Logf("eg.Wait() = %v", werr)
			}

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("backup() error = %v, want nil", err)
				}
				if h.after != 1 {
					t.Errorf("AfterProcess called %d times, want 1", h.after)
				}
				return
			}
			if !errors.Is(err, errSaveInProgress) {
				t.Fatalf("backup() error = %v, want errSaveInProgress", err)
			}
			if h.after != 0 {
				t.Errorf("AfterProcess called %d times, want 0", h.after)
			}
			if w.closed != 1 {
				t.Fatalf("writer closed %d times, want exactly 1", w.closed)
			}
			if !errors.Is(w.closeCtxErr, context.Canceled) {
				t.Errorf("writer context not cancelled at Close (ctx err %v): the mixed archive would be committed", w.closeCtxErr)
			}
		})
	}
}

// TestObserver_startBackupLoop_saveInProgressNotReported pins that a backup
// skipped because an in-place save is in progress aborts the upload (writer
// context cancelled at Close, no success hook) without reporting an error on
// the observer's channel: the runner would collect it and fail the exit, and
// after a restore or migrate() copy it would repeat until the agent's next save.
func TestObserver_startBackupLoop_saveInProgressNotReported(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)

	root := t.TempDir()
	idx := filepath.Join(root, originIndexDirName)
	if err := os.Mkdir(idx, 0o700); err != nil {
		t.Fatal(err)
	}
	mpath := writeMetadata(t, idx)
	past := time.Now().Add(-time.Hour)
	if err := os.Chtimes(mpath, past, past); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(idx, "ngt-meta.kvsdb"), []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}

	closed := make(chan struct{}, 1)
	w := &abortWriter{}
	h := &signalHook{before: closed}
	o := &observer{
		eg:           eg,
		dir:          root,
		metadataPath: filepath.Join(root, metadata.AgentMetadataFileName),
		storage:      abortWriterStorage{w: w},
		hooks:        []Hook{h},
	}
	lctx, lcancel := context.WithCancel(egctx)
	defer lcancel()
	ech, err := o.startBackupLoop(lctx)
	if err != nil {
		t.Fatalf("startBackupLoop() error = %v", err)
	}
	if err := o.requestBackup(lctx); err != nil {
		t.Fatalf("requestBackup() error = %v", err)
	}
	select {
	case <-closed:
	case <-ctx.Done():
		t.Fatal("the requested backup did not start in time")
	}
	// A second request is only taken once the first backup has returned, so
	// the first one has been handled (reported or skipped) when it starts.
	if err := o.requestBackup(lctx); err != nil {
		t.Fatalf("requestBackup() error = %v", err)
	}
	select {
	case <-closed:
	case <-ctx.Done():
		t.Fatal("the second backup did not start in time")
	}
	lcancel()

	for {
		select {
		case err, ok := <-ech:
			if !ok {
				if werr := eg.Wait(); werr != nil && !errors.Is(werr, context.Canceled) {
					t.Logf("eg.Wait() = %v", werr)
				}
				if h.after.Load() != 0 {
					t.Errorf("AfterProcess called %d times, want 0", h.after.Load())
				}
				return
			}
			// lcancel may interrupt the second backup, which is then reported
			// as cancelled; only the skip itself must not be reported.
			if err != nil && !errors.Is(err, context.Canceled) {
				t.Fatalf("observer reported %v, want the skipped backup not reported", err)
			}
		case <-ctx.Done():
			t.Fatal("the backup loop did not stop in time")
		}
	}
}

// signalHook signals each BeforeProcess call and counts AfterProcess calls.
type signalHook struct {
	before chan struct{}
	after  atomic.Int64
}

func (h *signalHook) BeforeProcess(ctx context.Context, _ *BackupInfo) (context.Context, error) {
	select {
	case h.before <- struct{}{}:
	default:
	}
	return ctx, nil
}

func (h *signalHook) AfterProcess(context.Context, *BackupInfo) error {
	h.after.Add(1)
	return nil
}

// time.NewTicker panics for a non-positive interval; timeout/5 is 0 for 1ns-4ns,
// and timeutil.Parse accepts such values.
func TestPostStopCheckInterval(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{name: "negative is clamped to 1ns", timeout: -time.Second, want: time.Nanosecond},
		{name: "zero is clamped to 1ns", timeout: 0, want: time.Nanosecond},
		{name: "1ns is clamped to 1ns", timeout: time.Nanosecond, want: time.Nanosecond},
		{name: "4ns is clamped to 1ns", timeout: 4 * time.Nanosecond, want: time.Nanosecond},
		{name: "5ns gives 1ns", timeout: 5 * time.Nanosecond, want: time.Nanosecond},
		{name: "10s gives a fifth", timeout: 10 * time.Second, want: 2 * time.Second},
		{name: "2m is capped at 2s", timeout: 2 * time.Minute, want: 2 * time.Second},
		{name: "1s gives a fifth", timeout: time.Second, want: 200 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := postStopCheckInterval(tt.timeout)
			if got != tt.want {
				t.Errorf("postStopCheckInterval(%v) = %v, want %v", tt.timeout, got, tt.want)
			}
			if got <= 0 {
				t.Errorf("interval must be positive, got %v", got)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		wantSo StorageObserver
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, StorageObserver, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotSo StorageObserver, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotSo, w.wantSo) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSo, w.wantSo)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotSo, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, gotSo, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			got, err := o.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_PostStop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			err := o.PostStop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_startTicker(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			got, err := o.startTicker(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_startBackupLoop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			got, err := o.startBackupLoop(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_onWrite(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		name string
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           name:"",
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           name:"",
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			err := o.onWrite(test.args.ctx, test.args.name)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_onCreate(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		name string
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           name:"",
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           name:"",
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			err := o.onCreate(test.args.ctx, test.args.name)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_isValidMetadata(t *testing.T) {
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		want bool
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			got, err := o.isValidMetadata()
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_terminate(t *testing.T) {
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			err := o.terminate()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_requestBackup(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			err := o.requestBackup(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_observer_backup(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		w               watch.Watcher
// 		eg              errgroup.Group
// 		storage         storage.Storage
// 		ch              chan struct{}
// 		dir             string
// 		metadataPath    string
// 		hooks           []Hook
// 		checkDuration   time.Duration
// 		postStopTimeout time.Duration
// 		watchEnabled    bool
// 		tickerEnabled   bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           w:nil,
// 		           eg:nil,
// 		           storage:nil,
// 		           ch:nil,
// 		           dir:"",
// 		           metadataPath:"",
// 		           hooks:nil,
// 		           checkDuration:nil,
// 		           postStopTimeout:nil,
// 		           watchEnabled:false,
// 		           tickerEnabled:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &observer{
// 				w:               test.fields.w,
// 				eg:              test.fields.eg,
// 				storage:         test.fields.storage,
// 				ch:              test.fields.ch,
// 				dir:             test.fields.dir,
// 				metadataPath:    test.fields.metadataPath,
// 				hooks:           test.fields.hooks,
// 				checkDuration:   test.fields.checkDuration,
// 				postStopTimeout: test.fields.postStopTimeout,
// 				watchEnabled:    test.fields.watchEnabled,
// 				tickerEnabled:   test.fields.tickerEnabled,
// 			}
//
// 			err := o.backup(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
