//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package service

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
)

var cfg = &config.NGT{
	Dimension:              2,
	DistanceType:           "l2",
	ObjectType:             "float",
	EnableInMemoryMode:     true,
	AutoIndexDurationLimit: "96h",
	AutoIndexCheckDuration: "96h",
	AutoSaveIndexDuration:  "96h",
	AutoIndexLength:        10000000000,
	KVSDB: &config.KVSDB{
		Concurrency: 1,
	},
}

const (
	maxIDNum                   = 10
	duplicateIDNum             = 1000
	maxCreateIndexNum          = 5
	createIndexPoolSize uint32 = 10000
)

func Test_ngt_parallel_delete_and_insert(t *testing.T) {
	n, err := New(cfg.Bind())
	if err != nil {
		t.Fatalf("failed to create ngt service: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n.Start(ctx)
	time.Sleep(10 * time.Millisecond)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	for dup := 0; dup < duplicateIDNum; dup++ {
		for i := int64(0); i < maxIDNum; i++ {
			i := i
			wg.Add(1)
			go func() {
				mu.Lock()
				defer mu.Unlock()
				defer wg.Done()
				c.Wait()

				uuid := strconv.FormatInt(i, 10)

				// NOTE: shoudl we check error?
				_ = n.Insert(uuid, []float32{float32(i), float32(i)})
				_ = n.Delete(uuid)
			}()
		}
	}

	wg.Add(1)
	go func() {
		mu.Lock()
		defer mu.Unlock()
		defer wg.Done()
		c.Wait()

		tic := time.NewTicker(time.Second)
		defer tic.Stop()

		for i := 0; i < maxCreateIndexNum; i++ {
			select {
			case <-tic.C:
				// NOTE: shoudl we check error?
				n.CreateIndex(ctx, createIndexPoolSize)
			}
		}
	}()

	time.Sleep(1 * time.Second)
	c.Broadcast()
	wg.Wait()

	for i := int64(0); i < maxIDNum; i++ {
		uuid := strconv.FormatInt(i, 10)
		if err := n.Insert(uuid, []float32{float32(i), float32(i)}); err != nil {
			t.Errorf("Insert error: %v", err)
		}
	}
}
