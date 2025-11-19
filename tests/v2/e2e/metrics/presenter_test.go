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

package metrics

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/net/grpc/codes"
)

var update = flag.Bool("update", false, "update golden files")

func TestSnapshotPresenter(t *testing.T) {
	t.Parallel()

	// Create a sample snapshot to test with
	snap := &GlobalSnapshot{
		Total:  100,
		Errors: 10,
		Latencies: &HistogramSnapshot{
			Total: 100,
			Sum:   float64(100 * time.Millisecond),
			Mean:  float64(1 * time.Millisecond),
			Min:   float64(100 * time.Microsecond),
			Max:   float64(10 * time.Millisecond),
		},
		QueueWaits: &HistogramSnapshot{
			Total: 100,
			Sum:   float64(50 * time.Millisecond),
			Mean:  float64(500 * time.Microsecond),
			Min:   float64(50 * time.Microsecond),
			Max:   float64(5 * time.Millisecond),
		},
		LatPercentiles: func() QuantileSketch {
			t, _ := NewTDigest(WithTDigestCompression(100))
			t.centroids = []centroid{{Mean: 1e6, Weight: 1}}
			return t
		}(),
		QWPercentiles: func() QuantileSketch {
			t, _ := NewTDigest(WithTDigestCompression(100))
			t.centroids = []centroid{{Mean: 5e5, Weight: 1}}
			return t
		}(),
		Codes: map[codes.Code]uint64{
			codes.OK:    90,
			codes.Aborted: 10,
		},
		Exemplars: []*item{
			{latency: 10 * time.Millisecond, requestID: "req-1"},
		},
	}

	p := NewSnapshotPresenter(snap)

	checkGoldenFile(t, "AsString.golden", p.AsString())

	json, err := p.AsJSON()
	if err != nil {
		t.Fatal(err)
	}
	checkGoldenFile(t, "AsJSON.golden", json)

	yaml, err := p.AsYAML()
	if err != nil {
		t.Fatal(err)
	}
	checkGoldenFile(t, "AsYAML.golden", yaml)

	csv, err := p.AsCSV()
	if err != nil {
		t.Fatal(err)
	}
	checkGoldenFile(t, "AsCSV.golden", csv)

	tsv, err := p.AsTSV()
	if err != nil {
		t.Fatal(err)
	}
	checkGoldenFile(t, "AsTSV.golden", tsv)
}

func checkGoldenFile(t *testing.T, goldenFile string, actual string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", goldenFile)
	if *update {
		err := os.MkdirAll(filepath.Dir(goldenPath), 0755)
		if err != nil {
			t.Fatalf("failed to create testdata dir: %v", err)
		}
		err = os.WriteFile(goldenPath, []byte(actual), 0644)
		if err != nil {
			t.Fatalf("failed to update golden file: %v", err)
		}
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	if string(golden) != actual {
		t.Errorf("output does not match golden file %s", goldenFile)
	}
}
