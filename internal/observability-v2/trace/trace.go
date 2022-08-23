package trace

import "context"

type Tracer interface {
	Start(ctx context.Context) error
}
