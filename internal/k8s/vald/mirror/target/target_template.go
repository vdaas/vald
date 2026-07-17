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

import "github.com/vdaas/vald/internal/k8s/vald"

func NewMirrorTargetTemplate(opts ...MirrorTargetTemplateOption) (*MirrorTarget, error) {
	mt := new(MirrorTarget)
	for _, opt := range append(defaultMirrorTargetTemplateOptions, opts...) {
		if err := opt(mt); err != nil {
			if abort, oerr := vald.SkipNonCriticalOptionError(err, opt); abort {
				return nil, oerr
			}
		}
	}
	return mt, nil
}
