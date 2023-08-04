package main

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/index/job/correction/usecase"


	"github.com/vdaas/vald/pkg/manager/index/config" // FIXME: あとで独自のconfigに切り替え
)

const (
	maxVersion = "v0.0.10"
	minVersion = "v0.0.0"
	name       = "index correction job"
)

func main() {
	// FIXME: demon前提なので基本的に止まらない。独自のrunnerを作る必要があるか
	if err := safety.RecoverFunc(func() error {
		return runner.Do(
			context.Background(),
			runner.WithName(name),
			runner.WithVersion(info.Version, maxVersion, minVersion),
			runner.WithConfigLoader(func(path string) (interface{}, *config.GlobalConfig, error) {
				cfg, err := config.NewConfig(path)
				if err != nil {
					return nil, nil, errors.Wrap(err, "failed to load "+name+"'s configuration")
				}
				return cfg, &cfg.GlobalConfig, nil
			}),
			runner.WithDaemonInitializer(func(cfg interface{}) (runner.Runner, error) {
				return usecase.New()
			}),
		)
	})(); err != nil {
		log.Fatal(err, info.Get())
		return
	}
}
