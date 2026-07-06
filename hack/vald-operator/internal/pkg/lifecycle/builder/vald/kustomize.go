package vald

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/util"
)

func (b *VrsBuilder) mergeOverlay(row client.Object) (*unstructured.Unstructured, error) {
	current, err := util.ConvertToUnstructured(row)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base to unstructured: %w", err)
	}

	// Copy the parsed default VRS before mutating; the cached version must remain intact.
	baseVrs := b.Config.DefaultVrs.Us.DeepCopy()
	baseVrs.SetName(current.GetName())
	baseVrs.SetNamespace(current.GetNamespace())

	patches := []unstructured.Unstructured{*current}

	patch, err := b.makeOverlayPatch(row)
	if err != nil {
		return nil, fmt.Errorf("failed to make overlay patch: %w", err)
	}

	if patch != nil {
		patches = append(patches, *patch)
	}

	allPatches := append([]unstructured.Unstructured{*baseVrs}, patches...)
	return util.PatchMultipleUnstructured(allPatches...)
}

func (b *VrsBuilder) makeOverlayPatch(row client.Object) (*unstructured.Unstructured, error) {
	patchRaw := b.CR.Spec.VectorEngine.Vald.Overlay.Raw
	if len(patchRaw) == 0 {
		return nil, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(patchRaw, &obj); err != nil {
		return nil, err
	}
	patch := &unstructured.Unstructured{Object: obj}
	patch.SetGroupVersionKind(b.GVK)
	patch.SetName(row.GetName())
	patch.SetNamespace(row.GetNamespace())

	return patch, nil
}
