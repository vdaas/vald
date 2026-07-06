package desired

import "context"

// Prop is a Condition-only ReadinessChecker: it observes an external property
// and reports the result. It does not produce Kubernetes resources, so it
// does not implement Builder.
type Prop struct {
	Check func() Result
}

func (dp *Prop) IsReady(_ context.Context) Result {
	return dp.Check()
}
