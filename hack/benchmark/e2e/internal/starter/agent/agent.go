package agent

import "github.com/vdaas/vald/hack/benchmark/e2e/internal/starter"

type server struct{}

func New(opts ...Option) starter.Starter {
	srv := new()
	return nil
}
