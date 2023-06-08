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
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func (j *job) exists(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking exists")
	eg, egctx := errgroup.New(ctx)
	for i := j.dataset.Range.Start; i <= j.dataset.Range.End; i++ {
		idx := i
		eg.Go(func() error {
			log.Debugf("[benchmark job] Start exists: iter = %d", i)
			err := j.limiter.Wait(egctx)
			if err != nil {
				log.Errorf("[benchmark job] limiter error is detected: %s", err.Error())
				if errors.Is(err, context.Canceled) {
					return nil
					// return errors.Join(err, context.Canceled)
				}
				select {
				case <-egctx.Done():
					return egctx.Err()
				case ech <- err:
				}
			}
			res, err := j.client.Exists(egctx, &payload.Object_ID{
				Id: strconv.Itoa(idx),
			})
			if err != nil {
				select {
				case <-egctx.Done():
					log.Errorf("[benchmark job] context error is detected: %s\t%s", err.Error(), egctx.Err())
					return nil
					// return errors.Join(err, egctx.Err())
				default:
					// if st, ok := status.FromError(err); ok {
					// 	log.Warnf("[benchmark job] exists error is detected: code = %d, msg = %s", st.Code(), err.Error())
					// }
				}
			}
			log.Debugf("[benchmark job] Finish exists: iter= %d \n%v\n", idx, res)
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Warnf("[benchmark job] exists RPC error is detected: err = %s", err.Error())
		return err
	}
	log.Info("[benchmark job] Finish benchmarking exists")
	return nil
}

func (j *job) getObject(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking getObject")
	// create data
	eg, egctx := errgroup.New(ctx)
	for i := j.dataset.Range.Start; i <= j.dataset.Range.End; i++ {
		log.Infof("[benchmark job] Start get object: iter = %d", i)
		ft := []*payload.Filter_Target{}
		if j.objectConfig != nil {
			for i, target := range j.objectConfig.FilterConfig.Targets {
				ft[i] = &payload.Filter_Target{
					Host: target.Host,
					Port: uint32(target.Port),
				}
			}
		}
		idx := i
		eg.Go(func() error {
			log.Debugf("[benchmark job] Start get object: iter = %d", idx)
			err := j.limiter.Wait(egctx)
			if err != nil {
				log.Errorf("[benchmark job] limiter error is detected: %s", err.Error())
				if errors.Is(err, context.Canceled) {
					// return errors.Join(err, context.Canceled)
					return nil
				}
				select {
				case <-egctx.Done():
					return egctx.Err()
				case ech <- err:
				}
			}
			res, err := j.client.GetObject(egctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: strconv.Itoa(idx),
				},
				Filters: &payload.Filter_Config{
					Targets: ft,
				},
			})
			if err != nil {
				select {
				case <-egctx.Done():
					log.Errorf("[benchmark job] context error is detected: %s\t%s", err.Error(), egctx.Err())
					// return errors.Join(err, egctx.Err())
					return nil
				default:
					// if st, ok := status.FromError(err); ok {
					// 	log.Warnf("[benchmark job] object error is detected: code = %d, msg = %s", st.Code(), err.Error())
					// }
				}
			}
			if res != nil {
				log.Infof("[benchmark get object job] iter=%d, Id=%s, Vec=%v", idx, res.GetId(), res.GetVector())
			}
			log.Debugf("[benchmark job] Finish get object: iter= %d \n%v\n", idx, res)
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Warnf("[benchmark job] object error is detected: err = %s", err.Error())
		return err
	}
	log.Info("[benchmark job] Finish benchmarking getObject")
	return nil
}
