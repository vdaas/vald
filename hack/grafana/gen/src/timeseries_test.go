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
package main

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/observability/metrics/mem"
	"github.com/vdaas/vald/internal/strings"
)

func TestAddMemstatsPanelsForcedGCQuery(t *testing.T) {
	t.Parallel()

	model, err := func() (dashboard.Dashboard, error) {
		builder := dashboard.NewDashboardBuilder("test")
		addMemstatsPanels(builder)
		return builder.Build()
	}()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Panels []struct {
			Title   string `json:"title"`
			Targets []struct {
				Expr string `json:"expr"`
			} `json:"targets"`
		} `json:"panels"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		title  string
		metric string
	}{
		{title: mem.NumGCMetricsDescription, metric: mem.NumGCMetricsName},
		{title: mem.NumForcedGCMetricsDescription, metric: mem.NumForcedGCMetricsName},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.title, func(tt *testing.T) {
			tt.Parallel()
			found := false
			for _, p := range got.Panels {
				if p.Title != test.title {
					continue
				}
				found = true
				if len(p.Targets) != 1 || !strings.HasPrefix(p.Targets[0].Expr, test.metric+"{") {
					tt.Errorf("panel %q expected a query of %s, got %+v", p.Title, test.metric, p.Targets)
				}
			}
			if !found {
				tt.Errorf("panel %q is not found", test.title)
			}
		})
	}
}
