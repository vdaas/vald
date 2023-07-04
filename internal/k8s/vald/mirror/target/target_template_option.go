// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

import "github.com/vdaas/vald/internal/errors"

type MirrorTargetOption func(*MirrorTarget) error

var defaultMirrorTargetOptions = []MirrorTargetOption{
	WithMirrorTargetLabels(map[string]string{
		"app.kubernetes.io/name":       "mirror-target",
		"app.kubernetes.io/managed-by": "gateway-mirror",
	}),
}

func WithMirrorTargetNamespace(ns string) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		if len(ns) != 0 {
			mt.ObjectMeta.Namespace = ns
		}
		return nil
	}
}

func WithMirrorTargetName(name string) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		if len(name) == 0 {
			return errors.NewErrCriticalOption("name", name)
		}
		mt.ObjectMeta.Name = name
		return nil
	}
}

func WithMirrorTargetStatus(st *MirrorTargetStatus) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		mt.Status = *st
		return nil
	}
}

func WithMirrorTargetLabels(labels map[string]string) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		if len(labels) != 0 {
			mt.ObjectMeta.Labels = labels
		}
		return nil
	}
}

func WithMirrorTargetColocation(n string) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		if len(n) == 0 {
			return errors.NewErrCriticalOption("colocation", n)
		}
		mt.Spec.Colocation = n
		return nil
	}
}

func WithMirrorTargetHost(n string) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		if len(n) == 0 {
			return errors.NewErrCriticalOption("host", n)
		}
		mt.Spec.Target.Host = n
		return nil
	}
}

func WithMirrorTargetPort(port int) MirrorTargetOption {
	return func(mt *MirrorTarget) error {
		if port <= 0 {
			return errors.NewErrCriticalOption("port", port)
		}
		mt.Spec.Target.Port = port
		return nil
	}
}
