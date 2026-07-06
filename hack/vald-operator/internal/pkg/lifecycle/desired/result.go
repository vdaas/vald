package desired

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Result is the outcome of an IsReady check, representing the state to be reflected
// in a metav1.Condition. The Type and Timestamp are filled in by the lifecycle framework.
type Result struct {
	Status  metav1.ConditionStatus
	Reason  string
	Message string
}

func Progressing(msg string) Result {
	return Result{Status: metav1.ConditionUnknown, Reason: "Progressing", Message: msg}
}

func Pending(msg string) Result {
	return Result{Status: metav1.ConditionUnknown, Reason: "Pending", Message: msg}
}

func Succeeded() Result {
	return Result{Status: metav1.ConditionTrue, Reason: "Succeeded"}
}

func Failed(err error) Result {
	msg := "failed"
	if err != nil {
		msg = "failed: " + err.Error()
	}
	return Result{Status: metav1.ConditionFalse, Reason: "Failed", Message: msg}
}
