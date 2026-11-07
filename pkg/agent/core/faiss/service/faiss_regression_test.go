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
	"slices"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

func writeGob(t *testing.T, path string, v any) {
	t.Helper()
	fi, err := file.Open(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if err := gob.NewEncoder(fi).Encode(v); err != nil {
		t.Fatal(errors.Join(err, fi.Close()))
	}
	if err := fi.Close(); err != nil {
		t.Fatal(err)
	}
}

func closeKVS(t *testing.T, f *faiss) {
	t.Helper()
	if f.kvs != nil {
		if err := f.kvs.Close(); err != nil {
			t.Error(err)
		}
	}
}

// Test_faiss_loadKVS_restoresIcnt verifies that the id counter resumes after
// the largest loaded id, so a restarted agent does not reuse ids that already
// label vectors of the loaded index.
func Test_faiss_loadKVS_restoresIcnt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ids  map[string]uint32
		name string
		icnt uint64
		want uint64
	}{
		{
			name: "next id follows the largest loaded id",
			ids:  map[string]uint32{"a": 0, "b": 5, "c": 2},
			want: 6,
		},
		{
			name: "stale counter is replaced by the loaded state",
			ids:  map[string]uint32{"a": 3},
			icnt: 100,
			want: 4,
		},
		{
			name: "empty kvsdb starts from zero",
			ids:  map[string]uint32{},
			want: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			ts := make(map[string]int64, len(test.ids))
			for k := range test.ids {
				ts[k] = 1
			}
			writeGob(t, file.Join(dir, kvsFileName), test.ids)
			writeGob(t, file.Join(dir, kvsTimestampFileName), ts)

			f := &faiss{fmap: make(map[string]int64), kvsdbConcurrency: 1, icnt: test.icnt}
			defer closeKVS(t, f)
			if err := f.loadKVS(t.Context(), dir, time.Minute); err != nil {
				t.Fatal(err)
			}
			if f.icnt != test.want {
				t.Errorf("icnt = %d, want %d", f.icnt, test.want)
			}
			if got := f.kvs.Len(); got != uint64(len(test.ids)) {
				t.Errorf("kvs.Len() = %d, want %d", got, len(test.ids))
			}
		})
	}
}

// Test_faiss_saveIndex_rollbackOnError verifies that a failed save does not
// leave the generation marked as saved, which would make every later
// SaveIndex a no-op until the next CreateIndex.
func Test_faiss_saveIndex_rollbackOnError(t *testing.T) {
	t.Parallel()
	f := new(faiss)
	f.isTrained.Store(true)
	f.nocie.Store(1)
	f.indexing.Store(true) // keeps saveIndex in its wait loop until ctx is done
	f.saving.Store(false)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	for i := range 2 {
		if err := f.saveIndex(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("attempt %d: err = %v, want context.Canceled", i, err)
		}
		if got := f.lastNocie.Load(); got != 0 {
			t.Fatalf("attempt %d: lastNocie = %d, want 0 (rolled back)", i, got)
		}
	}
}

// stubCore is a core.Faiss that records Remove calls and fails everything else
// the tests do not expect to reach the C library.
type stubCore struct {
	removed [][]int64
}

func (*stubCore) SaveIndex() error               { return errors.New("unexpected SaveIndex") }
func (*stubCore) SaveIndexWithPath(string) error { return errors.New("unexpected SaveIndexWithPath") }
func (*stubCore) Train(int, []float32) error     { return errors.New("unexpected Train") }
func (*stubCore) Add(int, []float32, []int64) (int, error) {
	return -1, errors.New("unexpected Add")
}

func (*stubCore) Search(int, int, int, []float32) ([]algorithm.SearchResult, error) {
	return nil, errors.New("unexpected Search")
}

func (s *stubCore) Remove(size int, ids []int64) (int, error) {
	s.removed = append(s.removed, slices.Clone(ids[:size]))
	return 0, nil
}
func (*stubCore) Close() {}

