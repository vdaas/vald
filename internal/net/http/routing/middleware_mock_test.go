package routing

import "github.com/vdaas/vald/internal/net/http/rest"

type middlewareMock struct {
	WrapFunc func(rest.Func) rest.Func
}

func (mo *middlewareMock) Wrap(f rest.Func) rest.Func {
	return mo.WrapFunc(f)
}
