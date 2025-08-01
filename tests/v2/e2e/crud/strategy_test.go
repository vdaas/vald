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

	agent "github.com/vdaas/vald/internal/client/v1/client/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	k8s "github.com/vdaas/vald/tests/v2/e2e/kubernetes"
	"github.com/vdaas/vald/tests/v2/e2e/kubernetes/portforward"
	"google.golang.org/grpc/metadata"
)

type runner struct {
	rootCtx context.Context
	client  vald.Client
	aclient agent.Client
	k8s     k8s.Client
}

func TestE2EStrategy(t *testing.T) {
	if cfg == nil || cfg.Strategies == nil {
		t.Fatal("test setting or strategies is nil, please add test configuration yaml file by -config option")
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var err error
	r := new(runner)
	if cfg.Kubernetes != nil {
		r.k8s, err = k8s.NewClient(cfg.Kubernetes.KubeConfig, "")
		if err != nil {
			t.Errorf("failed to create kubernetes client: %v", err)
		}
		if cfg.Kubernetes.PortForward.Enabled {
			if r.k8s == nil {
				t.Fatal("kubernetes client is nil")
			}

			pfd, err := portforward.New(
				portforward.WithAddress("localhost"),
				portforward.WithClient(r.k8s),
				portforward.WithNamespace(cfg.Kubernetes.PortForward.Namespace),
				portforward.WithServiceName(cfg.Kubernetes.PortForward.ServiceName),
				portforward.WithPorts(map[uint16]uint16{
					cfg.Kubernetes.PortForward.LocalPort.Port(): cfg.Kubernetes.PortForward.TargetPort.Port(),
				}),
			)
			if err != nil {
				if pfd != nil {
					pfd.Stop()
				}
				t.Fatalf("failed to construct portforward client: %v", err)
			}
			_, err = pfd.Start(ctx)
			if err != nil {
				if pfd != nil {
					pfd.Stop()
				}
				t.Fatalf("failed to start portforward: %v", err)
			}
			defer pfd.Stop()
		}
	}
	r.client, ctx, err = newClient(t, ctx, cfg.Metadata)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if r.client == nil {
		t.Fatal("gRPC E2E client is nil")
	}
	ech, err := r.client.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}
	r.aclient, err = agent.New(agent.WithValdClient(r.client))
	if err != nil {
		t.Fatalf("failed to create agent client: %v", err)
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
		err = r.client.Stop(ctx)
		if err != nil {
			t.Errorf("failed to stop client: %v", err)
		}
	}()
	t.Logf("connected addrs: %v", r.client.GRPCClient().ConnectedAddrs(ctx))
	t.Run("Run E2E V2 Scenarios", func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, cfg, cfg.FilePath, "e2e", func(ttt *testing.T, ctx context.Context) error {
			ttt.Helper()
			for i, st := range cfg.Strategies {
				r.processStrategy(ttt, ctx, i, st)
			}
			return nil
		}); err != nil {
			tt.Errorf("failed to process operations: %v", err)
		}
	})
}

func (r *runner) processStrategy(t *testing.T, ctx context.Context, idx int, st *config.Strategy) {
	t.Helper()
	if r == nil || st == nil {
		return
	}
	t.Run(fmt.Sprintf("#%d: strategy=%s", idx, st.Name), func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, st, st.Name, "strategy", func(ttt *testing.T, ctx context.Context) error {
			ttt.Helper()
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
					eg.Go(func() error {
						r.processOperation(ttt, egctx, i, op)
						return nil
					})
				}
			}
			return eg.Wait()
		}); err != nil {
			tt.Errorf("failed to process operations: %v", err)
		}
	})
}

