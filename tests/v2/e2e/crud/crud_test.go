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
	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
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

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	kclient, err = k8s.NewClient(cfg.Kubernetes.KubeConfig, "")
	if err != nil {
		log.Fatalf("failed to create kubernetes client: %v", err)
	}

	if cfg.Kubernetes.PortForward.Enabled {
		stop, _, err := k8s.Portforward(ctx, kclient,
			cfg.Kubernetes.PortForward.Namespace,
			cfg.Kubernetes.PortForward.PodName,
			cfg.Kubernetes.PortForward.LocalPort,
			cfg.Kubernetes.PortForward.TargetPort)
		if err != nil {
			if stop != nil {
				stop()
			}
			log.Fatalf("failed to portforward: %v", err)
		}
		defer stop()
	}

	ds, err = hdf5.HDF5ToDataset(cfg.Dataset.Name)
	if err != nil {
		log.Fatalf("failed to load dataset: %v", err)
	}
	client, ctx, err = newClient(ctx, cfg.Target, cfg.Metadata)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ech, err := client.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case err := <-ech:
			if err != nil {
				log.Fatalf("client daemon returned error: %v", err)
			}
		}
	}()
	defer func() {
		err = client.Stop(ctx)
		if err != nil {
			log.Fatalf("failed to stop client: %v", err)
		}
	}()

	os.Exit(m.Run())
}

func newClient(
	ctx context.Context, target *iconfig.GRPCClient, meta map[string]string,
) (client vald.Client, mctx context.Context, err error) {
	gopts, err := target.Opts()
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
	t.Logf("%v sleep for %s.", time.Now(), dur)
	time.Sleep(dur)
	t.Logf("%v sleep finished.", time.Now())
}

func recall(t *testing.T, resultIDs []string, neighbors []int) (recall float64) {
	t.Helper()
	ns := map[string]struct{}{}
	for _, n := range neighbors {
		ns[strconv.Itoa(n)] = struct{}{}
	}

	for _, r := range resultIDs {
		if _, ok := ns[r]; ok {
			recall++
		}
	}

	return recall / float64(len(neighbors))
}

func calculateRecall(t *testing.T, res *payload.Search_Response, idx int) (rc float64) {
	t.Helper()
	topKIDs := make([]string, 0, len(res.GetResults()))
	for _, d := range res.GetResults() {
		topKIDs = append(topKIDs, d.GetId())
	}

	if len(topKIDs) == 0 {
		t.Errorf("empty result is returned for test ID %s: %#v", res.GetRequestId(), topKIDs)
		return
	}
	rc = recall(t, topKIDs, ds.Neighbors[idx][:len(topKIDs)])
	return rc
}

func indexStatus(t *testing.T, ctx context.Context) {
	t.Helper()
	{
		res, err := client.IndexInfo(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexInfo %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexInfo %v", err)
			}
		}
		t.Logf("IndexInfo: %v", res.String())
	}
	{
		res, err := client.IndexDetail(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexDetail %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexDetail %v", err)
			}
		}
		t.Logf("IndexDetail: %v", res.String())
	}
	{
		res, err := client.IndexStatistics(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexStatistics %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexStatistics %v", err)
			}
		}
		t.Logf("IndexStatistics: %v", res.String())
	}
	{
		res, err := client.IndexStatisticsDetail(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexStatisticsDetail %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexStatisticsDetail %v", err)
			}
		}
		t.Logf("IndexStatisticsDetail: %v", res.String())
	}
}
