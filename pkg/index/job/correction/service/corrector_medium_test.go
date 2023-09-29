//go:build medium

// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/pkg/index/job/correction/config"
	"github.com/vdaas/vald/pkg/index/job/correction/service"
	"github.com/vdaas/vald/tests/e2e/hdf5"
	"github.com/vdaas/vald/tests/e2e/operation"
)

const (
	host = "vald-lb-gateway.default.svc.cluster.local"
	port = 8081
)

var ds *hdf5.Dataset

func init() {
	var err error
	ds, err = hdf5.HDF5ToDataset("/go/src/github.com/vdaas/vald/hack/benchmark/assets/dataset/fashion-mnist-784-euclidean.hdf5")
	if err != nil {
		panic(err)
	}
}

func TestCorrcotr(t *testing.T) {
	// TODO: ベクトルを
	op, err := operation.New(host, port)
	require.NoError(t, err)

	ctx := context.Background()
	err = op.Insert(t, ctx, operation.Dataset{
		Train: ds.Train[:100],
	})
	require.NoError(t, err)

	defer func(t *testing.T) {
		err = op.Remove(t, ctx, operation.Dataset{
			Train: ds.Train[:100],
		})
		require.NoError(t, err)
	}(t)

	discoverer, err := discoverer.New()
	require.NoError(t, err)

	corrector, err := service.New(&config.Data{}, discoverer)
	require.NoError(t, err)

	corrector.Start(ctx)

	// check the vector state of the cluster with GetObject
	// と思ったが、現状外部からindex replica数を知る方法がないので、先に
	// StreamListObjectをLBに作らないといけない
}
