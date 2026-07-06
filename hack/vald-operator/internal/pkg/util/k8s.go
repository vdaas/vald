package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func ConvertToUnstructured(c client.Object) (*unstructured.Unstructured, error) {
	o, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to unstructured: %w", err)
	}

	us := &unstructured.Unstructured{Object: o}
	return us, nil
}

func UpdateStatus(conditions *[]metav1.Condition, newCond metav1.Condition) {
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

func DeleteCondition(conditions *[]metav1.Condition, condType string) {
	for i, cond := range *conditions {
		if cond.Type == condType {
			*conditions = append((*conditions)[:i], (*conditions)[i+1:]...)
			return
		}
	}
}

func ToObjectSlice(list client.ObjectList) ([]client.Object, error) {
	raw, err := meta.ExtractList(list)
	if err != nil {
		return nil, fmt.Errorf("failed to convert list to object slice: %w", err)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("no objects found in the list")
	}
	var objs []client.Object
	for _, o := range raw {
		if obj, ok := o.(client.Object); ok {
			objs = append(objs, obj)
		}
	}
	return objs, nil
}
