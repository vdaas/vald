package valdoperatorrelease

import (
	builder "github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/builder/vald"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
)

func (d *Domain) ConditionWaitForCreateVrs(k8sClient client.Client, cfg *config.Config, capability builder.NodePoolCapability) *lifecycle.LifeCycle {
	cr := d.ValdOperatorRelease
	lc := &lifecycle.LifeCycle{}
	lc.Condition = lifecycle.Condition{
		Type:    lifecycle.ConditionWaitForCreateVrs,
		Message: "Waiting for VRS creation.",
	}
	// Resource implements both Builder and ReadinessChecker; share the same
	// instance so the list produced by Build can be re-read by IsReady.
	// Domain itself supplies the rule values the Builder needs at Build time.
	res := &desired.Resource{
		List:    &unstructured.UnstructuredList{},
		Client:  k8sClient,
		Builder: builder.NewVrsBuilder(cr, cfg, capability, d),
	}
	lc.Builder = res
	lc.Checker = res

	return lc
}
