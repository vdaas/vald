//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

package usecase

import (
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/sidecar/config"
	"github.com/vdaas/vald/pkg/agent/sidecar/usecase/initcontainer"
	"github.com/vdaas/vald/pkg/agent/sidecar/usecase/sidecar"
)

func New(cfg *config.Data) (r runner.Runner, err error) {
	switch config.SidecarMode(cfg.AgentSidecar.Mode) {
	case config.INITCONTAINER:
		return initcontainer.New(cfg)
	case config.SIDECAR:
	default:
	}

	return sidecar.New(cfg)
}
