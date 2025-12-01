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
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/strings"
	"sigs.k8s.io/yaml"
)

const (
	nanoToSec = 1e9
	barWidth  = 40
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

	p.renderSummary(&sb)
	if p.snapshot.Total > 1 {
		p.renderLatency(&sb)
		p.renderQueueWait(&sb)
	}
	p.renderStatusCodes(&sb)
	if p.snapshot.Total > 1 {
		p.renderExemplars(&sb)
	}

	return sb.String()
}

func (p *SnapshotPresenter) renderSummary(sb *strings.Builder) {
	s := p.snapshot
	total := s.Total
	errs := s.Errors

	totalDuration := s.LastUpdated.Sub(s.StartTime)

	// Errors (Fprint, Fprintf) are explicitly ignored for strings.Builder as it always returns nil error.
	_, _ = fmt.Fprint(sb, "\n--- Summary ---\n")
	_, _ = fmt.Fprintf(sb, "Total Requests:\t%d\n", total)
	_, _ = fmt.Fprintf(sb, "Total Duration:\t%s\n", totalDuration)
	if totalDuration.Seconds() > 0 {
		_, _ = fmt.Fprintf(sb, "Requests/sec:\t%.2f\n", float64(total)/totalDuration.Seconds())
	}
	_, _ = fmt.Fprintf(sb, "Errors:\t%d (%.2f%%)\n", errs, float64(errs)/float64(total)*percentMultiplier)
}

func (p *SnapshotPresenter) renderLatency(sb *strings.Builder) {
	_, _ = sb.WriteString(p.renderHistogram("Latency", p.snapshot.Latencies, p.snapshot.LatPercentiles))
}

func (p *SnapshotPresenter) renderQueueWait(sb *strings.Builder) {
	_, _ = sb.WriteString(p.renderHistogram("Queue Wait", p.snapshot.QueueWaits, p.snapshot.QWPercentiles))
}

func (p *SnapshotPresenter) renderStatusCodes(sb *strings.Builder) {
	s := p.snapshot
	total := s.Total
	_, _ = fmt.Fprint(sb, "\n--- Status Codes ---\n")
	if s.Codes != nil {
		codes := make([]codes.Code, 0, len(s.Codes))
		for code := range s.Codes {
			codes = append(codes, code)
		}
		slices.Sort(codes)
		for _, code := range codes {
			count := s.Codes[code]
			_, _ = fmt.Fprintf(sb, "\t- %s:\t%d (%.2f%%)\n", code.String(), count, float64(count)/float64(total)*percentMultiplier)
		}
	}
}

func (p *SnapshotPresenter) renderExemplars(sb *strings.Builder) {
	s := p.snapshot
	if s.ExemplarDetails != nil {
		renderExemplars := func(title string, items []*ExemplarItem) {
			if len(items) > 0 {
				_, _ = fmt.Fprintf(sb, "\n--- Exemplars (%s) ---\n", title)
				for _, ex := range items {
					status := ""
					if ex.Err != nil {
						status = " (Failed)"
					}
					_, _ = fmt.Fprintf(sb, "\t- RequestID:\t%s,\tLatency:\t%s%s\n", ex.RequestID, ex.Latency, status)
					if ex.Msg != "" {
						_, _ = fmt.Fprintf(sb, "\t  Message:\t%s\n", ex.Msg)
					}
				}
			}
		}
		renderExemplars("Slowest", s.ExemplarDetails.Slowest)
		renderExemplars("Fastest", s.ExemplarDetails.Fastest)
		renderExemplars("Average (Sampled)", s.ExemplarDetails.Average)
		renderExemplars("Failures", s.ExemplarDetails.Failures)
	} else if len(s.Exemplars) > 0 {
		_, _ = fmt.Fprintf(sb, "\n--- Exemplars (Top %d slowest requests) ---\n", len(s.Exemplars))
		for _, ex := range s.Exemplars {
			_, _ = fmt.Fprintf(sb, "\t- RequestID:\t%s,\tLatency:\t%s\n", ex.RequestID, ex.Latency)
			if ex.Msg != "" {
				_, _ = fmt.Fprintf(sb, "\t  Message:\t%s\n", ex.Msg)
			}
		}
	}
}

