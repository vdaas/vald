package vald

import (
	"fmt"
	"strings"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
)

func (b *VrsBuilder) makeName(str ...string) (string, error) {
	name := strings.Join(str, "-")
	if len(name) > 63 {
		return "", fmt.Errorf("the name %s is too long, must be 63 characters or less", name)
	}
	return name, nil
}

func (b *VrsBuilder) buildLabels(iKey int) map[string]string {
	return map[string]string{
		b.labelKey(nodePoolLabelType): string(b.CR.Spec.Infrastructure[iKey].Type),
		b.labelKey(nodePoolLabelRole): string(b.CR.Spec.Infrastructure[iKey].Role),
	}
}

func (b *VrsBuilder) buildLogging(ll string) *defaults.Logging {
	vald := b.CR.Spec.VectorEngine.Vald
	l := ll
	if l == "" {
		if vald.Defaults.LogLevel != "" {
			l = vald.Defaults.LogLevel
		} else {
			l = b.Config.VrsLogLevel
		}
	}
	return &defaults.Logging{
		Level:  l,
		Format: "raw",
		Logger: "glg",
	}
}

// labelKey renders a NodePool label key with the configured prefix. Exposed as a
// free function (rather than a VrsBuilder method) because tests build matching
// labels on fake Node objects before a builder exists; both call sites must use
// the same logic and prefix value.
func labelKey(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	return fmt.Sprintf("%s/%s", prefix, suffix)
}

func (b *VrsBuilder) labelKey(suffix string) string {
	return labelKey(b.Config.NodePoolLabelPrefix, suffix)
}
