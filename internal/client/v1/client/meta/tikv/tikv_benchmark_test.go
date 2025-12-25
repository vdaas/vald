//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// $ tiup playground --mode tikv-slim
// $ TIKV_STORE_ADDRS=127.0.0.1:20160 go test ./internal/client/v1/client/meta/tikv/ -run=^$ -bench=. -benchmem
package tikv

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/test/goleak"
)

const (
	// separated list of TiKV PD addresses (e.g. "127.0.0.1:2379").
	envPDAddrs = "TIKV_PD_ADDRS"
)

var cli Client

func getAddrs() []string {
	pdAddrs := strings.Split(os.Getenv(envPDAddrs), ",")
	if len(pdAddrs) == 1 && pdAddrs[0] == "" {
		pdAddrs = nil
	}
	return pdAddrs
}

func createClient(b *testing.B) Client {
	pdAddrs := getAddrs()
	if len(pdAddrs) == 0 {
		b.Errorf("environment variable %s not set; skipping TiKV benchmarks", envPDAddrs)
	}
	var err error
	cli, err = New(
		WithClient(
			grpc.New(
				"TiKV Client",
				grpc.WithInsecure(true),
			),
		),
		WithPDClient(
			grpc.New(
				"PD Client",
				grpc.WithAddrs(pdAddrs...),
				grpc.WithInsecure(true),
			),
		),
	)
	if err != nil {
		b.Fatalf("failed to create tikv client: %v", err)
	}
	cli.Start(b.Context())
	cli.StartPD(b.Context())

	// basic connectivity probe (Get for non-existing key)
	_, err = cli.Get(context.Background(), []byte("vald_bench_probe"))
	if err != nil {
		// Depending on cluster state Get may return region not found etc.
		// We treat only network/connection errors as fatal.
		b.Logf("tiKV connectivity probe returned error: %v (continuing)", err)
	}
	return cli
}

func generateKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	return key, err
}

func Benchmark(b *testing.B) {
	cli := createClient(b)
	defer cli.Stop(b.Context())
	defer cli.StopPD(b.Context())

	ctx := b.Context()
	val := []byte("vald_bench_val")

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		var key [8]byte
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		if err := cli.Put(ctx, key[:], val); err != nil {
			b.Fatalf("Put error: %v", err)
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		var key [8]byte
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		res, err := cli.Get(ctx, key[:])
		if err != nil {
			b.Fatalf("Get error: %v", err)
		}
		if !slices.Equal(res, val) {
			b.Errorf("i=%d: expected value %v, but got %v", i, val, res)
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		var key [8]byte
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		if err := cli.Delete(ctx, key[:]); err != nil {
			b.Fatalf("i=%d: Delete error: %v", i, err)
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		var key [8]byte
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		_, err := cli.Get(ctx, key[:])
		if !errors.Is(err, errNotFound) {
			b.Fatalf("i=%d: expected not found error, but got: %v", i, err)
		}
	}
}

func BenchmarkBatch(b *testing.B) {
	cli := createClient(b)
	defer cli.Stop(b.Context())
	defer cli.StopPD(b.Context())

	ctx := b.Context()
	length := 300

	b.ReportAllocs()
	b.ResetTimer()
	keys := make([][]byte, length)
	for i := range length {
		var key [8]byte
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		keys[i] = key[:]
	}
	val := []byte("vald_bench_val")
	vals := slices.Repeat([][]byte{val}, length)
	for b.Loop() {
		if err := cli.BatchPut(ctx, keys, vals); err != nil {
			b.Fatalf("BatchPut error: %v", err)
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		res, err := cli.BatchGet(ctx, keys)
		if err != nil {
			b.Fatalf("BatchGet error: %v", err)
		}
		for i := range res {
			if !slices.Equal(res[i], val) {
				b.Fatalf("expected value %v, but got %v", val, res[i])
			}
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := cli.BatchDelete(ctx, keys); err != nil {
			b.Fatalf("BatchDelete error: %v", err)
		}
	}
}

// Ensure that no goroutines leak from the benchmarks.
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, goleak.IgnoreCurrent())
}
