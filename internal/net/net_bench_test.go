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
	"testing"
)

// BenchmarkNet_Parse measures the latency and allocations of Parse across
// IP literals, hostnames with ports, and hostnames without ports.
func BenchmarkNet_Parse(b *testing.B) {
	ctx := context.Background()

	b.Run("IPLiteral_IPv4", func(b *testing.B) {
		const addr = "192.168.1.100:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _, _, _, _ = Parse(ctx, addr)
		}
	})

	b.Run("IPLiteral_IPv6", func(b *testing.B) {
		const addr = "[2001:db8::1]:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _, _, _, _ = Parse(ctx, addr)
		}
	})

	b.Run("Hostname_WithPort", func(b *testing.B) {
		const addr = "localhost:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _, _, _, _ = Parse(ctx, addr)
		}
	})

	b.Run("Hostname_WithoutPort", func(b *testing.B) {
		// Triggers the missing port error branch in Parse
		const addr = "localhost"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _, _, _, _ = Parse(ctx, addr)
		}
	})
}

// BenchmarkNet_SplitHostPort measures address splitting throughput and allocations.
func BenchmarkNet_SplitHostPort(b *testing.B) {
	b.Run("IPv4", func(b *testing.B) {
		const addr = "192.168.1.1:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = SplitHostPort(addr)
		}
	})

	b.Run("IPv6", func(b *testing.B) {
		const addr = "[2001:db8::1]:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = SplitHostPort(addr)
		}
	})

	b.Run("Hostname", func(b *testing.B) {
		const addr = "vald-agent.default.svc.cluster.local:8081"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = SplitHostPort(addr)
		}
	})

	b.Run("MissingPort", func(b *testing.B) {
		const addr = "vald-agent.default.svc.cluster.local"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = SplitHostPort(addr)
		}
	})
}

// BenchmarkNet_JoinHostPort measures address joining throughput and allocations.
func BenchmarkNet_JoinHostPort(b *testing.B) {
	b.Run("IPv4", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = JoinHostPort("192.168.1.1", 8080)
		}
	})

	b.Run("IPv6", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = JoinHostPort("2001:db8::1", 8080)
		}
	})

	b.Run("Hostname", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = JoinHostPort("vald-agent.default.svc.cluster.local", 8081)
		}
	})
}

// BenchmarkNet_IsIPLiteral measures the throughput of classifying IP literals.
func BenchmarkNet_IsIPLiteral(b *testing.B) {
	b.Run("IPv4_Literal", func(b *testing.B) {
		const host = "192.168.1.1"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = IsIPLiteral(host)
		}
	})

	b.Run("IPv4_WithPort", func(b *testing.B) {
		const host = "192.168.1.1:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = IsIPLiteral(host)
		}
	})

	b.Run("IPv6_Literal", func(b *testing.B) {
		const host = "2001:db8::1"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = IsIPLiteral(host)
		}
	})

	b.Run("IPv6_BracketedWithPort", func(b *testing.B) {
		const host = "[2001:db8::1]:8080"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = IsIPLiteral(host)
		}
	})

	b.Run("Hostname", func(b *testing.B) {
		const host = "vald-agent.default.svc.cluster.local"
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = IsIPLiteral(host)
		}
	})
}