// Test_faiss_CreateIndex_deleteBuffered verifies that deleting a vector that
// is still in the untrained add buffer removes it from that buffer, keeping
// addIds and addVecs aligned, instead of later adding it to the index as an
// orphan label.
func Test_faiss_CreateIndex_deleteBuffered(t *testing.T) {
	t.Parallel()
	vq, err := vqueue.New()
	if err != nil {
		t.Fatal(err)
	}
	sc := new(stubCore)
	// max(nlist, 2^m)*minPointsPerCentroid = 78 buffered vectors are needed to
	// train, so the three vectors below stay buffered.
	f := &faiss{
		dim:              2,
		nlist:            1,
		m:                1,
		fmap:             make(map[string]int64),
		kvs:              kvs.New(kvs.WithConcurrency(1)),
		vq:               vq,
		core:             sc,
		kvsdbConcurrency: 1,
	}
	defer closeKVS(t, f)
	f.indexing.Store(false)
	f.saving.Store(false)

	vecs := map[string][]float32{
		"a": {1, 1},
		"b": {2, 2},
		"c": {3, 3},
	}
	for _, uuid := range []string{"a", "b", "c"} {
		if err := f.Insert(uuid, vecs[uuid]); err != nil {
			t.Fatal(err)
		}
	}
	ctx := t.Context()
	if err := f.CreateIndex(ctx); err != nil {
		t.Fatal(err)
	}
	if f.isTrained.Load() || len(f.addIds) != 3 {
		t.Fatalf("setup: isTrained = %v, buffered = %d, want false and 3", f.isTrained.Load(), len(f.addIds))
	}

	if err := f.Delete("a"); err != nil {
		t.Fatal(err)
	}
	if err := f.CreateIndex(ctx); err != nil {
		t.Fatal(err)
	}

	if len(sc.removed) != 0 {
		t.Errorf("core.Remove called with %v for a buffered vector", sc.removed)
	}
	if len(f.addIds) != 2 || len(f.addVecs) != 2*f.dim {
		t.Fatalf("buffer: %d ids, %d floats, want 2 ids, %d floats", len(f.addIds), len(f.addVecs), 2*f.dim)
	}
	if got := f.addLen.Load(); got != 2 {
		t.Errorf("addLen = %d, want 2", got)
	}
	// Both remaining slots must belong to the live uuids b and c, each paired
	// with its own vector; with only two slots this also rules out an orphan.
	for _, uuid := range []string{"b", "c"} {
		oid, _, ok := f.kvs.Get(uuid)
		if !ok {
			t.Fatalf("kvs entry of %s is missing", uuid)
		}
		i := slices.Index(f.addIds, int64(oid))
		if i < 0 {
			t.Errorf("%s (oid %d) is not buffered: %v", uuid, oid, f.addIds)
			continue
		}
		if got := f.addVecs[i*f.dim : (i+1)*f.dim]; !slices.Equal(got, vecs[uuid]) {
			t.Errorf("buffered vector of %s = %v, want %v", uuid, got, vecs[uuid])
		}
	}
}

// Test_faiss_initFaiss_loadTimeout verifies that initFaiss returns on a load
// timeout instead of falling back while the abandoned load is still running.
// The timestamp kvsdb file is a FIFO with no writer, so loadKVS blocks in open
// and load deterministically times out.
func Test_faiss_initFaiss_loadTimeout(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		Faiss: &metadata.Faiss{IndexCount: 1},
	}); err != nil {
		t.Fatal(err)
	}
	writeGob(t, file.Join(dir, kvsFileName), map[string]uint32{"a": 0})
	fifo := file.Join(dir, kvsTimestampFileName)
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	// load() requires a non-empty faiss.index before it starts loading. A
	// corrupt one makes the index load fail fast, so only the kvsdb load
	// (blocked on the FIFO) keeps the concurrent load running past its timeout.
	if err := os.WriteFile(file.Join(dir, faissIndexFileName), []byte("corrupt"), 0o600); err != nil {
		t.Fatal(err)
	}

	eg, _ := errgroup.New(t.Context())
	f := &faiss{
		path:             dir,
		fmap:             make(map[string]int64),
		kvsdbConcurrency: 1,
		eg:               eg,
		minLit:           100 * time.Millisecond,
		maxLit:           100 * time.Millisecond,
		litFactor:        time.Millisecond,
	}

	err := f.initFaiss()
	if !errors.Is(err, errors.ErrIndexLoadTimeout) {
		t.Errorf("initFaiss() err = %v, want ErrIndexLoadTimeout", err)
	}

	// Release the abandoned loadKVS: an O_RDWR open of a FIFO does not block on
	// Linux and serves the reader whether it is already blocked in open or
	// opens later; the fd stays open until the abandoned load has finished.
	w, err := file.Open(fifo, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := gob.NewEncoder(w).Encode(map[string]int64{}); err != nil {
		t.Error(err)
	}
	done := make(chan error, 1)
	go func() { done <- eg.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Error(err)
		}
		closeKVS(t, f)
	case <-time.After(30 * time.Second):
		t.Error("abandoned load did not finish")
	}
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