func (r *runner) processOperation(
	t *testing.T, ctx context.Context, idx int, op *config.Operation,
) {
	t.Helper()
	if r == nil || op == nil {
		return
	}

	t.Run(fmt.Sprintf("#%d: operation=%s", idx, op.Name), func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, op, op.Name, "operation", func(ttt *testing.T, ctx context.Context) error {
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
	if r == nil || e == nil {
		return
	}

	t.Run(fmt.Sprintf("#%d: execution=%s type=%s mode=%s", idx, e.Name, e.Type, e.Mode), func(tt *testing.T) {
		if err := executeWithTimings(tt, ctx, e, e.Name, "execution", func(ttt *testing.T, ctx context.Context) error {
			ttt.Helper()
			switch e.Type {
			case config.OpSearch,
				config.OpSearchByID,
				config.OpLinearSearch,
				config.OpLinearSearchByID,
				config.OpInsert,
				config.OpUpdate,
				config.OpUpsert,
				config.OpRemove,
				config.OpRemoveByTimestamp,
				config.OpObject,
				config.OpListObject,
				config.OpTimestamp,
				config.OpExists:
				train, test, neighbors := getDatasetSlices(ttt, e)
				if e.BaseConfig != nil {
					start := time.Now()
					log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d, num: %d, offset: %d, parallelism: %d, qps: %d",
						e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx, e.Num, e.Offset, e.Parallelism, e.QPS)
					defer func() {
						log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d, num: %d, offset: %d, parallelism: %d, qps: %d",
							e.Name, time.Since(start).String(), e.Type, e.Mode, idx, e.Num, e.Offset, e.Parallelism, e.QPS)
					}()
				}
				switch e.Type {
				case config.OpSearch,
					config.OpSearchByID,
					config.OpLinearSearch,
					config.OpLinearSearchByID:
					r.processSearch(ttt, ctx, train, test, neighbors, e)
				case config.OpInsert,
					config.OpUpdate,
					config.OpUpsert,
					config.OpRemove,
					config.OpRemoveByTimestamp:
					r.processModification(ttt, ctx, train, e)
				case config.OpObject,
					config.OpListObject,
					config.OpTimestamp,
					config.OpExists:
					r.processObject(ttt, ctx, train, e)
				}
			case config.OpIndexInfo,
				config.OpIndexDetail,
				config.OpIndexStatistics,
				config.OpIndexStatisticsDetail,
				config.OpIndexProperty,
				config.OpFlush:
				start := time.Now()
				log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d",
					e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx)
				defer func() {
					log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d",
						e.Name, time.Since(start).String(), e.Type, e.Mode, idx)
				}()
				r.processIndex(ttt, ctx, e)
			case config.OpCreateIndex,
				config.OpSaveIndex,
				config.OpCreateAndSaveIndex:
				start := time.Now()
				log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d",
					e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx)
				defer func() {
					log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d",
						e.Name, time.Since(start).String(), e.Type, e.Mode, idx)
				}()
				r.processAgent(t, ctx, e)
			case config.OpKubernetes:
				if e.Kubernetes != nil {
					start := time.Now()
					log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d, kubernetes action: %s, kind: %s, namespace: %s, name: %s, status: %s",
						e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx, e.Kubernetes.Action, e.Kubernetes.Kind, e.Kubernetes.Namespace, e.Kubernetes.Name, e.Kubernetes.Status)
					defer func() {
						log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d, kubernetes action: %s, kind: %s, namespace: %s, name: %s, status: %s",
							e.Name, time.Since(start).String(), e.Type, e.Mode, idx, e.Kubernetes.Action, e.Kubernetes.Kind, e.Kubernetes.Namespace, e.Kubernetes.Name, e.Kubernetes.Status)
					}()
					r.processKubernetes(ttt, ctx, e)
				}
			case config.OpClient:
				// TODO implement gRPC client operation here, eg. start, stop, etc.
			case config.OpWait:
				// do nothing
			default:
				ttt.Errorf("unsupported operation type: %s detected during execution %d", e.Type, idx)
			}
			return nil
		}); err != nil {
			tt.Errorf("failed to process execution: %v", err)
		}
	})
}

func executeWithTimings[T interface {
	config.Timing
	config.Repeats
}](
	t *testing.T,
	ctx context.Context,
	cfg T,
	name, prefix string,
	fn func(*testing.T, context.Context) error,
) (err error) {
	t.Helper()
	if delay := cfg.GetDelay(); delay != "" {
		dur, err := delay.Duration()
		if err != nil {
			t.Errorf("failed to parse delay duration: %s, error: %v", delay, err)
		}
		if dur > 0 {
			log.Infof("delay is set to %s, this %s/%s will start after %s", delay, prefix, name, dur.String())
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(dur):
			}
		}
	}

	if timeout := cfg.GetTimeout(); timeout != "" {
		dur, err := timeout.Duration()
		if err != nil {
			t.Errorf("failed to parse timeout duration: %s, error: %v", timeout, err)
		}
		if dur > 0 {
			t.Logf("timeout is set to %s, this %s/%s will stop after %s", timeout, prefix, name, dur.String())
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, dur)
			defer cancel()
		}
	}

	err = executeWithRepeats(t, ctx, name, prefix, cfg.GetRepeats(), fn)

	if wait := cfg.GetWait(); wait != "" {
		dur, werr := wait.Duration()
		if werr != nil {
			t.Errorf("failed to parse wait duration: %s, error: %v", wait, werr)
			return err
		}
		if dur > 0 {
			log.Infof("\"%s.wait: %s\", wait configuration detected, this %s/%s is already finished, will wait for %s", prefix, wait, prefix, name, dur.String())
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(dur):
			}
		}
	}

	return err
}

func executeWithRepeats(
	t *testing.T,
	ctx context.Context,
	name, prefix string,
	repeats uint64,
	fn func(*testing.T, context.Context) error,
) (err error) {
	t.Helper()
	if repeats > 1 {
		for idx := range repeats {
			task := fmt.Sprintf("Repeat %s for %s (%d/%d)", prefix, name, idx+1, repeats)
			select {
			case <-ctx.Done():
				err = ctx.Err()
				log.Warnf("context canceled due to %s during %s", err.Error(), task)
				if err != nil && errors.IsNot(err, context.Canceled, context.DeadlineExceeded) {
					return err
				}
				return nil
			default:
			}
			log.Info(task)
			t.Run(task, func(tt *testing.T) {
				tt.Helper()
				ierr := fn(tt, ctx)
				if ierr != nil && errors.IsNot(ierr, context.Canceled, context.DeadlineExceeded) {
					err = errors.Join(err, ierr)
				}
			})
		}
		return err
	}
	return fn(t, ctx)
}

func newClient(
	t *testing.T, ctx context.Context, meta map[string]string,
) (client vald.Client, mctx context.Context, err error) {
	t.Helper()
	if cfg == nil || cfg.Target == nil {
		return nil, nil, errors.ErrGRPCTargetAddrNotFound
	}
	gopts, err := cfg.Target.Opts()
	if err != nil {
		return nil, nil, err
	}
	client, err = vald.New(
		vald.WithClient(
			grpc.New("E2E Strategy Testing Vald Client", gopts...),
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
