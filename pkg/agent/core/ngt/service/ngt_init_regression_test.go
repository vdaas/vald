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
	"encoding/gob"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
)

func writeGob(t *testing.T, path string, v any) {
	t.Helper()
	fi, err := file.Open(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	require.NoError(t, err)
	if err := gob.NewEncoder(fi).Encode(v); err != nil {
		t.Fatal(errors.Join(err, fi.Close()))
	}
	require.NoError(t, fi.Close())
}

func closeKVS(t *testing.T, n *ngt) {
	t.Helper()
	if n.kvs != nil {
		if err := n.kvs.Close(); err != nil {
			t.Error(err)
		}
	}
}

// Test_ngt_initNGT_loadTimeout verifies that initNGT returns on a load timeout
// instead of rebuilding the index while the abandoned load is still running.
// The timestamp kvsdb file is a FIFO with no writer, so loadKVS blocks in open
// and load deterministically times out.
func Test_ngt_initNGT_loadTimeout(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	require.NoError(t, metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		NGT: &metadata.NGT{IndexCount: 1},
	}))
	writeGob(t, file.Join(dir, kvsFileName), map[string]uint32{"a": 1})
	fifo := file.Join(dir, kvsTimestampFileName)
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))

	eg, _ := errgroup.New(t.Context())
	n := &ngt{
		path:             dir,
		fmap:             make(map[string]int64),
		kvsdbConcurrency: 1,
		eg:               eg,
		minLit:           100 * time.Millisecond,
		maxLit:           100 * time.Millisecond,
		litFactor:        time.Millisecond,
	}

	err := n.initNGT(core.WithDimension(2))
	if !errors.Is(err, errors.ErrIndexLoadTimeout) {
		t.Errorf("initNGT() err = %v, want ErrIndexLoadTimeout", err)
	}

	// Release the abandoned loadKVS: an O_RDWR open of a FIFO does not block on
	// Linux and serves the reader whether it is already blocked in open or
	// opens later; the fd stays open until the abandoned load has finished.
	w, err := file.Open(fifo, os.O_RDWR, 0)
	if err != nil {
		// a fallback removed the FIFO, so the abandoned load cannot be released.
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
		closeKVS(t, n)
		if n.core != nil {
			n.core.Close()
		}
	case <-time.After(30 * time.Second):
		t.Error("abandoned load did not finish")
	}
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

// saveTestIndex stores an NGT index of one vector at dir together with its
// metadata (claiming indexCount entries) and a kvsdb with the given uuids.
func saveTestIndex(t *testing.T, dir string, indexCount uint64, uuids map[string]uint32) {
	t.Helper()
	require.NoError(t, file.MkdirAll(dir, os.ModePerm))
	c, err := core.New(core.WithIndexPath(dir), core.WithDimension(2))
	require.NoError(t, err)
	_, err = c.Insert([]float32{1, 2})
	require.NoError(t, errors.Join(err, c.CreateIndex(1), c.SaveIndex()))
	c.Close()
	require.NoError(t, metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		NGT: &metadata.NGT{IndexCount: indexCount},
	}))
	ts := make(map[string]int64, len(uuids))
	for k := range uuids {
		ts[k] = 1
	}
	writeGob(t, file.Join(dir, kvsFileName), uuids)
	writeGob(t, file.Join(dir, kvsTimestampFileName), ts)
}

