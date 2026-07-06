package lifecycle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMakeCondition(t *testing.T) {
	tests := []struct {
		name       string
		result     desired.Result
		wantStatus metav1.ConditionStatus
		wantReason string
	}{
		{
			name:       "Succeeded",
			result:     desired.Succeeded(),
			wantStatus: metav1.ConditionTrue,
			wantReason: "Succeeded",
		},
		{
			name:       "Progressing",
			result:     desired.Progressing(""),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Progressing",
		},
		{
			name:       "Pending",
			result:     desired.Pending("waiting for VPC"),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name:       "Failed",
			result:     desired.Failed(assert.AnError),
			wantStatus: metav1.ConditionFalse,
			wantReason: "Failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Condition{Type: ConditionWaitForClusterCreate, Message: "default message"}
			got := c.MakeCondition(tt.result)
			assert.Equal(t, ConditionWaitForClusterCreate, got.Type)
			assert.Equal(t, tt.wantStatus, got.Status)
			assert.Equal(t, tt.wantReason, got.Reason)
			assert.False(t, got.LastTransitionTime.IsZero())
		})
	}
}

func TestMakeCondition_MessageFallback(t *testing.T) {
	c := Condition{Type: ConditionWaitForClusterCreate, Message: "default message"}

	// Resultにメッセージがない場合はConditionのデフォルトメッセージを使う
	assert.Equal(t, "default message", c.MakeCondition(desired.Progressing("")).Message)
	assert.Equal(t, "default message", c.MakeCondition(desired.Succeeded()).Message)

	// Resultにメッセージがある場合はそちらを優先する
	assert.Equal(t, "waiting for VPC", c.MakeCondition(desired.Pending("waiting for VPC")).Message)
	assert.Contains(t, c.MakeCondition(desired.Failed(assert.AnError)).Message, assert.AnError.Error())
}

func TestMakeCondition_StatusesAreDistinct(t *testing.T) {
	c := Condition{Type: ConditionWaitForClusterCreate, Message: "msg"}
	assert.NotEqual(t, c.MakeCondition(desired.Progressing("")).Status, c.MakeCondition(desired.Succeeded()).Status)
	assert.NotEqual(t, c.MakeCondition(desired.Progressing("")).Status, c.MakeCondition(desired.Failed(assert.AnError)).Status)
	assert.NotEqual(t, c.MakeCondition(desired.Succeeded()).Status, c.MakeCondition(desired.Failed(assert.AnError)).Status)
}
