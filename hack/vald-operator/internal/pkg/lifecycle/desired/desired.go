package desired

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Builder produces the desired Kubernetes objects for a single lifecycle phase.
// Implementations are expected to be pure functions of their captured inputs;
// any I/O needed to discover external state should be resolved by the caller
// and passed in at construction time.
type Builder interface {
	Build(ctx context.Context) (client.ObjectList, error)
}

// ReadinessChecker reports whether the phase has reached its desired state.
// It is independent of Builder — a phase may have only a check (no resources
// to create) or both.
type ReadinessChecker interface {
	IsReady(ctx context.Context) Result
}

type OperationResults map[string]controllerutil.OperationResult
