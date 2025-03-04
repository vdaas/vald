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
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestE2EStrategy(t *testing.T) {
	ctx := t.Context()
	for i, st := range cfg.Strategies {
		if st != nil {
			t.Run(fmt.Sprintf("#%d: strategy=%s", i, st.Name), func(tt *testing.T) {
				if st.Delay != "" {
					select {
					case <-ctx.Done():
						t.Fatal(ctx.Err())
					case <-time.After(st.Delay.DurationWithDefault(0)):
					}
				}
				if st.Timeout != "" {
					dur, err := st.Timeout.Duration()
					if err == nil && dur != 0 {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, dur)
						defer cancel()
					}
				}
				var eg errgroup.Group
				eg, ctx = errgroup.New(ctx)
				if st.Concurrency > 0 {
					eg.SetLimit(int(st.Concurrency))
				}
				for j, op := range st.Operations {
					eg.Go(safety.RecoverFunc(func() error {
						tt.Run(fmt.Sprintf("#%d: operation=%s", j, op.Name), func(ttt *testing.T) {
							if op.Delay != "" {
								select {
								case <-ctx.Done():
									t.Fatal(ctx.Err())
								case <-time.After(op.Delay.DurationWithDefault(0)):
								}
							}
							if op.Timeout != "" {
								dur, err := op.Timeout.Duration()
								if err == nil && dur != 0 {
									var cancel context.CancelFunc
									ctx, cancel = context.WithTimeout(ctx, dur)
									defer cancel()
								}
							}

							for k, e := range op.Executions {
								ttt.Run(fmt.Sprintf("#%d: execution=%s", k, e.Name), func(tttt *testing.T) {
									if e.Delay != "" {
										select {
										case <-ctx.Done():
											t.Fatal(ctx.Err())
										case <-time.After(e.Delay.DurationWithDefault(0)):
										}
									}
									if e.Timeout != "" {
										dur, err := e.Timeout.Duration()
										if err == nil && dur != 0 {
											var cancel context.CancelFunc
											ctx, cancel = context.WithTimeout(ctx, dur)
											defer cancel()
										}
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
									case config.OpIndexInfo:
									case config.OpIndexProperty:
									case config.OpKubernetes:
									case config.OpClient:
									case config.OpWait:
									default:
										tttt.Fatalf("unsupported operation type: %s", e.Type)
									}
									if e.Wait != "" {
										select {
										case <-ctx.Done():
											tttt.Fatal(ctx.Err())
										case <-time.After(e.Wait.DurationWithDefault(0)):
										}
									}
								})
							}

							if op.Wait != "" {
								select {
								case <-ctx.Done():
									ttt.Fatal(ctx.Err())
								case <-time.After(op.Wait.DurationWithDefault(0)):
								}
							}
						})
						return nil
					}))
				}
				eg.Wait()
				if st.Wait != "" {
					select {
					case <-ctx.Done():
						tt.Fatal(ctx.Err())
					case <-time.After(st.Wait.DurationWithDefault(0)):
					}
				}
			})
		}
	}
}
