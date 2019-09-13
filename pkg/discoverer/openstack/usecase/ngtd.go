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

package usecase

import (
	"context"

	"github.com/vdaas/vald/pkg/discoverer/openstack/config"
	"github.com/vdaas/vald/pkg/discoverer/openstack/service"
)

type Runner interface {
	Start(ctx context.Context) chan []error
}

type run struct {
	cfg    config.Data
	server service.Server
}

func New(cfg config.Data) (Runner, error) {
	return &run{
		cfg:    cfg,
		server: service.NewServer(nil),
	}, nil
}

func (t *run) Start(ctx context.Context) chan error {
	return t.server.ListenAndServe(ctx)
}
