//go:build e2e

//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"flag"
	"testing"

	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/dataset"
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
		params.WithArgumentFilters(
			func(s string) bool {
				return strings.HasPrefix(s, "-test.")
			},
		),
	).Parse()
	if fail || err != nil {
		log.Fatalf("failed to parse the parameters: %v", err)
	}
	// params filtered the -test.* arguments out of its own flag set above, so
	// they must be fed to testing's real flag.CommandLine (its flags are
	// already registered by testing.MainStart before TestMain runs) for
	// -test.run/-test.bench/-test.timeout/-test.v/... to take effect: m.Run
	// only calls flag.Parse when flag.Parsed() is still false, and the
	// previously used params.WithOverrideDefault(true) swapped
	// flag.CommandLine for params' already-parsed set, which silently left
	// every -test.* flag at its zero value (all tests always ran, no
	// benchmark could ever be selected and -timeout was never applied).
	testArgs := make([]string, 0, len(os.Args[1:]))
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-test.") {
			testArgs = append(testArgs, arg)
		}
	}
	// flag.CommandLine uses flag.ExitOnError, so a malformed -test.* flag
	// exits inside Parse itself; this branch only fires for flag.ErrHelp.
	if err := flag.CommandLine.Parse(testArgs); err != nil {
		log.Fatalf("failed to parse -test.* flags: %v", err)
	}
	if p.ShowVersion() {
		log.Info(info.Version)
		os.Exit(0)
	}
	if p.ConfigFilePath() == "" {
		log.Fatalf("config file is necessary for e2e tests")
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
	if cfg.Dataset != nil && cfg.Dataset.Name != "" {
		ds, err = dataset.ToDataset(cfg.Dataset)
		if err != nil {
			log.Fatalf("failed to load dataset: %v", err)
		}
	} else if cfg.Dataset != nil && cfg.Dataset.Dimension > 0 {
		// No HDF5 fixture requested (e.g. maximum vector dimension probing);
		// synthesize a single-vector dataset of the configured dimension instead.
		ds = newSyntheticDataset(cfg.Dataset.Dimension)
	} else {
		// dataset-less scenarios (e.g. operator verification) do not require hdf5 loading.
		log.Info("dataset name is empty, skipping dataset loading")
	}
	cfg.FilePath = fp
	os.Exit(m.Run())
}

// newSyntheticDataset builds a single-vector in-memory dataset for scenarios
// that need a vector of an arbitrary dimension without a pre-built HDF5
// fixture (e.g. maximum vector dimension probing, where fixtures would need
// to be multiple gigabytes in size for high dimensions).
func newSyntheticDataset(dim int) *hdf5.Dataset {
	vec := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
	// maxLen is set to 1 explicitly (matching the single synthesized vector)
	// so TrainCycle/TestCycle never fall into InitNoiseFunc's noise
	// generation path, which hdf5.ToDataset-backed datasets only enter once
	// num exceeds the size of the loaded fixture. This is a single-vector
	// dimension probe, so noise is an unintended side effect here.
	return hdf5.New([][]float32{vec}, [][]float32{vec}, [][]int{{0}}, 1)
}
