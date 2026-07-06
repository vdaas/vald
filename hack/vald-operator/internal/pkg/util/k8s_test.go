package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	condTypeReady  = "Ready"
	condTypeSynced = "Synced"
	condMsgAllGood = "all good"
)

func TestUpdateStatus(t *testing.T) {
	tests := []struct {
		name     string
		initial  []metav1.Condition
		newCond  metav1.Condition
		wantLen  int
		wantCond metav1.Condition
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

			UpdateStatus(&conditions, tt.newCond)

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

func TestToObjectSlice(t *testing.T) {
	t.Run("extracts items from UnstructuredList", func(t *testing.T) {
		obj := unstructured.Unstructured{}
		obj.SetName("foo")
		list := &unstructured.UnstructuredList{Items: []unstructured.Unstructured{obj}}

		objs, err := ToObjectSlice(list)
		require.NoError(t, err)
		assert.Len(t, objs, 1)
	})
	t.Run("returns error on empty list", func(t *testing.T) {
		list := &unstructured.UnstructuredList{}
		_, err := ToObjectSlice(list)
		assert.Error(t, err)
	})
}

func TestConvertToUnstructured(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "default"},
	}
	u, err := ConvertToUnstructured(pod)
	require.NoError(t, err)
	assert.Equal(t, "test-pod", u.GetName())
	assert.Equal(t, "default", u.GetNamespace())
}

func TestDeleteCondition(t *testing.T) {
	tests := []struct {
		name      string
		initial   []metav1.Condition
		condType  string
		wantLen   int
		wantTypes []string
	}{
		{
			name: "delete existing condition",
			initial: []metav1.Condition{
				{Type: condTypeReady, Status: metav1.ConditionTrue},
				{Type: condTypeSynced, Status: metav1.ConditionTrue},
			},
			condType:  condTypeReady,
			wantLen:   1,
			wantTypes: []string{condTypeSynced},
		},
		{
			name: "no-op when type not found",
			initial: []metav1.Condition{
				{Type: condTypeReady, Status: metav1.ConditionTrue},
			},
			condType:  "NotExist",
			wantLen:   1,
			wantTypes: []string{condTypeReady},
		},
		{
			name:      "no-op on empty slice",
			initial:   []metav1.Condition{},
			condType:  condTypeReady,
			wantLen:   0,
			wantTypes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := make([]metav1.Condition, len(tt.initial))
			copy(conditions, tt.initial)

			DeleteCondition(&conditions, tt.condType)

			assert.Len(t, conditions, tt.wantLen)
			types := make([]string, len(conditions))
			for i, c := range conditions {
				types[i] = c.Type
			}
			assert.Equal(t, tt.wantTypes, types)
		})
	}
}
