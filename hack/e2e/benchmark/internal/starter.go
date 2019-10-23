//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	"github.com/vdaas/vald/pkg/agent/ngt/usecase"
	"google.golang.org/grpc"
)

func StartAgentNGTServer(tb testing.TB, ctx context.Context, path string) {
	tb.Helper()
	cfg, err := config.NewConfig(path)
	if err != nil {
		tb.Errorf("failed to load config %s \t %s", path, err.Error())
		return
	}

	daemon, err := usecase.New(cfg)
	if err != nil {
		tb.Errorf("failed create daemon %s", err.Error())
		return
	}

	go func() {
		err = runner.Run(errgroup.Init(ctx), daemon)
		if err != nil {
			tb.Errorf("agent runnner returned error %s", err.Error())
		}
	}()
	time.Sleep(time.Second * 5)
}

func NewAgentClient(tb testing.TB, ctx context.Context, host string, port int) agent.AgentClient {
	tb.Helper()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		tb.Errorf("failed to connect %s:%d \t %s", host, port, err.Error())
		return nil
	}
	return agent.NewAgentClient(conn)
}
