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
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (j *job) update(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking update")
	// create data
	vecs := j.hdf5.GetByGroupName(j.dataset.Group)
	cfg := &payload.Update_Config{
		SkipStrictExistCheck:  j.updateConfig.SkipStrictExistCheck,
		DisableBalancedUpdate: j.updateConfig.DisableBalancedUpdate,
	}
	if j.timestamp > int64(0) {
		cfg.Timestamp = j.timestamp
	}
	eg, egctx := errgroup.New(ctx)
	eg.SetLimit(j.concurrencyLimit)
	for i := j.dataset.Range.Start; i <= j.dataset.Range.End; i++ {
		iter := i
		eg.Go(func() error {
			log.Debugf("[benchmark job] Start update: iter = %d", iter)
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
			res, err := j.client.Update(egctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:     strconv.Itoa(iter),
					Vector: addNoiseToVec(vecs[idx]),
				},
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
				log.Infof("[benchmark job] iter=%d, Name=%s, Uuid=%s, Ips=%v", iter, res.Name, res.Uuid, res.Ips)
			}
			log.Debugf("[benchmark job] Finish update: iter = %d", iter)
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Warnf("[benchmark job] update error is detected: err = %s", err.Error())
		return err
	}
	log.Info("[benchmark job] Finish benchmarking upsert")
	return nil
}
