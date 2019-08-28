package main

import (
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/discoverer/openstack/config"
	"github.com/vdaas/vald/pkg/discoverer/openstack/usecase"
)

func main() {
	// Try recover befor kill process for dump panic errors
	defer safety.Recover()

	log.Init(log.DefaultGlg())

	p, err := params.New(
		params.WithConfigFileDescription("openstack discoverer config file path"),
	).Parse()

	if err != nil {
		log.Fatal(err)
		return
	}

	if p.ShowVersion() {
		err = log.Infof("server version -> %s", config.GetVersion())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	var stg config.Data
	err := config.New(p.ConfigFilePath(), &stg)
	if err != nil {
		log.Fatal(err)
		return
	}

	if stg.Version != config.GetVersion() {
		log.Fatal(errors.ErrInvalidConfig)
		return
	}

	daemon, err := usecase.New(stg)
	if err != nil {
		log.Fatal(err)
		return
	}

	errs := runner.Run(context.Background(), daemon)
	if len(errs) > 0 {
		log.Fatal(errs)
		return
	}
}
