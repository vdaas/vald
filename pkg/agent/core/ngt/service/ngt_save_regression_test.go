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
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
)

// Test_ngt_saveIndex_cancelledCtxFullSnapshot verifies that a save started
// with a cancelled context still persists the whole kvs. kvs.Range stops
// silently on cancellation, so the snapshot used to be partial (here: empty)
// while kvsdb was truncated and the metadata IndexCount matched the truncated
// count, hiding the loss from the load-time IndexCount/2 > kvs.Len() check.
// This is the path of the Start loop's SaveIndex on SIGTERM and of a SaveIndex
// RPC whose client cancels.
func Test_ngt_saveIndex_cancelledCtxFullSnapshot(t *testing.T) {
	t.Parallel()
	const num = 100

	svc, err := New(&defaultConfig, WithIndexPath(t.TempDir()))
	if err != nil {
		t.Fatalf("failed to create ngt service: %v", err)
	}
	n, ok := svc.(*ngt)
	if !ok {
		t.Fatalf("New returned %T, want *ngt", svc)
	}
	defer func() {
		if cerr := n.Close(context.WithoutCancel(t.Context())); cerr != nil {
			t.Errorf("Close() err = %v", cerr)
		}
	}()

	for i := range num {
		vec := make([]float32, defaultConfig.Dimension)
		for j := range vec {
			vec[j] = float32(i*defaultConfig.Dimension + j + 1)
		}
		if err = n.Insert("uuid-"+strconv.Itoa(i), vec); err != nil {
			t.Fatalf("Insert(%d) err = %v", i, err)
		}
	}
	if err = n.CreateIndex(t.Context(), 4); err != nil {
		t.Fatalf("CreateIndex() err = %v", err)
	}
	if got := n.kvs.Len(); got != num {
		t.Fatalf("setup: kvs.Len() = %d, want %d", got, num)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err = n.saveIndex(ctx); err != nil {
		t.Fatalf("saveIndex() err = %v, want nil", err)
	}

	readGob := func(name string, v any) {
		t.Helper()
		fi, oerr := file.Open(file.Join(n.path, name), os.O_RDONLY, os.ModePerm)
		if oerr != nil {
			t.Fatal(oerr)
		}
		if derr := gob.NewDecoder(fi).Decode(v); derr != nil {
			t.Fatal(errors.Join(derr, fi.Close()))
		}
		if cerr := fi.Close(); cerr != nil {
			t.Fatal(cerr)
		}
	}
	m := make(map[string]uint32)
	readGob(kvsFileName, &m)
	if len(m) != num {
		t.Errorf("kvsdb entries = %d, want %d", len(m), num)
	}
	mt := make(map[string]int64)
	readGob(kvsTimestampFileName, &mt)
	if len(mt) != num {
		t.Errorf("kvs timestamp entries = %d, want %d", len(mt), num)
	}
	meta, err := metadata.Load(file.Join(n.path, metadata.AgentMetadataFileName))
	if err != nil {
		t.Fatal(err)
	}
	if meta.NGT == nil || meta.NGT.IndexCount != num {
		t.Errorf("metadata = %+v, want NGT.IndexCount %d", meta.NGT, num)
	}
}
