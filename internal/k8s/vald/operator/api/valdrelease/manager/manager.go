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
package manager

import (
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/defaults"
	v1 "k8s.io/api/core/v1"
)

const (
	componentLabelManagerIndex = "manager-index"
)

type Manager struct {
	Index Index `json:"index"`
}

type Index struct {
	Enabled                   bool                          `json:"enabled"`
	Logging                   *defaults.Logging             `json:"logging,omitempty"`
	Indexer                   Indexer                       `json:"indexer"`
	Saver                     *Saver                        `json:"saver,omitempty"`
	Creator                   *Creator                      `json:"creator,omitempty"`
	Resources                 *v1.ResourceRequirements      `json:"resources,omitempty"`
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	Affinity                  *v1.Affinity                  `json:"affinity,omitempty"`
	NodeSelector              map[string]string             `json:"nodeSelector,omitempty"`
	Tolerations               *[]v1.Toleration              `json:"tolerations,omitempty"`
}

type Indexer struct {
	AutoIndexCheckDuration     string `json:"auto_index_check_duration,omitempty"`
	AutoIndexLength            *int   `json:"auto_index_length,omitempty"`
	AutoIndexDurationLimit     string `json:"auto_index_duration_limit,omitempty"`
	AutoSaveIndexDurationLimit string `json:"auto_save_index_duration_limit,omitempty"`
	AutoSaveIndexWaitDuration  string `json:"auto_save_index_wait_duration,omitempty"`
	Concurrency                *int   `json:"concurrency,omitempty"`
}

type Saver struct {
	Enabled      bool              `json:"enabled"`
	Schedule     string            `json:"schedule,omitempty"`
	Suspend      bool              `json:"suspend"`
	Affinity     *v1.Affinity      `json:"affinity,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Tolerations  *[]v1.Toleration  `json:"tolerations,omitempty"`
	Image        *common.Image     `json:"image,omitempty"`
}

type Creator struct {
	Enabled      bool              `json:"enabled"`
	Concurrency  *int              `json:"concurrency,omitempty"`
	Schedule     string            `json:"schedule,omitempty"`
	Suspend      bool              `json:"suspend"`
	Affinity     *v1.Affinity      `json:"affinity,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Tolerations  *[]v1.Toleration  `json:"tolerations,omitempty"`
	Image        *common.Image     `json:"image,omitempty"`
}
