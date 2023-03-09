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
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func (j *job) search(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking search")
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
	// create data
	vecs := j.genVec(j.dataset)
	timeout, _ := time.ParseDuration(j.searchConfig.Timeout)
	cfg := &payload.Search_Config{
		Num:     uint32(j.searchConfig.Num),
		MinNum:  uint32(j.searchConfig.MinNum),
		Radius:  float32(j.searchConfig.Radius),
		Epsilon: float32(j.searchConfig.Epsilon),
		Timeout: timeout.Nanoseconds(),
	}
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
				if !errors.Is(err, context.Canceled) {
					ech <- errors.Wrap(err, ctx.Err().Error())
				} else {
					ech <- err
				}
			case ech <- err:
			}
			return err
		}
		lres[i] = res
	}
	// TODO: apply rpc from crd setting params
	sres := make([]*payload.Search_Response, len(vecs))
	log.Infof("[benchmark job] Start search")
	for i := 0; i < len(vecs); i++ {
		if len(vecs[i]) != j.dimension {
			log.Warn("len(vecs) ", len(vecs[i]), "is not matched with ", j.dimension)
			continue
		}
		err := j.limiter.Wait(ctx)
		if err != nil {
			errors.Is(context.Canceled, err)
			ech <- err
			break
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
					ech <- errors.Wrap(err, ctx.Err().Error())
				} else {
					ech <- err
				}
			case ech <- err:
				break
			}
		}
		sres[i] = res
	}
	recall := make([]float64, len(vecs))
	for i := 0; i < len(vecs); i++ {
		recall[i] = calcRecall(lres[i].Results, sres[i].Results)
		log.Info("[branch job] search recall: ", recall[i])
	}
	log.Info("[benchmark job] Finish benchmarking search")
	return nil
}
