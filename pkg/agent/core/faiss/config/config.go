//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package setting stores all server application settings
package config

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
)

// GlobalConfig is type alias for config.GlobalConfig.
type GlobalConfig = config.GlobalConfig

// Data represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Faiss represent faiss core configuration
	Faiss *config.Faiss `json:"faiss" yaml:"faiss"`
}

// NewConfig returns the Data struct or error from the given file path.
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

	if cfg.Faiss != nil {
		cfg.Faiss = cfg.Faiss.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	return cfg, nil
}
