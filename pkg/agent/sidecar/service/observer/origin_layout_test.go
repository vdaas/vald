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
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
)

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
