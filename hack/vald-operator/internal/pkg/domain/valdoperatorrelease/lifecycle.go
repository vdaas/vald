package valdoperatorrelease

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	builder "github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/builder/vald"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *Domain) NewLifeCycleFlow(k8sClient client.Client, cfg *config.Config, capability builder.NodePoolCapability) *lifecycle.Flow {
	cr := d.ValdOperatorRelease
	lcs := lifecycle.LifeCycles{
		*d.ConditionWaitForClusterCreate(),
		*d.ConditionWaitForCreateVrs(k8sClient, cfg, capability),
		*d.ConditionCompleted(),
	}

	return lifecycle.NewFlow(lcs, lcs.GetIndex(cr.Status.Phase))
}
