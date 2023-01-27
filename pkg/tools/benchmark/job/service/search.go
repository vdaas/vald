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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func (j *job) search(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking search")
	j.eg.Go(func() (err error) {
		if j.searchConfig == nil {
			err := errors.NewErrInvalidOption("searchConfig", j.searchConfig)
			select {
			case <-ctx.Done():
				if err != context.Canceled {
					ech <- errors.Wrap(err, ctx.Err().Error())
				} else {
					ech <- err
				}
			case ech <- err:
			}
			return err
		}

		vecs := j.hdf5.GetTest()
		if len(vecs) < j.dataset.Indexes {
			log.Infof("[benchmark job] update search iteration from %d to %d", j.dataset.Indexes, len(vecs))
			j.dataset.Indexes = len(vecs)
		}
		timeout, _ := time.ParseDuration(j.searchConfig.Timeout)
		cfg := &payload.Search_Config{
			Num:     uint32(j.searchConfig.Num),
			MinNum:  uint32(j.searchConfig.MinNum),
			Radius:  float32(j.searchConfig.Radius),
			Epsilon: float32(j.searchConfig.Epsilon),
			Timeout: timeout.Nanoseconds(),
		}
		for i := 0; i < j.dataset.Indexes; i++ {
			log.Infof("[benchmark job] Start search: iter = %d\n", i)
			lres, err := j.client.LinearSearch(ctx, &payload.Search_Request{
				Vector: vecs[i],
				Config: cfg,
			})
			if err != nil {
				select {
				case <-ctx.Done():
					if err != context.Canceled {
						ech <- errors.Wrap(err, ctx.Err().Error())
					} else {
						ech <- err
					}
				case ech <- err:
				}
				return err
			}
			bres := testing.Benchmark(func(b *testing.B) {
				b.Helper()
				b.ResetTimer()
				start := time.Now()
				sres, err := j.client.Search(ctx, &payload.Search_Request{
					Vector: vecs[i],
					Config: cfg,
				})
				if err != nil {
					select {
					case <-ctx.Done():
						if errors.Is(err, context.Canceled) {
							ech <- errors.Wrap(err, ctx.Err().Error())
						} else {
							ech <- err
						}
					case ech <- err:
						break
					}
				}
				latency := time.Since(start)
				recall := calcRecall(lres.Results, sres.Results)
				b.ReportMetric(recall, "recall")
				b.ReportMetric(float64(latency.Microseconds()), "latency")
			})
			// TODO: send metrics to the Prometeus
			log.Infof("[benchmark job] Finish search bench: iter= %d \n%#v\n", i, bres)
		}
		return nil
	})

	log.Info("[benchmark job] Finish benchmarking search")
	return nil
}
