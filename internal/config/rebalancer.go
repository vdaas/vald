//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package config providers configuration type and load configuration logic
package config

// RebalanceController represent rebalance controller configuration.
type RebalanceController struct {
	PodName                string  `yaml:"pod_name" json:"pod_name"`
	PodNamespace           string  `yaml:"pod_namespace" json:"pod_namespace"`
	RebalanceJobName       string  `yaml:"rebalance_job_name" json:"rebalance_job_name"`
	RebalanceJobNamespace  string  `yaml:"rebalance_job_namespace" json:"rebalance_job_namespace"`
	RebalanceJobTemplate   string  `yaml:"rebalance_job_template" json:"rebalance_job_template"`
	AgentName              string  `yaml:"agent_name" json:"agent_name"`
	AgentNamespace         string  `yaml:"agent_namespace" json:"agent_namespace"`
	AgentResourceType      string  `yaml:"agent_resource_type" json:"agent_resource_type"`
	ReconcileCheckDuration string  `yaml:"reconcile_check_duration" json:"reconcile_check_duration"`
	Tolerance              float64 `yaml:"tolerance" json:"tolerance"`
	RateThreshold          float64 `yaml:"rate_threshold" json:"rate_threshold"`
	LeaderElectionID       string  `yaml:"leader_election_id" json:"leader_election_id"`
}

// Bind binds rebalance controller configuration.
func (r *RebalanceController) Bind() *RebalanceController {
	r.PodName = GetActualValue(r.PodName)
	r.PodNamespace = GetActualValue(r.PodNamespace)
	r.RebalanceJobName = GetActualValue(r.RebalanceJobName)
	r.RebalanceJobNamespace = GetActualValue(r.RebalanceJobNamespace)
	r.RebalanceJobTemplate = GetActualValue(r.RebalanceJobTemplate)
	r.AgentName = GetActualValue(r.AgentName)
	r.AgentNamespace = GetActualValue(r.AgentNamespace)
	r.AgentResourceType = GetActualValue(r.AgentResourceType)
	r.ReconcileCheckDuration = GetActualValue(r.ReconcileCheckDuration)
	r.LeaderElectionID = GetActualValue(r.LeaderElectionID)

	return r
}

// RebalanceJob represent rebalance job configuration.
type RebalanceJob struct {
	// BlobStorage represent blob storage configurations.
	BlobStorage *Blob `yaml:"blob_storage" json:"blob_storage"`
	// BackupFilePath represent kvsdb file path.
	BackupFilePath string `yaml:"backup_file_path" json:"backup_file_path"`
	// RebalanceRate represent rate to rebalance data.
	RebalanceRate int `yaml:"rebalance_rate" json:"rebalance_rate"`
	// GatewayHost represent gateway host name.
	GatewayHost string `json:"gateway_host" yaml:"gateway_host"`
	// GatewayPort represent gateway port.
	GatewayPort int `json:"gateway_port" yaml:"gateway_port"`
	// GatewayClient represent gRPC client configuration.
	GatewayClient *GRPCClient `json:"gateway_client" yaml:"gateway_client"`
}

// Bind binds rebalance job configuration.
func (r *RebalanceJob) Bind() *RebalanceJob {
	r.BackupFilePath = GetActualValue(r.BackupFilePath)
	r.GatewayHost = GetActualValue(r.GatewayHost)

	if r.BlobStorage != nil {
		r.BlobStorage = r.BlobStorage.Bind()
	}

	if r.GatewayClient != nil {
		r.GatewayClient = r.GatewayClient.Bind()
	}

	return r
}
