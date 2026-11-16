// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package net

import (
	"context"
	"crypto/x509"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/tls"
)

func startEchoServer(b *testing.B) (addr string, closeFn func()) {
	b.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("failed to start echo listener: %v", err)
	}

	done := make(chan struct{})
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-done:
					return
				default:
					return
				}
			}
			go func(c net.Conn) {
				defer c.Close()
				buf := make([]byte, 256)
				for {
					n, err := c.Read(buf)
					if err != nil {
						return
					}
					if _, werr := c.Write(buf[:n]); werr != nil {
						return
					}
				}
			}(conn)
		}
	}()

	return ln.Addr().String(), func() {
		close(done)
		_ = ln.Close()
	}
}

// BenchmarkDialer_DialContext measures the latency and allocations of dialing
// with and without DNS caching, for both IP literals and hostnames.
func BenchmarkDialer_DialContext(b *testing.B) {
	addr, closeServer := startEchoServer(b)
	defer closeServer()

	_, port, err := SplitHostPort(addr)
	if err != nil {
		b.Fatalf("failed to split host port: %v", err)
	}
	localhostAddr := JoinHostPort("localhost", port)

	b.Run("IPLiteral_NoDNSCache", func(b *testing.B) {
		der, err := NewDialer(WithDialerTimeout("5s"))
		if err != nil {
			b.Fatalf("NewDialer error: %v", err)
		}
		ctx := context.Background()

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			conn, err := der.DialContext(ctx, "tcp", addr)
			if err != nil {
				b.Fatalf("DialContext error: %v", err)
			}
			_ = conn.Close()
		}
	})

	b.Run("IPLiteral_WithDNSCache", func(b *testing.B) {
		der, err := NewDialer(
			WithEnableDNSCache(),
			WithDNSRefreshDuration("1m"),
			WithDNSCacheExpiration("1m"),
			WithDialerTimeout("5s"),
		)
		if err != nil {
			b.Fatalf("NewDialer error: %v", err)
		}
		ctx := context.Background()

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			conn, err := der.DialContext(ctx, "tcp", addr)
			if err != nil {
				b.Fatalf("DialContext error: %v", err)
			}
			_ = conn.Close()
		}
	})

	b.Run("Hostname_WithDNSCache_CachedHit", func(b *testing.B) {
		der, err := NewDialer(
			WithEnableDNSCache(),
			WithDNSRefreshDuration("1m"),
			WithDNSCacheExpiration("1m"),
			WithDialerTimeout("5s"),
		)
		if err != nil {
			b.Fatalf("NewDialer error: %v", err)
		}
		ctx := context.Background()

		// Warm up first dial so that localhost address is parsed and cached in dnsCache
		warmConn, err := der.DialContext(ctx, "tcp", localhostAddr)
		if err != nil {
			b.Fatalf("warmup dial error: %v", err)
		}
		_ = warmConn.Close()

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			conn, err := der.DialContext(ctx, "tcp", localhostAddr)
			if err != nil {
				b.Fatalf("DialContext error: %v", err)
			}
			_ = conn.Close()
		}
	})

	b.Run("Hostname_Parallel_WithDNSCache", func(b *testing.B) {
		der, err := NewDialer(
			WithEnableDNSCache(),
			WithDNSRefreshDuration("1m"),
			WithDNSCacheExpiration("1m"),
			WithDialerTimeout("5s"),
		)
		if err != nil {
			b.Fatalf("NewDialer error: %v", err)
		}
		ctx := context.Background()

		// Warm up
		warmConn, err := der.DialContext(ctx, "tcp", localhostAddr)
		if err != nil {
			b.Fatalf("warmup dial error: %v", err)
		}
		_ = warmConn.Close()

		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				conn, err := der.DialContext(ctx, "tcp", localhostAddr)
				if err != nil {
					b.Errorf("DialContext error: %v", err)
					return
				}
				_ = conn.Close()
			}
		})
	})
}

// BenchmarkDialer_LookupCache measures lookup and round-robin selection latency
// directly on the dialerCache structure.
func BenchmarkDialer_LookupCache(b *testing.B) {
	c, err := cache.New[*dialerCache]()
	if err != nil {
		b.Fatalf("cache.New error: %v", err)
	}
	c.Set("localhost", &dialerCache{
		ips: []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4"},
	})
	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(c),
	)
	if err != nil {
		b.Fatalf("NewDialer error: %v", err)
	}
	d, ok := der.(*dialer)
	if !ok {
		b.Fatalf("dialer type assertion failed: %T", der)
	}
	ctx := context.Background()

	b.Run("Lookup_CacheHit", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dc, err := d.lookup(ctx, "localhost")
			if err != nil || dc == nil {
				b.Fatalf("lookup error: %v", err)
			}
		}
	})

	b.Run("RoundRobin_IP_Sequential", func(b *testing.B) {
		dc := &dialerCache{
			ips: []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = dc.IP()
		}
	})

	b.Run("RoundRobin_IP_Parallel", func(b *testing.B) {
		dc := &dialerCache{
			ips: []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
		}
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = dc.IP()
			}
		})
	})
}

// BenchmarkDialer_TLSHandshake measures TLS dial performance and validates the
// avoidance of redundant tls.Config clones.
func BenchmarkDialer_TLSHandshake(b *testing.B) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	addr := srv.Listener.Addr().String()
	pool := x509.NewCertPool()
	pool.AddCert(srv.Certificate())

	b.Run("RepeatedTLS_CachedServerName", func(b *testing.B) {
		cfg := &tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
		}
		der, err := NewDialer(
			WithTLS(cfg),
			WithDialerTimeout("10s"),
		)
		if err != nil {
			b.Fatalf("NewDialer error: %v", err)
		}
		ctx := context.Background()

		// Warm up first dial so that clone is cached in tlsConfigs
		warmConn, err := der.DialContext(ctx, "tcp", addr)
		if err != nil {
			b.Fatalf("warmup dial error: %v", err)
		}
		_ = warmConn.Close()

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			conn, err := der.DialContext(ctx, "tcp", addr)
			if err != nil {
				b.Fatalf("DialContext error: %v", err)
			}
			_ = conn.Close()
		}
	})

	b.Run("PreSetServerName_DirectReuse", func(b *testing.B) {
		host, _, _ := SplitHostPort(addr)
		cfg := &tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		}
		der, err := NewDialer(
			WithTLS(cfg),
			WithDialerTimeout("10s"),
		)
		if err != nil {
			b.Fatalf("NewDialer error: %v", err)
		}
		ctx := context.Background()

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			conn, err := der.DialContext(ctx, "tcp", addr)
			if err != nil {
				b.Fatalf("DialContext error: %v", err)
			}
			_ = conn.Close()
		}
	})
}
