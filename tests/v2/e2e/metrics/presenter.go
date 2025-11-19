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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/vdaas/vald/internal/net/grpc/codes"
	"gopkg.in/yaml.v2"
)

// SnapshotPresenter handles the formatting of a GlobalSnapshot into various output formats.
type SnapshotPresenter struct {
	snapshot *GlobalSnapshot
}

// NewSnapshotPresenter creates a new presenter for the given snapshot.
func NewSnapshotPresenter(snapshot *GlobalSnapshot) *SnapshotPresenter {
	return &SnapshotPresenter{
		snapshot: snapshot,
	}
}

// AsString returns a human-readable string representation of the snapshot.
func (p *SnapshotPresenter) AsString() string {
	if p.snapshot == nil || p.snapshot.Total == 0 {
		return "No data collected."
	}

	var sb strings.Builder
	s := p.snapshot
	total := s.Total
	errs := s.Errors
	totalDuration := time.Duration(s.Latencies.Sum)

	// --- Summary ---
	fmt.Fprint(&sb, "\n--- Summary ---\n")
	fmt.Fprintf(&sb, "Total Requests:\t%d\n", total)
	fmt.Fprintf(&sb, "Total Duration:\t%s\n", totalDuration)
	if totalDuration.Seconds() > 0 {
		fmt.Fprintf(&sb, "Requests/sec:\t%.2f\n", float64(total)/totalDuration.Seconds())
	}
	fmt.Fprintf(&sb, "Errors:\t%d (%.2f%%)\n", errs, float64(errs)/float64(total)*100)

	// --- Latency ---
	fmt.Fprint(&sb, "\n--- Latency ---\n")
	sb.WriteString(p.renderHistogram("Latency", s.Latencies, s.LatPercentiles))

	// --- Queue Wait ---
	fmt.Fprint(&sb, "\n--- Queue Wait ---\n")
	sb.WriteString(p.renderHistogram("Queue Wait", s.QueueWaits, s.QWPercentiles))

	// --- Status Codes ---
	fmt.Fprint(&sb, "\n--- Status Codes ---\n")
	if s.Codes != nil {
		codes := make([]codes.Code, 0, len(s.Codes))
		for code := range s.Codes {
			codes = append(codes, code)
		}
		slices.Sort(codes)
		for _, code := range codes {
			count := s.Codes[code]
			fmt.Fprintf(&sb, "\t- %s:\t%d (%.2f%%)\n", code.String(), count, float64(count)/float64(total)*100)
		}
	}

	// --- Exemplars ---
	fmt.Fprintf(&sb, "\n--- Exemplars (Top %d slowest requests) ---\n", len(s.Exemplars))
	for _, ex := range s.Exemplars {
		fmt.Fprintf(&sb, "\t- RequestID:\t%s,\tLatency:\t%s\n", ex.requestID, ex.latency)
	}

	return sb.String()
}

// AsJSON returns a JSON representation of the snapshot.
func (p *SnapshotPresenter) AsJSON() (string, error) {
	if p.snapshot == nil {
		return "null", nil
	}
	b, err := json.MarshalIndent(p.snapshot, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// AsYAML returns a YAML representation of the snapshot.
func (p *SnapshotPresenter) AsYAML() (string, error) {
	if p.snapshot == nil {
		return "null", nil
	}
	b, err := yaml.Marshal(p.snapshot)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// AsCSV returns a CSV representation of the snapshot's summary and percentiles.
func (p *SnapshotPresenter) AsCSV() (string, error) {
	return p.asSeparatedValue(',')
}

// AsTSV returns a TSV representation of the snapshot's summary and percentiles.
func (p *SnapshotPresenter) AsTSV() (string, error) {
	return p.asSeparatedValue('\t')
}

func (p *SnapshotPresenter) asSeparatedValue(separator rune) (string, error) {
	if p.snapshot == nil {
		return "", nil
	}

	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	writer.Comma = separator

	s := p.snapshot

	headers := []string{
		"Total", "Errors", "TotalDurationSec", "RPS", "ErrorRate",
		"LatencyMin", "LatencyMean", "LatencyMax",
		"LatencyP50", "LatencyP90", "LatencyP99",
		"QueueWaitMin", "QueueWaitMean", "QueueWaitMax",
		"QueueWaitP50", "QueueWaitP90", "QueueWaitP99",
	}
	writer.Write(headers)

	totalDuration := time.Duration(s.Latencies.Sum).Seconds()
	rps := 0.0
	if totalDuration > 0 {
		rps = float64(s.Total) / totalDuration
	}

	row := []string{
		fmt.Sprintf("%d", s.Total),
		fmt.Sprintf("%d", s.Errors),
		fmt.Sprintf("%.4f", totalDuration),
		fmt.Sprintf("%.2f", rps),
		fmt.Sprintf("%.4f", float64(s.Errors)/float64(s.Total)),
		fmt.Sprintf("%.4f", float64(s.Latencies.Min)/1e9),
		fmt.Sprintf("%.4f", float64(s.Latencies.Mean)/1e9),
		fmt.Sprintf("%.4f", float64(s.Latencies.Max)/1e9),
		fmt.Sprintf("%.4f", s.LatPercentiles.Quantile(0.5)/1e9),
		fmt.Sprintf("%.4f", s.LatPercentiles.Quantile(0.9)/1e9),
		fmt.Sprintf("%.4f", s.LatPercentiles.Quantile(0.99)/1e9),
		fmt.Sprintf("%.4f", float64(s.QueueWaits.Min)/1e9),
		fmt.Sprintf("%.4f", float64(s.QueueWaits.Mean)/1e9),
		fmt.Sprintf("%.4f", float64(s.QueueWaits.Max)/1e9),
		fmt.Sprintf("%.4f", s.QWPercentiles.Quantile(0.5)/1e9),
		fmt.Sprintf("%.4f", s.QWPercentiles.Quantile(0.9)/1e9),
		fmt.Sprintf("%.4f", s.QWPercentiles.Quantile(0.99)/1e9),
	}
	writer.Write(row)
	writer.Flush()
	return sb.String(), nil
}

// renderHistogram is a helper to render the histogram part of the string output.
func (p *SnapshotPresenter) renderHistogram(title string, h *HistogramSnapshot, q QuantileSketch) string {
	var sb strings.Builder

	if h != nil {
		fmt.Fprint(&sb, h.String())
	}
	if q != nil {
		fmt.Fprint(&sb, q.String())
	}
	if h != nil && len(h.Counts) > 0 {
		fmt.Fprint(&sb, "Histogram:\n")
		maxCount := uint64(0)
		for _, count := range h.Counts {
			if count > maxCount {
				maxCount = count
			}
		}
		for i, count := range h.Counts {
			var bar string
			if maxCount > 0 {
				bar = strings.Repeat("∎", int(float64(count)/float64(maxCount)*40))
			}
			var lowerBound, upperBound string
			if i == 0 {
				lowerBound = "0"
			} else {
				lowerBound = fmt.Sprintf("%.3f", float64(time.Duration(h.Bounds[i-1])))
			}
			if i == len(h.Bounds) {
				upperBound = "inf"
			} else {
				upperBound = fmt.Sprintf("%.3f", float64(time.Duration(h.Bounds[i])))
			}
			fmt.Fprintf(&sb, "\t%s - %s [%d]\t|%s\n", lowerBound, upperBound, count, bar)
		}
	}
	return sb.String()
}
