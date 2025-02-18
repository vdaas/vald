//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// package crud provides e2e tests using ann-benchmarks datasets
package crud

import (
	"os"
	"testing"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/hdf5"
)

var (
	cfg *config.Data
	ds  *hdf5.Dataset
)

func TestMain(m *testing.M) {
	var err error
	p, fail, err := params.New(
		params.WithName("vald/e2e"),
		params.WithOverrideDefault(true),
		params.WithArgumentFilters(
			func(s string) bool {
				return strings.HasPrefix(s, "-test.")
			},
		),
	).Parse()
	if fail || err != nil || p.ConfigFilePath() == "" {
		log.Fatalf("failed to parse the parameters: %v", err)
	}

	if testing.Short() {
		log.Info("skipping this pkg test when -short because e2e test takes a long time")
		os.Exit(0)
	}

	fp := p.ConfigFilePath()
	cfg, err = config.Load(fp)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Init(log.WithLevel(cfg.Logging.Level), log.WithFormat(cfg.Logging.Format))
	ds, err = hdf5.HDF5ToDataset(cfg.Dataset.Name)
	if err != nil {
		log.Fatalf("failed to load dataset: %v", err)
	}
	cfg.FilePath = fp
	os.Exit(m.Run())
}
