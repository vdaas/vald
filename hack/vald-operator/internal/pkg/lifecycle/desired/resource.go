package desired

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/util"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Resource implements both Builder and ReadinessChecker. Build delegates to an
// inner SubResourceBuilder, then IsReady fetches each produced object and
// reports readiness via the optional Check func.
type Resource struct {
	List    client.ObjectList
	Client  client.Client
	Builder SubResourceBuilder
	Check   func(obj client.Object) Result
}

type SubResourceBuilder interface {
	Build(ctx context.Context) (client.ObjectList, error)
}

func (drs *Resource) IsReady(ctx context.Context) Result {
	items, err := util.ToObjectSlice(drs.List)
	if err != nil {
		return Failed(err)
	}
	for _, obj := range items {
		if err := drs.Client.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
			if errors.IsNotFound(err) {
				return Progressing("")
			}
			return Failed(err)
		}
		if drs.Check != nil {
			result := drs.Check(obj)
			if result.Status != metav1.ConditionTrue {
				return result
			}
		}
	}
	return Succeeded()
}

func (drs *Resource) Build(ctx context.Context) (client.ObjectList, error) {
	var err error
	drs.List, err = drs.Builder.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build resource: %w", err)
	}
	if drs.List == nil {
		return nil, fmt.Errorf("built resource is nil")
	}
	return drs.List, nil
}
