//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// AgentSidecar represents the configuration for the agent sidecar.
type AgentSidecar struct {
	// Mode represents sidecar mode
	Mode string `yaml:"mode" json:"mode"`

	// WatchDir represents watch target directory for backup
	WatchDir string `yaml:"watch_dir" json:"watch_dir"`

	// WatchEnabled represent auto backup triggered by file changes is enabled or not
	WatchEnabled bool `yaml:"watch_enabled" json:"watch_enabled"`

	// AutoBackupEnabled represent auto backup triggered by timer is enabled or not
	AutoBackupEnabled bool `yaml:"auto_backup_enabled" json:"auto_backup_enabled"`

	// AutoBackupDuration represent checking loop duration for auto backup execution
	AutoBackupDuration string `yaml:"auto_backup_duration" json:"auto_backup_duration"`

	// PostStopTimeout represent timeout duration for file changing during post stop
	PostStopTimeout string `yaml:"post_stop_timeout" json:"post_stop_timeout"`

	// Filename represent backup filename
	Filename string `yaml:"filename" json:"filename"`

	// FilenameSuffix represent suffix of backup filename
	FilenameSuffix string `yaml:"filename_suffix" json:"filename_suffix"`

	// BlobStorage represent blob storage configurations
	BlobStorage *Blob `yaml:"blob_storage" json:"blob_storage"`

	// Compress represent compression configurations
	Compress *CompressCore `yaml:"compress" json:"compress"`

	// RestoreBackoffEnabled represent backoff enabled or not
	RestoreBackoffEnabled bool `yaml:"restore_backoff_enabled" json:"restore_backoff_enabled"`

	// RestoreBackoff represent backoff configurations for restoring process
	RestoreBackoff *Backoff `yaml:"restore_backoff" json:"restore_backoff"`

	// Client represent HTTP client configurations
	Client *Client `yaml:"client" json:"client"`
}

// Bind binds the actual data from the AgentSidecar receiver fields.
func (s *AgentSidecar) Bind() *AgentSidecar {
	s.Mode = GetActualValue(s.Mode)
	s.WatchDir = GetActualValue(s.WatchDir)
	s.AutoBackupDuration = GetActualValue(s.AutoBackupDuration)
	s.PostStopTimeout = GetActualValue(s.PostStopTimeout)
	s.Filename = GetActualValue(s.Filename)
	s.FilenameSuffix = GetActualValue(s.FilenameSuffix)

	if s.BlobStorage != nil {
		s.BlobStorage = s.BlobStorage.Bind()
	} else {
		s.BlobStorage = new(Blob)
	}

	if s.Compress != nil {
		s.Compress = s.Compress.Bind()
	} else {
		s.Compress = new(CompressCore)
	}

	if s.RestoreBackoff != nil {
		s.RestoreBackoff = s.RestoreBackoff.Bind()
	} else {
		s.RestoreBackoff = new(Backoff)
	}

	if s.Client != nil {
		s.Client = s.Client.Bind()
	} else {
		s.Client = new(Client)
	}

	return s
}
