// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package tls_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test"
	"google.golang.org/grpc"
)

func serverStarter(
	b *testing.B, hotReloadEnabled bool,
) (ctx context.Context, stop context.CancelFunc, addr string) {
	b.Helper()
	ctx, stop = context.WithCancel(b.Context())

	srv, err := starter.New(
		starter.WithConfig((&config.Servers{
			TLS: func() *config.TLS {
				baseTLS := &config.TLS{
					Cert:               test.GetTestdataPath("tls/server.crt"),
					Key:                test.GetTestdataPath("tls/server.key"),
					ClientAuth:         "noclientcert",
					InsecureSkipVerify: false,
					ServerName:         "vald.vdaas.org",
				}
				if hotReloadEnabled {
					hot := *baseTLS
					hot.HotReload = true
					return &hot
				}
				return baseTLS
			}(),
			Servers: []*config.Server{{
				Name: "bench-grpc",
				Mode: server.GRPC.String(),
			}},
		}).Bind()),
		starter.WithGRPC(func(sc *config.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(s *grpc.Server) {
					vald.RegisterIndexServer(s, new(mockIndexInfoServer))
				}),
			}
		}),
	)
	if err != nil {
		b.Error(err)
	}
	_ = srv.ListenAndServe(ctx)

	return ctx, stop, "127.0.0.1:50051"
}

func reloadTLSCerts(
	b *testing.B, targets []struct {
		dst  string
		srcs []string
	},
) (stop context.CancelFunc) {
	b.Helper()

	var ctx context.Context
	ctx, stop = context.WithCancel(b.Context())
	eg, egctx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()

		idx := 0
		for {
			select {
			case <-egctx.Done():
				return nil
			case <-tick.C:
				for _, t := range targets {
					if len(t.srcs) == 0 {
						continue
					}
					src := t.srcs[idx%len(t.srcs)]
					bs, err := file.ReadFile(src)
					if err != nil {
						continue
					}
					_, _ = file.OverWriteFile(egctx, t.dst, bytes.NewReader(bs), 0o600)
				}
				idx++
			}
		}
	}))

	return stop
}

type mockIndexInfoServer struct {
	vald.UnimplementedIndexServer
}

func (m mockIndexInfoServer) IndexInfo(
	ctx context.Context, _ *payload.Empty,
) (*payload.Info_Index_Count, error) {
	return &payload.Info_Index_Count{
		Stored: 100,
	}, nil
}
