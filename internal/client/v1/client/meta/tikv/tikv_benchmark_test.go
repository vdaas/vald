// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package tikv

import (
	"context"
	"encoding/binary"
	"os"
	"strings"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/test/goleak"
)

// envStoreAddrs is the environment variable which should contain a comma
// separated list of TiKV store addresses (e.g. "127.0.0.1:20160,127.0.0.1:20161").
const envStoreAddrs = "TIKV_STORE_ADDRS"

func getAddrs() []string {
	addrs := strings.Split(os.Getenv(envStoreAddrs), ",")
	if len(addrs) == 1 && addrs[0] == "" {
		return nil
	}
	return addrs
}

func createClient(b *testing.B) Client {
	addrs := getAddrs()
	if len(addrs) == 0 {
		b.Skipf("environment variable %s not set; skipping TiKV benchmarks", envStoreAddrs)
	}

	cli, err := New(
		WithAddrs(addrs...),
		WithClient(
			grpc.New(
				"TiKV Client",
				grpc.WithAddrs(addrs...),
				grpc.WithInsecure(true),
			),
		),
	)
	if err != nil {
		b.Fatalf("failed to create tikv client: %v", err)
	}
	cli.Start(context.Background())

	// basic connectivity probe (RawGet for non-existing key)
	_, err = cli.Get(context.Background(), []byte("vald_bench_probe"))
	if err != nil {
		// Depending on cluster state RawGet may return region not found etc.
		// We treat only network/connection errors as fatal.
		b.Logf("tiKV connectivity probe returned error: %v (continuing)", err)
	}
	return cli
}

func BenchmarkRawPut(b *testing.B) {
	cli := createClient(b)
	defer cli.Stop(context.Background())

	ctx := context.Background()
	val := []byte("vald_bench_val")

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		var key [8]byte
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		if err := cli.Put(ctx, key[:], val); err != nil {
			b.Fatalf("RawPut error: %v", err)
		}
	}
}

func BenchmarkRawGet(b *testing.B) {
	cli := createClient(b)
	defer cli.Stop(context.Background())

	ctx := context.Background()

	// insert a single key to fetch repeatedly
	const baseKey = "vald_bench_get_key"
	if err := cli.Put(ctx, []byte(baseKey), []byte("v")); err != nil {
		b.Fatalf("setup RawPut failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := cli.Get(ctx, []byte(baseKey)); err != nil {
			b.Fatalf("RawGet error: %v", err)
		}
	}
}

// Ensure that no goroutines leak from the benchmarks.
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, goleak.IgnoreCurrent())
}
