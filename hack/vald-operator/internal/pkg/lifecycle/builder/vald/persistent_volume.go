package vald

import "github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease"

// reflectPersistentVolume configures persistent storage on the agent when the
// CR requests it. CR-supplied StorageClass / AccessMode win; env-provided
// defaults (b.Config.DefaultStorageClass / b.Config.DefaultAccessMode) fill
// in anything the CR omits. The PV size formula is owned by the domain.
func (b *VrsBuilder) reflectPersistentVolume(v *valdrelease.ValdRelease) {
	pv := b.CR.Spec.VectorEngine.Vald.Agent.PersistentVolume
	if pv == nil || !pv.Enabled {
		return
	}

	sc := pv.StorageClass
	if sc == "" {
		sc = b.Config.DefaultStorageClass
	}
	am := pv.AccessMode
	if am == "" {
		am = b.Config.DefaultAccessMode
	}
	memoryBytes := v.Spec.Agent.Resources.Requests.Memory().Value()
	size := b.Rules.AgentPvSize(memoryBytes, b.Config.PvBufferRatio, b.Config.PvMinSizeBytes)
	v.Spec.Agent.SetPvEnable(sc, am, size)
}
