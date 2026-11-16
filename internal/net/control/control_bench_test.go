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

package control_test

import (
	"net"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/net/control"
)

// BenchmarkSocketController_Control measures the execution latency and allocations
// of the socket controller function during RawConn Control callbacks.
func BenchmarkSocketController_Control(b *testing.B) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("net.Listen error: %v", err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		conn, _ := ln.Accept()
		if conn != nil {
			<-done
			_ = conn.Close()
		}
	}()
	defer close(done)

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		b.Fatalf("net.Dial error: %v", err)
	}
	defer conn.Close()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		b.Fatal("conn is not *net.TCPConn")
	}
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		b.Fatalf("SyscallConn error: %v", err)
	}

	b.Run("Default_TCPNoDelayOnly", func(b *testing.B) {
		ctrl := control.New(control.TCPNoDelay, 0)
		fn := ctrl.GetControl()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fn("tcp", "127.0.0.1:0", rawConn)
		}
	})

	b.Run("ZeroFlags_Default", func(b *testing.B) {
		ctrl := control.New(0, 0)
		fn := ctrl.GetControl()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fn("tcp", "127.0.0.1:0", rawConn)
		}
	})

	b.Run("MultipleFlags_ActiveOnly", func(b *testing.B) {
		ctrl := control.New(
			control.ReusePort|
				control.ReuseAddr|
				control.TCPNoDelay|
				control.TCPQuickAck,
			0,
		)
		fn := ctrl.GetControl()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fn("tcp", "127.0.0.1:0", rawConn)
		}
	})

	b.Run("WithKeepAliveSeconds", func(b *testing.B) {
		ctrl := control.New(
			control.TCPNoDelay,
			control.KeepAliveSeconds(30*time.Second),
		)
		fn := ctrl.GetControl()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fn("tcp", "127.0.0.1:0", rawConn)
		}
	})
}

// BenchmarkSocketController_New measures initialization latency and allocations
// of creating socket controllers.
func BenchmarkSocketController_New(b *testing.B) {
	b.Run("Default", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = control.New(control.TCPNoDelay, 0)
		}
	})

	b.Run("ZeroFlags", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = control.New(0, 0)
		}
	})

	b.Run("MultipleFlags", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = control.New(
				control.ReusePort|
					control.ReuseAddr|
					control.TCPNoDelay|
					control.TCPFastOpen|
					control.TCPQuickAck,
				control.KeepAliveSeconds(30*time.Second),
			)
		}
	})
}

// BenchmarkSocketController_KeepAliveSeconds measures KeepAliveSeconds conversion performance.
func BenchmarkSocketController_KeepAliveSeconds(b *testing.B) {
	d := 30 * time.Second
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = control.KeepAliveSeconds(d)
	}
}
