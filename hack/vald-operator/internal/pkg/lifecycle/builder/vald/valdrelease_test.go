package vald

import (
	"context"
	"os"
	"testing"

	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func loadMvaldreleaseFromYAML(t *testing.T, path string) *v1.Mvaldrelease {
	data, err := os.ReadFile(path)
	assert.NoError(t, err)

	var cr v1.Mvaldrelease
	err = yaml.Unmarshal(data, &cr)
	assert.NoError(t, err)
	return &cr
}

func TestVrsBuilder_Validate(t *testing.T) {
	validCR := func() *v1.Mvaldrelease {
		return &v1.Mvaldrelease{
			Spec: v1.MvaldreleaseSpec{
				Infrastructure: []v1.MvaldreleaseInfra{
					{
						Role: "green",
						Clusters: []v1.DestClusters{
							{ID: "abc-123", Name: "cluster-a"},
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name    string
		cr      *v1.Mvaldrelease
		wantErr bool
	}{
		{
			name:    "nil CR",
			cr:      nil,
			wantErr: true,
		},
		{
			name: "empty infrastructure",
			cr: &v1.Mvaldrelease{
				Spec: v1.MvaldreleaseSpec{Infrastructure: nil},
			},
			wantErr: true,
		},
		{
			name: "infra with empty clusters",
			cr: &v1.Mvaldrelease{
				Spec: v1.MvaldreleaseSpec{
					Infrastructure: []v1.MvaldreleaseInfra{
						{Role: "green", Clusters: []v1.DestClusters{}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "cluster with empty Name",
			cr: func() *v1.Mvaldrelease {
				cr := validCR()
				cr.Spec.Infrastructure[0].Clusters[0].Name = ""
				return cr
			}(),
			wantErr: true,
		},
		{
			name:    "valid CR",
			cr:      validCR(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &VrsBuilder{CR: tt.cr}
			err := b.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVrsBuilder_Build(t *testing.T) {
	// Initialize config for testing
	config.DefaultVrsPath = "testdata/vrs.yaml"
	cfg, err := config.New()
	assert.NoError(t, err)

	cr := loadMvaldreleaseFromYAML(t, "testdata/mvrs.yaml")

	b := NewVrsBuilder(cr, cfg, AlwaysAvailable(), stubRules{})
	objList, err := b.Build(context.Background())
	assert.NoError(t, err)

	// debug output of the generated YAML
	//yamlData, err := yaml.Marshal(objList)
	//assert.NoError(t, err)
	//fmt.Printf("---\n%s\n", yamlData)

	// ゴールデンファイルのパス
	goldenPath := "testdata/vrs.golden.yaml"

	// Itemsの数
	assert.NotNil(t, objList)
	items := objList.(*unstructured.UnstructuredList).Items
	assert.Equal(t, 4, len(items))

	topItem := items[0]
	assert.Equal(t, cr.Namespace, topItem.GetNamespace())

	itemData, err := yaml.Marshal(topItem.Object)
	assert.NoError(t, err)
	expected, err := os.ReadFile(goldenPath)
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(itemData))
	// debug output of the first topItem to file
	err = os.WriteFile("testdata/vrs.golden.yaml", itemData, 0644)
	assert.NoError(t, err)
}
