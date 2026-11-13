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

package grpc

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	sweepDim    = 128
	sweepBase   = 20000
	sweepBatch  = 2000
	sweepNoAuto = "96h"
)

// newSweepAgent serves a real in-memory NGT agent (20k x 128 indexed) over
// loopback TCP gRPC and returns a client for it.
func newSweepAgent(b *testing.B) vald.Client {
	b.Helper()
	cfg := &config.NGT{
		Dimension:              sweepDim,
		DistanceType:           "l2",
		ObjectType:             "float",
		CreationEdgeSize:       50,
		SearchEdgeSize:         50,
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: sweepNoAuto,
		AutoIndexCheckDuration: sweepNoAuto,
		AutoSaveIndexDuration:  sweepNoAuto,
		AutoIndexLength:        10000000000,
		KVSDB:                  &config.KVSDB{Concurrency: 1},
		ErrorBufferLimit:       100,
	}
	n, err := service.New(cfg.Bind())
	if err != nil {
		b.Fatal(err)
	}
	ctx, cancel := context.WithCancel(b.Context())
	b.Cleanup(cancel)
	n.Start(ctx)
	for i, v := range vector.GaussianDistributedFloat32VectorGenerator(sweepBase, sweepDim) {
		if err := n.Insert("base-"+strconv.Itoa(i), v); err != nil {
			b.Fatal(err)
		}
	}
	if err := n.CreateIndex(ctx, 16); err != nil {
		b.Fatal(err)
	}
	h, err := New(WithNGT(n))
	if err != nil {
		b.Fatal(err)
	}
	gs := ggrpc.NewServer()
	vald.RegisterValdServer(gs, h)
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	go func() {
		if serr := gs.Serve(lis); serr != nil {
			b.Log(serr)
		}
	}()
	b.Cleanup(gs.Stop)
	conn, err := ggrpc.NewClient(lis.Addr().String(), ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if cerr := conn.Close(); cerr != nil {
			b.Error(cerr)
		}
	})
	return vald.NewValdClient(conn)
}

// BenchmarkAgentConcurrencySweep measures one agent's throughput (req/s) for
// Insert and Search (k=10) while P requests are kept in flight, to find the
// in-flight count per agent beyond which throughput stops growing: the P the
// LB multi-operation auto concurrency targets per agent.
func BenchmarkAgentConcurrencySweep(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	vc := newSweepAgent(b)
	queries := vector.GaussianDistributedFloat32VectorGenerator(sweepBatch, sweepDim)
	var seq atomic.Uint64
	ops := map[string]func(ctx context.Context, i int) error{
		"Insert": func(ctx context.Context, i int) error {
			_, err := vc.Insert(ctx, &payload.Insert_Request{
				Vector: &payload.Object_Vector{Id: "sweep-" + strconv.FormatUint(seq.Add(1), 10), Vector: queries[i]},
				Config: &payload.Insert_Config{SkipStrictExistCheck: true},
			})
			return err
		},
		"Search": func(ctx context.Context, i int) error {
			_, err := vc.Search(ctx, &payload.Search_Request{
				Vector: queries[i],
				Config: &payload.Search_Config{Num: 10, Epsilon: 0.1, Radius: -1, Timeout: int64(10 * time.Second)},
			})
			return err
		},
	}
	for _, op := range []string{"Insert", "Search"} {
		for _, p := range []int{1, 2, 4, 8, 16, 32, 64, 128, 256} {
			b.Run(op+"/P="+strconv.Itoa(p), func(b *testing.B) {
				do := ops[op]
				var total time.Duration
				b.ResetTimer()
				for range b.N {
					start := time.Now()
					eg, ectx := errgroup.New(b.Context())
					eg.SetLimit(p)
					for i := range sweepBatch {
						eg.Go(func() error { return do(ectx, i) })
					}
					if err := eg.Wait(); err != nil {
						b.Fatal(err)
					}
					total += time.Since(start)
				}
				b.ReportMetric(float64(sweepBatch*b.N)/total.Seconds(), "req/s")
			})
		}
	}
}
