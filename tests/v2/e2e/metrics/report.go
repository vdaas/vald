// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"bytes"
	"context"
	"math"
	"path/filepath"
	"time"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

type (
	// Report represents the top-level structured E2E V2 performance metrics JSON report.
	Report struct {
		Summary       *GlobalSnapshot  `json:"summary"`
		SchemaVersion string           `json:"schema_version"`
		Scenario      string           `json:"scenario"`
		Timestamp     string           `json:"timestamp"`
		Strategies    []StrategyReport `json:"strategies"`
	}

	// StrategyReport represents the aggregated metrics for a single test strategy.
	StrategyReport struct {
		Summary    *GlobalSnapshot   `json:"summary"`
		Name       string            `json:"name"`
		Operations []OperationReport `json:"operations"`
	}

	// OperationReport represents the aggregated metrics for a single operation within a strategy.
	OperationReport struct {
		Summary    *GlobalSnapshot   `json:"summary"`
		Name       string            `json:"name"`
		Executions []ExecutionReport `json:"executions"`
	}

	// ExecutionReport represents the metrics for a single execution step within an operation.
	ExecutionReport struct {
		Summary *GlobalSnapshot `json:"summary"`
		Name    string          `json:"name"`
		Type    string          `json:"type"`
		Mode    string          `json:"mode"`
		Index   int             `json:"index"`
	}
)

const (
	// ReportSchemaVersion is the current schema version of the structured E2E metrics report.
	ReportSchemaVersion = "v1"

	jsonExtension         = ".json"
	defaultReportFilePerm = os.FileMode(0o644)
)

// SnapshotFromCollector returns the GlobalSnapshot for col, or an empty GlobalSnapshot if col is nil.
func SnapshotFromCollector(col Collector) *GlobalSnapshot {
	if col == nil {
		return &GlobalSnapshot{}
	}
	if snap := col.GlobalSnapshot(); snap != nil {
		return snap
	}
	return &GlobalSnapshot{}
}

// WriteReport serializes report as indented JSON and writes it to path,
// creating parent directories if needed and overwriting any existing file.
func WriteReport(ctx context.Context, path string, report *Report) error {
	if strings.TrimSpace(path) == "" {
		return errors.ErrPathNotSpecified
	}
	if report == nil {
		return errors.New("cannot write nil metrics report")
	}
	if ctx == nil {
		return errors.NewErrContextNotProvided
	}
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "context canceled before writing metrics report")
	}

	normalized := normalizeReport(report)
	b, err := json.MarshalIndent(&normalized, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal metrics report")
	}
	b = append(b, '\n')

	if _, err := file.OverWriteFile(ctx, path, bytes.NewReader(b), defaultReportFilePerm); err != nil {
		return errors.Wrapf(err, "failed to write metrics report to %s", path)
	}
	return nil
}

// LoadReport reads, parses, and validates a single v1 metrics JSON report file.
func LoadReport(path string) (*Report, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.ErrPathNotSpecified
	}
	fi, exist, err := statPath(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat report %s", path)
	}
	if !exist {
		return nil, errors.ErrFileNotFound(path)
	}
	if fi.IsDir() {
		return nil, errors.Errorf("report path %s is a directory", path)
	}
	if fi.Size() == 0 {
		return nil, errors.Errorf("report file %s is empty", path)
	}

	data, err := file.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read report %s", path)
	}

	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, errors.Wrapf(err, "failed to decode report %s", path)
	}
	if report.SchemaVersion != ReportSchemaVersion {
		return nil, errors.Errorf("unsupported report schema_version %q: expected %q", report.SchemaVersion, ReportSchemaVersion)
	}
	return &report, nil
}

// LoadReportsFromDir scans dir non-recursively for *.json report files, returning them keyed by filename stem.
func LoadReportsFromDir(dir string) (map[string]*Report, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.ErrPathNotSpecified
	}
	fi, exist, err := statPath(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat report directory %s", dir)
	}
	if !exist {
		return nil, errors.ErrDirectoryNotFound(err, dir, fi)
	}
	if !fi.IsDir() {
		return nil, errors.Errorf("report directory path %s is not a directory", dir)
	}

	entries, err := file.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read report directory %s", dir)
	}

	reports := make(map[string]*Report, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), jsonExtension) {
			continue
		}
		name := entry.Name()
		report, err := LoadReport(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		reports[strings.TrimSuffix(name, jsonExtension)] = report
	}
	return reports, nil
}

func normalizeReport(report *Report) Report {
	out := *report
	if out.SchemaVersion == "" {
		out.SchemaVersion = ReportSchemaVersion
	}
	if out.Timestamp == "" {
		out.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if out.Summary == nil {
		out.Summary = &GlobalSnapshot{}
	}
	if out.Strategies == nil {
		out.Strategies = make([]StrategyReport, 0)
		return out
	}
	strategies := make([]StrategyReport, len(out.Strategies))
	for i, st := range out.Strategies {
		if st.Summary == nil {
			st.Summary = &GlobalSnapshot{}
		}
		if st.Operations == nil {
			st.Operations = make([]OperationReport, 0)
		} else {
			ops := make([]OperationReport, len(st.Operations))
			for j, op := range st.Operations {
				if op.Summary == nil {
					op.Summary = &GlobalSnapshot{}
				}
				if op.Executions == nil {
					op.Executions = make([]ExecutionReport, 0)
				} else {
					execs := make([]ExecutionReport, len(op.Executions))
					for k, ex := range op.Executions {
						if ex.Summary == nil {
							ex.Summary = &GlobalSnapshot{}
						}
						execs[k] = ex
					}
					op.Executions = execs
				}
				ops[j] = op
			}
			st.Operations = ops
		}
		strategies[i] = st
	}
	out.Strategies = strategies
	return out
}

func statPath(path string) (os.FileInfo, bool, error) {
	exist, fi, err := file.ExistsWithDetail(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, false, err
	}
	if !exist || fi == nil {
		return nil, false, nil
	}
	return fi, true, nil
}

func finiteOrZero(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return v
}

func isFinitePositive(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v > 0
}

func isFiniteNonNegative(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0
}
