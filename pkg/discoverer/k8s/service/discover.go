//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

// Package service manages the main logic of server.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/node"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/pkg/discoverer/k8s/model"
)

type Discoverer interface {
	GetServers() []model.Server
	Start(context.Context)
	Stop() error
}

type discoverer struct {
	ctrl k8s.Controller
}

func New(cfg *config.Discoverer) (Discoverer, error) {
	pw := pod.New()
	nw := node.New()
	d := &discoverer{
		ctrl: k8s.New(
			k8s.WithControllerName("vald k8s agent discoverer"),
			k8s.WithResourceController(pw),
			k8s.WithResourceController(nw),
		),
	}
	return nil, nil
}

func (d *discoverer) GetServers() []model.Server {
	return nil
}

func (d *discoverer) Start(ctx context.Context) <-chan error {
	return nil
}

func (d *discoverer) Stop() error {
	return nil
}
