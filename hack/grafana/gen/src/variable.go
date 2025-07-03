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

func addNamespaceVariable(builder *dashboard.DashboardBuilder) {
	builder.WithVariable(dashboard.NewQueryVariableBuilder("Namespace").
		Label("namespace").
		Options([]dashboard.VariableOption{}).
		Current(dashboard.VariableOption{
			Selected: cog.ToPtr(false),
			Text:     dashboard.StringOrArrayOfString{String: cog.ToPtr("default")},
			Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr("default")},
		}).
		Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(kube_pod_info, namespace)")}),
	)
}

func addIntervalVariable(builder *dashboard.DashboardBuilder, current string) {
	builder.WithVariable(dashboard.NewIntervalVariableBuilder("interval").
		Label("interval").
		Current(dashboard.VariableOption{
			Text:  dashboard.StringOrArrayOfString{String: cog.ToPtr(current)},
			Value: dashboard.StringOrArrayOfString{String: cog.ToPtr(current)},
		}).
		Values(dashboard.StringOrMap{String: cog.ToPtr("1m,2m,5m,10m,30m,1h,6h,12h,1d")}),
	)
}

func addOverviewVariables(builder *dashboard.DashboardBuilder) {
	addNamespaceVariable(builder)
	builder.
		WithVariable(dashboard.NewQueryVariableBuilder("ReplicaSet").
			Label("name").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(app_version_info, kubernetes_name)")}).
			IncludeAll(true).
			AllValue("vald-.*").
			Multi(true),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("PodName").
			Label("agent-pod").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(agent_core_ngt_is_indexing{kubernetes_name =~\"vald-agent\"}, target_pod)")}).
			IncludeAll(true).
			Multi(true),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("ContainerName").
			Label("agent-container").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(agent_core_ngt_is_indexing{kubernetes_name =~\"vald-agent\", target_pod=~\"$PodName\"}, container)")}).
			IncludeAll(true).
			AllValue(".*").
			Multi(true),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("ValdGatewayPodName").
			Label("gateway-pod").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(app_version_info{kubernetes_name=~\"vald-lb-gateway\"}, target_pod)")}).
			IncludeAll(true).
			Multi(true),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("ValdGatewayContainerName").
			Label("gateway-container").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(app_version_info{kubernetes_name =~\"vald-lb-gateway\", target_pod=~\"$ValdGatewayPodName\"}, container)")}).
			IncludeAll(true).
			AllValue(".*").
			Multi(true),
		)
	addIntervalVariable(builder, "30m")
}

func addBasicVariables(builder *dashboard.DashboardBuilder, serverName string) {
	builder.
		WithVariable(dashboard.NewQueryVariableBuilder("ReplicaSet").
			Label("name").
			Query(dashboard.StringOrMap{String: cog.ToPtr(fmt.Sprintf("label_values(%s{server_name=~\"%s\"}, kubernetes_name)", appInfo, serverName))}),
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
			Query(dashboard.StringOrMap{String: cog.ToPtr(fmt.Sprintf("label_values(%s{server_name=~\"%s\", kubernetes_name=\"$ReplicaSet\"}, target_pod)", appInfo, serverName))}),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("ContainerName").
			Label("container").
			Query(dashboard.StringOrMap{String: cog.ToPtr(fmt.Sprintf("label_values(%s{server_name=~\"%s\", kubernetes_name=\"$ReplicaSet\", target_pod=~\"$PodName\"}, container)", appInfo, serverName))}).
			IncludeAll(true).
			AllValue(".*").
			Multi(true),
		)
}

func addVariables(builder *dashboard.DashboardBuilder, serverName string, defaultInterval string) {
	addNamespaceVariable(builder)
	addBasicVariables(builder, serverName)
	addIntervalVariable(builder, defaultInterval)
}

func addHelmOperatorVariables(builder *dashboard.DashboardBuilder) {
	addNamespaceVariable(builder)
	builder.
		WithVariable(dashboard.NewQueryVariableBuilder("ReplicaSet").
			Label("name").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(helm_operator_build_info{namespace=~\"$Namespace\"}, app_kubernetes_io_name)")}),
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
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(helm_operator_build_info{namespace=~\"$Namespace\", app_kubernetes_io_name=~\"$ReplicaSet\"}, instance)")}),
		)
	addIntervalVariable(builder, "5m")
}

func addBenchmarkOperatorVariables(builder *dashboard.DashboardBuilder) {
	addNamespaceVariable(builder)
	addBasicVariables(builder, "benchmark operator.*")
	builder.
		WithVariable(dashboard.NewQueryVariableBuilder("JobReplicaSet").
			Label("job_name").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(app_version_info{server_name=~\"benchmark job.*\"},kubernetes_name)")}),
		).
		WithVariable(dashboard.NewQueryVariableBuilder("JobPodName").
			Label("job_pod").
			Query(dashboard.StringOrMap{String: cog.ToPtr("label_values(app_version_info{server_name=~\"benchmark job.*\", kubernetes_name=~\"$JobReplicaSet\"},target_pod)")}).
			IncludeAll(true).
			AllValue(".*"),
		)
	addIntervalVariable(builder, "5m")
}
