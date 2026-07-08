//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package service

import (
	"context"

	"github.com/vdaas/vald/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Phase names double as the condition types recorded on the ValdOperatorRelease
// status.
const (
	phaseWaitForClusterCreate = "WaitingClusterCreate"
	phaseWaitForCreateVrs     = "WaitingCreateVrs"
	phaseCompleted            = "Completed"
)

// result is the outcome of a phase check, representing the state to be
// reflected in a metav1.Condition. The Type and Timestamp are filled in by
// phase.condition.
type result struct {
	status  metav1.ConditionStatus
	reason  string
	message string
}

func progressing(msg string) result {
	return result{status: metav1.ConditionUnknown, reason: "Progressing", message: msg} //nolint:goconst
}

func pending(msg string) result {
	return result{status: metav1.ConditionUnknown, reason: "Pending", message: msg} //nolint:goconst
}

func succeeded() result {
	return result{status: metav1.ConditionTrue, reason: "Succeeded"} //nolint:goconst
}

func failed(err error) result {
	msg := "failed"
	if err != nil {
		msg = "failed: " + err.Error()
	}
	return result{status: metav1.ConditionFalse, reason: "Failed", message: msg} //nolint:goconst
}

// phase is one step of the reconcile pipeline. Either build or check (or
// both) may be set:
//   - build != nil : the phase produces Kubernetes objects to apply; fetch
//     then lists the existing objects eligible for pruning
//   - check != nil : the phase observes a property and reports readiness
//   - both nil     : a terminal/no-op phase (e.g. Completed)
type phase struct {
	build   func(ctx context.Context) ([]k8s.Object, error)
	fetch   func(ctx context.Context) ([]k8s.Object, error)
	check   func(ctx context.Context) result
	name    string
	message string
}

// terminal reports whether the phase has no work: nothing to build, nothing
// to check.
func (p *phase) terminal() bool {
	return p.build == nil && p.check == nil
}

// condition renders the phase result as a metav1.Condition, falling back to
// the phase's default message when the result carries none.
func (p *phase) condition(r result) metav1.Condition {
	message := r.message
	if message == "" {
		message = p.message
	}
	return metav1.Condition{
		Type:               p.name,
		Status:             r.status,
		LastTransitionTime: metav1.Now(),
		Reason:             r.reason,
		Message:            message,
	}
}

type phases []phase

// index returns the position of the phase whose name matches condType. An
// empty condType defaults to the first phase; an unknown one returns -1.
func (ps phases) index(condType string) int {
	if condType == "" {
		return 0
	}
	for i := range ps {
		if ps[i].name == condType {
			return i
		}
	}
	return -1
}
