// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
