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

var (
	activeCertPath string
	activeKeyPath  string
)

func init() {
	log.Init(log.WithLevel(level.ERROR.String()))
}

func serverStarter(b *testing.B, hot bool) (ctx context.Context, stop context.CancelFunc, addr string) {
	b.Helper()
	ctx, stop = context.WithCancel(b.Context())

    ln, err := net.Listen(net.TCP.String(), "127.0.0.1:0")
    if err != nil {
        b.Fatalf("listen: %v", err)
    }
    _, port, _ := net.SplitHostPort(ln.Addr().String())
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
	} else {
		activeCertPath, activeKeyPath = "", ""
	}

	stls, err := tls.NewServerConfig(
		tls.WithCert(certPath),
		tls.WithKey(keyPath),
		tls.WithClientAuth("noclientcert"),
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
				Host: "127.0.0.1",
                    Port: port,
				GRPC: &config.GRPC{},
			}},
		}).Bind()),
		starter.WithGRPC(func(sc *config.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(gs *grpc.Server) {
					vald.RegisterIndexServer(gs, mockIndexInfoServer{})
				}),
				server.WithGRPCOption(grpc.Creds(credentials.NewTLS(stls))),
			}
		}),
	)
	if err != nil {
		b.Error(err)
	}

	go func() { _ = srv.ListenAndServe(ctx) }()

    addr = net.JoinHostPort("127.0.0.1", port)
    deadline := time.Now().Add(3 * time.Second)
    for {
        dctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
        c, err := net.DialContext(dctx, net.TCP.String(), addr)
        cancel()
        if err == nil {
            _ = c.Close()
            break
        }
        if time.Now().After(deadline) {
            break
        }
        time.Sleep(50 * time.Millisecond)
    }
	return ctx, stop, addr
}

func reloadTLSCerts(b *testing.B) (stop context.CancelFunc) {
	b.Helper()

	var ctx context.Context
	ctx, stop = context.WithCancel(b.Context())
	eg, egctx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		tick := time.NewTicker(200 * time.Millisecond)
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

func (m mockIndexInfoServer) IndexInfo(context.Context, *payload.Empty) (*payload.Info_Index_Count, error) {
	return &payload.Info_Index_Count{Stored: 100}, nil
}

func runTLSHandshakePerOp(b *testing.B, hot bool) {
	ctx, stop, addr := serverStarter(b, hot)
	defer stop()

	ccfg, err := tls.NewClientConfig(
		tls.WithCa(test.GetTestdataPath("tls/ca.pem")),
		tls.WithServerName("vald.vdaas.org"),
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
			dctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			conn, err := grpc.DialContext(
				dctx, addr,
				grpc.WithTransportCredentials(credentials.NewTLS(ccfg)),
				grpc.WithBlock(),
				grpc.WithReturnConnectionError(),
			)
			cancel()
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

func Benchmark_TLS_HandshakePerOp_Static(b *testing.B)    { runTLSHandshakePerOp(b, false) }
func Benchmark_TLS_HandshakePerOp_HotReload(b *testing.B) { runTLSHandshakePerOp(b, true) }
