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
	"fmt"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	k8s "github.com/vdaas/vald/tests/v2/e2e/kubernetes"
)

type runner struct {
	client vald.Client
	k8s    k8s.Client
}

func TestE2EStrategy(t *testing.T) {
	if cfg == nil || cfg.Strategies == nil {
		t.Fatal("test setting or strategies is nil, please add test configuration yaml file by -config option")
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var (
		kclient k8s.Client
		err     error
	)
	if cfg.Kubernetes != nil {
		kclient, err = k8s.NewClient(cfg.Kubernetes.KubeConfig, "")
		if err != nil {
			t.Errorf("failed to create kubernetes client: %v", err)
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
				t.Errorf("failed to portforward: %v", err)
			}
			defer stop()
		}
	}

	client, ctx, err := newClient(ctx, cfg.Metadata)
	if err != nil {
		t.Fatal("failed to create client: %v", err)
	}
	ech, err := client.Start(ctx)
	if err != nil {
		t.Fatal("failed to start client: %v", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case err := <-ech:
			if err != nil {
				t.Errorf("client daemon returned error: %v", err)
			}
		}
	}()
	defer func() {
		err = client.Stop(ctx)
		if err != nil {
			t.Errorf("failed to stop client: %v", err)
		}
	}()

	r := &runner{
		client: client,
		k8s:    kclient,
	}

	for i, st := range cfg.Strategies {
		r.processStrategy(t, ctx, i, st)
	}
}

func (r *runner) processStrategy(t *testing.T, ctx context.Context, idx int, st *config.Strategy) {
	t.Helper()
	if st == nil {
		return
	}

	t.Run(fmt.Sprintf("#%d: strategy=%s", idx, st.Name), func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, st, "strategy", func(ttt *testing.T, ctx context.Context) error {
			eg, egctx := errgroup.New(ctx)
			if st.Concurrency > 0 {
				eg.SetLimit(int(st.Concurrency))
				ttt.Logf("concurrency is set to %d, the operations will execute concurrently with limit (%d)", st.Concurrency, st.Concurrency)
			} else {
				ttt.Logf("concurrency is not set, the operations will execute concurrently with no limit (%d)", len(st.Operations))
			}

			for i, op := range st.Operations {
				if op != nil {
					i, op := i, op
					eg.Go(safety.RecoverFunc(func() error {
						r.processOperation(ttt, egctx, i, op)
						return nil
					}))
				}
			}

			return eg.Wait()
		}); err != nil {
			tt.Errorf("failed to process operations: %v", err)
		}
	})
}

func (r *runner) processOperation(t *testing.T, ctx context.Context, idx int, op *config.Operation) {
	t.Helper()
	if op == nil {
		return
	}

	t.Run(fmt.Sprintf("#%d: operation=%s", idx, op.Name), func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, op, "operation", func(ttt *testing.T, ctx context.Context) error {
			ttt.Helper()
			for i, e := range op.Executions {
				r.processExecution(ttt, ctx, i, e)
			}
			return nil
		}); err != nil {
			tt.Errorf("failed to process operation: %v", err)
		}
	})
}

func (r *runner) processExecution(t *testing.T, ctx context.Context, idx int, e *config.Execution) {
	t.Helper()
	if e == nil {
		return
	}

	t.Run(fmt.Sprintf("#%d: execution=%s type=%s mode=%s", idx, e.Name, e.Type, e.Mode), func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, e, "execution", func(ttt *testing.T, ctx context.Context) error {
			ttt.Helper()
			ttt.Logf("operation type: %s", e.Type)

			switch e.Type {
			case config.OpSearch,
				config.OpSearchByID,
				config.OpLinearSearch,
				config.OpLinearSearchByID,
				config.OpInsert,
				config.OpUpdate,
				config.OpUpsert,
				config.OpRemove,
				config.OpObject:
				train, test, neighbors := getDatasetSlices(ttt, e)
				switch e.Type {
				case config.OpSearch, config.OpSearchByID, config.OpLinearSearch, config.OpLinearSearchByID:
					r.processSearch(ttt, ctx, test, train, neighbors, e)
				case config.OpInsert:
				case config.OpUpdate:
				case config.OpUpsert:
				case config.OpRemove:
				case config.OpObject:
				}
			case config.OpIndexInfo:
			case config.OpIndexProperty:
			case config.OpKubernetes:
			case config.OpClient:
			case config.OpWait:
			default:
				ttt.Errorf("unsupported operation type: %s detected during execution %d", e.Type, idx)
			}
			return nil
		}); err != nil {
			tt.Errorf("failed to process execution: %v", err)
		}
	})
}

func executeWithTimings[T config.Timing](t *testing.T, ctx context.Context, cfg T, prefix string, fn func(*testing.T, context.Context) error) error {
	t.Helper()

	if delay := cfg.GetDelay(); delay != "" {
		dur, err := delay.Duration()
		if err != nil {
			t.Errorf("failed to parse delay duration: %s, error: %v", delay, err)
		}
		if dur > 0 {
			t.Logf("delay is set to %s, this %s will start after %s", delay, prefix, dur.String())
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(dur):
			}
		}
	}

	var cancel context.CancelFunc = func() {}
	if timeout := cfg.GetTimeout(); timeout != "" {
		dur, err := timeout.Duration()
		if err != nil {
			t.Errorf("failed to parse timeout duration: %s, error: %v", timeout, err)
		}
		if dur > 0 {
			t.Logf("timeout is set to %s, this %s will stop after %s", timeout, prefix, dur.String())
			ctx, cancel = context.WithTimeout(ctx, dur)
		}
	}
	defer cancel()

	err := fn(t, ctx)

	if wait := cfg.GetWait(); wait != "" {
		dur, err := wait.Duration()
		if err != nil {
			t.Errorf("failed to parse wait duration: %s, error: %v", wait, err)
		}
		if dur > 0 {
			t.Logf("wait is set to %s, this %s is already finished, but will wait for %s", wait, prefix, dur.String())
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(dur):
			}
		}
	}

	return err
}

func getDatasetSlices(t *testing.T, e *config.Execution) (train, test [][]float32, neighbors [][]int) {
	t.Helper()
	if len(ds.Train) > int(e.Offset)+int(e.Num) {
		train = ds.Train[e.Offset : e.Offset+e.Num]
	} else {
		t.Errorf("train data is not enough, offset: %d, num: %d, total: %d", e.Offset, e.Num, len(ds.Train))
	}

	if len(ds.Test) > int(e.Offset)+int(e.Num) {
		test = ds.Test[e.Offset : e.Offset+e.Num]
	} else {
		t.Errorf("test data is not enough, offset: %d, num: %d, total: %d", e.Offset, e.Num, len(ds.Test))
	}

	if len(ds.Neighbors) > int(e.Offset)+int(e.Num) {
		neighbors = ds.Neighbors[e.Offset : e.Offset+e.Num]
	} else {
		t.Errorf("neighbor data is not enough, offset: %d, num: %d, total: %d", e.Offset, e.Num, len(ds.Neighbors))
	}
	return train, test, neighbors
}
