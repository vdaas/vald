package mvaldrelease

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makePhaseFlow() *lifecycle.Flow {
	lcs := lifecycle.LifeCycles{
		{Condition: lifecycle.Condition{Type: "phase-a"}, Checker: &desired.Prop{}},
		{Condition: lifecycle.Condition{Type: "phase-b"}, Checker: &desired.Prop{}},
		{Condition: lifecycle.Condition{Type: "phase-terminal"}}, // no Builder, no Checker
	}
	return lifecycle.NewFlow(lcs, 0)
}

func newDomain() *Domain {
	return &Domain{
		Mvaldrelease: &v1.Mvaldrelease{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		},
	}
}

func TestDomain_AdvanceToNextPhase_ToNonTerminal(t *testing.T) {
	d := newDomain()
	flow := makePhaseFlow()

	next := d.AdvanceToNextPhase(flow)

	assert.Equal(t, 1, next)
	assert.Equal(t, "phase-b", d.Status.Phase)
	assert.Equal(t, 1, d.Status.Progress.Completed)
	assert.Len(t, d.Status.Conditions, 1)
	assert.Equal(t, "phase-b", d.Status.Conditions[0].Type)
	assert.Equal(t, metav1.ConditionUnknown, d.Status.Conditions[0].Status, "non-terminal phases seed Progressing")
}

func TestDomain_AdvanceToNextPhase_ToTerminal(t *testing.T) {
	d := newDomain()
	flow := makePhaseFlow()
	flow.Current = 1 // currently at phase-b; next is phase-terminal

	next := d.AdvanceToNextPhase(flow)

	assert.Equal(t, 2, next)
	assert.Equal(t, "phase-terminal", d.Status.Phase)
	assert.Equal(t, 2, d.Status.Progress.Completed)
	assert.Equal(t, metav1.ConditionTrue, d.Status.Conditions[0].Status, "terminal phases seed Succeeded")
}

func TestDomain_AdvanceToNextPhase_NoNext(t *testing.T) {
	d := newDomain()
	flow := makePhaseFlow()
	flow.Current = 2 // already at the last phase

	next := d.AdvanceToNextPhase(flow)

	assert.Equal(t, -1, next)
	assert.Empty(t, d.Status.Phase)
	assert.Empty(t, d.Status.Conditions)
}
