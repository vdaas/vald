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

// AgentSidecar represents the configuration for the agent sidecar.
type AgentSidecar struct {
	// BlobStorage represents the blob storage configuration.
	BlobStorage *Blob `json:"blob_storage" yaml:"blob_storage"`
	// Client represents the client configuration.
	Client *Client `json:"client" yaml:"client"`
	// RestoreBackoff represents the restore backoff configuration.
	RestoreBackoff *Backoff `json:"restore_backoff" yaml:"restore_backoff"`
	// Compress represents the compress configuration.
	Compress *CompressCore `json:"compress" yaml:"compress"`
	// Filename represents the filename.
	Filename string `json:"filename" yaml:"filename"`
	// PostStopTimeout represents the post stop timeout duration.
	PostStopTimeout string `json:"post_stop_timeout" yaml:"post_stop_timeout"`
	// Mode represents the mode.
	Mode string `json:"mode" yaml:"mode"`
	// FilenameSuffix represents the filename suffix.
	FilenameSuffix string `json:"filename_suffix" yaml:"filename_suffix"`
	// AutoBackupDuration represents the auto backup duration.
	AutoBackupDuration string `json:"auto_backup_duration" yaml:"auto_backup_duration"`
	// WatchDir represents the watch directory.
	WatchDir string `json:"watch_dir" yaml:"watch_dir"`
	// AutoBackupEnabled enables auto backup.
	AutoBackupEnabled bool `json:"auto_backup_enabled" yaml:"auto_backup_enabled"`
	// RestoreBackoffEnabled enables restore backoff.
	RestoreBackoffEnabled bool `json:"restore_backoff_enabled" yaml:"restore_backoff_enabled"`
	// WatchEnabled enables watch.
	WatchEnabled bool `json:"watch_enabled" yaml:"watch_enabled"`
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
