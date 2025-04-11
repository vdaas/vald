// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func addVariables(builder *dashboard.DashboardBuilder, server_name string) {
	builder.WithVariable(dashboard.NewQueryVariableBuilder("Namespace").
		Label("namespace").
		Options([]dashboard.VariableOption{}).
		Current(dashboard.VariableOption{
			Selected: cog.ToPtr(false),
			Text:     dashboard.StringOrArrayOfString{String: cog.ToPtr("default")},
			Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr("default")},
		}).
		Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(kube_pod_info, namespace)")}),
	).
		WithVariable(dashboard.NewQueryVariableBuilder("ReplicaSet").
			Label("name").
			Query(dashboard.StringOrMap{String: cog.ToPtr(fmt.Sprintf("label_values(app_version_info{server_name=~\"%s\"}, kubernetes_name)", server_name))}),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("PodName").
			Label("pod").
			Current(dashboard.VariableOption{
				Selected: cog.ToPtr(false),
				Text:     dashboard.StringOrArrayOfString{String: cog.ToPtr("All")},
				Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr("$__all")},
			}).
			AllValue(".+").
			Multi(false).
			IncludeAll(true).
			Query(dashboard.StringOrMap{String: cog.ToPtr(fmt.Sprintf("label_values(app_version_info{server_name=~\"%s\", kubernetes_name=\"$ReplicaSet\"}, target_pod)", server_name))}),
		).
		WithVariable(dashboard.NewIntervalVariableBuilder("interval").
			Label("interval").
			Current(dashboard.VariableOption{
				Text:  dashboard.StringOrArrayOfString{String: cog.ToPtr("1m")},
				Value: dashboard.StringOrArrayOfString{String: cog.ToPtr("1m")},
			}).
			Values(dashboard.StringOrMap{String: cog.ToPtr("1m,2m,5m,10m,30m,1h,6h,12h,1d")}),
		)
}
