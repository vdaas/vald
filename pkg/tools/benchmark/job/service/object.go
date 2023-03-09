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
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (j *job) exists(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking exists")
	// create data
	for i := j.dataset.Range.Start; i <= j.dataset.Range.End; i++ {
		err := j.limiter.Wait(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				ech <- err
			}
		}
		res, err := j.client.Exists(ctx, &payload.Object_ID{
			Id: strconv.Itoa(i + j.dataset.Range.Start),
		})
		if err != nil {
			st, _, perr := status.ParseError(err, codes.Internal, "")
			if st.Code() == codes.NotFound {
				continue
			}
			if perr != nil {
				ech <- errors.Join(err, perr)
			}
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
		log.Infof("[benchmark exists job] iter=%d, Id=%s", i, res.GetId())
	}
	log.Info("[benchmark job] Finish benchmarking exists")
	return nil
}

func (j *job) getObject(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking getObject")
	if j.objectConfig == nil {
		log.Warnf("[benchmark job] No get object config is set: %v", j.objectConfig)
	}

	// create data
	vecs := j.genVec(j.dataset)
	log.Debug(j.dataset.Indexes, len(vecs))
	for i := 0; i < len(vecs); i++ {
		log.Infof("[benchmark job] Start getObject: iter = %d", i)
		err := j.limiter.Wait(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				ech <- err
			}
		}
		id := i + j.dataset.Indexes
		ft := []*payload.Filter_Target{}
		if j.objectConfig != nil {
			for i, target := range j.objectConfig.FilterConfig.Targets {
				ft[i] = &payload.Filter_Target{
					Host: target.Host,
					Port: uint32(target.Port),
				}
			}
		}
		res, err := j.client.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: strconv.Itoa(id),
			},
			Filters: &payload.Filter_Config{
				Targets: ft,
			},
		})
		log.Infof("[benchmark get object job] iter=%d, Id=%s, Vec=%v", i, res.GetId(), res.GetVector())
		if err != nil {
			st, _, perr := status.ParseError(err, codes.Internal, "")
			if st.Code() == codes.NotFound {
				continue
			}
			if perr != nil {
				ech <- errors.Join(err, perr)
			}
			select {
			case <-ctx.Done():
				if errors.Is(err, context.Canceled) {
					ech <- errors.Wrap(err, ctx.Err().Error())
				} else {
					st, _, err := status.ParseError(err, codes.NotFound, "hoge", i)
					if st.Code() != codes.NotFound {
						ech <- err
					}
					ech <- err
				}
			case ech <- err:
				break
			}
		}
	}
	log.Info("[benchmark job] Finish benchmarking getObject")
	return nil
}
