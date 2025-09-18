package tls

import (
    gotls "crypto/tls"
    "runtime"
    "testing"

    testdata "github.com/vdaas/vald/internal/test"
)

var certSink *gotls.Certificate

func Benchmark_GetCertificate_Static(b *testing.B) {
    cfg, err := NewServerConfig(
        WithCert(testdata.GetTestdataPath("tls/server.crt")),
        WithKey(testdata.GetTestdataPath("tls/server.key")),
        WithServerCertHotReload(false),
    )
    if err != nil {
        b.Fatalf("failed to create static server config: %v", err)
    }
    if len(cfg.Certificates) == 0 {
        b.Fatalf("no static certificate loaded")
    }
    get := func() *gotls.Certificate { return &cfg.Certificates[0] }

    _ = get()

    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        certSink = get()
    }
    runtime.KeepAlive(certSink)
}

func Benchmark_GetCertificate_HotReload(b *testing.B) {
    cfg, err := NewServerConfig(
        WithCert(testdata.GetTestdataPath("tls/server.crt")),
        WithKey(testdata.GetTestdataPath("tls/server.key")),
        WithServerCertHotReload(true),
    )
    if err != nil {
        b.Fatalf("failed to create hot-reload server config: %v", err)
    }
    if cfg.GetCertificate == nil {
        b.Fatalf("GetCertificate not set for hot reload config")
    }

    chi := &gotls.ClientHelloInfo{ServerName: "bench"}

    if c, err := cfg.GetCertificate(chi); err == nil {
        certSink = c
    }

    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        c, err := cfg.GetCertificate(chi)
        if err != nil {
            b.Fatalf("GetCertificate failed: %v", err)
        }
        certSink = c
    }
    runtime.KeepAlive(certSink)
}
