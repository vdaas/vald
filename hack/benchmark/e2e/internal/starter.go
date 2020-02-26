//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package internal

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	"github.com/vdaas/vald/pkg/agent/ngt/usecase"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

const (
	configuration = `
version: v0.0.0
server_config:
  servers:
  - name: agent-grpc
    host: 127.0.0.1
    port: 8082
    mode: GRPC
    probe_wait_time: "0s"
    http:
      shutdown_duration: "0s"
      handler_timeout: ""
      idle_timeout: ""
      read_header_timeout: ""
      read_timeout: ""
      write_timeout: ""
  - name: agent-rest
    host: 127.0.0.1
    port: 8081
    mode: REST
    probe_wait_time: "0s"
    http:
      shutdown_duration: "0s"
      handler_timeout: "60s"
      idle_timeout: "60s"
      read_header_timeout: "60s"
      read_timeout: "60s"
      write_timeout: "60s"
  startup_strategy:
  - agent-grpc
  - agent-rest
  shutdown_strategy:
  - agent-grpc
  - agent-rest
  full_shutdown_duration: 600s
  tls:
    enabled: false
ngt:
  dimension: 0
  bulk_insert_chunk_size: 10
  distance_type: unknown
  object_type: unknown
  creation_edge_size: 20
  search_edge_size: 10
  enable_in_memory_mode: true
`
)

var (
	baseCfg config.Data
	once    sync.Once
)

type Option func(*config.Data) error

func WithCreationEdgeSize(creationEdgeSize int) Option {
	return func(cfg *config.Data) error {
		cfg.NGT.CreationEdgeSize = creationEdgeSize
		return nil
	}
}

func WithSearchEdgeSize(searchEdgeSize int) Option {
	return func(cfg *config.Data) error {
		cfg.NGT.SearchEdgeSize = searchEdgeSize
		return nil
	}
}

func withHost(host, typ string) Option {
	return func(cfg *config.Data) error {
		for _, svr := range cfg.Server.Servers {
			if svr.Mode == typ {
				svr.Host = host
			}
		}
		return nil
	}
}

func withPort(port uint, typ string) Option {
	return func(cfg *config.Data) error {
		for _, svr := range cfg.Server.Servers {
			if svr.Mode == typ {
				svr.Port = port
			}
		}
		return nil
	}
}

func WithGRPCHost(host string) Option {
	return withHost(host, "GRPC")
}

func WithGRPCPort(port uint) Option {
	return withPort(port, "GRPC")
}

func WithRESTHost(host string) Option {
	return withHost(host, "REST")
}

func WithRESTPort(port uint) Option {
	return withPort(port, "REST")
}

func StartAgentNGTServer(tb testing.TB, ctx context.Context, d assets.Dataset, opts ...Option) {
	tb.Helper()

	once.Do(func() {
		sr := strings.NewReader(configuration)
		err := yaml.NewDecoder(sr).Decode(&baseCfg)
		if err != nil {
			tb.Errorf("failed to load config %s \t %s", d.Name(), err.Error())
		}
	})
	cfg := baseCfg
	cfg.NGT.Dimension = d.Dimension()
	cfg.NGT.DistanceType = d.DistanceType()
	cfg.NGT.ObjectType = d.ObjectType()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			tb.Error(err)
		}
	}

	daemon, err := usecase.New(&cfg)
	if err != nil {
		tb.Errorf("failed create daemon %s", err.Error())
		return
	}

	go func() {
		err = runner.Run(errgroup.Init(ctx), daemon, "agent-ngt")
		if err != nil {
			tb.Errorf("agent runner returned error %s", err.Error())
		}
	}()
	time.Sleep(time.Second * 5)
}

func NewAgentClient(tb testing.TB, ctx context.Context, address string) agent.AgentClient {
	tb.Helper()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		tb.Errorf("failed to connect %s \t %s", address, err.Error())
		return nil
	}
	return agent.NewAgentClient(conn)
}

func NewValdClient(tb testing.TB, ctx context.Context, address string) vald.ValdClient {
	tb.Helper()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		tb.Errorf("failed to connect %s \t %s", address, err.Error())
		return nil
	}
	return vald.NewValdClient(conn)
}
