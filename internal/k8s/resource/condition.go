// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetStatusCondition updates LastTransitionTime only when Status changes,
// unlike UpsertCondition. Delegates to
// k8s.io/apimachinery/pkg/api/meta.SetStatusCondition.
func SetStatusCondition(conditions *[]metav1.Condition, newCondition metav1.Condition) {
	apimeta.SetStatusCondition(conditions, newCondition)
}

// UpsertCondition inserts newCond into conditions or replaces the entry with
// the same Type. Unlike meta.SetStatusCondition it replaces the entry (and
// its LastTransitionTime) whenever Status, Reason or Message differ, and
// leaves the slice untouched when all three are identical.
func UpsertCondition(conditions *[]metav1.Condition, newCond metav1.Condition) {
	for i, cond := range *conditions {
		if cond.Type == newCond.Type {
			if cond.Status == newCond.Status &&
				cond.Reason == newCond.Reason &&
				cond.Message == newCond.Message {
				return
			}
			(*conditions)[i] = newCond
			return
		}
	}
	*conditions = append(*conditions, newCond)
}
