//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (j *job) search(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking search")
	// create data
	vecs := j.genVec(j.dataset)
	cfg := &payload.Search_Config{
		Num:     uint32(j.searchConfig.Num),
		MinNum:  uint32(j.searchConfig.MinNum),
		Radius:  float32(j.searchConfig.Radius),
		Epsilon: float32(j.searchConfig.Epsilon),
		Timeout: j.timeout.Nanoseconds(),
		AggregationAlgorithm: func() payload.Search_AggregationAlgorithm {
			if len(j.searchConfig.AggregationAlgorithm) > 0 {
				if v, ok := payload.Search_AggregationAlgorithm_value[j.searchConfig.AggregationAlgorithm]; ok {
					return payload.Search_AggregationAlgorithm(v)
				}
			}
			return 0
		}(),
	}
	sres := make([]*payload.Search_Response, len(vecs))
	log.Infof("[benchmark job] Start search")
	for i := 0; i < len(vecs); i++ {
		if len(vecs[i]) != j.dimension {
			log.Warn("len(vecs) ", len(vecs[i]), "is not matched with ", j.dimension)
			continue
		}
		err := j.limiter.Wait(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return errors.Join(err, context.Canceled)
			}
			ech <- err
		}
		res, err := j.client.Search(ctx, &payload.Search_Request{
			Vector: vecs[i],
			Config: cfg,
		})
		log.Infof("[benchmark job] search %d", i)
		if err != nil {
			select {
			case <-ctx.Done():
				if errors.Is(err, context.Canceled) {
					return errors.Join(err, context.Canceled)
				}
				select {
				case <-ctx.Done():
					return errors.Join(err, context.Canceled)
				case ech <- errors.Join(err, ctx.Err()):
				}
			default:
				st, _ := status.FromError(err)
				if st.Code() != codes.NotFound {
					log.Warnf("[benchmark job] search error is detected: code = %d, msg = %s", st.Code(), err.Error())
				}
			}
		}
		sres[i] = res
	}

	if j.searchConfig.EnableLinearSearch {
		lres := make([]*payload.Search_Response, len(vecs))
		for i := 0; i < len(vecs); i++ {
			if len(vecs[i]) != j.dimension {
				log.Warn("len(vecs) ", len(vecs[i]), "is not matched with ", j.dimension)
				continue
			}
			res, err := j.client.LinearSearch(ctx, &payload.Search_Request{
				Vector: vecs[i],
				Config: cfg,
			})
			if err != nil {
				select {
				case <-ctx.Done():
					if errors.Is(err, context.Canceled) {
						return errors.Join(err, context.Canceled)
					}
					select {
					case <-ctx.Done():
						return errors.Join(err, context.Canceled)
					case ech <- errors.Join(err, ctx.Err()):
					}
				default:
					st, _ := status.FromError(err)
					if st.Code() != codes.NotFound {
						log.Warnf("[benchmark job] linear search error is detected: code = %d, msg = %s", st.Code(), err.Error())
					}
				}
			}
			lres[i] = res
		}
		recall := make([]float64, len(vecs))
		for i := 0; i < len(vecs); i++ {
			recall[i] = calcRecall(lres[i], sres[i])
			log.Info("[branch job] search recall: ", recall[i])
		}
	}
	log.Info("[benchmark job] Finish benchmarking search")
	return nil
}