// Test_faiss_resetTrainState verifies that the state restored by a successful
// load is cleared before initFaiss falls back to a fresh core; otherwise
// CreateIndex would skip Train and Add vectors to an untrained index.
func Test_faiss_resetTrainState(t *testing.T) {
	t.Parallel()
	f := &faiss{icnt: 7}
	f.isTrained.Store(true)
	f.resetTrainState()
	if f.isTrained.Load() {
		t.Error("isTrained = true after reset, want false")
	}
	if f.icnt != 0 {
		t.Errorf("icnt = %d after reset, want 0", f.icnt)
	}
}

// Test_faiss_waitLoad verifies that a successful load always marks the index
// as trained (the load goroutine must not cancel ctx before reporting, which
// let ctx.Done() win the select and skip the success path), and that a timed
// out load keeps ErrIndexLoadTimeout even when storing the invalid metadata
// fails, so initFaiss does not start another load beside the abandoned one.
func Test_faiss_waitLoad(t *testing.T) {
	t.Parallel()
	errLoad := errors.New("load failed")
	tests := []struct {
		wantErr     error
		wait        func(release <-chan struct{}) error
		name        string
		timeout     time.Duration
		badMetaPath bool
		wantTrained bool
		wantInvalid bool
	}{
		{
			name:        "success marks the index as trained",
			timeout:     time.Minute,
			wait:        func(<-chan struct{}) error { return nil },
			wantTrained: true,
		},
		{
			name:    "load error is returned and the index stays untrained",
			timeout: time.Minute,
			wait:    func(<-chan struct{}) error { return errLoad },
			wantErr: errLoad,
		},
		{
			name:    "timeout stores invalid metadata and returns ErrIndexLoadTimeout",
			timeout: 10 * time.Millisecond,
			wait: func(release <-chan struct{}) error {
				<-release
				return nil
			},
			wantErr:     errors.ErrIndexLoadTimeout,
			wantInvalid: true,
		},
		{
			name:    "timeout keeps ErrIndexLoadTimeout when the metadata store fails",
			timeout: 10 * time.Millisecond,
			wait: func(release <-chan struct{}) error {
				<-release
				return nil
			},
			badMetaPath: true,
			wantErr:     errors.ErrIndexLoadTimeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			metaPath := file.Join(dir, metadata.AgentMetadataFileName)
			if tc.badMetaPath {
				// a path below a regular file fails with ENOTDIR, also as root.
				notDir := file.Join(dir, "not-a-dir")
				writeGob(t, notDir, map[string]uint32{})
				metaPath = file.Join(notDir, metadata.AgentMetadataFileName)
			}
			// the success path is repeated so that a select woken by a
			// concurrent ctx cancellation would be observed.
			runs := 1
			if tc.wantTrained {
				runs = 100
			}
			for range runs {
				eg, _ := errgroup.New(t.Context())
				f := &faiss{eg: eg}
				release := make(chan struct{})
				ctx, cancel := context.WithTimeout(t.Context(), tc.timeout)
				err := f.waitLoad(ctx, func() error { return tc.wait(release) }, dir, metaPath, tc.timeout)
				cancel()
				close(release)
				if tc.wantErr == nil && err != nil || tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
					t.Fatalf("waitLoad() err = %v, want %v", err, tc.wantErr)
				}
				if got := f.isTrained.Load(); got != tc.wantTrained {
					t.Fatalf("isTrained = %v, want %v", got, tc.wantTrained)
				}
				done := make(chan error, 1)
				go func() { done <- eg.Wait() }()
				select {
				case err := <-done:
					if err != nil {
						t.Error(err)
					}
				case <-time.After(30 * time.Second):
					t.Fatal("load goroutine did not finish")
				}
			}
			if tc.wantInvalid {
				meta, err := metadata.Load(metaPath)
				if err != nil {
					t.Fatal(err)
				}
				if !meta.IsInvalid {
					t.Error("metadata IsInvalid = false, want true")
				}
			}
		})
	}
}

