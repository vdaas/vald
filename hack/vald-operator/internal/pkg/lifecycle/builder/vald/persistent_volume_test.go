package vald

import (
	"testing"

	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func newAgentReleaseWithMemory(memory string) *valdrelease.ValdRelease {
	r := &valdrelease.ValdRelease{}
	r.Spec.Agent.Resources = &corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(memory),
		},
	}
	return r
}

func TestVrsBuilder_ReflectPersistentVolume_Disabled(t *testing.T) {
	b := &VrsBuilder{CR: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: &v1.AgentPersistentVolume{Enabled: false}},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)
	assert.Nil(t, r.Spec.Agent.PersistentVolume, "PV must not be set when Enabled=false")
}

func TestVrsBuilder_ReflectPersistentVolume_Nil(t *testing.T) {
	b := &VrsBuilder{CR: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: nil},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)
	assert.Nil(t, r.Spec.Agent.PersistentVolume, "PV must not be set when CR PersistentVolume is nil")
}

func TestVrsBuilder_ReflectPersistentVolume_FromCR(t *testing.T) {
	cfg := &config.Config{
		DefaultStorageClass: "fallback-sc",
		DefaultAccessMode:   "fallback-am",
		PvBufferRatio:       1.5,
		PvMinSizeBytes:      int64(1) << 30,
	}

	b := &VrsBuilder{Config: cfg, Rules: stubRules{}, CR: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: &v1.AgentPersistentVolume{
						Enabled:      true,
						StorageClass: "from-cr-sc",
						AccessMode:   "from-cr-am",
					}},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)

	assert.NotNil(t, r.Spec.Agent.PersistentVolume)
	assert.True(t, r.Spec.Agent.PersistentVolume.Enabled)
	assert.Equal(t, "from-cr-sc", r.Spec.Agent.PersistentVolume.StorageClass, "CR value must win")
	assert.Equal(t, "from-cr-am", r.Spec.Agent.PersistentVolume.AccessMode, "CR value must win")
}

func TestVrsBuilder_ReflectPersistentVolume_EnvFallback(t *testing.T) {
	cfg := &config.Config{
		DefaultStorageClass: "fallback-sc",
		DefaultAccessMode:   "fallback-am",
		PvBufferRatio:       1.5,
		PvMinSizeBytes:      int64(1) << 30,
	}

	b := &VrsBuilder{Config: cfg, Rules: stubRules{}, CR: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: &v1.AgentPersistentVolume{
						Enabled: true,
					}},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)

	assert.NotNil(t, r.Spec.Agent.PersistentVolume)
	assert.Equal(t, "fallback-sc", r.Spec.Agent.PersistentVolume.StorageClass, "env fallback when CR omits SC")
	assert.Equal(t, "fallback-am", r.Spec.Agent.PersistentVolume.AccessMode, "env fallback when CR omits AM")
}
