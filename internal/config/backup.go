//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// BackupManager represents the configuration for backup manager.
type BackupManager struct {
	Client *GRPCClient `json:"client" yaml:"client"`
}

// Bind binds the actual data from the BackupManager receiver fields.
func (b *BackupManager) Bind() *BackupManager {
	if b.Client != nil {
		b.Client = b.Client.Bind()
	} else {
		b.Client = newGRPCClientConfig()
	}
	return b
}
