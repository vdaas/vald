package mvaldrelease

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
)

func (d *Domain) ConditionCompleted() *lifecycle.LifeCycle {
	lc := &lifecycle.LifeCycle{}
	lc.Condition = lifecycle.Condition{
		Type:    lifecycle.ConditionCompleted,
		Message: "VRS creation completed successfully.",
	}
	// Terminal phase: no Builder and no Checker. Both default to nil.
	return lc
}
