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

package metrics

import (
	"encoding/csv"
	"maps"
	"math"
	"slices"
	"strconv"
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
	p.renderErrorDetails(&sb)

	return sb.String()
}

func (p *SnapshotPresenter) renderSummary(sb *strings.Builder) {
	s := p.snapshot
	total := s.Total
	errs := s.Errors

	totalDuration := s.LastUpdated.Sub(s.StartTime)

	sb.WriteString("\n--- Summary ---\n")
	sb.WriteString("Total Requests:\t")
	sb.WriteString(strconv.FormatUint(total, 10))
	sb.WriteByte('\n')

	sb.WriteString("Total Duration:\t")
	sb.WriteString(totalDuration.String())
	sb.WriteByte('\n')

	if totalDuration.Seconds() > 0 {
		sb.WriteString("Requests/sec:\t")
		sb.WriteString(strconv.FormatFloat(float64(total)/totalDuration.Seconds(), 'f', 2, 64))
		sb.WriteByte('\n')
	}

	sb.WriteString("Errors:\t")
	sb.WriteString(strconv.FormatUint(errs, 10))
	sb.WriteString(" (")
	sb.WriteString(strconv.FormatFloat(float64(errs)/float64(total)*percentMultiplier, 'f', 2, 64))
	sb.WriteString("%)\n")
}

func (p *SnapshotPresenter) renderLatency(sb *strings.Builder) {
	p.renderHistogram(sb, "Latency", p.snapshot.Latencies, p.snapshot.LatPercentiles)
}

func (p *SnapshotPresenter) renderQueueWait(sb *strings.Builder) {
	p.renderHistogram(sb, "Queue Wait", p.snapshot.QueueWaits, p.snapshot.QWPercentiles)
}

func (p *SnapshotPresenter) renderStatusCodes(sb *strings.Builder) {
	s := p.snapshot
	total := s.Total
	sb.WriteString("\n--- Status Codes ---\n")

	if s.Codes != nil {
		codes := make([]codes.Code, 0, len(s.Codes))
		for code := range s.Codes {
			codes = append(codes, code)
		}
		slices.Sort(codes)
		for _, code := range codes {
			count := s.Codes[code]
			sb.WriteString("\t- ")
			sb.WriteString(code.String())
			sb.WriteString(":\t")
			sb.WriteString(strconv.FormatUint(count, 10))
			sb.WriteString(" (")
			sb.WriteString(strconv.FormatFloat(float64(count)/float64(total)*percentMultiplier, 'f', 2, 64))
			sb.WriteString("%)\n")
		}
	}
}

func (p *SnapshotPresenter) renderErrorDetails(sb *strings.Builder) {
	s := p.snapshot
	if len(s.ErrorDetails) == 0 {
		return
	}
	sb.WriteString("\n--- Error Details ---\n")
	// Sort for stable output
	for _, msg := range slices.Sorted(maps.Keys(s.ErrorDetails)) {
		sb.WriteString("\t- " + msg + ":\t" + strconv.FormatUint(s.ErrorDetails[msg], 10))
		sb.WriteByte('\n')
	}
}

