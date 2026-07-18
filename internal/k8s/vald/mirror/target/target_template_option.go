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

import "github.com/vdaas/vald/internal/errors"

type MirrorTargetTemplateOption func(*MirrorTarget) error

var defaultMirrorTargetTemplateOptions = []MirrorTargetTemplateOption{
	WithMirrorTargetLabels(map[string]string{
		"app.kubernetes.io/name":       "mirror-target",
		"app.kubernetes.io/managed-by": "gateway-mirror",
	}),
}

func WithMirrorTargetNamespace(ns string) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		if ns != "" {
			mt.Namespace = ns
		}
		return nil
	}
}

func WithMirrorTargetName(name string) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		if name == "" {
			return errors.NewErrCriticalOption("name", name)
		}
		mt.Name = name
		return nil
	}
}

func WithMirrorTargetStatus(status *MirrorTargetStatus) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		mt.Status = *status
		return nil
	}
}

func WithMirrorTargetLabels(labels map[string]string) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		if len(labels) != 0 {
			mt.Labels = labels
		}
		return nil
	}
}

func WithMirrorTargetColocation(colocation string) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		if colocation == "" {
			return errors.NewErrCriticalOption("colocation", colocation)
		}
		mt.Spec.Colocation = colocation
		return nil
	}
}

func WithMirrorTargetHost(host string) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		if host == "" {
			return errors.NewErrCriticalOption("host", host)
		}
		mt.Spec.Target.Host = host
		return nil
	}
}

func WithMirrorTargetPort(port int) MirrorTargetTemplateOption {
	return func(mt *MirrorTarget) error {
		if port <= 0 {
			return errors.NewErrCriticalOption("port", port)
		}
		mt.Spec.Target.Port = port
		return nil
	}
}
