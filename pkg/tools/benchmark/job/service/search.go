//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"
	"math"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (j *job) search(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking search")
	// create data
	vecs := j.hdf5.GetByGroupName(j.dataset.Group)
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
	sres := make([]*payload.Search_Response, j.dataset.Indexes)
	eg, egctx := errgroup.New(ctx)
	eg.SetLimit(j.concurrencyLimit)
	for i := j.dataset.Range.Start; i <= j.dataset.Range.End; i++ {
		iter := i
		eg.Go(func() error {
			log.Debugf("[benchmark job] Start search: iter = %d", iter)
			err := j.limiter.Wait(egctx)
			if err != nil {
				log.Errorf("[benchmark job] limiter error is detected: %s", err.Error())
				if errors.Is(err, context.Canceled) {
					return nil
				}
				select {
				case <-egctx.Done():
					return egctx.Err()
				case ech <- err:
				}
			}
			loopCnt := math.Floor(float64(iter-1) / float64(len(vecs)))
			idx := iter - 1 - (len(vecs) * int(loopCnt))
			if len(vecs[idx]) != j.dimension {
				log.Warn("len(vecs) ", len(vecs[iter]), "is not matched with ", j.dimension)
				return nil
			}
			res, err := j.client.Search(egctx, &payload.Search_Request{
				Vector: vecs[idx],
				Config: cfg,
			})
			if err != nil {
				select {
				case <-egctx.Done():
					log.Errorf("[benchmark job] context error is detected: %s\t%s", err.Error(), egctx.Err())
					return nil
				default:
				}
			}
			if res != nil && j.searchConfig.EnableLinearSearch {
				sres[iter-j.dataset.Range.Start] = res
				log.Debugf("[benchmark job] Finish search: iter = %d, len = %d", iter, len(res.Results))
			} else {
				log.Debugf("[benchmark job] Finish search: iter = %d, res = %v", iter, res)
			}
			log.Debugf("[benchmark job] Finish search: iter = %d, len = %d", iter, len(res.Results))
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Warnf("[benchmark job] search error is detected: err = %s", err.Error())
		return err
	}
	if j.searchConfig.EnableLinearSearch {
		lres := make([]*payload.Search_Response, j.dataset.Indexes)
		for i := j.dataset.Range.Start; i <= j.dataset.Range.End; i++ {
			iter := i
			eg.Go(func() error {
				err := j.limiter.Wait(egctx)
				if err != nil {
					log.Errorf("[benchmark job] limiter error is detected: %s", err.Error())
					if errors.Is(err, context.Canceled) {
						return nil
					}
					select {
					case <-egctx.Done():
						return egctx.Err()
					case ech <- err:
					}
				}
				log.Debugf("[benchmark job] Start linear search: iter = %d", iter)
				loopCnt := math.Floor(float64(i-1) / float64(len(vecs)))
				idx := iter - 1 - (len(vecs) * int(loopCnt))
				if len(vecs[idx]) != j.dimension {
					log.Warn("len(vecs) ", len(vecs[idx]), "is not matched with ", j.dimension)
					return nil
				}
				res, err := j.client.LinearSearch(egctx, &payload.Search_Request{
					Vector: vecs[idx],
					Config: cfg,
				})
				if err != nil {
					select {
					case <-egctx.Done():
						log.Errorf("[benchmark job] context error is detected: %s\t%s", err.Error(), egctx.Err())
						return errors.Join(err, egctx.Err())
					default:
					}
				}
				if res != nil {
					lres[idx-j.dataset.Range.Start] = res
				}
				log.Debugf("[benchmark job] Finish linear search: iter = %d", iter)
				return nil
			})
		}
		err := eg.Wait()
		if err != nil {
			log.Warnf("[benchmark job] linear search error is detected: err = %s", err.Error())
			return err
		}
		recall := make([]float64, j.dataset.Indexes)
		cnt := float64(0)
		for i := 0; i < j.dataset.Indexes; i++ {
			recall[i] = calcRecall(lres[i], sres[i])
			log.Info("[branch job] search recall: ", recall[i])
			cnt += recall[i]
		}
		log.Info("[benchmark job] Total search recall: ", (cnt / float64(len(vecs))))
	}
	log.Info("[benchmark job] Finish benchmarking search")
	return nil
}
