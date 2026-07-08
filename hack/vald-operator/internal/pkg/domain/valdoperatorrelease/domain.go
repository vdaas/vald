package valdoperatorrelease

import (
	v1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
)

type Domain struct {
	*v1.ValdOperatorRelease `json:",inline" yaml:",inline"`
}

// Status初期化
func (d *Domain) InitProgress(lcs lifecycle.LifeCycles) {
	cr := d.ValdOperatorRelease
	d.Status.Progress.Total = len(lcs) - 1 // Completedで行うことは何もないのでTotal全体から1引かれる
	d.Status.Progress.Completed = lcs.GetIndex(cr.Status.Phase)
}
