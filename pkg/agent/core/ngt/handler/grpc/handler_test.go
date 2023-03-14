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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	info.Init("")
	goleak.VerifyTestMain(m)
}

func newIndexedNGTService(ctx context.Context, eg errgroup.Group, t request.ObjectType, dist vector.Distribution, num int, insertCfg *payload.Insert_Config,
	ngtCfg *config.NGT, ngtOpts []service.Option, overwriteIDs []string, overwriteVectors [][]float32,
) (service.NGT, error) {
	ngt, err := service.New(ngtCfg, append(ngtOpts, service.WithErrGroup(eg), service.WithEnableInMemoryMode(true))...)
	if err != nil {
		return nil, err
	}

	if num > 0 {
		// gen insert request
		reqs, err := request.GenMultiInsertReq(t, dist, num, ngtCfg.Dimension, insertCfg)
		if err != nil {
			return nil, err
		}

		// overwrite ID if needed
		for i, id := range overwriteIDs {
			reqs.Requests[i].Vector.Id = id
		}

		// overwrite Vectors if needed
		for i, v := range overwriteVectors {
			reqs.Requests[i].Vector.Vector = v
		}

		// insert and create index
		for _, req := range reqs.GetRequests() {
			err := ngt.Insert(req.GetVector().GetId(), req.GetVector().GetVector())
			if err != nil {
				return nil, err
			}
		}
		err = ngt.CreateIndex(ctx, 1000)
		if err != nil {
			return nil, err
		}
	}

	return ngt, nil
}
