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
	"github.com/vdaas/vald/internal/config"
)

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.Default `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// BackupManager represent backup manager configuration
	BackupManager *config.BackupManager `json:"backup" yaml:"backup"`

	// Compressor represent compressor configuration
	Compressor *config.Compressor `json:"compressor" yaml:"compressor"`
}

func NewConfig(path string) (cfg *Data, err error) {
	err = config.Read(path, &cfg)

	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfg.Bind()
	}

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	}

	if cfg.BackupManager != nil {
		cfg.BackupManager = cfg.BackupManager.Bind()
	}

	if cfg.Compressor != nil {
		cfg.Compressor = cfg.Compressor.Bind()
	}

	return cfg, nil
}
