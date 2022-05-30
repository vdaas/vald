package grpc

import "context"

type contextKey string

const grpcMethodContextKey contextKey = "grpc_method"

// WithGRPCMethod returns a copy of parent in which the method associated with key (grpcMethodContextKey).
func WithGRPCMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, grpcMethodContextKey, method)
}

// FromGRPCMethod returns the value associated with this context for key (grpcMethodContextKey).
func FromGRPCMethod(ctx context.Context) string {
	if v := ctx.Value(grpcMethodContextKey); v != nil {
		return v.(string)
	}
	return ""
}
