package tls

import (
    "bufio"
    "crypto/tls"
    "io"
    "net"
    "os"
    "path/filepath"
    "sync"
    "sync/atomic"
    "testing"
    "time"

    ilog "github.com/vdaas/vald/internal/log"
    "github.com/vdaas/vald/internal/log/level"
    testdata "github.com/vdaas/vald/internal/test"
)

func init() {
    // Reduce benchmark log noise: suppress DEBUG/WARN logs during bench runs.
    ilog.Init(ilog.WithLevel(level.ERROR.String()))
}

// startTLSEchoServer starts a simple TLS echo server and returns its address and a shutdown function.
func startTLSEchoServer(b *testing.B, scfg *tls.Config) (addr string, shutdown func()) {
    l, err := tls.Listen("tcp", "127.0.0.1:0", scfg)
    if err != nil {
        b.Fatalf("failed to listen: %v", err)
    }
    stop := make(chan struct{})
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            conn, err := l.Accept()
            if err != nil {
                select {
                case <-stop:
                    return
                default:
                }
                // transient error
                continue
            }
            go func(c net.Conn) {
                defer c.Close()
                br := bufio.NewReader(c)
                bw := bufio.NewWriter(c)
                buf := make([]byte, 32<<10) // up to 32KiB per read
                for {
                    c.SetDeadline(time.Now().Add(10 * time.Second))
                    n, err := br.Read(buf)
                    if n > 0 {
                        if _, werr := bw.Write(buf[:n]); werr != nil {
                            return
                        }
                        if werr := bw.Flush(); werr != nil {
                            return
                        }
                    }
                    if err != nil {
                        if err == io.EOF {
                            return
                        }
                        // read error or timeout
                        return
                    }
                }
            }(conn)
        }
    }()
    return l.Addr().String(), func() {
        close(stop)
        _ = l.Close()
        wg.Wait()
    }
}

func runEchoBench(b *testing.B, hot bool) {
    // Prepare server config
    scfg, err := NewServerConfig(
        WithCert(testdata.GetTestdataPath("tls/server.crt")),
        WithKey(testdata.GetTestdataPath("tls/server.key")),
        WithServerCertHotReload(hot),
    )
    if err != nil {
        b.Fatalf("server config: %v", err)
    }
    addr, shutdown := startTLSEchoServer(b, scfg)
    defer shutdown()

    // Prepare client config trusting the server certificate via CA
    ccfg, err := NewClientConfig(
        WithCa(testdata.GetTestdataPath("tls/ca.pem")),
    )
    if err != nil {
        b.Fatalf("client config: %v", err)
    }

    payload := []byte("0123456789abcdef0123456789abcdef")

    // Pre-establish connections per worker in RunParallel.
    // We establish inside the goroutine to bind one conn per worker.
    var totalNanos int64

    b.ReportAllocs()
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        conn, err := tls.Dial("tcp", addr, ccfg)
        if err != nil {
            b.Fatalf("dial: %v", err)
        }
        defer conn.Close()
        r := bufio.NewReader(conn)
        w := bufio.NewWriter(conn)
        buf := make([]byte, len(payload))
        for pb.Next() {
            start := time.Now()
            if _, err := w.Write(payload); err != nil {
                b.Fatalf("write: %v", err)
            }
            if err := w.Flush(); err != nil {
                b.Fatalf("flush: %v", err)
            }
            if _, err := io.ReadFull(r, buf); err != nil {
                b.Fatalf("read: %v", err)
            }
            atomic.AddInt64(&totalNanos, time.Since(start).Nanoseconds())
        }
    })
    b.StopTimer()

    // Report average latency and QPS.
    if b.N > 0 {
        avg := float64(totalNanos) / float64(b.N) // ns/op
        qps := 1e9 / avg
        b.ReportMetric(avg, "ns/op_avg")
        b.ReportMetric(qps, "qps")
    }
}

func Benchmark_EchoTLS_Static(b *testing.B)    { runEchoBench(b, false) }
func Benchmark_EchoTLS_HotReload(b *testing.B) { runEchoBench(b, true) }

// --- Handshake-included benchmarks (reconnect per op) ---

func runEchoBenchReconnect(b *testing.B, hot bool) {
    // Server
    scfg, err := NewServerConfig(
        WithCert(testdata.GetTestdataPath("tls/server.crt")),
        WithKey(testdata.GetTestdataPath("tls/server.key")),
        WithServerCertHotReload(hot),
    )
    if err != nil {
        b.Fatalf("server config: %v", err)
    }
    addr, shutdown := startTLSEchoServer(b, scfg)
    defer shutdown()

    // Client
    ccfg, err := NewClientConfig(WithCa(testdata.GetTestdataPath("tls/ca.pem")))
    if err != nil {
        b.Fatalf("client config: %v", err)
    }

    payload := []byte("0123456789abcdef0123456789abcdef")
    var totalNanos int64

    b.ReportAllocs()
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        buf := make([]byte, len(payload))
        for pb.Next() {
            start := time.Now()
            conn, err := tls.Dial("tcp", addr, ccfg) // handshake per op
            if err != nil {
                b.Fatalf("dial: %v", err)
            }
            r := bufio.NewReader(conn)
            w := bufio.NewWriter(conn)
            if _, err := w.Write(payload); err != nil {
                b.Fatalf("write: %v", err)
            }
            if err := w.Flush(); err != nil {
                b.Fatalf("flush: %v", err)
            }
            if _, err := io.ReadFull(r, buf); err != nil {
                b.Fatalf("read: %v", err)
            }
            _ = conn.Close()
            atomic.AddInt64(&totalNanos, time.Since(start).Nanoseconds())
        }
    })
    b.StopTimer()

    if b.N > 0 {
        avg := float64(totalNanos) / float64(b.N)
        qps := 1e9 / avg
        b.ReportMetric(avg, "ns/op_avg")
        b.ReportMetric(qps, "qps")
    }
}

