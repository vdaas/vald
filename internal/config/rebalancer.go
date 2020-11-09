//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	RebalanceJobName       string `yaml:"rebalance_job_name" json:"rebalance_job_name"`
	RebalanceJobNamespace  string `yaml:"rebalance_job_namespance" json:"rebalance_job_namespance"`
	AgentName              string `yaml:"agent_name" json:"agent_name"`
	AgentNamespace         string `yaml:"agent_namespace" json:"agent_namespace"`
	ReconcileCheckDuration string `yaml:"reconcile_check_duration" json:"reconcile_check_duration"`
	JobTemplatePath        string `yaml:"job_template_path" json:"job_template_path"`
	Tolerance              int    `yaml:"tolerance" json:"tolerance"`
}

// Bind binds rebalance controller configuration.
func (r *RebalanceController) Bind() *RebalanceController {
	r.RebalanceJobName = GetActualValue(r.RebalanceJobName)
	r.RebalanceJobNamespace = GetActualValue(r.RebalanceJobNamespace)
	r.AgentName = GetActualValue(r.AgentName)
	r.AgentNamespace = GetActualValue(r.AgentNamespace)
	r.ReconcileCheckDuration = GetActualValue(r.ReconcileCheckDuration)
	r.JobTemplatePath = GetActualValue(r.JobTemplatePath)

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
	// Client represent gRPC client configuration.
	Client *GRPCClient `json:"client" yaml:"client"`
}

// Bind binds rebalance job configuration.
func (r *RebalanceJob) Bind() *RebalanceJob {
	r.BackupFilePath = GetActualValue(r.BackupFilePath)
	r.GatewayHost = GetActualValue(r.GatewayHost)

	if r.BlobStorage != nil {
		r.BlobStorage = r.BlobStorage.Bind()
	}

	if r.Client != nil {
		r.Client = r.Client.Bind()
	}

	return r
}