// Test_ngt_initNGT_reloadsTooFewPrimary verifies that, in CoW mode, when the
// primary kvsdb is too small for its metadata and the secondary has even fewer
// entries, initNGT starts with the reloaded primary. The reload reports the
// same ErrIndicesAreTooFewComparedToMetadata as the first primary load; treating
// it as a failure fell through to a fresh empty index and kvsdb.
func Test_ngt_initNGT_reloadsTooFewPrimary(t *testing.T) {
	t.Parallel()
	base := t.TempDir()
	eg, _ := errgroup.New(t.Context())
	n := &ngt{
		basePath:          base,
		path:              file.Join(base, originIndexDirName),
		oldPath:           file.Join(base, oldIndexDirName),
		brokenPath:        file.Join(base, brokenIndexDirName),
		fmap:              make(map[string]int64),
		kvsdbConcurrency:  1,
		eg:                eg,
		minLit:            time.Minute,
		maxLit:            time.Minute,
		litFactor:         time.Millisecond,
		enableCopyOnWrite: true,
	}
	n.tmpPath.Store(t.TempDir())
	// primary: 2 entries for 10 indices (too few); secondary: 1 entry, consistent.
	saveTestIndex(t, n.path, 10, map[string]uint32{"a": 1, "b": 2})
	saveTestIndex(t, n.oldPath, 1, map[string]uint32{"a": 1})

	require.NoError(t, n.initNGT(core.WithDimension(2)))
	defer closeKVS(t, n)
	require.NotNil(t, n.core)
	defer n.core.Close()
	require.Equal(t, uint64(2), n.kvs.Len(), "the reloaded primary kvsdb must be kept")
	for _, uuid := range []string{"a", "b"} {
		_, _, ok := n.kvs.Get(uuid)
		require.True(t, ok, "kvs entry of %s is missing", uuid)
	}
}

// Test_ngt_initNGT_keepsTooFewPrimaryWithoutSecondary verifies that, in CoW
// mode, when the primary kvsdb is too small for its metadata and neither the
// secondary nor the non-CoW index can be loaded, initNGT starts with the
// reloaded primary instead of discarding it for a fresh empty index.
func Test_ngt_initNGT_keepsTooFewPrimaryWithoutSecondary(t *testing.T) {
	t.Parallel()
	base := t.TempDir()
	eg, _ := errgroup.New(t.Context())
	n := &ngt{
		basePath:          base,
		path:              file.Join(base, originIndexDirName),
		oldPath:           file.Join(base, oldIndexDirName),
		brokenPath:        file.Join(base, brokenIndexDirName),
		fmap:              make(map[string]int64),
		kvsdbConcurrency:  1,
		eg:                eg,
		minLit:            time.Minute,
		maxLit:            time.Minute,
		litFactor:         time.Millisecond,
		enableCopyOnWrite: true,
	}
	n.tmpPath.Store(t.TempDir())
	// primary: 2 entries for 10 indices (too few); no secondary, and base holds
	// no metadata, so both fallback loads fail.
	saveTestIndex(t, n.path, 10, map[string]uint32{"a": 1, "b": 2})

	require.NoError(t, n.initNGT(core.WithDimension(2)))
	defer closeKVS(t, n)
	require.NotNil(t, n.core)
	defer n.core.Close()
	require.Equal(t, uint64(2), n.kvs.Len(), "the reloaded primary kvsdb must be kept")
	for _, uuid := range []string{"a", "b"} {
		_, _, ok := n.kvs.Get(uuid)
		require.True(t, ok, "kvs entry of %s is missing", uuid)
	}
}

// closeCountCore counts Close calls on the wrapped core.
type closeCountCore struct {
	core.NGT
	closes atomic.Int32
}

func (c *closeCountCore) Close() {
	c.closes.Add(1)
	c.NGT.Close()
}

// Test_ngt_RegenerateIndexes_closesOldCore verifies that Flush frees the
// replaced core instead of leaking the whole C-side index on every call.
func Test_ngt_RegenerateIndexes_closesOldCore(t *testing.T) {
	t.Parallel()
	cfg := defaultConfig
	cfg.Dimension = 3
	nn, err := New(&cfg)
	require.NoError(t, err)
	n, ok := nn.(*ngt)
	require.True(t, ok)
	old := &closeCountCore{NGT: n.core}
	n.core = old

	require.NoError(t, n.RegenerateIndexes(t.Context()))
	require.Equal(t, int32(1), old.closes.Load(), "the replaced core must be closed once")
	require.NotSame(t, old, n.core)
	require.NotNil(t, n.core)
	require.NoError(t, n.Close(t.Context()))
	require.Equal(t, int32(1), old.closes.Load(), "Close must not close the replaced core again")
}
