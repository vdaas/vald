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

// Package config providers configuration type and load configuration logic
package config

// Logging represents Logging configuration.
type Logging struct {
	Logger string `json:"logger" yaml:"logger"`
	Level  string `json:"level"  yaml:"level"`
	Format string `json:"format" yaml:"format"`
}

// Bind returns Logging object whose every value is field value or envirionment value.
func (l *Logging) Bind() *Logging {
	l.Logger = GetActualValue(l.Logger)
	l.Level = GetActualValue(l.Level)
	l.Format = GetActualValue(l.Format)
	return l
}
