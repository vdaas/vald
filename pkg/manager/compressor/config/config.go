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

// Package setting stores all server application settings
package config

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
)

type GlobalConfig = config.GlobalConfig

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// BackupManager represent backup manager configuration
	BackupManager *config.BackupManager `json:"backup" yaml:"backup"`

	// Compressor represent compressor configuration
	Compressor *config.Compressor `json:"compressor" yaml:"compressor"`

	// Registerer represent registerer configuration
	Registerer *config.CompressorRegisterer `json:"registerer" yaml:"registerer"`
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

	if cfg.BackupManager != nil {
		cfg.BackupManager = cfg.BackupManager.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Compressor != nil {
		cfg.Compressor = cfg.Compressor.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Registerer != nil {
		cfg.Registerer = cfg.Registerer.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	return cfg, nil
}
