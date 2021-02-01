//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package ngtd provides ngtd starter  functionality
package ngtd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd"
	"github.com/yahoojapan/ngtd/kvs"
)

type ServerType = ngtd.ServerType

const (
	HTTP ServerType = 1
	GRPC ServerType = 2
)

type server struct {
	dim      int
	srvType  ServerType
	indexDir string
	port     int
}

func New(opts ...Option) starter.Starter {
	srv := new(server)
	for _, opt := range append(defaultOptions, opts...) {
		opt(srv)
	}
	return srv
}

func (ns *server) Run(ctx context.Context, tb testing.TB) func() {
	tb.Helper()

	if err := ns.createIndexDir(); err != nil {
		tb.Error(err)
	}

	gongt.SetDimension(ns.dim)

	db, err := kvs.NewGoLevel(ns.indexDir + "meta")
	if err != nil {
		tb.Error(err)
	}

	n, err := ngtd.NewNGTD(ns.indexDir+"ngt", db, ns.port)
	if err != nil {
		tb.Error(err)
	}

	go func() {
		if err := n.ListenAndServe(ns.srvType); err != nil {
			tb.Errorf("ngtd returned error: %s", err.Error())
		}
	}()

	time.Sleep(4 * time.Second)

	return func() {
		n.Stop()

		if err := ns.clearIndexDir(); err != nil {
			tb.Error(err)
		}
	}
}

func (ns *server) createIndexDir() error {
	if err := os.RemoveAll(ns.indexDir); err != nil {
		return err
	}

	if err := os.MkdirAll(ns.indexDir, os.ModeTemporary); err != nil {
		return err
	}

	return nil
}

func (ns *server) clearIndexDir() error {
	if err := os.RemoveAll(ns.indexDir + "meta"); err != nil {
		return err
	}

	if err := os.RemoveAll(ns.indexDir + "ngt"); err != nil {
		return err
	}

	return nil
}
