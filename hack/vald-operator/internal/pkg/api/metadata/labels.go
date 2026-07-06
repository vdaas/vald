package metadata

import v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"

const (
	ManagedByLabel         = "app.kubernetes.io/managed-by"
	SubResourceLabelSuffix = "/managed-resource"
)

func CreateSubResourceLabels(name string) map[string]string {
	return map[string]string{
		ManagedByLabel: v1.GroupVersion.Group,
		v1.GroupVersion.Group + SubResourceLabelSuffix: name,
	}
}
