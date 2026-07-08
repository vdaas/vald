package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
)

// stubBuilder returns a fixed object list (or error) and is enough to drive
// ResourceSyncer.Sync without standing up an envtest.
type stubBuilder struct {
	list client.ObjectList
	err  error
}

func (b *stubBuilder) Build(_ context.Context) (client.ObjectList, error) {
	return b.list, b.err
}

func newSchemeForSyncerTests(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))
	assert.NoError(t, controllerv1.AddToScheme(scheme))
	return scheme
}

func newOwner() *controllerv1.ValdOperatorRelease {
	return &controllerv1.ValdOperatorRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "owner",
			Namespace:  "default",
			UID:        types.UID("owner-uid"),
			Generation: 7,
		},
	}
}

func TestResourceSyncer_MakeKey_IncludesGVK(t *testing.T) {
	s := &ResourceSyncer{}
	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "n", Namespace: "ns"}}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	got := s.makeKey(cm)
	assert.Equal(t, "/v1/ConfigMap/ns/n", got)
}

func TestResourceSyncer_Sync_NilList(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewResourceSyncer(c, scheme)

	res, err := s.Sync(context.Background(), &stubBuilder{list: nil}, newOwner())
	assert.NoError(t, err)
	assert.Contains(t, res, "no_resources")
}

func TestResourceSyncer_Sync_BuilderError(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewResourceSyncer(c, scheme)

	_, err := s.Sync(context.Background(), &stubBuilder{err: fmt.Errorf("boom")}, newOwner())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestResourceSyncer_Sync_CreatesObjectsWithOwnerAndLabels(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewResourceSyncer(c, scheme)
	owner := newOwner()

	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "alpha", Namespace: "default"}}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	list := &corev1.ConfigMapList{Items: []corev1.ConfigMap{*cm}}
	list.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMapList"))

	res, err := s.Sync(context.Background(), &stubBuilder{list: list}, owner)
	assert.NoError(t, err)
	assert.Contains(t, res, "/v1/ConfigMap/default/alpha")

	got := &corev1.ConfigMap{}
	assert.NoError(t, c.Get(context.Background(), client.ObjectKey{Name: "alpha", Namespace: "default"}, got))
	assert.Equal(t, "7", got.Labels["managed-generation"])
	assert.Len(t, got.OwnerReferences, 1)
	assert.Equal(t, types.UID("owner-uid"), got.OwnerReferences[0].UID)
}

func TestResourceSyncer_Sync_PrunesOrphans(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newOwner()

	// Pre-existing orphan that is owned by `owner` but not in the new Build.
	orphan := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orphan",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         controllerv1.GroupVersion.String(),
					Kind:               "ValdOperatorRelease",
					Name:               owner.Name,
					UID:                owner.UID,
					Controller:         boolPtr(true),
					BlockOwnerDeletion: boolPtr(true),
				},
			},
		},
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(orphan).Build()
	s := NewResourceSyncer(c, scheme)

	// New Build emits a different ConfigMap; pruning should remove orphan.
	keep := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "keep", Namespace: "default"}}
	keep.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	list := &corev1.ConfigMapList{Items: []corev1.ConfigMap{*keep}}
	list.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMapList"))

	res, err := s.Sync(context.Background(), &stubBuilder{list: list}, owner)
	assert.NoError(t, err)
	assert.Contains(t, res, "/v1/ConfigMap/default/keep")
	assert.Equal(t, desired.OperationResults{
		"/v1/ConfigMap/default/keep":   res["/v1/ConfigMap/default/keep"],
		"/v1/ConfigMap/default/orphan": "pruned",
	}, res)

	err = c.Get(context.Background(), client.ObjectKey{Name: "orphan", Namespace: "default"}, &corev1.ConfigMap{})
	assert.Error(t, err, "orphan should have been deleted")
}

func boolPtr(b bool) *bool { return &b }
