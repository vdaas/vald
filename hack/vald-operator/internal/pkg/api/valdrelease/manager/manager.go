// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package manager

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
	v1 "k8s.io/api/core/v1"
)

const (
	componentLabelManagerIndex = "manager-index"
)

// +schemagen:begin

type Manager struct {
	Index Index `json:"index"`
}

type Index struct {
	Enabled                   bool                          `json:"enabled"`
	Logging                   *defaults.Logging             `json:"logging,omitempty"`
	Indexer                   Indexer                       `json:"indexer,omitempty"`
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
	Image        string            `json:"image,omitempty"`
}

type Creator struct {
	Enabled      bool              `json:"enabled"`
	Concurrency  *int              `json:"concurrency,omitempty"`
	Schedule     string            `json:"schedule,omitempty"`
	Suspend      bool              `json:"suspend"`
	Affinity     *v1.Affinity      `json:"affinity,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Tolerations  *[]v1.Toleration  `json:"tolerations,omitempty"`
	Image        string            `json:"image,omitempty"`
}
