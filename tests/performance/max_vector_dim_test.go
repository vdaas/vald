//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package performance

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

const (
	maxBit       = 32
	freeMemLimit = 500 // Limit of free memory remaining(MB)
)

func init_ngt_service(dim int) (service.NGT, error) {
	cfg := &config.NGT{
		Dimension:              dim,
		DistanceType:           ngt.L2.String(),
		ObjectType:             ngt.Float.String(),
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: "96h",
		AutoIndexCheckDuration: "96h",
		AutoSaveIndexDuration:  "96h",
		AutoIndexLength:        10000000000,
		KVSDB: &config.KVSDB{
			Concurrency: 1,
		},
		VQueue: &config.VQueue{
			InsertBufferPoolSize: 100,
			DeleteBufferPoolSize: 100,
		},
	}
	ngt, err := service.New(cfg.Bind())
	if err != nil {
		return nil, err
	}
	return ngt, nil
}

func parse(raw string) (key string, value int) {
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	keyValue := strings.Split(text, ":")
	val := 0
	if keyValue[1] != "" {
		val, err := strconv.Atoi(keyValue[1])
		if err != nil {
			panic(err)
		}
		return keyValue[0], val
	}
	return keyValue[0], val
}

// Test for investigation of max dimension size for agent handler
func TestMaxDimInsert(t *testing.T) {
	t.Helper()
	eg, ctx := errgroup.New(context.Background())
	mu := sync.Mutex{}
	// Get the above the limit of bit (2~32)
	bits := make([]int, 0, maxBit-1)
	ticker := time.NewTicker(5 * time.Second)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				if ok := mu.TryLock(); !ok {
					mu.Unlock()
				}
				return nil
			case <-ticker.C:
				f, err := os.Open("/proc/meminfo")
				if err != nil {
					panic(err)
				}
				bufio.NewScanner(f)
				scanner := bufio.NewScanner(f)
				var MemFree int
				for scanner.Scan() {
					key, value := parse(scanner.Text())
					switch key {
					case "MemFree":
						MemFree = value
					}
				}
				err = f.Close()
				if err != nil {
					panic(err)
				}
				if MemFree/1024 < freeMemLimit {
					t.Logf("MemFree reaches the limit: current : %dMb, limit : %dMb", MemFree/1024, freeMemLimit)
					return errors.New("Memory Limit")
				}
			}
		}
	})
	eg.Go(func() error {
		for bit := 2; bit <= maxBit; bit++ {
			select {
			case <-ctx.Done():
				t.Log("canceld")
				return nil
			default:
				dim := 1 << bit
				if bit == maxBit {
					dim--
				}
				if dim > ngt.VectorDimensionSizeLimit {
					t.Fatal(errors.ErrInvalidDimensionSize(dim, ngt.VectorDimensionSizeLimit))
				}
				t.Logf("Start test: dimension = %d (bit = %d)", dim, bit)
				ngt, err := init_ngt_service(dim)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					t.Errorf("[Fail] Create NGT service: %#v", err)
					return err
				}
				vec := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
				err = ngt.Insert(strconv.Itoa(dim), vec)
				if err != nil {
					t.Errorf("Insert error: %#v", err)
					return err
				}
				t.Logf("Insert is finished: dimension = %d (bit = %d)", dim, bit)
				err = ngt.CreateIndex(ctx, 10)
				if err != nil {
					t.Errorf("CreateIndex is failed: %#v", err)
					return err
				}
				t.Logf("CreateIndex is finished: dimension = %d (bit = %d)", dim, bit)
				err = ngt.Close(ctx)
				if err != nil {
					t.Errorf("NGT close error: %#v", err)
					return err
				}
				mu.Lock()
				bits = append(bits, bit)
				mu.Unlock()
				t.Logf("All processes are finished: dimension = %d (bit = %d)", dim, bit)
			}
			t.Log("Wait for memory release")
			time.Sleep(30 * time.Second)
		}
		return nil
	})
	eg.Wait()
	// Get the max bit, which the environment finish process, from bits
	var max_bit int
	for _, v := range bits {
		if max_bit < v {
			max_bit = v
		}
	}
	t.Logf("Max bit is %d", max_bit)
}

// Test for investigation of max dimension size for agent handler with gRPC
func TestMaxDimInsertGRPC(t *testing.T) {
	// MaxUint64 cannot be used due to overflows
	t.Helper()
	eg, ctx := errgroup.New(context.Background())
	mu := sync.Mutex{}
	// Get the above the limit of bit (2~32)
	bits := make([]int, 0, maxBit-1)
	ticker := time.NewTicker(5 * time.Second)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				if ok := mu.TryLock(); !ok {
					mu.Unlock()
				}
				return nil
			case <-ticker.C:
				f, err := os.Open("/proc/meminfo")
				if err != nil {
					panic(err)
				}
				bufio.NewScanner(f)
				scanner := bufio.NewScanner(f)
				var MemFree int
				for scanner.Scan() {
					key, value := parse(scanner.Text())
					switch key {
					case "MemFree":
						MemFree = value
					}
				}
				err = f.Close()
				if err != nil {
					panic(err)
				}
				if MemFree/1024 < freeMemLimit {
					t.Logf("MemFree reaches the limit: current : %dMb, limit : %dMb", MemFree/1024, freeMemLimit)
					return errors.New("Memory Limit")
				}
			}
		}
	})
	eg.Go(func() error {
		for bit := 2; bit <= maxBit; bit++ {
			select {
			case <-ctx.Done():
				t.Log("canceld")
				return nil
			default:
				dim := 1 << bit
				if bit == maxBit {
					dim--
				}
				if dim > ngt.VectorDimensionSizeLimit {
					t.Fatal(errors.ErrInvalidDimensionSize(dim, ngt.VectorDimensionSizeLimit))
				}
				t.Logf("Start test: dimension = %d (bit = %d)", dim, bit)
				ngt, err := init_ngt_service(dim)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					t.Errorf("[Fail] Create NGT service: %#v", err)
					return err
				}
				s, err := grpc.New(grpc.WithNGT(ngt))
				if err != nil {
					t.Errorf("[Error] Failed to create grpc service: %#v", s)
					return err
				}
				vec := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
				req := &payload.Insert_Request{
					Vector: &payload.Object_Vector{
						Id:     strconv.Itoa(dim),
						Vector: vec,
					},
					Config: &payload.Insert_Config{
						SkipStrictExistCheck: false,
					},
				}
				_, err = s.Insert(ctx, req)
				if err != nil {
					t.Errorf("[Error] Failed to Insert Vector (Dim = %d): %#v", dim, err)
					return err
				}
				t.Logf("Insert is finished: dimension = %d (bit = %d)", dim, bit)
				_, err = s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: 10,
				})
				if err != nil {
					t.Errorf("CreateIndex is failed: %#v", err)
					return err
				}
				t.Logf("CreateIndex is finished: dimension = %d (bit = %d)", dim, bit)
				err = ngt.Close(ctx)
				if err != nil {
					return err
				}
				mu.Lock()
				bits = append(bits, bit)
				mu.Unlock()
				t.Logf("All processes are finished: dimension = %d (bit = %d)", dim, bit)
			}
			t.Log("Wait for memory release")
			time.Sleep(30 * time.Second)
		}
		return nil
	})
	eg.Wait()
	// Get the max bit, which the environment finish process, from bits
	var max_bit int
	for _, v := range bits {
		if max_bit < v {
			max_bit = v
		}
	}
	t.Logf("Max bit is %d", max_bit)
}
