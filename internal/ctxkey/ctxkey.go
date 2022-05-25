package ctxkey

import "context"

type contextKey string

const grpcMethodContextKey contextKey = "grpc_method"

// WithGRPCMethod returns a copy of parent in which the method associated with key (grpcMethodContextKey).
func WithGRPCMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, grpcMethodContextKey, method)
}

// FromGRPCMethod returns the value associated with this context for key (grpcMethodContextKey)
func FromGRPCMethod(ctx context.Context) string {
	return ctx.Value(grpcMethodContextKey).(string)
}

const backoffNameContextKey contextKey = "backoff_name"

// WithBackoffName returns a copy of parent in which the method associated with key (backoffNameContextKey).
func WithBackoffName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, backoffNameContextKey, name)
}

// FromBackoffName returns the value associated with this context for key (backoffNameContextKey)
func FromBackoffName(ctx context.Context) string {
	return ctx.Value(backoffNameContextKey).(string)
}
