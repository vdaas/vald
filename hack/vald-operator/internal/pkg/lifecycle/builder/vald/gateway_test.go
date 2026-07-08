package vald

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	corev1 "k8s.io/api/core/v1"
	apires "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func minimalVor() *v1.ValdOperatorRelease {
	return &v1.ValdOperatorRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test-ns",
		},
		Spec: v1.ValdOperatorReleaseSpec{
			Infrastructure: []v1.ValdOperatorReleaseInfra{
				{
					Role:   "green",
					Active: true,
					Clusters: []v1.DestClusters{
						{Name: "cluster-1"},
					},
					NodePools: v1.NodePools{
						v1.NodePoolTypeGeneral: v1.NodePool{
							Replicas:        1,
							MachineResource: v1.MachineResource{Cpu: "2", Memory: "4Gi"},
						},
					},
				},
			},
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{
						Ngt: v1.Ngt{Dimension: 128},
					},
				},
			},
		},
	}
}

func initConfig(t *testing.T) *config.Config {
	t.Helper()
	config.DefaultVrsPath = "testdata/vrs.yaml"
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("config load: %v", err)
	}
	return cfg
}

func newTestBuilder(cr *v1.ValdOperatorRelease, cfg *config.Config) *VrsBuilder {
	return NewVrsBuilder(cr, cfg, AlwaysAvailable(), stubRules{})
}

// stubRules mirrors the production Domain rules so builder tests stay
// independent of the domain/valdoperatorrelease package (which would otherwise
// introduce an import cycle). Keep this in sync with Domain.
type stubRules struct{}

func (stubRules) ResolveAgentNodePool(infra v1.ValdOperatorReleaseInfra) AgentNodePoolSpec {
	gn := infra.NodePools.GetNodePool(v1.NodePoolTypeGeneral)
	an := infra.NodePools.GetNodePool(v1.NodePoolTypeValdAgent)
	if an == nil || an.Replicas == 0 {
		return AgentNodePoolSpec{
			NodeCount:       gn.Replicas,
			MachineResource: gn.MachineResource.GetResourceList(),
		}
	}
	return AgentNodePoolSpec{
		NodeCount:       an.Replicas,
		MachineResource: an.MachineResource.GetResourceList(),
	}
}

func (stubRules) AgentPvSize(memoryBytes int64, pvBufferRatio float64, pvMinSizeBytes int64) string {
	const gi = int64(1) << 30
	size := max(int64(float64(memoryBytes)*pvBufferRatio), pvMinSizeBytes)
	sizeGi := (size + gi - 1) / gi
	return apires.NewQuantity(sizeGi*gi, apires.BinarySI).String()
}

func TestBuildLb_IngressEnabled(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.Ingress = &v1.GatewayIngress{
		Enabled: true,
		Host:    "foo.example.com",
	}
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.True(t, lb.Ingress.Enabled)
	assert.Equal(t, "foo.example.com", lb.Ingress.Host)
}

func TestBuildLb_IngressDisabledWhenNil(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.False(t, lb.Ingress.Enabled)
	assert.Empty(t, lb.Ingress.Host)
}

func TestBuildLb_ServiceTypeDefault(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, corev1.ServiceTypeNodePort, lb.ServiceType)
}

func TestBuildLb_ServiceTypeLoadBalancer(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.ServiceType = "LoadBalancer"
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, corev1.ServiceTypeLoadBalancer, lb.ServiceType)
}

func TestBuildLb_ServiceTypeClusterIP(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.ServiceType = "ClusterIP"
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, corev1.ServiceTypeClusterIP, lb.ServiceType)
}
