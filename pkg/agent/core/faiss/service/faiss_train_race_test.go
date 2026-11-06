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
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

// blockingTrainCore is a saveOKCore whose Train blocks until released and
// whose Add succeeds, so a CreateIndex run can be held right before it updates
// isTrained, trainSize and the add buffer.
type blockingTrainCore struct {
	entered chan struct{}
	release chan struct{}
	saveOKCore
}

func (c *blockingTrainCore) Train(int, []float32) error {
	close(c.entered)
	select {
	case <-c.release:
		return nil
	case <-time.After(30 * time.Second):
		return errors.New("Train was not released")
	}
}

func (*blockingTrainCore) Add(n int, _ []float32, _ []int64) (int, error) { return n, nil }

// Test_faiss_trainState_concurrentReaders verifies that the training state and
// the add buffer size, which CreateIndex updates under cimu, can be read by
// GetTrainSize, SaveIndex and the lock-free pre-check of another CreateIndex
// while that run is in progress. Run with -race (make test): the readers are
// not ordered with the writes made after Train returns.
func Test_faiss_trainState_concurrentReaders(t *testing.T) {
	t.Parallel()
	vq, err := vqueue.New()
	if err != nil {
		t.Fatal(err)
	}
	c := &blockingTrainCore{entered: make(chan struct{}), release: make(chan struct{})}
	f := &faiss{
		dim:              2,
		nlist:            1,
		m:                1,
		path:             t.TempDir(),
		fmap:             make(map[string]int64),
		kvs:              kvs.New(kvs.WithConcurrency(1)),
		vq:               vq,
		core:             c,
		kvsdbConcurrency: 1,
	}
	defer closeKVS(t, f)
	f.indexing.Store(false)
	f.saving.Store(false)
	// max(nlist, 2^m)*minPointsPerCentroid vectors are needed to train.
	n := 2 * minPointsPerCentroid
	for i := range n {
		if err := f.Insert("uuid-"+strconv.Itoa(i), []float32{float32(i), float32(i)}); err != nil {
			t.Fatal(err)
		}
	}
	// lastNocie != nocie, so saveIndex reads isTrained instead of returning early.
	f.nocie.Store(1)

	done := make(chan error, 1)
	go func() { done <- f.CreateIndex(t.Context()) }()
	select {
	case <-c.entered:
	case err := <-done:
		t.Fatalf("CreateIndex returned before Train: %v", err)
	case <-time.After(30 * time.Second):
		t.Fatal("CreateIndex did not reach Train")
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	close(c.release)
	if got := f.GetTrainSize(); got != 0 && got != n {
		t.Errorf("GetTrainSize() = %d during CreateIndex, want 0 or %d", got, n)
	}
	if err := f.CreateIndex(ctx); err != nil &&
		!errors.Is(err, context.Canceled) && !errors.Is(err, errors.ErrUncommittedIndexNotFound) {
		t.Errorf("concurrent CreateIndex() err = %v", err)
	}
	if err := f.saveIndex(ctx); err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("concurrent saveIndex() err = %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("CreateIndex() err = %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("CreateIndex did not finish")
	}
	if !f.isTrained.Load() {
		t.Error("isTrained = false after training, want true")
	}
	if got := f.GetTrainSize(); got != n {
		t.Errorf("GetTrainSize() = %d, want %d", got, n)
	}
	if got := f.addLen.Load(); got != 0 || len(f.addVecs) != 0 {
		t.Errorf("add buffer: addLen = %d, %d floats, want empty", got, len(f.addVecs))
	}
}
