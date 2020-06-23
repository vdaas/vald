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

// Operation is type of implemented load test.
type Operation uint8

// Operation method.
const (
	Unknown Operation = iota
	Insert
	Search
)

// OperationMethod convert string to Operation.
func OperationMethod(s string) Operation {
	switch strings.ToLower(s) {
	case "insert":
		return Insert
	case "search":
		return Search
	default:
		return Unknown
	}
}

// String convert Operation to string.
func (o Operation) String() string {
	switch o {
	case Insert:
		return "insert"
	case Search:
		return "search"
	default:
		return "unknown operation"
	}
}

// Data represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`
	Addr                string             `json:"addr" yaml:"addr"`
	Method              string             `json:"method" yaml:"method"`
	Dataset             string             `json:"dataset" yaml:"dataset"`
	Concurrency         int                `json:"concurrency" yaml:"concurrency"`
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
	cfg.Method = config.GetActualValue(cfg.Method)
	cfg.Dataset = config.GetActualValue(cfg.Dataset)
	cfg.ProgressDuration = config.GetActualValue(cfg.ProgressDuration)

	return cfg, nil
}
