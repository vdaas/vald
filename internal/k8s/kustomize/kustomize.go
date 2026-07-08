//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package kustomize merges unstructured Kubernetes objects through an
// in-memory kustomize run, applying strategic-merge patch semantics.
package kustomize

import (
	"strconv"

	"github.com/vdaas/vald/internal/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/yaml"
)

// Patch refers to a patch file inside a Kustomization.
type Patch struct {
	Path string `yaml:"path"`
}

// Kustomization is the minimal kustomization.yaml model used by Merge.
type Kustomization struct {
	Resources []string `yaml:"resources"`
	Patches   []Patch  `yaml:"patches"`
}

// Merge overlays the given objects in order: the first element is the base
// and every following element is applied to it as a strategic-merge patch.
// A single input is returned as a deep copy; an empty input returns nil.
func Merge(overlays ...unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if len(overlays) == 0 {
		return nil, nil
	}

	if len(overlays) == 1 {
		result := overlays[0].DeepCopy()
		return result, nil
	}

	base := overlays[0]
	patches := overlays[1:]

	fsInMemory := filesys.MakeFsInMemory()

	datas := make(map[string]unstructured.Unstructured)
	datas["base.yaml"] = base

	for i, patch := range patches {
		datas["patch"+strconv.Itoa(i)+".yaml"] = patch
	}

	fs, err := writeKustomizeFiles(fsInMemory, datas, len(patches))
	if err != nil {
		return nil, err
	}

	opts := krusty.MakeDefaultOptions()
	k := krusty.MakeKustomizer(opts)
	resMap, err := k.Run(fs, "/")
	if err != nil {
		return nil, errors.Wrap(err, "failed to run kustomize")
	}

	// Merge overlays exactly one object, so anything other than a single
	// resource means the inputs disagreed (e.g. a patch targeting a different
	// object) and silently returning the last or an empty object would hide
	// the mistake.
	resources := resMap.Resources()
	if len(resources) != 1 {
		return nil, errors.Errorf("kustomize merge produced %d resources, want exactly 1", len(resources))
	}
	obj, err := resources[0].Map()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert resource to map")
	}
	return &unstructured.Unstructured{Object: obj}, nil
}

func createKustomizationForPatches(patchCount int) ([]byte, error) {
	patches := make([]Patch, patchCount)
	for i := range patchCount {
		patches[i] = Patch{Path: "patch" + strconv.Itoa(i) + ".yaml"}
	}

	ks := Kustomization{
		Resources: []string{"base.yaml"},
		Patches:   patches,
	}
	return yaml.Marshal(&ks)
}

func writeKustomizeFiles(
	fs filesys.FileSystem, data map[string]unstructured.Unstructured, patchCount int,
) (filesys.FileSystem, error) {
	files, err := makeKustomizeFiles(data, patchCount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make kustomize files")
	}

	for name, data := range files {
		if err := fs.WriteFile(name, data); err != nil {
			return nil, errors.Wrapf(err, "failed to write file %s", name)
		}
	}
	return fs, nil
}

func makeKustomizeFiles(
	files map[string]unstructured.Unstructured, patchCount int,
) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for name, obj := range files {
		data, err := yaml.Marshal(obj.Object)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal object %s", name)
		}
		result[name] = data
	}

	kustomizationYaml, err := createKustomizationForPatches(patchCount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kustomization")
	}
	result["kustomization.yaml"] = kustomizationYaml
	return result, nil
}
