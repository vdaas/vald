package backoff

import "context"

type contextKey string

const backoffNameContextKey contextKey = "backoff_name"

// WithBackoffName returns a copy of parent in which the method associated with key (backoffNameContextKey).
func WithBackoffName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, backoffNameContextKey, name)
}

// FromBackoffName returns the value associated with this context for key (backoffNameContextKey).
func FromBackoffName(ctx context.Context) string {
	if val := ctx.Value(backoffNameContextKey); val != nil {
		return val.(string)
	}
	return ""
}
