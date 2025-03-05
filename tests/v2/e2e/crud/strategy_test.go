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
			t.Fatalf("failed to create kubernetes client: %v", err)
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
				t.Fatalf("failed to portforward: %v", err)
			}
			defer stop()
		}
	}

	client, ctx, err := newClient(ctx, cfg.Metadata)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	ech, err := client.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case err := <-ech:
			if err != nil {
				t.Fatalf("client daemon returned error: %v", err)
			}
		}
	}()
	defer func() {
		err = client.Stop(ctx)
		if err != nil {
			t.Fatalf("failed to stop client: %v", err)
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
		if st.Delay != "" {
			dur, err := st.Delay.Duration()
			if err != nil {
				tt.Fatalf("failed to parse delay duration: %s, error: %v", st.Delay, err)
			}
			if dur > 0 {
				tt.Logf("delay is set to %s, this strategy will start after %s", st.Delay, dur.String())
				select {
				case <-ctx.Done():
					tt.Fatal(ctx.Err())
				case <-time.After(dur):
				}
			}
		}
		if st.Timeout != "" {
			dur, err := st.Timeout.Duration()
			if err != nil {
				tt.Fatalf("failed to parse timeout duration: %s, error: %v", st.Timeout, err)
			}
			if dur > 0 {
				tt.Logf("timeout is set to %s, this strategy will stop after %s", st.Timeout, dur.String())
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, dur)
				defer cancel()
			}
		}
		eg, egctx := errgroup.New(ctx)
		if st.Concurrency > 0 {
			eg.SetLimit(int(st.Concurrency))
			tt.Logf("concurrency is set to %d, the operations will execute concurrently with limit (%d)", st.Concurrency, st.Concurrency)
		} else {
			tt.Logf("concurrency is not set, the operations will execute concurrently with no limit (%d)", len(st.Operations))
		}
		for i, op := range st.Operations {
			if op != nil {
				eg.Go(safety.RecoverFunc(func() error {
					r.processOperation(tt, egctx, i, op)
					return nil
				}))
			}
		}
		err := eg.Wait()
		if err != nil {
			tt.Fatalf("failed to execute operations: %v", err)
		}
		if st.Wait != "" {
			dur, err := st.Wait.Duration()
			if err != nil {
				tt.Fatalf("failed to parse wait duration: %s, error: %v", st.Wait, err)
			}
			if dur > 0 {
				tt.Logf("wait is set to %s, this strategy is already finished, but will wait for %s", st.Wait, dur.String())
				select {
				case <-ctx.Done():
					tt.Fatal(ctx.Err())
				case <-time.After(dur):
				}
			}
		}
	})
}

func (r *runner) processOperation(t *testing.T, ctx context.Context, idx int, op *config.Operation) {
	t.Helper()
	if op == nil {
		return
	}
	t.Run(fmt.Sprintf("#%d: operation=%s", idx, op.Name), func(tt *testing.T) {
		if op.Delay != "" {
			dur, err := op.Delay.Duration()
			if err != nil {
				tt.Fatalf("failed to parse delay duration: %s, error: %v", op.Delay, err)
			}
			if dur > 0 {
				tt.Logf("delay is set to %s, this operation will start after %s", op.Delay, dur.String())
				select {
				case <-ctx.Done():
					tt.Fatal(ctx.Err())
				case <-time.After(dur):
				}
			}
		}
		if op.Timeout != "" {
			dur, err := op.Timeout.Duration()
			if err != nil {
				tt.Fatalf("failed to parse timeout duration: %s, error: %v", op.Timeout, err)
			}
			if dur > 0 {
				tt.Logf("timeout is set to %s, this operation will stop after %s", op.Timeout, dur.String())
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, dur)
				defer cancel()
			}
		}
		for i, e := range op.Executions {
			r.processExecution(tt, ctx, i, e)
		}
		if op.Wait != "" {
			dur, err := op.Wait.Duration()
			if err != nil {
				tt.Fatalf("failed to parse wait duration: %s, error: %v", op.Wait, err)
			}
			if dur > 0 {
				tt.Logf("wait is set to %s, this operation is already finished, but will wait for %s", op.Wait, dur.String())
				select {
				case <-ctx.Done():
					tt.Fatal(ctx.Err())
				case <-time.After(dur):
				}
			}
		}
	})
	return
}

func (r *runner) processExecution(t *testing.T, ctx context.Context, idx int, e *config.Execution) {
	t.Helper()
	if e == nil {
		return
	}
	t.Run(fmt.Sprintf("#%d: execution=%s", idx, e.Name), func(tt *testing.T) {
		if e.Delay != "" {
			dur, err := e.Delay.Duration()
			if err != nil {
				tt.Fatalf("failed to parse delay duration: %s, error: %v", e.Delay, err)
			}
			if dur > 0 {
				tt.Logf("delay is set to %s, this execution will start after %s", e.Delay, dur.String())
				select {
				case <-ctx.Done():
					tt.Fatal(ctx.Err())
				case <-time.After(dur):
				}
			}
		}
		if e.Timeout != "" {
			dur, err := e.Timeout.Duration()
			if err != nil {
				tt.Fatalf("failed to parse timeout duration: %s, error: %v", e.Timeout, err)
			}
			if dur > 0 {
				tt.Logf("timeout is set to %s, this execution will stop after %s", e.Timeout, dur.String())
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, dur)
				defer cancel()
			}
		}
		tt.Logf("operation type: %s", e.Type)
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
			var (
				train, test [][]float32
				neighbors   [][]int
			)
			if len(ds.Train) > int(e.Offset)+int(e.Num) {
				train = ds.Train[e.Offset : e.Offset+e.Num]
			} else {
				tt.Fatalf("train data is not enough, offset: %d, num: %d, total: %d", e.Offset, e.Num, len(ds.Train))
			}
			if len(ds.Test) > int(e.Offset)+int(e.Num) {
				test = ds.Test[e.Offset : e.Offset+e.Num]
			} else {
				tt.Fatalf("test data is not enough, offset: %d, num: %d, total: %d", e.Offset, e.Num, len(ds.Test))
			}
			if len(ds.Neighbors) > int(e.Offset)+int(e.Num) {
				neighbors = ds.Neighbors[e.Offset : e.Offset+e.Num]
			} else {
				tt.Fatalf("neighbor data is not enough, offset: %d, num: %d, total: %d", e.Offset, e.Num, len(ds.Neighbors))
			}
			switch e.Type {
			case config.OpSearch:
			case config.OpSearchByID:
			case config.OpLinearSearch:
			case config.OpLinearSearchByID:
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
			tt.Fatalf("unsupported operation type: %s detected during execution %d", e.Type, idx)
		}
		if e.Wait != "" {
			dur, err := e.Wait.Duration()
			if err != nil {
				tt.Fatalf("failed to parse wait duration: %s, error: %v", e.Wait, err)
			}
			if dur > 0 {
				tt.Logf("wait is set to %s, this execution is already finished, but will wait for %s", e.Wait, dur.String())
				select {
				case <-ctx.Done():
					tt.Fatal(ctx.Err())
				case <-time.After(dur):
				}
			}
		}
	})
}
