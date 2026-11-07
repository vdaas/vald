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
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
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