// saveOKCore is a stubCore whose SaveIndexWithPath succeeds, so saveIndex
// reaches the metadata write.
type saveOKCore struct {
	stubCore
}

func (*saveOKCore) SaveIndexWithPath(string) error { return nil }

// Test_faiss_saveIndex_cancelledCtxFullSnapshot verifies that a save started
// with a cancelled context still persists the whole kvs. kvs.Range stops
// silently on cancellation, so the snapshot used to be partial (here: empty)
// while kvsdb was truncated and the metadata IndexCount matched the truncated
// count, hiding the loss from the load-time consistency check.
func Test_faiss_saveIndex_cancelledCtxFullSnapshot(t *testing.T) {
	t.Parallel()
	const n = 100
	dir := t.TempDir()
	f := &faiss{
		path: dir,
		fmap: make(map[string]int64),
		kvs:  kvs.New(kvs.WithConcurrency(1)),
		core: new(saveOKCore),
	}
	defer closeKVS(t, f)
	f.isTrained.Store(true)
	f.indexing.Store(false)
	f.saving.Store(false)
	f.nocie.Store(1)
	for i := range n {
		f.kvs.Set("uuid-"+strconv.Itoa(i), uint32(i), int64(i+1))
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := f.saveIndex(ctx); err != nil {
		t.Fatalf("saveIndex() err = %v, want nil", err)
	}

	readGob := func(name string, v any) {
		t.Helper()
		fi, err := file.Open(file.Join(dir, name), os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		if err := gob.NewDecoder(fi).Decode(v); err != nil {
			t.Fatal(errors.Join(err, fi.Close()))
		}
		if err := fi.Close(); err != nil {
			t.Fatal(err)
		}
	}
	m := make(map[string]uint32)
	readGob(kvsFileName, &m)
	if len(m) != n {
		t.Errorf("kvsdb entries = %d, want %d", len(m), n)
	}
	mt := make(map[string]int64)
	readGob(kvsTimestampFileName, &mt)
	if len(mt) != n {
		t.Errorf("kvs timestamp entries = %d, want %d", len(mt), n)
	}
	meta, err := metadata.Load(file.Join(dir, metadata.AgentMetadataFileName))
	if err != nil {
		t.Fatal(err)
	}
	if meta.Faiss == nil || meta.Faiss.IndexCount != n {
		t.Errorf("metadata = %+v, want Faiss.IndexCount %d", meta.Faiss, n)
	}
}

// pathRecordingCore is a stubCore whose SaveIndexWithPath writes a stand-in
// index file and records its path, failing like the faiss writer when the path
// is a directory or its parent directory is missing.
type pathRecordingCore struct {
	stubCore
	saved []string
}

func (c *pathRecordingCore) SaveIndexWithPath(path string) error {
	c.saved = append(c.saved, path)
	// os.OpenFile (unlike file.Open) does not create missing parent
	// directories, matching faiss write_index.
	fi, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fi.WriteString("index")
	return errors.Join(err, fi.Close())
}

// Test_faiss_saveIndex_indexFile verifies that the index is saved as a file
// inside the index directory and that the directory exists before the index
// writer runs. The directory itself used to be passed to faiss write_index,
// which cannot open a directory, so no index was ever persisted.
func Test_faiss_saveIndex_indexFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		kvs  map[string]uint32
		name string
	}{
		{
			name: "empty kvs creates the missing index directory",
		},
		{
			name: "kvsdb and index file are written to the same directory",
			kvs:  map[string]uint32{"a": 0, "b": 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dir := file.Join(t.TempDir(), "index")
			c := new(pathRecordingCore)
			f := &faiss{
				path: dir,
				fmap: make(map[string]int64),
				kvs:  kvs.New(kvs.WithConcurrency(1)),
				core: c,
			}
			defer closeKVS(t, f)
			f.isTrained.Store(true)
			f.indexing.Store(false)
			f.saving.Store(false)
			f.nocie.Store(1)
			for k, id := range test.kvs {
				f.kvs.Set(k, id, 1)
			}

			if err := f.saveIndex(t.Context()); err != nil {
				t.Fatalf("saveIndex() err = %v, want nil", err)
			}
			want := file.Join(dir, faissIndexFileName)
			if !slices.Equal(c.saved, []string{want}) {
				t.Errorf("SaveIndexWithPath paths = %v, want [%s]", c.saved, want)
			}
			if exist, fi, err := file.ExistsWithDetail(want); !exist || err != nil || fi == nil || fi.IsDir() {
				t.Errorf("index file %s: exist = %v, info = %v, err = %v, want a regular file", want, exist, fi, err)
			}
			if _, err := metadata.Load(file.Join(dir, metadata.AgentMetadataFileName)); err != nil {
				t.Errorf("metadata.Load() err = %v, want nil", err)
			}
		})
	}
}

