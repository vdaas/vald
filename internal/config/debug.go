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

// Package config providers configuration type and load configuration logic
package config

type Debug struct {
	// Profile represent profiling the server
	Profile struct {
		Enable bool    `yaml:"enable" json:"enable"`
		Server *Server `yaml:"server" json:"server"`
	} `yaml:"profile" json:"profile"`

	// Log represent the server enable debug log or not.
	Log struct {
		Level string `yaml:"level" json:"level"`
		Mode  string `yaml:"mode" json:"mode"`
	} `yaml:"log" json:"log"`
}

func (d *Debug) Bind() *Debug {
	if d.Profile.Server != nil {
		d.Profile.Server.Bind()
	}
	d.Log.Level = GetActualValue(d.Log.Level)
	d.Log.Mode = GetActualValue(d.Log.Mode)
	return d
}
