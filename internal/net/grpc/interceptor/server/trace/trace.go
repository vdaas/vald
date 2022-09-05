package trace

import (
	"github.com/vdaas/vald/internal/net/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

func TraceInterceptor() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

func TraceStreamInterceptor() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor()
}
