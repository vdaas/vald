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

package main

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
)

func addIndexCorrectionPanels(builder *dashboard.DashboardBuilder) {
	builder.
		WithPanel(indexCorrectionPanel("Checked Index Count", "index_job_correction_checked_index_count")).
		WithPanel(indexCorrectionPanel("Corrected Index Count", "index_job_correction_corrected_old_index_count")).
		WithPanel(indexCorrectionPanel("Corrected Replication Count", "index_job_correction_corrected_replication_count"))
}

func indexCorrectionPanel(title string, metric string) *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title(title).
		WithTarget(prometheusQuery(
			fmt.Sprintf(`%s{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"}`, metric),
		).Format("table")).
		ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"mean"})).
		Span(8).Height(heightMedium)
}
