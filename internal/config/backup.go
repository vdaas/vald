//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

type BackupManager struct {
	Host   string      `json:"host" yaml:"host"`
	Port   int         `json:"port" yaml:"port"`
	Addr   string      `json:"addr" yaml:"addr"`
	Client *GRPCClient `json:"client" yaml:"client"`
}

func (b *BackupManager) Bind() *BackupManager {
	if len(b.Addr) == 0 {
		b.Addr = fmt.Sprintf("%s, %d", b.Host, b.Port)
	}
	if b.Client != nil {
		b.Client = b.Client.Bind()
	}
	return b
}
