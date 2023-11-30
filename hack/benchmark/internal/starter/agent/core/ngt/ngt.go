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

// Package ngt provides ngt agent starter  functionality
package ngt

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/core/ngt/config"
	"github.com/vdaas/vald/pkg/agent/core/ngt/usecase"
)

const name = "agent-ngt"

type server struct {
	cfg    *config.Data
	client vald.Client
}

func New(opts ...Option) starter.Starter {
	srv := new(server)
	for _, opt := range append(defaultOptions, opts...) {
		opt(srv)
	}
	return srv
}

func (s *server) Run(ctx context.Context, tb testing.TB) func() {
	tb.Helper()

	info.Init(name)

	daemon, err := usecase.New(s.cfg)
	if err != nil {
		tb.Fatal(err)
	}

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := runner.Run(ctx, daemon, name)
		if err != nil {
			tb.Fatalf("agent runner returned error %s", err.Error())
		}
	}()

	time.Sleep(5 * time.Second)

	ech, err := s.client.Start(ctx)
	if err != nil {
		tb.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-ech:
				tb.Error(err)
			}
		}
	}()

	return func() {
		cancel()
		s.client.Stop(ctx)
		wg.Wait()
	}
}
