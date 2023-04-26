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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (j *job) remove(ctx context.Context, ech chan error) error {
	log.Info("[benchmark job] Start benchmarking remove")
	// create data
	vecs := j.genVec(j.dataset)
	cfg := &payload.Remove_Config{
		SkipStrictExistCheck: j.removeConfig.SkipStrictExistCheck,
	}
	if j.timestamp > int64(0) {
		cfg.Timestamp = j.timestamp
	}
	for i := 0; i < len(vecs); i++ {
		log.Infof("[benchmark job] Start remove: iter = %d", i)
		err := j.limiter.Wait(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return errors.Join(err, context.Canceled)
			}
			ech <- err
		}
		res, err := j.client.Remove(ctx, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: strconv.Itoa(i),
			},
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
				log.Warnf("[benchmark job] remove error is detected: code = %d, msg = %s", st.Code(), err.Error())
			}
		}
		log.Infof("[benchmark job] Finish remove: iter= %d \n%v", i, res)
	}

	log.Info("[benchmark job] Finish benchmarking remove")
	return nil
}
