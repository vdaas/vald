package lifecycle

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Condition struct {
	Type    string
	Message string
}

const (
	ConditionWaitForClusterCreate string = "WaitingClusterCreate"
	ConditionWaitForCreateVrs     string = "WaitingCreateVrs"
	ConditionCompleted            string = "Completed"
)

func (c Condition) MakeCondition(r desired.Result) metav1.Condition {
	message := r.Message
	if message == "" {
		message = c.Message
	}
	return metav1.Condition{
		Type:               c.Type,
		Status:             r.Status,
		LastTransitionTime: metav1.Now(),
		Reason:             r.Reason,
		Message:            message,
	}
}
