//go:build e2e

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	k8s "github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/portforward"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/retry"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
	"google.golang.org/grpc/metadata"
)

type runner struct {
	rootCtx context.Context
	client  vald.Client
	aclient agent.Client
	k8s     k8s.Client
	clients *resource.Clients
}

func TestE2EStrategy(t *testing.T) {
	runE2EStrategy(test.NewNode(t))
}

// BenchmarkE2EStrategy drives the exact same scenario configuration as
// TestE2EStrategy through the benchmark harness: every strategy, operation
// and execution becomes a b.Run sub-benchmark, and each execution-level pass
// is wrapped in b.Loop (see executeWithTimings) so ns/op reports the time of
// one full configured execution (including its repeats), while the scenario's
// delay/wait sleeps stay outside the measured window. Run it with
// -benchtime 1x (or a small fixed count) since one iteration already
// performs the execution's full configured request load.
func BenchmarkE2EStrategy(b *testing.B) {
	runE2EStrategy(test.NewNode(b))
}

// runE2EStrategy is deliberately non-generic: test.NewNode captured the
// concrete testing type at the two entry points above, so the whole
// strategy/operation/execution tree below works on the type-erased
// test.Node (which is itself a testing.TB) and process* stay plain methods.
func runE2EStrategy(t test.Node) {
	if cfg == nil || cfg.Strategies == nil {
		t.Fatal("test setting or strategies is nil, please add test configuration yaml file by -config option")
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var err error
	r := new(runner)
	if cfg.Kubernetes != nil {
		r.k8s, err = k8s.New(k8s.WithKubeConfigPath(cfg.Kubernetes.KubeConfig))
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
	r.clients = resource.NewClients(r.k8s)
	if cfg.Target != nil {
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
		// Register the Stop defer right after Start succeeds so that a failure
		// in the following setup (e.g. agent.New -> t.Fatalf) cannot leak the
		// started client and its background goroutines.
		defer func() {
			err = r.client.Stop(ctx)
			if err != nil {
				t.Errorf("failed to stop client: %v", err)
			}
		}()
		r.aclient, err = agent.New(agent.WithValdClient(r.client))
		if err != nil {
			t.Fatalf("failed to create agent client: %v", err)
		}
		errgroup.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case err := <-ech:
				if err != nil {
					t.Errorf("client daemon returned error: %v", err)
				}
			}
			return nil
		})
		t.Logf("connected addrs: %v", r.client.GRPCClient().ConnectedAddrs(ctx))
	} else {
		// scenarios such as operator verification only use kubernetes/http operations and do not require a gRPC target.
		t.Log("gRPC target is not configured, skipping gRPC client setup")
	}
	t.Run("Run E2E V2 Scenarios", func(tt test.Node) {
		if err := executeWithTimings(tt, ctx, cfg, cfg.FilePath, "e2e", false, func(ttt test.Node, ctx context.Context) error {
			ttt.Helper()
			for i, st := range cfg.Strategies {
				col := r.processStrategy(ttt, ctx, i, st)
				if cfg.Metrics != nil && cfg.Metrics.Enabled && cfg.Collector != nil && col != nil {
					cfg.Strategies[i].Collector = col
					if err := col.MergeInto(cfg.Collector); err != nil {
						ttt.Errorf("failed to merge strategy collector: %v", err)
					}
				}
			}
			return nil
		}); err != nil {
			tt.Errorf("failed to process operations: %v", err)
		}
		if cfg.Metrics != nil && cfg.Metrics.Enabled && cfg.Collector != nil {
			snapshot := cfg.Collector.GlobalSnapshot()
			log.Infof("Global Metrics for %s:\n%s", cfg.FilePath, snapshot)
		}
	})
}

func (r *runner) processStrategy(
	t test.Node, ctx context.Context, idx int, st *config.Strategy,
) (col metrics.Collector) {
	t.Helper()
	if r == nil || st == nil {
		return nil
	}
	col = st.Collector
	t.Run(fmt.Sprintf("#%d: strategy=%s", idx, st.Name), func(tt test.Node) {
		if err := executeWithTimings(tt, ctx, st, st.Name, "strategy", false, func(ttt test.Node, ctx context.Context) error {
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
						c := r.processOperation(ttt, egctx, st.Name, i, op)
						if st.Metrics != nil && st.Metrics.Enabled && col != nil && c != nil {
							st.Operations[i].Collector = c
							if err := c.MergeInto(col); err != nil {
								ttt.Logf("failed to merge operation for collector: %v and %v error: %v", col, c, err)
							}
						}
						return nil
					})
				}
			}
			return eg.Wait()
		}); err != nil {
			tt.Errorf("failed to process operations: %v", err)
		}
		if st.Metrics != nil && st.Metrics.Enabled && col != nil {
			snapshot := col.GlobalSnapshot()
			log.Infof("Strategy Metrics for %s:\n%s", st.Name, snapshot)
		}
	})
	return col
}

