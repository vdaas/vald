package lifecycle

import "github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"

// LifeCycle is one phase of the reconcile pipeline. Either Builder or Checker
// (or both) may be set:
//   - Builder != nil  : the phase produces Kubernetes objects to apply
//   - Checker != nil  : the phase observes a property and reports readiness
//   - both nil        : a terminal/no-op phase (e.g. Completed)
type LifeCycle struct {
	Condition Condition
	Builder   desired.Builder
	Checker   desired.ReadinessChecker
}

type LifeCycles []LifeCycle

type Flow struct {
	LifeCycles LifeCycles
	Current    int
}

func NewFlow(lc LifeCycles, c int) *Flow {
	return &Flow{
		LifeCycles: lc,
		Current:    c,
	}
}

func (lcs *LifeCycles) GetIndex(ct string) int {
	if ct == "" {
		return 0 // Default to first phase if no type is provided
	}
	for i, cond := range *lcs {
		if cond.Condition.Type == ct {
			return i
		}
	}
	return -1 // Not found
}

func (f *Flow) GetNext() int {
	if f.Current+1 < len(f.LifeCycles) {
		return f.Current + 1
	}
	return -1 // No next phase
}