// Test_faiss_load_missingIndexFile verifies that a directory with valid
// metadata and kvsdb but without the index file is reported as
// ErrIndexFileNotFound, so initFaiss takes its usual missing-index fallback.
// The directory used to be handed to faiss read_index instead.
func Test_faiss_load_missingIndexFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		Faiss: &metadata.Faiss{IndexCount: 1},
	}); err != nil {
		t.Fatal(err)
	}
	writeGob(t, file.Join(dir, kvsFileName), map[string]uint32{"a": 0})
	writeGob(t, file.Join(dir, kvsTimestampFileName), map[string]int64{"a": 1})

	eg, _ := errgroup.New(t.Context())
	f := &faiss{
		path:             dir,
		fmap:             make(map[string]int64),
		kvsdbConcurrency: 1,
		eg:               eg,
		minLit:           time.Minute,
		maxLit:           time.Minute,
		litFactor:        time.Millisecond,
	}
	defer closeKVS(t, f)

	if err := f.load(t.Context(), dir); !errors.Is(err, errors.ErrIndexFileNotFound) {
		t.Errorf("load() err = %v, want ErrIndexFileNotFound", err)
	}
}

// Test_moveAsideLegacyIndexFile verifies that a regular file left at the index
// path by an older agent, which wrote the faiss index to the path itself, is
// renamed next to the path so the index directory can be created and saved
// again; the file must be kept. Before, MkdirAll failed on it on every save.
func Test_moveAsideLegacyIndexFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		setup      func(t *testing.T, path string)
		name       string
		wantBackup bool
	}{
		{
			name:  "missing path is left alone",
			setup: func(*testing.T, string) {},
		},
		{
			name: "existing directory is left alone",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if err := file.MkdirAll(path, os.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "legacy index file is moved aside and saving succeeds",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if err := os.WriteFile(path, []byte("legacy"), os.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			wantBackup: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parent := t.TempDir()
			path := file.Join(parent, "index")
			test.setup(t, path)

			if err := moveAsideLegacyIndexFile(path); err != nil {
				t.Fatalf("moveAsideLegacyIndexFile() err = %v, want nil", err)
			}

			entries, err := file.ReadDir(parent)
			if err != nil {
				t.Fatal(err)
			}
			var backups []string
			for _, e := range entries {
				if strings.HasPrefix(e.Name(), "index.legacy-") {
					backups = append(backups, e.Name())
				}
			}
			if !test.wantBackup {
				if len(backups) != 0 {
					t.Errorf("backups = %v, want none", backups)
				}
				return
			}
			if len(backups) != 1 {
				t.Fatalf("backups = %v, want exactly one", backups)
			}
			b, err := file.ReadFile(file.Join(parent, backups[0]))
			if err != nil || string(b) != "legacy" {
				t.Errorf("backup content = %q, err = %v, want %q", b, err, "legacy")
			}

			c := new(pathRecordingCore)
			f := &faiss{
				path: path,
				fmap: make(map[string]int64),
				kvs:  kvs.New(kvs.WithConcurrency(1)),
				core: c,
			}
			defer closeKVS(t, f)
			f.isTrained.Store(true)
			f.indexing.Store(false)
			f.saving.Store(false)
			f.nocie.Store(1)
			if err := f.saveIndex(t.Context()); err != nil {
				t.Errorf("saveIndex() after moving the legacy file aside err = %v, want nil", err)
			}
		})
	}
}
