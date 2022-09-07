//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
)

type Mode uint8

const (
	SIDECAR Mode = 1 + iota
	INITCONTAINER
)

func (m Mode) String() string {
	switch m {
	case SIDECAR:
		return "sidecar"
	case INITCONTAINER:
		return "initcontainer"
	}
	return "unknown"
}

func SidecarMode(m string) Mode {
	switch strings.ToLower(m) {
	case "sidecar":
		return SIDECAR
	case "initcontainer":
		return INITCONTAINER
	}
	return 0
}

type GlobalConfig = config.GlobalConfig

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Sidecar represent agent storage sync sidecar service configuration
	AgentSidecar *config.AgentSidecar `json:"agent_sidecar" yaml:"agent_sidecar"`
}

func NewConfig(path string) (cfg *Data, err error) {
	cfg = new(Data)

	err = config.Read(path, &cfg)

	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfg.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Observability != nil {
		cfg.Observability = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(config.Observability).Bind()
	}

	if cfg.AgentSidecar != nil {
		cfg.AgentSidecar = cfg.AgentSidecar.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	return cfg, nil
}
