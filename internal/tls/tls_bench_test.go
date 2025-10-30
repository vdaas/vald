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
//
// Benchmark can be run by
// make certs/gen
// go test ./internal/tls -run=^$ -bench=. -benchmem
//
package tls_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	localhost = "127.0.0.1"
)

var (
	activeCertPath string
	activeKeyPath  string
)

func init() {
	log.Init(log.WithLevel(level.ERROR.String()))
}

func serverStarter(
	b *testing.B, hot bool,
) (ctx context.Context, stop context.CancelFunc, addr string) {
	b.Helper()
	ctx, stop = context.WithCancel(b.Context())

	// Get a free port by listening on port 0 and closing the listener immediately
	ln, err := net.Listen(net.TCP.String(), net.JoinHostPort(localhost, 0))
	if err != nil {
		b.Fatalf("listen: %v", err)
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		b.Fatalf("split host port: %v", err)
	}
	_ = ln.Close()

	certPath := test.GetTestdataPath("tls/server.crt")
	keyPath := test.GetTestdataPath("tls/server.key")
	if hot {
		dir := b.TempDir()
		activeCertPath = filepath.Join(dir, "active.crt")
		activeKeyPath = filepath.Join(dir, "active.key")
		_, _ = file.CopyFile(ctx, certPath, activeCertPath)
		_, _ = file.CopyFile(ctx, keyPath, activeKeyPath)
		certPath, keyPath = activeCertPath, activeKeyPath
	}

	stls, err := tls.NewServerConfig(
		tls.WithCert(certPath),
		tls.WithKey(keyPath),
		tls.WithServerCertHotReload(hot),
	)
	if err != nil {
		b.Fatalf("server TLS config: %v", err)
	}

	srv, err := starter.New(
		starter.WithConfig((&config.Servers{
			Servers: []*config.Server{{
				Name: "bench-grpc",
				Mode: server.GRPC.String(),
				Host: localhost,
				Port: port,
				GRPC: &config.GRPC{},
			}},
		}).Bind()),
		starter.WithGRPC(func(_ *config.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(gs *grpc.Server) {
					vald.RegisterIndexServer(gs, mockIndexInfoServer{})
				}),
				server.WithGRPCOption(grpc.Creds(credentials.NewTLS(stls))),
			}
		}),
	)
	if err != nil {
		b.Fatalf("starter initialization failed: %v", err)
	}

	if err := srv.ListenAndServe(ctx); err != nil {
		b.Logf("ListenAndServe: %v", err)
	}

	addr = net.JoinHostPort(localhost, port)
	return ctx, stop, addr
}

func reloadTLSCerts(b *testing.B) (stop context.CancelFunc) {
	b.Helper()

	var ctx context.Context
	ctx, stop = context.WithCancel(b.Context())
	eg, egctx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		tick := time.NewTicker(time.Millisecond)
		defer tick.Stop()
		srcA := test.GetTestdataPath("tls/server.crt")
		srcB := test.GetTestdataPath("tls/server2.crt")
		if !file.Exists(srcB) {
			srcB = srcA
		}
		useA := false
		for {
			select {
			case <-egctx.Done():
				return nil
			case <-tick.C:
				if activeCertPath == "" {
					continue
				}
				if useA {
					_, _ = file.CopyFile(egctx, srcA, activeCertPath)
				} else {
					_, _ = file.CopyFile(egctx, srcB, activeCertPath)
				}
				useA = !useA
			}
		}
	}))

	return stop
}

type mockIndexInfoServer struct {
	vald.UnimplementedIndexServer
}

func (_ mockIndexInfoServer) IndexInfo(
	context.Context, *payload.Empty,
) (*payload.Info_Index_Count, error) {
	return &payload.Info_Index_Count{Stored: 100}, nil
}

func runTLSHandshakePerOp(b *testing.B, hot bool) {
	ctx, stop, addr := serverStarter(b, hot)
	defer stop()

	ccfg, err := tls.NewClientConfig(
		tls.WithCa(test.GetTestdataPath("tls/ca.pem")),
	)
	if err != nil {
		b.Fatalf("client tls: %v", err)
	}

	var stopReload context.CancelFunc
	if hot {
		stopReload = reloadTLSCerts(b)
		defer stopReload()
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(ccfg)))
			if err != nil {
				b.Fatalf("dial: %v", err)
			}
			_, err = vald.NewIndexClient(conn).IndexInfo(ctx, &payload.Empty{})
			_ = conn.Close()
			if err != nil {
				b.Fatalf("IndexInfo: %v", err)
			}
		}
	})
}

func BenchmarkTLSHandshakePerOpStatic(b *testing.B)    { runTLSHandshakePerOp(b, false) }
func BenchmarkTLSHandshakePerOpHotReload(b *testing.B) { runTLSHandshakePerOp(b, true) }
