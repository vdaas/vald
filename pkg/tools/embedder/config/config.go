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
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
)

type GlobalConfig = iconf.GlobalConfig

type OpenAI struct {
	Token string `json:"token" yaml:"token"`
	Model string `json:"model" yaml:"model"`
}

func (o *OpenAI) Bind() *OpenAI {
	o.Token = iconf.GetActualValue(o.Token)
	o.Model = iconf.GetActualValue(o.Model)
	return o
}

type LLM struct {
	Provider string  `json:"provider" yaml:"provider"`
	OpenAI   *OpenAI `json:"openai"   yaml:"openai"`
}

func (l *LLM) Bind() *LLM {
	l.Provider = iconf.GetActualValue(l.Provider)
	if l.OpenAI != nil {
		l.OpenAI = l.OpenAI.Bind()
	}
	return l
}

type Data struct {
	GlobalConfig  `json:",inline" yaml:",inline"`
	Server        *iconf.Servers       `json:"server_config" yaml:"server_config"`
	Observability *iconf.Observability `json:"observability" yaml:"observability"`
	Client        *iconf.GRPCClient    `json:"client"        yaml:"client"`
	Meta          *iconf.Meta          `json:"meta"          yaml:"meta"`
	LLM           *LLM                 `json:"llm"           yaml:"llm"`
}

func NewConfig(path string) (cfg *Data, err error) {
	cfg = new(Data)
	if err = iconf.Read(path, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.ErrInvalidConfig
	}
	cfg.Bind()
	if cfg.Server == nil || cfg.Client == nil || cfg.LLM == nil {
		return nil, errors.ErrInvalidConfig
	}
	cfg.Server = cfg.Server.Bind()
	cfg.Client = cfg.Client.Bind()
	if cfg.Observability != nil {
		cfg.Observability = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(iconf.Observability).Bind()
	}
	if cfg.Meta != nil {
		cfg.Meta = cfg.Meta.Bind()
	}
	cfg.LLM = cfg.LLM.Bind()
	if cfg.LLM.Provider == "" {
		cfg.LLM.Provider = "openai"
	}
	if cfg.LLM.Provider != "openai" {
		return nil, errors.New("unsupported llm provider: " + cfg.LLM.Provider)
	}
	if cfg.LLM.OpenAI == nil || cfg.LLM.OpenAI.Token == "" {
		return nil, errors.New("llm.openai.token is required")
	}
	return cfg, nil
}
