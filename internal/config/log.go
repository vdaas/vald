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

// Package config providers configuration type and load configuration logic
package config

import (
	kglg "github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/glg"
)

type Log struct {
	Type   string `json:"type" yaml:"type"`
	Level  string `json:"level" yaml:"level"`
	Mode   string `json:"mode" yaml:"mode"`
	Format string `json:"format" yaml:"format"`
}

func (l Log) Bind() Log {
	l.Type = GetActualValue(l.Type)
	l.Level = GetActualValue(l.Level)
	l.Mode = GetActualValue(l.Mode)
	l.Format = GetActualValue(l.Format)
	return l
}

func (l Log) Opts() (opts []log.Option) {
	switch l.Type {
	case "zap":
		// TODO(@funapy)
		fallthrough

	case "glg":
		fallthrough

	default:
		gopts := []glg.Option{
			glg.WithLevel(l.Level),
			glg.WithMode(l.Mode),
		}

		if l.Format == "json" {
			gopts = append(gopts, glg.WithEnableJSON())
		}

		opts = []log.Option{
			log.WithLogger(
				glg.New(kglg.Get(), gopts...),
			),
		}
	}
	return
}
