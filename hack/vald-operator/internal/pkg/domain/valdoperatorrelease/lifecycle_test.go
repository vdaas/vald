package valdoperatorrelease

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	builder "github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/builder/vald"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// testConfig returns an empty *config.Config sufficient for condition tests that
// don't actually invoke Build/IsReady.
func testConfig() *config.Config {
	return &config.Config{}
}

func newDomainWithPhase(phase string) *Domain {
	return &Domain{
		ValdOperatorRelease: &v1.ValdOperatorRelease{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
			Status: v1.ValdOperatorReleaseStatus{
				Phase: phase,
			},
		},
	}
}

func TestConditionCompleted(t *testing.T) {
	d := newDomainWithPhase("")
	lc := d.ConditionCompleted()

	assert.Equal(t, lifecycle.ConditionCompleted, lc.Condition.Type)
	assert.Nil(t, lc.Builder)
	assert.Nil(t, lc.Checker)
}

func TestConditionWaitForCreateVrs(t *testing.T) {
	d := newDomainWithPhase("")
	lc := d.ConditionWaitForCreateVrs(nil, testConfig(), builder.AlwaysAvailable())

	assert.Equal(t, lifecycle.ConditionWaitForCreateVrs, lc.Condition.Type)
	assert.NotNil(t, lc.Builder)
	assert.NotNil(t, lc.Checker)
}

func TestInitProgress(t *testing.T) {
	tests := []struct {
		name          string
		phase         string
		wantTotal     int
		wantCompleted int
	}{
		{
			name:          "phase at start",
			phase:         string(lifecycle.ConditionWaitForClusterCreate),
			wantTotal:     2,
			wantCompleted: 0,
		},
		{
			name:          "phase at second step",
			phase:         string(lifecycle.ConditionWaitForCreateVrs),
			wantTotal:     2,
			wantCompleted: 1,
		},
		{
			name:          "phase at completed",
			phase:         string(lifecycle.ConditionCompleted),
			wantTotal:     2,
			wantCompleted: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newDomainWithPhase(tt.phase)
			lcs := lifecycle.LifeCycles{
				*d.ConditionWaitForClusterCreate(),
				*d.ConditionWaitForCreateVrs(nil, testConfig(), builder.AlwaysAvailable()),
				*d.ConditionCompleted(),
			}
			d.InitProgress(lcs)
			assert.Equal(t, tt.wantTotal, d.Status.Progress.Total)
			assert.Equal(t, tt.wantCompleted, d.Status.Progress.Completed)
		})
	}
}

func TestNewLifeCycleFlow(t *testing.T) {
	d := newDomainWithPhase(string(lifecycle.ConditionWaitForClusterCreate))
	flow := d.NewLifeCycleFlow(nil, testConfig(), builder.AlwaysAvailable())

	assert.NotNil(t, flow)
	assert.Equal(t, 0, flow.Current)
}
