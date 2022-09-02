package trace

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

var (
	UnaryClientInterceptor  = otelgrpc.UnaryClientInterceptor
	StreamClientInterceptor = otelgrpc.StreamClientInterceptor
)
