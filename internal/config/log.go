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

import "strings"

type logFormat uint8

const (
	JSON logFormat = iota
	YAML
)

func (lf logFormat) String() string {
	switch lf {
	case JSON:
		return "json"
	case YAML:
		return "yaml"
	default:
		return "unknown"
	}
}

type logMode uint8

const (
	GLG logMode = iota
)

func (lm logMode) Mode(mode string) logMode {
	switch strings.ToLower(mode) {
	case "glg":
		return GLG
	default:
		return GLG
	}
}

type Log struct {
	Mode   string `json:"mode" yaml:"mode"`
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
}

func (l *Log) Bind() *Log {
	l.Level = GetActualValue(l.Level)
	l.Mode = GetActualValue(l.Mode)
	l.Format = GetActualValue(l.Format)
	return l
}