// AsJSON returns a JSON representation of the snapshot.
func (p *SnapshotPresenter) AsJSON() (string, error) {
	if p.snapshot == nil {
		return "null", nil
	}
	b, err := json.MarshalIndent(p.snapshot, "", "\t")
	if err != nil {
		return "", err
	}
	return conv.Btoa(b), nil
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
	}
	if s.LatPercentiles != nil {
		for _, q := range s.LatPercentiles.Quantiles() {
			headers = append(headers, fmt.Sprintf("LatencyP%d", int(q*percentMultiplier)))
		}
	}
	headers = append(headers, "QueueWaitMin", "QueueWaitMean", "QueueWaitMax")
	if s.QWPercentiles != nil {
		for _, q := range s.QWPercentiles.Quantiles() {
			headers = append(headers, fmt.Sprintf("QueueWaitP%d", int(q*percentMultiplier)))
		}
	}
	// G104: Handle error
	if err := writer.Write(headers); err != nil {
		return "", errors.Join(err, errors.New("failed to write headers"))
	}

	totalDuration := s.LastUpdated.Sub(s.StartTime).Seconds()
	rps := 0.0
	if totalDuration > 0 {
		rps = float64(s.Total) / totalDuration
	}

	errorRate := 0.0
	if s.Total > 0 {
		errorRate = float64(s.Errors) / float64(s.Total)
	}

	latMin, latMean, latMax := 0.0, 0.0, 0.0
	if s.Latencies != nil {
		latMin = float64(s.Latencies.Min) / nanoToSec
		latMean = float64(s.Latencies.Mean) / nanoToSec
		latMax = float64(s.Latencies.Max) / nanoToSec
	}

	row := []string{
		fmt.Sprintf("%d", s.Total),
		fmt.Sprintf("%d", s.Errors),
		fmt.Sprintf("%.4f", totalDuration),
		fmt.Sprintf("%.2f", rps),
		fmt.Sprintf("%.4f", errorRate),
		fmt.Sprintf("%.4f", latMin),
		fmt.Sprintf("%.4f", latMean),
		fmt.Sprintf("%.4f", latMax),
	}
	if s.LatPercentiles != nil {
		for _, q := range s.LatPercentiles.Quantiles() {
			row = append(row, fmt.Sprintf("%.4f", s.LatPercentiles.Quantile(q)/nanoToSec))
		}
	}

	qwMin, qwMean, qwMax := 0.0, 0.0, 0.0
	if s.QueueWaits != nil {
		qwMin = float64(s.QueueWaits.Min) / nanoToSec
		qwMean = float64(s.QueueWaits.Mean) / nanoToSec
		qwMax = float64(s.QueueWaits.Max) / nanoToSec
	}

	row = append(row,
		fmt.Sprintf("%.4f", qwMin),
		fmt.Sprintf("%.4f", qwMean),
		fmt.Sprintf("%.4f", qwMax),
	)
	if s.QWPercentiles != nil {
		for _, q := range s.QWPercentiles.Quantiles() {
			row = append(row, fmt.Sprintf("%.4f", s.QWPercentiles.Quantile(q)/nanoToSec))
		}
	}

	if err := writer.Write(row); err != nil {
		return "", errors.Join(err, errors.New("failed to write row"))
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// renderHistogram is a helper to render the histogram part of the string output.
func (p *SnapshotPresenter) renderHistogram(title string, h *HistogramSnapshot, q TDigest) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "\n--- %s ---\n", title)

	// Helper to convert nanoseconds to duration string
	fmtDur := func(ns float64) string {
		return time.Duration(math.Round(ns)).String()
	}

	if h != nil {
		// Manually format HistogramSnapshot to handle units correctly
		if h.Total == 0 {
			_, _ = fmt.Fprint(&sb, "No data collected.\n")
		} else {
			_, _ = fmt.Fprintf(
				&sb,
				"\tMean:\t%s\tStdDev:\t%s\tMin:\t%s\tMax:\t%s\tTotal:\t%d\n",
				fmtDur(h.Mean),
				fmtDur(h.StdDev),
				fmtDur(h.Min),
				fmtDur(h.Max),
				h.Total,
			)
		}
	}
	if q != nil {
		// TDigest stores raw values (nanoseconds).
		// We rely on predefined quantiles for reporting.
		qs := q.Quantiles()
		if len(qs) > 0 {
			_, _ = fmt.Fprint(&sb, "Percentiles:\n")
			for _, quantile := range qs {
				val := q.Quantile(quantile)
				_, _ = fmt.Fprintf(&sb, "\tP%g:\t%s\n", quantile*percentMultiplier, fmtDur(val))
			}
		}
	}
	if h != nil && len(h.Counts) > 0 {
		_, _ = fmt.Fprint(&sb, "Histogram:\n")
		maxCount := uint64(0)
		for _, count := range h.Counts {
			if count > maxCount {
				maxCount = count
			}
		}
		for i := 0; i < len(h.Counts); i++ {
			count := h.Counts[i]
			// Consolidate consecutive zero buckets
			if count == 0 {
				startIdx := i
				endIdx := i
				for j := i + 1; j < len(h.Counts); j++ {
					if h.Counts[j] == 0 {
						endIdx = j
					} else {
						break
					}
				}
				// Render consolidated zero bucket
				var lowerBound, upperBound string
				if startIdx == 0 {
					lowerBound = "0"
				} else if startIdx-1 < len(h.Bounds) {
					lowerBound = fmtDur(h.Bounds[startIdx-1])
				} else {
					lowerBound = "?"
				}

				if endIdx >= len(h.Bounds) {
					upperBound = "inf"
				} else {
					upperBound = fmtDur(h.Bounds[endIdx])
				}
				_, _ = fmt.Fprintf(&sb, "\t%s - %s [0]\t|\n", lowerBound, upperBound)
				i = endIdx
				continue
			}

			var bar string
			if maxCount > 0 {
				bar = strings.Repeat("∎", int(float64(count)/float64(maxCount)*barWidth))
			}
			var lowerBound, upperBound string
			if i == 0 {
				lowerBound = "0"
			} else {
				if i-1 < len(h.Bounds) {
					lowerBound = fmtDur(h.Bounds[i-1])
				} else {
					lowerBound = "?"
				}
			}
			if i >= len(h.Bounds) {
				upperBound = "inf"
			} else {
				upperBound = fmtDur(h.Bounds[i])
			}
			_, _ = fmt.Fprintf(&sb, "\t%s - %s [%d]\t|%s\n", lowerBound, upperBound, count, bar)
		}
	}
	return sb.String()
}
