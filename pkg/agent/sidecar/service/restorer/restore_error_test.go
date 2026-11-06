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

package restorer

import (
	"archive/tar"
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

// Names of the files in the test tar archives.
const (
	fileA = "a.txt"
	fileB = "b.txt"
)

// errAfterReader yields data and then returns err (io.EOF when err is nil).
type errAfterReader struct {
	err  error
	data []byte
}

func (r *errAfterReader) Read(p []byte) (int, error) {
	if len(r.data) > 0 {
		n := copy(p, r.data)
		r.data = r.data[n:]
		return n, nil
	}
	if r.err == nil {
		return 0, io.EOF
	}
	return 0, r.err
}

func (*errAfterReader) Close() error { return nil }

type readerStorage struct {
	storage.Storage
	rc io.ReadCloser
}

func (s readerStorage) Reader(context.Context) (io.ReadCloser, error) {
	return s.rc, nil
}

// TestRestorer_restore_propagatesReadError pins that a read error from the blob
// reader reaches the tar consumer. The archive is cut right after a complete
// entry; with a plain pw.Close() the consumer saw a clean io.EOF at the entry
// boundary and restore() reported success for a truncated archive.
func TestRestorer_restore_propagatesReadError(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("blob read failed")

	archive := func(t *testing.T, trailer bool) []byte {
		t.Helper()
		content := []byte("hello")
		buf := new(bytes.Buffer)
		tw := tar.NewWriter(buf)
		if err := tw.WriteHeader(&tar.Header{
			Name:     "f.txt",
			Mode:     0o644,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(content); err != nil {
			t.Fatal(err)
		}
		if trailer {
			// Close writes the end-of-archive marker.
			if err := tw.Close(); err != nil {
				t.Fatal(err)
			}
		} else if err := tw.Flush(); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}

	tests := []struct {
		readErr error
		wantErr error
		name    string
		trailer bool
	}{
		{name: "read error after a complete entry is propagated", trailer: false, readErr: errBoom, wantErr: errBoom},
		{name: "complete archive without error succeeds", trailer: true, readErr: nil, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, egctx := errgroup.New(ctx)

			r := &restorer{
				eg:  eg,
				dir: t.TempDir(),
				storage: readerStorage{
					rc: &errAfterReader{data: archive(t, tt.trailer), err: tt.readErr},
				},
			}

			// restore is bounded by ctx: both the pipe reads and the tar loop honour it.
			err := r.restore(egctx)
			if ctx.Err() != nil {
				t.Fatalf("restore did not finish in time: %v", ctx.Err())
			}
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("restore() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("restore() error = %v, want %v", err, tt.wantErr)
			}

			// The copy goroutine hands its error to restore() through the pipe
			// only. Returning it to the errgroup, which is the process-wide group
			// in the init container, cancelled the daemon context and ended the
			// restore retry loop after a single failed attempt.
			waited := make(chan error, 1)
			go func() { waited <- eg.Wait() }()
			select {
			case err := <-waited:
				if err != nil {
					t.Errorf("eg.Wait() = %v, want nil", err)
				}
			case <-ctx.Done():
				t.Fatalf("eg.Wait did not return: %v", ctx.Err())
			}
			cancel()
		})
	}
}

// tarEntry is one regular file of a test archive; content may be shorter than
// size to cut the archive in the middle of the entry.
type tarEntry struct {
	name    string
	content []byte
	size    int64
}

func buildTar(t *testing.T, entries []tarEntry, trailer bool) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, e := range entries {
		if err := tw.WriteHeader(&tar.Header{
			Name:     e.name,
			Mode:     0o644,
			Size:     e.size,
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(e.content); err != nil {
			t.Fatal(err)
		}
	}
	if trailer {
		if err := tw.Close(); err != nil {
			t.Fatal(err)
		}
	} else if n := len(entries); n > 0 && int64(len(entries[n-1].content)) == entries[n-1].size {
		// Pad the last complete entry; a cut entry is left as is (its bytes
		// are already in buf, and Flush would reject the missing remainder).
		if err := tw.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	return buf.Bytes()
}

func runRestore(t *testing.T, dir string, data []byte, readErr error) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)
	r := &restorer{
		eg:      eg,
		dir:     dir,
		storage: readerStorage{rc: &errAfterReader{data: data, err: readErr}},
	}
	err := r.restore(egctx)
	if ctx.Err() != nil {
		t.Fatalf("restore did not finish in time: %v", ctx.Err())
	}
	cancel()
	if werr := eg.Wait(); werr != nil && !errors.Is(werr, context.Canceled) {
		t.Logf("eg.Wait() = %v", werr)
	}
	return err
}

// TestRestorer_restore_failedAttemptCleansUp pins that a failed attempt removes
// the files it wrote, including a partially written one. Previously the
// leftovers stayed in the directory and the retry stopped at the first one
// (ErrFileAlreadyExists treated as "already restored"), reporting success for
// an incomplete restore.
func TestRestorer_restore_failedAttemptCleansUp(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("blob read failed")
	full1 := []byte("first-file")
	full2 := []byte("second-file-content")

	tests := []struct {
		name    string
		partial []tarEntry
	}{
		{
			name:    "complete entry then read error",
			partial: []tarEntry{{name: fileA, content: full1, size: int64(len(full1))}},
		},
		{
			name: "complete entry then error inside the next entry",
			partial: []tarEntry{
				{name: fileA, content: full1, size: int64(len(full1))},
				{name: fileB, content: full2[:4], size: int64(len(full2))},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()

			// The error may surface as errBoom or as the cancellation it causes,
			// depending on which reader sees it first; only failure matters here.
			if err := runRestore(t, dir, buildTar(t, tt.partial, false), errBoom); err == nil {
				t.Fatal("first restore() succeeded, want an error")
			}
			for _, name := range []string{fileA, fileB} {
				if file.Exists(file.Join(dir, name)) {
					t.Errorf("%s left behind by the failed attempt", name)
				}
			}

			complete := buildTar(t, []tarEntry{
				{name: fileA, content: full1, size: int64(len(full1))},
				{name: fileB, content: full2, size: int64(len(full2))},
			}, true)
			if err := runRestore(t, dir, complete, nil); err != nil {
				t.Fatalf("retry restore() error = %v, want nil", err)
			}
			for name, want := range map[string][]byte{fileA: full1, fileB: full2} {
				got, err := file.ReadFile(file.Join(dir, name))
				if err != nil {
					t.Fatalf("read %s after retry: %v", name, err)
				}
				if !bytes.Equal(got, want) {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
		})
	}
}

// TestRestorer_restore_keepsPreExistingData pins that real pre-existing data is
// still left untouched: restore stops at the existing file and reports success.
func TestRestorer_restore_keepsPreExistingData(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	existing := []byte("local-data")
	if _, err := file.WriteFile(t.Context(), file.Join(dir, fileA), bytes.NewReader(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	remote := []byte("remote")
	data := buildTar(t, []tarEntry{{name: fileA, content: remote, size: int64(len(remote))}}, true)
	if err := runRestore(t, dir, data, nil); err != nil {
		t.Fatalf("restore() error = %v, want nil", err)
	}
	got, err := file.ReadFile(file.Join(dir, fileA))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, existing) {
		t.Errorf("a.txt = %q, want pre-existing %q", got, existing)
	}
}

// TestRestorer_restore_existingFileStopsCopyCleanly pins that restore() stopping
// early on an already restored file is not an error for the copy goroutine. It
// closes the pipe while the goroutine is still writing the rest of the archive;
// returning the resulting io.ErrClosedPipe recorded it on the errgroup, which is
// the process-wide group in the init container, so every restart over existing
// index files exited non-zero.
func TestRestorer_restore_existingFileStopsCopyCleanly(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)

	dir := t.TempDir()
	if _, err := file.OverWriteFile(ctx, file.Join(dir, "f.txt"), bytes.NewReader([]byte("old")), 0o600); err != nil {
		t.Fatal(err)
	}
	// The large second entry keeps the copy goroutine blocked on the pipe after
	// restore() has stopped at the first one.
	big := make([]byte, 256<<10)
	data := buildTar(t, []tarEntry{
		{name: "f.txt", content: []byte("new"), size: 3},
		{name: "g.bin", content: big, size: int64(len(big))},
	}, true)

	r := &restorer{
		eg:      eg,
		dir:     dir,
		storage: readerStorage{rc: &errAfterReader{data: data}},
	}
	if err := r.restore(egctx); err != nil {
		t.Fatalf("restore() error = %v, want nil", err)
	}

	waited := make(chan error, 1)
	go func() { waited <- eg.Wait() }()
	select {
	case err := <-waited:
		if err != nil {
			t.Errorf("eg.Wait() = %v, want nil: an early stop of restore is not a copy failure", err)
		}
	case <-ctx.Done():
		t.Fatalf("eg.Wait did not return: %v", ctx.Err())
	}
}
