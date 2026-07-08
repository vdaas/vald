package valdoperatorrelease

import (
	"fmt"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
)

func (d *Domain) ConditionWaitForClusterCreate() *lifecycle.LifeCycle {
	cr := d.ValdOperatorRelease
	lc := &lifecycle.LifeCycle{}
	lc.Condition = lifecycle.Condition{
		Type:    lifecycle.ConditionWaitForClusterCreate,
		Message: "Waiting for Cluster Creation.",
	}
	lc.Checker = &desired.Prop{
		Check: func() desired.Result {
			if len(cr.Spec.Infrastructure) == 0 {
				return desired.Failed(fmt.Errorf("infrastructure configuration is missing"))
			}
			for _, infra := range cr.Spec.Infrastructure {
				if len(infra.Clusters) == 0 {
					return desired.Failed(fmt.Errorf("no clusters defined in configuration"))
				}
				for _, cluster := range infra.Clusters {
					if cluster.ID == "" || cluster.Name == "" {
						return desired.Pending(fmt.Sprintf("waiting for cluster to be provisioned: %#v", cluster))
					}
				}
			}
			return desired.Succeeded()
		},
	}

	return lc
}
