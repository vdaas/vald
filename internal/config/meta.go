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

import "fmt"

type Meta struct {
	Host   string      `json:"host" yaml:"host"`
	Port   int         `json:"port" yaml:"port"`
	Client *GRPCClient `json:"client" yaml:"client"`
}

func (m *Meta) Bind() *Meta {
	m.Host = GetActualValue(m.Host)

	if m.Client != nil {
		m.Client.Bind()
	} else {
		m.Client = newGRPCClientConfig()
	}
	if len(m.Host) != 0 {
		m.Client.Addrs = append(m.Client.Addrs, fmt.Sprintf("%s:%d", m.Host, m.Port))
	}
	return m
}
