package vald

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
)

func (b *VrsBuilder) buildDefaults() defaults.Defaults {
	vald := b.CR.Spec.VectorEngine.Vald
	return defaults.Defaults{
		Logging: b.buildLogging(vald.Defaults.LogLevel),
	}
}
