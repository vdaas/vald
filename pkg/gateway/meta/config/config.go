//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package config

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
)

type (
	GlobalConfig = config.GlobalConfig
	Server       = config.Server
)

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Client represent gateway client configuration to lb gateway
	Client *config.GRPCClient `json:"client" yaml:"client"`

	// MetadataStore represent metadata store (e.g. TiKV PD) configuration
	MetadataStore *MetadataStoreConfig `json:"metadata_store" yaml:"metadata_store"`
}

// MetadataStoreConfig represents the configuration for the metadata store.
type MetadataStoreConfig struct {
	// Addrs represent metadata store addresses (e.g. TiKV PD addresses)
	Addrs []string `json:"addrs" yaml:"addrs"`
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

	if cfg.Client != nil {
		cfg.Client = cfg.Client.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.MetadataStore == nil || len(cfg.MetadataStore.Addrs) == 0 {
		return nil, errors.ErrInvalidConfig
	}

	return cfg, nil
}
