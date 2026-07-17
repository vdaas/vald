// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package target

import (
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald"
	mirrv1 "github.com/vdaas/vald/internal/k8s/vald/mirror/api/v1"
)

type (
	MirrorTargetWatcher k8s.ResourceController
	MirrorTarget        = mirrv1.ValdMirrorTarget
	MirrorTargetStatus  = mirrv1.MirrorTargetStatus
	MirrorTargetPhase   = mirrv1.MirrorTargetPhase
)

const (
	MirrorTargetPhasePending      = mirrv1.MirrorTargetPending
	MirrorTargetPhaseConnected    = mirrv1.MirrorTargetConnected
	MirrorTargetPhaseDisconnected = mirrv1.MirrorTargetDisconnected
	MirrorTargetPhaseUnknown      = mirrv1.MirrorTargetUnknown
)

type Endpoint struct {
	Colocation string
	Host       string
	Phase      MirrorTargetPhase
	Port       int
}

// New creates a MirrorTargetWatcher that lists ValdMirrorTarget resources on
// every reconcile and reports them, converted into Endpoint values keyed by
// name, to the callback registered via WithOnReconcileFunc. Option
// application errors abort construction only when critical; any other option
// error is logged as a warning and skipped.
func New(opts ...Option) (MirrorTargetWatcher, error) {
	return vald.NewListWatcher(
		mirrv1.AddToScheme,
		func(m *mirrv1.ValdMirrorTarget) Endpoint {
			return Endpoint{
				Colocation: m.Spec.Colocation,
				Host:       m.Spec.Target.Host,
				Port:       m.Spec.Target.Port,
				Phase:      m.Status.Phase,
			}
		},
		vald.SkipNonCriticalOptionError,
		opts...,
	)
}
