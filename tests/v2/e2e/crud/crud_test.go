//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// package crud provides e2e tests using ann-benchmarks datasets
package crud

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/e2e/hdf5"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	k8s "github.com/vdaas/vald/tests/v2/e2e/kubernetes"
	"google.golang.org/grpc/metadata"
)

var (
	cfg *config.Data

	ctx     context.Context
	client  vald.Client
	kclient k8s.Client

	ds *hdf5.Dataset
)

func TestMain(m *testing.M) {
	log.Init()
	var err error
	p, fail, err := params.New(
		params.WithName("vald/e2e"),
		params.WithOverrideDefault(true),
		params.WithArgumentFilters(
			func(s string) bool {
				return strings.HasPrefix(s, "-test.")
			},
		),
	).Parse()
	if fail || err != nil || p.ConfigFilePath() == "" {
		log.Fatalf("failed to parse the parameters: %v", err)
	}

	if testing.Short() {
		log.Info("skipping this pkg test when -short because e2e test takes a long time")
		os.Exit(0)
	}

	cfg, err = config.Load(p.ConfigFilePath())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ds, err = hdf5.HDF5ToDataset(cfg.Dataset.Name)
	if err != nil {
		log.Fatalf("failed to load dataset: %v", err)
	}
	os.Exit(m.Run())
}

func newClient(
	ctx context.Context, meta map[string]string,
) (client vald.Client, mctx context.Context, err error) {
	if cfg == nil || cfg.Target == nil {
		return nil, nil, errors.ErrGRPCTargetAddrNotFound
	}
	gopts, err := cfg.Target.Opts()
	if err != nil {
		return nil, nil, err
	}
	client, err = vald.New(
		vald.WithClient(
			grpc.New(gopts...),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	if meta != nil {
		mctx = metadata.NewOutgoingContext(ctx, metadata.New(meta))
	}
	return client, mctx, nil
}

func sleep(t *testing.T, dur time.Duration) {
	t.Helper()
	t.Logf("%v sleep for %s.", time.Now(), dur)
	time.Sleep(dur)
	t.Logf("%v sleep finished.", time.Now())
}
