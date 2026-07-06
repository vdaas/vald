package desired_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
)

// --- Result constructors ---

func TestProgressing(t *testing.T) {
	r := desired.Progressing("building")
	assert.Equal(t, metav1.ConditionUnknown, r.Status)
	assert.Equal(t, "Progressing", r.Reason)
	assert.Equal(t, "building", r.Message)
}

func TestPending(t *testing.T) {
	r := desired.Pending("waiting for external system")
	assert.Equal(t, metav1.ConditionUnknown, r.Status)
	assert.Equal(t, "Pending", r.Reason)
	assert.Equal(t, "waiting for external system", r.Message)
}

func TestSucceeded(t *testing.T) {
	r := desired.Succeeded()
	assert.Equal(t, metav1.ConditionTrue, r.Status)
	assert.Equal(t, "Succeeded", r.Reason)
	assert.Empty(t, r.Message)
}

func TestFailed(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		r := desired.Failed(fmt.Errorf("something broke"))
		assert.Equal(t, metav1.ConditionFalse, r.Status)
		assert.Equal(t, "Failed", r.Reason)
		assert.Contains(t, r.Message, "something broke")
	})
	t.Run("nil error", func(t *testing.T) {
		r := desired.Failed(nil)
		assert.Equal(t, metav1.ConditionFalse, r.Status)
		assert.Equal(t, "Failed", r.Reason)
		assert.Equal(t, "failed", r.Message)
	})
}

// --- Prop ---

func TestProp_IsReady(t *testing.T) {
	want := desired.Succeeded()
	p := &desired.Prop{
		Check: func() desired.Result { return want },
	}
	got := p.IsReady(context.Background())
	assert.Equal(t, want, got)
}

// --- Resource ---

func newConfigMap(name, ns string) *v1.ConfigMap {
	return &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns}}
}

func TestResource_IsReady_NotFound(t *testing.T) {
	fc := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
	cm := newConfigMap("test", "default")
	list := &v1.ConfigMapList{Items: []v1.ConfigMap{*cm}}

	r := &desired.Resource{List: list, Client: fc}
	result := r.IsReady(context.Background())
	assert.Equal(t, metav1.ConditionUnknown, result.Status)
	assert.Equal(t, "Progressing", result.Reason)
}

func TestResource_IsReady_AllFound(t *testing.T) {
	cm := newConfigMap("test", "default")
	fc := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(cm).Build()
	list := &v1.ConfigMapList{Items: []v1.ConfigMap{*cm}}

	r := &desired.Resource{List: list, Client: fc}
	result := r.IsReady(context.Background())
	assert.Equal(t, metav1.ConditionTrue, result.Status)
	assert.Equal(t, "Succeeded", result.Reason)
}

func TestResource_IsReady_CheckFails(t *testing.T) {
	cm := newConfigMap("test", "default")
	fc := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(cm).Build()
	list := &v1.ConfigMapList{Items: []v1.ConfigMap{*cm}}

	customResult := desired.Pending("waiting")
	r := &desired.Resource{
		List:   list,
		Client: fc,
		Check:  func(_ client.Object) desired.Result { return customResult },
	}
	result := r.IsReady(context.Background())
	assert.Equal(t, customResult, result)
}

func TestResource_Build(t *testing.T) {
	want := &v1.ConfigMapList{Items: []v1.ConfigMap{{}}}

	builder := &fakeBuilder{list: want}
	r := &desired.Resource{Builder: builder}

	got, err := r.Build(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestResource_Build_NilList(t *testing.T) {
	builder := &fakeBuilder{list: nil}
	r := &desired.Resource{Builder: builder}

	_, err := r.Build(context.Background())
	assert.Error(t, err)
}

type fakeBuilder struct {
	list client.ObjectList
	err  error
}

func (f *fakeBuilder) Build(_ context.Context) (client.ObjectList, error) {
	return f.list, f.err
}