func (r *runner) processOperation(
	t test.Node, ctx context.Context, strategyName string, idx int, op *config.Operation,
) (col metrics.Collector) {
	t.Helper()
	if r == nil || op == nil {
		return nil
	}
	col = op.Collector
	t.Run(fmt.Sprintf("#%d: operation=%s", idx, op.Name), func(tt test.Node) {
		if err := executeWithTimings(tt, ctx, op, op.Name, "operation", false, func(ttt test.Node, ctx context.Context) error {
			ttt.Helper()
			for i, e := range op.Executions {
				c := r.processExecution(ttt, ctx, strategyName, op.Name, i, e)
				if op.Metrics != nil && op.Metrics.Enabled && col != nil && c != nil {
					op.Executions[i].Collector = c
					if err := c.MergeInto(col); err != nil {
						ttt.Errorf("failed to merge execution collector: %v", err)
					}
				}
			}
			return nil
		}); err != nil {
			tt.Errorf("failed to process operation: %v", err)
		}
		if op.Metrics != nil && op.Metrics.Enabled && col != nil {
			snapshot := col.GlobalSnapshot()
			log.Infof("Operation Metrics for %s/%s:\n%s", strategyName, op.Name, snapshot)
			logRecallAndQPS(tt, fmt.Sprintf("%s/%s", strategyName, op.Name), col)
		}
	})
	return col
}

func (r *runner) processExecution(
	t test.Node, ctx context.Context, strategyName, opName string, idx int, e *config.Execution,
) (col metrics.Collector) {
	t.Helper()
	if r == nil || e == nil {
		return nil
	}

	t.Run(fmt.Sprintf("#%d: execution=%s type=%s mode=%s", idx, e.Name, e.Type, e.Mode), func(tt test.Node) {
		if err := executeWithTimings(tt, ctx, e, e.Name, "execution", true, func(ttt test.Node, ctx context.Context) error {
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
				if r.client == nil {
					ttt.Fatalf("gRPC client is not initialized, target configuration is required for %s operation", e.Type)
				}
				train, test, neighbors := getDatasetSlices(ttt, e)
				// getDatasetSlices returns nil cycles when no dataset is available
				// (Bind allows scenarios without a top-level dataset, e.g. operator
				// verification). These operations all consume the dataset, so fail
				// cleanly here instead of letting a nil cycle panic deep in the gRPC
				// helpers and abort the whole test binary.
				if train == nil {
					ttt.Fatalf("dataset is not initialized; a top-level dataset configuration is required for %s operation", e.Type)
				}
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
					// processSearch takes the query-vector cycle first (test) and the
					// indexed-vector cycle second (train); passing them swapped made
					// every Search/LinearSearch query a train vector, so recall@k
					// against the dataset's test ground truth was structurally zero.
					return r.processSearch(ttt, ctx, test, train, neighbors, e)
				case config.OpInsert,
					config.OpUpdate,
					config.OpUpsert,
					config.OpRemove,
					config.OpRemoveByTimestamp:
					return r.processModification(ttt, ctx, train, e)
				case config.OpObject,
					config.OpListObject,
					config.OpTimestamp,
					config.OpExists:
					return r.processObject(ttt, ctx, train, e)
				}
			case config.OpIndexInfo,
				config.OpIndexDetail,
				config.OpIndexStatistics,
				config.OpIndexStatisticsDetail,
				config.OpIndexProperty,
				config.OpFlush:
				if r.client == nil {
					ttt.Fatalf("gRPC client is not initialized, target configuration is required for %s operation", e.Type)
				}
				start := time.Now()
				log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d",
					e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx)
				defer func() {
					log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d",
						e.Name, time.Since(start).String(), e.Type, e.Mode, idx)
				}()
				return r.processIndex(ttt, ctx, e)
			case config.OpCreateIndex,
				config.OpSaveIndex,
				config.OpCreateAndSaveIndex:
				if r.aclient == nil {
					ttt.Fatalf("agent gRPC client is not initialized, target configuration is required for %s operation", e.Type)
				}
				start := time.Now()
				log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d",
					e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx)
				defer func() {
					log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d",
						e.Name, time.Since(start).String(), e.Type, e.Mode, idx)
				}()
				return r.processAgent(ttt, ctx, e)
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
			case config.OpHTTP:
				if e.HTTP != nil {
					start := time.Now()
					log.Infof("started %s execution at %s, type: %s, mode: %s, execution: %d, http method: %s, url: %s",
						e.Name, start.Format("2006-01-02 15:04:05"), e.Type, e.Mode, idx, e.HTTP.Method, e.HTTP.URL)
					defer func() {
						log.Infof("finished %s execution in %s, type: %s, mode: %s, execution: %d, http method: %s, url: %s",
							e.Name, time.Since(start).String(), e.Type, e.Mode, idx, e.HTTP.Method, e.HTTP.URL)
					}()
					r.processHTTP(ttt, ctx, e)
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
			tt.Errorf("failed to process execution %d %s (type: %s, mode: %s): %v", idx, e.Name, e.Type, e.Mode, err)
		}
		if e.Metrics != nil && e.Metrics.Enabled && e.Collector != nil {
			snapshot := e.Collector.GlobalSnapshot()
			log.Infof("Execution Metrics for %s/%s/%s:\n%s", strategyName, opName, e.Name, snapshot)
			logRecallAndQPS(tt, fmt.Sprintf("%s/%s/%s", strategyName, opName, e.Name), e.Collector)
		}
	})
	return e.Collector
}

