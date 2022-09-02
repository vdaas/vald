// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package operation

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/io"
)

func (o *operation) Search(ctx context.Context, b *testing.B, ds assets.Dataset) {
	b.ResetTimer()
	b.Run("Search", func(b *testing.B) {
		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			v, err := ds.Train(i % ds.TrainSize())
			if err != nil {
				b.Error(err)
				continue
			}
			_, err = o.client.Search(ctx, &payload.Search_Request{
				Vector: v.([]float32),
				Config: cfg,
			})
			if err != nil {
				grpcError(b, err)
			}
		}
	})
}

func (o *operation) StreamSearch(ctx context.Context, b *testing.B, ds assets.Dataset) {
	b.ResetTimer()
	b.Run("StreamSearch", func(b *testing.B) {
		sc, err := o.client.StreamSearch(ctx)
		if err != nil {
			b.Fatal(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		go func() {
			defer wg.Done()

			for {
				res, err := sc.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					grpcError(b, err)
					continue
				}
				if res.GetResponse() == nil {
					b.Error("returned response is nil")
				}
			}
		}()

		for i := 0; i < b.N; i++ {
			idx := i % ds.TrainSize()
			v, err := ds.Train(idx)
			if err != nil {
				b.Error(err)
				continue
			}
			err = sc.Send(&payload.Search_Request{
				Vector: v.([]float32),
				Config: cfg,
			})
			if err != nil {
				b.Error(err)
			}
		}

		if err := sc.CloseSend(); err != nil {
			b.Fatal(err)
		}
		wg.Wait()
	})
}

func (o *operation) SearchByID(ctx context.Context, b *testing.B, maxIdNum int) {
	b.Run("SearchByID", func(b *testing.B) {
		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := o.client.SearchByID(ctx, &payload.Search_IDRequest{
				Id:     strconv.Itoa(i % maxIdNum),
				Config: cfg,
			})
			if err != nil {
				grpcError(b, err)
			}
		}
	})
}

func (o *operation) StreamSearchByID(ctx context.Context, b *testing.B, maxIdNum int) {
	b.ResetTimer()
	b.Run("StreamSearchByID", func(b *testing.B) {
		sc, err := o.client.StreamSearchByID(ctx)
		if err != nil {
			b.Fatal(err)
		}
		wg := sync.WaitGroup{}
		wg.Add(1)

		cfg := &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.1,
		}
		b.ResetTimer()

		go func() {
			defer wg.Done()

			for {
				res, err := sc.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					grpcError(b, err)
					continue
				}
				if res.GetResponse() == nil {
					b.Error("returned response is nil")
				}
			}
		}()

		for i := 0; i < b.N; i++ {
			err = sc.Send(&payload.Search_IDRequest{
				Id:     strconv.Itoa(i % maxIdNum),
				Config: cfg,
			})
			if err != nil {
				b.Fatal(err)
			}
		}

		if err := sc.CloseSend(); err != nil {
			b.Fatal(err)
		}
		wg.Wait()
	})
}
