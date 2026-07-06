package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestNodePools_GetNodePool(t *testing.T) {
	pool := NodePool{Name: "general-pool", Replicas: 2}
	pools := NodePools{
		NodePoolTypeGeneral: pool,
	}

	t.Run("existing type returns pool", func(t *testing.T) {
		got := pools.GetNodePool(NodePoolTypeGeneral)
		assert.NotNil(t, got)
		assert.Equal(t, pool.Name, got.Name)
	})

	t.Run("non-existent type returns nil", func(t *testing.T) {
		got := pools.GetNodePool(NodePoolTypeValdAgent)
		assert.Nil(t, got)
	})

	t.Run("empty NodePools returns nil", func(t *testing.T) {
		var empty NodePools
		got := empty.GetNodePool(NodePoolTypeGeneral)
		assert.Nil(t, got)
	})
}

func TestMachineResource_GetResourceList(t *testing.T) {
	t.Run("with storage", func(t *testing.T) {
		mr := &MachineResource{
			Cpu:     "500m",
			Memory:  "1Gi",
			Storage: "10Gi",
		}
		rl := mr.GetResourceList()
		assert.Equal(t, resource.MustParse("500m"), rl[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("1Gi"), rl[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("10Gi"), rl[corev1.ResourceStorage])
	})

	t.Run("without storage omits storage key", func(t *testing.T) {
		mr := &MachineResource{
			Cpu:    "1",
			Memory: "2Gi",
		}
		rl := mr.GetResourceList()
		assert.Equal(t, resource.MustParse("1"), rl[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("2Gi"), rl[corev1.ResourceMemory])
		_, hasStorage := rl[corev1.ResourceStorage]
		assert.False(t, hasStorage)
	})
}
