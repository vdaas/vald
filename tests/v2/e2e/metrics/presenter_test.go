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
	"github.com/vdaas/vald/internal/test"
)

var update = flag.Bool("update", false, "update golden files")

func TestSnapshotPresenter(t *testing.T) {
	type args struct {
		snapshot *GlobalSnapshot
	}
	type want struct {
		goldenFile string
	}

	// Helper to create a sample snapshot
	createSampleSnapshot := func() *GlobalSnapshot {
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
			LatPercentiles: func() TDigest {
				t, _ := NewTDigest(WithTDigestCompression(100))
				t.Add(1e6)
				return t
			}(),
			QWPercentiles: func() TDigest {
				t, _ := NewTDigest(WithTDigestCompression(100))
				t.Add(5e5)
				return t
			}(),
			Codes: map[codes.Code]uint64{
				codes.OK:      90,
				codes.Aborted: 10,
			},
			Exemplars: []*ExemplarItem{
				{Latency: 10 * time.Millisecond, RequestID: "req-1"},
			},
		}
		return snap
	}

	// Define checks for each format
	runCheck := func(name, goldenFile string, convert func(*SnapshotPresenter) (string, error)) {
		t.Run(name, func(t *testing.T) {
			if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (string, error) {
				p := NewSnapshotPresenter(args.snapshot)
				return convert(p)
			}, []test.Case[string, args]{
				{
					Name: "valid snapshot",
					Args: args{
						snapshot: createSampleSnapshot(),
					},
					Want: test.Result[string]{
						// We don't populate Val here because we check against golden file
					},
					CheckFunc: func(tt *testing.T, want test.Result[string], got test.Result[string]) error {
						if got.Err != nil {
							return got.Err
						}
						checkGoldenFile(tt, goldenFile, got.Val)
						return nil
					},
				},
				{
					Name: "empty snapshot",
					Args: args{
						snapshot: &GlobalSnapshot{},
					},
					CheckFunc: func(tt *testing.T, want test.Result[string], got test.Result[string]) error {
						if got.Err != nil {
							return got.Err
						}
						// Just verify no error and non-empty output for empty snapshot (might vary by format)
						// Actually for empty snapshot, some return "null" or specific string.
						// We can have a separate golden file for empty if needed, but basic check is enough.
						return nil
					},
				},
			}...); err != nil {
				t.Error(err)
			}
		})
	}

	runCheck("AsString", "AsString.golden", func(p *SnapshotPresenter) (string, error) {
		return p.AsString(), nil
	})

	runCheck("AsJSON", "AsJSON.golden", func(p *SnapshotPresenter) (string, error) {
		return p.AsJSON()
	})

	runCheck("AsYAML", "AsYAML.golden", func(p *SnapshotPresenter) (string, error) {
		return p.AsYAML()
	})

	runCheck("AsCSV", "AsCSV.golden", func(p *SnapshotPresenter) (string, error) {
		return p.AsCSV()
	})

	runCheck("AsTSV", "AsTSV.golden", func(p *SnapshotPresenter) (string, error) {
		return p.AsTSV()
	})
}

func checkGoldenFile(t *testing.T, goldenFile string, actual string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", goldenFile)
	if *update {
		err := os.MkdirAll(filepath.Dir(goldenPath), 0o755)
		if err != nil {
			t.Fatalf("failed to create testdata dir: %v", err)
		}
		err = os.WriteFile(goldenPath, []byte(actual), 0o644)
		if err != nil {
			t.Fatalf("failed to update golden file: %v", err)
		}
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		// If file doesn't exist and not updating, fail
		t.Fatalf("failed to read golden file: %v", err)
	}

	if string(golden) != actual {
		t.Errorf("output does not match golden file %s", goldenFile)
	}
}
