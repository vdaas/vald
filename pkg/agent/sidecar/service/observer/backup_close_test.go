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
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

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
