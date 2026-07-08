package valdoperatorrelease

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/util"
)

// AdvanceToNextPhase moves Status.Phase / Progress.Completed to the next
// phase in the flow if one exists. The seeded condition is Succeeded for
// terminal (Builder==nil && Checker==nil) phases and Progressing otherwise.
//
// Returns the index of the new current phase, or -1 if no next phase
// existed.
func (d *Domain) AdvanceToNextPhase(flow *lifecycle.Flow) int {
	next := flow.GetNext()
	if next < 0 {
		return -1
	}
	nx := &flow.LifeCycles[next]
	d.Status.Phase = nx.Condition.Type
	d.Status.Progress.Completed = next

	var seed desired.Result
	if nx.Builder == nil && nx.Checker == nil {
		seed = desired.Succeeded()
	} else {
		seed = desired.Progressing("")
	}
	util.UpdateStatus(&d.Status.Conditions, nx.Condition.MakeCondition(seed))
	return next
}
