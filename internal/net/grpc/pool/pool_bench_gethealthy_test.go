package pool

import (
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Benchmark_GetHealthyConn(b *testing.B) {
	addr := "localhost:5003"
	defer ListenAndServe(b, addr)()

	ctx := b.Context()
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(10),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		b.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond) // Wait for connections to be ready

	po, ok := p.(*pool)
	if !ok {
		b.Fatal("not a *pool")
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	for b.Loop() {
		po.getHealthyConn(ctx)
	}
}

func BenchmarkParallel_GetHealthyConn(b *testing.B) {
	addr := "localhost:5004"
	defer ListenAndServe(b, addr)()

	ctx := b.Context()
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(10),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		b.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond) // Wait for connections to be ready

	po, ok := p.(*pool)
	if !ok {
		b.Fatal("not a *pool")
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			po.getHealthyConn(ctx)
		}
	})
}
