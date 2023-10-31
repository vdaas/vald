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

	// Creator represents auto indexing service configurations.
	Creator *config.IndexCreator `json:"creator" yaml:"creator"`
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

	if cfg.Creator != nil {
		_ = cfg.Creator.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}
	return cfg, nil
}
