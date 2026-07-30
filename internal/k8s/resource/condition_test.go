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
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	condTypeReady  = "Ready"
	condTypeSynced = "Synced"
	condMsgAllGood = "all good"
)

func TestUpsertCondition(t *testing.T) {
	tests := []struct {
		name     string
		newCond  metav1.Condition
		wantCond metav1.Condition
		initial  []metav1.Condition
		wantLen  int
	}{
		{
			name:     "add to empty slice",
			initial:  []metav1.Condition{},
			newCond:  metav1.Condition{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK"},
			wantLen:  1,
			wantCond: metav1.Condition{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK"},
		},
		{
			name: "add new type to non-empty slice",
			initial: []metav1.Condition{
				{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK"},
			},
			newCond:  metav1.Condition{Type: condTypeSynced, Status: metav1.ConditionTrue, Reason: condTypeSynced},
			wantLen:  2,
			wantCond: metav1.Condition{Type: condTypeSynced, Status: metav1.ConditionTrue, Reason: condTypeSynced},
		},
		{
			name: "update existing condition with different status",
			initial: []metav1.Condition{
				{Type: condTypeReady, Status: metav1.ConditionFalse, Reason: "NotReady", Message: "initializing"},
			},
			newCond:  metav1.Condition{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK"},
			wantLen:  1,
			wantCond: metav1.Condition{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK"},
		},
		{
			name: "no update when condition is identical",
			initial: []metav1.Condition{
				{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK", Message: condMsgAllGood},
			},
			newCond:  metav1.Condition{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK", Message: condMsgAllGood},
			wantLen:  1,
			wantCond: metav1.Condition{Type: condTypeReady, Status: metav1.ConditionTrue, Reason: "OK", Message: condMsgAllGood},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := make([]metav1.Condition, len(tt.initial))
			copy(conditions, tt.initial)

			UpsertCondition(&conditions, tt.newCond)

			assert.Len(t, conditions, tt.wantLen)

			var found *metav1.Condition
			for i := range conditions {
				if conditions[i].Type == tt.newCond.Type {
					found = &conditions[i]
					break
				}
			}
			assert.NotNil(t, found)
			assert.Equal(t, tt.wantCond.Status, found.Status)
			assert.Equal(t, tt.wantCond.Reason, found.Reason)
			assert.Equal(t, tt.wantCond.Message, found.Message)
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestSetStatusCondition(t *testing.T) {
// 	type args struct {
// 		conditions   *[]metav1.Condition
// 		newCondition metav1.Condition
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           conditions:nil,
// 		           newCondition:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           conditions:nil,
// 		           newCondition:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			SetStatusCondition(test.args.conditions, test.args.newCondition)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
