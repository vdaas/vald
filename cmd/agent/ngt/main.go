//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//


// Package main provides program main
package main

import (
	"context"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	ver "github.com/vdaas/vald/internal/version"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	"github.com/vdaas/vald/pkg/agent/ngt/usecase"
)

const (
	// version represent the version
	version    = "v0.0.1"
	maxVersion = "v0.0.10"
	minVersion = "v0.0.0"
)

func main() {
	defer safety.RecoverWithError(nil)

	log.Init(log.DefaultGlg())

	p, err := params.New(
		params.WithConfigFileDescription("agent config file path"),
	).Parse()

	if err != nil {
		log.Fatal(err)
		return
	}

	if p.ShowVersion() {
		log.Infof("server version -> %s", version)
		return
	}

	cfg, err := config.NewConfig(p.ConfigFilePath())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = ver.Check(cfg.Version, maxVersion, minVersion)
	if err != nil {
		log.Fatal(err)
		return
	}

	daemon, err := usecase.New(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = runner.Run(errgroup.Init(context.Background()), daemon)
	if err != nil {
		log.Fatal(err)
	}
}
