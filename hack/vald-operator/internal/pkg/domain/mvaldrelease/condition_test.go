package mvaldrelease

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newDomainWithInfra(infras []v1.MvaldreleaseInfra) *Domain {
	return &Domain{
		Mvaldrelease: &v1.Mvaldrelease{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
			Spec: v1.MvaldreleaseSpec{
				Infrastructure: infras,
			},
		},
	}
}

func TestConditionWaitForClusterCreate(t *testing.T) {
	tests := []struct {
		name       string
		domain     *Domain
		wantStatus metav1.ConditionStatus
		wantReason string
	}{
		{
			name:       "empty infrastructure",
			domain:     newDomainWithInfra(nil),
			wantStatus: metav1.ConditionFalse,
			wantReason: "Failed",
		},
		{
			name: "infra with no clusters",
			domain: newDomainWithInfra([]v1.MvaldreleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{}},
			}),
			wantStatus: metav1.ConditionFalse,
			wantReason: "Failed",
		},
		{
			name: "cluster with empty Id",
			domain: newDomainWithInfra([]v1.MvaldreleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{{ID: "", Name: "cluster-a"}}},
			}),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name: "cluster with empty Name",
			domain: newDomainWithInfra([]v1.MvaldreleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{{ID: "abc-123", Name: ""}}},
			}),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name: "valid clusters",
			domain: newDomainWithInfra([]v1.MvaldreleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{{ID: "abc-123", Name: "cluster-a"}}},
			}),
			wantStatus: metav1.ConditionTrue,
			wantReason: "Succeeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := tt.domain.ConditionWaitForClusterCreate()
			assert.NotNil(t, lc)
			assert.NotNil(t, lc.Checker)

			result := lc.Checker.IsReady(context.Background())
			assert.Equal(t, tt.wantStatus, result.Status)
			assert.Equal(t, tt.wantReason, result.Reason)
		})
	}
}

func TestConditionWaitForClusterCreate_PendingIsNotFailed(t *testing.T) {
	// cluster.ID == "" means the external system is still creating the cluster. Should be Pending, not Failed.
	domain := newDomainWithInfra([]v1.MvaldreleaseInfra{
		{Role: "green", Clusters: []v1.DestClusters{{ID: "", Name: "cluster-a"}}},
	})
	lc := domain.ConditionWaitForClusterCreate()
	result := lc.Checker.IsReady(context.Background())
	assert.Equal(t, desired.Pending("").Status, result.Status)
	assert.NotEqual(t, desired.Failed(nil).Status, result.Status)
}
