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

func (j *job) upsert(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking upsert")
	if j.upsertConfig == nil {
		err := errors.NewErrInvalidOption("upsertConfig", j.upsertConfig)
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
	timestamp, _ := strconv.Atoi(j.upsertConfig.Timestamp)
	cfg := &payload.Upsert_Config{
		SkipStrictExistCheck:  j.upsertConfig.SkipStrictExistCheck,
		Timestamp:             int64(timestamp),
		DisableBalancedUpdate: j.upsertConfig.DisableBalancedUpdate,
	}
	for i := 0; i < len(vecs); i++ {
		log.Infof("[benchmark job] Start upsert: iter = %d", i)
		err := j.limiter.Wait(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				ech <- err
			}
		}
		id := i + j.dataset.Range.Start
		res, err := j.client.Upsert(ctx, &payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:     strconv.Itoa(id),
				Vector: vecs[i],
			},
			Config: cfg,
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
		log.Infof("[benchmark job] iter=%d, Name=%s, Uuid=%s, Ips=%v", i, res.Name, res.Uuid, res.Ips)
	}

	log.Info("[benchmark job] Finish benchmarking upsert")
	return nil
}
