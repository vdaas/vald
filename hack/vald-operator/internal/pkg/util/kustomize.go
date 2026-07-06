package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/yaml"
)

type Patch struct {
	Path string `yaml:"path"`
}

type Kustomization struct {
	Resources []string `yaml:"resources"`
	Patches   []Patch  `yaml:"patches"`
}

func PatchMultipleUnstructured(patches ...unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if len(patches) == 0 {
		return nil, nil
	}

	if len(patches) == 1 {
		result := patches[0].DeepCopy()
		return result, nil
	}

	base := patches[0]
	actualPatches := patches[1:]

	fsInMemory := filesys.MakeFsInMemory()

	datas := make(map[string]unstructured.Unstructured)
	datas["base.yaml"] = base

	for i, patch := range actualPatches {
		datas[fmt.Sprintf("patch%d.yaml", i)] = patch
	}

	fs, err := writeKustomizeFiles(fsInMemory, datas, len(actualPatches))
	if err != nil {
		return nil, err
	}

	opts := krusty.MakeDefaultOptions()
	k := krusty.MakeKustomizer(opts)
	resMap, err := k.Run(fs, "/")
	if err != nil {
		return nil, fmt.Errorf("failed to run kustomize: %w", err)
	}

	us := &unstructured.Unstructured{}
	for _, r := range resMap.Resources() {
		obj, err := r.Map()
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource to map: %w", err)
		}
		us.Object = obj
	}
	return us, nil
}

func createKustomizationForPatches(patchCount int) ([]byte, error) {
	patches := make([]Patch, patchCount)
	for i := 0; i < patchCount; i++ {
		patches[i] = Patch{Path: fmt.Sprintf("patch%d.yaml", i)}
	}

	ks := Kustomization{
		Resources: []string{"base.yaml"},
		Patches:   patches,
	}
	return yaml.Marshal(&ks)
}

func writeKustomizeFiles(fs filesys.FileSystem, data map[string]unstructured.Unstructured, patchCount int) (filesys.FileSystem, error) {
	files, err := makeKustomizeFiles(data, patchCount)
	if err != nil {
		return nil, fmt.Errorf("failed to make kustomize files: %w", err)
	}

	for name, data := range files {
		if err := fs.WriteFile(name, data); err != nil {
			return nil, fmt.Errorf("failed to write file %s: %w", name, err)
		}
	}
	return fs, nil
}

func makeKustomizeFiles(files map[string]unstructured.Unstructured, patchCount int) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for name, obj := range files {
		data, err := yaml.Marshal(obj.Object)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal object %s: %w", name, err)
		}
		result[name] = data
	}

	kustomizationYaml, err := createKustomizationForPatches(patchCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create kustomization: %w", err)
	}
	result["kustomization.yaml"] = kustomizationYaml
	return result, nil
}
