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

// Package setting stores all server application settings
package config

import (
	"strings"

	"github.com/vdaas/vald/internal/config"
)

// GlobalConfig is type alias of config.GlobalConfig.
type GlobalConfig = config.GlobalConfig

// Operation is operation type of implemented load test.
type Operation uint8

// Operation method definition.
const (
	UnknownOperation Operation = iota
	Insert
	StreamInsert
	Search
	StreamSearch
)

// OperationMethod converts string to Operation.
func OperationMethod(s string) Operation {
	switch strings.ToLower(s) {
	case "insert":
		return Insert
	case "streaminsert":
		return StreamInsert
	case "search":
		return Search
	case "streamsearch":
		return StreamSearch
	default:
		return UnknownOperation
	}
}

// String converts Operation to string.
func (o Operation) String() string {
	switch o {
	case Insert:
		return "Insert"
	case StreamInsert:
		return "StreamInsert"
	case Search:
		return "Search"
	case StreamSearch:
		return "StreamSearch"
	default:
		return "Unknown operation"
	}
}

// Service is service type of implemented load test.
type Service uint8

// Service definitions.
const (
	UnknownService Service = iota
	Agent
	Gateway
)

// ServiceMethod converts string to Service.
func ServiceMethod(s string) Service {
	switch strings.ToLower(s) {
	case "agent":
		return Agent
	case "gateway":
		return Gateway
	default:
		return UnknownService
	}
}

// String converts Service to string.
func (s Service) String() string {
	switch s {
	case Agent:
		return "Agent"
	case Gateway:
		return "Gateway"
	default:
		return "Unknown service"
	}
}

// Data represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`
	Addr                string             `json:"addr" yaml:"addr"`
	Service             string             `json:"service" yaml:"service"`
	Operation           string             `json:"operation" yaml:"operation"`
	Dataset             string             `json:"dataset" yaml:"dataset"`
	Concurrency         int                `json:"concurrency" yaml:"concurrency"`
	BatchSize           int                `json:"batch_size" yaml:"batch_size"`
	ProgressDuration    string             `json:"progress_duration" yaml:"progress_duration"`
	Client              *config.GRPCClient `json:"client" yaml:"client"`
}

// NewConfig returns loaded configuration from file.
func NewConfig(path string) (cfg *Data, err error) {
	err = config.Read(path, &cfg)

	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfg.Bind()
	}
	if cfg.Client != nil {
		cfg.Client.Bind()
	}

	cfg.Addr = config.GetActualValue(cfg.Addr)
	cfg.Operation = config.GetActualValue(cfg.Operation)
	cfg.Dataset = config.GetActualValue(cfg.Dataset)
	cfg.ProgressDuration = config.GetActualValue(cfg.ProgressDuration)
	cfg.Service = config.GetActualValue(cfg.Service)

	return cfg, nil
}
