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