func (p *SnapshotPresenter) renderExemplars(sb *strings.Builder) {
	s := p.snapshot
	if s.ExemplarDetails != nil {
		renderExemplars := func(title string, items []*ExemplarItem) {
			if len(items) > 0 {
				sb.WriteString("\n--- Exemplars (")
				sb.WriteString(title)
				sb.WriteString(") ---\n")

				for _, ex := range items {
					status := ""
					if ex.Err != nil {
						status = " (Failed)"
					}
					sb.WriteString("\t- RequestID:\t")
					sb.WriteString(ex.RequestID)
					sb.WriteString(",\tLatency:\t")
					sb.WriteString(ex.Latency.String())
					sb.WriteString(status)
					sb.WriteByte('\n')

					if ex.Msg != "" {
						sb.WriteString("\t  Message:\t")
						sb.WriteString(ex.Msg)
						sb.WriteByte('\n')
					}
				}
			}
		}
		renderExemplars("Slowest", s.ExemplarDetails.Slowest)
		renderExemplars("Fastest", s.ExemplarDetails.Fastest)
		renderExemplars("Average (Sampled)", s.ExemplarDetails.Average)
		renderExemplars("Failures", s.ExemplarDetails.Failures)
	} else if len(s.Exemplars) > 0 {
		sb.WriteString("\n--- Exemplars (Top ")
		sb.WriteString(strconv.Itoa(len(s.Exemplars)))
		sb.WriteString(" slowest requests) ---\n")

		for _, ex := range s.Exemplars {
			sb.WriteString("\t- RequestID:\t")
			sb.WriteString(ex.RequestID)
			sb.WriteString(",\tLatency:\t")
			sb.WriteString(ex.Latency.String())
			sb.WriteByte('\n')

			if ex.Msg != "" {
				sb.WriteString("\t  Message:\t")
				sb.WriteString(ex.Msg)
				sb.WriteByte('\n')
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
			headers = append(headers, "LatencyP"+strconv.Itoa(int(q*percentMultiplier)))
		}
	}
	headers = append(headers, "QueueWaitMin", "QueueWaitMean", "QueueWaitMax")
	if s.QWPercentiles != nil {
		for _, q := range s.QWPercentiles.Quantiles() {
			headers = append(headers, "QueueWaitP"+strconv.Itoa(int(q*percentMultiplier)))
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
		strconv.FormatUint(s.Total, 10),
		strconv.FormatUint(s.Errors, 10),
		strconv.FormatFloat(totalDuration, 'f', 4, 64),
		strconv.FormatFloat(rps, 'f', 2, 64),
		strconv.FormatFloat(errorRate, 'f', 4, 64),
		strconv.FormatFloat(latMin, 'f', 4, 64),
		strconv.FormatFloat(latMean, 'f', 4, 64),
		strconv.FormatFloat(latMax, 'f', 4, 64),
	}
	if s.LatPercentiles != nil {
		for _, q := range s.LatPercentiles.Quantiles() {
			row = append(row, strconv.FormatFloat(s.LatPercentiles.Quantile(q)/nanoToSec, 'f', 4, 64))
		}
	}

	qwMin, qwMean, qwMax := 0.0, 0.0, 0.0
	if s.QueueWaits != nil {
		qwMin = float64(s.QueueWaits.Min) / nanoToSec
		qwMean = float64(s.QueueWaits.Mean) / nanoToSec
		qwMax = float64(s.QueueWaits.Max) / nanoToSec
	}

	row = append(row,
		strconv.FormatFloat(qwMin, 'f', 4, 64),
		strconv.FormatFloat(qwMean, 'f', 4, 64),
		strconv.FormatFloat(qwMax, 'f', 4, 64),
	)
	if s.QWPercentiles != nil {
		for _, q := range s.QWPercentiles.Quantiles() {
			row = append(row, strconv.FormatFloat(s.QWPercentiles.Quantile(q)/nanoToSec, 'f', 4, 64))
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
func (p *SnapshotPresenter) renderHistogram(
	sb *strings.Builder, title string, h *HistogramSnapshot, q TDigest,
) {
	sb.WriteString("\n--- ")
	sb.WriteString(title)
	sb.WriteString(" ---\n")

	// Helper to convert nanoseconds to duration string
	fmtDur := func(ns float64) string {
		return time.Duration(math.Round(ns)).String()
	}

	if h != nil {
		// Manually format HistogramSnapshot to handle units correctly
		if h.Total == 0 {
			sb.WriteString("No data collected.\n")
		} else {
			sb.WriteString("\tMean:\t")
			sb.WriteString(fmtDur(h.Mean))
			sb.WriteString("\tStdDev:\t")
			sb.WriteString(fmtDur(h.StdDev))
			sb.WriteString("\tMin:\t")
			sb.WriteString(fmtDur(h.Min))
			sb.WriteString("\tMax:\t")
			sb.WriteString(fmtDur(h.Max))
			sb.WriteString("\tTotal:\t")
			sb.WriteString(strconv.FormatUint(h.Total, 10))
			sb.WriteByte('\n')
		}
	}
	if q != nil {
		// TDigest stores raw values (nanoseconds).
		// We rely on predefined quantiles for reporting.
		qs := q.Quantiles()
		if len(qs) > 0 {
			sb.WriteString("Percentiles:\n")
			for _, quantile := range qs {
				val := q.Quantile(quantile)
				sb.WriteString("\tP")
				sb.WriteString(strconv.FormatFloat(quantile*percentMultiplier, 'g', -1, 64))
				sb.WriteString(":\t")
				sb.WriteString(fmtDur(val))
				sb.WriteByte('\n')
			}
		}
	}
	if h != nil && len(h.Counts) > 0 {
		sb.WriteString("Histogram:\n")
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
				sb.WriteString("\t")
				sb.WriteString(lowerBound)
				sb.WriteString(" - ")
				sb.WriteString(upperBound)
				sb.WriteString(" [0]\t|\n")

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

			sb.WriteString("\t")
			sb.WriteString(lowerBound)
			sb.WriteString(" - ")
			sb.WriteString(upperBound)
			sb.WriteString(" [")
			sb.WriteString(strconv.FormatUint(count, 10))
			sb.WriteString("]\t|")
			sb.WriteString(bar)
			sb.WriteByte('\n')
		}
	}
}