func executeWithTimings[T interface {
	config.Timing
	config.Repeater
}](
	t test.Node,
	ctx context.Context,
	cfg T,
	name, prefix string,
	measured bool,
	fn func(test.Node, context.Context) error,
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

	var timeoutDur time.Duration
	if timeout := cfg.GetTimeout(); timeout != "" {
		dur, terr := timeout.Duration()
		if terr != nil {
			t.Errorf("failed to parse timeout duration: %s, error: %v", timeout, terr)
		}
		timeoutDur = dur
	}

	if measured && test.IsBenchmark(t) {
		// Benchmark mode: one measured iteration = one full configured
		// execution pass (including its repeats). test.Measured drives
		// test.Loop (b.Loop), which confines the benchmark timer to the loop
		// body, so the delay above and the wait below stay unmeasured, gives
		// every iteration a fresh timeout window and joins per-iteration
		// errors. Only the execution level passes measured=true:
		// strategy/operation levels are grouping nodes (they call Run) and
		// must not loop.
		if timeoutDur > 0 {
			t.Logf("timeout is set to %s, each benchmark iteration of this %s/%s will stop after %s", timeoutDur, prefix, name, timeoutDur.String())
		}
		err = test.Measured(ctx, t, timeoutDur, func(ictx context.Context) error {
			return executeWithRepeats(t, ictx, name, prefix, cfg.GetRepeats(), fn)
		})
	} else {
		// The timeout bounds only the execution body: the trailing wait below
		// must keep the parent ctx, or an execution that succeeds close to its
		// deadline gets its cool-down wait cut short and reports a spurious
		// context.DeadlineExceeded instead of its real (possibly nil) result.
		ectx := ctx
		if timeoutDur > 0 {
			t.Logf("timeout is set to %s, this %s/%s will stop after %s", timeoutDur, prefix, name, timeoutDur.String())
			var cancel context.CancelFunc
			ectx, cancel = context.WithTimeout(ctx, timeoutDur)
			defer cancel()
		}
		err = executeWithRepeats(t, ectx, name, prefix, cfg.GetRepeats(), fn)
	}

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

// executeWithRepeats adapts the scenario's declarative repeat configuration
// to retry.Do — see internal/test/retry for the exact per-mode exit and
// error semantics (including ModeSuccess doubling as a bounded wait, which
// max_vector_dim.yaml's ResourceExhausted branch relies on) — and threads
// per-attempt logging through Policy.OnRetry.
func executeWithRepeats(
	t test.Node,
	ctx context.Context,
	name, prefix string,
	repeats *config.Repeats,
	fn func(test.Node, context.Context) error,
) error {
	t.Helper()
	task := fmt.Sprintf("%s for %s", prefix, name)
	policy := retry.Policy{
		OnRetry: func(attempt uint64, err error) {
			log.Warnf("failed to finish %s (attempt %d), error: %v, will retry", task, attempt, err)
		},
	}
	if repeats != nil && repeats.Enabled {
		task = fmt.Sprintf("Repeat %s, ExitCondition: %s", task, repeats.ExitCondition)
		switch repeats.ExitCondition {
		case config.Success:
			policy.Mode = retry.ModeSuccess
		case config.Count:
			policy.Mode = retry.ModeCount
			policy.Count = repeats.Count
		case config.Timeout:
			policy.Mode = retry.ModeTimeout
		default:
			// Fail loudly instead of silently degrading to a single attempt
			// (retry.ModeOnce is the Policy zero value) when a scenario
			// enables repeats but omits or misspells exit_condition.
			t.Errorf("unknown repeats exit_condition %q for %s, ignoring the repeat configuration", repeats.ExitCondition, task)
		}
		if wait := repeats.Interval; wait != "" {
			dur, werr := wait.Duration()
			if werr != nil {
				t.Errorf("failed to parse repeat interval: %s, error: %v", wait, werr)
				return nil
			}
			policy.Interval = dur
		}
	}
	var attempts uint64
	err := retry.Do(ctx, policy, func(ictx context.Context) error {
		attempts++
		log.Infof("%s (attempt %d)", task, attempts)
		return fn(t, ictx)
	})
	if policy.Mode != retry.ModeOnce {
		log.Infof("finished %s after %d attempt(s), error: %v", task, attempts, err)
	}
	return err
}

func newClient(
	t testing.TB, ctx context.Context, meta map[string]string,
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
