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

func serverStarter(b *testing.B, hotReloadEnabled bool) (ctx context.Context, stop context.CancelFunc, addr string) {
	b.Helper()
	// TODO start server and return close func
	ctx, stop = context.WithCancel(b.Context())

	// Bind to an ephemeral port to obtain a free port, close it, then start the server on that port.
	ln, err := net.Listen(net.TCP.String(), "127.0.0.1:0")
	if err != nil {
		b.Fatalf("listen: %v", err)
	}
	host, port, _ := net.SplitHostPort(ln.Addr().String())
	_ = ln.Close()
	addr = net.JoinHostPort(host, port)

	srv, err := starter.New(
		starter.WithConfig((&config.Servers{
			TLS: func() *config.TLS {
				if hotReloadEnabled {
					return &config.TLS{
						// write hot reload version here
						Cert:               test.GetTestdataPath("tls/server.crt"),
						Key:                test.GetTestdataPath("tls/server.key"),
						ClientAuth:         "noclientcert",
						InsecureSkipVerify: false,
						ServerName:         "vald.vdaas.org",
						// HotReload:          true,
					}
				}
				return &config.TLS{
					// write something here
					Cert:               test.GetTestdataPath("tls/server.crt"),
					Key:                test.GetTestdataPath("tls/server.key"),
					ClientAuth:         "noclientcert",
					InsecureSkipVerify: false,
					ServerName:         "vald.vdaas.org",
					// HotReload:        false,
				}
			}(),
			// configure this
			Servers: []*config.Server{{
				Name: "bench-grpc",
				Mode: server.GRPC.String(),
				Host: host,
				Port: port,
				GRPC: &config.GRPC{},
			}},
		}).Bind()), // set server configuration here
		starter.WithGRPC(func(sc *config.Server) []server.Option {
			// Build server TLS via internal/tls (toggle hot-reload here since config.TLS does not expose the flag).
			stls, serr := tls.NewServerConfig(
				tls.WithCert(test.GetTestdataPath("tls/server.crt")),
				tls.WithKey(test.GetTestdataPath("tls/server.key")),
				tls.WithClientAuth("noclientcert"),
				tls.WithServerCertHotReload(hotReloadEnabled),
			)
			opts := []server.Option{
				server.WithGRPCRegistFunc(func(s *grpc.Server) {
					vald.RegisterIndexServer(s, new(mockIndexInfoServer))
				}),
			}
			if serr == nil {
				opts = append(opts, server.WithGRPCOption(grpc.Creds(credentials.NewTLS(stls))))
			} else {
				// Keep server running without TLS if building TLS failed (adjust to b.Error if desired).
				b.Logf("failed to build TLS config for server: %v", serr)
			}
			return opts
		}),
	)
	if err != nil {
		b.Error(err)
	}
	_ = srv.ListenAndServe(ctx)

	return ctx, stop, addr
}

func reloadTLSCerts(b *testing.B) (stop context.CancelFunc) {
	b.Helper()

	var ctx context.Context
	ctx, stop = context.WithCancel(b.Context())
	eg, egctx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		tick := time.NewTicker(time.Second * 5)
		// TODO: reload certs periodically
		for {
			select {
			case <-egctx.Done():
				return nil
			case <-tick.C:
				// reload here using internal/file package
				// Read server.crt and server2.crt alternately, and overwrite the same destination path (no rename/copy).
				dst := test.GetTestdataPath("tls/server.crt")
				srcA := test.GetTestdataPath("tls/server.crt")
				srcB := test.GetTestdataPath("tls/server2.crt")
				if !file.Exists(srcB) {
					continue
				}
				next := srcB
				if cur, err := file.ReadFile(dst); err == nil {
					if a, err := file.ReadFile(srcA); err == nil && bytes.Equal(cur, a) {
						next = srcB
					} else {
						next = srcA
					}
				}
				if bs, err := file.ReadFile(next); err == nil {
					_, _ = file.WriteFile(egctx, dst, bytes.NewReader(bs), 0o600)
				}
			}
		}
	}))

	return stop
}

type mockIndexInfoServer struct {
	vald.UnimplementedIndexServer
}

func (m mockIndexInfoServer) IndexInfo(ctx context.Context, _ *payload.Empty) (*payload.Info_Index_Count, error) {
	return &payload.Info_Index_Count{
		Stored: 100,
	}, nil
}