func Benchmark_EchoTLS_HandshakePerOp_Static(b *testing.B)    { runEchoBenchReconnect(b, false) }
func Benchmark_EchoTLS_HandshakePerOp_HotReload(b *testing.B) { runEchoBenchReconnect(b, true) }

// --- Rotation-under-load benchmark ---

// copyFile copies src to dst.
func copyFile(src, dst string) error {
    in, err := os.ReadFile(src)
    if err != nil {
        return err
    }
    dir := filepath.Dir(dst)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return err
    }
    tmp := dst + ".tmp"
    if err := os.WriteFile(tmp, in, 0o644); err != nil {
        return err
    }
    return os.Rename(tmp, dst)
}

// atomicReplace writes data to path atomically via temp+rename.
func atomicReplace(path string, data []byte) error {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return err
    }
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, data, 0o644); err != nil {
        return err
    }
    return os.Rename(tmp, path)
}

// startCertReloader toggles active cert file between valid and invalid
// to simulate rotation and short error windows. Returns stop func.
func startCertReloader(b *testing.B, activeCert, activeKey string) (stop func()) {
    validCert := testdata.GetTestdataPath("tls/server.crt")
    validKey := testdata.GetTestdataPath("tls/server.key")
    invalidCert := testdata.GetTestdataPath("tls/invalid-server.crt")

    // Prepare initial valid files
    if err := copyFile(validCert, activeCert); err != nil {
        b.Fatalf("copy cert: %v", err)
    }
    if err := copyFile(validKey, activeKey); err != nil {
        b.Fatalf("copy key: %v", err)
    }

    done := make(chan struct{})
    go func() {
        tick := time.NewTicker(50 * time.Millisecond)
        defer tick.Stop()
        for {
            select {
            case <-done:
                return
            case <-tick.C:
                // 1) brief invalid window
                if data, err := os.ReadFile(invalidCert); err == nil {
                    _ = atomicReplace(activeCert, data)
                }
                // key stays valid; LoadX509KeyPair will fail, server falls back to last good cert
                time.Sleep(5 * time.Millisecond)
                // 2) restore valid cert
                if data, err := os.ReadFile(validCert); err == nil {
                    _ = atomicReplace(activeCert, data)
                }
            }
        }
    }()
    return func() { close(done) }
}

func runEchoBenchRotation(b *testing.B) {
    dir := b.TempDir()
    activeCert := filepath.Join(dir, "active.crt")
    activeKey := filepath.Join(dir, "active.key")
    stopRotate := startCertReloader(b, activeCert, activeKey)
    defer stopRotate()

    // Server with hot reload pointed at active files
    scfg, err := NewServerConfig(
        WithCert(activeCert),
        WithKey(activeKey),
        WithServerCertHotReload(true),
    )
    if err != nil {
        b.Fatalf("server config: %v", err)
    }
    addr, shutdown := startTLSEchoServer(b, scfg)
    defer shutdown()

    // Client trusts CA
    ccfg, err := NewClientConfig(WithCa(testdata.GetTestdataPath("tls/ca.pem")))
    if err != nil {
        b.Fatalf("client config: %v", err)
    }

    payload := []byte("0123456789abcdef0123456789abcdef")
    var totalNanos int64

    b.ReportAllocs()
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        conn, err := tls.Dial("tcp", addr, ccfg)
        if err != nil {
            b.Fatalf("dial: %v", err)
        }
        defer conn.Close()
        r := bufio.NewReader(conn)
        w := bufio.NewWriter(conn)
        buf := make([]byte, len(payload))
        for pb.Next() {
            start := time.Now()
            if _, err := w.Write(payload); err != nil {
                b.Fatalf("write: %v", err)
            }
            if err := w.Flush(); err != nil {
                b.Fatalf("flush: %v", err)
            }
            if _, err := io.ReadFull(r, buf); err != nil {
                b.Fatalf("read: %v", err)
            }
            atomic.AddInt64(&totalNanos, time.Since(start).Nanoseconds())
        }
    })
    b.StopTimer()

    if b.N > 0 {
        avg := float64(totalNanos) / float64(b.N)
        qps := 1e9 / avg
        b.ReportMetric(avg, "ns/op_avg")
        b.ReportMetric(qps, "qps")
    }
}

func Benchmark_EchoTLS_HotReload_WithRotation(b *testing.B) { runEchoBenchRotation(b) }
