// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package config

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
)

// GlobalConfig is a type alias of config.GlobalConfig representing application base configurations.
type GlobalConfig = config.GlobalConfig

// Data represents the application configurations.
type Data struct {
	// GlobalConfig represents application base configurations.
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represents observability configurations.
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Deletion represents auto indexing service configurations.
	Deletion *config.IndexDeletor `json:"deletor" yaml:"deletor"`
}

// NewConfig load configurations from file path.
func NewConfig(path string) (cfg *Data, err error) {
	cfg = new(Data)

	if err = config.Read(path, &cfg); err != nil {
		return nil, err
	}

	if cfg != nil {
		_ = cfg.GlobalConfig.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Server != nil {
		_ = cfg.Server.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Observability != nil {
		_ = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(config.Observability).Bind()
	}

	if cfg.Deletion != nil {
		_ = cfg.Deletion.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}
	return cfg, nil
}
